package base

import (
	"io"
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

const (
	// this is cap size of Base's buffef([]byte)
	DEFAULT_BUF_CAP = 512
)

var (
	// tag for Marshal/Unmarshal
	FBS_TAG_NAME string = "fbs"
)

// Base is Base Object of byte buffer for flatbuffers
// read from r and store bytes.
// Diffs has jounals for writing
type Base struct {
	r     io.Reader
	bytes []byte
	Diffs []Diff
}

// Diff is journal for create/update
type Diff struct {
	Offset int
	bytes  []byte
}

// Root manage top node of data.
type Root struct {
	*Node
}

// Node base struct.
// Pos is start position of buffer
// ValueInfos is information of pointted fields
type Node struct {
	*Base
	Pos        int
	Size       int
	VTable     []uint16
	TLen       uint16
	ValueInfos []ValueInfo
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
func NewNode(b *Base, pos int) *Node {
	return NewNode2(b, pos, false)
}

// NewNode2 ...  provide skip to initialize Vtable
func NewNode2(b *Base, pos int, noVTable bool) *Node {
	node := &Node{Base: b, Pos: pos, Size: -1}
	if !noVTable {
		node.vtable()
	}
	return node
}

// NewBase initialize Base struct via buffer(buf)
func NewBase(buf []byte) *Base {
	return &Base{bytes: buf}
}

// NewBaseByIO initialize , make Base
//   this dosent use []byte, use io.Reader
func NewBaseByIO(rio io.Reader, cap int) *Base {
	b := &Base{r: rio, bytes: make([]byte, 0, cap)}
	return b
}

// NextBase provide next root flatbuffers
// this is mainly for streaming data.
func (b *Base) NextBase(skip int) *Base {
	newBase := &Base{
		r:     b.r,
		bytes: b.bytes,
	}
	newBase.bytes = newBase.bytes[skip:]
	if cap(newBase.bytes) < DEFAULT_BUF_CAP {
		newBase.Diffs = append(newBase.Diffs, Diff{Offset: cap(newBase.bytes), bytes: make([]byte, 0, DEFAULT_BUF_CAP)})
	}
	if cap(b.bytes) > skip {
		b.bytes = b.bytes[:skip]
	}

	return newBase
}

// HasIoReader ... Base buffer is reading io.Reader or not.
func (b *Base) HasIoReader() bool {
	return b.r != nil
}

// R ... is access buffer data
// Base cannot access byte buffer directly.
func (b *Base) R(off int) []byte {

	n, e := loncha.IndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})
	if e == nil && n >= 0 {
		return b.Diffs[n].bytes[(off - b.Diffs[n].Offset):]
	}
	if off+32 >= len(b.bytes) {
		b.expandBuf(off - len(b.bytes) + 32)
	}

	return b.bytes[off:]
}

func (b *Base) expandBuf(plus int) error {
	if !b.HasIoReader() {
		return nil
	}
	l := len(b.bytes)
	b.bytes = b.bytes[:l+plus]
	n, err := io.ReadAtLeast(b.r, b.bytes[l:], plus)
	if n < plus || err != nil {
		b.bytes = b.bytes[:l+n]
		return ERR_READ_BUFFER
	}
	return nil
}

// LenBuf ... size of buffer
func (b *Base) LenBuf() int {

	if len(b.Diffs) < 1 {
		return len(b.bytes)
	}
	if len(b.bytes) < b.Diffs[len(b.Diffs)-1].Offset+len(b.Diffs[len(b.Diffs)-1].bytes) {
		return b.Diffs[len(b.Diffs)-1].Offset + len(b.Diffs[len(b.Diffs)-1].bytes)
	}
	return len(b.bytes)

}

// RawBufInfo ... capacity/length infomation of Base's buffer
type RawBufInfo struct {
	Len int
	Cap int
}

// BufInfo is buffer detail infomation
// information of Base.bytes is stored to 0 indexed
//  information of Base.Diffs is stored to 1 indexed
type BufInfo [2]RawBufInfo

// BufInfo ... return infos(buffer detail information)
func (b *Base) BufInfo() (infos BufInfo) {

	infos[0].Len = len(b.bytes)
	infos[0].Cap = cap(b.bytes)

	infos[1].Len = infos[0].Len
	infos[1].Cap = infos[0].Cap
	for _, diff := range b.Diffs {
		if len(diff.bytes)+diff.Offset > infos[1].Len {
			infos[1].Len = diff.Offset + len(diff.bytes)
		}
		if cap(diff.bytes)+diff.Offset > infos[1].Cap {
			infos[1].Cap = diff.Offset + cap(diff.bytes)
		}
	}
	return
}

