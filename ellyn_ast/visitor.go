package ellyn_ast

import (
	"bytes"
	"fmt"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"go/ast"
	"go/token"
	"sort"
)

type insert struct {
	offset        int
	content       []byte
	contentGetter func() string
	priority      int
}

type FileVisitor struct {
	fileId       uint32
	relativePath string
	fset         *token.FileSet
	content      []byte
	prog         *Program
	inserts      []*insert
}

func (f *FileVisitor) WriteTo(absPath string) {
	newContent := f.mergeInserts()
	utils.OS.WriteTo(absPath, newContent)
}

func (f *FileVisitor) mergeInserts() []byte {
	sort.Slice(f.inserts, func(i, j int) bool {
		if f.inserts[i].offset != f.inserts[j].offset {
			return f.inserts[i].offset-f.inserts[j].offset < 0
		}
		return f.inserts[i].priority-f.inserts[j].priority < 0
	})
	pre := 0
	var res []byte
	for _, item := range f.inserts {
		if item.offset != pre {
			res = append(res, f.content[pre:item.offset]...)
			pre = item.offset
		}
		content := item.content
		if item.contentGetter != nil {
			content = utils.String.String2bytes(item.contentGetter())
		}
		res = append(res, content...)
	}
	res = append(res, f.content[pre:]...)
	return res
}

func (f *FileVisitor) insert(offset int, content string, priority int) {
	f.inserts = append(f.inserts, &insert{
		offset:   offset,
		content:  utils.String.String2bytes(content),
		priority: priority})
}

func (f *FileVisitor) insertAtPost(offset int, contentGetter func() string, priority int) {
	f.inserts = append(f.inserts, &insert{
		offset:        offset,
		contentGetter: contentGetter,
		priority:      priority,
	})
}

func (f *FileVisitor) insertBlockVisit(beginPos, endPos token.Pos) {
	fmt.Printf("block begin:%d,end:%d\n", f.offset(beginPos), f.offset(endPos))
	f.addBlock(beginPos, endPos)
}

