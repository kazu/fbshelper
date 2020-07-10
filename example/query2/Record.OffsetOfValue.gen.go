// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call per Field
go run github . com / cheekybits / genny gen "Record=Root OffsetOfValue=Version Int32=int64 CommonNode=base.CommonNode 3=0 False=false" ;
go run github . com / cheekybits / genny gen "Record=Root OffsetOfValue=Index   Int32=Index CommonNode=Index        3=1 False=false" ;
*/

// import (
// 	"github.com/cheekybits/genny/generic"
// 	b "github.com/kazu/fbshelper/query/base"
// )

var (
	Record_OffsetOfValue_3 int = base.AtoiNoErr(Atoi("3"))
	Record_OffsetOfValue   int = Record_OffsetOfValue_3
)

// (field inedx, field type) -> Record_IdxToType
var DUMMY_Record_OffsetOfValue bool = SetRecordFields("Record", "OffsetOfValue", "Int32", Record_OffsetOfValue_3)

// RecordSetNodeToIdx(Record_OffsetOfValue, base.NameToTypeEnum("Int32")) &&
// RecordSetIdxToName("Int32", Record_OffsetOfValue)

func (node Record) OffsetOfValue() (result *CommonNode) {
	result = NewCommonNode()
	common := node.FieldAt(Record_OffsetOfValue_3)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}