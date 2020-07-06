package base_test

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/kazu/loncha"

	flatbuffers "github.com/google/flatbuffers/go"
	query "github.com/kazu/fbshelper/example/query"
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
	//, vFn func(b *flatbuffers.Builder) flatbuffers.UOffsetT) flatbuffers.UOffsetT{

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

	q := query.OpenByBuf(buf)

	assert.Equal(t, int32(1), q.Version())
	assert.Equal(t, uint64(12), q.Index().File().Id())
	assert.Equal(t, []byte("root_test1.json"), q.Index().File().Name())
	assert.Equal(t, int64(456), q.Index().File().IndexAt())
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
			fq := query.OpenByBuf(buf).Index().File()
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
			fq := query.Open(bytes.NewReader(buf), 512).Index().File()
			e := fq.Unmarshal(&file)

			assert.NoError(t, e)
			assert.Equal(t, tt.ID, file.ID)
			assert.Equal(t, tt.Name, file.Name)
			assert.Equal(t, tt.IndexAt, file.IndexAt)
		})
	}
}

func Test_QueryFbs(t *testing.T) {
	buf := MakeRootFileFbs(12, "root_test.json", 456)
	root := query.OpenByBuf(buf)
	idx := root.Index()
	assert.Equal(t, uint64(12), idx.File().Id())
	assert.Equal(t, "root_test.json", string(idx.File().Name()))

	buf2 := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})

	root = query.OpenByBuf(buf2)
	assert.Equal(t, len(buf2), root.Len())
	z := root.Index().IndexString().Maps().Last()

	assert.Equal(t, uint64(1), z.Value().FileId())
	assert.Equal(t, int64(2), z.Value().Offset())

	assert.Equal(t, int32(234), root.Index().IndexString().Size())
	assert.Equal(t, 2, root.Index().IndexString().Maps().Count())
}

func Test_QueryNext(t *testing.T) {
	buf := MakeRootRecord(512)
	buf2 := append(buf, MakeRootRecord(513)...)

	root := query.Open(bytes.NewReader(buf2), base.DEFAULT_BUF_CAP)
	z := root.Index().InvertedMapNum()

	record := z.Value()
	val := z.FieldAt(0)
	assert.NotNil(t, val)
	assert.Equal(t, int64(512), root.Index().InvertedMapNum().Key())
	assert.Equal(t, int64(2), record.Offset())
	assert.Equal(t, uint64(1), record.FileId())
	assert.Equal(t, len(buf), root.Len())
	assert.Equal(t, len(buf), root.Len())
	assert.Equal(t, len(buf), z.Info().Pos+z.Info().Size)

	root2 := root.Next()

	assert.Equal(t, int64(513), root2.Index().InvertedMapNum().Key())
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

	q := query.Open(bytes.NewReader(buf), 512)

	list := q.Index().IndexString().Maps()
	n := list.Count()
	_ = n

	infos := map[string]base.Info{}

	info := list.First().ValueInfo(0)
	infos["Maps[0].0"] = base.Info(info)
	for i := 0; i < list.First().CountOfField(); i++ {
		tmpInfo := list.First().ValueInfo(i)
		infos[fmt.Sprintf("Maps[0].0.%d", i)] = base.Info(tmpInfo)

	}

	info2 := list.Last().Value().Info()
	infos["Maps[1].1"] = base.Info(info2)

	assert.Equal(t, true, info.Pos < info2.Pos)

}

func Test_SearchInfo(t *testing.T) {

	buf := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})

	q := query.Open(bytes.NewReader(buf), 512)
	cond := func(pos int, info base.Info) bool {
		return info.Pos <= pos && (info.Pos+info.Size) > pos
		//return true
	}

	result := []base.NodePath{}
	infos := []base.Info{}
	recFn := func(s base.NodePath, info base.Info) {
		result = append(result, s)
		infos = append(infos, info)
	}

	q.SearchInfo(109, recFn, cond)
	sortInfos := make([]base.Info, len(infos))
	copy(sortInfos, infos)

	sort.Slice(sortInfos, func(i, j int) bool {
		return sortInfos[i].Pos < sortInfos[j].Pos
	})

	loncha.Uniq(&sortInfos, func(i int) interface{} {
		return fmt.Sprintf("%d.%d", sortInfos[i].Pos, sortInfos[i].Size)
	})

	for i := 0; i < len(infos); i++ {
		fmt.Printf("%v\t%+v\n", result[i], infos[i])
	}

	assert.True(t, len(infos) > 0)

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

	q := query.Open(bytes.NewReader(buf), 512)
	_ = q
	assert.NotNil(t, buf)
	assert.Equal(t, 2, q.Index().NumList().Num().Count())
	assert.Equal(t, int32(345), q.Index().NumList().Num().At(0))
	assert.Equal(t, int32(584), q.Index().NumList().Num().At(1))

}
