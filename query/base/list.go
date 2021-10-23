package base

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/kazu/fbshelper/query/dump"
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
		l.dataW.Write(l.R(pos, Size(size))[:size])
		l.dCur = i
		l.dLen += size
	}
	last := l.LenBuf()
	_ = last

	//remove writed data area
	bytes := l.R(0, Size(first.Node.Pos-vSize))
	bytes = bytes[0 : first.Node.Pos-vSize : first.Node.Pos-vSize]
	l.IO = l.IO.NewFromBytes(bytes)

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
		l.dataW.Write(elm.R(pos, Size(size))[:size])
		//remove writed data area
		bytes := l.R(0, Size(pos+size))
		if len(bytes) >= pos+size {
			bytes = bytes[0 : pos+size : pos+size]
			l.IO = NewBase(bytes)
		}
		return
	}
	wStart := l.dataStart() + l.dLen

	if w, ok := l.dataW.(io.WriterAt); ok {

		//w.WriteAt(elm.R(0), int64(l.dLen+pos))
		cur := l.dLen //wStart

		//FIXME: set size ?
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
	nList.IO = dst.NewFromBytes(make([]byte, 4+int(newLen)*4))
	flatbuffers.WriteUint32(nList.U(0, 4), newLen)

	cur2 := 0
	_ = cur2
	for i := 0; i < int((*List)(dst.CommonNode).VLen()); i++ {
		cur2 = 4 + i*4
		flatbuffers.WriteUint32(nList.U(4+i*4, 4),
			flatbuffers.GetUint32(dst.R(dst.NodeList.ValueInfo.Pos+i*4, Size(4)))+(*List)(src.CommonNode).VLen()*4)
	}

	cur := int((*List)(dst.CommonNode).VLen()) * 4

	for i := 0; i < int((*List)(src.CommonNode).VLen()); i++ {
		flatbuffers.WriteUint32(nList.U(cur+i*4+4, 4),
			flatbuffers.GetUint32(src.R(src.NodeList.ValueInfo.Pos+i*4, Size(4)))+uint32(dst.dLen))
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

func (node *List) dup() (l *List) {

	l = &List{}
	l.NodeList = node.NodeList
	l.Name = node.Name
	l.IO = node.IO.Dup()
	l.NodeList.ValueInfo = ValueInfo(node.InfoSlice())

	return
}

func (node *List) toCommonNode() *CommonNode {
	return (*CommonNode)(node)
}

// AtWihoutError ... ignore error
func (node *List) AtWihoutError(i int) *CommonNode {
	r, _ := node.At(i)
	return r
}

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
		nNode = NewNode2(node.IO, ptr, true)
	} else {
		_ = node.R(ptr + 3)

		// recovered := false
		// defer func() {
		// 	if recovered {
		// 		return
		// 	}

		// 	if err := recover(); err != nil {
		// 		fmt.Printf("ptr=0x%x(%d) start=0x%x pos=0x%x\n", ptr, ptr,
		// 			ptr-node.NodeList.ValueInfo.Pos+4,
		// 			node.NodeList.ValueInfo.Pos,
		// 		)
		// 		dumpPos := node.NodeList.ValueInfo.Pos - 4
		// 		node.Dump(dumpPos, OptDumpSize(64))
		// 		panic(err)
		// 	}
		// }()

		nPos := ptr + int(flatbuffers.GetUint32(node.R(ptr, Size(4))))
		if node.R(nPos) == nil {
			return nil, ERR_NOT_FOUND
		}

		nNode = NewNode(node.IO, ptr+int(flatbuffers.GetUint32(node.R(ptr))))
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
	// info := node.InfoSlice()
	cnt := node.Count()
	_ = cnt
	// _ = info
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
	return flatbuffers.GetUint32(node.R(node.NodeList.ValueInfo.Pos-flatbuffers.SizeUOffsetT, Size(4)))
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
		vInfo = FbsStringInfo(NewNode(node.IO, ptr+int(flatbuffers.GetUint32(node.R(ptr, Size(4))))))
	} else {

		// vInfos := make([]Info, 0, node.Count())

		// node.Select(func(e *CommonNode) bool {
		// 	vInfos = append(vInfos, e.Info())
		// 	return true
		// })

		if elm, err := node.Last(); err == nil {
			vInfo = elm.Info()
		}
	}
	info.VLen = flatbuffers.GetUint32(node.R(node.NodeList.ValueInfo.Pos-4, Size(4)))
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
		start := ptr + int(flatbuffers.GetUint32(node.R(ptr, Size(4))))
		size := info.Size
		if i+1 < int(node.NodeList.ValueInfo.Pos) {
			size = ptr + 4 + int(flatbuffers.GetUint32(node.R(ptr+4, Size(4)))) - start
		}
		cInfo := Info{Pos: start, Size: size}
		if condFn(pos, info) {

			fn(NodePath{Name: node.Name, Idx: i}, cInfo)
		}
	}

}

