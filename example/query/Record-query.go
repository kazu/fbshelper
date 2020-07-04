
// Code generated by genmaps.go; DO NOT EDIT.
// template file is https://github.com/kazu/fbshelper/blob/master/template/query.go.tmpl github.com/kazu/fbshelper/template/query.go.tmpl 
//   https://github.com/kazu/fbshelper/blob/master/template/union.query.go.tmpl


package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

const (
    DUMMY_Record = flatbuffers.VtableMetadataFields
)

const (
        Record_FileId =     0
        Record_Offset =     1
        Record_Size =     2
        Record_OffsetOfValue =     3
        Record_ValueSize =     4
)


type FbsRecord struct {
	*base.Node
}


func (node FbsRecord) Info() base.Info {
    info := base.Info{Pos: node.Pos, Size: -1}
    size := 0
        size += base.SizeOfuint64
        size += base.SizeOfint64
        size += base.SizeOfint64
        size += base.SizeOfint32
        size += base.SizeOfint32
    info.Size = size

    return info

}

func (node FbsRecord) ValueInfo(i int) base.ValueInfo {

    switch i {
    case 0:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOfuint64
    case 1:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOfint64
    case 2:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOfint64
    case 3:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOfint32
    case 4:
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOfint32
     }
     return node.ValueInfos[i]
}

func (node FbsRecord) FieldAt(i int) interface{} {

    switch i {
    case 0:
        return node.FileId()
    case 1:
        return node.Offset()
    case 2:
        return node.Size()
    case 3:
        return node.OffsetOfValue()
    case 4:
        return node.ValueSize()
     }
     return nil
}






func (node FbsRecord) FileId() uint64 {
    buf := node.Bytes
    pos := node.Pos
    return uint64(flatbuffers.GetUint64(buf[pos:]))
}


func (node FbsRecord) Offset() int64 {
    buf := node.Bytes
    pos := node.Pos
                pos += base.SizeOfuint64
    return int64(flatbuffers.GetInt64(buf[pos:]))
}


func (node FbsRecord) Size() int64 {
    buf := node.Bytes
    pos := node.Pos
                pos += base.SizeOfuint64
                pos += base.SizeOfint64
    return int64(flatbuffers.GetInt64(buf[pos:]))
}


func (node FbsRecord) OffsetOfValue() int32 {
    buf := node.Bytes
    pos := node.Pos
                pos += base.SizeOfuint64
                pos += base.SizeOfint64
                pos += base.SizeOfint64
    return int32(flatbuffers.GetInt32(buf[pos:]))
}


func (node FbsRecord) ValueSize() int32 {
    buf := node.Bytes
    pos := node.Pos
                pos += base.SizeOfuint64
                pos += base.SizeOfint64
                pos += base.SizeOfint64
                pos += base.SizeOfint32
    return int32(flatbuffers.GetInt32(buf[pos:]))
}




func (node FbsRecord) CountOfField() int {
    return 5
}