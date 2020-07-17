// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import (
	"github.com/kazu/fbshelper/query/base"
	"github.com/kazu/fbshelper/query/log"
)

/*
must call 1 times per Table / struct ( IndexString ) ;
*/

type IndexString struct {
	*base.CommonNode
}

func emptyIndexString() *IndexString {
	return &IndexString{CommonNode: &base.CommonNode{}}
}

var IndexString_IdxToType map[int]int = map[int]int{}
var IndexString_IdxToTypeGroup map[int]int = map[int]int{}
var IndexString_IdxToName map[int]string = map[int]string{}
var IndexString_NameToIdx map[string]int = map[string]int{}

var DUMMP_IndexStringFalse bool = base.SetNameIsStrunct("IndexString", base.ToBool("False"))

func SetIndexStringFields(nName, fName, fType string, fNum int) bool {

	base.RequestSettingNameFields(nName, fName, fType, fNum)

	enumFtype, ok := base.NameToType[fType]
	if ok {
		IndexStringSetIdxToType(fNum, enumFtype)
	}
	//FIXME: basic type only store?

	IndexStringSetIdxToName(fNum, fType)

	grp := IndexStringGetTypeGroup(fType)
	IndexStringSetTypeGroup(fNum, grp)

	IndexString_IdxToName[fNum] = fType

	IndexString_NameToIdx[fName] = fNum
	base.SetNameToIdx("IndexString", IndexString_NameToIdx)

	return true

}
func IndexStringSetIdxToName(i int, s string) {
	IndexString_IdxToName[i] = s

	base.SetIdxToName("IndexString", IndexString_IdxToName)
}

func IndexStringSetIdxToType(k, v int) bool {
	IndexString_IdxToType[k] = v
	base.SetIdxToType("IndexString", IndexString_IdxToType)
	return true
}

func IndexStringSetTypeGroup(k, v int) bool {
	IndexString_IdxToTypeGroup[k] = v
	base.SetdxToTypeGroup("IndexString", IndexString_IdxToTypeGroup)
	return true
}

func IndexStringGetTypeGroup(s string) (result int) {
	return base.GetTypeGroup(s)
}

func (node IndexString) commonNode() *base.CommonNode {
	if node.CommonNode == nil {
		log.Log(log.LOG_WARN, func() log.LogArgs {
			return log.F("CommonNode not found IndexString")
		})
	} else if len(node.CommonNode.Name) == 0 || len(node.CommonNode.IdxToType) == 0 {
		node.CommonNode.Name = "IndexString"
		node.CommonNode.IdxToType = IndexString_IdxToType
		node.CommonNode.IdxToTypeGroup = IndexString_IdxToTypeGroup
	}
	return node.CommonNode
}
func (node IndexString) SearchInfo(pos int, fn base.RecFn, condFn base.CondFn) {

	node.commonNode().SearchInfo(pos, fn, condFn)

}

func (node IndexString) Info() (info base.Info) {

	return node.commonNode().Info()

}

func (node IndexString) IsLeafAt(j int) bool {

	return node.commonNode().IsLeafAt(j)

}

func (node IndexString) CountOfField() int {
	return len(IndexString_IdxToType)
}

func (node IndexString) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)
}

func (node IndexString) FieldAt(idx int) *base.CommonNode {
	return node.commonNode().FieldAt(idx)
}

type IndexStringWithErr struct {
	*IndexString
	Err error
}

func IndexStringSingle(node *IndexString, e error) IndexStringWithErr {
	return IndexStringWithErr{IndexString: node, Err: e}
}

func NewIndexString() *IndexString {
	node := emptyIndexString()
	node.NodeList = &base.NodeList{}
	node.CommonNode.Name = "IndexString"
	node.Init()

	return node
}

func (node IndexString) FieldGroups() map[int]int {
	return IndexString_IdxToTypeGroup
}

func (node IndexString) Root() (Root, error) {
	if !node.InRoot() {
		return Root{}, log.ERR_NO_INCLUDE_ROOT
	}
	root := toRoot(node.Base)
	return root, nil
}
