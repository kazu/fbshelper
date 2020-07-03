
package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsIndexString struct {
	*base.Node
}


type FbsIndexStringMaps struct {
    *base.Node
    VPos   uint32
	VLen   uint32
	VStart uint32
}
func (node FbsIndexString) Size() int32 {
	        if node.VTable[0] == 0 {    
		        return int32(0)
	        }
            pos := node.Pos + int(node.VTable[0])
            return int32(flatbuffers.GetInt32(node.Bytes[pos:]))
}
//func (node FbsIndexString) Maps() {
func (node FbsIndexString) Maps() FbsIndexStringMaps {
	        if node.VTable[1] == 0 {
                return FbsIndexStringMaps{}
	        }
            buf := node.Bytes
            vPos := uint32(node.Pos + int(node.VTable[1]))
            vLenOff := flatbuffers.GetUint32(buf[vPos:])
            vLen := flatbuffers.GetUint32(buf[vPos+vLenOff:])
            start := vPos + vLenOff + flatbuffers.SizeUOffsetT

            return FbsIndexStringMaps{
                Node: base.NewNode(node.Base, node.Pos),
                VPos:   vPos,
		        VLen:   vLen,
		        VStart: start,
            }
}




func (node FbsIndexStringMaps) At(i int) FbsInvertedMapString {
    if i > int(node.VLen) || i < 0 {
		return FbsInvertedMapString{}
	}

	buf := node.Bytes
	ptr := node.VStart + uint32(i-1)*4
	return FbsInvertedMapString{Node: base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:])))}
}


func (node FbsIndexStringMaps) First() FbsInvertedMapString {
	return node.At(0)
}


func (node FbsIndexStringMaps) Last() FbsInvertedMapString {
	return node.At(int(node.VLen))
}

func (node FbsIndexStringMaps) Select(fn func(m FbsInvertedMapString) bool) []FbsInvertedMapString {

	result := make([]FbsInvertedMapString, 0, int(node.VLen))
	for i := 0; i < int(node.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node FbsIndexStringMaps) Find(fn func(m FbsInvertedMapString) bool) FbsInvertedMapString{

	for i := 0; i < int(node.VLen); i++ {
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
	return int(node.VLen)
}