func (f *FileVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		f.addAgentImport(n.Name.End())
	case *ast.BlockStmt:
		// If it's a switch or select, the body is a list of case clauses; don't tag the block itself.
		if len(n.List) > 0 {
			switch n.List[0].(type) {
			case *ast.CaseClause: // switch
				for _, n := range n.List {
					clause := n.(*ast.CaseClause)
					f.addCounters(clause.Colon+1, clause.Colon+1, clause.End(), clause.Body, false)
				}
				return f
			case *ast.CommClause: // select
				for _, n := range n.List {
					clause := n.(*ast.CommClause)
					f.addCounters(clause.Colon+1, clause.Colon+1, clause.End(), clause.Body, false)
				}
				return f
			}
		}
		f.addCounters(n.Lbrace, n.Lbrace+1, n.Rbrace+1, n.List, true) // +1 to step past closing brace.
	case *ast.IfStmt:
		if n.Init != nil {
			ast.Walk(f, n.Init)
		}
		ast.Walk(f, n.Cond)
		ast.Walk(f, n.Body)
		if n.Else == nil {
			return nil
		}
		// The elses are special, because if we have
		//	if x {
		//	} else if y {
		//	}
		// we want to cover the "if y". To do this, we need a place to drop the counter,
		// so we add a hidden block:
		//	if x {
		//	} else {
		//		if y {
		//		}
		//	}
		elseOffset := f.findText(n.Body.End(), "else")
		if elseOffset < 0 {
			panic("lost else")
		}
		newElseStart := elseOffset + 4
		f.insert(newElseStart, "{", 1)
		f.insert(f.offset(n.Else.End()), "}", 1)

		// We just created a block, now walk it.
		// Adjust the position of the new block to start after
		// the "else". That will cause it to follow the "{"
		// we inserted above.
		pos := f.fset.File(n.Body.End()).Pos(newElseStart)
		switch stmt := n.Else.(type) {
		case *ast.IfStmt:
			block := &ast.BlockStmt{
				Lbrace: pos,
				List:   []ast.Stmt{stmt},
				Rbrace: stmt.End(),
			}
			n.Else = block
		case *ast.BlockStmt:
			stmt.Lbrace = pos
		default:
			panic("unexpected node type in if")
		}
		ast.Walk(f, n.Else)
		return nil
	case *ast.SelectStmt:
		// Don't annotate an empty select - creates a syntax error.
		if n.Body == nil || len(n.Body.List) == 0 {
			return nil
		}
	case *ast.SwitchStmt:
		// Don't annotate an empty switch - creates a syntax error.
		if n.Body == nil || len(n.Body.List) == 0 {
			if n.Init != nil {
				ast.Walk(f, n.Init)
			}
			if n.Tag != nil {
				ast.Walk(f, n.Tag)
			}
			return nil
		}
	case *ast.TypeSwitchStmt:
		// Don't annotate an empty type switch - creates a syntax error.
		if n.Body == nil || len(n.Body.List) == 0 {
			if n.Init != nil {
				ast.Walk(f, n.Init)
			}
			ast.Walk(f, n.Assign)
			return nil
		}
	case *ast.FuncDecl:
		// Don't annotate functions with blank names - they cannot be executed.
		// Similarly for bodyless funcs.
		if n.Name.Name == "_" || n.Body == nil {
			return nil
		}
		fname := n.Name.Name
		// Skip AddUint32 and StoreUint32 if we're instrumenting
		// sync/atomic itself in atomic mode (out of an abundance of
		// caution), since as part of the instrumentation process we
		// add calls to AddUint32/StoreUint32, and we don't want to
		// somehow create an infinite loop.
		//
		// Note that in the current implementation (Go 1.20) both
		// routines are assembly stubs that forward calls to the
		// runtime/internal/atomic equivalents, hence the infinite
		// loop scenario is purely theoretical (maybe if in some
		// future implementation one of these functions might be
		// written in Go). See #57445 for more details.
		if atomicOnAtomic() && (fname == "AddUint32" || fname == "StoreUint32") {
			return nil
		}
		// Determine proper function or method name.
		if r := n.Recv; r != nil && len(r.List) == 1 {
			t := r.List[0].Type
			star := ""
			if p, _ := t.(*ast.StarExpr); p != nil {
				t = p.X
				star = "*"
			}
			if p, _ := t.(*ast.Ident); p != nil {
				fname = star + p.Name + "." + fname
			}
		}
		fmt.Printf("func %s\n", fname)
		f.addFuncByDecl(fname, n)
		ast.Walk(f, n.Body)
		return nil
	case *ast.FuncLit:

		// Hack: function literals aren't named in the go/ast representation,
		// and we don't know what name the compiler will choose. For now,
		// just make up a descriptive name.
		pos := n.Pos()
		p := f.fset.File(pos).Position(pos)
		fname := fmt.Sprintf("func.L%d.C%d", p.Line, p.Column)
		fmt.Printf("func %s\n", fname)
		f.addFuncByLint(fname, n)
		ast.Walk(f, n.Body)
		// todo get parent func
		return nil
	case *ast.GoStmt:
		switch expr := n.Call.Fun.(type) {
		case *ast.Ident:
			f.insert(f.offset(expr.Pos()), "func(){", 1)
			f.insert(f.offset(n.Call.End()), "}()", 1)
		case *ast.FuncLit:
		}
	}
	return f
}

func (f *FileVisitor) addAgentImport(pos token.Pos) {
	f.insert(f.fset.Position(pos).Offset, fmt.Sprintf(";import \"%s\";", f.prog.sdkImportPkgPath), 1)
}

func (f *FileVisitor) addFuncByLint(fName string, lit *ast.FuncLit) {
	f.addFunc(fName, f.fset.Position(lit.Pos()), f.fset.Position(lit.End()), f.fset.Position(lit.Body.Pos()))
}

func (f *FileVisitor) addFuncByDecl(fName string, decl *ast.FuncDecl) {
	f.addFunc(fName, f.fset.Position(decl.Pos()), f.fset.Position(decl.End()), f.fset.Position(decl.Body.Pos()))
}

func (f *FileVisitor) addFunc(fName string, begin, end token.Position, bodyBegin token.Position) {
	fc := f.prog.addMethod(f.fileId, fName, begin, end)
	f.insert(bodyBegin.Offset+1,
		fmt.Sprintf(
			`_ellynCtx := ellyn_agent.Agent.GetCtx();ellyn_agent.Agent.Push(_ellynCtx,%d);defer ellyn_agent.Agent.Pop(_ellynCtx);`, fc.Id),
		1)
}

