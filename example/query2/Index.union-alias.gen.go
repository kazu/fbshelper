// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package query

import "github.com/kazu/fbshelper/query/base"

/*
genny must be called per Index ;
*/

var DUMMP_IndexIndexNum bool = base.SetAlias("Index", "IndexNum")

func (node Index) IndexNum() IndexNum {
	//result := IndexNum{CommonNode: node.CommonNode}
	result := IndexNum{}
	result.CommonNode = &CommonNode{}
	result.NodeList = node.NodeList
	result.CommonNode.Name = "IndexNum"
	result.FetchIndex()
	return result
}

func IndexFromIndexNum(v *IndexNum) *Index {
	result := &Index{}
	result.CommonNode = v.CommonNode
	result.FetchIndex()
	return result
}

/*
genny must be called per Index ;
*/

var DUMMP_IndexIndexString bool = base.SetAlias("Index", "IndexString")

func (node Index) IndexString() IndexString {
	//result := IndexString{CommonNode: node.CommonNode}
	result := IndexString{}
	result.CommonNode = &CommonNode{}
	result.NodeList = node.NodeList
	result.CommonNode.Name = "IndexString"
	result.FetchIndex()
	return result
}

func IndexFromIndexString(v *IndexString) *Index {
	result := &Index{}
	result.CommonNode = v.CommonNode
	result.FetchIndex()
	return result
}

/*
genny must be called per Index ;
*/

var DUMMP_IndexFile bool = base.SetAlias("Index", "File")

func (node Index) File() File {
	//result := File{CommonNode: node.CommonNode}
	result := File{}
	result.CommonNode = &CommonNode{}
	result.NodeList = node.NodeList
	result.CommonNode.Name = "File"
	result.FetchIndex()
	return result
}

func IndexFromFile(v *File) *Index {
	result := &Index{}
	result.CommonNode = v.CommonNode
	result.FetchIndex()
	return result
}

/*
genny must be called per Index ;
*/

var DUMMP_IndexInvertedMapNum bool = base.SetAlias("Index", "InvertedMapNum")

func (node Index) InvertedMapNum() InvertedMapNum {
	//result := InvertedMapNum{CommonNode: node.CommonNode}
	result := InvertedMapNum{}
	result.CommonNode = &CommonNode{}
	result.NodeList = node.NodeList
	result.CommonNode.Name = "InvertedMapNum"
	result.FetchIndex()
	return result
}

func IndexFromInvertedMapNum(v *InvertedMapNum) *Index {
	result := &Index{}
	result.CommonNode = v.CommonNode
	result.FetchIndex()
	return result
}

/*
genny must be called per Index ;
*/

var DUMMP_IndexInvertedMapString bool = base.SetAlias("Index", "InvertedMapString")

func (node Index) InvertedMapString() InvertedMapString {
	//result := InvertedMapString{CommonNode: node.CommonNode}
	result := InvertedMapString{}
	result.CommonNode = &CommonNode{}
	result.NodeList = node.NodeList
	result.CommonNode.Name = "InvertedMapString"
	result.FetchIndex()
	return result
}

func IndexFromInvertedMapString(v *InvertedMapString) *Index {
	result := &Index{}
	result.CommonNode = v.CommonNode
	result.FetchIndex()
	return result
}

/*
genny must be called per Index ;
*/

var DUMMP_IndexNumList bool = base.SetAlias("Index", "NumList")

func (node Index) NumList() NumList {
	//result := NumList{CommonNode: node.CommonNode}
	result := NumList{}
	result.CommonNode = &CommonNode{}
	result.NodeList = node.NodeList
	result.CommonNode.Name = "NumList"
	result.FetchIndex()
	return result
}

func IndexFromNumList(v *NumList) *Index {
	result := &Index{}
	result.CommonNode = v.CommonNode
	result.FetchIndex()
	return result
}

/*
genny must be called per Index ;
*/

var DUMMP_IndexHoges bool = base.SetAlias("Index", "Hoges")

func (node Index) Hoges() Hoges {
	//result := Hoges{CommonNode: node.CommonNode}
	result := Hoges{}
	result.CommonNode = &CommonNode{}
	result.NodeList = node.NodeList
	result.CommonNode.Name = "Hoges"
	result.FetchIndex()
	return result
}

func IndexFromHoges(v *Hoges) *Index {
	result := &Index{}
	result.CommonNode = v.CommonNode
	result.FetchIndex()
	return result
}
