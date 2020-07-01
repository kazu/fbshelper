package info_test

import (
	"bytes"

	"github.com/davecgh/go-spew/spew"
	"github.com/kazu/fbshelper/example/vfs_schema"

	"io"
	"io/ioutil"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/kazu/fbshelper/info"
	"github.com/stretchr/testify/assert"
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
	_ = recodOpt
	fileOpt := info.Option{
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
				Nest: fileOpt,
			},
		},
	}
}

func Test_GetFbsRootInfo(t *testing.T) {

	info.ENABLE_LOG_DEBUG = true
	buf := MakeRootFileFbs(12, "root_file.json", 456)

	assert.NotNil(t, buf)

	f := LoadRootFileFbs(bytes.NewBuffer(buf))
	assert.NotNil(t, f)

	fbsInfo := info.GetFbsRootInfo(buf, RootFileOption())
	info.Debugf("info %s\n", spew.Sdump(fbsInfo))

	assert.Equal(t, len(buf), int(fbsInfo.Length))

	fbsInfo.FetchAll(buf, RootFileOption())
	assert.Equal(t, fbsInfo.Table[0], uint64(1))

}
