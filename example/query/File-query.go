
// Code generated by genmaps.go; DO NOT EDIT.
// template file is https://github.com/kazu/fbshelper/blob/master/template/query.go.tmpl github.com/kazu/fbshelper/template/query.go.tmpl 
//   https://github.com/kazu/fbshelper/blob/master/template/union.query.go.tmpl


package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsFile struct {
	*base.Node
}


func (node FbsFile) Id() uint64 {
	        if node.VTable[0] == 0 {    
		        return uint64(0)
	        }
            pos := node.Pos + int(node.VTable[0])
            return uint64(flatbuffers.GetUint64(node.Bytes[pos:]))
}
func (node FbsFile) Name() []byte {
	        if node.VTable[1] == 0 {
                return nil
	        }
            buf := node.Bytes
	        pos := uint32(node.Pos + int(node.VTable[1]))
	        sLenOff := flatbuffers.GetUint32(buf[pos:])
	        sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])
	        start := pos + sLenOff + flatbuffers.SizeUOffsetT

            return buf[start:start+sLen]
}
func (node FbsFile) IndexAt() int64 {
	        if node.VTable[2] == 0 {    
		        return int64(0)
	        }
            pos := node.Pos + int(node.VTable[2])
            return int64(flatbuffers.GetInt64(node.Bytes[pos:]))
}