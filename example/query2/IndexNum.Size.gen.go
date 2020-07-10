// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call per Field
go run github . com / cheekybits / genny gen "IndexNum=Root Size=Version Int32=int64 CommonNode=base.CommonNode 0=0 False=false" ;
go run github . com / cheekybits / genny gen "IndexNum=Root Size=Index   Int32=Index CommonNode=Index        0=1 False=false" ;
*/

// import (
// 	"github.com/cheekybits/genny/generic"
// 	b "github.com/kazu/fbshelper/query/base"
// )

var (
	IndexNum_Size_0 int = base.AtoiNoErr(Atoi("0"))
	IndexNum_Size   int = IndexNum_Size_0
)

// (field inedx, field type) -> IndexNum_IdxToType
var DUMMY_IndexNum_Size bool = SetIndexNumFields("IndexNum", "Size", "Int32", IndexNum_Size_0)

// IndexNumSetNodeToIdx(IndexNum_Size, base.NameToTypeEnum("Int32")) &&
// IndexNumSetIdxToName("Int32", IndexNum_Size)

func (node IndexNum) Size() (result *CommonNode) {
	result = NewCommonNode()
	common := node.FieldAt(IndexNum_Size_0)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}