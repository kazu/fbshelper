// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package vfs_schema

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Hoges struct {
	_tab flatbuffers.Table
}

func GetRootAsHoges(buf []byte, offset flatbuffers.UOffsetT) *Hoges {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Hoges{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Hoges) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Hoges) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Hoges) Files(obj *File, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Hoges) FilesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func HogesStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func HogesAddFiles(builder *flatbuffers.Builder, files flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(files), 0)
}
func HogesStartFilesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func HogesEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}