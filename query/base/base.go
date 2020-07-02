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

func FbsString(node *Node) []byte {
	buf := node.Bytes
	pos := uint32(node.Pos + int(node.VTable[0]))
	sLenOff := flatbuffers.GetUint32(buf[pos:])
	sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])
	start := pos + sLenOff + flatbuffers.SizeUOffsetT

	return buf[start : start+sLen]
}
