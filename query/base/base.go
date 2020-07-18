package base

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/kazu/loncha"

	flatbuffers "github.com/google/flatbuffers/go"
	log "github.com/kazu/fbshelper/query/log"
)

const (
	// DEFAULT_BUF_CAP .. cap size of Base's buffef([]byte)
	DEFAULT_BUF_CAP = 512
	// DEFAULT_NODE_BUF_CAP ... cap size of Base's buffer via CreateNode
	DEFAULT_NODE_BUF_CAP = 64
)

var (
	// tag for Marshal/Unmarshal
	FBS_TAG_NAME string = "fbs"
)

var (
	ERR_MUST_POINTER    error = log.ERR_MUST_POINTER
	ERR_INVALID_TYPE    error = log.ERR_INVALID_TYPE
	ERR_NOT_FOUND       error = log.ERR_NOT_FOUND
	ERR_READ_BUFFER     error = log.ERR_READ_BUFFER
	ERR_MORE_BUFFER     error = log.ERR_MORE_BUFFER
	ERR_NO_SUPPORT      error = log.ERR_NO_SUPPORT
	ERR_INVALID_INDEX   error = log.ERR_INVALID_INDEX
	ERR_NO_INCLUDE_ROOT error = log.ERR_NO_INCLUDE_ROOT
)

type LogLevel = log.LogLevel

var CurrentLogLevel LogLevel
var LogW io.Writer = os.Stderr

const (
	LOG_ERROR LogLevel = iota
	LOG_WARN
	LOG_DEBUG
)

type LogArgs struct {
	Fmt  string
	Infs []interface{}
}

type LogFn func() LogArgs

func SetLogLevel(l LogLevel) {
	CurrentLogLevel = l
}

func F(s string, v ...interface{}) LogArgs {
	return LogArgs{Fmt: s, Infs: v}
}

// if no output , not eval args
//  Log(LOG_DEBUG, func() LogArgs { return F("test %d \n", 1) })
func Log(l LogLevel, fn LogFn) {

	if CurrentLogLevel < l {
		return
	}

	var b strings.Builder
	switch l {
	case LOG_DEBUG:
		b.WriteString("D: ")
	case LOG_WARN:
		b.WriteString("W: ")
	case LOG_ERROR:
		b.WriteString("E: ")
	default:
		b.WriteString(" ")
	}
	args := fn()
	fmt.Fprintf(&b, args.Fmt, args.Infs...)
	io.WriteString(LogW, b.String())

	return

}

// Base is Base Object of byte buffer for flatbuffers
// read from r and store bytes.
// Diffs has jounals for writing
type Base struct {
	r       io.Reader
	bytes   []byte
	RDiffs  []Diff
	Diffs   []Diff
	dirties []Dirty
}

type Dirty struct {
	Pos int
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

func IsMatchBit(i, j int) bool {
	if (i & j) > 0 {
		return true
	}
	return false
}

// NewNode ... this provide creation of Node
// Node share buffer in same tree.
//    Base is buffer
//    pos is start position in Base's buffer
func NewNode(b *Base, pos int) *Node {
	return NewNode2(b, pos, false)
}

// NewNode2 ...  provide skip to initialize Vtable
func NewNode2(b *Base, pos int, noLoadVTable bool) *Node {
	node := &Node{Base: b, Pos: pos, Size: -1}
	if !noLoadVTable {
		node.preLoadVtable()
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
	if cap(b.bytes) > skip {
		b.bytes = b.bytes[:skip]
	}

	if len(b.bytes) <= skip && len(b.RDiffs) > 0 {
		//FIXME: impelement
		Log(LOG_WARN, func() LogArgs {
			return F("NextBase(): require RDiff. but not implemented\n")
		})
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
	if len(b.Diffs) == 0 {
		return b.readerR(off)
	}

	n, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})
	if e == nil && n >= 0 {
		return b.Diffs[n].bytes[(off - b.Diffs[n].Offset):]
	}

	return b.bytes[off:]
}

func (b *Base) readerR(off int) []byte {

	n, e := loncha.LastIndexOf(b.RDiffs, func(i int) bool {
		return b.RDiffs[i].Offset <= off && off < (b.RDiffs[i].Offset+len(b.RDiffs[i].bytes))
	})
	if e == nil && n >= 0 {
		return b.RDiffs[n].bytes[(off - b.RDiffs[n].Offset):]
	}

	if off+32 >= len(b.bytes) {
		b.expandBuf(off - len(b.bytes) + 32)
	}

	return b.bytes[off:]
}

func (b *Base) C(off, size int, src []byte) error {

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})
	if e == nil && sn >= 0 && b.Diffs[sn].Offset == off && len(b.Diffs[sn].bytes) == size {
		b.Diffs[sn].bytes = src[:size]
		return nil
	}

	diff := b.D(off, size)
	diff.bytes = src[:size]
	return nil
}

