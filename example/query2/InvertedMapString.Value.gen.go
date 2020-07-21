// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
genny must be called per Field
*/

var (
	InvertedMapString_Value_1 int = base.AtoiNoErr(Atoi("1"))
	InvertedMapString_Value   int = InvertedMapString_Value_1
)

// (field inedx, field type) -> InvertedMapString_IdxToType
var DUMMY_InvertedMapString_Value bool = SetInvertedMapStringFields("InvertedMapString", "Value", "Record", InvertedMapString_Value_1)

func (node InvertedMapString) Value() (result *Record) {
	result = emptyRecord()
	common := node.FieldAt(InvertedMapString_Value_1)
	if common.Node == nil {
		result = NewRecord()
		node.SetFieldAt(InvertedMapString_Value_1, result.SelfAsCommonNode())
		common = node.FieldAt(InvertedMapString_Value_1)
	}

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}

func (node InvertedMapString) SetValue(v *base.CommonNode) error {

	return node.CommonNode.SetFieldAt(1, v)
}
