package base_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	query "github.com/kazu/fbshelper/example/query2"
	query2 "github.com/kazu/fbshelper/example/query2"
	"github.com/kazu/fbshelper/example/vfs_schema"
	"github.com/kazu/fbshelper/query/base"
	log "github.com/kazu/fbshelper/query/log"
)

func Benchmark_DirectBuf(b *testing.B) {

	benchs := []struct {
		name    string
		configs []Config
	}{
		{
			name:    "directbuf block=8k  merge-interval-8",
			configs: []Config{loopCnt(b.N), blockSize(8192), useMerge(true), mInterval(8), useNoLayer(true)},
		},
		{
			name:    "directbuf block=4k  merge-interval-8",
			configs: []Config{loopCnt(b.N), blockSize(4096), useMerge(true), mInterval(8), useNoLayer(true)},
		},
		{
			name:    "directbuf block=512  merge-interval-8",
			configs: []Config{loopCnt(b.N), blockSize(512), useMerge(true), mInterval(8), useNoLayer(true)},
		},
		{
			name:    "directbuf block=4k merge-interval-4",
			configs: []Config{loopCnt(b.N), blockSize(4096), useMerge(true), mInterval(4), useNoLayer(true)},
		},
		{
			name:    "directbuf block=512 merge-interval-4",
			configs: []Config{loopCnt(b.N), blockSize(512), useMerge(true), mInterval(4), useNoLayer(true)},
		},
		{
			name:    "directbuf block=4k merge-interval-2",
			configs: []Config{blockSize(4096), useMerge(true), mInterval(2), useNoLayer(true)},
		},
		{
			name:    "directbuf block=512 merge-interval-2",
			configs: []Config{blockSize(512), useMerge(true), mInterval(2), useNoLayer(true)},
		},
	}

	for _, bb := range benchs {
		b.ResetTimer()
		b.Run(bb.name, func(b *testing.B) {
			b.StartTimer()
			_, deferfn := MakeDirectBufFileList(append(bb.configs, loopCnt(b.N))...)
			b.StopTimer()
			deferfn()
		})

	}

}

func MakeRootWithOutRecord(key uint64, fID uint64, offset int64, size int64) *query.Root {

	root := query.NewRoot()

	root.SetVersion(query.FromInt32(1))
	root.WithHeader()

	inv := query.NewInvertedMapNum()
	inv.SetKey(query.FromInt64(int64(key)))

	rec := query.NewRecord()
	rec.SetFileId(query.FromUint64(fID))
	rec.SetOffset(query.FromInt64(offset))
	rec.SetSize(query.FromInt64(size))
	rec.SetOffsetOfValue(query.FromInt32(0))
	rec.SetValueSize(query.FromInt32(0))

	root.SetRecord(rec)

	root.SetIndexType(query.FromByte(byte(vfs_schema.IndexInvertedMapNum)))

	root.SetIndex(&query.Index{CommonNode: inv.CommonNode})

	root.Flatten()
	root.Base = base.NewNoLayer(root.Base)
	return root
}

func MakeRootWithOutFile() *query.Root {

	root := query.NewRoot()

	root.SetVersion(query.FromInt32(1))
	root.WithHeader()

	//hoge := query2.NewHoges()

	// inv := query.NewInvertedMapNum()
	// inv.SetKey(query.FromInt64(int64(key)))

	// rec := query.NewRecord()
	// rec.SetFileId(query.FromUint64(fID))
	// rec.SetOffset(query.FromInt64(offset))
	// rec.SetSize(query.FromInt64(size))
	// rec.SetOffsetOfValue(query.FromInt32(0))
	// rec.SetValueSize(query.FromInt32(0))

	//root.SetRecord(rec)

	root.SetIndexType(query.FromByte(byte(vfs_schema.IndexHoges)))

	//root.SetIndex(&query.Index{CommonNode: hoge.CommonNode})

	root.Flatten()
	root.Base = base.NewNoLayer(root.Base)
	return root
}

