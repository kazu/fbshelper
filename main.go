package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kazu/fbshelper/fbsparser"
	"github.com/kazu/lonacha/structer"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprint(os.Stderr, "usage: fbshelper fbsfile templatefile outputdir")
		return
	}

	fbsfile := os.Args[1]
	tmplate := os.Args[2]
	outDir := os.Args[3]

	bytes, err := ioutil.ReadFile(fbsfile)
	if err != nil {
		fmt.Fprint(os.Stderr, "cannot read file"+fbsfile)
	}

	parser := &fbsparser.Parser{Buffer: string(bytes)}

	parser.Init()
	err = parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error %s err=%s", fbsfile, err)

	}

	parser.Execute()
	parser.Fbs.FinilizeForFbs()

	for _, info := range parser.Fbs.Structs {

		newSrc, err := FromTemplate(info, tmplate)
		if err == nil {
			output := filepath.Join(outDir, info.Name+".fbshelper.go")
			ioutil.WriteFile(output, []byte(newSrc), 0644)
		} else {
			fmt.Fprint(os.Stderr, err.Error()+"\n")
		}
	}
}

func FromTemplate(info fbsparser.StructInfo, path string) (out string, err error) {
	tmpStr := structer.LoadFile(path)

	funcMap := template.FuncMap{
		"isMessage": IsMessage,
		"isSlice":   IsSlice,
		"toCamel":   ToCamelCase,
	}

	t := template.Must(template.New("info").Funcs(funcMap).Parse(tmpStr))
	s := &strings.Builder{}
	err = t.Execute(s, info)
	out = s.String()
	return
}

func IsMessage(s string) bool {
	if s[:2] == `[]` {
		if strings.ToUpper(s[2:3]) == s[2:3] {
			return true
		}
		return false
	}
	if strings.ToUpper(s[:1]) == s[:1] {
		return true
	}
	return false
}

func IsSlice(s string) bool {
	//fmt.Printf("IsSlice s=%s result=%v\n", s, s[:2] == `[]`)

	if s[:2] == `[]` {
		return true
	}
	return false
}

func ToCamelCase(s string) string {
	if s[:2] == `[]` {
		return fbsparser.ToCamelCase(s[2:])
	}
	return fbsparser.ToCamelCase(s)
}
