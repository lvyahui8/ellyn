package ellyn_ast

import "go/token"

type GoFunc struct {
	id    uint32
	name  string
	begin token.Position
	end   token.Position
}
