// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call 1 times per Table / struct ( InvertedMapNum ) ;
*/

type InvertedMapNum struct {
	*base.CommonNode
}

func NewInvertedMapNum() *InvertedMapNum {
	return &InvertedMapNum{CommonNode: &base.CommonNode{}}
}

var InvertedMapNum_IdxToType map[int]int = map[int]int{}
var InvertedMapNum_IdxToTypeGroup map[int]int = map[int]int{}
var InvertedMapNum_IdxToName map[int]string = map[int]string{}
var InvertedMapNum_NameToIdx map[string]int = map[string]int{}

var DUMMP_InvertedMapNumFalse bool = base.SetNameIsStrunct("InvertedMapNum", base.ToBool("False"))

func SetInvertedMapNumFields(nName, fName, fType string, fNum int) bool {

	base.RequestSettingNameFields(nName, fName, fType, fNum)

	enumFtype, ok := base.NameToType[fType]
	if ok {
		InvertedMapNumSetIdxToType(fNum, enumFtype)
	}
	//FIXME: basic type only store?

	InvertedMapNumSetIdxToName(fNum, fType)

	grp := InvertedMapNumGetTypeGroup(fType)
	InvertedMapNumSetTypeGroup(fNum, grp)

	InvertedMapNum_IdxToName[fNum] = fType

	InvertedMapNum_NameToIdx[fName] = fNum
	base.SetNameToIdx("InvertedMapNum", InvertedMapNum_NameToIdx)

	return true

}
func InvertedMapNumSetIdxToName(i int, s string) {
	InvertedMapNum_IdxToName[i] = s

	base.SetIdxToName("InvertedMapNum", InvertedMapNum_IdxToName)
}

func InvertedMapNumSetIdxToType(k, v int) bool {
	InvertedMapNum_IdxToType[k] = v
	base.SetIdxToType("InvertedMapNum", InvertedMapNum_IdxToType)
	return true
}

func InvertedMapNumSetTypeGroup(k, v int) bool {
	InvertedMapNum_IdxToTypeGroup[k] = v
	base.SetdxToTypeGroup("InvertedMapNum", InvertedMapNum_IdxToTypeGroup)
	return true
}

func InvertedMapNumGetTypeGroup(s string) (result int) {
	return base.GetTypeGroup(s)
}

func (node InvertedMapNum) commonNode() *base.CommonNode {
	if node.CommonNode == nil {
		base.Log(base.LOG_WARN, func() base.LogArgs {
			return base.F("CommonNode not found InvertedMapNum")
		})
	} else if len(node.CommonNode.Name) == 0 || len(node.CommonNode.IdxToType) == 0 {
		node.CommonNode.Name = "InvertedMapNum"
		node.CommonNode.IdxToType = InvertedMapNum_IdxToType
		node.CommonNode.IdxToTypeGroup = InvertedMapNum_IdxToTypeGroup
	}
	return node.CommonNode
}
func (node InvertedMapNum) SearchInfo(pos int, fn base.RecFn, condFn base.CondFn) {

	node.commonNode().SearchInfo(pos, fn, condFn)

}

func (node InvertedMapNum) Info() (info base.Info) {

	return node.commonNode().Info()

}

func (node InvertedMapNum) IsLeafAt(j int) bool {

	return node.commonNode().IsLeafAt(j)

}

func (node InvertedMapNum) CountOfField() int {
	return len(InvertedMapNum_IdxToType)
}

func (node InvertedMapNum) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)
}

func (node InvertedMapNum) FieldAt(idx int) *base.CommonNode {
	return node.commonNode().FieldAt(idx)
}
