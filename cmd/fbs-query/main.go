package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	tdata "github.com/kazu/fbshelper/query/template"

	"github.com/kazu/fbshelper/fbsparser"
)

var RootName string = "_root"

func main() {

	if len(os.Args) < 3 {
		fmt.Fprint(os.Stderr, "usage: fbs-query fbsfile outputdir\n")
		return
	}

	fbsfile := os.Args[1]
	tmplate := tdata.SturctTeamplate()
	tmplateunion := tdata.UnionTeamplate()

	outDir := os.Args[2]

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

	RootName = parser.Fbs.RootType

	for _, info := range parser.Fbs.Structs {
		newSrc, err := FromTemplate(info, tmplate)
		if err == nil {
			output := filepath.Join(outDir, info.Name+"-query.go")
			ioutil.WriteFile(output, []byte(newSrc), 0644)
		} else {
			fmt.Fprint(os.Stderr, err.Error()+"\n")
		}
	}
	for _, union := range parser.Fbs.Unions {
		newSrc, err := FromTemplate(union, tmplateunion)
		if err == nil {
			output := filepath.Join(outDir, union.Name+"-quuery.go")
			ioutil.WriteFile(output, []byte(newSrc), 0644)
		} else {
			fmt.Fprint(os.Stderr, err.Error()+"\n")
		}
	}

}

func FromTemplate(info interface{}, tmpStr string) (out string, err error) {
	//func FromTemplate(info interface{}, path string) (out string, err error) {
	//tmpStr := structer.LoadFile(path)

	funcMap := template.FuncMap{
		"isMessage":  IsMessage,
		"isSlice":    IsSlice,
		"toCamel":    ToCamelCase,
		"isUnion":    Search,
		"search":     Search,
		"toBareType": ToBareType,
		"isRoot":     IsRoot,
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

func IsRoot(s string) bool {
	return s == RootName
}

func ToCamelCase(s string) string {
	if s[:2] == `[]` {
		return fbsparser.ToCamelCase(s[2:])
	}
	return fbsparser.ToCamelCase(s)
}

func Search(m map[string]bool, s string) bool {
	for key, value := range m {
		if key == s {
			return value
		}
	}
	return false
}

func ToBareType(st string) (s string) {
	s = st
	if s[:2] == `[]` {
		s = s[2:]
	}

	if s[:1] == `*` {
		s = s[1:]
	}

	return
}
