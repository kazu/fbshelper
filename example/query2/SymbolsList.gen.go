// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

type SymbolsList struct { // genny
	*CommonNode
}

// Symbols genny
func NewSymbolsList() *SymbolsList {

	list := emptySymbolsList()
	list.NodeList = &base.NodeList{}
	list.CommonNode.Name = "[]Symbols"

	(*base.List)(list.CommonNode).InitList()
	return list
}

func emptySymbolsList() *SymbolsList {
	return &SymbolsList{CommonNode: &base.CommonNode{}}
}

func (node SymbolsList) At(i int) (result *Symbols, e error) {
	result = &Symbols{}
	result.CommonNode, e = (*base.List)(node.CommonNode).At(i)
	return
}

func (node SymbolsList) AtWihoutError(i int) (result *Symbols) {
	result, e := node.At(i)
	if e != nil {
		result = nil
	}
	return
}

func (node SymbolsList) SetAt(i int, v *Symbols) error {
	return (*base.List)(node.CommonNode).SetAt(i, v.CommonNode)
}
func (node SymbolsList) Add(v SymbolsList) error {
	return (*base.List)(node.CommonNode).Add((*base.List)(v.CommonNode))
}

func (node SymbolsList) Range(start, last int) *SymbolsList {
	l := (*base.List)(node.CommonNode).New(base.OptRange(start, last))
	if l == nil {
		return nil
	}
	return &SymbolsList{CommonNode: (*base.CommonNode)(l)}
}

func (node SymbolsList) First() (result *Symbols, e error) {
	return node.At(0)
}

func (node SymbolsList) Last() (result *Symbols, e error) {
	return node.At(int(node.NodeList.ValueInfo.VLen) - 1)
}

func (node SymbolsList) Select(fn func(*Symbols) bool) (result []*Symbols) {
	result = make([]*Symbols, 0, int(node.NodeList.ValueInfo.VLen))
	commons := (*base.List)(node.CommonNode).Select(func(cm *CommonNode) bool {
		return fn(&Symbols{CommonNode: cm})
	})
	for _, cm := range commons {
		result = append(result, &Symbols{CommonNode: cm})
	}
	return result
}

func (node SymbolsList) Find(fn func(*Symbols) bool) *Symbols {
	result := &Symbols{}
	result.CommonNode = (*base.List)(node.CommonNode).Find(func(cm *CommonNode) bool {
		return fn(&Symbols{CommonNode: cm})
	})
	return result
}

func (node SymbolsList) All() []*Symbols {
	return node.Select(func(*Symbols) bool { return true })
}

func (node SymbolsList) Count() int {
	return int(node.NodeList.ValueInfo.VLen)
}

func (node SymbolsList) SwapAt(i, j int) error {
	return (*List)(node.CommonNode).SwapAt(i, j)
}

func (node SymbolsList) SortBy(less func(i, j int) bool) error {
	return (*List)(node.CommonNode).SortBy(less)
}

// Search ... binary search
func (node SymbolsList) Search(fn func(*Symbols) bool) *Symbols {
	result := &Symbols{}

	i := (*base.List)(node.CommonNode).SearchIndex(int((*base.List)(node.CommonNode).VLen()), func(cm *CommonNode) bool {
		return fn(&Symbols{CommonNode: cm})
	})
	if i < int((*base.List)(node.CommonNode).VLen()) {
		result, _ = node.At(i)
	}

	return result
}

func (node SymbolsList) SearchIndex(fn func(*Symbols) bool) int {

	i := (*base.List)(node.CommonNode).SearchIndex(int((*base.List)(node.CommonNode).VLen()), func(cm *CommonNode) bool {
		return fn(&Symbols{CommonNode: cm})
	})
	if i < int((*base.List)(node.CommonNode).VLen()) {
		return i
	}

	return -1
}
