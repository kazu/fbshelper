
package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsFiles struct {
	*base.Node
}


type FbsFilesFiles struct {
    *base.Node
    VPos   uint32
	VLen   uint32
	VStart uint32
}
//func (node FbsFiles) Files() {
func (node FbsFiles) Files() FbsFilesFiles {
	        if node.VTable[0] == 0 {
                return FbsFilesFiles{}
	        }
            buf := node.Bytes
            vPos := uint32(node.Pos + int(node.VTable[0]))
            vLenOff := flatbuffers.GetUint32(buf[vPos:])
            vLen := flatbuffers.GetUint32(buf[vPos+vLenOff:])
            start := vPos + vLenOff + flatbuffers.SizeUOffsetT

            return FbsFilesFiles{
                Node: base.NewNode(node.Base, node.Pos),
                VPos:   vPos,
		        VLen:   vLen,
		        VStart: start,
            }
}




func (node FbsFilesFiles) At(i int) FbsFile {
    if i > int(node.VLen) || i < 0 {
		return FbsFile{}
	}

	buf := node.Bytes
	ptr := node.VStart + uint32(i-1)*4
	return FbsFile{Node: base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:])))}
}


func (node FbsFilesFiles) First() FbsFile {
	return node.At(0)
}


func (node FbsFilesFiles) Last() FbsFile {
	return node.At(int(node.VLen))
}

func (node FbsFilesFiles) Select(fn func(m FbsFile) bool) []FbsFile {

	result := make([]FbsFile, 0, int(node.VLen))
	for i := 0; i < int(node.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node FbsFilesFiles) Find(fn func(m FbsFile) bool) FbsFile{

	for i := 0; i < int(node.VLen); i++ {
		if m := node.At(i); fn(m) {
			return m
		}
	}
	return FbsFile{}
}

func (node FbsFilesFiles) All() []FbsFile {
	return node.Select(func(m FbsFile) bool { return true })
}

func (node FbsFilesFiles) Count() int {
	return int(node.VLen)
}