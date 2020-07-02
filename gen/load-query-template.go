// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "not specified tempalte directory")
		return
	}
	templateDir := os.Args[1]
	base := "query.go.tmpl"
	bufStruct, e := ioutil.ReadFile(filepath.Join(templateDir, base))
	if e != nil {
		fmt.Fprint(os.Stderr, "not load template file")
	}
	bufUnion, e := ioutil.ReadFile(filepath.Join(templateDir, "union."+base))
	if e != nil {
		fmt.Fprint(os.Stderr, "not load template file")
	}

	template := `package query 

	func SturctTeamplate() string {
		return %s
	}

	func UnionTeamplate() string {
		return %s 
	}
`
	f, e := os.Create("../query/template/load.go")
	fmt.Fprintf(f, template, Quote(string(bufStruct)), Quote(string(bufUnion)))
	f.Close()

}

func Quote(s string) string {
	return "`" + s + "`"
}
