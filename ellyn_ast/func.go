package ellyn_ast

import "go/token"

type GoFunc struct {
	begin token.Position
	end   token.Position
}
