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

func NewCommonNode() *base.CommonNode {
	return &base.CommonNode{}
}

func Open(r io.Reader, cap int) RootType {
	return RootType{CommonNode: base.Open(r, cap)}
}

func OpenByBuf(buf []byte) RootType {
	return RootType{CommonNode: base.OpenByBuf(buf)}
}

func (node RootType) Next() RootType {
	start := node.Len()

	if node.LenBuf()+4 < start {
		return node
	}

	newBase := node.Base.NextBase(start)

	root := RootType{}
	root.Node = base.NewNode(newBase, int(flatbuffers.GetUOffsetT(newBase.R(0))))
	return root
}

func (node RootType) HasNext() bool {

	return node.LenBuf()+4 < node.Len()
}

func Atoi(s string) (int, error) {

	n, err := strconv.Atoi(s)
	if err != nil && s == "FieldNum" {
		err = nil
	}
	return n, err
}
