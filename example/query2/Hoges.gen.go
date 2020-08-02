// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import (
	"github.com/kazu/fbshelper/query/base"
	"github.com/kazu/fbshelper/query/log"
)

/*
must call 1 times per Table / struct ( Hoges ) ;
*/

type Hoges struct {
	*base.CommonNode
}

func emptyHoges() *Hoges {
	return &Hoges{CommonNode: &base.CommonNode{}}
}

var Hoges_IdxToType map[int]int = map[int]int{}
var Hoges_IdxToTypeGroup map[int]int = map[int]int{}
var Hoges_IdxToName map[int]string = map[int]string{}
var Hoges_NameToIdx map[string]int = map[string]int{}

var DUMMP_HogesFalse bool = base.SetNameIsStrunct("Hoges", base.ToBool("False"))

func SetHogesFields(nName, fName, fType string, fNum int) bool {

	base.RequestSettingNameFields(nName, fName, fType, fNum)

	enumFtype, ok := base.NameToType[fType]
	if ok {
		HogesSetIdxToType(fNum, enumFtype)
	}
	//FIXME: basic type only store?

	HogesSetIdxToName(fNum, fType)

	grp := HogesGetTypeGroup(fType)
	HogesSetTypeGroup(fNum, grp)

	Hoges_IdxToName[fNum] = fType

	Hoges_NameToIdx[fName] = fNum
	base.SetNameToIdx("Hoges", Hoges_NameToIdx)

	return true

}
func HogesSetIdxToName(i int, s string) {
	Hoges_IdxToName[i] = s

	base.SetIdxToName("Hoges", Hoges_IdxToName)
}

func HogesSetIdxToType(k, v int) bool {
	Hoges_IdxToType[k] = v
	base.SetIdxToType("Hoges", Hoges_IdxToType)
	return true
}

func HogesSetTypeGroup(k, v int) bool {
	Hoges_IdxToTypeGroup[k] = v
	base.SetdxToTypeGroup("Hoges", Hoges_IdxToTypeGroup)
	return true
}

func HogesGetTypeGroup(s string) (result int) {
	return base.GetTypeGroup(s)
}

func (node Hoges) commonNode() *base.CommonNode {
	if node.CommonNode == nil {
		log.Log(log.LOG_WARN, func() log.LogArgs {
			return log.F("CommonNode not found Hoges")
		})
	} else if len(node.CommonNode.Name) == 0 || len(node.CommonNode.IdxToType) == 0 {
		node.CommonNode.Name = "Hoges"
		node.CommonNode.IdxToType = Hoges_IdxToType
		node.CommonNode.IdxToTypeGroup = Hoges_IdxToTypeGroup
	}
	return node.CommonNode
}
func (node Hoges) SearchInfo(pos int, fn base.RecFn, condFn base.CondFn) {

	node.commonNode().SearchInfo(pos, fn, condFn)

}

func (node Hoges) Info() (info base.Info) {

	return node.commonNode().Info()

}

func (node Hoges) IsLeafAt(j int) bool {

	return node.commonNode().IsLeafAt(j)

}

func (node Hoges) CountOfField() int {
	return len(Hoges_IdxToType)
}

func (node Hoges) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)
}

func (node Hoges) FieldAt(idx int) *base.CommonNode {
	return node.commonNode().FieldAt(idx)
}

type HogesWithErr struct {
	*Hoges
	Err error
}

func HogesSingle(node *Hoges, e error) HogesWithErr {
	return HogesWithErr{Hoges: node, Err: e}
}

func NewHoges() *Hoges {
	base.ApplyRequestNameFields()
	node := emptyHoges()
	node.NodeList = &base.NodeList{}
	node.CommonNode.Name = "Hoges"
	node.Init()

	return node
}

func (node Hoges) FieldGroups() map[int]int {
	return Hoges_IdxToTypeGroup
}

func (node Hoges) Root() (Root, error) {
	if !node.InRoot() {
		return Root{}, log.ERR_NO_INCLUDE_ROOT
	}
	root := toRoot(node.Base)
	return root, nil
}