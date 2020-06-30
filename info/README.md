# fbshelper/info 

## introduction

this is simple program getting [flatbuffers] data infomation.

# usage 

```go

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
    fileOpt := info.Option{
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
    
    var buf byte[]
    //... reading fbs encoded data to buf

    fbsInfo := info.GetFbsRootInfo(buf, fileOpt)

    fmt.Printf("data length=%d\n", fbsInfo.Length)

```

### Authorization code grant

[flatbuffers]:https://google.github.io/flatbuffers/
