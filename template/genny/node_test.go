package query

import (
	"testing"

	b "github.com/kazu/fbshelper/query/base"
)

func Test_searchInfo(t *testing.T) {

	a := NewNodeName()
	//NodeNameFieldsNumInit()
	cond := func(pos int, info b.Info) bool {
		return info.Pos <= pos && (info.Pos+info.Size) > pos
		//return true
	}
	recFn := func(s b.NodePath, info b.Info) {
		//		result = append(result, s)
		//		infos = append(infos, info)
	}
	a.SearchInfo(0, recFn, cond)

}
