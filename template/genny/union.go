package query

/*
    must call 1 times per NodeName
	go run github.com/cheekybits/genny gen "UnionName=Index"

*/

import (
	"github.com/kazu/fbshelper/query/base"
)

// import (
// 	b "github.com/kazu/fbshelper/query/base"
// )

type UnionName struct {
	*base.CommonNode
}

func NewUnionName() *UnionName {
	result := &UnionName{CommonNode: &base.CommonNode{}}
	result.Name = "UnionName"
	return result
}

func (node UnionName) Member(i int) interface{} {
	return nil
}
