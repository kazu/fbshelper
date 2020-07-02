package info_test

import (
	"bytes"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	vfs_schema "github.com/kazu/fbshelper/example/vfs_schema"

	"io"
	"io/ioutil"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/kazu/fbshelper/info"
	"github.com/stretchr/testify/assert"

	query "github.com/kazu/fbshelper/query"
)

type File struct {
	id       uint64
	name     string
	index_at int64
}

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

func MakeInvertedMapString(b *flatbuffers.Builder, key string) flatbuffers.UOffsetT {
	//, vFn func(b *flatbuffers.Builder) flatbuffers.UOffsetT) flatbuffers.UOffsetT{

	fkey := b.CreateString(key)
	vfs_schema.InvertedMapStringStart(b)
	vfs_schema.InvertedMapStringAddKey(b, fkey)
	vfs_schema.InvertedMapStringAddValue(b, vfs_schema.CreateRecord(b, 1, 2, 3, 4, 5))
	return vfs_schema.InvertedMapStringEnd(b)
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

func LoadRootFileFbs(rio io.Reader) *File {

	raws, e := ioutil.ReadAll(rio)
	if e != nil {
		return nil
	}
	vRoot := vfs_schema.GetRootAsRoot(raws, 0)

	uTable := new(flatbuffers.Table)
	vRoot.Index(uTable)
	fbsFile := new(vfs_schema.File)
	fbsFile.Init(uTable.Bytes, uTable.Pos)

	return &File{id: fbsFile.Id(), name: string(fbsFile.Name()), index_at: fbsFile.IndexAt()}
}

func RecordOption() info.Option {
	return info.Option{
		Maps: []info.OptionType{
			info.OptionType{
				Key:  "file_id",
				Size: 8,
			},
			info.OptionType{
				Key:  "offset",
				Size: 8,
			},
			info.OptionType{
				Key:  "offset_of_value",
				Size: 4,
			},
			info.OptionType{
				Key:  "value_size",
				Size: 4,
			},
		},
	}
}

func FileOption() info.Option {
	return info.Option{
		Maps: []info.OptionType{
			info.OptionType{
				Key:  "file_id",
				Size: 8,
			},
			info.OptionType{
				Key:  "name",
				Size: 0,
			},
			info.OptionType{
				Key:  "index_at",
				Size: 8,
			},
		},
	}
}

func InvertedMapStringOption() info.Option {
	return info.Option{
		Maps: []info.OptionType{
			info.OptionType{
				Key:  "key",
				Size: -1,
			},
			info.OptionType{
				Key:  "Value",
				Size: -1,
				// Nest ?
			},
		},
	}
}

func IndexStringOption() info.Option {
	return info.Option{
		Maps: []info.OptionType{
			info.OptionType{
				Key:  "size",
				Size: 4,
			},
			info.OptionType{
				Key:      "maps",
				Size:     -1,
				IsVector: true,
				//
			},
		},
	}
}

func RootOption() info.Option {
	return info.Option{
		Maps: []info.OptionType{
			info.OptionType{
				Key:  "version",
				Size: 4,
			},
			info.OptionType{
				Key:  "index_type",
				Size: 1,
			},
			info.OptionType{
				Key:  "index",
				Size: -1,
				//	Nest: fileOpt,
			},
		},
	}
}

func RootFileOption() info.Option {
	recordOpt := RecordOption()
	_ = recordOpt
	fileOpt := FileOption()
	rootOpt := RootOption()
	rootOpt.Maps[2].Nest = fileOpt
	return rootOpt
}

func RootIndexStringOption() info.Option {

	recordOpt := RecordOption()

	invMapStr := InvertedMapStringOption()
	invMapStr.Maps[1].Nest = recordOpt

	idxStrOpt := IndexStringOption()
	idxStrOpt.Maps[1].Nest = invMapStr

	rootOpt := RootOption()
	rootOpt.Maps[2].Nest = idxStrOpt
	return rootOpt
}
func int64Align(n int) int {

	if diff := n % 8; diff > 0 {
		return n + (8 - diff)
	}
	return n

}

func Test_GetFbsRootInfo(t *testing.T) {

	info.ENABLE_LOG_DEBUG = true
	buf := MakeRootFileFbs(12, "root_test.json", 456)

	assert.NotNil(t, buf)

	f := LoadRootFileFbs(bytes.NewBuffer(buf))
	assert.NotNil(t, f)

	fbsInfo := info.GetFbsRootInfo(buf, RootFileOption())
	info.Debugf("info %s\n", spew.Sdump(fbsInfo))

	lastb := buf[len(buf)-8:]
	_ = lastb

	assert.Equal(t, len(buf), int64Align(int(fbsInfo.Length)))
	//assert.Equal(t, buf[94:96], buf[92:94])

	fbsInfo.FetchAll(buf, RootFileOption())
	assert.Equal(t, fbsInfo.Table[0], uint64(1))

}

func Test_GetFbsRootInfo_RootIndexString(t *testing.T) {
	info.ENABLE_LOG_DEBUG = true

	buf := MakeRootIndexString(func(b *flatbuffers.Builder) flatbuffers.UOffsetT {
		return MakeIndexString(b, func(b *flatbuffers.Builder, i int) flatbuffers.UOffsetT {
			return MakeInvertedMapString(b, fmt.Sprintf("     %d", i))
		})
	})
	assert.NotNil(t, buf)

	fbsInfo := info.GetFbsRootInfo(buf, RootIndexStringOption())
	info.Debugf("info %s\n", spew.Sdump(fbsInfo))
	fbsInfo.FetchAll(buf, RootIndexStringOption())

	assert.Equal(t, len(buf), int64Align(int(fbsInfo.Length)))

}

func Test_QueryFbs(t *testing.T) {
	buf := MakeRootFileFbs(12, "root_test.json", 456)
	root := query.OpenByBuf(buf)
	idx := root.Index()
	assert.Equal(t, uint64(12), idx.File().Id())
	assert.Equal(t, "root_test.json", string(idx.File().Name()))
}
