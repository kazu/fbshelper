package info_test

import (
	"github.com/kazu/fbshelper/example/vfs_schema"
	"bytes"

	"testing"
	"io"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
	"github.com/kazu/fbshelper/info"
	
)

type File struct {
	id       uint64
	name     string
	index_at int64
}

func MakeRootFileFbs(id uint64, name string, index_at int64) []byte {

	b := flatbuffers.NewBuilder(0)
	fname := b.CreateString(name)

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

func RootFileOption() info.Option {

	recodOpt := info.Option{
		Maps: []info.OptionType{
			info.OptionType{
				Key: "file_id",
				Size: 8,
			},
			info.OptionType{
				Key: "offset",
				Size: 8,
			},
			info.OptionType{
				Key: "offset_of_value",
				Size: 4,
			},
			info.OptionType{
				Key: "value_size",
				Size: 4,
			},
		},
	}


	return info.Option{
		Maps: []info.OptionType{
			info.OptionType{
				Key: "version",
				Size: 4,
			},
			info.OptionType{
				Key: "index_type",
				Size: 1,
			},
			info.OptionType{
				Key: "index",
				Size: -1,
				Nest: recodOpt,
			},
		},
	}
}


func Test_FbsFile(t *testing.T) {

	buf := MakeRootFileFbs(12, "root_file.json", 456)

	assert.NotNil(t, buf)

	f := LoadRootFileFbs(bytes.NewBuffer(buf))
	assert.NotNil(t, f)


	info := info.GetFbsRootInfo(buf, RootFileOption())
	spew.Dump(info)
	assert.Equal(t, info.Table[0], uint64(1))
}

/*
func Test_StreamFbs(t *testing.T) {

	r := vfs.NewRecord(1000, 1234, 12345)

	buf := r.ToFbs(uint64(12))
	assert.Equal(t, 88, len(buf))

	r2 := vfs.NewRecord(1001, 2234, 123)
	buf = append(buf, r2.ToFbs(uint64(15))...)
	n := flatbuffers.GetUOffsetT(buf[0:4])

	bbuf := bytes.NewBuffer(buf)
	// bbuf := bytes.NewBuffer(buf[0:GetFbsSize(buf)])

	r3 := vfsindex.RecordFromFbs(bbuf)
	info := GetFbsRootInfo(buf)
	spew.Dump(info)

	assert.Equal(t, info.table[0], uint64(1))
	assert.NotNil(t, r3)
	//assert.Equal(t, 10, GetFbsSize(buf))
	assert.Equal(t, uint32(16), uint32(n))
}
*/