// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package vfs_schema

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type IndexString struct {
	_tab flatbuffers.Table
}

func GetRootAsIndexString(buf []byte, offset flatbuffers.UOffsetT) *IndexString {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &IndexString{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *IndexString) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *IndexString) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *IndexString) Size() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IndexString) MutateSize(n int32) bool {
	return rcv._tab.MutateInt32Slot(4, n)
}

func (rcv *IndexString) Maps(obj *InvertedMapString, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *IndexString) MapsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func IndexStringStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func IndexStringAddSize(builder *flatbuffers.Builder, size int32) {
	builder.PrependInt32Slot(0, size, 0)
}
func IndexStringAddMaps(builder *flatbuffers.Builder, maps flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(maps), 0)
}
func IndexStringStartMapsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func IndexStringEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
