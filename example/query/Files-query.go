
// Code generated by genmaps.go; DO NOT EDIT.
// template file is https://github.com/kazu/fbshelper/blob/master/template/query.go.tmpl github.com/kazu/fbshelper/template/query.go.tmpl 
//   https://github.com/kazu/fbshelper/blob/master/template/union.query.go.tmpl


package vfs_schema

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

const (
    DUMMY_Files = flatbuffers.VtableMetadataFields
)

type FbsFiles struct {
	*base.Node
}


type FbsFilesFiles struct {
    *base.NodeList
}
func (node FbsFiles) Info() base.Info {

    info := base.Info{Pos: node.Pos, Size: -1}
    for i := 0; i < len(node.VTable); i++ {
        vInfo := node.ValueInfo(i)
        if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
            info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
        }
    }
    return info    
}

func (node FbsFiles) ValueInfo(i int) base.ValueInfo {

    switch i {
    case 0:
         if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPosList(i)
        }
        node.ValueInfos[i].Size = node.Files().Info().Size
     }
     return node.ValueInfos[i]
}





func (node FbsFiles) Files() FbsFilesFiles {
    if node.VTable[0] == 0 {
        return FbsFilesFiles{}
    }
    nodelist :=  node.ValueList(0)
    return FbsFilesFiles{
                NodeList: &nodelist,
    }
}






func (node FbsFilesFiles) At(i int) FbsFile {
    if i > int(node.ValueInfo.VLen) || i < 0 {
		return FbsFile{}
	}

	buf := node.Bytes
	ptr := uint32(node.ValueInfo.Pos + (i-1)*4)
	return FbsFile{Node: base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:])))}
}


func (node FbsFilesFiles) First() FbsFile {
	return node.At(0)
}


func (node FbsFilesFiles) Last() FbsFile {
	return node.At(int(node.ValueInfo.VLen))
}

func (node FbsFilesFiles) Select(fn func(m FbsFile) bool) []FbsFile {

	result := make([]FbsFile, 0, int(node.ValueInfo.VLen))
	for i := 0; i < int(node.ValueInfo.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node FbsFilesFiles) Find(fn func(m FbsFile) bool) FbsFile{

	for i := 0; i < int(node.ValueInfo.VLen); i++ {
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
	return int(node.ValueInfo.VLen)
}

func (node FbsFilesFiles) Info() base.Info {

    info := base.Info{Pos: node.ValueInfo.Pos, Size: -1}
    vInfo := node.Last().Info()



    if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
        info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
    }
    return info
}
