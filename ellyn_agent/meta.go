package ellyn_agent

import (
	_ "embed"
	"fmt"
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
	"reflect"
	"strconv"
	"strings"
)

func initMetaData() {
	packages = initCsvData[*Package](packagesData)
	files = initCsvData[*File](filesData)
	methods = initCsvData[*Method](filesData)
	blocks = initCsvData[*Block](blocksData)
}

type CsvRow interface {
	parse(cols []string)
	encodeRow() string
}

func initCsvData[T CsvRow](compressedContent []byte) []T {
	if len(compressedContent) == 0 {
		return nil
	}
	csvContent := utils.String.Bytes2string(utils.Gzip.UnCompress(compressedContent))
	lines := strings.Split(csvContent, "\n")
	var res []T
	for _, line := range lines {
		cols := strings.Split(line, ",")
		var t T
		val := reflect.ValueOf(t)
		val.Set(reflect.New(val.Elem().Type()))
		t.parse(cols)
		res = append(res, t)
	}
	return res
}

func parseId(col string) uint32 {
	id, err := strconv.ParseUint(col, 10, 32)
	asserts.IsNil(err)
	return uint32(id)
}

//go:embed meta/packages.dat
var packagesData []byte

var packages []*Package

type Package struct {
	Id   uint32
	Name string
	Path string
}

func (p *Package) encodeRow() string {
	return fmt.Sprintf("%d,%s,%s", p.Id, p.Name, p.Path)
}

func (p *Package) parse(cols []string) {
	p.Id = parseId(cols[0])
	p.Name = cols[1]
	p.Path = cols[2]
}

//go:embed meta/files.dat
var filesData []byte

var files []*File

type File struct {
	FileId       uint32
	RelativePath string
}

func (f *File) encodeRow() string {
	return fmt.Sprintf("%d,%s", f.FileId, f.RelativePath)
}

func (f *File) parse(cols []string) {
	f.FileId = parseId(cols[0])
	f.RelativePath = cols[1]
}

//go:embed meta/methods.dat
var methodsData []byte

var methods []*Method

type Method struct {
	Id             uint32
	FullName       string // 完成函数名
	FileId         uint32 // 所在文件id
	PackageId      uint32 // 包id
	Blocks         []*Block
	BlockCnt       int
	ArgsTypeList   []reflect.Type
	ReturnTypeList []reflect.Type
}

func (m *Method) encodeRow() string {
	return fmt.Sprintf("%d,%s,%d,%d,%d", m.Id, m.FullName, m.FileId, m.PackageId, m.BlockCnt)
}

func (m *Method) parse(cols []string) {
	m.Id = parseId(cols[0])
	m.FullName = cols[1]
	m.FileId = parseId(cols[2])
	m.PackageId = parseId(cols[3])
	m.BlockCnt = int(parseId(cols[4]))
	m.Blocks = make([]*Block, m.BlockCnt)
}

func newMethodBlockBits(methodId uint32) *collections.BitMap {
	return collections.NewBitMap(uint(len(methods[methodId].Blocks)))
}

//go:embed meta/blocks.dat
var blocksData []byte

var blocks []*Block

type Block struct {
	Id           uint32
	MethodId     uint32
	MethodOffset int
}

func (b *Block) encodeRow() string {
	return fmt.Sprintf("%d,%d,%d", b.Id, b.MethodId, b.MethodOffset)
}

func (b *Block) parse(cols []string) {
	b.Id = parseId(cols[0])
	b.MethodId = parseId(cols[1])
	b.MethodOffset = int(parseId(cols[2]))
	method := methods[b.MethodId]
	method.Blocks[b.MethodOffset] = b
}
