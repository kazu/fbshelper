package base

import (
	"io"
	"reflect"

	flatbuffers "github.com/google/flatbuffers/go"
	//. "github.com/kazu/fbshelper/query/error"
)

type CommonNode struct {
	*NodeList
	Name           string
	IdxToType      map[int]int
	IdxToTypeGroup map[int]int
}

func (node *CommonNode) Info() (info Info) {

	nName := node.Name

	if node.NodeList == nil || node.Node == nil {
		//FIXME
		return info
	}
	info.Pos = node.Node.Pos
	info.Size = -1
	if IsStructName[nName] {
		size := 0
		for i := 0; i < len(node.IdxToTypeGroup); i++ {
			size += TypeToSize[node.IdxToType[i]]
		}
		info.Size = size
		return info
	}
	for i := 0; i < len(node.VTable); i++ {
		vInfo := node.ValueInfo(i)
		if info.Pos+info.Size < vInfo.Pos+vInfo.Size {
			info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
		}
	}
	return info
}

func (node *CommonNode) ValueInfo(idx int) ValueInfo {

	if IsStructName[node.Name] {
		if len(node.ValueInfos) > idx {
			return node.ValueInfos[idx]
		}
		node.ValueInfos = make([]ValueInfo, 0, node.CountOfField())
		info := ValueInfo{Pos: node.Node.Pos, Size: 0}
		for i := 0; i < node.CountOfField(); i++ {
			info.Pos += info.Size
			info.Size = TypeToSize[node.IdxToType[i]]
			node.ValueInfos = append(node.ValueInfos, info)
		}
	}

	grp := node.IdxToTypeGroup[idx]

	if IsFieldStruct(grp) {
		if node.ValueInfos[idx].IsNotReady() {
			node.ValueInfoPos(idx)
		}
		node.ValueInfos[idx].Size = node.FieldAt(idx).Info().Size

	} else if IsFieldUnion(grp) {
		if node.ValueInfos[idx].IsNotReady() {
			node.ValueInfoPosTable(idx)
		}
		if node.CanFollow(idx) {
			node.ValueInfos[idx].Size = node.FollowUnion(idx).Info().Size
		}

	} else if IsFieldBytes(grp) {
		if node.ValueInfos[idx].IsNotReady() {
			node.ValueInfoPosBytes(idx)
		}

	} else if IsFieldSlice(grp) {
		if node.ValueInfos[idx].IsNotReady() {
			node.ValueInfoPosList(idx)
		}

		node.ValueInfos[idx].Size = node.FieldAt(idx).InfoSlice().Size

	} else if IsFieldTable(grp) {
		if node.ValueInfos[idx].IsNotReady() {
			node.ValueInfoPosTable(idx)
		}

		node.ValueInfos[idx].Size = node.FieldAt(idx).Info().Size

	} else if IsFieldBasicType(grp) {
		if node.ValueInfos[idx].IsNotReady() {
			node.ValueInfoPos(idx)
		}

		node.ValueInfos[idx].Size = TypeToSize[node.IdxToType[idx]]

	} else {
		Log(LOG_ERROR, func() LogArgs {
			return F("Invalid %s.%s idx=%d\n", node.Name, All_IdxToName[node.Name][idx], idx)
		})
	}

	return node.ValueInfos[idx]
}
func (node *CommonNode) CanFollow(idx int) bool {

	return node.FieldAt(idx-1).Byte() > 0
}

func (node *CommonNode) FollowUnion(idx int) *CommonNode {
	idxOfAlias := node.FieldAt(idx-1).Byte() - 1
	if int(idxOfAlias) >= len(UnionAlias[All_IdxToName[node.Name][idx]]) {
		Log(LOG_WARN, func() LogArgs {
			return F("Invalid union=%s aliases=%+v idx=%d\n",
				node.Name, All_IdxToName[node.Name], idxOfAlias)
		})
		return nil
	}

	newName := UnionAlias[All_IdxToName[node.Name][idx]][idxOfAlias]

	return &CommonNode{
		NodeList:       &NodeList{Node: node.Node},
		Name:           newName,
		IdxToType:      All_IdxToType[newName],
		IdxToTypeGroup: All_IdxToTypeGroup[newName],
	}

}

