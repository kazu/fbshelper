
// Code generated by genmaps.go; DO NOT EDIT.
// template file is https://github.com/kazu/fbshelper/blob/master/template/query.go.tmpl github.com/kazu/fbshelper/template/query.go.tmpl 
//   https://github.com/kazu/fbshelper/blob/master/template/union.query.go.tmpl


package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
    "reflect"
)

const (
    DUMMY_Symbol = flatbuffers.VtableMetadataFields
)

const (
        Symbol_Key =     0
)

var Symbol_FieldEnum = map[string]int{
        "Key": Symbol_Key,
}



type FbsSymbol struct {
	*base.Node
}


type FbsSymbolKey struct {
    *base.NodeList
}
func (node FbsSymbol) Info() base.Info {

    info := base.Info{Pos: node.Pos, Size: -1}
    for i := 0; i < len(node.VTable); i++ {
        vInfo := node.ValueInfo(i)
        if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
            info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
        }
    }
    return info    
}


func (node FbsSymbol) IsLeafAt(i int) bool {
    switch i {
    case 0:
        return false
    }
    return false
}


func (node FbsSymbol) ValueInfo(i int) base.ValueInfo {

    switch i {
    case 0:
         if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPosList(i)
        }
        node.ValueInfos[i].Size = node.Key().Info().Size
     }
     return node.ValueInfos[i]
}

func (node FbsSymbol) FieldAt(i int) interface{} {

    switch i {
    case 0:
        return node.Key()
     }
     return nil
}


// Unmarsla parse flatbuffers data and store the result
// in the value point to by v, if v is ni or not pointer,
// Unmarshal returns an ERR_MUST_POINTER, ERR_INVALID_TYPE
func (node FbsSymbol) Unmarshal(v interface{}) error {

    return node.Node.Unmarshal(v, func(s string, rv reflect.Value) error {
        
        switch Symbol_FieldEnum[s] {
        }
        return nil
    })

}




func (node FbsSymbol) Key() FbsSymbolKey {
    if node.VTable[0] == 0 {
        return FbsSymbolKey{}
    }
    nodelist :=  node.ValueList(0)
    return FbsSymbolKey{
                NodeList: &nodelist,
    }
}




func (node FbsSymbol) CountOfField() int {
    return 1
}


type Fbsbytes []byte    
            
            


func (node FbsSymbolKey) At(i int) Fbsbytes {
    if i >= int(node.ValueInfo.VLen) || i < 0 {
		return Fbsbytes{}
	}

    ptr := int(node.ValueInfo.Pos) + i*4
    return Fbsbytes(base.FbsString(base.NewNode(node.Base, ptr+ int(flatbuffers.GetUint32( node.R(ptr) )))))
}


func (node FbsSymbolKey) First() Fbsbytes {
	return node.At(0)
}


func (node FbsSymbolKey) Last() Fbsbytes {
	return node.At(int(node.ValueInfo.VLen)-1)
}

func (node FbsSymbolKey) Select(fn func(m Fbsbytes) bool) []Fbsbytes {

	result := make([]Fbsbytes, 0, int(node.ValueInfo.VLen))
	for i := 0; i < int(node.ValueInfo.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node FbsSymbolKey) Find(fn func(m Fbsbytes) bool) Fbsbytes{

	for i := 0; i < int(node.ValueInfo.VLen); i++ {
		if m := node.At(i); fn(m) {
			return m
		}
	}
	return Fbsbytes{}
}

func (node FbsSymbolKey) All() []Fbsbytes {
	return node.Select(func(m Fbsbytes) bool { return true })
}

func (node FbsSymbolKey) Count() int {
	return int(node.ValueInfo.VLen)
}

func (node FbsSymbolKey) Info() base.Info {

    info := base.Info{Pos: node.ValueInfo.Pos, Size: -1}
	ptr := int(node.ValueInfo.Pos) + (int(node.ValueInfo.VLen)-1)*4

    vInfo := base.FbsStringInfo(base.NewNode(node.Base, ptr+   int(flatbuffers.GetUint32( node.R(ptr) ))))



    if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
        info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
    }
    return info
}
