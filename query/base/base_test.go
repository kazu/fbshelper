package base_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/kazu/fbshelper/example/vfs_schema"
	flatbuffers "github.com/google/flatbuffers/go"
	query "github.com/kazu/fbshelper/query/base"
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
func TestBase(t *testing.T) {


	buf := MakeRootFileFbs(12, "root_test1.json", 456)
	assert.NotNil(t, buf)

	q := query.OpenByBuf(buf)

	assert.Equal(t, int32(1),         	q.Version())
	assert.Equal(t, uint64(12),       	q.Index().File().Id(),)
	assert.Equal(t, "root_test1.json", 	q.Index().File().Name())
	assert.Equal(t, int64(456),			q.Index().File().IndexAt())
}