func (list *List) Add(slist *List) error {

	elm, err := slist.First()
	if err != nil {
		//FIXME error
		return fmt.Errorf("(*List).Append(): fail slist.First() err=%s", err)
	}
	g := GetTypeGroup(elm.Name)

	if IsFieldBasicType(g) || IsFieldStruct(g) {
		return list.addStructList(slist)
	}

	if IsFieldUnion(g) {
		return ERR_NO_SUPPORT
	}

	if IsFieldTable(g) {
		return list.addTableList(slist)
	}

	return ERR_NO_SUPPORT
}

func (list *List) vlenAndTotals(slist *List) (vlens []int, totals []int) {

	vlens = make([]int, 0, 2)
	totals = make([]int, 0, 2)

	for _, node := range []*List{list, slist} {
		if node.NodeList.ValueInfo.Pos == 0 || node.NodeList.ValueInfo.VLen == 0 || node.NodeList.ValueInfo.Size == 0 {
			node.NodeList.ValueInfo = ValueInfo(node.InfoSlice())
		}
		if node.NodeList.ValueInfo.VLen == 0 {
			node.NodeList.ValueInfo.Size = 0
		}
		vlens = append(vlens, int(node.NodeList.ValueInfo.VLen))
		totals = append(totals, int(node.NodeList.ValueInfo.Size))
	}
	return
}

func (node *List) vlenTotal() (int, int) {
	if node.NodeList.ValueInfo.Pos == 0 || node.NodeList.ValueInfo.VLen == 0 || node.NodeList.ValueInfo.Size == 0 {
		node.NodeList.ValueInfo = ValueInfo(node.InfoSlice())
	}
	if node.NodeList.ValueInfo.VLen == 0 {
		node.NodeList.ValueInfo.Size = 0
	}

	return int(node.NodeList.ValueInfo.VLen), node.NodeList.ValueInfo.Size

}

func (list *List) addStructList(slist *List) error {

	for _, elm := range slist.All() {
		if e := list.setStructAt(list.Count(), elm); e != nil {
			return e
		}
	}

	return nil

}

func VlensTotals(lists []*List) ([]int, []int) {

	vlens := make([]int, len(lists))
	totals := make([]int, len(lists))

	for i, list := range lists {
		vlens[i], totals[i] = list.vlenTotal()
	}
	return vlens, totals
}

func (list *List) updateOffsetToTable(lastIdx, vlen, total, header_extend, vSize int) int {
	toTable := func(idx int) int {
		return int(list.NodeList.ValueInfo.Pos) + idx*4
	}
	// +4 all offset to table in current list. add new one element
	for i := 0; i < vlen; i++ {
		off := flatbuffers.GetUint32(list.R(toTable(i), Size(4)))
		off += 4
		flatbuffers.WriteUint32(list.U(toTable(i), 4), off)
	}

	// position to store new Data
	toData := toTable(0) + total + header_extend
	flatbuffers.WriteUint32(list.U(toTable(lastIdx), 4), uint32(toData+vSize-toTable(lastIdx)))

	return toData
}

