
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
    DUMMY_Symbols = flatbuffers.VtableMetadataFields
)

const (
        Symbols_Symbols =     0
)

var Symbols_FieldEnum = map[string]int{
        "Symbols": Symbols_Symbols,
}



type FbsSymbols struct {
	*base.Node
}


type FbsSymbolsSymbols struct {
    *base.NodeList
}

func (node FbsSymbols) SearchInfo(pos int, fn RecFn, condFn CondFn) {

	info := node.Info()

	if condFn(pos, info) {
		fn(base.NodePath{Name: "Symbols", Idx: -1}, info)
	}else{
        return
    }

	for i := 0; i < node.CountOfField(); i++ {
		if node.IsLeafAt(i) {
			fInfo := base.Info(node.ValueInfo(i))
			if condFn(pos, fInfo) {
				fn(base.NodePath{Name: "Symbols", Idx: i}, fInfo)
			}
			continue
		}
        switch i {
        case 0:
                node.Symbols().SearchInfo(pos, fn, condFn)    
        default:
			base.Log(base.LOG_ERROR, func() base.LogArgs {
				return F("node must be Noder")
			})
        }

	}

}
func (node FbsSymbols) Info() base.Info {

    info := base.Info{Pos: node.Pos, Size: -1}
    for i := 0; i < len(node.VTable); i++ {
        vInfo := node.ValueInfo(i)
        if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
            info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
        }
    }
    return info    
}


func (node FbsSymbols) IsLeafAt(i int) bool {
    switch i {
    case 0:
        return false
    }
    return false
}
func (node FbsSymbols) ValueInfo(i int) base.ValueInfo {

    switch i {
    case 0:
         if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPosList(i)
        }
        node.ValueInfos[i].Size = node.Symbols().Info().Size
     }
     return node.ValueInfos[i]
}


func (node FbsSymbols) FieldAt(i int) interface{} {

    switch i {
    case 0:
        return node.Symbols()
     }
     return nil
}


// Unmarsla parse flatbuffers data and store the result
// in the value point to by v, if v is ni or not pointer,
// Unmarshal returns an ERR_MUST_POINTER, ERR_INVALID_TYPE
func (node FbsSymbols) Unmarshal(v interface{}) error {

    return node.Node.Unmarshal(v, func(s string, rv reflect.Value) error {
        
        switch Symbols_FieldEnum[s] {
        }
        return nil
    })

}




func (node FbsSymbols) Symbols() FbsSymbolsSymbols {
    if node.VTable[0] == 0 {
        return FbsSymbolsSymbols{}
    }
    nodelist :=  node.ValueList(0)
    return FbsSymbolsSymbols{
                NodeList: &nodelist,
    }
}




func (node FbsSymbols) CountOfField() int {
    return 1
}



        
         
               

func (node FbsSymbolsSymbols) At(i int) FbsSymbol {
    if i >= int(node.ValueInfo.VLen) || i < 0 {
        return FbsSymbol{}    
	}

    ptr := int(node.ValueInfo.Pos) + i*4 
	return FbsSymbol{Node: base.NewNode(node.Base, ptr + int(flatbuffers.GetUint32( node.R(ptr) )))}
}


func (node FbsSymbolsSymbols) First() FbsSymbol {
	return node.At(0)
}


func (node FbsSymbolsSymbols) Last() FbsSymbol {
	return node.At(int(node.ValueInfo.VLen)-1)
}

func (node FbsSymbolsSymbols) Select(fn func(m FbsSymbol) bool) []FbsSymbol {

	result := make([]FbsSymbol, 0, int(node.ValueInfo.VLen))
	for i := 0; i < int(node.ValueInfo.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node FbsSymbolsSymbols) Find(fn func(m FbsSymbol) bool) FbsSymbol{

	for i := 0; i < int(node.ValueInfo.VLen); i++ {
		if m := node.At(i); fn(m) {
			return m
		}
	}
    return FbsSymbol{}   

}

func (node FbsSymbolsSymbols) All() []FbsSymbol {
	return node.Select(func(m FbsSymbol) bool { return true })
}

func (node FbsSymbolsSymbols) Count() int {
	return int(node.ValueInfo.VLen)
}

func (node FbsSymbolsSymbols) Info() base.Info {

    info := base.Info{Pos: node.ValueInfo.Pos, Size: -1}
    vInfo := node.Last().Info()



    if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
        info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
    }
    return info
}

func (node FbsSymbolsSymbols) SearchInfo(pos int, fn RecFn, condFn CondFn) {

    info := node.Info()

    if condFn(pos, info) {
        fn(base.NodePath{Name: "Symbols.Symbols", Idx: -1}, info)
	}else{
        return
    }

    var v interface{}
    for _, cNode := range node.All() {
        v = cNode
        if vv, ok := v.(Searcher); ok {    
                vv.SearchInfo(pos, fn, condFn)
        }else{
            goto NO_NODE
        }
    }
    return
    

NO_NODE:     
    for i := 0 ; i < int(node.ValueInfo.VLen) ; i++ {
        ptr := int(node.ValueInfo.Pos) + i*4
        start := ptr + int(flatbuffers.GetUint32( node.R(ptr) ))
        size := info.Size
        if i + 1 < int(node.ValueInfo.Pos) {
            size = ptr+4 + int(flatbuffers.GetUint32( node.R(ptr+4) )) - start
        }
        cInfo := base.Info{Pos: start, Size: size}
        if condFn(pos, info) {
            fn(base.NodePath{Name: "Symbols.Symbols", Idx: i}, cInfo)
        }
    }
}
