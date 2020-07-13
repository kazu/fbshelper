// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

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

type CommonNodeList = base.CommonNode

func emptyCommonNode() *base.CommonNode {
	return &base.CommonNode{}
}

func emptyCommonNodeList() *CommonNodeList {
	return &CommonNodeList{}
}

func Open(r io.Reader, cap int) Root {
	result := Root{CommonNode: base.Open(r, cap)}
	result.CommonNode.Name = "Root"
	base.SetRootName(result.CommonNode.Name)
	result.FetchIndex()
	return result
	//return Root{CommonNode: base.Open(r, cap)}
}

func OpenByBuf(buf []byte) Root {
	result := Root{CommonNode: base.OpenByBuf(buf)}
	result.CommonNode.Name = "Root"
	base.SetRootName(result.CommonNode.Name)
	result.FetchIndex()
	return result
}

func toRoot(b *base.Base) Root {
	common := &CommonNode{}
	common.NodeList = &base.NodeList{}
	common.Node = base.NewNode(b, int(flatbuffers.GetUOffsetT(b.R(0))))
	root := Root{CommonNode: common}
	root.CommonNode.Name = "Root"
	root.FetchIndex()
	return root
}

func (node Root) Next() Root {
	start := node.Len()

	if node.LenBuf()+4 < start {
		return node
	}

	newBase := node.Base.NextBase(start)

	root := Root{}
	root.CommonNode = emptyCommonNode()
	root.NodeList = &base.NodeList{}
	root.Node = base.NewNode(newBase, int(flatbuffers.GetUOffsetT(newBase.R(0))))
	return root
}

func (node Root) InsertBuf(pos, size int) {
	node.CommonNode.InsertBuf(pos, size)
}

func (node Root) HasNext() bool {

	return node.LenBuf()+4 < node.Len()
}

func Atoi(s string) (int, error) {

	n, err := strconv.Atoi(s)
	if err != nil && s == "FieldNum" {
		err = nil
	}
	return n, err
}
