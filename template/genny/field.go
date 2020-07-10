package query

import (
	"github.com/cheekybits/genny/generic"
	"github.com/kazu/fbshelper/query/base"
)

/*
must call per Field
	go run github.com/cheekybits/genny gen "NodeName=Root FieldName=Version FieldType=int64 ResultType=base.CommonNode FieldNum=0 IsStruct=false"
	go run github.com/cheekybits/genny gen "NodeName=Root FieldName=Index   FieldType=Index ResultType=Index        FieldNum=1 IsStruct=false"
*/

// import (
// 	"github.com/cheekybits/genny/generic"
// 	b "github.com/kazu/fbshelper/query/base"
// )

type FieldName generic.Type

//type FieldNum generic.Number
type FieldType generic.Type

//type ResultType generic.Type

var (
	NodeName_FieldName_FieldNum int = base.AtoiNoErr(Atoi("FieldNum"))
	NodeName_FieldName          int = NodeName_FieldName_FieldNum
)

// (field inedx, field type) ->  NodeName_IdxToType
var DUMMY_NodeName_FieldName bool = SetNodeNameFields("NodeName", "FieldName", "FieldType", NodeName_FieldName_FieldNum)

// NodeNameSetNodeToIdx(NodeName_FieldName, base.NameToTypeEnum("FieldType")) &&
// 							NodeNameSetIdxToName("FieldType", NodeName_FieldName)

func (node NodeName) FieldName() (result *ResultType) {
	result = NewResultType()
	common := node.FieldAt(NodeName_FieldName_FieldNum)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}
