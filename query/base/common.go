package base

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	flatbuffers "github.com/google/flatbuffers/go"
	log "github.com/kazu/fbshelper/query/log"
	"github.com/kazu/loncha"
)

// CommonNode ... is node of flatbuffers.
type CommonNode struct {
	*NodeList
	Name           string
	IdxToType      map[int]int
	IdxToTypeGroup map[int]int

	// for sorting of List
	lessFn func(i, j int) bool
}

// Info ... return information of node
func (node *CommonNode) Info() (info Info) {

	nName := node.Name

	if node.NodeList == nil || node.Node == nil {
		//FIXME
		return info
	}
	info.Pos = node.Node.Pos
	info.Size = -1

	grp := GetTypeGroup(nName)
	if IsFieldBasicType(grp) {
		info.Size = TypeToSize[NameToTypeEnum(nName)]
		return info
	}

	if IsStructName[nName] {
		size := 0
		for i := 0; i < len(node.IdxToTypeGroup); i++ {
			size += TypeToSize[node.IdxToType[i]]
		}
		info.Size = size
		return info
	}
	for i := 0; i < len(node.IdxToTypeGroup); i++ {
		vInfo := node.ValueInfo(i)
		if info.Pos+info.Size < vInfo.Pos+vInfo.Size {
			info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
		}
	}
	return info
}

// ValueInfo ... return infomation of field/attribute of node.
func (node *CommonNode) ValueInfo(idx int) ValueInfo {

	info := ValueInfo{Pos: node.Node.Pos, Size: 0}
	if IsStructName[node.Name] {
		for i := 0; i < idx; i++ {
			info.Pos += info.Size
			info.Size = TypeToSize[node.IdxToType[i]]
		}
		return info
	}

	grp := node.IdxToTypeGroup[idx]

	if IsFieldStruct(grp) {
		info.Pos = node.VirtualTable(idx)
		info.Size = node.FieldAt(idx).Info().Size

	} else if IsFieldUnion(grp) {
		info.Pos = node.Table(idx)

		if node.CanFollow(idx) {
			info.Size = node.FollowUnion(idx).Info().Size
		}

	} else if IsFieldBytes(grp) {
		info = node.ValueInfoPosList(idx)
		info.Size = (*List)(node.FieldAt(idx)).InfoSlice().Size

	} else if IsFieldSlice(grp) {
		info = node.ValueInfoPosList(idx)
		if node.FieldAt(idx).NodeList.Node == nil {
			// FIXME remove later
			// node.ValueInfoPosList(idx)
			// node.FieldAt(idx)
			// Log(LOG_WARN, func() LogArgs {
			// 	return F(
			// 		"node not found: node=%s.FieldAt(%d)", node.Name, idx)
			// })
			return info
		}

		info.Size = (*List)(node.FieldAt(idx)).InfoSlice().Size

	} else if IsFieldTable(grp) {
		info.Pos = node.Table(idx)
		info.Size = node.FieldAt(idx).Info().Size

	} else if IsFieldBasicType(grp) {
		info.Pos = node.VirtualTable(idx)
		info.Size = TypeToSize[node.IdxToType[idx]]

	} else {
		Log(LOG_ERROR, func() LogArgs {
			return F("Invalid %s.%s idx=%d\n", node.Name, All_IdxToName[node.Name][idx], idx)
		})
	}

	return info
}

// CanFollow ... check instance on Union.
func (node *CommonNode) CanFollow(idx int) bool {

	return node.FieldAt(idx-1).Byte() > 0
}

// FollowUnion ... return instance via Union.
func (node *CommonNode) FollowUnion(idx int) *CommonNode {
	idxOfAlias := node.FieldAt(idx-1).Byte() - 1
	if int(idxOfAlias) >= len(UnionAlias[All_IdxToName[node.Name][idx]]) {
		Log(LOG_WARN, func() LogArgs {
			return F(
				"Invalid union=%s aliases=%+v idx=%d\n",
				node.Name, All_IdxToName[node.Name], idxOfAlias)
		})
		return nil
	}

	newName := UnionAlias[All_IdxToName[node.Name][idx]][idxOfAlias]
	union := node.FieldAt(idx)

	next := &CommonNode{
		NodeList:       union.NodeList,
		Name:           newName,
		IdxToType:      union.IdxToType,
		IdxToTypeGroup: union.IdxToTypeGroup,
	}
	next.FetchIndex()
	return next
}

