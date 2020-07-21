package query

/*
   genny must be called 1 times per NodeName
*/

import (
	"github.com/kazu/fbshelper/query/base"
)

type UnionName struct {
	*base.CommonNode
}

func emptyUnionName() *UnionName {
	result := &UnionName{CommonNode: &base.CommonNode{}}
	result.Name = "UnionName"
	return result
}

func (node UnionName) Member(i int) interface{} {
	return nil
}

func NewUnionName() *UnionName {
	return emptyUnionName()
}
