// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call 1 times per Table / struct ( Symbols ) ;
*/

type Symbols struct {
	*base.CommonNode
}

func NewSymbols() *Symbols {
	return &Symbols{CommonNode: &base.CommonNode{}}
}

var Symbols_IdxToType map[int]int = map[int]int{}
var Symbols_IdxToTypeGroup map[int]int = map[int]int{}
var Symbols_IdxToName map[int]string = map[int]string{}
var Symbols_NameToIdx map[string]int = map[string]int{}

var DUMMP_SymbolsFalse bool = base.SetNameIsStrunct("Symbols", base.ToBool("False"))

func SetSymbolsFields(nName, fName, fType string, fNum int) bool {
	enumFtype, ok := base.NameToType[fType]
	if ok {
		SymbolsSetIdxToType(fNum, enumFtype)
	}
	//FIXME: basic type only store?

	SymbolsSetIdxToName(fNum, fType)

	grp := SymbolsGetTypeGroup(fType)
	SymbolsSetTypeGroup(fNum, grp)

	Symbols_IdxToName[fNum] = fType

	Symbols_NameToIdx[fName] = fNum
	base.SetNameToIdx("Symbols", Symbols_NameToIdx)

	return true

}
func SymbolsSetIdxToName(i int, s string) {
	Symbols_IdxToName[i] = s

	base.SetIdxToName("Symbols", Symbols_IdxToName)
}

func SymbolsSetIdxToType(k, v int) bool {
	Symbols_IdxToType[k] = v
	base.SetIdxToType("Symbols", Symbols_IdxToType)
	return true
}

func SymbolsSetTypeGroup(k, v int) bool {
	Symbols_IdxToTypeGroup[k] = v
	base.SetdxToTypeGroup("Symbols", Symbols_IdxToTypeGroup)
	return true
}

func SymbolsGetTypeGroup(s string) (result int) {
	return base.GetTypeGroup(s)
}

func (node Symbols) commonNode() *base.CommonNode {
	if node.CommonNode == nil {
		base.Log(base.LOG_WARN, func() base.LogArgs {
			return base.F("CommonNode not found Symbols")
		})
	} else if len(node.CommonNode.Name) == 0 || len(node.CommonNode.IdxToType) == 0 {
		node.CommonNode.Name = "Symbols"
		node.CommonNode.IdxToType = Symbols_IdxToType
		node.CommonNode.IdxToTypeGroup = Symbols_IdxToTypeGroup
	}
	return node.CommonNode
}
func (node Symbols) SearchInfo(pos int, fn base.RecFn, condFn base.CondFn) {

	node.commonNode().SearchInfo(pos, fn, condFn)

}

func (node Symbols) Info() (info base.Info) {

	return node.commonNode().Info()

}

func (node Symbols) IsLeafAt(j int) bool {

	return node.commonNode().IsLeafAt(j)

}

func (node Symbols) CountOfField() int {
	return len(Symbols_IdxToType)
}

func (node Symbols) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)
}

func (node Symbols) FieldAt(idx int) *base.CommonNode {
	return node.commonNode().FieldAt(idx)
}
