package query

import "github.com/kazu/fbshelper/query/base"

type ListType struct { // genny
	*CommonNode
}

// NodeName genny
func NewListType() *ListType {

	list := emptyListType()
	list.NodeList = &base.NodeList{}
	list.CommonNode.Name = "[]NodeName"

	(*base.List)(list.CommonNode).InitList()
	return list
}

func emptyListType() *ListType {
	return &ListType{CommonNode: &base.CommonNode{}}
}

func (node ListType) At(i int) (result *NodeName, e error) {
	result = &NodeName{}
	result.CommonNode, e = (*base.List)(node.CommonNode).At(i)
	return
}

func (node ListType) AtWihoutError(i int) (result *NodeName) {
	result, e := node.At(i)
	if e != nil {
		result = nil
	}
	return
}

func (node ListType) SetAt(i int, v *NodeName) error {
	return (*base.List)(node.CommonNode).SetAt(i, v.CommonNode)
}

func (node ListType) First() (result *NodeName, e error) {
	return node.At(0)
}

func (node ListType) Last() (result *NodeName, e error) {
	return node.At(int(node.NodeList.ValueInfo.VLen) - 1)
}

func (node ListType) Select(fn func(*NodeName) bool) (result []*NodeName) {
	result = make([]*NodeName, 0, int(node.NodeList.ValueInfo.VLen))
	commons := (*base.List)(node.CommonNode).Select(func(cm *CommonNode) bool {
		return fn(&NodeName{CommonNode: cm})
	})
	for _, cm := range commons {
		result = append(result, &NodeName{CommonNode: cm})
	}
	return result
}

func (node ListType) Find(fn func(*NodeName) bool) *NodeName {
	result := &NodeName{}
	result.CommonNode = (*base.List)(node.CommonNode).Find(func(cm *CommonNode) bool {
		return fn(&NodeName{CommonNode: cm})
	})
	return result
}

func (node ListType) All() []*NodeName {
	return node.Select(func(*NodeName) bool { return true })
}

func (node ListType) Count() int {
	return int(node.NodeList.ValueInfo.VLen)
}

func (node ListType) SwapAt(i, j int) error {
	return (*List)(node.CommonNode).SwapAt(i, j)
}

func (node ListType) SortBy(less func(i, j int) bool) error {
	return (*List)(node.CommonNode).SortBy(less)
}

// Search ... binary search
func (node ListType) Search(fn func(*NodeName) bool) *NodeName {
	result := &NodeName{}

	i := (*base.List)(node.CommonNode).SearchIndex(int((*base.List)(node.CommonNode).VLen()), func(cm *CommonNode) bool {
		return fn(&NodeName{CommonNode: cm})
	})
	if i < int((*base.List)(node.CommonNode).VLen()) {
		result, _ = node.At(i)
	}

	return result
}

func (node ListType) SearchIndex(fn func(*NodeName) bool) int {

	i := (*base.List)(node.CommonNode).SearchIndex(int((*base.List)(node.CommonNode).VLen()), func(cm *CommonNode) bool {
		return fn(&NodeName{CommonNode: cm})
	})
	if i < int((*base.List)(node.CommonNode).VLen()) {
		return i
	}

	return -1
}
