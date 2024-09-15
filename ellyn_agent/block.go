package ellyn_agent

import _ "embed"

//go:embed blocks.dat
var blocksData []byte

type allBlocks []*Block

type Block struct {
	Id uint
}

func initBlocks() {

}
