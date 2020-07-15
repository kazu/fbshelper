package base_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"unsafe"

	flatbuffers "github.com/google/flatbuffers/go"
	query2 "github.com/kazu/fbshelper/example/query2"
	"github.com/kazu/fbshelper/example/vfs_schema"
	"github.com/kazu/fbshelper/query/base"
	"github.com/stretchr/testify/assert"
)

func MakeRootFileFbs(id uint64, name string, index_at int64) []byte {

	b := flatbuffers.NewBuilder(130)
	fname := b.CreateString(name)
	zz := b.Bytes[120:]
	_ = zz

	vfs_schema.FileStart(b)
	vfs_schema.FileAddId(b, id)
	vfs_schema.FileAddName(b, fname)
	vfs_schema.FileAddIndexAt(b, index_at)
	fbFile := vfs_schema.FileEnd(b)

	vfs_schema.RootStart(b)
	vfs_schema.RootAddVersion(b, 1)
	vfs_schema.RootAddIndexType(b, vfs_schema.IndexFile)
	vfs_schema.RootAddIndex(b, fbFile)
	vfs_schema.RootAddRecord(b, vfs_schema.CreateRecord(b, id, 12, 34, 56, 78))

	b.Finish(vfs_schema.RootEnd(b))
	return b.FinishedBytes()
}

func MakeRootRecordFbs(id uint64, name string, index_at int64) []byte {

	b := flatbuffers.NewBuilder(0)

	vfs_schema.RootStart(b)
	vfs_schema.RootAddVersion(b, 0)
	vfs_schema.RootAddIndexType(b, vfs_schema.IndexFile)
	vfs_schema.RootAddRecord(b, vfs_schema.CreateRecord(b, id, 12, 34, 56, 78))
	b.Finish(vfs_schema.RootEnd(b))
	return b.FinishedBytes()
}

func MakeRootFileFbsNoVersion(id uint64, name string, index_at int64) []byte {

	b := flatbuffers.NewBuilder(130)
	fname := b.CreateString(name)
	zz := b.Bytes[120:]
	_ = zz

	vfs_schema.FileStart(b)
	vfs_schema.FileAddId(b, id)
	vfs_schema.FileAddName(b, fname)
	vfs_schema.FileAddIndexAt(b, index_at)
	fbFile := vfs_schema.FileEnd(b)

	vfs_schema.RootStart(b)
	//vfs_schema.RootAddVersion(b, 1)
	vfs_schema.RootAddIndexType(b, vfs_schema.IndexFile)
	vfs_schema.RootAddIndex(b, fbFile)
	vfs_schema.RootAddRecord(b, vfs_schema.CreateRecord(b, id, 12, 34, 56, 78))

	b.Finish(vfs_schema.RootEnd(b))
	return b.FinishedBytes()
}

func MakeRootNumList(fn func(b *flatbuffers.Builder) flatbuffers.UOffsetT) []byte {
	b := flatbuffers.NewBuilder(0)
	numList := fn(b)

	vfs_schema.RootStart(b)
	vfs_schema.RootAddVersion(b, 1)
	vfs_schema.RootAddIndexType(b, vfs_schema.IndexNumList)
	vfs_schema.RootAddIndex(b, numList)
	vfs_schema.RootAddRecord(b, vfs_schema.CreateRecord(b, 666, 12, 34, 56, 78))

	b.Finish(vfs_schema.RootEnd(b))
	return b.FinishedBytes()
}

func MakeNumList(b *flatbuffers.Builder, fn func(i int) int32) flatbuffers.UOffsetT {

	vfs_schema.NumListStartNumVector(b, 2)
	// vector store reverse
	b.PrependInt32(fn(1))
	b.PrependInt32(fn(0))
	nums := b.EndVector(2)

	vfs_schema.NumListStart(b)
	vfs_schema.NumListAddNum(b, nums)
	return vfs_schema.IndexStringEnd(b)

}

func MakeRootIndexString(fn func(b *flatbuffers.Builder) flatbuffers.UOffsetT) []byte {
	b := flatbuffers.NewBuilder(0)

	fbIndex := fn(b)

	vfs_schema.RootStart(b)
	vfs_schema.RootAddVersion(b, 1)
	vfs_schema.RootAddIndexType(b, vfs_schema.IndexIndexString)
	vfs_schema.RootAddIndex(b, fbIndex)
	vfs_schema.RootAddRecord(b, vfs_schema.CreateRecord(b, 666, 12, 34, 56, 78))

	b.Finish(vfs_schema.RootEnd(b))
	return b.FinishedBytes()

}

