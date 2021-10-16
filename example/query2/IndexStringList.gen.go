// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

type IndexStringList struct { // genny
	*CommonNode
}

// IndexString genny
func NewIndexStringList() *IndexStringList {

	list := emptyIndexStringList()
	list.NodeList = &base.NodeList{}
	list.CommonNode.Name = "[]IndexString"

	(*base.List)(list.CommonNode).InitList()
	return list
}

func emptyIndexStringList() *IndexStringList {
	return &IndexStringList{CommonNode: &base.CommonNode{}}
}

func (node IndexStringList) At(i int) (result *IndexString, e error) {
	result = &IndexString{}
	result.CommonNode, e = (*base.List)(node.CommonNode).At(i)
	return
}

func (node IndexStringList) SetAt(i int, v *IndexString) error {
	return (*base.List)(node.CommonNode).SetAt(i, v.CommonNode)
}

func (node IndexStringList) First() (result *IndexString, e error) {
	return node.At(0)
}

func (node IndexStringList) Last() (result *IndexString, e error) {
	return node.At(int(node.NodeList.ValueInfo.VLen) - 1)
}

func (node IndexStringList) Select(fn func(*IndexString) bool) (result []*IndexString) {
	result = make([]*IndexString, 0, int(node.NodeList.ValueInfo.VLen))
	commons := (*base.List)(node.CommonNode).Select(func(cm *CommonNode) bool {
		return fn(&IndexString{CommonNode: cm})
	})
	for _, cm := range commons {
		result = append(result, &IndexString{CommonNode: cm})
	}
	return result
}

func (node IndexStringList) Find(fn func(*IndexString) bool) *IndexString {
	result := &IndexString{}
	result.CommonNode = (*base.List)(node.CommonNode).Find(func(cm *CommonNode) bool {
		return fn(&IndexString{CommonNode: cm})
	})
	return result
}

func (node IndexStringList) All() []*IndexString {
	return node.Select(func(*IndexString) bool { return true })
}

func (node IndexStringList) Count() int {
	return int(node.NodeList.ValueInfo.VLen)
}

func (node IndexStringList) SwapAt(i, j int) error {
	return (*List)(node.CommonNode).SwapAt(i, j)
}

func (node IndexStringList) SortBy(less func(i, j int) bool) error {
	return (*List)(node.CommonNode).SortBy(less)
}

// Search ... binary search
func (node IndexStringList) Search(fn func(*IndexString) bool) *IndexString {
	result := &IndexString{}

	i := (*base.List)(node.CommonNode).SearchIndex(int((*base.List)(node.CommonNode).VLen()), func(cm *CommonNode) bool {
		return fn(&IndexString{CommonNode: cm})
	})
	if i < int((*base.List)(node.CommonNode).VLen()) {
		result, _ = node.At(i)
	}

	return result
}

func (node IndexStringList) SearchIndex(fn func(*IndexString) bool) int {

	i := (*base.List)(node.CommonNode).SearchIndex(int((*base.List)(node.CommonNode).VLen()), func(cm *CommonNode) bool {
		return fn(&IndexString{CommonNode: cm})
	})
	if i < int((*base.List)(node.CommonNode).VLen()) {
		return i
	}

	return -1
}
