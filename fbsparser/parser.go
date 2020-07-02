package fbsparser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/kazu/loncha/structer"
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
	RootType  string
	Structs   []StructInfo
	Unions    []Union
	Enums     []Enum

	// store instantly
	Fields     []structer.FieldInfo
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
		for _, fInfo := range info.Fields {
			fname := fInfo.Name
			oftype := fInfo.Type
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
	for i, info := range fbs.Structs {
		//newFields := info.Fields
		newFields := []structer.FieldInfo{}
		for _, fInfo := range info.Fields {
			fname := fInfo.Name
			ftype := fInfo.Type

			if ftype[:2] == `[]` && strings.ToUpper(ftype[:3]) == ftype[:3] {
				//newFields = append(newFields, structer.FieldInfo{Name: fname, Type: fmt.Sprintf("[]*%sMessage", ftype[2:])})
				newFields = append(newFields, structer.FieldInfo{Name: fname, Type: fmt.Sprintf("[]%s", ftype[2:])})
				if fbs.isUnion(ftype[2:]) {
					info.IsUnion[fname] = true
					//unionIdx = append(unionIdx, idx)
				}
				continue
			}
			if ftype[:2] == `[]` {
				newFields = append(newFields, fInfo)
				continue
			}
			if strings.ToUpper(ftype[:1]) == ftype[:1] {
				if fbs.isUnion(ftype) {
					info.IsUnion[fname] = true
					newFields = append(newFields, structer.FieldInfo{Name: fname + "Type", Type: "byte"})
				}
				//newFields = append(newFields, structer.FieldInfo{Name: fname, Type: fmt.Sprintf("*%sMessage", ftype)})
				newFields = append(newFields, structer.FieldInfo{Name: fname, Type: fmt.Sprintf("%s", ftype)})
				continue
			} else {
				newFields = append(newFields, fInfo)
			}
		}
		fbs.Structs[i].Fields = newFields

	}
	fbs.SearchNotTableMessage()
}

func (fbs *Fbs) SetNameSpace(s string) {
	fbs.NameSpace = s
}

func (fbs *Fbs) SetRootType(s string) {
	fbs.RootType = s
}

func (fbs *Fbs) ExtractStruct(isTable bool) {

	structInfo := StructInfo{}
	structInfo.PkgName = fbs.NameSpace
	structInfo.Name = fbs.typeName
	structInfo.IsTable = isTable
	structInfo.Fields = []structer.FieldInfo{}
	structInfo.IsSlice = map[string]bool{}
	structInfo.IsMessage = map[string]bool{}
	structInfo.IsUnion = map[string]bool{}
	structInfo.IsStruct = map[string]bool{}

	structInfo.Fields = fbs.Fields

	fbs.Structs = append(fbs.Structs, structInfo)

	fbs.Fields = []structer.FieldInfo{}
}

func (fbs *Fbs) SetTypeName(s string) {
	fbs.typeName = s
}

func (fbs *Fbs) FieldName(s string) {
	fbs.fname = toCamelCase(s)
}
func (fbs *Fbs) EnumName(s string) {
	fbs.fname = s
}

func (fbs *Fbs) SetType(s string) {
	fbs.fvalue = s
}

func (fbs *Fbs) NewExtractField() {
	if fbs.Fields == nil {
		fbs.Fields = []structer.FieldInfo{}
	}
	fbs.Fields = append(fbs.Fields, structer.FieldInfo{Name: fbs.fname, Type: fbs.fvalue})
}

func (fbs *Fbs) NewExtractFieldWithValue() {
	fbs.Fields = append(fbs.Fields, structer.FieldInfo{Name: fbs.fname, Type: fbs.fvalue})
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
	for _, f := range fbs.Fields {
		union.Aliases = append(union.Aliases, f.Name)
	}

	fbs.Unions = append(fbs.Unions, union)

	fbs.Fields = []structer.FieldInfo{}
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
