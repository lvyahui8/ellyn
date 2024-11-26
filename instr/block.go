package instr

// 基于此改造 https://github.com/golang/go/blob/master/src/cmd/cover/cover.go
import (
	"bytes"
	"go/ast"
	"go/token"
)

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
