package agent

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"github.com/lvyahui8/ellyn/sdk/common/logging"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
	"reflect"
	"strconv"
	"strings"
)

// meta 元数据
var meta embed.FS

// getDat 读取元数据文件
func getDat(file string) []byte {
	data, _ := meta.ReadFile(MetaRelativePath + "/" + file)
	return data
}

// initMetaData 初始化元数据
func initMetaData() {
	packages = initCsvData[*Package](getDat(MetaPackages))
	files = initCsvData[*File](getDat(MetaFiles))
	methods = initCsvData[*Method](getDat(MetaMethods))
	blocks = initCsvData[*Block](getDat(MetaBlocks))
	log.InfoKV(logging.Empty().Int("packages", len(packages)).
		Int("files", len(files)).
		Int("methods", len(methods)).
		Int("blocks", len(blocks)))
}

// CsvRow 元数据文件行的编码解码定义。不同元数据文件分别为不同的实现类
type CsvRow interface {
	// parse 解码行，从一行csv文件中解码出数据
	parse(cols []string)
	// encodeRow 编码行，将数据编码成csv文件中的一行
	encodeRow() string
}

// initCsvData 从压缩csv数据中反序列化出对象实例
func initCsvData[T CsvRow](compressedContent []byte) []T {
	if len(compressedContent) == 0 {
		return nil
	}
	csvContent := utils.String.Bytes2string(utils.Gzip.UnCompress(compressedContent))
	lines := strings.Split(csvContent, "\n")
	var res []T
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		cols := strings.Split(line, ",")
		var t T
		val := reflect.New(reflect.ValueOf(t).Type().Elem())
		t = val.Interface().(T)
		t.parse(cols)
		res = append(res, t)
	}
	return res
}

// EncodeCsvRows 将对象实例列表编码成压缩的csv数据
func EncodeCsvRows[T CsvRow](rows []T) []byte {
	var res []byte
	for _, row := range rows {
		res = append(res, row.encodeRow()...)
		res = append(res, '\n')
	}
	return utils.Gzip.Compress(res)
}

func parseUint32(col string) uint32 {
	id, err := strconv.ParseUint(col, 10, 32)
	asserts.IsNil(err)
	return uint32(id)
}

// Pos 定义在文件中的位置
type Pos struct {
	// Offset 文件偏移量，字节数
	Offset int `json:"offset"`
	// Line 所在行
	Line int `json:"line"`
	// Column 所在列
	Column int `json:"column"`
}

// NewPos 初始化行
func NewPos(offset, line, column int) *Pos {
	return &Pos{offset, line, column}
}

func (p *Pos) String() string {
	return fmt.Sprintf("L%dC%d:%d", p.Line, p.Column, p.Offset)
}

// ParsePos 解码Pos , 格式:L%dC%d:%d
func ParsePos(encodedPos string) *Pos {
	colIdx := strings.Index(encodedPos, "C")
	offsetIdx := strings.Index(encodedPos, ":")
	return &Pos{
		Line:   int(parseUint32(encodedPos[1:colIdx])),
		Column: int(parseUint32(encodedPos[colIdx+1 : offsetIdx])),
		Offset: int(parseUint32(encodedPos[offsetIdx+1:])),
	}
}

// packages 目标项目所有的package
var packages []*Package

// Package 代表golang package
type Package struct {
	Id uint32
	// Name pkg名，如ellyn
	Name string
	// Path Pkg全路径，即写代码时的Import path. 如：github.com/lvyahui8/ellyn
	Path string
	// Dir pkg在本地文件系统的绝对路径
	Dir string
}

func NewPackage(dir, path string) *Package {
	items := strings.Split(path, "/")
	name := items[len(items)-1]
	return &Package{Dir: dir, Name: name, Path: path}
}

func (p *Package) encodeRow() string {
	return fmt.Sprintf("%d,%s,%s", p.Id, p.Name, p.Path)
}

func (p *Package) parse(cols []string) {
	p.Id = parseUint32(cols[0])
	p.Name = cols[1]
	p.Path = cols[2]
}

// files 目标项目所有的源码文件
var files []*File

// File 代表golang源码文件
type File struct {
	// FileId 文件id，遍历的过程中生成
	FileId uint32
	// PackageId golang package id，遍历的过程中生成
	PackageId uint32
	// RelativePath 在目标项目中的相对路径
	RelativePath string
	// LineNum 文件总行数
	LineNum int
}

