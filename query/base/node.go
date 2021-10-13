package base

import (
	"reflect"

	flatbuffers "github.com/google/flatbuffers/go"
)

// Node base struct.
// Pos is start position of buffer
// ValueInfos is information of pointted fields
type Node struct {
	Base
	Pos  int
	Size int
	TLen uint16
}

// NodeList is struct for Vector(list) Node
type NodeList struct {
	*Node
	ValueInfo
}

// ValueInfo is information of flatbuffer's table/struct field
type ValueInfo struct {
	Pos  int
	Size int
	VLen uint32
}

// Info is infotion of flatbuffer's table/struct
type Info ValueInfo

// NodePath is path for traverse data tree
type NodePath struct {
	Name string
	Idx  int
}

// NewNode ... this provide creation of Node
// Node share buffer in same tree.
//    Base is buffer
//    pos is start position in Base's buffer
func NewNode(b Base, pos int) *Node {
	return NewNode2(b, pos, false)
}

// NewNode2 ...  provide skip to initialize Vtable
func NewNode2(b Base, pos int, noLoadVTable bool) *Node {
	node := &Node{Base: b, Pos: pos, Size: -1}
	if !noLoadVTable {
		node.preLoadVtable()
	}
	return node
}

func (n *Node) vtable() {
	n.preLoadVtable()
}

func (n *Node) preLoadVtable() {

	vOffset := int(flatbuffers.GetUOffsetT(n.R(n.Pos)))
	vPos := int(n.Pos) - vOffset
	vLen := int(flatbuffers.GetVOffsetT(n.R(vPos)))
	n.TLen = uint16(flatbuffers.GetVOffsetT(n.R(vPos + 2)))

	for cur := vPos + 4; cur < vPos+vLen; cur += 2 {
		if len(n.R(cur)) < 2 {
			a := n.R(cur)
			_ = a
			panic("invalid preLoadVtable")
		}
		flatbuffers.GetVOffsetT(n.R(cur))
	}
}

// FbsString ... return []bytes for string data
func FbsString(node *Node) []byte {
	pos := node.VirtualTable(0)
	sLenOff := int(flatbuffers.GetUint32(node.R(pos)))
	sLen := int(flatbuffers.GetUint32(node.R(pos + sLenOff)))
	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	return node.R(start)[:sLen]
}

// FbsStringInfo ... return Node Infomation for FbsString
func FbsStringInfo(node *Node) Info {
	pos := node.VirtualTable(0)
	sLenOff := int(flatbuffers.GetUint32(node.R(pos)))
	sLen := flatbuffers.GetUint32(node.R(pos + sLenOff))
	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	return Info{Pos: start, Size: int(sLen)}
}

// IsNotReady ... already set ValueInfo (field information)
func (info ValueInfo) IsNotReady() bool {
	return info.Pos < 1

}

// ValueInfoPosBytes ... etching vtable position infomation for []byte
func (node *Node) ValueInfoPosBytes(vIdx int) (info ValueInfo) {

	info = node.ValueInfoPosList(vIdx)
	info.Size = int(info.VLen)
	return info
}

// ValueInfoPosList ... etching vtable position infomation for flatbuffers vector
func (node *Node) ValueInfoPosList(vIdx int) (info ValueInfo) {

	vPos := node.VirtualTable(vIdx)
	if node.VirtualTableIsZero(vIdx) {
		return info
	}

	vLenOff := int(flatbuffers.GetUint32(node.R(vPos)))
	vLen := flatbuffers.GetUint32(node.R(vPos + vLenOff))
	start := vPos + vLenOff + flatbuffers.SizeUOffsetT

	info.Pos = int(start)
	info.VLen = vLen

	return info

}

func (node *Node) ValueTable(vIdx int) *Node {
	return NewNode(node.Base,
		node.Table(vIdx))
}

func (node *Node) ValueStruct(vIdx int) *Node {

	return NewNode2(node.Base, node.VirtualTable(vIdx), true)

}

func (node *Node) ValueList(vIdx int) NodeList {

	info := node.ValueInfoPosList(vIdx)

	return NodeList{Node: NewNode(node.Base, node.Pos),
		ValueInfo: info}
}

type UnmarshalFn func(string, reflect.Value) error

