package query

/*
must call per UnionName
　genny gen "NodeName=Root UnionName=Version AliasName=IndexString,IndexNum"
*/

import (
	b "github.com/kazu/fbshelper/query/base"
)

//type AliasName generic.Type

var DUMMP_NodeNameAliasName bool = b.SetAlias("UnionName", "AliasName")
