// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
must call 1 times per Table / struct ( Symbol ) ;
*/

type Symbol struct {
	*base.CommonNode
}

func emptySymbol() *Symbol {
	return &Symbol{CommonNode: &base.CommonNode{}}
}

var Symbol_IdxToType map[int]int = map[int]int{}
var Symbol_IdxToTypeGroup map[int]int = map[int]int{}
var Symbol_IdxToName map[int]string = map[int]string{}
var Symbol_NameToIdx map[string]int = map[string]int{}

var DUMMP_SymbolFalse bool = base.SetNameIsStrunct("Symbol", base.ToBool("False"))

func SetSymbolFields(nName, fName, fType string, fNum int) bool {

	base.RequestSettingNameFields(nName, fName, fType, fNum)

	enumFtype, ok := base.NameToType[fType]
	if ok {
		SymbolSetIdxToType(fNum, enumFtype)
	}
	//FIXME: basic type only store?

	SymbolSetIdxToName(fNum, fType)

	grp := SymbolGetTypeGroup(fType)
	SymbolSetTypeGroup(fNum, grp)

	Symbol_IdxToName[fNum] = fType

	Symbol_NameToIdx[fName] = fNum
	base.SetNameToIdx("Symbol", Symbol_NameToIdx)

	return true

}
func SymbolSetIdxToName(i int, s string) {
	Symbol_IdxToName[i] = s

	base.SetIdxToName("Symbol", Symbol_IdxToName)
}

func SymbolSetIdxToType(k, v int) bool {
	Symbol_IdxToType[k] = v
	base.SetIdxToType("Symbol", Symbol_IdxToType)
	return true
}

func SymbolSetTypeGroup(k, v int) bool {
	Symbol_IdxToTypeGroup[k] = v
	base.SetdxToTypeGroup("Symbol", Symbol_IdxToTypeGroup)
	return true
}

func SymbolGetTypeGroup(s string) (result int) {
	return base.GetTypeGroup(s)
}

func (node Symbol) commonNode() *base.CommonNode {
	if node.CommonNode == nil {
		base.Log(base.LOG_WARN, func() base.LogArgs {
			return base.F("CommonNode not found Symbol")
		})
	} else if len(node.CommonNode.Name) == 0 || len(node.CommonNode.IdxToType) == 0 {
		node.CommonNode.Name = "Symbol"
		node.CommonNode.IdxToType = Symbol_IdxToType
		node.CommonNode.IdxToTypeGroup = Symbol_IdxToTypeGroup
	}
	return node.CommonNode
}
func (node Symbol) SearchInfo(pos int, fn base.RecFn, condFn base.CondFn) {

	node.commonNode().SearchInfo(pos, fn, condFn)

}

func (node Symbol) Info() (info base.Info) {

	return node.commonNode().Info()

}

func (node Symbol) IsLeafAt(j int) bool {

	return node.commonNode().IsLeafAt(j)

}

func (node Symbol) CountOfField() int {
	return len(Symbol_IdxToType)
}

func (node Symbol) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)
}

func (node Symbol) FieldAt(idx int) *base.CommonNode {
	return node.commonNode().FieldAt(idx)
}

func (node Symbol) Root() Root {
	return toRoot(node.Base)
}

type SymbolWithErr struct {
	*Symbol
	Err error
}

func SymbolSingle(node *Symbol, e error) SymbolWithErr {
	return SymbolWithErr{Symbol: node, Err: e}
}
