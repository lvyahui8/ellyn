package ellyn_ast

type Handler interface {
	FileFilter(program Program, pkg Package, file SourceFile) bool
}
