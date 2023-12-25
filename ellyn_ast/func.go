package ellyn_ast

import "go/token"

type GoFunc struct {
	file  *SourceFile
	begin token.Position
	end   token.Position
}
