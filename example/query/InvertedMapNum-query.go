
// Code generated by genmaps.go; DO NOT EDIT.
// template file is https://github.com/kazu/fbshelper/blob/master/template/query.go.tmpl github.com/kazu/fbshelper/template/query.go.tmpl 
//   https://github.com/kazu/fbshelper/blob/master/template/union.query.go.tmpl


package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

const (
    DUMMY_InvertedMapNum = flatbuffers.VtableMetadataFields
)

const (
        InvertedMapNum_Key =     0
        InvertedMapNum_Value =     1
)


type FbsInvertedMapNum struct {
	*base.Node
}


func (node FbsInvertedMapNum) Info() base.Info {

    info := base.Info{Pos: node.Pos, Size: -1}
    for i := 0; i < len(node.VTable); i++ {
        vInfo := node.ValueInfo(i)
        if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
            info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
        }
    }
    return info    
}

func (node FbsInvertedMapNum) ValueInfo(i int) base.ValueInfo {

    switch i {
    case 0:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOfint64
    case 1:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = node.Value().Info().Size
     }
     return node.ValueInfos[i]
}

func (node FbsInvertedMapNum) FieldAt(i int) interface{} {

    switch i {
    case 0:
        return node.Key()
    case 1:
        return node.Value()
     }
     return nil
}






func (node FbsInvertedMapNum) Key() int64 {
    if node.VTable[0] == 0 {
        return int64(0)
    }
    return int64(flatbuffers.GetInt64(node.ValueNormal(0)))
}


func (node FbsInvertedMapNum) Value() FbsRecord {    
   if node.VTable[1] == 0 {
        return FbsRecord{}  
    }
    return FbsRecord{Node: node.ValueStruct(1)}

}




func (node FbsInvertedMapNum) CountOfField() int {
    return 2
}