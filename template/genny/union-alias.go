package query

/*
genny must be called per UnionName
*/

import (
	"github.com/kazu/fbshelper/query/base"
)

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

func UnionNameFromAliasName(v *AliasName) *UnionName {
	result := &UnionName{}
	result.CommonNode = v.CommonNode
	result.FetchIndex()
	return result
}
