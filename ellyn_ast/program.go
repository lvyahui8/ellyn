package ellyn_ast

type Program struct {
	mainPkg   Package
	rootPkg   Package
	pkgMap    map[string]Package
	allFuncs  []*GoFunc
	allBlocks []*Block
}

func NewProgram(mainPkgDir string) *Program {
	return nil
}
