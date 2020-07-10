// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call per Field
go run github . com / cheekybits / genny gen "Root=Root Version=Version Int32=int64 CommonNode=base.CommonNode 0=0 False=false" ;
go run github . com / cheekybits / genny gen "Root=Root Version=Index   Int32=Index CommonNode=Index        0=1 False=false" ;
*/

// import (
// 	"github.com/cheekybits/genny/generic"
// 	b "github.com/kazu/fbshelper/query/base"
// )

var (
	Root_Version_0 int = base.AtoiNoErr(Atoi("0"))
	Root_Version   int = Root_Version_0
)

// (field inedx, field type) -> Root_IdxToType
var DUMMY_Root_Version bool = SetRootFields("Root", "Version", "Int32", Root_Version_0)

// RootSetNodeToIdx(Root_Version, base.NameToTypeEnum("Int32")) &&
// RootSetIdxToName("Int32", Root_Version)

func (node Root) Version() (result *CommonNode) {
	result = NewCommonNode()
	common := node.FieldAt(Root_Version_0)

	result.Name = common.Name
	result.NodeList = common.NodeList
	result.IdxToType = common.IdxToType
	result.IdxToTypeGroup = common.IdxToTypeGroup

	return
}