func MakeIndexString(b *flatbuffers.Builder, fn func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT) flatbuffers.UOffsetT {

	data := fn(b, 0)
	data2 := fn(b, 1)

	vfs_schema.IndexStringStartMapsVector(b, 2)
	b.PrependUOffsetT(data)
	b.PrependUOffsetT(data2)
	maps := b.EndVector(2)

	vfs_schema.IndexStringStart(b)
	vfs_schema.IndexStringAddMaps(b, maps)
	vfs_schema.IndexStringAddSize(b, 234)
	return vfs_schema.IndexStringEnd(b)

}

func MakeInvertedMapString(b *flatbuffers.Builder, key string) flatbuffers.UOffsetT {

	fkey := b.CreateString(key)
	vfs_schema.InvertedMapStringStart(b)
	vfs_schema.InvertedMapStringAddKey(b, fkey)
	vfs_schema.InvertedMapStringAddValue(b, vfs_schema.CreateRecord(b, 1, 2, 3, 4, 5))
	return vfs_schema.InvertedMapStringEnd(b)
}

func MakeRootRecord(key uint64) []byte {

	b := flatbuffers.NewBuilder(0)

	vfs_schema.InvertedMapNumStart(b)
	vfs_schema.InvertedMapNumAddKey(b, int64(key))
	vfs_schema.InvertedMapNumAddValue(b, vfs_schema.CreateRecord(b, 1, 2, 3, 0, 0))
	iMapNum := vfs_schema.InvertedMapNumEnd(b)

	vfs_schema.RootStart(b)
	vfs_schema.RootAddVersion(b, 1)
	vfs_schema.RootAddIndexType(b, vfs_schema.IndexInvertedMapNum)
	vfs_schema.RootAddIndex(b, iMapNum)
	vfs_schema.RootAddRecord(b, vfs_schema.CreateRecord(b, 666, 12, 34, 56, 78))

	b.Finish(vfs_schema.RootEnd(b))

	return b.FinishedBytes()

}

type File struct {
	ID      uint64 `fbs:"Id"`
	Name    []byte `fbs:"Name"`
	IndexAt int64  `fbs:"IndexAt"`
}

func TestBase(t *testing.T) {

	buf := MakeRootFileFbs(12, "root_test1.json", 456)
	assert.NotNil(t, buf)

	q := query2.OpenByBuf(buf)

	assert.Equal(t, int32(1), q.Version().Int32())
	assert.Equal(t, uint64(12), q.Index().File().Id().Uint64())
	assert.Equal(t, []byte("root_test1.json"), q.Index().File().Name().Bytes())
	assert.Equal(t, int64(456), q.Index().File().IndexAt().Int64())
}

type FileTest struct {
	TestName string
	ID       uint64
	Name     []byte
	IndexAt  int64
}

func DataRootFileTest() []FileTest {

	return []FileTest{
		{"first", 14, []byte("root_test1.json"), 755},
		{"second", 12, []byte("root_test6.json"), 238},
		{"third", 68789, []byte("root_test6 .json"), 6789},
	}

}

func TestUnmarshal(t *testing.T) {

	tests := DataRootFileTest()
	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			buf := MakeRootFileFbs(tt.ID, string(tt.Name), tt.IndexAt)
			file := File{}
			fq := query2.OpenByBuf(buf).Index().File()
			e := fq.Unmarshal(&file)

			assert.NoError(t, e)
			assert.Equal(t, tt.ID, file.ID)
			assert.Equal(t, tt.Name, file.Name)
			assert.Equal(t, tt.IndexAt, file.IndexAt)
		})
	}
}

func TestOpen(t *testing.T) {
	tests := DataRootFileTest()
	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			buf := MakeRootFileFbs(tt.ID, string(tt.Name), tt.IndexAt)
			file := File{}
			fq := query2.Open(bytes.NewReader(buf), 512).Index().File()
			e := fq.Unmarshal(&file)

			assert.NoError(t, e)
			assert.Equal(t, tt.ID, file.ID)
			assert.Equal(t, tt.Name, file.Name)
			assert.Equal(t, tt.IndexAt, file.IndexAt)
		})
	}
}

