package query

/*
    must call 1 times per NodeName
	go run github.com/cheekybits/genny gen "UnionName=Index"

*/

import (
	b "github.com/kazu/fbshelper/query/base"
)

type UnionName struct {
	*b.Node
}

// func InitUnionName() {
// 	b.UnionAlias["UnionName"] = make([]string, 0, 10)
// }

func (node NodeName) Member(i int) interface{} {
	return nil
}