func Benchmark_AddList(b *testing.B) {

	o := log.CurrentLogLevel
	base.SetL2Current(log.LOG_WARN, base.FLT_IS)
	defer base.SetL2Current(o, base.FLT_NORMAL)

	appendFileListUseAdd := func(dlist, slist *query.FileList) {
		dlist.Add(*slist)
	}

	appendFileListUseSetAt := func(dlist, slist *query.FileList) {
		for _, f := range slist.All() {
			dlist.SetAt(dlist.Count(), f)
		}
	}
	_ = appendFileListUseSetAt

	benchs := []struct {
		name       string
		isFileList bool
		isNoLayer  bool
		useRoot    bool
		useFlatten bool
		useMovOff  bool
		cnt        int
		Fn         func(dlist, slist *query.FileList)
	}{
		{"add filelist_1001 use AddAt() ", true, true, true, false, true, 1001, appendFileListUseAdd},
		{"add filelist_1001 use SetAt() ", true, true, true, false, true, 1001, appendFileListUseSetAt},
		{"add filelist_1001 use AddAt() ", true, true, false, false, false, 1001, appendFileListUseAdd},
		{"add filelist_1001 use SetAt() ", true, true, false, false, false, 1001, appendFileListUseSetAt},
		{"add filelist_1001 use AddAt() ", true, true, false, true, false, 1001, appendFileListUseAdd},
		{"add filelist_1001 use SetAt() ", true, true, false, true, false, 1001, appendFileListUseSetAt},
		{"add filelist_1001 use AddAt() ", true, true, false, false, true, 1001, appendFileListUseAdd},
		{"add filelist_1001 use SetAt() ", true, true, false, false, true, 1001, appendFileListUseSetAt},
		{"add filelist_1001 use AddAt() ", true, true, false, true, true, 1001, appendFileListUseAdd},
		{"add filelist_1001 use SetAt() ", true, true, false, true, true, 1001, appendFileListUseSetAt},
		{"add filelist_402  use AddAt()", true, true, false, false, true, 402, appendFileListUseAdd},
		{"add filelist_402  use SetAt()", true, true, false, false, true, 402, appendFileListUseSetAt},
		{"add filelist_202  use AddAt()", true, true, false, false, true, 202, appendFileListUseAdd},
		{"add filelist_202  use SetAt()", true, true, false, false, true, 202, appendFileListUseSetAt},
	}

	flag2str := func(flags ...bool) string {
		var b strings.Builder
		for _, flag := range flags {
			if flag {
				b.WriteString("T")
				continue
			}
			b.WriteString("F")
		}
		return b.String()
	}

	for _, bb := range benchs {

		b.Run(fmt.Sprintf("%s flag=%s", bb.name, flag2str(bb.isFileList, bb.isNoLayer, bb.useRoot, bb.useFlatten, bb.useMovOff)), func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				var root *query2.Root

				slist := MakeFileList(bb.isNoLayer, bb.cnt/2, uint64(rand.Intn(bb.cnt)), uint64(rand.Intn(bb.cnt/10)), int64(rand.Intn(bb.cnt/10)), "benchAdding")
				dlist := MakeFileList(bb.isNoLayer, bb.cnt/2, uint64(rand.Intn(bb.cnt)), uint64(rand.Intn(bb.cnt/10)), int64(rand.Intn(bb.cnt/10)), "benchAdding")

				if bb.useRoot {
					root = MakeRootWithOutFile()

					hoges := query2.NewHoges()
					hoges.Base = base.NewNoLayer(hoges.Base)
					hoges.SetFiles(dlist)
					root.SetIndex(&query.Index{CommonNode: hoges.CommonNode})
					dlist = root.Index().Hoges().Files()
				}
				if bb.useFlatten {
					dlist.Flatten()
					//slist.Flatten()
				}

				base.OptUseMovOff(bb.useMovOff)(&base.CurrentGlobalConfig)

				b.StartTimer()

				bb.Fn(dlist, slist)

				b.StopTimer()

			}
		})
	}

}