func (f *FileVisitor) addBlock(begin, end token.Pos) {
	block := f.prog.addBlock(f.fileId, f.fset.Position(begin), f.fset.Position(end))
	f.insertAtPost(block.Begin.Offset, func() string {
		return fmt.Sprintf("ellyn_agent.Agent.VisitBlock(_ellynCtx,%d);", block.MethodOffset)
	}, 2)
}

// offset translates a token position into a 0-indexed byte offset.
func (f *FileVisitor) offset(pos token.Pos) int {
	return f.fset.Position(pos).Offset
}

// findText finds text in the original source, starting at pos.
// It correctly skips over comments and assumes it need not
// handle quoted strings.
// It returns a byte offset within f.src.
func (f *FileVisitor) findText(pos token.Pos, text string) int {
	b := []byte(text)
	start := f.offset(pos)
	i := start
	s := f.content
	for i < len(s) {
		if bytes.HasPrefix(s[i:], b) {
			return i
		}
		if i+2 <= len(s) && s[i] == '/' && s[i+1] == '/' {
			for i < len(s) && s[i] != '\n' {
				i++
			}
			continue
		}
		if i+2 <= len(s) && s[i] == '/' && s[i+1] == '*' {
			for i += 2; ; i++ {
				if i+2 > len(s) {
					return 0
				}
				if s[i] == '*' && s[i+1] == '/' {
					i += 2
					break
				}
			}
			continue
		}
		i++
	}
	return -1
}

func (f *FileVisitor) addCounters(pos, insertPos, blockEnd token.Pos, list []ast.Stmt, extendToClosingBrace bool) {
	// Special case: make sure we add a counter to an empty block. Can't do this below
	// or we will add a counter to an empty statement list after, say, a return statement.
	if len(list) == 0 {
		//f.edit.Insert(f.offset(insertPos), f.newCounter(insertPos, blockEnd, 0)+";")
		f.insertBlockVisit(insertPos, blockEnd)
		return
	}
	// Make a copy of the list, as we may mutate it and should leave the
	// existing list intact.
	list = append([]ast.Stmt(nil), list...)
	// We have a block (statement list), but it may have several basic blocks due to the
	// appearance of statements that affect the flow of control.
	for {
		// Find first statement that affects flow of control (break, continue, if, etc.).
		// It will be the last statement of this basic block.
		var last int
		end := blockEnd
		for last = 0; last < len(list); last++ {
			stmt := list[last]
			end = f.statementBoundary(stmt)
			if f.endsBasicSourceBlock(stmt) {
				// If it is a labeled statement, we need to place a counter between
				// the label and its statement because it may be the target of a goto
				// and thus start a basic block. That is, given
				//	foo: stmt
				// we need to create
				//	foo: ; stmt
				// and mark the label as a block-terminating statement.
				// The result will then be
				//	foo: COUNTER[n]++; stmt
				// However, we can't do this if the labeled statement is already
				// a control statement, such as a labeled for.
				if label, isLabel := stmt.(*ast.LabeledStmt); isLabel && !f.isControl(label.Stmt) {
					newLabel := *label
					newLabel.Stmt = &ast.EmptyStmt{
						Semicolon: label.Stmt.Pos(),
						Implicit:  true,
					}
					end = label.Pos() // Previous block ends before the label.
					list[last] = &newLabel
					// Open a gap and drop in the old statement, now without a label.
					list = append(list, nil)
					copy(list[last+1:], list[last:])
					list[last+1] = label.Stmt
				}
				last++
				extendToClosingBrace = false // Block is broken up now.
				break
			}
		}
		if extendToClosingBrace {
			end = blockEnd
		}
		if pos != end { // Can have no source to cover if e.g. blocks abut.
			f.insertBlockVisit(insertPos, end)
		}

		list = list[last:]
		if len(list) == 0 {
			break
		}
		pos = list[0].Pos()
		insertPos = pos
	}
}

