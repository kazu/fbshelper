package query

import (
	"github.com/cheekybits/genny/generic"
	"github.com/kazu/fbshelper/query/base"
)

/*
genny must be called per Field
*/

type FieldName generic.Type

type FieldType generic.Type

var (
	NodeName_FieldName_FieldNum int = base.AtoiNoErr(Atoi("FieldNum"))
	NodeName_FieldName          int = NodeName_FieldName_FieldNum
)

// (field inedx, field type) ->  NodeName_IdxToType
var DUMMY_NodeName_FieldName bool = SetNodeNameFields("NodeName", "FieldName", "FieldType", NodeName_FieldName_FieldNum)

func (node NodeName) FieldName() (result *ResultType) {
	result = NewResultType()
	common := node.FieldAt(NodeName_FieldName_FieldNum)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}
