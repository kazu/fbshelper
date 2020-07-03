
package vfs_schema

import (
	//flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type FbsIndex struct {
    *base.Node
}
func(node FbsIndex) IndexNum() FbsIndexNum {

    return FbsIndexNum{Node: node.Node}
}
func(node FbsIndex) IndexString() FbsIndexString {

    return FbsIndexString{Node: node.Node}
}
func(node FbsIndex) File() FbsFile {

    return FbsFile{Node: node.Node}
}
func(node FbsIndex) InvertedMapNum() FbsInvertedMapNum {

    return FbsInvertedMapNum{Node: node.Node}
}
func(node FbsIndex) InvertedMapString() FbsInvertedMapString {

    return FbsInvertedMapString{Node: node.Node}
}