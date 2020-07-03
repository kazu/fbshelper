
package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsRoot struct {
	*base.Node
}


type FbsRootRootIndex struct {
    *base.Node
}
func OpenByBuf(buf []byte) FbsRoot {
	return FbsRoot{
		Node: base.NewNode(&base.Base{Bytes: buf}, int(flatbuffers.GetUOffsetT(buf))),
	}
}
func (node FbsRoot) Version() int32 {
	        if node.VTable[0] == 0 {    
		        return int32(0)
	        }
            pos := node.Pos + int(node.VTable[0])
            return int32(flatbuffers.GetInt32(node.Bytes[pos:]))
}
func (node FbsRoot) IndexType() byte {
	        if node.VTable[1] == 0 {    
		        return byte(0)
	        }
            pos := node.Pos + int(node.VTable[1])
            return byte(flatbuffers.GetByte(node.Bytes[pos:]))
}
func (node FbsRoot) Index() FbsIndex {
	        if node.VTable[2] == 0 {
                return FbsIndex{}
	        }
            pos := node.Pos + int(node.VTable[2])
            return FbsIndex{Node: base.NewNode(node.Base, int(flatbuffers.GetUint32(node.Bytes[pos:]))+pos)}
}