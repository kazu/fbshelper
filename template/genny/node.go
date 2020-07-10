package query

import (
	"github.com/kazu/fbshelper/query/base"
	_ "github.com/kazu/fbshelper/query/log"
)


// import (
// 	b "github.com/kazu/fbshelper/query/base"
// 	. "github.com/kazu/fbshelper/query/error"
// )


/*
    must call 1 times per NodeName
	go run github.com/cheekybits/genny gen "NodeName=Root IsStruct=false"
*/

type NodeName struct {
	*base.CommonNode
}

func NewNodeName() *NodeName {
	return &NodeName{CommonNode: &base.CommonNode{}}
}

var NodeName_IdxToType map[int]int = map[int]int{}
var NodeName_IdxToTypeGroup map[int]int = map[int]int{}
var NodeName_IdxToName map[int]string = map[int]string{}
var NodeName_NameToIdx map[string]int = map[string]int{}


var DUMMP_NodeNameIsStruct bool = base.SetNameIsStrunct("NodeName", base.ToBool("IsStruct"))

func SetNodeNameFields(nName, fName, fType string, fNum int) bool {
	enumFtype, ok := base.NameToType[fType]
	if ok {
		NodeNameSetIdxToType(fNum, enumFtype)
	}
	//FIXME: 基本型以外は無視?

	NodeNameSetIdxToName(fNum, fType)

	grp := NodeNameGetTypeGroup(fType)
	NodeNameSetTypeGroup(fNum, grp)

	NodeName_IdxToName[fNum] = fType

	NodeName_NameToIdx[fName] = fNum
	base.SetNameToIdx("NodeName", NodeName_NameToIdx)

	return true

}
func NodeNameSetIdxToName(i int, s string) {
	NodeName_IdxToName[i] = s

	base.SetIdxToName("NodeName", NodeName_IdxToName)
}

func NodeNameSetIdxToType(k, v int) bool {
	NodeName_IdxToType[k] = v
	base.SetIdxToType("NodeName", NodeName_IdxToType)
	return true
}

func NodeNameSetTypeGroup(k, v int) bool {
	NodeName_IdxToTypeGroup[k] = v
	base.SetdxToTypeGroup("NodeName", NodeName_IdxToTypeGroup)
	return true
}

func NodeNameGetTypeGroup(s string) (result int) {
	return base.GetTypeGroup(s)
}

func (node NodeName) commonNode() *base.CommonNode {
	if node.CommonNode == nil {
		base.Log(base.LOG_WARN, func() base.LogArgs {
			return base.F("CommonNode not found NodeName")
		})
	}else if len(node.CommonNode.Name) == 0 || len(node.CommonNode.IdxToType) == 0 {
		node.CommonNode.Name = "NodeName"
		node.CommonNode.IdxToType = NodeName_IdxToType
		node.CommonNode.IdxToTypeGroup = NodeName_IdxToTypeGroup
	}
	return node.CommonNode 
}
func (node NodeName) SearchInfo(pos int, fn base.RecFn, condFn base.CondFn) {

	node.commonNode().SearchInfo(pos, fn, condFn)
	// info := node.Info()
	// if condFn(pos, info) {
	//     fn(base.NodePath{Name: "NodeName", Idx: -1}, info)
	// }else{
	//     return
	// }

	// for i := 0; i < node.CountOfField(); i++ {
	// 	g :=  NodeName_IdxToTypeGroup[i]
	//     if node.IsLeafAt(i) {
	//         fInfo := base.Info(node.ValueInfo(i))
	//         if condFn(pos, fInfo) {
	//             fn(base.NodePath{Name: "NodeName", Idx: i}, fInfo)
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

// if base.IsStructName["NodeName"] {
// } else if base.IsFieldStruct(i) {
// } else if base.IsFieldUnion(i) {
// } else if base.IsFieldBytes(i) {
// } else if base.IsFieldSlice(i) {
// } else if base.IsFieldTable(i) {
// } else if base.IsFieldBasicType(i) {
// } else {
// 	base.Log(base.LOG_ERROR, func() base.LogArgs {
// 		return base.F("Invalid %s.%s idx=%d\n", "NodeName", "FieldName", i)
// 	})
// }

func (node NodeName) Info() (info base.Info) {

	return node.commonNode().Info()

	// if node.Node == nil {
	// 	node.Node = &base.Node{}
	// }

	// info.Pos = node.Pos
	// info.Size = -1
	// if base.IsStructName["NodeName"] {
	// 	size := 0
	// 	for i :=0; i < node.CountOfField(); i++ {
	// 		size += base.TypeToSize[NodeName_IdxToType[i]]
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

func (node NodeName) IsLeafAt(j int) bool {

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
	// 			return base.F("Invalid %s.%s idx=%d\n", "NodeName", "FieldName", i)
	// 	})
	// }
	// return false
}

func (node NodeName) CountOfField() int {
	return len(NodeName_IdxToType)
}

func (node NodeName) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)

	// if base.IsStructName["NodeName"] {
	// 	if len(node.ValueInfos) > i {
	// 		return node.ValueInfos[i]
	// 	}
	// 	node.ValueInfos = make([]base.ValueInfo, 0, node.CountOfField())
	// 	info := base.ValueInfo{Pos: node.Pos, Size: 0}
	// 	for i :=0; i < node.CountOfField(); i++ {
	// 		info.Pos += info.Size
	// 		info.Size = base.TypeToSize[NodeName_IdxToType[i]]
	// 		node.ValueInfos = append(node.ValueInfos,  info)
	// 	}
	// }

	// grp :=  NodeName_IdxToTypeGroup[i]

	// if base.IsFieldStruct(grp) {
	// 	if node.ValueInfos[i].IsNotReady() {
	// 			node.ValueInfoPos(i)
	// 	}

	// 	fTypeStr := NodeName_IdxToName[j]
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
	// 			return base.F("Invalid %s.%s idx=%d\n", "NodeName", "FieldName", i)
	// 		})
	// }

	// return base.ValueInfo{}
}

func (node NodeName) FieldAt(idx int) (*base.CommonNode) {
	//return node.commonNode().FieldAt(idx).NodeList
	return node.commonNode().FieldAt(idx)
}

// func (node NodeName) SizeAsStruct() int {
// 	if base.IsStructName["NodeName"] {
// 		size := 0
// 		for i :=0; i < node.CountOfField(); i++ {
// 			size += base.TypeToSize[NodeName_IdxToType[i]]
// 		}
// 		return size
// 	}
// }p
