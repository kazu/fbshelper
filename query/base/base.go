package base

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Base struct {
	Bytes []byte
	Diffs []Diff
}

type Diff struct {
	Offset int
	Bytes  []byte
}

type Root struct {
	*Node
}

type Node struct {
	*Base
	Pos    int
	VTable []uint16
	TLen   uint16
}

func NewNode(b *Base, pos int) *Node {
	node := &Node{Base: b, Pos: pos}
	node.vtable()
	return node
}

func (n *Node) vtable() {
	if len(n.VTable) > 0 {
		return
	}
	buf := n.Bytes
	vOffset := uint32(flatbuffers.GetUOffsetT(buf[n.Pos:]))
	vPos := uint32(n.Pos) - vOffset
	vLen := uint16(flatbuffers.GetVOffsetT(buf[vPos:]))
	n.TLen = uint16(flatbuffers.GetVOffsetT(buf[vPos+2:]))

	for cur := vPos + 4; cur < vPos+uint32(vLen); cur += 2 {
		n.VTable = append(n.VTable, uint16(flatbuffers.GetVOffsetT(buf[cur:])))
	}
}

// must be generated

type FbsRoot struct {
	*Node
}

type FbsRootIndex struct {
	*Node
}

type FbsFile struct {
	*Node
}

func OpenByBuf(buf []byte) *FbsRoot {
	return &FbsRoot{
		Node: NewNode(&Base{Bytes: buf}, int(flatbuffers.GetUOffsetT(buf))),
	}
}

func (node *FbsRoot) Version() int32 {
	if node.VTable[0] == 0 {
		return 0
	}
	pos := node.Pos + int(node.VTable[0])

	return flatbuffers.GetInt32(node.Bytes[pos:])
}

func (root *FbsRoot) Index() *FbsRootIndex {
	if root.VTable[2] == 0 {
		return nil
	}
	pos := root.Pos + int(root.VTable[2])
	return &FbsRootIndex{Node: NewNode(root.Base, int(flatbuffers.GetUint32(root.Bytes[pos:]))+pos)}
}

func (r *FbsRootIndex) File() *FbsFile {
	return &FbsFile{Node: r.Node}
}

func (node *FbsFile) Id() uint64 {
	if node.VTable[0] == 0 {
		return 0
	}
	pos := node.Pos + int(node.VTable[0])

	return uint64(flatbuffers.GetUint64(node.Bytes[pos:]))
}

func (node *FbsFile) Name() string {
	if node.VTable[1] == 0 {
		return ""
	}
	buf := node.Bytes
	pos := uint32(node.Pos + int(node.VTable[1]))
	sLenOff := flatbuffers.GetUint32(buf[pos:])
	sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])
	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	return string(buf[start : start+sLen])
}

func (node *FbsFile) IndexAt() int64 {
	if node.VTable[2] == 0 {
		return 0
	}
	pos := node.Pos + int(node.VTable[2])

	return flatbuffers.GetInt64(node.Bytes[pos:])
}
