// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call per Field
go run github . com / cheekybits / genny gen "Record=Root FileId=Version Uint64=int64 CommonNode=base.CommonNode 0=0 False=false" ;
go run github . com / cheekybits / genny gen "Record=Root FileId=Index   Uint64=Index CommonNode=Index        0=1 False=false" ;
*/

// import (
// 	"github.com/cheekybits/genny/generic"
// 	b "github.com/kazu/fbshelper/query/base"
// )

var (
	Record_FileId_0 int = base.AtoiNoErr(Atoi("0"))
	Record_FileId   int = Record_FileId_0
)

// (field inedx, field type) -> Record_IdxToType
var DUMMY_Record_FileId bool = SetRecordFields("Record", "FileId", "Uint64", Record_FileId_0)

// RecordSetNodeToIdx(Record_FileId, base.NameToTypeEnum("Uint64")) &&
// RecordSetIdxToName("Uint64", Record_FileId)

func (node Record) FileId() (result *CommonNode) {
	result = NewCommonNode()
	common := node.FieldAt(Record_FileId_0)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}
