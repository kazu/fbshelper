package gen

//go:generate go run github.com/jteeuwen/go-bindata/go-bindata -pkg query -o ../query/template/genny-load.go ../template/genny
//go:generate go run load-query-template.go ../template