func (list *List) moveOffsetToTable(size int) {
	if CurrentGlobalConfig.useNewMovOfftToTable {
		list.movOffToTable(size)
		return
	}

	list.oldMovOffToTable(size)
}

func (list *List) oldMovOffToTable(size int) {

	cnt := list.Count()

	startToTable := int(list.NodeList.ValueInfo.Pos)
	for toTable := startToTable; toTable < startToTable+cnt*4; toTable += 4 {
		off := flatbuffers.GetUint32(list.R(toTable, Size(4)))
		off += uint32(size)
		flatbuffers.WriteUint32(list.U(toTable, 4), off)
	}
}

func (list *List) movOffToTable(size int) {

	cnt := list.Count()

	startToTable := int(list.NodeList.ValueInfo.Pos)
	sizeToTable := cnt * 4

	posToOff := map[int]uint32{}

	for toTable := startToTable; toTable < startToTable+sizeToTable; toTable += 4 {
		off := flatbuffers.GetUint32(list.R(toTable, Size(4)))
		off += uint32(size)
		posToOff[toTable-startToTable] = off
	}

	bufs := list.U(startToTable, sizeToTable)

	if len(bufs) != sizeToTable {
		panic(fmt.Sprintf("movOffToTable len=%d cap=%d sizeToTable=%d\n",
			len(bufs), cap(bufs), sizeToTable))
	}

	for pos, off := range posToOff {
		flatbuffers.WriteUint32(bufs[pos:], off)
		roff := flatbuffers.GetUint32(list.R(pos+startToTable, Size(4)))
		if off != roff && bytes.Equal(bufs[pos:pos+4], list.R(pos+startToTable, Size(4))[:4]) {
			panic(fmt.Sprintf("bufs=%+v data=%+v\n", bufs[pos:pos+4], list.R(pos+startToTable, Size(4))[0:4]))
		}
	}

	// must equal pointer
	if &list.R(startToTable)[0] != &bufs[0] {
		panic(fmt.Sprintf("movOffToTable %p %p\n",
			&list.R(startToTable)[0], &bufs[0]))
	}

}

