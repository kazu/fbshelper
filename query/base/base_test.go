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
	query "github.com/kazu/fbshelper/example/query2"
	query2 "github.com/kazu/fbshelper/example/query2"
	"github.com/kazu/fbshelper/example/vfs_schema"
	"github.com/kazu/fbshelper/query/base"
	log "github.com/kazu/fbshelper/query/log"
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

	// rBufInfo := root.BufInfo()
	// r2BufInfo := root2.BufInfo()
	// _, _ = rBufInfo, r2BufInfo
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
	root, _ = first.Root()
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

	base.SetLogLevel(base.LOG_DEBUG)
	buf := MakeRootFileFbsNoVersion(12, "root_test1.json", 456)

	root := query2.OpenByBuf(buf)

	common := query2.FromUint64(13)

	root.Record().SetFieldAt(0, common)
	assert.Equal(t, uint64(13), root.Record().FileId().Uint64(), "edit Root.Record.FileId")

	e := root.SetFieldAt(2, common)
	assert.Error(t, e)

	version := query2.FromInt32(1)
	oVersion := root.Version().Int32()
	root.SetFieldAt(0, version)

	assert.Equal(t, int32(0), oVersion)
	assert.Equal(t, int32(1), root.Version().Int32(), "edit Root.Version")

	oRoot := query2.OpenByBuf(buf)
	nFile := query2.NewFile()

	fbsUint64 := query2.FromUint64(13)
	fbsInt64 := query2.FromInt64(55)

	nFile.SetFieldAt(0, fbsUint64)
	nFile.SetFieldAt(2, fbsInt64)
	nFile.SetFieldAt(1, base.FromBytes([]byte("file-bytes")))

	name := nFile.Name()
	name.SetAt(0, query.FromByte([]byte("Z")[0]))
	name.Base.Merge()
	nFile.SetName(name)

	assert.Equal(t, root.Index().File().Id().Uint64(), oRoot.Index().File().Id().Uint64())
	assert.Equal(t, uint64(13), nFile.Id().Uint64())
	assert.Equal(t, []byte("Zile-bytes"), nFile.Name().Bytes(), "edit File.Name")

	root.SetFieldAt(2, nFile.CommonNode)

	assert.NotEqual(t, oRoot.Index().File().Id().Uint64(), root.Index().File().Id().Uint64())
	assert.Equal(t, fbsUint64.Uint64(), root.Index().File().Id().Uint64())
	assert.Equal(t, fbsInt64.Int64(), root.Index().File().IndexAt().Int64())
	assert.Equal(t, []byte("Zile-bytes"), root.Index().File().Name().Bytes(), "edit Root.Index.File.Name")

	buf = MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})

	root = query2.OpenByBuf(buf)

	cnt := root.Index().IndexString().Maps().Count()
	_ = cnt

	// insert new element to vector
	invs := query2.NewInvertedMapString()
	rec := query2.NewRecord()
	rec.SetFileId(query2.FromUint64(9))
	rec.SetOffset(query2.FromInt64(8))
	rec.SetSize(query2.FromInt64(7))
	rec.SetOffsetOfValue(query2.FromInt32(6))
	rec.SetValueSize(query2.FromInt32(5))
	invs.SetValue(rec.CommonNode)
	invs.SetKey(base.FromBytes([]byte("inverted")))

	maps := root.Index().IndexString().Maps()
	query2.RootFromCommon(maps.CommonNode).AllTree().DumpAll(0, os.Stdout)
	maps.SetAt(cnt, invs)

	// replace new element to vector
	invs = query2.NewInvertedMapString()
	rec = query2.NewRecord()
	rec.SetFileId(query2.FromUint64(19))
	rec.SetOffset(query2.FromInt64(18))
	rec.SetSize(query2.FromInt64(17))
	rec.SetOffsetOfValue(query2.FromInt32(16))
	rec.SetValueSize(query2.FromInt32(15))
	invs.SetValue(rec.CommonNode)
	invs.SetKey(base.FromBytes([]byte("inverted2")))

	maps.SetAt(1, invs)

	root2 := query2.RootFromCommon(maps.CommonNode)
	root2.Dedup()

	assert.Equal(t, cnt+1, root2.Index().IndexString().Maps().Count(), "root.IndexString.Maps.Count")
	assert.Equal(t, cnt+1, 3)
	assert.Equal(t, []byte("inverted"), query2.InvertedMapStringSingle(root2.Index().IndexString().Maps().At(2)).Key().Bytes())
	assert.Equal(t, []byte("inverted2"), query2.InvertedMapStringSingle(root2.Index().IndexString().Maps().At(1)).Key().Bytes())

	query2.RootFromCommon(maps.CommonNode).AllTree().DumpAll(0, os.Stdout)

}

