package query

import (
	"github.com/kazu/fbshelper/query/base"
	"github.com/kazu/fbshelper/query/log"
)

/*
   must call 1 times per Table/struct(NodeName)
*/

type NodeName struct {
	*base.CommonNode
}

func emptyNodeName() *NodeName {
	return &NodeName{CommonNode: &base.CommonNode{}}
}

var NodeName_IdxToType map[int]int = map[int]int{}
var NodeName_IdxToTypeGroup map[int]int = map[int]int{}
var NodeName_IdxToName map[int]string = map[int]string{}
var NodeName_NameToIdx map[string]int = map[string]int{}

var DUMMP_NodeNameIsStruct bool = base.SetNameIsStrunct("NodeName", base.ToBool("IsStruct"))

func SetNodeNameFields(nName, fName, fType string, fNum int) bool {

	base.RequestSettingNameFields(nName, fName, fType, fNum)

	enumFtype, ok := base.NameToType[fType]
	if ok {
		NodeNameSetIdxToType(fNum, enumFtype)
	}
	//FIXME: basic type only store?

	NodeNameSetIdxToName(fNum, fType)

	grp := NodeNameGetTypeGroup(fType)
	NodeNameSetTypeGroup(fNum, grp)

	NodeName_IdxToName[fNum] = fType

	NodeName_NameToIdx[fName] = fNum
	base.SetNameToIdx("NodeName", NodeName_NameToIdx)

	return true

}
func NodeNameSetIdxToName(i int, s string) {
	NodeName_IdxToName[i] = s

	base.SetIdxToName("NodeName", NodeName_IdxToName)
}

func NodeNameSetIdxToType(k, v int) bool {
	NodeName_IdxToType[k] = v
	base.SetIdxToType("NodeName", NodeName_IdxToType)
	return true
}

func NodeNameSetTypeGroup(k, v int) bool {
	NodeName_IdxToTypeGroup[k] = v
	base.SetdxToTypeGroup("NodeName", NodeName_IdxToTypeGroup)
	return true
}

func NodeNameGetTypeGroup(s string) (result int) {
	return base.GetTypeGroup(s)
}

func (node NodeName) commonNode() *base.CommonNode {
	if node.CommonNode == nil {
		log.Log(log.LOG_WARN, func() log.LogArgs {
			return log.F("CommonNode not found NodeName")
		})
	} else if len(node.CommonNode.Name) == 0 || len(node.CommonNode.IdxToType) == 0 {
		node.CommonNode.Name = "NodeName"
		node.CommonNode.IdxToType = NodeName_IdxToType
		node.CommonNode.IdxToTypeGroup = NodeName_IdxToTypeGroup
	}
	return node.CommonNode
}
func (node NodeName) SearchInfo(pos int, fn base.RecFn, condFn base.CondFn) {

	node.commonNode().SearchInfo(pos, fn, condFn)

}

func (node NodeName) Info() (info base.Info) {

	return node.commonNode().Info()

}

func (node NodeName) IsLeafAt(j int) bool {

	return node.commonNode().IsLeafAt(j)

}

func (node NodeName) CountOfField() int {
	return len(NodeName_IdxToType)
}

func (node NodeName) ValueInfo(i int) base.ValueInfo {
	return node.commonNode().ValueInfo(i)
}

func (node NodeName) FieldAt(idx int) *base.CommonNode {
	return node.commonNode().FieldAt(idx)
}

func (src NodeName) Equal(dst NodeName) bool {
	for i := 0; i < src.CountOfField(); i++ {
		if !src.FieldAt(i).Equal(dst.FieldAt(i)) {
			return false
		}
	}
	return true
}

type NodeNameWithErr struct {
	*NodeName
	Err error
}

func NodeNameSingle(node *NodeName, e error) NodeNameWithErr {
	return NodeNameWithErr{NodeName: node, Err: e}
}

func NewNodeName() *NodeName {
	base.ApplyRequestNameFields()
	node := emptyNodeName()
	node.NodeList = &base.NodeList{}
	node.CommonNode.Name = "NodeName"
	node.Init()

	return node
}

func (node NodeName) FieldGroups() map[int]int {
	return NodeName_IdxToTypeGroup
}

func (node NodeName) Root() (RootType, error) {
	if !node.InRoot() {
		return RootType{}, log.ERR_NO_INCLUDE_ROOT
	}
	root := toRoot(node.IO)
	return root, nil
}
