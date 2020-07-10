// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call per Field
go run github . com / cheekybits / genny gen "IndexString=Root Size=Version Int32=int64 CommonNode=base.CommonNode 0=0 False=false" ;
go run github . com / cheekybits / genny gen "IndexString=Root Size=Index   Int32=Index CommonNode=Index        0=1 False=false" ;
*/

// import (
// 	"github.com/cheekybits/genny/generic"
// 	b "github.com/kazu/fbshelper/query/base"
// )

var (
	IndexString_Size_0 int = base.AtoiNoErr(Atoi("0"))
	IndexString_Size   int = IndexString_Size_0
)

// (field inedx, field type) -> IndexString_IdxToType
var DUMMY_IndexString_Size bool = SetIndexStringFields("IndexString", "Size", "Int32", IndexString_Size_0)

// IndexStringSetNodeToIdx(IndexString_Size, base.NameToTypeEnum("Int32")) &&
// IndexStringSetIdxToName("Int32", IndexString_Size)

func (node IndexString) Size() (result *CommonNode) {
	result = NewCommonNode()
	common := node.FieldAt(IndexString_Size_0)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}
