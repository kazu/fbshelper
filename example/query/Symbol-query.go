
package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsSymbol struct {
	*base.Node
}


type FbsSymbolKey struct {
    *base.Node
    VPos   uint32
	VLen   uint32
	VStart uint32
}
//func (node FbsSymbol) Key() {
func (node FbsSymbol) Key() FbsSymbolKey {
	        if node.VTable[0] == 0 {
                return FbsSymbolKey{}
	        }
            buf := node.Bytes
            vPos := uint32(node.Pos + int(node.VTable[0]))
            vLenOff := flatbuffers.GetUint32(buf[vPos:])
            vLen := flatbuffers.GetUint32(buf[vPos+vLenOff:])
            start := vPos + vLenOff + flatbuffers.SizeUOffsetT

            return FbsSymbolKey{
                Node: base.NewNode(node.Base, node.Pos),
                VPos:   vPos,
		        VLen:   vLen,
		        VStart: start,
            }
}


type Fbsbytes []byte    
            
            


func (node FbsSymbolKey) At(i int) Fbsbytes {
    if i > int(node.VLen) || i < 0 {
		return Fbsbytes{}
	}

	buf := node.Bytes
	ptr := node.VStart + uint32(i-1)*4
    return Fbsbytes(base.FbsString(base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:])))))
}


func (node FbsSymbolKey) First() Fbsbytes {
	return node.At(0)
}


func (node FbsSymbolKey) Last() Fbsbytes {
	return node.At(int(node.VLen))
}

func (node FbsSymbolKey) Select(fn func(m Fbsbytes) bool) []Fbsbytes {

	result := make([]Fbsbytes, 0, int(node.VLen))
	for i := 0; i < int(node.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node FbsSymbolKey) Find(fn func(m Fbsbytes) bool) Fbsbytes{

	for i := 0; i < int(node.VLen); i++ {
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
	return int(node.VLen)
}