func TestOpen2(t *testing.T) {
	tests := DataRootFileTest()
	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			buf := MakeRootFileFbs(tt.ID, string(tt.Name), tt.IndexAt)
			file := File{}
			root := query2.Open(bytes.NewReader(buf), 512)
			fq := root.Index().File()
			e := fq.Unmarshal(&file)

			assert.NoError(t, e)
			assert.Equal(t, tt.ID, fq.Id().Uint64())
			assert.Equal(t, tt.IndexAt, fq.IndexAt().Int64())
			assert.Equal(t, tt.ID, file.ID)
			assert.Equal(t, tt.Name, file.Name)
			assert.Equal(t, tt.IndexAt, file.IndexAt)
		})
	}
}

func Test_QueryFbs(t *testing.T) {
	buf := MakeRootFileFbs(12, "root_test.json", 456)
	root := query2.OpenByBuf(buf)
	idx := root.Index()
	assert.Equal(t, uint64(12), idx.File().Id().Uint64())
	assert.Equal(t, "root_test.json", string(idx.File().Name().Bytes()))

	buf2 := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})
	root2 := query2.OpenByBuf(buf2)
	assert.Equal(t, len(buf2), root2.Len())
	z, e := root2.Index().IndexString().Maps().Last()
	_, _ = z, e

	assert.Equal(t, uint64(1), z.Value().FileId().Uint64())
	assert.Equal(t, int64(2), z.Value().Offset().Int64())

	assert.Equal(t, int32(234), root2.Index().IndexString().Size().Int32())

	assert.Equal(t, 2, root2.Index().IndexString().Maps().Count())

}

func Test_QueryNext(t *testing.T) {
	buf := MakeRootRecord(512)
	buf2 := append(buf, MakeRootRecord(513)...)

	root := query2.Open(bytes.NewReader(buf2), base.DEFAULT_BUF_CAP)
	z := root.Index().InvertedMapNum()

	record := z.Value()
	rPos1 := record.Node.Pos

	val := z.FieldAt(0)
	assert.NotNil(t, val)

	assert.Equal(t, int64(512), root.Index().InvertedMapNum().Key().Int64())
	assert.Equal(t, int64(2), record.Offset().Int64())
	assert.Equal(t, rPos1, record.Node.Pos)
	assert.Equal(t, uint64(1), record.FileId().Uint64())
	assert.Equal(t, rPos1, record.Node.Pos)
	assert.Equal(t, len(buf), root.Len())
	assert.Equal(t, len(buf), root.Len())
	assert.Equal(t, len(buf), z.Info().Pos+z.Info().Size)

	root2 := root.Next()

	assert.Equal(t, int64(513), root2.Index().InvertedMapNum().Key().Int64())
	assert.False(t, root2.HasNext())
	len1 := root.LenBuf()

	rBufInfo := root.BufInfo()
	r2BufInfo := root2.BufInfo()
	_, _ = rBufInfo, r2BufInfo
	assert.Equal(t, len(buf), len1)

}

func Test_RootIndexStringInfoPos(t *testing.T) {

	buf := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})

	q := query2.Open(bytes.NewReader(buf), 512)

	list := q.Index().IndexString().Maps()
	n := list.Count()
	_ = n

	infos := map[string]base.Info{}

	lFirst, _ := list.First()
	info := lFirst.ValueInfo(0)
	infos["Maps[0].0"] = base.Info(info)

	for i := 0; i < lFirst.CountOfField(); i++ {
		tmpInfo := lFirst.ValueInfo(i)
		infos[fmt.Sprintf("Maps[0].0.%d", i)] = base.Info(tmpInfo)

	}
	lLast, _ := list.Last()
	info2 := lLast.Value().Info()
	infos["Maps[1].1"] = base.Info(info2)

	assert.Equal(t, true, info.Pos < info2.Pos)

}

func Test_TraverseInfo(t *testing.T) {

	buf := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})

	type Result struct {
		name   string
		pNode  *base.CommonNode
		idx    int
		parent int
		child  int
		size   int
	}

	results := []Result{}

	recFn := func(node *base.CommonNode, idx, parent, child, size int) {
		//result[node.Name] = []int{parent, child, idx, size}
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
	var pos int
	cond := func(parent, child, size int) bool {
		//return parent <= pos && pos <= child
		return true
	}
	q := query2.Open(bytes.NewReader(buf), 512)
	pos = q.Index().IndexString().Size().Node.Pos

	q.TraverseInfo(pos, recFn, cond)
	assert.True(t, len(results) > 0)
	for _, result := range results {
		fmt.Printf("%+v\n", result)
	}
}

