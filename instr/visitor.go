package instr

import (
	"fmt"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

const varNamePrefix = "_ellynVar"
const notCollectedVarName = "ellyn_agent.NotCollected"

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
	inserts := f.inserts
	if len(inserts) == 1 {
		// only import
		inserts = nil
	}
	for _, item := range inserts {
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
	//fmt.Printf("block begin:%d,end:%d\n", f.offset(beginPos), f.offset(endPos))
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
		//fmt.Printf("func %s\n", fname)
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
		//fmt.Printf("func %s\n", fname)
		f.addFuncByLint(fname, n)
		ast.Walk(f, n.Body)
		// todo get parent func
		return nil
	case *ast.GoStmt:
		f.wrapGo(n)
	}
	return f
}

func (f *FileVisitor) wrapGo(n *ast.GoStmt) {
	f.insert(f.offset(n.Pos()),
		"{_ellynCtxId,_ellynFromMethod := _ellynCtx.Snapshot();", 100)
	initCtxCode := "ellyn_agent.Agent.InitCtx(_ellynCtxId,_ellynFromMethod);"
	switch expr := n.Call.Fun.(type) {
	case *ast.Ident:
		f.insert(f.offset(expr.Pos()), "func(){"+initCtxCode, 1)
		f.insert(f.offset(n.Call.End()), "}()", 1)
	case *ast.FuncLit:
		f.insert(f.offset(expr.Body.Lbrace)+1, initCtxCode, 0)
	}
	f.insert(f.offset(n.End()), "}", 2)
}

func (f *FileVisitor) parseVarLists() {
	// - 参数列表和返回值列表
	// - 未命名需要增加命名，要考虑匿名参数
	// - 判断参数类型，考虑到值传递参数拷贝对性能影响，只拷贝小对象
	// 	- 传递：基础类型、slice、pointer、iface、 string（限制长度）、error
	//  - 不传递：Array,struct,eface、方法等

}

func (f *FileVisitor) addAgentImport(pos token.Pos) {
	f.insert(f.fset.Position(pos).Offset, fmt.Sprintf(";import \"%s\";", f.prog.sdkImportPkgPath), 1)
}

func (f *FileVisitor) addFuncByLint(fName string, lit *ast.FuncLit) {
	f.addFunc(fName, f.fset.Position(lit.Pos()), f.fset.Position(lit.End()), f.fset.Position(lit.Body.Pos()), lit.Type)
}

func (f *FileVisitor) addFuncByDecl(fName string, decl *ast.FuncDecl) {
	f.addFunc(fName, f.fset.Position(decl.Pos()), f.fset.Position(decl.End()), f.fset.Position(decl.Body.Pos()), decl.Type)
}

func (f *FileVisitor) modifyParamsAndResults(funcType *ast.FuncType) (params []string, results []string) {
	_, params = f.modifyVarList(funcType.Params, "Param")
	cnt, results := f.modifyVarList(funcType.Results, "Ret")
	if cnt > 0 && len(funcType.Results.List) == 1 {
		// 返回值列表只有1个并且之前是匿名，命名后要加括号
		f.insert(f.offset(funcType.Results.Pos()), "(", 0)
		f.insert(f.offset(funcType.Results.End()), ")", 0)
	}
	return
}

func (f *FileVisitor) modifyVarList(list *ast.FieldList, namePrefix string) (modifiedCnt int, nameList []string) {
	if list == nil {
		return
	}
	for _, item := range list.List {
		// 判断类型，go属于值传递（指针、引用本身也是一个值），对于可能存在大量拷贝的类型不做收集
		// 遍历变量列表，对未命名的进行命名
		notCollectedType := false
		// 统一传指针
		addressOf := "&"
		switch item.Type.(type) {
		case *ast.FuncType:
			notCollectedType = true
		case *ast.StarExpr:
			addressOf = ""
		}
		if len(item.Names) == 0 {
			// 匿名， 生成变量名并插入
			varName := fmt.Sprintf("%s%s%d", varNamePrefix, namePrefix, modifiedCnt)
			modifiedCnt++
			f.insert(f.offset(item.Pos()), varName+" ", 1)

			if notCollectedType {
				nameList = append(nameList, notCollectedVarName)
			} else {
				nameList = append(nameList, addressOf+varName)
			}
		} else {
			for _, n := range item.Names {
				if n.Name == "_" || notCollectedType {
					nameList = append(nameList, notCollectedVarName)
				} else {
					nameList = append(nameList, addressOf+n.Name)
				}
			}
		}
	}
	return
}

func (f *FileVisitor) addFunc(fName string, begin, end, bodyBegin token.Position, funcType *ast.FuncType) {
	fc := f.prog.addMethod(f.fileId, fName, begin, end, funcType)

	format := "_ellynCtx,_ellynCollect,_ellynCleaner := ellyn_agent.Agent.GetCtx();" +
		"if _ellynCleaner != nil { defer _ellynCleaner() };" +
		"if _ellynCollect { ellyn_agent.Agent.Push(_ellynCtx,%d,%s);" +
		"defer ellyn_agent.Agent.Pop(_ellynCtx,%s) };"
	var args []any
	args = append(args, fc.Id)

	if f.prog.conf.NoArgs {
		args = append(args, "nil", "nil")
	} else {
		params, results := f.modifyParamsAndResults(funcType)
		args = append(args, fmt.Sprintf("[]any{%s}", strings.Join(params, ",")))
		args = append(args, fmt.Sprintf("[]any{%s}", strings.Join(results, ",")))
	}
	f.insert(bodyBegin.Offset+1,
		fmt.Sprintf(format, args...), 1)
}

func (f *FileVisitor) addBlock(begin, end token.Pos) {
	block := f.prog.addBlock(f.fileId, f.fset.Position(begin), f.fset.Position(end))
	f.insertAtPost(block.Begin.Offset, func() string {
		return fmt.Sprintf("if _ellynCollect { ellyn_agent.Agent.Mark(_ellynCtx,%d,%d) };", block.MethodOffset, block.Id)
	}, 2)
}

// offset translates a token position into a 0-indexed byte offset.
func (f *FileVisitor) offset(pos token.Pos) int {
	return f.fset.Position(pos).Offset
}
