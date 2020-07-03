
package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsInvertedMapString struct {
	*base.Node
}


func (node FbsInvertedMapString) Key() []byte {
	        if node.VTable[0] == 0 {
                return nil
	        }
            buf := node.Bytes
	        pos := uint32(node.Pos + int(node.VTable[0]))
	        sLenOff := flatbuffers.GetUint32(buf[pos:])
	        sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])
	        start := pos + sLenOff + flatbuffers.SizeUOffsetT

            return buf[start:start+sLen]
}
func (node FbsInvertedMapString) Value() FbsRecord {
	        if node.VTable[1] == 0 {
                return FbsRecord{}
	        }
            pos := node.Pos + int(node.VTable[1])
            return FbsRecord{Node: base.NewNode(node.Base, int(flatbuffers.GetUint32(node.Bytes[pos:]))+pos)}
}