func dump(tree *base.Tree) string {
	if tree.Parent != nil {
		return fmt.Sprintf("{Type:\t%s,\tPos:\t%d\tParentPos:\t%d}\n",
			tree.Node.Name,
			tree.Pos(),
			tree.Parent.Pos(),
		)
	}
	return fmt.Sprintf("{Type:\t%s,\tPos:\t%d}\n",
		tree.Node.Name,
		tree.Pos(),
	)
}
func dumpTrees(trees []*base.Tree) string {

	var b strings.Builder
	for _, tree := range trees {
		fmt.Fprint(&b, dump(tree))
	}
	return b.String()
}

func dumpAll(i int, tree *base.Tree, w io.Writer) {
	/*if i > 10 {
		return
	}*/
	for j := 0; j < i*2; j++ {
		io.WriteString(w, "\t")
	}
	fmt.Fprintf(w, "%s", dump(tree))
	for _, child := range tree.Childs {
		dumpAll(i+1, child, w)
	}
}

func Test_AllTree(t *testing.T) {
	buf := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})

	type Tree = base.Tree

	q := query2.Open(bytes.NewReader(buf), 512)
	tree := q.AllTree()

	//var b strings.Builder

	dumpAll(0, tree, os.Stdout)
	assert.Equal(t, 4, len(tree.Childs))

}

func Test_FindTree(t *testing.T) {
	type Tree = base.Tree

	checkFn := func(pos int, tree *Tree) bool {
		return tree.Parent != nil && tree.Parent.Pos() <= pos && pos <= tree.Pos()
	}

	tests := []struct {
		Pos         int
		ResultLen   int
		ResultCheck func(int, *Tree) bool
	}{
		{31, 7, checkFn},
		{40, 6, checkFn},
		{44, 5, checkFn},
		{48, 5, checkFn},
		{50, 4, checkFn},
		{60, 3, checkFn},
		{70, 1, checkFn},
		{76, 3, checkFn},
		{80, 2, checkFn},
		{88, 1, checkFn},
		{126, 5, checkFn},
	}

	buf := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test pos=%d", tt.Pos), func(t *testing.T) {
			q := query2.Open(bytes.NewReader(buf), 512)
			results := make([]*Tree, 0)
			resultCh := q.FindTree(func(t *Tree) bool {
				return t.Parent != nil && t.Parent.Pos() <= tt.Pos && tt.Pos <= t.Pos()
			})
			for {
				t, ok := <-resultCh
				if !ok {
					break
				}
				results = append(results, t)
				// debug only
				//fmt.Print(dump(t))
			}
			assert.Equal(t, tt.ResultLen, len(results), fmt.Sprintf("pos=%d results=%s", tt.Pos, dumpTrees(results)))
			for _, tree := range results {
				assert.True(t, tt.ResultCheck(tt.Pos, tree))
			}
		})
	}

}

func Test_DirectSturct(t *testing.T) {

	type FileTest struct {
		FileId        uint64
		Offset        int64
		Size          int64
		OffsetOfValue int32
		ValueSize     int32
	}
	a := FileTest{
		FileId:        1,
		Offset:        2,
		Size:          3,
		OffsetOfValue: 4,
		ValueSize:     5,
	}

	buf := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})
	root := query2.OpenByBuf(buf)
	pos := query2.InvertedMapStringSingle(root.Index().IndexString().Maps().Last()).Value().Node.Pos
	_ = buf
	//fmt.Printf("top=0x%p fileId=0x%p ValueSize=0x%p \n", &a, &(a.FileId), &(a.ValueSize))
	b := (*FileTest)(unsafe.Pointer(&buf[pos]))
	elm, _ := root.Index().IndexString().Maps().Last()
	diff := int(unsafe.Offsetof(a.ValueSize)) //int(unsafe.Pointer(&a.ValueSize)) - int(unsafe.Pointer(&a.FileId))
	assert.Equal(t, elm.Value().ValueSize().Node.Pos-elm.Value().FileId().Node.Pos, diff)
	assert.Equal(t, elm.Value().FileId().Uint64(), b.FileId)
}

