
package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsSymbols struct {
	*base.Node
}


type FbsSymbolsSymbols struct {
    *base.Node
    VPos   uint32
	VLen   uint32
	VStart uint32
}
//func (node FbsSymbols) Symbols() {
func (node FbsSymbols) Symbols() FbsSymbolsSymbols {
	        if node.VTable[0] == 0 {
                return FbsSymbolsSymbols{}
	        }
            buf := node.Bytes
            vPos := uint32(node.Pos + int(node.VTable[0]))
            vLenOff := flatbuffers.GetUint32(buf[vPos:])
            vLen := flatbuffers.GetUint32(buf[vPos+vLenOff:])
            start := vPos + vLenOff + flatbuffers.SizeUOffsetT

            return FbsSymbolsSymbols{
                Node: base.NewNode(node.Base, node.Pos),
                VPos:   vPos,
		        VLen:   vLen,
		        VStart: start,
            }
}




func (node FbsSymbolsSymbols) At(i int) FbsSymbol {
    if i > int(node.VLen) || i < 0 {
		return FbsSymbol{}
	}

	buf := node.Bytes
	ptr := node.VStart + uint32(i-1)*4
	return FbsSymbol{Node: base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:])))}
}


func (node FbsSymbolsSymbols) First() FbsSymbol {
	return node.At(0)
}


func (node FbsSymbolsSymbols) Last() FbsSymbol {
	return node.At(int(node.VLen))
}

func (node FbsSymbolsSymbols) Select(fn func(m FbsSymbol) bool) []FbsSymbol {

	result := make([]FbsSymbol, 0, int(node.VLen))
	for i := 0; i < int(node.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node FbsSymbolsSymbols) Find(fn func(m FbsSymbol) bool) FbsSymbol{

	for i := 0; i < int(node.VLen); i++ {
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
	return int(node.VLen)
}