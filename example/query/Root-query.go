
// Code generated by genmaps.go; DO NOT EDIT.
// template file is https://github.com/kazu/fbshelper/blob/master/template/query.go.tmpl github.com/kazu/fbshelper/template/query.go.tmpl 
//   https://github.com/kazu/fbshelper/blob/master/template/union.query.go.tmpl


package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
    "reflect"
    "io"
)

const (
    DUMMY_Root = flatbuffers.VtableMetadataFields
)

const (
        Root_Version =     0
        Root_IndexType =     1
        Root_Index =     2
)

var Root_FieldEnum = map[string]int{
        "Version": Root_Version,
        "IndexType": Root_IndexType,
        "Index": Root_Index,
}



type FbsRoot struct {
	*base.Node
}


type FbsRootRootIndex struct {
    *base.Node
}
func Open(r io.Reader, cap int) FbsRoot {
    b := base.NewBaseByIO(r, 512)
    
    return FbsRoot{
		Node: base.NewNode(b, int(flatbuffers.GetUOffsetT( b.R(0) ))),
	}
}


func OpenByBuf(buf []byte) FbsRoot {
	return FbsRoot{
		Node: base.NewNode(base.NewBase(buf), int(flatbuffers.GetUOffsetT(buf))),
	}
}

func (node FbsRoot) Len() int {
    info := node.Info()
    size := info.Pos + info.Size

    if (size % 8) == 0 {
        return size
    }

    return size + (8 - (size % 8)) 
}

func (node FbsRoot) Next() FbsRoot {
    start := node.Len()

    if node.LenBuf() + 4 < start {
        return node
    }

    return FbsRoot{
		Node: base.NewNode(node.Base, start + int(flatbuffers.GetUOffsetT( node.R(start)  ))),
	}
}

func (node FbsRoot) HasNext() bool {

    return node.LenBuf() + 4 < node.Len()
}
func (node FbsRoot) Info() base.Info {

    info := base.Info{Pos: node.Pos, Size: -1}
    for i := 0; i < len(node.VTable); i++ {
        vInfo := node.ValueInfo(i)
        if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
            info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
        }
    }
    return info    
}

func (node FbsRoot) ValueInfo(i int) base.ValueInfo {

    switch i {
    case 0:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOfint32
    case 1:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOfbyte
    case 2:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPosTable(i)
        }
        //node.ValueInfos[i].Size = node.Index().Info(i-1).Size
        // 
        // typeIdx := node.FieldAt(i-1).(EnumIndex)
        // node.ValueInfos[i].Size = node.Index().Info(typeIdx).Size
        
        
        eIdx := int(node.IndexType())
        node.ValueInfos[i].Size = node.Index().Info(eIdx).Size
     }
     return node.ValueInfos[i]
}

func (node FbsRoot) FieldAt(i int) interface{} {

    switch i {
    case 0:
        return node.Version()
    case 1:
        return node.IndexType()
    case 2:
        return node.Index()
     }
     return nil
}


// Unmarsla parse flatbuffers data and store the result
// in the value point to by v, if v is ni or not pointer,
// Unmarshal returns an ERR_MUST_POINTER, ERR_INVALID_TYPE
func (node FbsRoot) Unmarshal(v interface{}) error {

    return node.Node.Unmarshal(v, func(s string, rv reflect.Value) error {
        
        switch Root_FieldEnum[s] {
        case Root_Version:
            //return node.Version()
            rv.Set(reflect.ValueOf(  node.Version() ))
        case Root_IndexType:
            //return node.IndexType()
            rv.Set(reflect.ValueOf(  node.IndexType() ))
        }
        return nil
    })

}




func (node FbsRoot) Version() int32 {
    if node.VTable[0] == 0 {
        return int32(0)
    }
    return int32(flatbuffers.GetInt32(node.ValueNormal(0)))
}


func (node FbsRoot) IndexType() byte {
    if node.VTable[1] == 0 {
        return byte(0)
    }
    return byte(flatbuffers.GetByte(node.ValueNormal(1)))
}



func (node FbsRoot) Index() FbsIndex {
    if node.VTable[2] == 0 {
        return FbsIndex{}  
    }
    return FbsIndex{Node: node.ValueTable(2)}
}




func (node FbsRoot) CountOfField() int {
    return 3
}