func (list *List) addTableList(alists ...*List) error {

	if alists[0].Count() == 0 {
		return nil
	}

	if list.Count() == 0 {
		oalist := alists[0]
		alist := oalist.dup()
		list.IO = alist.IO.Dup()
		list.NodeList.ValueInfo.Pos = alist.NodeList.ValueInfo.Pos
		list.NodeList.ValueInfo = ValueInfo(list.InfoSlice())
		return nil
	}

	vlen, _ := list.vlenTotal()

	list.NodeList.ValueInfo = ValueInfo(list.InfoSlice())
	oldImpl := list.Impl()

	firstElm, e := list.At(list.indexOnMaxPos(false))
	lastElm, _ := list.At(list.indexOnMaxPos(true))
	_ = lastElm
	if e != nil {
		return fmt.Errorf("addTableList(): cannot found first element m=%s", e)
	}

	dataEnd := list.NodeList.ValueInfo.Pos + list.NodeList.ValueInfo.Size
	vtableStart := firstElm.Node.Pos - firstElm.VirtualTableLen()
	headerEnd := list.NodeList.ValueInfo.Pos + 4*list.Count()

	oalist := alists[0]
	alist := oalist.dup()

	alist.NodeList.ValueInfo = ValueInfo(alist.InfoSlice())

	alastElm, e := alist.At(alist.indexOnMaxPos(true))
	if e != nil {
		return fmt.Errorf("addTableList(): cannot found alast element m=%s", e)
	}
	afirstElm, e := alist.At(alist.indexOnMaxPos(false))
	if e != nil {
		return fmt.Errorf("addTableList(): cannot found first element m=%s", e)
	}

	aDataEnd := alastElm.Node.Pos + alist.NodeList.ValueInfo.Size
	aVlen := int(alist.NodeList.VLen)
	aSizeOfHeader := alist.Count() * 4

	avTableStart := afirstElm.Node.Pos - afirstElm.VirtualTableLen()

	o := CurrentLogLevel
	SetLogLevel(LOG_WARN)
	defer SetLogLevel(o)

	adumper := dump.New(alist,
		func(v interface{}) string {
			alist := v.(*List)
			a := alist.Impl().Dump(0, OptDumpSize(500))
			return a
		},
		[]dump.FuncInfo{
			{"alist", dump.TypeVariable},
			{"InsertSpace", dump.TypeVariable},
			{"(avTableStart, dataEnd-vtableStart, false)", dump.TypeParam},
		})
	dumper := dump.New(list,
		func(v interface{}) string {
			alist := v.(*List)
			return alist.Impl().Dump(0, OptDumpSize(500))
		},
		[]dump.FuncInfo{
			{"alist", dump.TypeVariable},
			{"InsertSpace", dump.TypeVariable},
			{"(avTableStart, dataEnd-vtableStart, false)", dump.TypeParam},
		})

	// make space for list data (first element vtable -> list data last
	adumper.DumpWithFlag(L2isEnable(L2_DEBUG_IS), "B: make space for list data")
	alist.toCommonNode().InsertSpace(avTableStart, dataEnd-vtableStart, false)
	adumper.DumpWithFlag(L2isEnable(L2_DEBUG_IS), "A: make space for list data InsertSpace(0x%x,0x%x,%v)", avTableStart, dataEnd-vtableStart, false)

	if !alist.toCommonNode().InRoot() {
		alist.moveOffsetToTable(dataEnd - vtableStart)
		adumper.DumpWithFlag(L2isEnable(L2_DEBUG_IS), "A: move inc offset to Talbe 0x%x", dataEnd-vtableStart)
	}

	dumper.DumpWithFlag(L2isEnable(L2_DEBUG_IS), "B: make space for headers to insert alist")
	list.toCommonNode().InsertSpace(headerEnd, aSizeOfHeader, false)
	if !list.toCommonNode().InRoot() {
		list.moveOffsetToTable(aSizeOfHeader)
	}
	dumper.DumpWithFlag(L2isEnable(L2_DEBUG_IS), "A: make space for headers to insert alist moveOffsetToTable(0x%x,0x%x,%v)", headerEnd, aSizeOfHeader, false)

	ptrIdx := func(idx int) int {
		return int(list.NodeList.ValueInfo.Pos) + idx*4
	}

	dumper.DumpWithFlag(L2isEnable(L2_DEBUG_IS), "B: merge list and write vector len")
	adumper.DumpWithFlag(L2isEnable(L2_DEBUG_IS), "B: merge list and write vector len")

	flatbuffers.WriteUint32(list.U(ptrIdx(-1), 4), uint32(vlen+aVlen))
	list.Copy(alist.IO, alist.NodeList.ValueInfo.Pos, aSizeOfHeader, headerEnd, 0)

	dumper.DumpWithFlag(L2isEnable(L2_DEBUG_IS), "A: merge list and write vector len Copy(0x%x,0x%x,0x%x,0x%x)", alist.NodeList.ValueInfo.Pos, aSizeOfHeader, headerEnd, 0)

	list.Copy(alist.IO, avTableStart+dataEnd-vtableStart, aDataEnd-avTableStart, dataEnd+aSizeOfHeader, 0)
	dumper.DumpWithFlag(L2isEnable(L2_DEBUG_IS), "A: merge list data Copy(0x%x,0x%x,0x%x,0x%x)", avTableStart+dataEnd-vtableStart, aDataEnd-avTableStart, dataEnd+aSizeOfHeader, 0)

	Log2(L2_DEBUG_IS, L2fmt("alist dump history\n %s\n", adumper.String()))
	Log2(L2_DEBUG_IS, L2fmt("list dump history\n %s\n", dumper.String()))

	if list.IO.Type() == BASE_NO_LAYER || list.IO.Type() == BASE_DOUBLE_LAYER {
		oldImpl.overwrite(list.Impl())
	}

	// adumper.Finish()
	// dumper.Finish()

	list.NodeList.ValueInfo = ValueInfo(list.InfoSlice())

	return nil
}

