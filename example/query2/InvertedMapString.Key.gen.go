// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
genny must be called per Field
*/

var (
	InvertedMapString_Key_0 int = base.AtoiNoErr(Atoi("0"))
	InvertedMapString_Key   int = InvertedMapString_Key_0
)

// (field inedx, field type) -> InvertedMapString_IdxToType
var DUMMY_InvertedMapString_Key bool = SetInvertedMapStringFields("InvertedMapString", "Key", "[]byte", InvertedMapString_Key_0)

func (node InvertedMapString) Key() (result *CommonNode) {
	result = NewCommonNode()
	common := node.FieldAt(InvertedMapString_Key_0)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}
