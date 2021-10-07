package base

import (
	"io"
	"math/rand"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/kazu/fbshelper/query/log"
	"github.com/kazu/loncha"
)

// CommonList ... return List for separated data area.
type CommonList struct {
	*CommonNode
	dataW io.Writer
	dCur  int
	dLen  int
}

// DataLen ... return data.
func (l *CommonList) DataLen() int {
	return l.dLen
}

// SetDataWriter ... set io.Writer for writing data area.
func (l *CommonList) SetDataWriter(w io.Writer) {
	l.dCur = 0
	l.dataW = w
	l.dLen = 0
}

// DataWriter ... return DataWriter
func (l *CommonList) DataWriter() io.Writer {

	return l.dataW
}

// WriteDataAll ... write all data to datawriter
func (l *CommonList) WriteDataAll() (e error) {

	if l.CommonNode == nil || l.Count() == 0 {
		return
	}
	l.Merge()

	first, e := (*List)(l.CommonNode).At(0)
	if e != nil {
		return e
	}

	vSize := first.CountOfField()*2 + 4
	first.Node.Size = first.Info().Size

	g := GetTypeGroup(first.Name)
	if !IsFieldTable(g) {
		return log.ERR_NO_SUPPORT
	}

	for i, elm := range (*List)(l.CommonNode).All() {
		pos := elm.Node.Pos - vSize
		elm.Node.Size = elm.Info().Size
		size := elm.Node.Size + vSize
		l.dataW.Write(l.R(pos)[:size])
		l.dCur = i
		l.dLen += size
	}
	last := l.LenBuf()
	_ = last

	//remove writed data area
	bytes := l.R(0)
	bytes = bytes[0 : first.Node.Pos-vSize : first.Node.Pos-vSize]
	l.Base = l.Base.NewFromBytes(bytes)

	return nil

}

// At ... return Element of list
func (l *CommonList) At(i int) (*CommonNode, error) {
	if l.dataW == nil {
		return (*List)(l.CommonNode).At(i)
	}
	// FIXME
	return nil, log.ERR_NO_SUPPORT
}

// SetAt ... Set Element to list
func (l *CommonList) SetAt(i int, elm *CommonNode) error {
	if l.dataW == nil {
		return (*List)(l.CommonNode).SetAt(i, elm)
	}
	if l.Count() != i {
		return log.ERR_NO_SUPPORT
	}

	vSize := elm.CountOfField()*2 + 4
	elm.Node.Size = elm.Info().Size
	pos := elm.Node.Pos - vSize
	size := elm.Node.Size + vSize
	//l.dataW.Write(elm.R(pos)[:size])

	//elm.Merge()

	(*List)(l.CommonNode).SetAt(i, elm)
	last := l.LenBuf()
	_ = last
	l.WriteElm(elm, pos, size)
	l.dCur = i

	return nil
}

func (l *CommonList) dataStart() int {

	return l.NodeList.ValueInfo.Pos + int((*List)(l.CommonNode).VLen())*4

}

// WriteElm ... Set Element to list
func (l *CommonList) WriteElm(elm *CommonNode, pos, size int) {

	defer func() {
		l.dLen += size
	}()
	if len(l.GetDiffs()) == 0 {
		l.dataW.Write(elm.R(pos)[:size])
		//remove writed data area
		bytes := l.R(0)
		if len(bytes) >= pos+size {
			bytes = bytes[0 : pos+size : pos+size]
			l.Base = NewBase(bytes)
		}
		return
	}
	wStart := l.dataStart() + l.dLen

	if w, ok := l.dataW.(io.WriterAt); ok {

		//w.WriteAt(elm.R(0), int64(l.dLen+pos))
		cur := l.dLen //wStart

		w.WriteAt(elm.R(0), int64(cur))

		for i := range elm.GetDiffs() {
			bytes := elm.GetDiffs()[i].bytes
			offset := elm.GetDiffs()[i].Offset
			if len(bytes)+offset > size {
				log.Log(LOG_ERROR, func() log.LogArgs {
					return log.F("WriteElm: invalid  len(bytes)=%d + offset=%d > size=%d\n", len(bytes), offset, size)
				})
			}
			w.WriteAt(bytes, int64(cur+pos+offset))
		}
	}
	// remove written data
	diffs := l.GetDiffs()
	loncha.Delete(&diffs, func(i int) bool {
		return wStart+pos+size <= diffs[i].Offset+len(diffs[i].bytes) || wStart <= diffs[i].Offset
	})
	l.SetDiffs(diffs)
	return
}

