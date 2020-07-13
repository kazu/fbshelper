// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
genny must be called per Field
*/

var (
	IndexString_Maps_1 int = base.AtoiNoErr(Atoi("1"))
	IndexString_Maps   int = IndexString_Maps_1
)

// (field inedx, field type) -> IndexString_IdxToType
var DUMMY_IndexString_Maps bool = SetIndexStringFields("IndexString", "Maps", "[]InvertedMapString", IndexString_Maps_1)

func (node IndexString) Maps() (result *InvertedMapStringList) {
	result = emptyInvertedMapStringList()
	common := node.FieldAt(IndexString_Maps_1)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}
