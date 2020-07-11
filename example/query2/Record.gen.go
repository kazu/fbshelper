// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call 1 times per Table / struct ( Record ) ;
*/

type Record struct {
	*base.CommonNode
}

func NewRecord() *Record {
	return &Record{CommonNode: &base.CommonNode{}}
}

var Record_IdxToType map[int]int = map[int]int{}
var Record_IdxToTypeGroup map[int]int = map[int]int{}
var Record_IdxToName map[int]string = map[int]string{}
var Record_NameToIdx map[string]int = map[string]int{}

var DUMMP_RecordTrue bool = base.SetNameIsStrunct("Record", base.ToBool("True"))

func SetRecordFields(nName, fName, fType string, fNum int) bool {

	base.RequestSettingNameFields(nName, fName, fType, fNum)

	enumFtype, ok := base.NameToType[fType]
	if ok {
		RecordSetIdxToType(fNum, enumFtype)
	}
	//FIXME: basic type only store?

	RecordSetIdxToName(fNum, fType)

	grp := RecordGetTypeGroup(fType)
	RecordSetTypeGroup(fNum, grp)

	Record_IdxToName[fNum] = fType

	Record_NameToIdx[fName] = fNum
	base.SetNameToIdx("Record", Record_NameToIdx)

	return true

}
func RecordSetIdxToName(i int, s string) {
	Record_IdxToName[i] = s

	base.SetIdxToName("Record", Record_IdxToName)
}

func RecordSetIdxToType(k, v int) bool {
	Record_IdxToType[k] = v
	base.SetIdxToType("Record", Record_IdxToType)
	return true
}

func RecordSetTypeGroup(k, v int) bool {
	Record_IdxToTypeGroup[k] = v
	base.SetdxToTypeGroup("Record", Record_IdxToTypeGroup)
	return true
}

func RecordGetTypeGroup(s string) (result int) {
	return base.GetTypeGroup(s)
}

func (node Record) commonNode() *base.CommonNode {
	if node.CommonNode == nil {
		base.Log(base.LOG_WARN, func() base.LogArgs {
			return base.F("CommonNode not found Record")
		})
	} else if len(node.CommonNode.Name) == 0 || len(node.CommonNode.IdxToType) == 0 {
		node.CommonNode.Name = "Record"
		node.CommonNode.IdxToType = Record_IdxToType
		node.CommonNode.IdxToTypeGroup = Record_IdxToTypeGroup
	}
	return node.CommonNode
}
func (node Record) SearchInfo(pos int, fn base.RecFn, condFn base.CondFn) {

	node.commonNode().SearchInfo(pos, fn, condFn)

}

func (node Record) Info() (info base.Info) {

	return node.commonNode().Info()

}

func (node Record) IsLeafAt(j int) bool {

	return node.commonNode().IsLeafAt(j)

}

func (node Record) CountOfField() int {
	return len(Record_IdxToType)
}

func (node Record) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)
}

func (node Record) FieldAt(idx int) *base.CommonNode {
	return node.commonNode().FieldAt(idx)
}
