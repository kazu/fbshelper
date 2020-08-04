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
	result = emptyResultType()
	common := node.FieldAt(NodeName_FieldName_FieldNum)
	if common.Node == nil {
		result = NewResultType()
		node.SetFieldAt(NodeName_FieldName_FieldNum, result.SelfAsCommonNode())
		common = node.FieldAt(NodeName_FieldName_FieldNum)
	}

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}

//func (node NodeName) SetFieldName(v *base.CommonNode) error {
func (node NodeName) SetFieldName(v *ResultType) error {

	return node.CommonNode.SetFieldAt(FieldNum, v.SelfAsCommonNode())
}
