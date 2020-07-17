// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
genny must be called per Field
*/

var (
	Root_Record_3 int = base.AtoiNoErr(Atoi("3"))
	Root_Record   int = Root_Record_3
)

// (field inedx, field type) -> Root_IdxToType
var DUMMY_Root_Record bool = SetRootFields("Root", "Record", "Record", Root_Record_3)

func (node Root) Record() (result *Record) {
	result = emptyRecord()
	common := node.FieldAt(Root_Record_3)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}

func (node Root) SetRecord(v *base.CommonNode) error {

	return node.CommonNode.SetFieldAt(3, v)
}
