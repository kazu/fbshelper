package query

// go run github.com/cheekybits/genny gen "RootType=Root "
import (
	"io"
	"strconv"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/kazu/fbshelper/query/base"
)

// import (
// 	"io"
// 	"strconv"

// 	flatbuffers "github.com/google/flatbuffers/go"
// 	"github.com/kazu/fbshelper/query/base"
// 	b "github.com/kazu/fbshelper/query/base"
// )

//type RootType geneic.Type

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
