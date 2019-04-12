package fbsparser

import (
	"github.com/kazu/lonacha/structer"
)

type Fbs struct {
	NameSpace  string
	typeName   string
	Structs    []structer.StructInfo
	Fields     map[string]string
	fname      string
	fvalue     string
	isRepeated bool
}

func (fbs *Fbs) SetNameSpace(s string) {
	fbs.NameSpace = s
}

func (fbs *Fbs) ExtractStruct() {

	structInfo := structer.StructInfo{
		PkgName: fbs.NameSpace,
		Name:    fbs.typeName,
	}
	structInfo.Fields = fbs.Fields

	fbs.Structs = append(fbs.Structs, structInfo)

	fbs.Fields = map[string]string{}
}

func (fbs *Fbs) SetTypeName(s string) {
	fbs.typeName = s
}

func (fbs *Fbs) FieldNaame(s string) {
	fbs.fname = s
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
	fbs.fvalue = "[]" + s
}