// Add ... if src.Name and dst.Name is the same and list, join header
//       only support element is table .
func (dst *CommonList) Add(src *CommonList) (nList *CommonList, e error) {

	if dst.Name != src.Name {
		return nil, ERR_NO_SUPPORT
	}

	g := GetTypeGroup(dst.Name)

	if !IsFieldTable(g) {
		return nil, ERR_NO_SUPPORT
	}

	if dst.dataW == nil || src.dataW == nil {
		return nil, ERR_NO_SUPPORT
	}

	newLen := (*List)(dst.CommonNode).VLen() + (*List)(src.CommonNode).VLen()

	nList = &CommonList{
		CommonNode: &CommonNode{},
		dataW:      dst.dataW,
		dCur:       dst.dCur + src.dCur,
		dLen:       dst.dLen + src.dLen,
	}
	nList.NodeList = &NodeList{}
	nList.CommonNode.Name = dst.Name
	(*List)(nList.CommonNode).InitList()

	//	a := dst.New(make([]byte, 4+int(newLen)*4))
	nList.Base = dst.NewFromBytes(make([]byte, 4+int(newLen)*4))
	flatbuffers.WriteUint32(nList.U(0, 4), newLen)

	cur2 := 0
	_ = cur2
	for i := 0; i < int((*List)(dst.CommonNode).VLen()); i++ {
		cur2 = 4 + i*4
		flatbuffers.WriteUint32(nList.U(4+i*4, 4),
			flatbuffers.GetUint32(dst.R(dst.NodeList.ValueInfo.Pos+i*4))+(*List)(src.CommonNode).VLen()*4)
	}

	cur := int((*List)(dst.CommonNode).VLen()) * 4

	for i := 0; i < int((*List)(src.CommonNode).VLen()); i++ {
		flatbuffers.WriteUint32(nList.U(cur+i*4+4, 4),
			flatbuffers.GetUint32(src.R(src.NodeList.ValueInfo.Pos+i*4))+uint32(dst.dLen))
	}
	nList.Node.Pos = 4
	nList.NodeList.ValueInfo = dst.NodeList.ValueInfo
	nList.NodeList.ValueInfo.VLen = (*List)(nList.CommonNode).VLen()
	nList.NodeList.ValueInfo.Size += src.NodeList.ValueInfo.Size - 4

	seeker, ok := dst.dataW.(io.Seeker)
	if !ok {
		return nil, log.ERR_NO_SUPPORT
	}
	n, e := seeker.Seek(0, io.SeekEnd)
	//	n, e := seeker.Seek(int64(dst.dLen), io.SeekStart)
	_ = n
	if e != nil {
		return nil, e
	}

	srcDataR, ok := src.dataW.(io.Reader)
	if !ok {
		return nil, log.ERR_NO_SUPPORT
	}

	written, e := io.Copy(dst.dataW, srcDataR)
	if e != nil {
		return nil, e
	}
	if written != int64(src.dLen) {
		return nil, log.ERR_INVLIAD_WRITE_SIZE
	}

	return
}

// VLen ... return Len of Virtual Table.
func (l *CommonList) VLen() uint32 {
	return (*List)(l.CommonNode).VLen()
}

// List ... List Node as CommonNode.
type List CommonNode

