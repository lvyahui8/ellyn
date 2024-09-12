package ellyn_ast

import "go/token"

type Block struct {
	id    uint32
	fc    *GoFunc
	begin token.Position
	end   token.Position
}
