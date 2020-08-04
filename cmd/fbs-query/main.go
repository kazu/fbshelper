package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	tdata "github.com/kazu/fbshelper/query/template"

	"github.com/kazu/fbshelper/fbsparser"
	base "github.com/kazu/fbshelper/query/base"
)

var RootName string = "_root"
var NameisStruct map[string]bool = map[string]bool{}

/*
var _bindata = map[string]func() ([]byte, error){
	"../template/genny/dummy.go": template_genny_dummy_go,
	"../template/genny/field.go": template_genny_field_go,
	"../template/genny/node.go": template_genny_node_go,
	"../template/genny/node_test.go": template_genny_node_test_go,
	"../template/genny/root.go": template_genny_root_go,
	"../template/genny/union-alias.go": template_genny_union_alias_go,
	"../template/genny/union.go": template_genny_union_go,
}
*/

func writeTmp(dir string, file string, data []byte) error {
	return ioutil.WriteFile(filepath.Join(dir, file), data, os.ModePerm)
}

var GennyLineTemplate string = `//go:generate go run github.com/cheekybits/genny -in=%s -out=%s  gen `

func MakeGenGennies(parser *fbsparser.Parser, dir string, opt GennyOpt) (result string, e error) {

	UseFilter := true
	if opt.filter == "" {
		UseFilter = false
	}
	filter := opt.filter

	var bb strings.Builder
	b := &bb
	odir := "../"
	tdir := filepath.Join(dir, "genny")

	e = os.MkdirAll(tdir, os.ModePerm)
	if e != nil {
		return "", e
	}

	defer func() {
		//os.RemoveAll(tdir)
	}()
	for _, fpath := range tdata.AssetNames() {
		basePath := filepath.Base(fpath)
		data, e := tdata.Asset(fpath)
		if e != nil {
			return "", e
		}

		e = writeTmp(tdir, basePath, data)
		if e != nil {
			return "", e
		}
	}
	tmpHead := `
	package main
	import (
		_ "github.com/cheekybits/genny"
	)
	
	`

	fmt.Fprintf(b, "%s\n\n", tmpHead)
	tdir = "../genny/"
	odir = "../"

	oPath := filepath.Join(odir, "root.go")
	iPath := filepath.Join(tdir, "root.go")

	fmt.Fprintf(b, GennyLineTemplate, iPath, oPath)
	fmt.Fprintf(b, ` "RootType=%s"`, RootName)
	fmt.Fprintf(b, "\n")

	for _, info := range parser.Fbs.Structs {
		if UseFilter && filter != info.Name {
			continue
		}
		iPath = filepath.Join(tdir, "node.go")
		oPath = filepath.Join(odir, info.Name+".gen.go")

		fmt.Fprintf(b, GennyLineTemplate, iPath, oPath)
		fmt.Fprintf(b, ` "NodeName=%s`, info.Name)
		fmt.Fprintf(b, ` RootType=%s`, RootName)
		fmt.Fprintf(b, ` IsStruct=%v"`, IsStruct(info.Name))

		fmt.Fprintf(b, "\n")

		iPath = filepath.Join(tdir, "list.go")
		oPath = filepath.Join(odir, info.Name+"List.gen.go")

		fmt.Fprintf(b, GennyLineTemplate, iPath, oPath)
		fmt.Fprintf(b, ` "NodeName=%s`, info.Name)
		fmt.Fprintf(b, ` ListType=%s"`, info.Name+"List")
		fmt.Fprintf(b, "\n")

		iPath = filepath.Join(tdir, "field.go")
		for idx, field := range info.Fields {
			oPath = filepath.Join(odir, info.Name+"."+field.Name+".gen.go")
			fmt.Fprintf(b, GennyLineTemplate, iPath, oPath)
			fmt.Fprintf(b, ` "NodeName=%s`, info.Name)
			fmt.Fprintf(b, " FieldName=%s", field.Name)
			fmt.Fprintf(b, " FieldType=%s", field.Type)

			if _, ok := base.NameToType[field.Type]; ok {
				fmt.Fprintf(b, " ResultType=%s", "CommonNode")
			} else if IsBasicTypeSlice(field.Type) {
				fmt.Fprintf(b, " ResultType=%s", "List")
			} else if IsSlice(field.Type) {
				fmt.Fprintf(b, " ResultType=%s", ConvResultType(field.Type))
			} else {
				fmt.Fprintf(b, " ResultType=%s", field.Type)
			}
			fmt.Fprintf(b, " FieldNum=%d", idx)
			fmt.Fprintf(b, ` IsStruct=%v"`, IsStruct(field.Type))
			fmt.Fprintf(b, "\n")
		}

	}

	for _, union := range parser.Fbs.Unions {
		if UseFilter && filter != union.Name {
			continue
		}

		iPath = filepath.Join(tdir, "union.go")
		oPath = filepath.Join(odir, union.Name+".union.gen.go")
		fmt.Fprintf(b, GennyLineTemplate, iPath, oPath)
		fmt.Fprintf(b, ` "UnionName=%s"`, union.Name)
		fmt.Fprintf(b, "\n")

		iPath = filepath.Join(tdir, "union-alias.go")
		oPath = filepath.Join(odir, union.Name+".union-alias.gen.go")
		fmt.Fprintf(b, GennyLineTemplate, iPath, oPath)

		fmt.Fprintf(b, `"UnionName=%s AliasName=%s"`, union.Name,
			strings.Join(union.Aliases, ","))
		fmt.Fprintf(b, "\n")
	}

	return b.String(), e
}
func ConvResultType(o string) string {
	s := strings.ReplaceAll(o, "[]byte", "CommonNode")
	if len(s) < 2 {
		return s
	}
	if s[:2] == "[]" {
		s = s[2:] + "List"
	}
	return s
}

