# fbs-query 

generate routine easy accessing flatbuffers data.

## usage sample

### generate source code
```console
$ cat index.fbs
namespace vfs_schema;

union Index {
  File,
  Files
}

table File {
    id:uint64; // file inode number
    name:string;
    index_at:int64;
}

table Files {
    datas:[File];
}

table Root {
    version:int;
    index:Index;
}

root_type Root;

$ fbs-query index.fbs tmp/
```


### read flatbuffers data in flac's generated code

```go
import (
    query "./tmp"
)

// buf is []byte

vRoot := vfs_schema.GetRootAsRoot(buf, 0)

uTable := new(flatbuffers.Table)
vRoot.Index(uTable)
fbsFile := new(vfs_schema.File)
fbsFile.Init(uTable.Bytes, uTable.Pos)
  
fbsFile.Id()
fbsFile.Name()
fbsFile.IndexAt()
```


### read via fbs-query's generated code

```go
q := query.OpenByBuf(buf)
q.Len()
q.Index().File().Id()
q.Index().File().Name()
q.Index().File().IndexAt()

```

### read vector data 

```go
fbs := query.Open(ioReader)
fbs.Files().First()
fbs.Files().Last()
fbs.Files().Len()
fbs.Files().At(1)
fbs.Files().All()
fbs.Files().Select(func(m query.FbsFile) bool {
    return m.Id() == 10
})
// for streaming data
fbs.Next().Files().First()

```

### unmarshal 

```go

f := struct{
    ID      uint64 `fbs:"Id"`
	Name    []byte `fbs:"Name"`
	IndexAt int64  `fbs:"IndexAt"`
}{}

fbs := query.OpenByBuf(buf)
fbs.Files().First().Unmarshal(&f)

```


## TODO

- [x] change base.Base when call Next()
- [ ] support basic type slice ( []int, ... ), define FbsType ?
- [ ] support writing
- [ ] marshal
- [ ] unmarshal nested Table/Struct