func (node *List) setStructAt(idx int, elm *CommonNode) error {

	if idx > int(node.NodeList.ValueInfo.VLen) {
		return ERR_INVALID_INDEX
	}

	vlen, total := node.vlenTotal()

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
		elm.IO, elm.Node.Pos, elm.Node.Size,
		ptr, extend)

	total += extend
	if vlen == idx {
		vlen++
		flatbuffers.WriteUint32(node.U(node.NodeList.ValueInfo.Pos-4, 4), uint32(vlen))
	}
	node.NodeList.ValueInfo = ValueInfo(node.InfoSlice())
	return nil

}

func (node *List) setTableAt(idx int, elm *CommonNode) error {

	if idx > int(node.NodeList.ValueInfo.VLen) {
		return ERR_INVALID_INDEX
	}

	vlen, total := node.vlenTotal()

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
	toData := 0
	oldImpl := node.Impl()

	if header_extend > 0 {
		flatbuffers.WriteUint32(header, uint32(total-vlen*4+vSize))
		// extend space of setting offset to elm's table
		(*CommonNode)(node).InsertBuf(ptr, 4)
		toData = node.updateOffsetToTable(idx, vlen, total, header_extend, vSize)
		vlen++
	} else {
		toData = oElm.Node.Pos - vSize
	}

	t := node.IO.Type()
	_ = t

	if body_extend > 0 {
		// store vtable for new element
		(*CommonNode)(node).InsertSpace(toData, body_extend, false)
		if header_extend == 0 {
			for i := idx + 1; i < vlen; i++ {
				off := flatbuffers.GetUint32(node.R(ptrIdx(i), Size(4)))
				off += uint32(body_extend)
				flatbuffers.WriteUint32(node.U(ptrIdx(i), 4), off)
			}
		}
	}
	// update vlen
	flatbuffers.WriteUint32(node.U(ptrIdx(-1), 4), uint32(vlen))
	node.Copy(elm.IO,
		elm.Node.Pos-vSize, elm.Node.Size+vSize,
		toData, 0)
	node.NodeList.ValueInfo = ValueInfo(node.InfoSlice())
	t = node.IO.Type()

	if body_extend == 0 {
		return nil
	}

	if node.IO.Type() == BASE_NO_LAYER || node.IO.Type() == BASE_DOUBLE_LAYER {
		oldImpl.overwrite(node.Impl())
	}

	return nil

}

