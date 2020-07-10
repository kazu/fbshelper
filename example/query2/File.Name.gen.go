// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
genny must be called per Field
*/

var (
	File_Name_1 int = base.AtoiNoErr(Atoi("1"))
	File_Name   int = File_Name_1
)

// (field inedx, field type) -> File_IdxToType
var DUMMY_File_Name bool = SetFileFields("File", "Name", "[]byte", File_Name_1)

func (node File) Name() (result *CommonNode) {
	result = NewCommonNode()
	common := node.FieldAt(File_Name_1)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}
