package query

// this is dummy genny template file
import (
	"github.com/kazu/fbshelper/query/base"
	b "github.com/kazu/fbshelper/query/base"
)

/*
type ResultType struct {
	*b.NodeList
}
*/
type ResultType NodeName

type RootType struct {
	*b.CommonNode
}

type AliasName struct {
	*b.CommonNode
}

func NewResultType() *ResultType {
	return &ResultType{CommonNode: &base.CommonNode{}}
}