func Test_SliceBasicType(t *testing.T) {

	buf := MakeRootNumList(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeNumList(b, func(i int) int32 {
			switch i {
			case 0:
				return 345
			case 1:
				return 584
			default:
				return 999
			}
		})
	})

	q := query2.Open(bytes.NewReader(buf), 512)
	_ = q

	type CommonNodeWithErr struct {
		*base.CommonNode
		e error
	}
	withe := func(c *base.CommonNode, e error) CommonNodeWithErr {
		return CommonNodeWithErr{
			CommonNode: c,
			e:          e,
		}
	}

	assert.NotNil(t, buf)
	assert.Equal(t, 2, q.Index().NumList().Num().Count())
	assert.Equal(t, int32(345), withe(q.Index().NumList().Num().At(0)).Int32())
	assert.Equal(t, int32(584), withe(q.Index().NumList().Num().At(1)).Int32())

}

func Test_RootSetVersion(t *testing.T) {

	buf := MakeRootNumList(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeNumList(b, func(i int) int32 {
			switch i {
			case 0:
				return 345
			case 1:
				return 584
			default:
				return 999
			}
		})
	})

	q := query2.Open(bytes.NewReader(buf), 512)
	assert.Equal(t, int32(1), q.Version().Int32())
	q.Version().SetInt32(2)
	assert.Equal(t, int32(2), q.Version().Int32())
	assert.NotNil(t, buf)
}

func TestSetInt64(t *testing.T) {
	tests := DataRootFileTest()
	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			buf := MakeRootFileFbs(tt.ID, string(tt.Name), tt.IndexAt)
			file := File{}
			_ = file
			root := query2.Open(bytes.NewReader(buf), 512)
			fq := root.Index().File()
			assert.Equal(t, tt.ID, fq.Id().Uint64())
			assert.Equal(t, tt.IndexAt, fq.IndexAt().Int64())
			fq.IndexAt().SetInt64(tt.IndexAt + 2)
			assert.Equal(t, tt.IndexAt+2, fq.IndexAt().Int64())
			fq.Id().SetUint64(tt.ID + 2)
			assert.Equal(t, tt.ID+2, fq.Id().Uint64())
		})
	}
}

func Test_MakeRootFileNoVersion(t *testing.T) {

	buf := MakeRootFileFbsNoVersion(123, "aaa", 987)

	root := query2.OpenByBuf(buf)
	assert.Equal(t, int64(0), root.Version().Int64())
	assert.Equal(t, []byte("aaa"), root.Index().File().Name().Bytes())
	assert.Equal(t, int64(123), root.Index().File().Id().Int64())
}

func Test_InsertBuf(t *testing.T) {
	buf := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})

	//root := query2.OpenByBuf(buf)
	root := query2.Open(bytes.NewReader(buf), 512)
	root_old := root
	last := query2.InvertedMapStringSingle(root.Index().IndexString().Maps().Last())
	olastPos := last.Key().Node.Pos
	oRecordPos := last.Value().Node.Pos

	first := query2.InvertedMapStringSingle(root.Index().IndexString().Maps().First())
	oLen := root.Len()
	first.InsertBuf(126, 8)
	root = first.Root()
	last = query2.InvertedMapStringSingle(root.Index().IndexString().Maps().Last())

	assert.Equal(t, oLen, root.Len()-8)
	assert.Equal(t, olastPos, last.Key().Node.Pos-8)
	assert.Equal(t, oRecordPos, last.Value().Node.Pos-8)

	check := func(t *testing.T, args ...*query2.Root) {
		size := len(args)
		results := make([][]interface{}, 4)
		for i := range results {
			results[i] = make([]interface{}, size)
		}

		for i, root := range args {
			results[0][i] = root.Version().Int32()
			results[1][i] = root.Index().IndexString().Size().Int32()
			results[2][i] =
				query2.InvertedMapStringSingle(
					root.Index().IndexString().Maps().First(),
				).Value().FileId().Uint64()
			results[3][i] =
				query2.InvertedMapStringSingle(
					root.Index().IndexString().Maps().Last(),
				).Value().FileId().Uint64()
		}
		for i := range results {
			assert.Equal(t, results[i][0], results[i][1], i)
		}
	}

	check(t, &root_old, &root)
}

