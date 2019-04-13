package fbsparser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/kazu/lonacha/structer"
)

type StructInfo struct {
	structer.StructInfo
	IsMessage map[string]bool
	IsSlice   map[string]bool
}

type Fbs struct {
	NameSpace string
	typeName  string
	Structs   []StructInfo
	Unions    []Union
	Enums     []Enum

	// store instantly
	Fields     map[string]string
	fname      string
	fvalue     string
	isRepeated bool
}

type Union struct {
	Name    string
	Aliases []string
}

type Enum struct {
	Name    string
	Type    string
	Defines map[string]string
}

func (fbs *Fbs) FinilizeForFbs() {
	for _, info := range fbs.Structs {
		for fname, ftype := range info.Fields {
			// slice of fbs type
			if ftype[:2] == `[]` && strings.ToUpper(ftype[:3]) == ftype[:3] {
				info.Fields[fname] = fmt.Sprintf("[]*%sMessage", ftype[2:])
				continue
			}
			if strings.ToUpper(ftype[:1]) == ftype[:1] {
				info.Fields[fname] = fmt.Sprintf("*%sMessage", ftype)
			}
		}
	}
}

func (fbs *Fbs) SetNameSpace(s string) {
	fbs.NameSpace = s
}

func (fbs *Fbs) ExtractStruct() {

	structInfo := StructInfo{}
	structInfo.PkgName = fbs.NameSpace
	structInfo.Name = fbs.typeName
	structInfo.Fields = map[string]string{}
	structInfo.IsSlice = map[string]bool{}
	structInfo.IsMessage = map[string]bool{}

	structInfo.Fields = fbs.Fields

	fbs.Structs = append(fbs.Structs, structInfo)

	fbs.Fields = map[string]string{}
}

func (fbs *Fbs) SetTypeName(s string) {
	fbs.typeName = s
}

func (fbs *Fbs) FieldNaame(s string) {
	fbs.fname = toCamelCase(s)
}
func (fbs *Fbs) SetType(s string) {
	fbs.fvalue = s
}

func (fbs *Fbs) NewExtractField() {
	if fbs.Fields == nil {
		fbs.Fields = map[string]string{}
	}
	fbs.Fields[fbs.fname] = fbs.fvalue
}

func (fbs *Fbs) NewExtractFieldWithValue() {
	fbs.Fields[fbs.fname] = fbs.fvalue
}

func (fbs *Fbs) SetRepeated(s string) {
	//	fbs.isRepeated = true
	fbs.fvalue = "[]" + fbs.fvalue
}

func (fbs *Fbs) NewUnion(s string) {
	union := Union{
		Name:    s,
		Aliases: []string{},
	}
	for key, _ := range fbs.Fields {
		union.Aliases = append(union.Aliases, key)
	}

	fbs.Unions = append(fbs.Unions, union)

	fbs.Fields = map[string]string{}
}

func toCamelCase(st string) string {
	s := strings.Split(st, "_")
	for i := 0; i < len(s); i++ {
		w := []rune(strings.ToLower(s[i]))
		w[0] = unicode.ToUpper(w[0])
		s[i] = string(w)
	}
	return strings.Join(s, "")
}

func ToCamelCase(st string) string {
	return toCamelCase(st)
}