// SetAt ... Set Element to list.
//           if elm type is  variable length (examply Table). this operation is heavy.
//           to add list , should use Add()
func (node *List) SetAt(idx int, elm *CommonNode) error {

	if idx > int(node.NodeList.ValueInfo.VLen) {
		return ERR_INVALID_INDEX
	}

	g := GetTypeGroup(elm.Name)
	//ptr := int(node.NodeList.ValueInfo.Pos) + vlen*4

	if IsFieldBasicType(g) || IsFieldStruct(g) {
		return node.setStructAt(idx, elm)

	}
	if IsFieldUnion(g) {
		return ERR_NO_SUPPORT
	}
	if IsFieldTable(g) {
		return node.setTableAt(idx, elm)
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

	if i == j {
		return nil
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

	iOffset := int(flatbuffers.GetUint32(node.R(iPtr, Size(4))))
	jOffset := int(flatbuffers.GetUint32(node.R(jPtr, Size(4))))
	iOffset -= (jPtr - iPtr)
	jOffset -= (iPtr - jPtr)

	flatbuffers.WriteUint32(node.U(iPtr, size), uint32(jOffset))
	flatbuffers.WriteUint32(node.U(jPtr, size), uint32(iOffset))

	return nil

}

func (node *List) Len() int           { return int(node.VLen()) }
func (node *List) Swap(i, j int)      { node.SwapAt(i, j) }
func (node *List) Less(i, j int) bool { return node.lessFn(i, j) }

// SortBy ... sort of List
func (node *List) SortBy(less func(i, j int) bool) error {
	if !node.IsList() {
		return log.ERR_NO_SUPPORT
	}

	var o func(i, j int) bool
	node.lessFn, o = less, node.lessFn
	sort.Sort(node)
	node.lessFn = o

	return nil
}

// IsSorted ... return true if list is sorted
func (node *List) IsSorted(less func(i, j int) bool) (result bool) {
	result = false
	if !node.IsList() {
		return
	}

	var o func(i, j int) bool
	node.lessFn, o = less, node.lessFn
	result = sort.IsSorted(node)
	node.lessFn = o
	return
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

func (node *List) toTable(i int) int {

	return int(node.NodeList.ValueInfo.Pos) + i*4

}

func (node *List) new() (nList *List) {

	nList = &List{}
	nList.NodeList = &NodeList{}
	nList.Name = node.Name

	nList.InitList()
	switch node.IO.Type() {
	case BASE_IMPL:
		nList.IO = nList.IO.Impl()
	case BASE_NO_LAYER:
		nList.IO = NewNoLayer(nList.IO)
	case BASE_DOUBLE_LAYER:
		nList.IO = NewDoubleLayer(nList.IO)
	}
	return nList
}

// New ... return new list
//              if ListOptionge was not set, return empty list
//              if OptRange was set, return sub list by range
func (node *List) New(optFns ...ListOpt) (sub *List) {

	if len(optFns) == 0 {
		return node.new()
	}
	opt := DefaultListOpt()
	idx, cnt, e := node.setupOptRange(&opt, optFns...)
	if e != nil {
		return nil
	}

	if node.Count() <= idx {
		return nil
	}

	sub = node.new()

	sizeToTable := cnt * 4

	dumper := dump.New(sub,
		func(v interface{}) string {
			sub := v.(*List)
			return sub.Impl().Dump(0, OptDumpSize(500))
		},
		[]dump.FuncInfo{
			{"sub", dump.TypeVariable},
			{"SubList", dump.TypeVariable},
		})

	odumper := dump.New(sub,
		func(v interface{}) string {
			//sub := v.(*List)
			return node.Impl().Dump(0, OptDumpSize(500))
		},
		[]dump.FuncInfo{
			{"sub", dump.TypeVariable},
			{"SubList", dump.TypeVariable},
		})
	odumper.Dump("list ")
	Log2(L2OptFlag(LOG_DEBUG, FLT_NORMAL), L2fmt("list dump \n %s\n", odumper.String()))

	//copy header
	bufs := sub.U(sub.toTable(0), 4*cnt)

	posToOff := map[int]uint32{}

	for toTable := node.toTable(idx); toTable < node.toTable(idx)+sizeToTable; toTable += 4 {
		off := flatbuffers.GetUint32(node.R(toTable, Size(4)))
		posToOff[toTable-node.toTable(idx)] = off
	}
	Log2(L2OptFlag(LOG_DEBUG, FLT_NORMAL),
		L2fmt("copy header 0x%x size=0x%x -> 0x%x\n",
			node.toTable(idx), sizeToTable, sub.toTable(0)))

	dStart := -1
	dEnd := -1
	// find data range
	isTable := true
	for i := idx; i < idx+cnt; i++ {
		elm, begin, e := node.atWithBegin(i)
		if i == idx && isTable && elm.Node.Pos == begin {
			isTable = false
		}
		if e != nil {
			Log2(L2OptFlag(LOG_WARN, FLT_NORMAL), L2fmt("list.At(%d) err=%s", i, e.Error()))
			continue
		}
		if i == idx {
			dStart = begin
		}

		if dStart > begin {
			dStart = begin
		}
		if dEnd < elm.Node.Pos+elm.Info().Size {
			dEnd = elm.Node.Pos + elm.Info().Size
		}
	}

	// skip space in vector header
	skipSize := (node.Count() - idx - cnt) * 4

	// skip data space before idx
	skipSize += dStart - node.toTable(node.Count())

	if !isTable {
		goto NO_TABLE
	}

	Log2(L2OptFlag(LOG_DEBUG, FLT_NORMAL),
		L2fmt("sub.Copy(0x%x, 0x%x, 0x%x, 0)\n",
			dStart, dEnd-dStart,
			sub.toTable(0)+sizeToTable),
	)
	Log2(L2OptFlag(LOG_DEBUG, FLT_NORMAL), L2fmt("skipsize =0x%x\n", skipSize))
	sub.Copy(node.IO.Dup(), dStart, dEnd-dStart, sub.toTable(0)+sizeToTable, 0)

	for pos, off := range posToOff {
		flatbuffers.WriteUint32(bufs[pos:], off-uint32(skipSize))
	}

	goto FINISH
NO_TABLE:

	sub.Copy(node.IO.Dup(), dStart, dEnd-dStart, sub.toTable(0), 0)

FINISH:

	flatbuffers.WriteUint32(sub.R(sub.Node.Pos-4, Size(4)), uint32(cnt))

	dumper.Dump("A: SubList()")

	Log2(L2OptFlag(LOG_DEBUG, FLT_NORMAL), L2fmt("sub dump history\n %s\n", dumper.String()))
	sub.NodeList.ValueInfo = ValueInfo(sub.InfoSlice())
	return sub

}

func (node *List) indexOnMaxPos(t bool) (idx int) {

	pos := 0
	idx = -1
	for i, elm := range node.All() {
		if i == 0 {
			if t {
				pos = elm.Node.Pos + elm.Info().Size
			} else {
				pos = elm.Node.Pos
			}
			idx = i
			continue
		}

		if t {
			if epos := elm.Node.Pos + elm.Info().Size; epos > pos {
				pos = epos
				idx = i
			}
		} else {
			if epos := elm.Node.Pos; epos < pos {
				pos = epos
				idx = i
			}
		}
	}
	return
}

func (node *List) atWithBegin(i int) (*CommonNode, int, error) {

	elm, e := node.At(i)
	if e != nil {
		return elm, -1, e
	}

	grp := node.typeGroupOfChild()
	if IsFieldBasicType(grp) || IsFieldStruct(grp) {
		return elm, elm.Node.Pos, nil
	}

	return elm, elm.Node.Pos - elm.VirtualTableLen(), nil

}

func (node *List) typeGroupOfChild() int {

	tName := node.Name[2:]
	return GetTypeGroup(tName)

}

type ListOption struct {
	start, last int
}

type ListOpt func(*ListOption)

type ListOpts []ListOpt

func (opts ListOpts) Apply(opt *ListOption) {

	for _, o := range opts {
		o(opt)
	}

}
func DefaultListOpt() ListOption {

	return ListOption{start: -1, last: -1}

}

func OptRange(start, last int) ListOpt {

	return func(s *ListOption) {
		s.start = start
		s.last = last
	}
}

func (o *ListOption) hasError() error {

	if o.last < o.start {
		return errors.New("out of range")
	}

	if o.start < 0 {
		return errors.New("out of range")
	}

	return nil

}

func (o *ListOption) toParam() (idx, cnt int) {

	return o.start, o.last - o.start + 1
}

func (node *List) setupOpt(opt *ListOption, optFns ...ListOpt) error {

	ListOpts(optFns).Apply(opt)
	if opt.hasError() != nil {
		return opt.hasError()
	}
	return nil
}

func (node *List) setupOptRange(opt *ListOption, optFns ...ListOpt) (int, int, error) {

	if e := node.setupOpt(opt, optFns...); e != nil {
		return -1, -1, e
	}
	idx, cnt := opt.toParam()

	return idx, cnt, nil

}