func Test_U(t *testing.T) {

	node := base.NewNode2(base.NewBase(make([]byte, 8)), 0, true)

	buf := node.U(0, 8)
	buf[0] = 1

	assert.Equal(t, byte(1), node.R(0)[0])
}

func Test_BaseCopy(t *testing.T) {

	obuf1 := []byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	obuf2 := []byte{5, 6, 7, 0, 5, 6, 7, 0, 5, 6, 7, 0, 5, 6, 7, 0}

	tests := []struct {
		SrcOff  int
		SrcSize int
		DstOff  int
		Extend  int
	}{
		{4, 6, 6, 6},
		{0, 16, 0, 16},  // add to front
		{0, 16, 16, 16}, // add to last
		{4, 4, 4, 0},    // overwrite
		{4, 4, 0, 0},    // overwrite front
		{4, 4, 12, 0},   // overwrite back
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Base.Copy(%v)", tt), func(t *testing.T) {
			buf1 := make([]byte, len(obuf1))
			buf2 := make([]byte, len(obuf2))
			copy(buf1, obuf1)
			copy(buf2, obuf2)

			node := base.NewNode2(base.NewBase(buf1), 0, true)
			node2 := base.NewNode2(base.NewBase(buf2), 0, true)

			node.Copy(node2.Base, tt.SrcOff, tt.SrcSize, tt.DstOff, tt.Extend)

			assert.Equal(t, node.R(tt.DstOff)[0], node2.R(tt.SrcOff)[0])
			assert.Equal(t, node.R(tt.DstOff + 3)[0], node2.R(tt.SrcOff + 3)[0],
				fmt.Sprintf("node.R(%d)=%v  node2.R(%d)=%v ", tt.DstOff+3, node.R(tt.DstOff+3), tt.SrcOff+5, node2.R(tt.SrcOff+5)))
			assert.Equal(t, len(obuf1)+tt.Extend, node.LenBuf())

		})

	}

}

func Test_SetFieldAt(t *testing.T) {

	buf := MakeRootFileFbsNoVersion(12, "root_test1.json", 456)

	root := query2.OpenByBuf(buf)

	common := query2.FromUint64(13)
	common.SetUint64(13)

	root.Record().SetFieldAt(0, common)
	assert.Equal(t, uint64(13), root.Record().FileId().Uint64(), "edit Root.Record.FileId")

	e := root.SetFieldAt(2, common)
	assert.Error(t, e)

	version := query2.FromInt32(1)
	version.SetInt32(1)
	oVersion := root.Version().Int32()
	root.SetFieldAt(0, version)

	assert.Equal(t, int32(0), oVersion)
	assert.Equal(t, int32(1), root.Version().Int32(), "edit Root.Version")

	oRoot := query2.OpenByBuf(buf)
	nFile := query2.NewFile()

	fbsUint64 := query2.FromUint64(13)
	fbsUint64.SetUint64(13)
	fbsInt64 := query2.FromInt64(55)
	fbsInt64.SetInt64(55)

	nFile.SetFieldAt(0, fbsUint64)
	nFile.SetFieldAt(2, fbsInt64)

	assert.Equal(t, root.Index().File().Id().Uint64(), oRoot.Index().File().Id().Uint64())

	root.SetFieldAt(2, nFile.CommonNode)

	assert.NotEqual(t, oRoot.Index().File().Id().Uint64(), root.Index().File().Id().Uint64())
	assert.Equal(t, fbsUint64.Uint64(), root.Index().File().Id().Uint64())
	assert.Equal(t, fbsInt64.Int64(), root.Index().File().IndexAt().Int64())

}

func Test_NewRootIndexString(t *testing.T) {
	//root := NewRoot()
	//root.Version().SetInt32(2)
	//root.SetVersion(FromInt32(2))
	//root.SetAt(0, FromInt32(2))/;

	// root = CreateRoot(2, nil)

	// idxStr := NewIndexString()
	// idxStr.Size().SetInt32(2)

	// idxStr := CreateIndexString(2, nil)

	// invStr := NewInvertedMapString()

	// rec1 := NewRecord()
	// rec1 = CreateRecord(1, 2, 3, 4, 5)

	// invStr.Value().Set(rec1)
	// root.Index().IndexString().Set(idxStr)
	// root.Index().IndexString().InvertedMapString().Maps().Add(invStr)

}
