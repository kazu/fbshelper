package base

import (
	"errors"
	"reflect"

	"github.com/kazu/loncha"

	flatbuffers "github.com/google/flatbuffers/go"
)

const (
	SizeOfbool        = 1
	SizeOfSizeOfint8  = 1
	SizeOfSizeOfint16 = 2
	SizeOfuint16      = 2
	SizeOfint32       = 4
	SizeOfuint32      = 4
	SizeOfint64       = 8
	SizeOfuint64      = 8
	SizeOffloat32     = 4
	SizeOffloat64     = 8
	SizeOfuint8       = 1
	SizeOfbyte        = 1
)

var (
	ERR_MUST_POINTER error = errors.New("parameter must be pointer")
	ERR_INVALID_TYPE error = errors.New("parameter invalid type(must be struct or map[string]interface)")
	ERR_NOT_FOUND    error = errors.New("data is not found")
)

var (
	FBS_TAG_NAME string = "fbs"
)

type Base struct {
	Bytes []byte
	Diffs []Diff
}

type Diff struct {
	Offset int
	Bytes  []byte
}

type Root struct {
	*Node
}

type Node struct {
	*Base
	Pos        int
	Size       int
	VTable     []uint16
	TLen       uint16
	ValueInfos []ValueInfo
}

type NodeList struct {
	*Node
	ValueInfo
}

type ValueInfo struct {
	Pos  int
	Size int
	VLen uint32
}

type Info ValueInfo

func NewNode(b *Base, pos int) *Node {
	return NewNode2(b, pos, false)
}

func NewNode2(b *Base, pos int, noVTable bool) *Node {
	node := &Node{Base: b, Pos: pos, Size: -1}
	if !noVTable {
		node.vtable()
	}
	return node
}

func (b *Base) R(int off) []byte {

	n, e := loncha.IndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].Bytes))
	})
	if e == nil || n >= 0 {
		return b.Diffs[i].Bytes[(off - b.Diffs[i].Offset):]
	}
	return b.Bytes[off:]
}

func (n *Node) vtable() {
	if len(n.VTable) > 0 {
		return
	}
	buf := n.Bytes
	vOffset := uint32(flatbuffers.GetUOffsetT(buf[n.Pos:]))
	vPos := uint32(n.Pos) - vOffset
	vLen := uint16(flatbuffers.GetVOffsetT(buf[vPos:]))
	n.TLen = uint16(flatbuffers.GetVOffsetT(buf[vPos+2:]))

	for cur := vPos + 4; cur < vPos+uint32(vLen); cur += 2 {
		n.VTable = append(n.VTable, uint16(flatbuffers.GetVOffsetT(buf[cur:])))
	}
	n.ValueInfos = make([]ValueInfo, len(n.VTable))
}

func FbsString(node *Node) []byte {
	buf := node.Bytes
	pos := uint32(node.Pos + int(node.VTable[0]))
	sLenOff := flatbuffers.GetUint32(buf[pos:])
	sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])
	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	return buf[start : start+sLen]
}

func FbsStringInfo(node *Node) Info {

	buf := node.Bytes
	pos := uint32(node.Pos + int(node.VTable[0]))
	sLenOff := flatbuffers.GetUint32(buf[pos:])
	sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])
	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	return Info{Pos: int(start), Size: int(sLen)}
}

func (info ValueInfo) IsNotReady() bool {
	return info.Pos < 1

}

func (node *Node) ValueInfoPos(vIdx int) ValueInfo {
	if node.VTable[vIdx] == 0 {
		node.ValueInfos[vIdx].Pos = -1
		node.ValueInfos[vIdx].Size = -1
		return node.ValueInfos[vIdx]
	}
	node.ValueInfos[vIdx].Pos = node.Pos + int(node.VTable[vIdx])
	return node.ValueInfos[vIdx]
}

func (node *Node) ValueInfoPosBytes(vIdx int) ValueInfo {
	if node.VTable[vIdx] == 0 {
		node.ValueInfos[vIdx].Pos = -1
		node.ValueInfos[vIdx].Size = -1
		return node.ValueInfos[vIdx]
	}
	buf := node.Bytes
	pos := uint32(node.Pos + int(node.VTable[vIdx]))
	sLenOff := flatbuffers.GetUint32(buf[pos:])
	sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])
	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	node.ValueInfos[vIdx].Pos = int(start)
	node.ValueInfos[vIdx].Size = int(sLen)
	return node.ValueInfos[vIdx]
}

func (node *Node) ValueInfoPosTable(vIdx int) ValueInfo {
	if node.VTable[vIdx] == 0 {
		node.ValueInfos[vIdx].Pos = -1
		node.ValueInfos[vIdx].Size = -1
		return node.ValueInfos[vIdx]
	}

	pos := node.Pos + int(node.VTable[vIdx])
	start := int(flatbuffers.GetUint32(node.Bytes[pos:])) + pos
	node.ValueInfos[vIdx].Pos = start

	return node.ValueInfos[vIdx]
}

func (node *Node) ValueInfoPosList(vIdx int) ValueInfo {
	if node.VTable[vIdx] == 0 {
		node.ValueInfos[vIdx].Pos = -1
		node.ValueInfos[vIdx].Size = -1
		return node.ValueInfos[vIdx]
	}
	buf := node.Bytes
	vPos := uint32(node.Pos + int(node.VTable[vIdx]))
	vLenOff := flatbuffers.GetUint32(buf[vPos:])
	vLen := flatbuffers.GetUint32(buf[vPos+vLenOff:])
	start := vPos + vLenOff + flatbuffers.SizeUOffsetT

	node.ValueInfos[vIdx].Pos = int(start)
	node.ValueInfos[vIdx].VLen = vLen

	return node.ValueInfos[vIdx]

}

func (node *Node) ValueNormal(vIdx int) []byte {
	if node.ValueInfos[vIdx].Pos < 1 {
		node.ValueInfoPos(vIdx)
	}
	return node.Bytes[node.ValueInfos[vIdx].Pos:]
}

func (node *Node) ValueBytes(vIdx int) []byte {
	if node.ValueInfos[vIdx].Pos < 1 {
		node.ValueInfoPosBytes(vIdx)
	}
	valInfo := node.ValueInfos[vIdx]
	return node.Bytes[valInfo.Pos : valInfo.Pos+valInfo.Size]

}

func (node *Node) ValueTable(vIdx int) *Node {
	if node.ValueInfos[vIdx].Pos < 1 {
		node.ValueInfoPosTable(vIdx)
	}

	return NewNode(node.Base, node.ValueInfos[vIdx].Pos)
}

func (node *Node) ValueStruct(vIdx int) *Node {
	if node.ValueInfos[vIdx].Pos < 1 {
		node.ValueInfoPos(vIdx)
	}

	return NewNode2(node.Base, node.ValueInfos[vIdx].Pos, true)
}

func (node *Node) ValueList(vIdx int) NodeList {

	if node.ValueInfos[vIdx].Pos < 1 {
		node.ValueInfoPosList(vIdx)
	}

	return NodeList{Node: NewNode(node.Base, node.Pos),
		ValueInfo: node.ValueInfos[vIdx]}
}

type UnmarshalFn func(string, reflect.Value) error

func (node *Node) Unmarshal(ptr interface{}, setter UnmarshalFn) error {

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
