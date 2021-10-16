// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

type FileList struct { // genny
	*CommonNode
}

// File genny
func NewFileList() *FileList {

	list := emptyFileList()
	list.NodeList = &base.NodeList{}
	list.CommonNode.Name = "[]File"

	(*base.List)(list.CommonNode).InitList()
	return list
}

func emptyFileList() *FileList {
	return &FileList{CommonNode: &base.CommonNode{}}
}

func (node FileList) At(i int) (result *File, e error) {
	result = &File{}
	result.CommonNode, e = (*base.List)(node.CommonNode).At(i)
	return
}

func (node FileList) AtWihoutError(i int) (result *File) {
	result, e := node.At(i)
	if e != nil {
		result = nil
	}
	return
}

func (node FileList) SetAt(i int, v *File) error {
	return (*base.List)(node.CommonNode).SetAt(i, v.CommonNode)
}
func (node FileList) Add(v FileList) error {
	return (*base.List)(node.CommonNode).Add((*base.List)(v.CommonNode))
}

func (node FileList) First() (result *File, e error) {
	return node.At(0)
}

func (node FileList) Last() (result *File, e error) {
	return node.At(int(node.NodeList.ValueInfo.VLen) - 1)
}

func (node FileList) Select(fn func(*File) bool) (result []*File) {
	result = make([]*File, 0, int(node.NodeList.ValueInfo.VLen))
	commons := (*base.List)(node.CommonNode).Select(func(cm *CommonNode) bool {
		return fn(&File{CommonNode: cm})
	})
	for _, cm := range commons {
		result = append(result, &File{CommonNode: cm})
	}
	return result
}

func (node FileList) Find(fn func(*File) bool) *File {
	result := &File{}
	result.CommonNode = (*base.List)(node.CommonNode).Find(func(cm *CommonNode) bool {
		return fn(&File{CommonNode: cm})
	})
	return result
}

func (node FileList) All() []*File {
	return node.Select(func(*File) bool { return true })
}

func (node FileList) Count() int {
	return int(node.NodeList.ValueInfo.VLen)
}

func (node FileList) SwapAt(i, j int) error {
	return (*List)(node.CommonNode).SwapAt(i, j)
}

func (node FileList) SortBy(less func(i, j int) bool) error {
	return (*List)(node.CommonNode).SortBy(less)
}

// Search ... binary search
func (node FileList) Search(fn func(*File) bool) *File {
	result := &File{}

	i := (*base.List)(node.CommonNode).SearchIndex(int((*base.List)(node.CommonNode).VLen()), func(cm *CommonNode) bool {
		return fn(&File{CommonNode: cm})
	})
	if i < int((*base.List)(node.CommonNode).VLen()) {
		result, _ = node.At(i)
	}

	return result
}

func (node FileList) SearchIndex(fn func(*File) bool) int {

	i := (*base.List)(node.CommonNode).SearchIndex(int((*base.List)(node.CommonNode).VLen()), func(cm *CommonNode) bool {
		return fn(&File{CommonNode: cm})
	})
	if i < int((*base.List)(node.CommonNode).VLen()) {
		return i
	}

	return -1
}
