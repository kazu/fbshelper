package query

import (
	"strconv"

	"github.com/cheekybits/genny/generic"
	b "github.com/kazu/fbshelper/query/base"
)

/*
must call per Field
　genny gen "NodeName=Root FieldName=Version FieldType=int64 FieldNum=0 IsStruct=false"
　genny gen "NodeName=Root FieldName=Index FieldType=Index FieldNum=1 IsStruct=false"
*/

type FieldName generic.Type

//type FieldNum generic.Number
type FieldType generic.Type

var (
	FieldName_FieldNum int = b.AtoiNoErr(strconv.Atoi("FieldNum"))
	NodeName_FieldName int = FieldName_FieldNum
)

// (field inedx, field type) ->  NodeName_IdxToType
var DUMMY_NodeName_FieldName bool = SetNodeToIdx(NodeName_FieldName, b.NameToTypeEnum("FieldType"))