func (b *Base) Copy(src *Base, srcOff, size, dstOff, extend int) {

	if len(b.bytes) > dstOff {
		diff := Diff{Offset: dstOff, bytes: b.bytes[dstOff:]}
		b.Diffs = append(b.Diffs, diff)
		b.bytes = b.bytes[:dstOff]
	}

	for i, diff := range b.Diffs {
		if diff.Offset >= dstOff {
			diff.Offset += extend
		}
		b.Diffs[i] = diff
	}

	if len(src.bytes) > srcOff {
		nSize := len(src.bytes[srcOff:])
		if nSize > size {
			nSize = size
		}

		diff := Diff{Offset: dstOff, bytes: src.bytes[srcOff : srcOff+nSize]}
		b.Diffs = append(b.Diffs, diff)
	}
	for _, diff := range src.Diffs {
		if diff.Offset >= srcOff {
			nDiff := diff
			nDiff.Offset -= srcOff
			nDiff.Offset += dstOff
			b.Diffs = append(b.Diffs, nDiff)
		}
	}
	return
}

func (b *Base) D(off, size int) *Diff {

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})

	if e == nil && sn >= 0 {
		// if b.Diffs[sn].Offset+len(b.Diffs[sn].bytes) < off+size {
		// 	Log(LOG_ERROR, func() LogArgs {
		// 		return F("D(%d,%d) Invalid diff=%d\n",
		// 			off, size, )
		// 	})
		// }else
		if b.Diffs[sn].Offset+len(b.Diffs[sn].bytes) <= off+size {
			diff := Diff{Offset: off}
			b.Diffs = append(b.Diffs, diff)
			idx := len(b.Diffs) - 1
			return &b.Diffs[idx]
		}

		diffbefore := b.Diffs[sn]
		//MENTION: should increase cap ?
		diff := Diff{Offset: off}
		diffafter := Diff{Offset: off + size, bytes: diffbefore.bytes[off-diffbefore.Offset+size:]}
		diffbefore.bytes = diffbefore.bytes[:off+size-diffbefore.Offset]
		b.Diffs[sn] = diffbefore
		b.Diffs = append(b.Diffs[:sn+1],
			append([]Diff{diffafter}, b.Diffs[sn+1:]...)...)
		b.Diffs = append(b.Diffs, diff)
		return &b.Diffs[len(b.Diffs)-1]
	}

	if len(b.bytes) > off {
		diff := Diff{Offset: off, bytes: b.bytes[off:]}
		b.Diffs = append(b.Diffs, diff)
		b.bytes = b.bytes[:off]
	}

	//MENTION: should increase cap ?
	diff := Diff{Offset: off}
	b.Diffs = append(b.Diffs, diff)
	idx := len(b.Diffs) - 1

	return &b.Diffs[idx]
}

func (b *Base) U(off, size int) []byte {

	diff := b.D(off, size)
	diff.bytes = make([]byte, size)
	return diff.bytes
}

