package ellyn_agent

import (
	_ "embed"
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"reflect"
)

//go:embed methods.dat
var methodsData []byte

var allMethods []*Method

type Method struct {
	Id             uint32
	Blocks         []*Block
	ArgsTypeList   []reflect.Type
	ReturnTypeList []reflect.Type
}

func newMethodBlockBits(methodId uint32) *collections.BitMap {
	return collections.NewBitMap(uint(len(allMethods[methodId].Blocks)))
}

func initMethods() {

}