// SearchInfo ... collect Infomation via traversing.
// if TraverseInfo is work. removed
func (node *CommonNode) SearchInfo(pos int, fn RecFn, condFn CondFn) {
	info := node.Info()

	if condFn(pos, info) {
		fn(NodePath{Name: node.Name, Idx: -1}, info)
	} else {
		return
	}

	for i := 0; i < len(node.IdxToTypeGroup); i++ {
		g := node.IdxToTypeGroup[i]

		if node.IsLeafAt(i) {
			fInfo := Info(node.ValueInfo(i))
			if condFn(pos, fInfo) {
				fn(NodePath{Name: node.Name, Idx: i}, fInfo)
			}
			continue
		}
		if IsMatchBit(g, FieldTypeStruct) {
			node.FieldAt(i).SearchInfo(pos, fn, condFn)
		} else if IsMatchBit(g, FieldTypeUnion) && node.CanFollow(i) {
			node.FollowUnion(i).SearchInfo(pos, fn, condFn)
		} else if IsMatchBit(g, FieldTypeSlice) && IsMatchBit(g, FieldTypeBasic1) {
			(*List)(node.FieldAt(i)).SearchInfoSlice(pos, fn, condFn)
		} else if IsMatchBit(g, FieldTypeSlice) {
			(*List)(node.FieldAt(i)).SearchInfoSlice(pos, fn, condFn)
		} else if IsMatchBit(g, FieldTypeTable) {
			node.FieldAt(i).SearchInfo(pos, fn, condFn)
		} else if IsMatchBit(g, FieldTypeBasic) {
		} else {
			Log(LOG_ERROR, func() LogArgs {
				return F("node must be Noder")
			})
		}
	}

}

// IsLeafAt ... field has child or not.
func (node *CommonNode) IsLeafAt(j int) bool {
	nName := node.Name
	tGroup := node.IdxToTypeGroup[j]

	if IsFieldStruct(tGroup) {

		return false
	} else if IsFieldUnion(tGroup) {

		return false
	} else if IsFieldBytes(tGroup) {

		return true
	} else if IsFieldSlice(tGroup) {

		return false
	} else if IsFieldTable(tGroup) {

		return false
	} else if IsFieldBasicType(tGroup) {
		return true
	} else {
		Log(LOG_ERROR, func() LogArgs {
			return F("Invalid Node=%s tGroup=%d\n", nName, tGroup)
		})
	}
	return false

}

// CountOfField ... return count of field/attribute
func (node *CommonNode) CountOfField() int {
	return len(node.IdxToTypeGroup)
}

// FieldAt ... return field by field index
func (node *CommonNode) FieldAt(idx int) (cNode *CommonNode) {

	node.clearValueInfoOnDirty()

	result := &NodeList{}
	grp := node.IdxToTypeGroup[idx]

	if IsStructName[node.Name] {
		pos := node.Node.Pos
		for i := 0; i < node.CountOfField(); i++ {
			if i < idx {
				pos += TypeToSize[node.IdxToType[i]]
			} else {
				break
			}
		}
		result.Node = NewNode2(node.Node.IO, pos, true)

		goto RESULT
	}

	if node.VirtualTableIsZero(idx) {
		goto RESULT
	}

	if IsFieldStruct(grp) {
		result.Node = node.ValueStruct(idx)

	} else if IsFieldUnion(grp) {
		result.Node = node.ValueTable(idx)

	} else if IsFieldBytes(grp) {

		valInfo := node.ValueInfoPosBytes(idx)
		valInfo.VLen = uint32(valInfo.Size)

		nNode := NewNode2(node.IO, valInfo.Pos, true)
		nNode.Size = valInfo.Size
		result.Node = nNode
		result.ValueInfo = valInfo

	} else if IsFieldSlice(grp) {
		nodeList := node.ValueList(idx)
		result = &nodeList

	} else if IsFieldTable(grp) {
		result.Node = node.ValueTable(idx)

	} else if IsFieldBasicType(grp) {
		result.Node = NewNode2(node.IO, node.VirtualTable(idx), true)
		result.Node.Size = TypeToSize[node.IdxToType[idx]]

	} else {
		Log(LOG_ERROR, func() LogArgs {
			return F("FieldAt: Invalid Node=%s idx=%d\n", node.Name, idx)
		})
	}
RESULT:
	cNode = &CommonNode{}
	cNode.NodeList = result
	name := All_IdxToName[node.Name][idx]
	cNode.Name = name
	cNode.FetchIndex()
	// a := node.Base.Impl()
	// b := cNode.Base
	// _, _ = a, b

	if !IsStructName[node.Name] && !node.VirtualTableIsZero(idx) && !node.IO.Impl().Equal(cNode.IO.Impl()) {
		Log(LOG_DEBUG, func() LogArgs {
			return F("FieldAt: Invalid Node=%s idx=%d\n", node.Name, idx)
		})
	}

	return cNode
}

// FetchIndex ... fetching index
func (node *CommonNode) FetchIndex() {
	if _, ok := All_IdxToName[node.Name]; ok {
		node.IdxToType = All_IdxToType[node.Name]
		node.IdxToTypeGroup = All_IdxToTypeGroup[node.Name]
	}
}

