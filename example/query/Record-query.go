
package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsRecord struct {
	*base.Node
}


func (node FbsRecord) FileId() uint64 {
	        if node.VTable[0] == 0 {    
		        return uint64(0)
	        }
            pos := node.Pos + int(node.VTable[0])
            return uint64(flatbuffers.GetUint64(node.Bytes[pos:]))
}
func (node FbsRecord) Offset() int64 {
	        if node.VTable[1] == 0 {    
		        return int64(0)
	        }
            pos := node.Pos + int(node.VTable[1])
            return int64(flatbuffers.GetInt64(node.Bytes[pos:]))
}
func (node FbsRecord) Size() int64 {
	        if node.VTable[2] == 0 {    
		        return int64(0)
	        }
            pos := node.Pos + int(node.VTable[2])
            return int64(flatbuffers.GetInt64(node.Bytes[pos:]))
}
func (node FbsRecord) OffsetOfValue() int32 {
	        if node.VTable[3] == 0 {    
		        return int32(0)
	        }
            pos := node.Pos + int(node.VTable[3])
            return int32(flatbuffers.GetInt32(node.Bytes[pos:]))
}
func (node FbsRecord) ValueSize() int32 {
	        if node.VTable[4] == 0 {    
		        return int32(0)
	        }
            pos := node.Pos + int(node.VTable[4])
            return int32(flatbuffers.GetInt32(node.Bytes[pos:]))
}