func (n *Node) vtable() {
	if len(n.VTable) > 0 {
		return
	}
	vOffset := int(flatbuffers.GetUOffsetT(n.R(n.Pos)))
	vPos := int(n.Pos) - vOffset
	vLen := int(flatbuffers.GetVOffsetT(n.R(vPos)))
	n.TLen = uint16(flatbuffers.GetVOffsetT(n.R(vPos + 2)))

	for cur := vPos + 4; cur < vPos+vLen; cur += 2 {
		n.VTable = append(n.VTable, uint16(flatbuffers.GetVOffsetT(n.R(cur))))
	}
	n.ValueInfos = make([]ValueInfo, len(n.VTable))
}

// FbsString ... return []bytes for string data
func FbsString(node *Node) []byte {
	pos := node.Pos + int(node.VTable[0])
	sLenOff := int(flatbuffers.GetUint32(node.R(pos)))
	sLen := int(flatbuffers.GetUint32(node.R(pos + sLenOff)))
	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	return node.R(start)[:sLen]
}

// FbsStringInfo ... return Node Infomation for FbsString
func FbsStringInfo(node *Node) Info {

	pos := node.Pos + int(node.VTable[0])
	sLenOff := int(flatbuffers.GetUint32(node.R(pos)))
	sLen := flatbuffers.GetUint32(node.R(pos + sLenOff))
	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	return Info{Pos: start, Size: int(sLen)}
}

// IsNotReady ... already set ValueInfo (field information)
func (info ValueInfo) IsNotReady() bool {
	return info.Pos < 1

}

// ValueInfoPos ... fetching vtable position infomation
func (node *Node) ValueInfoPos(vIdx int) ValueInfo {
	if node.VTable[vIdx] == 0 {
		node.ValueInfos[vIdx].Pos = -1
		node.ValueInfos[vIdx].Size = -1
		return node.ValueInfos[vIdx]
	}
	node.ValueInfos[vIdx].Pos = node.Pos + int(node.VTable[vIdx])
	return node.ValueInfos[vIdx]
}

// ValueInfoPosBytes ... etching vtable position infomation for []byte
func (node *Node) ValueInfoPosBytes(vIdx int) ValueInfo {
	if node.VTable[vIdx] == 0 {
		node.ValueInfos[vIdx].Pos = -1
		node.ValueInfos[vIdx].Size = -1
		return node.ValueInfos[vIdx]
	}
	pos := node.Pos + int(node.VTable[vIdx])
	sLenOff := int(flatbuffers.GetUint32(node.R(pos)))
	sLen := flatbuffers.GetUint32(node.R(pos + sLenOff))

	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	node.ValueInfos[vIdx].Pos = start
	node.ValueInfos[vIdx].Size = int(sLen)
	return node.ValueInfos[vIdx]
}

// ValueInfoPosTable ... etching vtable position infomation for flatbuffers table
func (node *Node) ValueInfoPosTable(vIdx int) ValueInfo {
	if node.VTable[vIdx] == 0 {
		node.ValueInfos[vIdx].Pos = -1
		node.ValueInfos[vIdx].Size = -1
		return node.ValueInfos[vIdx]
	}

	pos := node.Pos + int(node.VTable[vIdx])
	start := int(flatbuffers.GetUint32(node.R(pos))) + pos

	node.ValueInfos[vIdx].Pos = start

	return node.ValueInfos[vIdx]
}

// ValueInfoPosList ... etching vtable position infomation for flatbuffers vector
func (node *Node) ValueInfoPosList(vIdx int) ValueInfo {
	if node.VTable[vIdx] == 0 {
		node.ValueInfos[vIdx].Pos = -1
		node.ValueInfos[vIdx].Size = -1
		return node.ValueInfos[vIdx]
	}
	vPos := node.Pos + int(node.VTable[vIdx])
	vLenOff := int(flatbuffers.GetUint32(node.R(vPos)))
	vLen := flatbuffers.GetUint32(node.R(vPos + vLenOff))
	start := vPos + vLenOff + flatbuffers.SizeUOffsetT

	node.ValueInfos[vIdx].Pos = int(start)
	node.ValueInfos[vIdx].VLen = vLen

	return node.ValueInfos[vIdx]

}

func (node *Node) ValueNormal(vIdx int) []byte {
	if node.ValueInfos[vIdx].Pos < 1 {
		node.ValueInfoPos(vIdx)
	}
	return node.R(node.ValueInfos[vIdx].Pos)
}

func (node *Node) ValueBytes(vIdx int) []byte {
	if node.ValueInfos[vIdx].Pos < 1 {
		node.ValueInfoPosBytes(vIdx)
	}
	valInfo := node.ValueInfos[vIdx]

	return node.R(valInfo.Pos)[:valInfo.Size]

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