// At ... return Element of list
func (node *List) At(i int) (*CommonNode, error) {

	if !(*CommonNode)(node).IsList() {
		return nil, ERR_NO_SUPPORT
	}
	if i >= int(node.VLen()) || i < 0 {
		//if i >= int(node.NodeList.ValueInfo.VLen) || i < 0 {
		return nil, ERR_INVALID_INDEX
	}

	tName := node.Name[2:]
	grp := GetTypeGroup(tName)

	ptr := int(node.NodeList.ValueInfo.Pos) + i*4
	if i > 0 && (IsFieldBasicType(grp) || IsFieldStruct(grp)) {
		first, _ := node.First()
		size := first.Info().Size
		ptr = int(node.NodeList.ValueInfo.Pos) + i*size
	}

	var nNode *Node
	if IsFieldBasicType(grp) || IsFieldStruct(grp) {
		nNode = NewNode2(node.Base, ptr, true)
	} else {
		_ = node.R(ptr + 3)
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

// First ... First Element in List
func (node *List) First() (*CommonNode, error) {
	return node.At(0)
}

// Last ... Last Element in List
func (node *List) Last() (*CommonNode, error) {
	return node.At(int(node.VLen()) - 1)
}

// Select ... Select Elements by condtion function
func (node *List) Select(fn func(m *CommonNode) bool) []*CommonNode {
	result := make([]*CommonNode, 0, int(node.NodeList.ValueInfo.VLen))
	for i := 0; i < int(node.NodeList.ValueInfo.VLen); i++ {
		if m, err := node.At(i); err == nil && fn(m) {
			result = append(result, m)
		}
	}
	return result
}

// Find ... Find Element by condtion function
func (node *List) Find(fn func(m *CommonNode) bool) *CommonNode {

	for i := 0; i < int(node.NodeList.ValueInfo.VLen); i++ {
		if m, e := node.At(i); e == nil && fn(m) {
			return m
		}
	}

	return nil
}

// All ... All Element by condtion function
func (node *List) All() []*CommonNode {
	return node.Select(func(m *CommonNode) bool { return true })
}

// VLen ... return Length of flatbuffers's Virtual Table
func (node *List) VLen() uint32 {
	return flatbuffers.GetUint32(node.R(node.NodeList.ValueInfo.Pos - flatbuffers.SizeUOffsetT))
}

// InfoSlice ... return infomation of List
func (node *List) InfoSlice() Info {
	info := Info{Pos: node.NodeList.ValueInfo.Pos, Size: -1}
	tName := node.Name[2:]
	var vInfo Info

	grp := GetTypeGroup(tName)
	if IsFieldBasicType(grp) {
		ptr := int(node.NodeList.ValueInfo.Pos)
		size := int(node.NodeList.ValueInfo.VLen) * TypeToSize[NameToTypeEnum(tName)]
		vInfo = Info{Pos: ptr, Size: size}
	} else if IsFieldBytes(grp) {
		ptr := int(node.NodeList.ValueInfo.Pos) + (int(node.NodeList.ValueInfo.VLen)-1)*4
		vInfo = FbsStringInfo(NewNode(node.Base, ptr+int(flatbuffers.GetUint32(node.R(ptr)))))
	} else {
		if elm, err := node.Last(); err == nil {
			vInfo = elm.Info()
		}
	}
	info.VLen = flatbuffers.GetUint32(node.R(node.NodeList.ValueInfo.Pos - 4))
	if info.VLen == 0 {
		info.Size = 0
	}
	if info.Pos+info.Size < vInfo.Pos+vInfo.Size {
		info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
	}
	return info
}

// SearchInfoSlice ... search list information
func (node *List) SearchInfoSlice(pos int, fn RecFn, condFn CondFn) {

	info := (*CommonNode)(node).Info()

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

// SetAt ... Set Element to list
func (node *List) SetAt(idx int, elm *CommonNode) error {

	if idx > int(node.NodeList.ValueInfo.VLen) {
		return ERR_INVALID_INDEX
	}

	if node.NodeList.ValueInfo.Pos == 0 || node.NodeList.ValueInfo.VLen == 0 || node.NodeList.ValueInfo.Size == 0 {
		node.NodeList.ValueInfo = ValueInfo(node.InfoSlice())
	}
	if node.NodeList.ValueInfo.VLen == 0 {
		node.NodeList.ValueInfo.Size = 0
	}
	vlen := int(node.NodeList.ValueInfo.VLen)
	total := node.NodeList.ValueInfo.Size

	g := GetTypeGroup(elm.Name)
	//ptr := int(node.NodeList.ValueInfo.Pos) + vlen*4

	if IsFieldBasicType(g) || IsFieldStruct(g) {
		// new element
		if elm.Node.Size <= 0 {
			elm.Node.Size = elm.Info().Size
		}
		ptr := int(node.NodeList.ValueInfo.Pos) + idx*elm.Node.Size
		//node.insertBuf(node.NodeList.ValueInfo.Pos + total, elm.Node.Size)
		extend := 0
		if vlen == idx {
			extend = elm.Node.Size
		}
		node.Copy(
			elm.Base, elm.Node.Pos, elm.Node.Size,
			ptr, extend)

		total += extend
		if vlen == idx {
			vlen++
			flatbuffers.WriteUint32(node.U(node.NodeList.ValueInfo.Pos-4, 4), uint32(vlen))
		}
		node.NodeList.ValueInfo = ValueInfo(node.InfoSlice())
		return nil
	} else if IsFieldUnion(g) {
		return ERR_NO_SUPPORT
	} else if IsFieldTable(g) {
		// new element
		vSize := elm.CountOfField()*2 + 4
		oSize := 0
		var oElm *CommonNode
		if idx < vlen {
			oElm, _ = node.At(idx)
			oSize = oElm.Info().Size
		}
		if elm.Node.Size <= 0 {
			elm.Node.Size = elm.Info().Size
		}
		ptrIdx := func(idx int) int {
			return int(node.NodeList.ValueInfo.Pos) + idx*4
		}
		ptr := ptrIdx(idx)

		header := make([]byte, 4)

		//flatbuffers.WriteUint32(off, uint32(total-vlen*4+vSize))

		header_extend := 0
		body_extend := 0
		if vlen == idx {
			header_extend = 4
		}

		if header_extend > 0 {
			body_extend = elm.Node.Size + vSize
		} else {
			body_extend = elm.Node.Size - oSize
		}
		if body_extend < 0 {
			body_extend = 0
		}
		dstPos := 0
		if header_extend > 0 {
			flatbuffers.WriteUint32(header, uint32(total-vlen*4+vSize))
			(*CommonNode)(node).InsertBuf(ptr, 4)
			for i := 0; i < vlen; i++ {
				dataPtr := ptrIdx(i)
				_ = dataPtr
				off := flatbuffers.GetUint32(node.R(ptrIdx(i)))
				off += 4
				flatbuffers.WriteUint32(node.U(ptrIdx(i), 4), off)
			}
			dstPos = ptrIdx(0) + total + header_extend
			flatbuffers.WriteUint32(node.U(ptrIdx(idx), 4), uint32(dstPos+vSize-ptr))
			vlen++
		} else {
			dstPos = oElm.Node.Pos - vSize
		}

		if body_extend > 0 {
			(*CommonNode)(node).InsertSpace(dstPos, body_extend, false)
			if header_extend == 0 {
				for i := idx + 1; i < vlen; i++ {
					off := flatbuffers.GetUint32(node.R(ptrIdx(i)))
					off += uint32(body_extend)
					flatbuffers.WriteUint32(node.U(ptrIdx(i), 4), off)
				}
			}
		}
		// update vlen
		flatbuffers.WriteUint32(node.U(ptrIdx(-1), 4), uint32(vlen))
		node.Copy(elm.Base,
			elm.Node.Pos-vSize, elm.Node.Size+vSize,
			dstPos, 0)
		node.NodeList.ValueInfo = ValueInfo(node.InfoSlice())
		return nil
	}

	return ERR_NO_SUPPORT
}

// InitList ... initlize List.
func (node *List) InitList() error {

	// write vLen == 0
	node.Node = NewNode2(NewBase(make([]byte, 4)), 4, true)
	node.NodeList.ValueInfo.Pos = 4
	return nil
}

func (node *List) isDirectList() bool {
	return node.isDirect()
}

func (node *List) isDirect() bool {

	if !node.IsList() {
		return false
	}

	tName := node.Name[2:]

	grp := GetTypeGroup(tName)
	if IsFieldBasicType(grp) || IsFieldStruct(grp) {
		return true
	}

	return false
}

func (node *List) offsetOfList() int {

	return node.offsetOf()
}

func (node *List) offsetOf() int {

	if !node.IsList() {
		return -1
	}

	if !node.isDirectList() {
		return 4
	}

	first, _ := node.First()
	return first.Info().Size

}

// SwapAt ... swap data in i and j in List.
func (node *List) SwapAt(i, j int) error {
	if !node.IsList() {
		return log.ERR_NO_SUPPORT
	}

	if i >= int(node.VLen()) || j >= int(node.VLen()) {
		return log.ERR_INVALID_INDEX
	}

	if node.isDirectList() {
		size := node.offsetOfList()
		node.Flatten()

		iPtr := int(node.NodeList.ValueInfo.Pos) + i*size
		jPtr := int(node.NodeList.ValueInfo.Pos) + j*size

		tmp := make([]byte, size)
		iData := node.U(iPtr, size)[:size:size]
		jData := node.U(jPtr, size)[:size:size]
		copy(tmp, iData)
		copy(iData, jData)
		copy(jData, tmp)

		return nil
	}

	size := flatbuffers.SizeUOffsetT

	iPtr := int(node.NodeList.ValueInfo.Pos) + i*size
	jPtr := int(node.NodeList.ValueInfo.Pos) + j*size

	iOffset := int(flatbuffers.GetUint32(node.R(iPtr)))
	jOffset := int(flatbuffers.GetUint32(node.R(jPtr)))
	iOffset -= (jPtr - iPtr)
	jOffset -= (iPtr - jPtr)

	flatbuffers.WriteUint32(node.U(iPtr, size), uint32(jOffset))
	flatbuffers.WriteUint32(node.U(jPtr, size), uint32(iOffset))

	return nil

}

// SortBy ... sort of List
func (node *List) SortBy(less func(i, j int) bool) error {
	if !node.IsList() {
		return log.ERR_NO_SUPPORT
	}

	left, right := 0, int(node.VLen()-1)

	return node.quicksort(left, right, less)

}

func (node *List) quicksort(left, right int, less func(i, j int) bool) error {
	if !node.IsList() {
		return log.ERR_NO_SUPPORT
	}

	len := (right - left + 1)
	if len < 2 {
		return nil
	}

	pivot := rand.Int() % len

	node.SwapAt(pivot, right)

	for i := 0; i < len; i++ {
		if less(i, right) {
			node.SwapAt(left, i)
			left++
		}
	}
	node.SwapAt(left, right)

	node.quicksort(0, left-1, less)
	node.quicksort(left+1, len-2-left, less)

	return nil
}

// SearchIndex ... binary search
// copy/modify from golang.org/src/sort/search.go
func (node *List) SearchIndex(n int, fn func(c *CommonNode) bool) int {
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if cNode, err := node.At(h); err == nil && !fn(cNode) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

// IsList ... return true if List is true
func (node *List) IsList() bool {
	return (*CommonNode)(node).IsList()
}

// SelfAsCommonNode ... return self CommonNode used by trick genny.
func (node *List) SelfAsCommonNode() *CommonNode {
	return (*CommonNode)(node).SelfAsCommonNode()
}

// Count ... return count of element in List
func (node *List) Count() int {
	// MENTION: if dosent work, enable comment out routine
	//return int(node.NodeList.ValueInfo.VLen)
	return int(node.VLen())
}

// FromByteList ... []byte in flatbuffers.
func FromByteList(bytes []byte) *List {

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

	return (*List)(common)
}
