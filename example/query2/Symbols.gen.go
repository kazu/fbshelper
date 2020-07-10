// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

// import (
// 	b "github.com/kazu/fbshelper/query/base"
// 	. "github.com/kazu/fbshelper/query/error"
// )

/*
must call 1 times per Symbols ;
go run github . com / cheekybits / genny gen "Symbols=Root False=false" ;
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
	//FIXME: 基本型以外は無視?

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
	// info := node.Info()
	// if condFn(pos, info) {
	// fn(base.NodePath{Name: "Symbols", Idx: -1}, info)
	// }else{
	//     return
	// }

	// for i := 0; i < node.CountOfField(); i++ {
	// g := Symbols_IdxToTypeGroup[i]
	//     if node.IsLeafAt(i) {
	//         fInfo := base.Info(node.ValueInfo(i))
	//         if condFn(pos, fInfo) {
	// fn(base.NodePath{Name: "Symbols", Idx: i}, fInfo)
	//         }
	//         continue
	//     }
	// 	if base.IsMatchBit(g, base.FieldTypeStruct) {
	// 		  node.FieldAt(i).SearchInfo(pos, fn, condFn)
	//     } else if base.IsMatchBit(g, base.FieldTypeUnion) {
	// 		  mNode, _ := node.FieldAt(i).Member(int(node.FieldAt(i - 1).Byte())).(base.Noder)

	// 		  mNode.SearchInfo(pos, fn, condFn)
	// 	} else if base.IsMatchBit(g, base.FieldTypeSlice) && base.IsMatchBit(g, base.FieldTypeBasic1) {
	// 	} else if base.IsMatchBit(g, base.FieldTypeSlice) {
	// 		  node.FieldAt(i).SearchInfo(pos, fn, condFn)
	// 	} else if base.IsMatchBit(g, base.FieldTypeTable) {
	// 		  node.FieldAt(i).SearchInfo(pos, fn, condFn)
	// 	} else if base.IsMatchBit(g, base.FieldTypeBasic) {
	//     } else {
	// 		  base.Log(base.LOG_ERROR, func() base.LogArgs {
	//             return base.F("node must be Noder")
	//         })
	// 	}
	// }
}

// if base.FalseName["Symbols"] {
// } else if base.IsFieldStruct(i) {
// } else if base.IsFieldUnion(i) {
// } else if base.IsFieldBytes(i) {
// } else if base.IsFieldSlice(i) {
// } else if base.IsFieldTable(i) {
// } else if base.IsFieldBasicType(i) {
// } else {
// 	base.Log(base.LOG_ERROR, func() base.LogArgs {
// return base.F("Invalid %s.%s idx=%d\n", "Symbols", "FieldName", i)
// 	})
// }

func (node Symbols) Info() (info base.Info) {

	return node.commonNode().Info()

	// if node.Node == nil {
	// 	node.Node = &base.Node{}
	// }

	// info.Pos = node.Pos
	// info.Size = -1
	// if base.FalseName["Symbols"] {
	// 	size := 0
	// 	for i :=0; i < node.CountOfField(); i++ {
	// size += base.TypeToSize[Symbols_IdxToType[i]]
	// 	}
	// 	info.Size = size
	// 	return info
	// }

	// for i := 0; i < len(node.VTable); i++ {
	//     vInfo := node.ValueInfo(i)
	//     if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
	//         info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
	//     }
	// }
	// return info

}

func (node Symbols) IsLeafAt(j int) bool {

	return node.commonNode().IsLeafAt(j)

	// if base.IsFieldStruct(i) {

	// 	return false
	// } else if base.IsFieldUnion(i) {

	// 	return false
	// } else if base.IsFieldBytes(i) {

	// 	return true
	// } else if base.IsFieldSlice(i) {

	// 	return false
	// } else if base.IsFieldTable(i) {

	// 	return false
	// } else if base.IsFieldBasicType(i) {
	// 	return true
	// } else {
	// 	base.Log(base.LOG_ERROR, func() base.LogArgs {
	// return base.F("Invalid %s.%s idx=%d\n", "Symbols", "FieldName", i)
	// 	})
	// }
	// return false
}

func (node Symbols) CountOfField() int {
	return len(Symbols_IdxToType)
}

func (node Symbols) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)

	// if base.FalseName["Symbols"] {
	// 	if len(node.ValueInfos) > i {
	// 		return node.ValueInfos[i]
	// 	}
	// 	node.ValueInfos = make([]base.ValueInfo, 0, node.CountOfField())
	// 	info := base.ValueInfo{Pos: node.Pos, Size: 0}
	// 	for i :=0; i < node.CountOfField(); i++ {
	// 		info.Pos += info.Size
	// info.Size = base.TypeToSize[Symbols_IdxToType[i]]
	// 		node.ValueInfos = append(node.ValueInfos,  info)
	// 	}
	// }

	// grp := Symbols_IdxToTypeGroup[i]

	// if base.IsFieldStruct(grp) {
	// 	if node.ValueInfos[i].IsNotReady() {
	// 			node.ValueInfoPos(i)
	// 	}

	// fTypeStr := Symbols_IdxToName[j]
	// 	idxToType = All_IdxToType[fTypeStr]
	// 	size := 0
	// 	for nextIdx := 0 ; nextIdx < len(idxToType);  nextIdx++ {
	// 		size +=base.TypeToSize[idxToType[nextIdx]]
	// 	}
	// 	node.ValueInfos[i].Size = size

	// } else if base.IsFieldUnion(i) {

	// } else if base.IsFieldBytes(i) {
	// } else if base.IsFieldSlice(i) {
	// } else if base.IsFieldTable(i) {
	// } else if base.IsFieldBasicType(i) {
	// } else {
	// 		base.Log(base.LOG_ERROR, func() base.LogArgs {
	// return base.F("Invalid %s.%s idx=%d\n", "Symbols", "FieldName", i)
	// 		})
	// }

	// return base.ValueInfo{}
}

func (node Symbols) FieldAt(idx int) *base.CommonNode {
	//return node.commonNode().FieldAt(idx).NodeList
	return node.commonNode().FieldAt(idx)
}

// func (node Symbols) SizeAsStruct() int {
// if base.FalseName["Symbols"] {
// 		size := 0
// 		for i :=0; i < node.CountOfField(); i++ {
// size += base.TypeToSize[Symbols_IdxToType[i]]
// 		}
// 		return size
// 	}