func RunGenerate(dir, s string, isNotDryRun bool) error {

	os.MkdirAll(dir, os.ModePerm)

	genPath := filepath.Join(dir, "tmp-gen-fbs-query.go")
	f, err := os.Create(genPath)
	if err != nil {
		return err
	}
	f.WriteString(s)
	f.Close()
	fmt.Fprintf(os.Stdout, "created %s\n", genPath)
	fmt.Fprintf(os.Stdout, "runnnig ... %s\n", genPath)
	if isNotDryRun {
		out, err := exec.Command("go", "generate", "-v", dir).Output()
		fmt.Println(out)
		if err != nil {
			return err
		}
		os.RemoveAll(dir)
		os.RemoveAll(filepath.Join(dir, "../", "genny"))
	}
	fmt.Fprintf(os.Stdout, "run(dry-run=%v) go generate %s\n", !isNotDryRun, genPath)

	return nil
}

type GennyOpt struct {
	filter  string
	verbose bool
}

func GennyMode(parser *fbsparser.Parser, outDir string, opt GennyOpt) {

	s, e := MakeGenGennies(parser, outDir, opt)
	if e != nil {
		fmt.Fprintf(os.Stderr, "make gen,go filer %s err=%v", e)
		return
	}
	if opt.verbose {
		fmt.Println(s)
	}
	e = RunGenerate(filepath.Join(outDir, "/gen"), s, true)
	if e != nil {
		fmt.Fprintf(os.Stderr, "run gen.go fail err=%v", e)
		return
	}
	return

}

func TextTemplateMode(parser *fbsparser.Parser, outDir, tmplate, tmplateunion string) {
	for _, info := range parser.Fbs.Structs {
		newSrc, err := FromTemplate(info, tmplate)
		if err == nil {
			output := filepath.Join(outDir, info.Name+"-query.go")
			ioutil.WriteFile(output, []byte(comment()+newSrc), 0644)
		} else {
			fmt.Fprint(os.Stderr, err.Error()+"\n")
		}
	}
	for _, union := range parser.Fbs.Unions {
		newSrc, err := FromTemplate(union, tmplateunion)
		if err == nil {
			output := filepath.Join(outDir, union.Name+"-quuery.go")
			ioutil.WriteFile(output, []byte(comment()+newSrc), 0644)
		} else {
			fmt.Fprint(os.Stderr, err.Error()+"\n")
		}
	}

}

const Usage string = `
fbs query access generator
Usage:
  fbs-query 

  Flags:
  	-h      help for gist
	-fbs	input fbs file
	-out 	output directory
`

func main() {

	var filter, fbsfile, outDir string
	var isHelp, verbose bool
	flag.StringVar(&filter, "filter", "", "generate filter")
	flag.StringVar(&fbsfile, "fbs", "", "input flatbuffers schema")
	flag.StringVar(&outDir, "out", "", "output directory")
	flag.BoolVar(&isHelp, "h", false, "Help about command")
	flag.BoolVar(&verbose, "v", false, "verbose output")

	flag.Parse()
	if isHelp {
		fmt.Println(Usage)
	}

	if len(os.Args) < 3 {
		fmt.Println(Usage)
		return
	}
	//fmt.Println(filter, fbsfile, outDir)
	//return

	tmplate := tdata.SturctTeamplate()
	tmplateunion := tdata.UnionTeamplate()
	_, _ = tmplate, tmplateunion

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
		if !info.IsTable {
			NameisStruct[info.Name] = true
		}
	}

	GennyMode(parser, outDir, GennyOpt{filter: filter, verbose: verbose})
}

func comment() string {
	return `
// Code generated by genmaps.go; DO NOT EDIT.
// template file is https://github.com/kazu/fbshelper/blob/master/template/query.go.tmpl github.com/kazu/fbshelper/template/query.go.tmpl 
//   https://github.com/kazu/fbshelper/blob/master/template/union.query.go.tmpl

`
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
		"add":        Add,
		"isStruct":   IsStruct,
	}

	t := template.Must(template.New("info").Funcs(funcMap).Parse(tmpStr))
	s := &strings.Builder{}
	err = t.Execute(s, info)
	out = s.String()
	return
}
func Add(i, j int) int {
	return i + j
}

func IsStruct(s string) bool {
	return NameisStruct[s]
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

func IsBasicTypeSlice(s string) bool {
	return IsSlice(s) && base.HasNameType(s[2:])
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
