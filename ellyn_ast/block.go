package ellyn_ast

import "go/token"

type Block struct {
	fc    *GoFunc
	begin token.Position
	end   token.Position
}
