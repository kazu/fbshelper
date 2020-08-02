// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

type NumListList struct { // genny
	*CommonNode
}

// NumList genny
func NewNumListList() *NumListList {

	list := emptyNumListList()
	list.NodeList = &base.NodeList{}
	list.CommonNode.Name = "[]NumList"

	list.InitList()
	return list
}

func emptyNumListList() *NumListList {
	return &NumListList{CommonNode: &base.CommonNode{}}
}

func (node NumListList) At(i int) (result *NumList, e error) {
	result = &NumList{}
	result.CommonNode, e = node.CommonNode.At(i)
	return
}

func (node NumListList) SetAt(i int, v *NumList) error {
	return node.CommonNode.SetAt(i, v.CommonNode)
}

func (node NumListList) First() (result *NumList, e error) {
	return node.At(0)
}

func (node NumListList) Last() (result *NumList, e error) {
	return node.At(int(node.NodeList.ValueInfo.VLen) - 1)
}

func (node NumListList) Select(fn func(*NumList) bool) (result []*NumList) {
	result = make([]*NumList, 0, int(node.NodeList.ValueInfo.VLen))
	commons := node.CommonNode.Select(func(cm *CommonNode) bool {
		return fn(&NumList{CommonNode: cm})
	})
	for _, cm := range commons {
		result = append(result, &NumList{CommonNode: cm})
	}
	return result
}

func (node NumListList) Find(fn func(*NumList) bool) *NumList {
	result := &NumList{}
	result.CommonNode = node.CommonNode.Find(func(cm *CommonNode) bool {
		return fn(&NumList{CommonNode: cm})
	})
	return result
}

func (node NumListList) All() []*NumList {
	return node.Select(func(*NumList) bool { return true })
}

func (node NumListList) Count() int {
	return int(node.NodeList.ValueInfo.VLen)
}

// Search ... binary search
func (node NumListList) Search(fn func(*NumList) bool) *NumList {
	result := &NumList{}

	i := node.CommonNode.SearchIndex(int(node.VLen()), func(cm *CommonNode) bool {
		return fn(&NumList{CommonNode: cm})
	})
	if i < int(node.VLen()) {
		result, _ = node.At(i)
	}

	return result
}
