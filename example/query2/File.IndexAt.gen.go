// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
genny must be called per Field
*/

var (
	File_IndexAt_2 int = base.AtoiNoErr(Atoi("2"))
	File_IndexAt   int = File_IndexAt_2
)

// (field inedx, field type) -> File_IdxToType
var DUMMY_File_IndexAt bool = SetFileFields("File", "IndexAt", "Int64", File_IndexAt_2)

func (node File) IndexAt() (result *CommonNode) {
	result = NewCommonNode()
	common := node.FieldAt(File_IndexAt_2)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}
