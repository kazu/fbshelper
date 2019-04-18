package fbsparser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/kazu/lonacha/structer"
)

type StructInfo struct {
	structer.StructInfo
	IsTable   bool
	IsMessage map[string]bool
	IsSlice   map[string]bool
	IsUnion   map[string]bool
	IsStruct  map[string]bool
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
	uName      string
}

type Union struct {
	PkgName string
	Name    string
	Aliases []string
}

type Enum struct {
	Name    string
	Type    string
	Defines map[string]string
}

func (fbs *Fbs) isUnion(s string) bool {

	/*for _, union := range fbs.Unions {
		for _, alias := range union.Aliases {
			if s == alias {
				return true
			}
		}
	}*/

	for _, union := range fbs.Unions {
		if union.Name == s {
			fmt.Printf("union=%+v \n", fbs.Unions)
			fmt.Printf("union.Name=%s %s\n", union.Name, s)
			return true
		}
	}

	return false
}

func (fbs *Fbs) SearchNotTableMessage() {
	for _, info := range fbs.Structs {
		for fname, oftype := range info.Fields {
			ftype := oftype
			if ftype[:2] == `[]` {
				ftype = ftype[2:]
			}
			if ftype[:1] == `*` {
				ftype = ftype[1:]
			}
			for _, sinfo := range fbs.Structs {
				if info.Name == sinfo.Name {
					continue
				}
				if sinfo.Name+"Message" == ftype && !sinfo.IsTable {
					info.IsStruct[fname] = true
					break
				}
			}
		}
	}
}

func (fbs *Fbs) FinilizeForFbs() {
	for _, info := range fbs.Structs {
		newFields := info.Fields
		for fname, ftype := range info.Fields {
			//info.IsUnion[fname] = false
			// slice of fbs type

			if ftype[:2] == `[]` && strings.ToUpper(ftype[:3]) == ftype[:3] {
				newFields[fname] = fmt.Sprintf("[]*%sMessage", ftype[2:])
				if fbs.isUnion(ftype[2:]) {
					info.IsUnion[fname] = true
				}
				continue
			}
			if ftype[:2] == `[]` {
				continue
			}
			if strings.ToUpper(ftype[:1]) == ftype[:1] {
				newFields[fname] = fmt.Sprintf("*%sMessage", ftype)
				if fbs.isUnion(ftype) {
					info.IsUnion[fname] = true
				}
				continue
			}
		}
		info.Fields = newFields
	}
	fbs.SearchNotTableMessage()
}

func (fbs *Fbs) SetNameSpace(s string) {
	fbs.NameSpace = s
}

func (fbs *Fbs) ExtractStruct(isTable bool) {

	structInfo := StructInfo{}
	structInfo.PkgName = fbs.NameSpace
	structInfo.Name = fbs.typeName
	structInfo.IsTable = isTable
	structInfo.Fields = map[string]string{}
	structInfo.IsSlice = map[string]bool{}
	structInfo.IsMessage = map[string]bool{}
	structInfo.IsUnion = map[string]bool{}
	structInfo.IsStruct = map[string]bool{}

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
	if s == "" {
		fbs.fvalue = "[]" + fbs.fvalue
	} else {
		fbs.fvalue = "[]" + s
	}
}

func (fbs *Fbs) UnionName(s string) {
	fbs.uName = s
}

func (fbs *Fbs) NewUnion() {
	union := Union{
		PkgName: fbs.NameSpace,
		Name:    fbs.uName,
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
