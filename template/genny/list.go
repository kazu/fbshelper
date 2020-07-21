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

	list.InitList()
	return list
}

func emptyListType() *ListType {
	return &ListType{CommonNode: &base.CommonNode{}}
}

func (node ListType) At(i int) (result *NodeName, e error) {
	result = &NodeName{}
	result.CommonNode, e = node.CommonNode.At(i)
	return
}

func (node ListType) SetAt(i int, v *NodeName) error {
	return node.CommonNode.SetAt(i, v.CommonNode)
}

func (node ListType) First() (result *NodeName, e error) {
	return node.At(0)
}

func (node ListType) Last() (result *NodeName, e error) {
	return node.At(int(node.NodeList.ValueInfo.VLen) - 1)
}

func (node ListType) Select(fn func(*NodeName) bool) (result []*NodeName) {
	result = make([]*NodeName, 0, int(node.NodeList.ValueInfo.VLen))
	commons := node.CommonNode.Select(func(cm *CommonNode) bool {
		return fn(&NodeName{CommonNode: cm})
	})
	for _, cm := range commons {
		result = append(result, &NodeName{CommonNode: cm})
	}
	return result
}

func (node ListType) Find(fn func(*NodeName) bool) *NodeName {
	result := &NodeName{}
	result.CommonNode = node.CommonNode.Find(func(cm *CommonNode) bool {
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
