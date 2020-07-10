package query

/*
must call per UnionName
　genny gen "NodeName=Root UnionName=Version AliasName=IndexString,IndexNum"
*/

import (
	"github.com/kazu/fbshelper/query/base"
)

// import (
// 	b "github.com/kazu/fbshelper/query/base"
// )

//type AliasName generic.Type

var DUMMP_UnionNameAliasName bool = base.SetAlias("UnionName", "AliasName")

func (node UnionName) AliasName() AliasName {
	//result := AliasName{CommonNode: node.CommonNode}
	result := AliasName{}
	result.CommonNode = &CommonNode{}
	result.NodeList = node.NodeList
	result.CommonNode.Name = "AliasName"
	result.FetchIndex()
	return result
}