func (node *Node) Unmarshal(ptr interface{}, setter UnmarshalFn) error {

	return node.unmarshal(ptr, setter)
}

func (node *Node) unmarshal(ptr interface{}, setter UnmarshalFn) error {

	rv := reflect.ValueOf(ptr)
	_ = rv

	if rv.Kind() != reflect.Ptr {
		return ERR_MUST_POINTER
	}
	z := rv.Elem().Kind()
	_ = z
	if rv.Elem().Kind() != reflect.Struct {
		return ERR_INVALID_TYPE
	}

	t := rv.Elem().Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tName, ok := field.Tag.Lookup(FBS_TAG_NAME)
		if !ok {
			continue
		}
		if err := setter(tName, rv.Elem().Field(i)); err != nil {
			return err
		}
	}
	return nil
}

type Value struct {
	*Node
	S int
	E int
}

func (v Value) String() string {
	return string(v.R(v.S)[:v.E])
}

// mock
func (nList *NodeList) Member(i int) interface{} {
	return nil
}

type FieldsDefile struct {
	IdxToTyoe      map[int]int
	IdxToTypeGroup map[int]int
}

func (node *Node) Byte() byte {
	if node == nil {
		return 0
	}

	return flatbuffers.GetByte(node.R(node.Pos))
}

func (node *Node) Bool() bool {
	if node == nil {
		return false
	}

	return flatbuffers.GetBool(node.R(node.Pos))
}

func (node *Node) Uint8() uint8 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetUint8(node.R(node.Pos))
}

func (node *Node) Uint16() uint16 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetUint16(node.R(node.Pos))
}

func (node *Node) Uint32() uint32 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetUint32(node.R(node.Pos))
}

func (node *Node) Uint64() uint64 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetUint64(node.R(node.Pos))
}

func (node *Node) Int8() int8 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetInt8(node.R(node.Pos))
}

func (node *Node) Int16() int16 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetInt16(node.R(node.Pos))
}

func (node *Node) Int32() int32 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetInt32(node.R(node.Pos))
}

func (node *Node) Int64() int64 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetInt64(node.R(node.Pos))
}

func (node *Node) Float32() float32 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetFloat32(node.R(node.Pos))
}

func (node *Node) Float64() float64 {
	if node == nil {
		return 0
	}

	return flatbuffers.GetFloat64(node.R(node.Pos))
}

func (node *Node) Bytes() []byte {
	if node == nil {
		return nil
	}
	return node.R(node.Pos)[:node.Size]
}

// VirtualTableIsZero ... return checking VTable is empty
func (node *Node) VirtualTableIsZero(idx int) bool {

	return node.VirtualTable(idx) == node.Pos
}

// VirtualTable ... return VTable.
func (node *Node) VirtualTable(idx int) int {

	voff := flatbuffers.GetUint32(node.R(node.Pos))
	vPos := node.Pos - int(voff)

	if node.ShouldCheckBound() && node.LenBuf() <= int(vPos)+4+idx*2 {
		return node.Pos
	}

	tOffset := flatbuffers.GetUint16(node.R(int(vPos) + 4 + idx*2))
	return node.Pos + int(tOffset)
}

// TableLen ... return table length in VTable.
func (node *Node) TableLen() int {
	voff := flatbuffers.GetUint32(node.R(node.Pos))
	vPos := node.Pos - int(voff)
	return int(flatbuffers.GetUint16(node.R(int(vPos) + 2)))
}

func (node *Node) VirtualTableLen() int {
	voff := flatbuffers.GetUint32(node.R(node.Pos))
	vPos := node.Pos - int(voff)
	return int(flatbuffers.GetUint16(node.R(int(vPos))))
}

func (node *Node) Table(idx int) int {
	pos := node.VirtualTable(idx)
	if node.VirtualTableIsZero(idx) {
		return -1
	}
	return pos + int(flatbuffers.GetUint32(node.R(pos)))
}

func (node *NodeList) clearValueInfoOnDirty() {
	node.ClearValueInfoOnDirty(node)
}

func (node *Node) BaseToNoLayer() {
	if _, already := node.Base.(NoLayer); already {
		return
	}
	node.Base = NewNoLayer(node.Base)
}
