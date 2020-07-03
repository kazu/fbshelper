
package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsInvertedMapNum struct {
	*base.Node
}


func (node FbsInvertedMapNum) Key() int64 {
	        if node.VTable[0] == 0 {    
		        return int64(0)
	        }
            pos := node.Pos + int(node.VTable[0])
            return int64(flatbuffers.GetInt64(node.Bytes[pos:]))
}
func (node FbsInvertedMapNum) Value() FbsRecord {
	        if node.VTable[1] == 0 {
                return FbsRecord{}
	        }
            pos := node.Pos + int(node.VTable[1])
            return FbsRecord{Node: base.NewNode(node.Base, int(flatbuffers.GetUint32(node.Bytes[pos:]))+pos)}
}