func (node *CommonNode) SearchInfo(pos int, fn RecFn, condFn CondFn) {
	info := node.Info()

	if condFn(pos, info) {
		fn(NodePath{Name: node.Name, Idx: -1}, info)
	} else {
		return
	}

	for i := 0; i < len(node.IdxToTypeGroup); i++ {
		g := node.IdxToTypeGroup[i]

		if node.IsLeafAt(node.IdxToTypeGroup[i]) {
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
			node.FieldAt(i).SearchInfoSlice(pos, fn, condFn)
		} else if IsMatchBit(g, FieldTypeSlice) {
			node.FieldAt(i).SearchInfoSlice(pos, fn, condFn)
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

func (node *CommonNode) CountOfField() int {
	return len(node.IdxToType)
}

func (node *CommonNode) FieldAt(idx int) (cNode *CommonNode) {
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
		node.Node.Pos = pos
		result.Node = node.Node

		goto RESULT
	}

	if node.VTable[idx] == 0 {
		goto RESULT
	}

	if IsFieldStruct(grp) {
		result.Node = node.ValueStruct(idx)

	} else if IsFieldUnion(grp) {
		result.Node = node.ValueTable(idx)

	} else if IsFieldBytes(grp) {
		if node.ValueInfos[idx].Pos < 1 {
			node.ValueInfoPosBytes(idx)
		}
		valInfo := node.ValueInfos[idx]

		nNode := NewNode2(node.Base, node.ValueInfos[idx].Pos, true)
		nNode.Size = valInfo.Size
		result.Node = nNode

	} else if IsFieldSlice(grp) {
		nodeList := node.ValueList(idx)
		result = &nodeList

	} else if IsFieldTable(grp) {
		result.Node = node.ValueTable(idx)

	} else if IsFieldBasicType(grp) {
		if node.ValueInfos[idx].Pos < 1 {
			node.ValueInfoPos(idx)
		}
		result.Node = NewNode2(node.Base, node.ValueInfos[idx].Pos, true)
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
	return cNode
}

func (node *CommonNode) FetchIndex() {
	if _, ok := All_IdxToName[node.Name]; ok {
		node.IdxToType = All_IdxToType[node.Name]
		node.IdxToTypeGroup = All_IdxToTypeGroup[node.Name]
	}
}

func (node *CommonNode) IsList() bool {
	return node.NodeList.ValueInfo.VLen > 0
}

func (node *CommonNode) At(i int) (*CommonNode, error) {

	if !node.IsList() {
		return nil, ERR_NO_SUPPORT
	}

	if i >= int(node.NodeList.ValueInfo.VLen) || i < 0 {
		return nil, ERR_INVALID_INDEX
	}

	tName := node.Name[2:]
	grp := GetTypeGroup(tName)

	ptr := int(node.NodeList.ValueInfo.Pos) + i*4

	var nNode *Node
	if IsFieldBasicType(grp) {
		nNode = NewNode2(node.Base, ptr, true)
	} else {
		nNode = NewNode(node.Base, ptr+int(flatbuffers.GetUint32(node.R(ptr))))
	}

	cNode := &CommonNode{}
	cNode.NodeList = &NodeList{}
	cNode.Node = nNode
	cNode.Name = tName
	if _, ok := All_IdxToName[tName]; ok {
		cNode.IdxToType = All_IdxToType[tName]
		cNode.IdxToTypeGroup = All_IdxToTypeGroup[tName]
	}

	return cNode, nil

}

func (node *CommonNode) First() (*CommonNode, error) {
	return node.At(0)
}

func (node *CommonNode) Last() (*CommonNode, error) {
	return node.At(int(node.NodeList.ValueInfo.VLen) - 1)
}

func (node *CommonNode) Select(fn func(m *CommonNode) bool) []*CommonNode {
	result := make([]*CommonNode, 0, int(node.NodeList.ValueInfo.VLen))
	for i := 0; i < int(node.NodeList.ValueInfo.VLen); i++ {
		if m, err := node.At(i); err == nil && fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node *CommonNode) Find(fn func(m *CommonNode) bool) *CommonNode {

	for i := 0; i < int(node.NodeList.ValueInfo.VLen); i++ {
		if m, e := node.At(i); e == nil && fn(m) {
			return m
		}
	}

	return nil
}

func (node *CommonNode) All() []*CommonNode {
	return node.Select(func(m *CommonNode) bool { return true })
}

func (node *CommonNode) InfoSlice() Info {
	info := Info{Pos: node.NodeList.ValueInfo.Pos, Size: -1}
	tName := node.Name[2:]
	var vInfo Info

	grp := GetTypeGroup(tName)
	if IsFieldBasicType(grp) {
		ptr := int(node.NodeList.ValueInfo.Pos) + (int(node.NodeList.ValueInfo.VLen)-1)*4
		vInfo = Info{Pos: ptr, Size: TypeToSize[NameToTypeEnum(tName)]}
	} else if IsFieldBytes(grp) {
		ptr := int(node.NodeList.ValueInfo.Pos) + (int(node.NodeList.ValueInfo.VLen)-1)*4
		vInfo = FbsStringInfo(NewNode(node.Base, ptr+int(flatbuffers.GetUint32(node.R(ptr)))))
	} else {
		if elm, err := node.Last(); err == nil {
			vInfo = elm.Info()
		}
	}

	if info.Pos+info.Size < vInfo.Pos+vInfo.Size {
		info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
	}
	return info
}

func (node *CommonNode) SearchInfoSlice(pos int, fn RecFn, condFn CondFn) {

	info := node.Info()

	if condFn(pos, info) {
		fn(NodePath{Name: node.Name, Idx: -1}, info)
	} else {
		return
	}

	var v interface{}
	for _, cNode := range node.All() {
		v = cNode
		if vv, ok := v.(Searcher); ok {
			vv.SearchInfo(pos, fn, condFn)
		} else {
			goto NO_NODE
		}
	}
	return

NO_NODE:
	for i := 0; i < int(node.NodeList.ValueInfo.VLen); i++ {
		ptr := int(node.NodeList.ValueInfo.Pos) + i*4
		start := ptr + int(flatbuffers.GetUint32(node.R(ptr)))
		size := info.Size
		if i+1 < int(node.NodeList.ValueInfo.Pos) {
			size = ptr + 4 + int(flatbuffers.GetUint32(node.R(ptr+4))) - start
		}
		cInfo := Info{Pos: start, Size: size}
		if condFn(pos, info) {

			fn(NodePath{Name: node.Name, Idx: i}, cInfo)
		}
	}

}

func Open(r io.Reader, cap int) (node *CommonNode) {

	b := NewBaseByIO(r, 512)

	node = &CommonNode{}
	node.NodeList = &NodeList{}
	node.Node = NewNode(b, int(flatbuffers.GetUOffsetT(b.R(0))))
	return node
}

func OpenByBuf(buf []byte) *CommonNode {
	node := &CommonNode{}
	node.NodeList = &NodeList{}
	node.Node = NewNode(NewBase(buf), int(flatbuffers.GetUOffsetT(buf)))
	return node
}

func (node *CommonNode) Len() int {
	info := node.Info()
	size := info.Pos + info.Size

	if (size % 8) == 0 {
		return size
	}

	return size + (8 - (size % 8))
}

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
			rv.SetBytes(cNode.R(cNode.Node.Pos)[:cNode.Node.Size])
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