func (f *File) encodeRow() string {
	return fmt.Sprintf("%d,%d,%s,%d", f.FileId, f.PackageId, f.RelativePath, f.LineNum)
}

func (f *File) parse(cols []string) {
	f.FileId = parseUint32(cols[0])
	f.PackageId = parseUint32(cols[1])
	f.RelativePath = cols[2]
	f.LineNum = int(parseUint32(cols[3]))
}

// VarDef 代表变量声明
type VarDef struct {
	// Names 变量名列表
	Names []string
	// Type 变量类型的文字表示
	Type string
}

// VarDefList 参数列表，主要表示方法的出参和入参
type VarDefList struct {
	// list 参数列表
	list []*VarDef
	// idx2name 所处位置的变量名
	idx2name []string
	// idx2type 所处位置的变量类型
	idx2type []string
}

// NewVarDefList 基于多个参数声明创建列表
func NewVarDefList(list []*VarDef) *VarDefList {
	res := &VarDefList{
		list: list,
	}
	if list == nil {
		return res
	}
	for _, def := range list {
		for _, name := range def.Names {
			res.idx2type = append(res.idx2type, def.Type)
			res.idx2name = append(res.idx2name, name)
		}
	}
	return res
}

func (vdl *VarDefList) Encode() string {
	var list []string
	for _, def := range vdl.list {
		list = append(list, fmt.Sprintf("{%s}%s", strings.Join(def.Names, ":"), def.Type))
	}
	return strings.Join(list, ";")
}

func (vdl *VarDefList) Type(idx int) string {
	return vdl.idx2type[idx]
}

func (vdl *VarDefList) Name(idx int) string {
	return vdl.idx2name[idx]
}

func (vdl *VarDefList) Count() int {
	return len(vdl.idx2name)
}

func (vdl *VarDefList) MarshalJSON() ([]byte, error) {
	return json.Marshal(vdl.Encode())
}

func decodeVarDef(str string) *VarDefList {
	if str == "" {
		return NewVarDefList(nil)
	}
	items := strings.Split(str, ";")
	var list []*VarDef
	for _, item := range items {
		idx := strings.Index(item, "}")
		list = append(list, &VarDef{strings.Split(item[1:idx], ":"), item[idx+1 : len(item)]})
	}
	return NewVarDefList(list)
}

// methods 目标项目所有方法
var methods []*Method

type Method struct {
	Id         uint32
	Name       string
	FullName   string // 完整函数名
	FileId     uint32 // 所在文件id
	PackageId  uint32 // 包id
	Blocks     []*Block
	BlockCnt   int
	Begin      *Pos
	End        *Pos
	ArgsList   *VarDefList
	ReturnList *VarDefList
}

func (m *Method) encodeRow() string {
	return fmt.Sprintf("%d,%s,%d,%d,%d,%s,%s,%s,%s",
		m.Id, m.FullName, m.FileId, m.PackageId, len(m.Blocks), m.Begin, m.End,
		m.ArgsList.Encode(),
		m.ReturnList.Encode())
}

func (m *Method) parse(cols []string) {
	m.Id = parseUint32(cols[0])
	m.FullName = cols[1]
	m.FileId = parseUint32(cols[2])
	m.PackageId = parseUint32(cols[3])
	m.BlockCnt = int(parseUint32(cols[4]))
	m.Blocks = make([]*Block, m.BlockCnt)
	m.Begin = ParsePos(cols[5])
	m.End = ParsePos(cols[6])
	m.ArgsList = decodeVarDef(cols[7])
	m.ReturnList = decodeVarDef(cols[8])
}

func newMethodBlockFlags(methodId uint32) []bool {
	return make([]bool, methods[methodId].BlockCnt)
}

var blocks []*Block

type Block struct {
	Id           uint32
	FileId       uint32
	MethodId     uint32
	MethodOffset int
	Begin        *Pos
	End          *Pos
}

func (b *Block) encodeRow() string {
	return fmt.Sprintf("%d,%d,%d,%s,%s,%d", b.Id, b.MethodId, b.MethodOffset, b.Begin, b.End, b.FileId)
}

func (b *Block) parse(cols []string) {
	b.Id = parseUint32(cols[0])
	b.MethodId = parseUint32(cols[1])
	b.MethodOffset = int(parseUint32(cols[2]))
	method := methods[b.MethodId]
	method.Blocks[b.MethodOffset] = b
	b.Begin = ParsePos(cols[3])
	b.End = ParsePos(cols[4])
	b.FileId = parseUint32(cols[5])
}
