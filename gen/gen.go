package gen

import (
	_ "github.com/jteeuwen/go-bindata"
)

//go:generate go run github.com/jteeuwen/go-bindata/go-bindata -pkg query -o ../query/template/genny-load.go ../template/genny
