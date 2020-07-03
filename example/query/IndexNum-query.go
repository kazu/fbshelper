
package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsIndexNum struct {
	*base.Node
}


type FbsIndexNumMaps struct {
    *base.Node
    VPos   uint32
	VLen   uint32
	VStart uint32
}
func (node FbsIndexNum) Size() int32 {
	        if node.VTable[0] == 0 {    
		        return int32(0)
	        }
            pos := node.Pos + int(node.VTable[0])
            return int32(flatbuffers.GetInt32(node.Bytes[pos:]))
}
//func (node FbsIndexNum) Maps() {
func (node FbsIndexNum) Maps() FbsIndexNumMaps {
	        if node.VTable[1] == 0 {
                return FbsIndexNumMaps{}
	        }
            buf := node.Bytes
            vPos := uint32(node.Pos + int(node.VTable[1]))
            vLenOff := flatbuffers.GetUint32(buf[vPos:])
            vLen := flatbuffers.GetUint32(buf[vPos+vLenOff:])
            start := vPos + vLenOff + flatbuffers.SizeUOffsetT

            return FbsIndexNumMaps{
                Node: base.NewNode(node.Base, node.Pos),
                VPos:   vPos,
		        VLen:   vLen,
		        VStart: start,
            }
}




func (node FbsIndexNumMaps) At(i int) FbsInvertedMapNum {
    if i > int(node.VLen) || i < 0 {
		return FbsInvertedMapNum{}
	}

	buf := node.Bytes
	ptr := node.VStart + uint32(i-1)*4
	return FbsInvertedMapNum{Node: base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:])))}
}


func (node FbsIndexNumMaps) First() FbsInvertedMapNum {
	return node.At(0)
}


func (node FbsIndexNumMaps) Last() FbsInvertedMapNum {
	return node.At(int(node.VLen))
}

func (node FbsIndexNumMaps) Select(fn func(m FbsInvertedMapNum) bool) []FbsInvertedMapNum {

	result := make([]FbsInvertedMapNum, 0, int(node.VLen))
	for i := 0; i < int(node.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node FbsIndexNumMaps) Find(fn func(m FbsInvertedMapNum) bool) FbsInvertedMapNum{

	for i := 0; i < int(node.VLen); i++ {
		if m := node.At(i); fn(m) {
			return m
		}
	}
	return FbsInvertedMapNum{}
}

func (node FbsIndexNumMaps) All() []FbsInvertedMapNum {
	return node.Select(func(m FbsInvertedMapNum) bool { return true })
}

func (node FbsIndexNumMaps) Count() int {
	return int(node.VLen)
}