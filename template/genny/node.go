package query

import (
	//"github.com/cheekybits/genny/generic"
	b "github.com/kazu/fbshelper/query/base"
)
/*
    must call 1 times per NodeName
	go run github.com/cheekybits/genny gen "NodeName=Root IsStruct=false"

*/ 
//type NodeName generic.Type
//type FieldType geneic.Type
//type FieldNumType generic.Number
//type FieldsNums generic.Number
//type IsStruct generic.Number

type NodeName struct{
	*b.Node
}

var NodeName_IdxToType map[int]int = map[int]int{}
var NodeName_IdxToTypeGroup map[int]int = map[int]int{}

var DUMMP_NodeNameIsStruct bool = b.SetNameIsStrunct("NodeName", b.ToBool("IsStruct"))

func SetNodeToIdx(k, v int ) bool {
	NodeName_IdxToType[k] = v
	return true
}

func SetTypeGroup(k, v int) bool {
	NodeName_IdxToTypeGroup[k] = v
	return true

}

func GetTypeGroup(s string) (result int) {

	result = 0 
	if enum, ok := b.NameToType[s]; ok {
		return b.TypeToGroup[enum]
	}

	if _, ok := b.UnionAlias[s]; ok {
		return b.FieldTypeUnion
	}
	if s[0:2] == "[]byte" {
		return b.FieldTypeSlice | b.FieldTypeBasic1
	}

	if s[0:2] == "[]" {
		result |= b.FieldTypeSlice
	}

	if enum, ok := b.NameToType[s[2:]]; ok {
		result |= b.TypeToGroup[enum]
		return
	}

	result |= b.FieldTypeTable
	return 

}


func (node NodeName) SearchInfo(pos int, fn b.RecFn, condFn b.CondFn) {
	info := node.Info()

	if condFn(pos, info) {
        fn(b.NodePath{Name: "NodeName", Idx: -1}, info)
    }else{
        return
	}
	
	for i := 0; i < node.CountOfField(); i++ {
		g :=  NodeName_IdxToTypeGroup[i]
        if node.IsLeafAt(i) {
            fInfo := b.Info(node.ValueInfo(i))
            if condFn(pos, fInfo) {
                fn(b.NodePath{Name: "NodeName", Idx: i}, fInfo)
            }
            continue
        }
		if b.IsMatchBit(g, b.FieldTypeStruct) {
			  node.FieldAt(i).SearchInfo(pos, fn, condFn)
        } else if b.IsMatchBit(g, b.FieldTypeUnion) {
			  mNode, _ := node.FieldAt(i).Member(int(node.FieldAt(i - 1).Byte())).(b.Noder)

			  mNode.SearchInfo(pos, fn, condFn)
		} else if b.IsMatchBit(g, b.FieldTypeSlice) && b.IsMatchBit(g, b.FieldTypeBasic1) {
		} else if b.IsMatchBit(g, b.FieldTypeSlice) {	
			  node.FieldAt(i).SearchInfo(pos, fn, condFn)
		} else if b.IsMatchBit(g, b.FieldTypeTable) {	    
			  node.FieldAt(i).SearchInfo(pos, fn, condFn)
		} else if b.IsMatchBit(g, b.FieldTypeBasic) {
	    } else {
			  b.Log(b.LOG_ERROR, func() b.LogArgs {
                return b.F("node must be Noder")
            })
		}
	}
}


// if b.IsFieldStruct(i) {
// } else if b.IsFieldUnion(i) {
// } else if b.IsFieldBytes(i) {
// } else if b.IsFieldSlice(i) {
// } else if b.IsFieldTable(i) {	
// } else if b.IsFieldBasicType(i) {
// } else {	
// 	b.Log(b.LOG_ERROR, func() b.LogArgs {
// 		return b.F("Invalid %s.%s idx=%d\n", "NodeName", "FieldName", i)
// 	})
// }

func (node NodeName) Info() (info b.Info) {
	if node.Node == nil {
		node.Node = &b.Node{}
	}

	info.Pos = node.Pos
	info.Size = -1
	if b.IsStructName["NodeName"] {	
		size := 0
		for i :=0; i < node.CountOfField(); i++ {
			size += b.TypeToSize[NodeName_IdxToType[i]]
		}
		info.Size = size
		return info
	}

    for i := 0; i < len(node.VTable); i++ {
        vInfo := node.ValueInfo(i)
        if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
            info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
        }
	}
	return info

}


func (node NodeName) IsLeafAt(j int) bool {
		i :=  NodeName_IdxToTypeGroup[j]
	if b.IsFieldStruct(i) {
			
		return false
	} else if b.IsFieldUnion(i) {
			
		
		return false
	} else if b.IsFieldBytes(i) {
		
		return true	
	} else if b.IsFieldSlice(i) {
		
		return false
	} else if b.IsFieldTable(i) {	

		return false
	} else if b.IsFieldBasicType(i) {
		return true
	} else {	
		b.Log(b.LOG_ERROR, func() b.LogArgs {
				return b.F("Invalid %s.%s idx=%d\n", "NodeName", "FieldName", i)
		})
	}
	return false
}

// mock
func (node NodeName) CountOfField() int {
	return len(NodeName_IdxToType)
}
	

func (node NodeName) ValueInfo(i int) b.ValueInfo {
	return b.ValueInfo{}
}

func (node NodeName) FieldAt(i int) NodeName {
	return NodeName{}
}