// FIXME:  support base.RDiff
func (b *Base) expandBuf(plus int) error {
	if !b.HasIoReader() && cap(b.bytes)-len(b.bytes) < plus {
		return nil
	}

	if !b.HasIoReader() {
		return nil
	}
	if cap(b.bytes) < len(b.bytes)+plus {
		diff := Diff{Offset: len(b.bytes), bytes: make([]byte, 0, DEFAULT_BUF_CAP)}
		n, err := io.ReadAtLeast(b.r, diff.bytes, plus)
		if n < plus || err != nil {
			diff.bytes = diff.bytes[:n]
			if n > 0 {
				b.RDiffs = append(b.RDiffs, diff)
			}
			return ERR_READ_BUFFER
		}
		b.RDiffs = append(b.RDiffs, diff)
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

	max := len(b.bytes)
	for i := range b.Diffs {
		if b.Diffs[i].Offset+len(b.Diffs[i].bytes) > max {
			max = b.Diffs[i].Offset + len(b.Diffs[i].bytes)
		}
	}

	return max

}
func (b *Base) Merge() {

	nbytes := make([]byte, b.LenBuf())

	copy(nbytes, b.bytes)
	for i := range b.Diffs {
		copy(nbytes[b.Diffs[i].Offset:], b.Diffs[i].bytes)
	}
	b.bytes = nbytes
	b.Diffs = []Diff{}

}

// RawBufInfo ... capacity/length infomation of Base's buffer
type RawBufInfo struct {
	Len int
	Cap int
}

// BufInfo is buffer detail infomation
// information of Base.bytes is stored to 0 indexed
// information of Base.Diffs is stored to 1 indexed
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
	n.preLoadVtable()
}

func (n *Node) preLoadVtable() {

	vOffset := int(flatbuffers.GetUOffsetT(n.R(n.Pos)))
	vPos := int(n.Pos) - vOffset
	vLen := int(flatbuffers.GetVOffsetT(n.R(vPos)))
	n.TLen = uint16(flatbuffers.GetVOffsetT(n.R(vPos + 2)))

	for cur := vPos + 4; cur < vPos+vLen; cur += 2 {
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
func (b *Base) insertBuf(pos, size int) *Base {
	return b.insertSpace(pos, size, true)
}
func (b *Base) insertSpace(pos, size int, isCreate bool) *Base {

	newBase := &Base{
		r:     b.r,
		bytes: b.bytes,
	}

	newBase.Diffs = make([]Diff, len(b.Diffs), cap(b.Diffs))

	copy(newBase.Diffs, b.Diffs)

	for i, diff := range newBase.Diffs {
		if diff.Offset < pos && diff.Offset+len(diff.bytes) > pos {
			newBase.Diffs = append(newBase.Diffs[:i],
				append([]Diff{
					Diff{Offset: diff.Offset, bytes: diff.bytes[:pos-diff.Offset]},
					Diff{Offset: pos, bytes: diff.bytes[pos-diff.Offset:]}},
					b.Diffs[i+1:]...)...)
		}
	}

	for i := range newBase.Diffs {
		if newBase.Diffs[i].Offset >= pos {
			newBase.Diffs[i].Offset += size
		}
	}

	if len(newBase.bytes) > pos {
		newBase.Diffs = append(newBase.Diffs,
			Diff{Offset: pos + size, bytes: newBase.bytes[pos:]})
		newBase.bytes = newBase.bytes[:pos]
	}
	if isCreate {
		newBase.Diffs = append(newBase.Diffs,
			Diff{Offset: pos, bytes: make([]byte, size)})
	}
	return newBase
}

func (node *Node) VirtualTableIsZero(idx int) bool {

	return node.VirtualTable(idx) == node.Pos
}

func (node *Node) VirtualTable(idx int) int {

	voff := flatbuffers.GetUint32(node.R(node.Pos))
	vPos := node.Pos - int(voff)
	tOffset := flatbuffers.GetUint16(node.R(int(vPos) + 4 + idx*2))
	return node.Pos + int(tOffset)
}

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
	return pos + int(flatbuffers.GetUint32(node.R(pos)))
}
