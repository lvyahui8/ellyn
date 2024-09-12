package ellyn_ast

import "go/token"

type GoFunc struct {
	id    uint32
	begin token.Position
	end   token.Position
}
