
// Code generated by genmaps.go; DO NOT EDIT.
// template file is https://github.com/kazu/fbshelper/blob/master/template/query.go.tmpl github.com/kazu/fbshelper/template/query.go.tmpl 
//   https://github.com/kazu/fbshelper/blob/master/template/union.query.go.tmpl


package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

const (
    DUMMY_IndexString = flatbuffers.VtableMetadataFields
)

type FbsIndexString struct {
	*base.Node
}


type FbsIndexStringMaps struct {
    *base.NodeList
}
func (node FbsIndexString) Info() base.Info {

    info := base.Info{Pos: node.Pos, Size: -1}
    for i := 0; i < len(node.VTable); i++ {
        vInfo := node.ValueInfo(i)
        if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
            info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
        }
    }
    return info    
}

func (node FbsIndexString) ValueInfo(i int) base.ValueInfo {

    switch i {
    case 0:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOfint32
    case 1:
         if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPosList(i)
        }
        node.ValueInfos[i].Size = node.Maps().Info().Size
     }
     return node.ValueInfos[i]
}





func (node FbsIndexString) Size() int32 {
    if node.VTable[0] == 0 {
        return int32(0)
    }
    return int32(flatbuffers.GetInt32(node.ValueNormal(0)))
}


func (node FbsIndexString) Maps() FbsIndexStringMaps {
    if node.VTable[1] == 0 {
        return FbsIndexStringMaps{}
    }
    nodelist :=  node.ValueList(1)
    return FbsIndexStringMaps{
                NodeList: &nodelist,
    }
}






func (node FbsIndexStringMaps) At(i int) FbsInvertedMapString {
    if i > int(node.ValueInfo.VLen) || i < 0 {
		return FbsInvertedMapString{}
	}

	buf := node.Bytes
	ptr := uint32(node.ValueInfo.Pos + (i-1)*4)
	return FbsInvertedMapString{Node: base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:])))}
}


func (node FbsIndexStringMaps) First() FbsInvertedMapString {
	return node.At(0)
}


func (node FbsIndexStringMaps) Last() FbsInvertedMapString {
	return node.At(int(node.ValueInfo.VLen))
}

func (node FbsIndexStringMaps) Select(fn func(m FbsInvertedMapString) bool) []FbsInvertedMapString {

	result := make([]FbsInvertedMapString, 0, int(node.ValueInfo.VLen))
	for i := 0; i < int(node.ValueInfo.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node FbsIndexStringMaps) Find(fn func(m FbsInvertedMapString) bool) FbsInvertedMapString{

	for i := 0; i < int(node.ValueInfo.VLen); i++ {
		if m := node.At(i); fn(m) {
			return m
		}
	}
	return FbsInvertedMapString{}
}

func (node FbsIndexStringMaps) All() []FbsInvertedMapString {
	return node.Select(func(m FbsInvertedMapString) bool { return true })
}

func (node FbsIndexStringMaps) Count() int {
	return int(node.ValueInfo.VLen)
}

func (node FbsIndexStringMaps) Info() base.Info {

    info := base.Info{Pos: node.ValueInfo.Pos, Size: -1}
    vInfo := node.Last().Info()



    if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
        info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
    }
    return info
}