// IsList ... return true if is List Node.
func (node *CommonNode) IsList() bool {
	if node.Name == "" && node.NodeList.ValueInfo.VLen > 0 {
		return true
	}
	if node.Name[0:2] == "[]" {
		return true
	}
	// if node.NodeList.ValueInfo.VLen > 0 {
	// 	return true
	// }
	return false
	//FIXME ?
	//return node.VLen() > 0
}

// List ... return List instance if node is List.
func (node *CommonNode) List() *List {

	if !node.IsList() {
		return nil
	}

	return (*List)(node)
}

// Open ... Open New CommonNode with io.Reader
func Open(r io.Reader, cap int, opts ...Option) (node *CommonNode) {

	if len(opts) > 0 {
		opts[0](&DefaultOption)
	}

	ApplyRequestNameFields()

	b := NewBaseByIO(r, 512)

	node = &CommonNode{}
	node.NodeList = &NodeList{}
	node.Node = NewNode(b, int(flatbuffers.GetUOffsetT(b.R(0, Size(4)))))
	return node
}

// OpenByBuf ... Open New CommonNode with byte buffer
func OpenByBuf(buf []byte, opts ...Option) *CommonNode {

	if len(opts) > 0 {
		opts[0](&DefaultOption)
	}

	ApplyRequestNameFields()

	node := &CommonNode{}
	node.NodeList = &NodeList{}
	node.Node = NewNode(NewBase(buf), int(flatbuffers.GetUOffsetT(buf)))
	return node
}

// Len ... return len of byte size.
func (node *CommonNode) Len() int {
	info := node.Info()
	size := info.Pos + info.Size

	if (size % 8) == 0 {
		return size
	}

	return size + (8 - (size % 8))
}

// Unmarshal ... unmarshal to struct. only work fbs's struct node.
func (node *CommonNode) Unmarshal(v interface{}) error {

	return node.unmarshal(v, func(s string, rv reflect.Value) error {
		z := node.CountOfField()
		_ = z
		i, ok := All_NameToIdx[node.Name][s]
		if !ok {
			return nil
		}

		grp := node.IdxToTypeGroup[i]
		if IsFieldStruct(grp) {

		} else if IsFieldUnion(grp) {

		} else if IsFieldBytes(grp) {
			cNode := node.FieldAt(i)
			rv.SetBytes(cNode.R(cNode.Node.Pos, Size(cNode.Node.Size))[:cNode.Node.Size])
		} else if IsFieldBasicType(grp) {
			cNode := node.FieldAt(i)
			rv.Set(reflect.ValueOf(cNode).MethodByName(cNode.Name).Call([]reflect.Value{})[0])
		}

		return nil
	})
}

func (node *CommonNode) value() interface{} {

	if _, ok := NameToType[node.Name]; !ok {
		return reflect.ValueOf(node)
	}

	return reflect.ValueOf(node).MethodByName(node.Name).Call([]reflect.Value{})[0]

}

