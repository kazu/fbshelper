# fbs-query 

generate routine easy accessing flatbuffers data.

generate source code
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


you read flatbuffers data in flac's generated code

```go
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


read via fbs-query's generated code

```go
q := query.OpenByBuf(buf)
q.Len()
q.Index().File().Id()
q.Index().File().Name()
q.Index().File().IndexAt()

```

read vector data 

```go
query.OpenByBuf(buf).Files().First()
query.OpenByBuf(buf).Files().Last()
query.OpenByBuf(buf).Files().Len()
query.OpenByBuf(buf).Files().At(1)
query.OpenByBuf(buf).Files().All()
query.OpenByBuf(buf).Files().Select(func(m query.FbsFile) bool {
    return m.Id() == 10
})

```
