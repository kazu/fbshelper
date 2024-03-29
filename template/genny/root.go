package query

import (
	"io"
	"strconv"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/kazu/fbshelper/query/base"
)

type Bool = base.Bool
type Byte = base.Byte
type Int8 = base.Int8
type Int16 = base.Int16
type Int32 = base.Int32
type Int64 = base.Int64
type Uint8 = base.Uint8
type Uint16 = base.Uint16
type Uint32 = base.Uint32
type Uint64 = base.Uint64
type Float32 = base.Float32
type Float64 = base.Float64

type CommonNode = base.CommonNode
type List = base.List

type CommonNodeList = base.CommonNode

func emptyCommonNode() *base.CommonNode {
	return &base.CommonNode{}
}

func NewCommonNode() *base.CommonNode {
	return emptyCommonNode()
}

func emptyCommonNodeList() *CommonNodeList {
	return &CommonNodeList{}
}

func NewCommonNodeList() *CommonNodeList {
	return emptyCommonNodeList()
}

func emptyList() *base.List {
	return (*base.List)(&base.CommonNode{})
}

func NewList() *List {
	return emptyList()
}

func Open(r io.Reader, cap int, opts ...base.Option) RootType {
	result := RootType{CommonNode: base.Open(r, cap, opts...)}
	result.CommonNode.Name = "RootType"
	base.SetRootName(result.CommonNode.Name)
	result.FetchIndex()
	return result
}

func OpenByBuf(buf []byte, opts ...base.Option) RootType {
	result := RootType{CommonNode: base.OpenByBuf(buf, opts...)}
	result.CommonNode.Name = "RootType"
	base.SetRootName(result.CommonNode.Name)
	result.FetchIndex()
	return result
}

func toRoot(b base.IO) RootType {
	common := &CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Node = base.NewNode(b, int(flatbuffers.GetUOffsetT(b.R(0, base.Size(4)))))
	root := RootType{CommonNode: common}
	root.CommonNode.Name = "RootType"
	root.FetchIndex()
	return root
}

func RootFromCommon(c *base.CommonNode) RootType {
	return RootType{
		CommonNode: c.RootCommon(),
	}
}

func (node RootType) Next() RootType {
	start := node.Len()

	if node.LenBuf() < start {
		return node
	}

	newBase := node.IO.Next(start)

	root := RootType{}
	root.CommonNode = emptyCommonNode()
	root.NodeList = &base.NodeList{}
	root.Node = base.NewNode(newBase, int(flatbuffers.GetUOffsetT(newBase.R(0, base.Size(4)))))
	return root
}

func (node RootType) InsertBuf(pos, size int) {
	node.CommonNode.InsertBuf(pos, size)
}

func (node RootType) HasNext() bool {

	return node.LenBuf()+4 < node.Len()
}

func (node RootType) WithHeader() RootType {
	if !node.InRoot() {
		header := make([]byte, 8)
		pos := node.Node.Pos + 8
		flatbuffers.WriteUint32(header, uint32(pos))
		node.Copy(base.NewBase(header), 0, 8, 0, 8)

		node.Node.Pos += 8
	}
	return node
}

func Atoi(s string) (int, error) {

	n, err := strconv.Atoi(s)
	if err != nil && s == "FieldNum" {
		err = nil
	}
	return n, err
}

func FromBool(v bool) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Bool"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfbool)), 0, true)
	common.Node.Size = base.SizeOfbool
	common.SetBool(v)
	return common
}

func FromByte(v byte) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Byte"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfbyte)), 0, true)
	common.Node.Size = base.SizeOfbyte
	common.SetByte(v)
	return common
}

func FromInt8(v int8) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Int8"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfint8)), 0, true)
	common.Node.Size = base.SizeOfint8
	common.SetInt8(v)
	return common
}

func FromInt16(v int16) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Int16"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfint16)), 0, true)
	common.Node.Size = base.SizeOfint16
	common.SetInt16(v)
	return common
}

func FromInt32(v int32) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Int32"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfint32)), 0, true)
	common.Node.Size = base.SizeOfint32
	common.SetInt32(v)
	return common
}

func FromInt64(v int64) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Int64"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfint64)), 0, true)
	common.Node.Size = base.SizeOfint64
	common.SetInt64(v)
	return common
}

func FromUint8(v uint8) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Uint8"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfuint8)), 0, true)
	common.Node.Size = base.SizeOfuint8
	common.SetUint8(v)
	return common
}

func FromUint16(v uint16) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Uint16"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfuint16)), 0, true)
	common.Node.Size = base.SizeOfuint16
	common.SetUint16(v)
	return common
}

func FromUint32(v uint32) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Uint32"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfuint32)), 0, true)
	common.Node.Size = base.SizeOfuint32
	common.SetUint32(v)
	return common
}

func FromUint64(v uint64) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Uint64"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOfuint64)), 0, true)
	common.Node.Size = base.SizeOfuint64
	common.SetUint64(v)
	return common
}

func FromFloat32(v float32) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Float32"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOffloat32)), 0, true)
	common.Node.Size = base.SizeOffloat32
	common.SetFloat32(v)
	return common
}

func FromFloat64(v float64) *CommonNode {

	common := &base.CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Name = "Float64"
	common.Node = base.NewNode2(base.NewBase(make([]byte, base.SizeOffloat64)), 0, true)
	common.Node.Size = base.SizeOffloat64
	common.SetFloat64(v)
	return common
}