func (node *CommonNode) SetBool(v bool) error {

	if node.Name != "Bool" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteBool(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetByte(v byte) error {

	if node.Name != "Byte" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteByte(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetInt8(v int8) error {

	if node.Name != "Int8" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteInt8(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetInt16(v int16) error {

	if node.Name != "Int16" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteInt16(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetInt32(v int32) error {

	if node.Name != "Int32" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteInt32(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetInt64(v int64) error {

	if node.Name != "Int64" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteInt64(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetUint8(v uint8) error {

	if node.Name != "Uint8" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteUint8(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetUint16(v uint16) error {

	if node.Name != "Uint16" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteUint16(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetUint32(v uint32) error {

	if node.Name != "Uint32" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteUint32(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetUint64(v uint64) error {

	if node.Name != "Uint64" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteUint64(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetFloat32(v float32) error {

	if node.Name != "Float32" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteFloat32(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

func (node *CommonNode) SetFloat64(v float64) error {

	if node.Name != "Float64" {
		return ERR_INVALID_TYPE
	}
	if node.Node.Size <= 0 {
		return ERR_INVALID_TYPE
	}

	flatbuffers.WriteFloat64(node.U(node.Node.Pos, node.Node.Size), v)

	return nil

}

// TraverseInfo ... traverse Info/ValueInfo of all node.
func (node *CommonNode) TraverseInfo(pos int, fn TraverseRec, condFn TraverseCond) {

	for i := 0; i < len(node.IdxToTypeGroup); i++ {

		if node.IsLeafAt(i) {
			fInfo := Info(node.ValueInfo(i))
			if condFn(node.Node.Pos, fInfo.Pos, fInfo.Size) {
				fn(node, i, node.Node.Pos, fInfo.Pos, fInfo.Size)
			}
			continue
		}

		// fixme: not require  to get data via FieldAt()
		var next *CommonNode
		g := node.IdxToTypeGroup[i]
		if !IsFieldUnion(g) {
			next = node.FieldAt(i)
		} else {
			next = node.FollowUnion(i)
		}
		var nPos int
		if IsFieldSlice(g) {
			nPos = next.NodeList.ValueInfo.Pos
		} else {
			if next == nil || next.Node == nil {
				// empty node
				continue
			}
			nPos = next.Node.Pos
		}

		if condFn(node.Node.Pos, nPos, -1) {
			fn(node, i, node.Node.Pos, nPos, -1)
		}

		if IsFieldSlice(g) {
			next.TraverseInfoSlice(pos, fn, condFn)
		} else {
			if IsFieldUnion(g) || IsFieldTable(g) {
				if node.VirtualTable(i) == node.Table(i) {
					continue
				}
			}
			next.TraverseInfo(pos, fn, condFn)
		}
	}
}

// TraverseInfoSlice ... traverse Info/ValueInfo of all list node.
func (node *CommonNode) TraverseInfoSlice(pos int, fn TraverseRec, condFn TraverseCond) {

	var v interface{}
	for i, cNode := range (*List)(node).All() {

		if condFn(node.NodeList.ValueInfo.Pos, cNode.Node.Pos, -1) {
			fn(node, i, node.NodeList.ValueInfo.Pos, cNode.Node.Pos, -1)
		}

		v = cNode
		if vv, ok := v.(Searcher); ok {
			vv.TraverseInfo(pos, fn, condFn)
		} else {
			goto NO_NODE
		}
	}
	return

NO_NODE:

	for i := 0; i < int(node.NodeList.ValueInfo.VLen); i++ {
		ptr := int(node.NodeList.ValueInfo.Pos) + i*4
		start := ptr + int(flatbuffers.GetUint32(node.R(ptr, Size(4))))
		// size := info.Size
		// if i+1 < int(node.NodeList.ValueInfo.Pos) {
		// 	size = ptr + 4 + int(flatbuffers.GetUint32(node.R(ptr+4))) - start
		// }
		if condFn(node.Node.Pos, start, -1) {
			fn(node, i, node.NodeList.ValueInfo.Pos, start, -1)
		}
	}
}

// Count ... count of List node.
func (node *CommonNode) Count() int {
	return int(node.NodeList.ValueInfo.VLen)
}

// Tree ... return node tree.
type Tree struct {
	Node   *CommonNode
	Parent *Tree
	Childs []*Tree
}

// Pos ... return offset position in buffer.
func (t Tree) Pos() int {

	if t.Node.NodeList.ValueInfo.Pos > 0 {
		return t.Node.NodeList.ValueInfo.Pos
	}
	if t.Node.Node == nil {
		return -1
	}
	return t.Node.Node.Pos
}

// Size ... return Size  in buffer.
func (t Tree) Size() int {

	if t.Node.NodeList.ValueInfo.Pos > 0 {
		return t.Node.NodeList.ValueInfo.Size
	}
	if t.Node.Node == nil {
		return -1
	}
	return t.Node.Node.Size
}

// Dump ... dump of node information
func (tree Tree) Dump() string {
	if tree.Parent != nil {
		return fmt.Sprintf("{Type:\t%s,\tPos:\t%d,\tSize:\t%d,\tParentPos:\t%d}\n",
			tree.Node.Name,
			tree.Pos(),
			tree.Size(),
			tree.Parent.Pos(),
		)
	}
	return fmt.Sprintf("{Type:\t%s,\tPos:\t%dtSize:\t%d}\n",
		tree.Node.Name,
		tree.Pos(),
		tree.Size(),
	)
}

// DumpAll ... dump of all node information
func (tree Tree) DumpAll(i int, w io.Writer) {

	for j := 0; j < i*2; j++ {
		io.WriteString(w, "\t")
	}
	fmt.Fprintf(w, "%s", tree.Dump())
	for _, child := range tree.Childs {
		child.DumpAll(i+1, w)
	}
}

// AllTree ... return all Tree
func (node *CommonNode) AllTree() *Tree {

	type Result struct {
		name   string
		pNode  *CommonNode
		idx    int
		parent int
		child  int
		size   int
	}

	results := []Result{}

	recFn := func(node *CommonNode, idx, parent, child, size int) {
		results = append(results,
			Result{
				name:   node.Name,
				pNode:  node,
				idx:    idx,
				parent: parent,
				child:  child,
				size:   size,
			})
	}

	cond := func(parent, child, size int) bool { return true }
	node.TraverseInfo(0, recFn, cond)

	trees := make([]*Tree, 0, len(results)+1)
	trees = append(trees, &Tree{Node: node})

	for _, result := range results {
		tree := &Tree{}
		if result.pNode.NodeList.ValueInfo.Pos > 0 {
			tree.Node, _ = (*List)(result.pNode).At(result.idx)
		} else {
			tree.Node = result.pNode.FieldAt(result.idx)
		}
		trees = append(trees, tree)
	}

	for i := range results {
		j, e := loncha.IndexOf(trees, func(idx int) bool {
			if trees[idx].Node.NodeList.ValueInfo.Pos > 0 {
				return trees[idx].Node.NodeList.ValueInfo.Pos == results[i].parent
			}
			return trees[idx].Node.Node != nil && (trees[idx].Node.Node.Pos == results[i].parent)
		})
		if e != nil {
			fmt.Fprintf(os.Stderr, "NOT FOUND parent=%d\n", results[i].parent)
			continue
		}
		trees[j].Childs = append(trees[j].Childs, trees[i+1])
		trees[i+1].Parent = trees[j]
	}

	return trees[0]

}

type TreeCond func(*Tree) bool

func findtree(tree *Tree, cond TreeCond, out chan *Tree) {

	for _, child := range tree.Childs {
		if cond(child) {
			out <- child
		}
		findtree(child, cond, out)
	}

}

// FindTree ... return Tree by condtion function.
func (node *CommonNode) FindTree(cond TreeCond) <-chan *Tree {

	ch := make(chan *Tree, 10)
	tree := node.AllTree()

	go func() {
		findtree(tree, cond, ch)
		close(ch)
	}()

	return ch

}

func (node *CommonNode) root() *CommonNode {
	common := &CommonNode{}
	common.NodeList = &NodeList{}
	common.Node = NewNode(node.IO, int(flatbuffers.GetUOffsetT(node.R(0, Size(4)))))
	common.Name = RootName
	common.FetchIndex()
	return common
}

// RootCommon ... return Root Node.
func (node *CommonNode) RootCommon() *CommonNode {
	return node.root()
}

// InRoot ... return true if buffer include Root Node.
func (node *CommonNode) InRoot() bool {

	pos := int(flatbuffers.GetVOffsetT(node.R(0, Size(2))))
	if len(node.R(pos, Size(4))) < 4 {
		return false
	}
	if node.Node.Pos < 8 {
		return false
	}

	return pos != int(flatbuffers.GetUOffsetT(node.R(pos, Size(4))))

}

// InsertBuf ... insert buffer in pos position.
func (node *CommonNode) InsertBuf(pos, size int) {

	node.InsertSpace(pos, size, true)
}

// InsertSpace ... insert buffer in pos position.
func (node *CommonNode) InsertSpace(pos, size int, isInsert bool) {

	newBase := node.IO.insertSpace(pos, size, isInsert)

	defer func() {
		if node.IsList() {
			if node.NodeList.ValueInfo.Pos > pos {
				node.NodeList.ValueInfo.Pos += size
			}
		} else if node.Node.Pos > pos {
			node.Node.Pos += size
		}
		node.IO = newBase
		//node.Base.Impl().overwrite(newBase.Impl())

	}()

	//FIXME: if dont has root, not update vtable
	if !node.InRoot() {
		return
	}

	root := node.root()
	ch := root.FindTree(func(t *Tree) bool {
		return t.Parent != nil &&
			t.Node.Node != nil &&
			t.Parent.Node.Node != nil &&
			t.Parent.Pos() <= pos && pos <= t.Pos()
	})

	for {
		cTree, ok := <-ch
		if !ok {
			break
		}
		tree := cTree.Parent
		if tree.Node.Name == node.Name && tree.Node.Node.Pos == node.Node.Pos {
			continue
		}
		tree.Node.IO = newBase
		//tree.Node.Base.Impl().overwrite(newBase.Impl())
		idx, err := loncha.IndexOf(tree.Childs, func(i int) bool {
			return tree.Childs[i] == cTree
		})
		if err != nil {
			Log(LOG_ERROR, func() LogArgs {
				return F("invalid Tree e=%s\n", err.Error())
			})
			continue
		}
		tree.Node.movePos(idx, pos, size)
		// fIXME
		// oldBase.Diffs = newBase.Diffs
		// oldBase.bytes = newBase.bytes
		// oldBase.dirties = newBase.dirties
		// tree.Node.Base = oldBase
		// node.Base = oldBase
	}
	return
}

func (node *CommonNode) movePosOnList(i, pos, size int) {

	if !node.IsList() {
		Log(LOG_ERROR, func() LogArgs {
			return F("movePosOnList: Invalid Node=%s idx=%d\n", node.Name, i)
		})
		return
	}

	if i >= int(node.NodeList.ValueInfo.VLen) || i < 0 {
		Log(LOG_ERROR, func() LogArgs {
			return F("movePosOnList: Invalid Index Node=%s idx=%d\n", node.Name, i)
		})
		return
	}

	tName := node.Name[2:]
	grp := GetTypeGroup(tName)
	ptr := int(node.NodeList.ValueInfo.Pos) + i*4

	if IsFieldBasicType(grp) {
		Log(LOG_ERROR, func() LogArgs {
			return F("movePosOnList: BasicType must not be moved Node=%s idx=%d\n", node.Name, i)
		})
		return
	} else {
		nextOff := flatbuffers.GetUint32(node.R(ptr, Size(4)))
		nextOff += uint32(size)
		flatbuffers.WriteUint32(node.U(ptr, SizeOfuint32), nextOff)
	}

	node.IO.AddDirty(Dirty{Pos: node.NodeList.ValueInfo.Pos})

}

func (node *CommonNode) movePos(idx, pos, size int) {
	if node.IsList() {
		node.movePosOnList(idx, pos, size)
		return
	}

	if IsStructName[node.Name] {
		// FIXME: checking struct
		return
	}
	grp := node.IdxToTypeGroup[idx]

	if IsFieldUnion(grp) || IsFieldTable(grp) || IsFieldSlice(grp) || IsFieldBytes(grp) {
		// FIXME: VTable[idx] is 0 pattern
		cPos := node.VirtualTable(idx)

		nextOff := flatbuffers.GetUint32(node.R(cPos, Size(4)))
		if cPos+int(nextOff) < pos {
			Log(LOG_WARN, func() LogArgs {
				return F("%s.movePos(%d, %d, %d) skip tableption latger %d\n",
					node.Name, idx, pos, size, cPos+int(nextOff))
			})
			return
		}
		Log(LOG_DEBUG, func() LogArgs {
			return F("%s.movePos(%d, %d, %d)  oldpos=%d\n",
				node.Name, idx, pos, size, cPos+int(nextOff))
		})

		nextOff += uint32(size)
		flatbuffers.WriteUint32(node.U(cPos, SizeOfuint32), nextOff)
	} else if IsFieldStruct(grp) || IsFieldBasicType(grp) {
		Log(LOG_WARN, func() LogArgs {
			return F("MovePos: fixme must move Node=%s idx=%d\n", node.Name, idx)
		})
	} else {
		Log(LOG_ERROR, func() LogArgs {
			return F("MovePos: Invalid Node=%s idx=%d\n", node.Name, idx)
		})
	}

	node.AddDirty(Dirty{Pos: node.Node.Pos})

}

// IsRoot ... return true if root node.
func (node *CommonNode) IsRoot() bool {

	if RootName == node.Name {
		return true
	}
	return false
}

// Init ... initialize as CommonNode.
func (node *CommonNode) Init() error {

	node.IdxToType = All_IdxToType[node.Name]
	node.IdxToTypeGroup = All_IdxToTypeGroup[node.Name]

	grp := GetTypeGroup(node.Name)
	cntOfField := len(node.IdxToTypeGroup)

	vLen := 0
	tLen := 0
	tLenExt := 0

	if IsStructName[node.Name] {
		for i := 0; i < cntOfField; i++ {
			tLen += TypeToSize[node.IdxToType[i]]
		}
	} else if IsFieldSlice(grp) {
		tLen += 4
	} else {
		vLen += 4 + cntOfField*2
		tLen += 4
		for i := 0; i < cntOfField; i++ {
			grp := node.IdxToTypeGroup[i]
			if IsFieldBasicType(grp) {
				tLenExt += TypeToSize[node.IdxToType[i]]
			} else {
				tLenExt += 4
			}
		}
	}

	node.Node = NewNode2(NewBase(make([]byte, tLen+vLen, tLen+vLen+tLenExt)), vLen, true)

	flatbuffers.WriteUint32(node.R(node.Node.Pos, Size(4)), uint32(vLen))

	if vLen > 0 {
		flatbuffers.WriteUint16(node.R(0, Size(2)), uint16(vLen))
		flatbuffers.WriteUint16(node.R(2, Size(2)), uint16(tLen))
	}
	return nil
}

func (node *CommonNode) insertVTable(idx, size int) int {
	if node.TLen != uint16(node.TableLen()) {
		node.TLen = uint16(node.TableLen())
	}
	log.Log(log.LOG_DEBUG, func() log.LogArgs {
		return log.F(
			"node.insertVTable(%d, %d) before info=\n%s\n",
			idx, size, node.DumpTableWithVTable(),
		)
	})

	node.InsertBuf(node.Node.Pos+int(node.TLen), size)
	vPos := node.Node.Pos - int(flatbuffers.GetUOffsetT(node.R(node.Node.Pos, Size(4))))

	// write VTable[idx]
	flatbuffers.WriteVOffsetT(node.U(vPos+4+idx*2, 2), flatbuffers.VOffsetT(node.TLen))
	// Write TLen in VTable
	flatbuffers.WriteVOffsetT(node.U(vPos+2, 2), flatbuffers.VOffsetT(node.TLen+uint16(size)))

	wPos := node.Node.Pos + int(node.TLen)

	for i := 0; i < node.CountOfField(); i++ {
		if i == idx {
			continue
		}
		g := node.IdxToTypeGroup[i]
		if !IsFieldSlice(g) && (IsFieldBasicType(g) || IsFieldStruct(g)) {
			continue
		}

		if node.Table(i) >= wPos {
			// pos := node.VirtualTable(i)
			// offset := flatbuffers.GetUint32(node.R(pos))
			// flatbuffers.WriteUint32(node.U(pos, 4), offset+uint32(size))
			node.movePos(i, wPos, size)
		}
	}
	log.Log(log.LOG_DEBUG, func() log.LogArgs {
		return log.F(
			"node.insertVTable(%d, %d) after info=\n%s\n",
			idx, size, node.DumpTableWithVTable(),
		)
	})

	return wPos

}

// IsSameType ... return true if type information of field collect.
//FIXME: should improve for performance
func (node *CommonNode) IsSameType(idx int, fNode *CommonNode) bool {
	if All_IdxToName[node.Name][idx] == fNode.Name {
		return true
	}

	return IsFieldUnion(node.IdxToTypeGroup[idx]) &&
		loncha.Contain(UnionAlias[All_IdxToName[node.Name][idx]], func(i int) bool {
			return UnionAlias[All_IdxToName[node.Name][idx]][i] == fNode.Name
		})
}

// SetFieldAt ... set Node to current Node field/attribute
func (node *CommonNode) SetFieldAt(idx int, fNode *CommonNode) error {

	if len(node.IdxToTypeGroup) <= idx {
		return ERR_INVALID_INDEX
	}
	if !node.IsSameType(idx, fNode) {
		return ERR_INVALID_TYPE
	}

	//	grp := node.IdxToTypeGroup[idx]

	if IsStructName[node.Name] {
		pos := node.Node.Pos
		for i := 0; i < node.CountOfField(); i++ {
			if i < idx {
				pos += TypeToSize[node.IdxToType[i]]
			} else {
				size := TypeToSize[node.IdxToType[i]]
				node.Copy(fNode.IO, fNode.Node.Pos, size, pos, 0)
				return nil
			}
		}
	}

	g := node.IdxToTypeGroup[idx]

	if IsFieldSlice(g) {
		size := 4
		if fNode.Node.Size <= 0 {
			fNode.Node.Size = (*List)(fNode).InfoSlice().Size //fNode.ValueInfo(idx).Size
			if fNode.Node.Size < 0 {
				panic("!!!")
			}
		}

		var oFieldSize int = 0
		if node.VirtualTableIsZero(idx) {
			// update VTable in buffer. wPos is postion in table,
			wPos := node.insertVTable(idx, size)
			// write offset to vLen (list start)
			flatbuffers.WriteUint32(node.U(wPos, 4), uint32(4))
			// insert space (vLen + data )
			//node.InsertBuf(wPos+4, 4 + fNode.Node.Size)
			node.preLoadVtable()
		} else {
			oFieldSize = node.ValueInfo(idx).Size
		}
		extend := fNode.Node.Size - oFieldSize
		dstPos := node.Table(idx)
		node.Copy(fNode.IO,
			fNode.NodeList.ValueInfo.Pos-4, fNode.NodeList.ValueInfo.Size+4,
			dstPos, extend)

		return nil
		//node.Copy(fNode.Base, fNode.ValueInfo.Pos-vSize, fNode.Node.Size+vSize, dstPos, extend)
	} else if IsFieldStruct(g) || IsFieldBasicType(g) {
		var size int

		if IsFieldStruct(g) {
			size = fNode.Info().Size
		} else {
			size = TypeToSize[node.IdxToType[idx]]
		}

		node.preLoadVtable()

		if !node.VirtualTableIsZero(idx) {
			//if node.VTable[idx] > 0 {
			node.Copy(fNode.IO,
				fNode.Node.Pos, size,
				node.VirtualTable(idx), 0)

			return nil
		}
		wPos := node.insertVTable(idx, size)
		node.Copy(fNode.IO, fNode.Node.Pos, size, wPos, 0)

		node.preLoadVtable()
		return nil

	} else if IsFieldUnion(g) || IsFieldTable(g) {
		size := 4

		if fNode.Node.Size <= 0 {
			fNode.Node.Size = fNode.Info().Size //fNode.ValueInfo(idx).Size
			if fNode.Node.Size <= 0 {
				panic("!!!")
			}
		}
		vSize := fNode.CountOfField()*2 + 4
		oFieldSize := node.ValueInfo(idx).Size
		if node.VirtualTableIsZero(idx) {
			// update VTable in buffer.
			wPos := node.insertVTable(idx, size)
			// write tlen in Vtable
			node.InsertBuf(wPos+4, vSize+fNode.Node.Size)
			flatbuffers.WriteUint32(node.U(wPos, 4), uint32(vSize)+4) // +2 require ?
			node.Copy(fNode.IO, fNode.Node.Pos-vSize, fNode.Node.Size+vSize,
				wPos+4, fNode.Node.Size+vSize)
			node.preLoadVtable()
			return nil
		}
		extend := fNode.Node.Size - oFieldSize //node.ValueInfo(idx).Size
		//FIXME: should support shrink buffer
		if extend < 0 {
			extend = 0
		}
		dstPos := node.Table(idx) - vSize
		node.Copy(fNode.IO, fNode.Node.Pos-vSize, fNode.Node.Size+vSize, dstPos, extend)

		return nil
	}
	return ERR_NO_SUPPORT
}

// FromBytes ... return fbs data for []byte
// Deprecated: use FromByteList()
func FromBytes(bytes []byte) *CommonNode {

	buf := make([]byte, len(bytes)+4)
	flatbuffers.WriteUint32(buf, uint32(len(bytes)))
	copy(buf[4:], bytes)

	common := &CommonNode{}
	common.NodeList = &NodeList{}
	common.Name = "[]byte"
	common.Node = NewNode2(NewBase(buf), 4, true)
	common.NodeList.ValueInfo.Pos = 4
	common.NodeList.ValueInfo.VLen = uint32(len(bytes))
	common.NodeList.ValueInfo.Size = len(bytes)

	return common
}

// DumpTableWithVTable ... Dump flatbuffers Table and VTable information.
func (node *CommonNode) DumpTableWithVTable() string {

	g := GetTypeGroup(node.Name)
	if IsFieldBasicType(g) || IsFieldStruct(g) || IsFieldSlice(g) {
		return "{}"
	}

	var b strings.Builder

	vPos := node.Node.Pos - int(flatbuffers.GetUOffsetT(node.R(node.Node.Pos, Size(4))))
	//vPos := node.Node.Pos - int(vOff)

	fmt.Fprintf(&b, "Node=%s\n\tVTable{Pos:\t%d, VLen:\t\t%d, TLen:\t\t%d, \n\t\t",
		node.Name, vPos, node.VirtualTableLen(), node.TableLen())

	for i := 0; i < node.CountOfField(); i++ {
		if node.LenBuf() <= vPos+4+i*4 {
			fmt.Fprintf(&b, "Off(%d):\t\tempty,", i)
			continue
		}
		fmt.Fprintf(&b, "Off(%d):\t\t\t%d\t,  ",
			i, flatbuffers.GetUint16(node.R(int(vPos)+4+i*2, Size(2))))
	}
	fmt.Fprintf(&b, "}\n\t")

	min := func(i, j int) int {
		if i < j {
			return i
		}
		return j
	}

	fmt.Fprintf(&b, " Table{Pos:\t%d, \n\t\t", node.Node.Pos)
	for i := 0; i < node.CountOfField(); i++ {
		pos := node.VirtualTable(i)
		fmt.Fprintf(&b, "At(idx=%d,Pos=%d):\t%v, ",
			i, pos, node.R(pos)[:min(4, len(node.R(pos)))])
	}
	fmt.Fprintf(&b, "}\n")
	return b.String()

}

// SelfAsCommonNode ... return self CommonNode used by trick genny.
func (node *CommonNode) SelfAsCommonNode() *CommonNode {
	return node
}

// Equal ... return boolan if value is same.
func (n1 *CommonNode) Equal(n2 *CommonNode) bool {

	isSlice := func(n *CommonNode) bool {
		return n.NodeList.ValueInfo.Pos > 0
	}
	if isSlice(n1) {
		for i := 0; i < int(n1.NodeList.ValueInfo.VLen); i++ {
			l1, err := (*List)(n1).At(i)
			if err != nil {
				return false
			}
			l2, err := (*List)(n2).At(i)
			if err != nil {
				return false
			}
			if !l1.Equal(l2) {
				return false
			}
		}
		return true
	}

	n1cnt := n1.CountOfField()
	_ = n1cnt
	for i := 0; i < n1.CountOfField(); i++ {
		if n1.IdxToTypeGroup[i] != n2.IdxToTypeGroup[i] {
			return false
		}
		if IsFieldUnion(n1.IdxToTypeGroup[i]) {
			if !n1.FollowUnion(i).Equal(n2.FollowUnion(i)) {
				return false
			}
			continue
		}

		if !IsFieldBasicType(n1.IdxToTypeGroup[i]) {
			if !n1.FieldAt(i).Equal(n2.FieldAt(i)) {
				return false
			}
			continue
		}
		size := TypeToSize[n1.IdxToType[i]]

		f1 := n1.FieldAt(i)
		f2 := n2.FieldAt(i)

		// for debug
		if f1.Node == nil || f2.Node == nil {
			n1.FieldAt(i)
			n2.FieldAt(i)
		}

		if len(f1.R(f1.Node.Pos)) < size {
			f1.R(f1.Node.Pos, Size(1))
			return false
		}
		if len(f2.R(f2.Node.Pos)) < size {
			return false
		}

		if !bytes.Equal(f1.R(f1.Node.Pos, Size(size))[:size], f2.R(f2.Node.Pos, Size(size))[:size]) {
			return false
		}
	}
	return true
}