func Test_NewRootIndexString(t *testing.T) {
	log.SetLogLevel(log.LOG_DEBUG)

	var e error
	root := query2.NewRoot()
	root.WithHeader()

	root.SetVersion(query2.FromInt32(3))
	root.SetIndexType(query2.FromByte(2))

	idxStr := query2.NewIndexString()
	idxStr.SetSize(query2.FromInt32(3))

	assert.Equal(t, query2.FromInt32(3).Int32(), idxStr.Size().Int32())

	inv := query2.NewInvertedMapString()
	inv.SetKey(base.FromBytes([]byte("inv-str-key")))

	assert.Equal(t, []byte("inv-str-key"), inv.Key().Bytes())

	rec := query2.NewRecord()
	rec.SetFileId(query2.FromUint64(9))
	rec.SetOffset(query2.FromInt64(8))
	rec.SetSize(query2.FromInt64(7))
	rec.SetOffsetOfValue(query2.FromInt32(6))
	rec.SetValueSize(query2.FromInt32(5))

	inv.SetValue(rec.CommonNode)

	assert.Equal(t, []byte("inv-str-key"), inv.Key().Bytes())
	assert.Equal(t, query2.FromUint64(9).Uint64(), inv.Value().FileId().Uint64())

	maps := query2.NewInvertedMapStringList()
	maps.SetAt(0, inv)

	inv, e = maps.At(0)

	assert.NoError(t, e)
	assert.Equal(t, []byte("inv-str-key"), inv.Key().Bytes())
	assert.Equal(t, query2.FromUint64(9).Uint64(), inv.Value().FileId().Uint64())

	idxStr.SetMaps(maps.CommonNode)

	inv, e = idxStr.Maps().At(0)

	assert.NoError(t, e)
	assert.Equal(t, query2.FromInt32(3).Int32(), idxStr.Size().Int32())
	assert.Equal(t, []byte("inv-str-key"), inv.Key().Bytes())
	assert.Equal(t, query2.FromUint64(9).Uint64(), inv.Value().FileId().Uint64())

	root.SetIndex(idxStr.CommonNode)

	idx := root.Index().IndexString()
	idxStr = &idx
	inv, e = idxStr.Maps().At(0)

	assert.NoError(t, e)
	assert.Equal(t, query2.FromInt32(3).Int32(), idxStr.Size().Int32())
	assert.Equal(t, []byte("inv-str-key"), inv.Key().Bytes())
	assert.Equal(t, query2.FromUint64(9).Uint64(), inv.Value().FileId().Uint64())

}

type MyWriter struct {
	Buf []byte
}

func (w *MyWriter) Write(p []byte) (n int, e error) {

	w.Buf = append(w.Buf, p...)
	return len(p), nil
}

func (w *MyWriter) WriteAt(p []byte, offset int64) (n int, e error) {
	oldlen := len(w.Buf)
	if len(w.Buf) < int(offset)+len(p) {
		w.Buf = w.Buf[:int(offset)+len(p)]
	}
	off := int(offset)
	//w.Buf[off : off+len(p)] = p
	copy(w.Buf[off:off+len(p)], p)

	if oldlen < int(offset)+len(p) && len(w.Buf) > int(offset)+len(p) {
		fmt.Println("???")
	}

	return len(p), nil
}

func Test_CommonList(t *testing.T) {

	//files := query2.NewFiles()
	flist := query2.NewFileList()

	for i := 0; i < 3; i++ {
		file := query2.NewFile()
		file.SetId(query2.FromUint64(10 + uint64(i)))
		file.SetIndexAt(query2.FromInt64(2000 + int64(i)))
		file.SetName(base.FromBytes([]byte("namedayo")))
		flist.SetAt(i, file)
	}

	flist.Merge()
	w := &MyWriter{}
	w.Buf = make([]byte, 0, 4096)

	cl := base.CommonList{}
	cl.CommonNode = flist.CommonNode

	cl.SetDataWriter(w)
	cl.WriteDataAll()
	cnt := cl.Count()

	for i := 0; i < 3; i++ {
		file := query2.NewFile()
		file.SetId(query2.FromUint64(20 + uint64(i)))
		file.SetIndexAt(query2.FromInt64(3000 + int64(i)))
		file.SetName(base.FromBytes([]byte("namedayoadd")))
		cl.SetAt(cnt+i, file.CommonNode)
	}

	cl.Merge()
	flist = query2.NewFileList()

	flist.NodeList.ValueInfo = cl.NodeList.ValueInfo
	newBytes := cl.Base.R(0)[:cl.LenBuf()]
	newBytes = append(newBytes, w.Buf...)
	flist.Base = base.NewBase(newBytes)
	flist.Node.Pos = cl.Node.Pos

	assert.Equal(t, 6, flist.Count())
	file, _ := flist.At(3)

	assert.True(t, len(w.Buf) > 10)
	assert.Equal(t, uint64(20), file.Id().Uint64())
	assert.Equal(t, []byte("namedayoadd"), file.Name().Bytes())

}

type Lazy func() interface{}

func Test_LazyLog(t *testing.T) {

	b := 0

	a := Lazy{
		b += 1
		return 1 + 1
	}
	Log(false, "aaaa %d ", a)

	assert.Equal(t, 0, b)
}

//type Infn func() interface{}
func Log(run bool, s string, fns ...(func() interface{})) {
	if run {
		args := []interface{}{}
		for _, fn := range fns {
			args = append(args, fn())
		}
		fmt.Printf(s, args...)
	}
}
