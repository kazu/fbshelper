// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call 1 times per Table / struct ( IndexNum ) ;
*/

type IndexNum struct {
	*base.CommonNode
}

func emptyIndexNum() *IndexNum {
	return &IndexNum{CommonNode: &base.CommonNode{}}
}

var IndexNum_IdxToType map[int]int = map[int]int{}
var IndexNum_IdxToTypeGroup map[int]int = map[int]int{}
var IndexNum_IdxToName map[int]string = map[int]string{}
var IndexNum_NameToIdx map[string]int = map[string]int{}

var DUMMP_IndexNumFalse bool = base.SetNameIsStrunct("IndexNum", base.ToBool("False"))

func SetIndexNumFields(nName, fName, fType string, fNum int) bool {

	base.RequestSettingNameFields(nName, fName, fType, fNum)

	enumFtype, ok := base.NameToType[fType]
	if ok {
		IndexNumSetIdxToType(fNum, enumFtype)
	}
	//FIXME: basic type only store?

	IndexNumSetIdxToName(fNum, fType)

	grp := IndexNumGetTypeGroup(fType)
	IndexNumSetTypeGroup(fNum, grp)

	IndexNum_IdxToName[fNum] = fType

	IndexNum_NameToIdx[fName] = fNum
	base.SetNameToIdx("IndexNum", IndexNum_NameToIdx)

	return true

}
func IndexNumSetIdxToName(i int, s string) {
	IndexNum_IdxToName[i] = s

	base.SetIdxToName("IndexNum", IndexNum_IdxToName)
}

func IndexNumSetIdxToType(k, v int) bool {
	IndexNum_IdxToType[k] = v
	base.SetIdxToType("IndexNum", IndexNum_IdxToType)
	return true
}

func IndexNumSetTypeGroup(k, v int) bool {
	IndexNum_IdxToTypeGroup[k] = v
	base.SetdxToTypeGroup("IndexNum", IndexNum_IdxToTypeGroup)
	return true
}

func IndexNumGetTypeGroup(s string) (result int) {
	return base.GetTypeGroup(s)
}

func (node IndexNum) commonNode() *base.CommonNode {
	if node.CommonNode == nil {
		base.Log(base.LOG_WARN, func() base.LogArgs {
			return base.F("CommonNode not found IndexNum")
		})
	} else if len(node.CommonNode.Name) == 0 || len(node.CommonNode.IdxToType) == 0 {
		node.CommonNode.Name = "IndexNum"
		node.CommonNode.IdxToType = IndexNum_IdxToType
		node.CommonNode.IdxToTypeGroup = IndexNum_IdxToTypeGroup
	}
	return node.CommonNode
}
func (node IndexNum) SearchInfo(pos int, fn base.RecFn, condFn base.CondFn) {

	node.commonNode().SearchInfo(pos, fn, condFn)

}

func (node IndexNum) Info() (info base.Info) {

	return node.commonNode().Info()

}

func (node IndexNum) IsLeafAt(j int) bool {

	return node.commonNode().IsLeafAt(j)

}

func (node IndexNum) CountOfField() int {
	return len(IndexNum_IdxToType)
}

func (node IndexNum) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)
}

func (node IndexNum) FieldAt(idx int) *base.CommonNode {
	return node.commonNode().FieldAt(idx)
}

func (node IndexNum) Root() Root {
	return toRoot(node.Base)
}

type IndexNumWithErr struct {
	*IndexNum
	Err error
}

func IndexNumSingle(node *IndexNum, e error) IndexNumWithErr {
	return IndexNumWithErr{IndexNum: node, Err: e}
}

func NewIndexNum() *IndexNum {
	node := emptyIndexNum()
	node.NodeList = &base.NodeList{}
	node.CommonNode.Name = "IndexNum"
	node.Init()

	return node
}

func (node IndexNum) FieldGroups() map[int]int {
	return IndexNum_IdxToTypeGroup
}