// statementBoundary finds the location in s that terminates the current basic
// block in the source.
func (f *FileVisitor) statementBoundary(s ast.Stmt) token.Pos {
	// Control flow statements are easy.
	switch s := s.(type) {
	case *ast.BlockStmt:
		// Treat blocks like basic blocks to avoid overlapping counters.
		return s.Lbrace
	case *ast.IfStmt:
		found, pos := hasFuncLiteral(s.Init)
		if found {
			return pos
		}
		found, pos = hasFuncLiteral(s.Cond)
		if found {
			return pos
		}
		return s.Body.Lbrace
	case *ast.ForStmt:
		found, pos := hasFuncLiteral(s.Init)
		if found {
			return pos
		}
		found, pos = hasFuncLiteral(s.Cond)
		if found {
			return pos
		}
		found, pos = hasFuncLiteral(s.Post)
		if found {
			return pos
		}
		return s.Body.Lbrace
	case *ast.LabeledStmt:
		return f.statementBoundary(s.Stmt)
	case *ast.RangeStmt:
		found, pos := hasFuncLiteral(s.X)
		if found {
			return pos
		}
		return s.Body.Lbrace
	case *ast.SwitchStmt:
		found, pos := hasFuncLiteral(s.Init)
		if found {
			return pos
		}
		found, pos = hasFuncLiteral(s.Tag)
		if found {
			return pos
		}
		return s.Body.Lbrace
	case *ast.SelectStmt:
		return s.Body.Lbrace
	case *ast.TypeSwitchStmt:
		found, pos := hasFuncLiteral(s.Init)
		if found {
			return pos
		}
		return s.Body.Lbrace
	}
	// If not a control flow statement, it is a declaration, expression, call, etc. and it may have a function literal.
	// If it does, that's tricky because we want to exclude the body of the function from this block.
	// Draw a line at the start of the body of the first function literal we find.
	// TODO: what if there's more than one? Probably doesn't matter much.
	found, pos := hasFuncLiteral(s)
	if found {
		return pos
	}
	return s.End()
}

// endsBasicSourceBlock reports whether s changes the flow of control: break, if, etc.,
// or if it's just problematic, for instance contains a function literal, which will complicate
// accounting due to the block-within-an expression.
func (f *FileVisitor) endsBasicSourceBlock(s ast.Stmt) bool {
	switch s := s.(type) {
	case *ast.BlockStmt:
		// Treat blocks like basic blocks to avoid overlapping counters.
		return true
	case *ast.BranchStmt:
		return true
	case *ast.ForStmt:
		return true
	case *ast.IfStmt:
		return true
	case *ast.LabeledStmt:
		return true // A goto may branch here, starting a new basic block.
	case *ast.RangeStmt:
		return true
	case *ast.SwitchStmt:
		return true
	case *ast.SelectStmt:
		return true
	case *ast.TypeSwitchStmt:
		return true
	case *ast.ExprStmt:
		// Calls to panic change the flow.
		// We really should verify that "panic" is the predefined function,
		// but without type checking we can't and the likelihood of it being
		// an actual problem is vanishingly small.
		if call, ok := s.X.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" && len(call.Args) == 1 {
				return true
			}
		}
	}
	found, _ := hasFuncLiteral(s)
	return found
}

// isControl reports whether s is a control statement that, if labeled, cannot be
// separated from its label.
func (f *FileVisitor) isControl(s ast.Stmt) bool {
	switch s.(type) {
	case *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.SelectStmt, *ast.TypeSwitchStmt:
		return true
	}
	return false
}

func atomicOnAtomic() bool {
	return true
}

// funcLitFinder implements the ast.Visitor pattern to find the location of any
// function literal in a subtree.
type funcLitFinder token.Pos

func (f *funcLitFinder) Visit(node ast.Node) (w ast.Visitor) {
	if f.found() {
		return nil // Prune search.
	}
	switch n := node.(type) {
	case *ast.FuncLit:
		*f = funcLitFinder(n.Body.Lbrace)
		return nil // Prune search.
	}
	return f
}

func (f *funcLitFinder) found() bool {
	return token.Pos(*f) != token.NoPos
}

// hasFuncLiteral reports the existence and position of the first func literal
// in the node, if any. If a func literal appears, it usually marks the termination
// of a basic block because the function body is itself a block.
// Therefore we draw a line at the start of the body of the first function literal we find.
// TODO: what if there's more than one? Probably doesn't matter much.
func hasFuncLiteral(n ast.Node) (bool, token.Pos) {
	if n == nil {
		return false, 0
	}
	var literal funcLitFinder
	ast.Walk(&literal, n)
	return literal.found(), token.Pos(literal)
}
