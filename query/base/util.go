package base

import "github.com/kazu/loncha"

//. "github.com/kazu/fbshelper/query/error"

// type CondFn func(int) bool
// type ItrFn func(int)

// type FieldIter struct {
// 	Cond []CondFn
// 	Fn   []ItrFn
// }

// func IteratorField(cnt int, conds []CondFn, fnList ...ItrFn) {

// 	for j := 0 ; j < cnt; j++ {
// 		for i := 0; i < len(conds); i++ {
// 			if conds[i](j) {
// 				fnList[i](j)
// 				goto NEXT_J
// 			}
// 		}
// NEXT_J:
// 	}

// }

const (
	SizeOfbool    = 1
	SizeOfint8    = 1
	SizeOfint16   = 2
	SizeOfuint16  = 2
	SizeOfint32   = 4
	SizeOfuint32  = 4
	SizeOfint64   = 8
	SizeOfuint64  = 8
	SizeOffloat32 = 4
	SizeOffloat64 = 8
	SizeOfuint8   = 1
	SizeOfbyte    = 1
)

const (
	TypeBool = iota
	TypeByte
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64
	TypeFloat32
	TypeFloat64
)

const (
	FieldTypeNone = 1 << iota
	FieldTypeStruct
	FieldTypeTable
	FieldTypeUnion
	FieldTypeSlice
	FieldTypeBasic1
	FieldTypeBasic2
	FieldTypeBasic4
	FieldTypeBasic8
)

type Bool bool
type Byte byte
type Int8 int8
type Int16 int16
type Int32 int32
type Int64 int64
type Uint8 uint8
type Uint16 uint16
type Uint32 uint32
type Uint64 uint64
type Float32 float32
type Float64 float64

var NameToType map[string]int = map[string]int{
	"bool":    TypeBool,
	"byte":    TypeByte,
	"int8":    TypeInt8,
	"int16":   TypeInt16,
	"int32":   TypeInt32,
	"int64":   TypeInt64,
	"unt8":    TypeUint8,
	"uint16":  TypeUint16,
	"uint32":  TypeUint32,
	"uint64":  TypeUint64,
	"float32": TypeFloat32,
	"float64": TypeFloat64,
	"Bool":    TypeBool,
	"Byte":    TypeByte,
	"Int8":    TypeInt8,
	"Int16":   TypeInt16,
	"Int32":   TypeInt32,
	"Int64":   TypeInt64,
	"Unt8":    TypeUint8,
	"Uint16":  TypeUint16,
	"Uint32":  TypeUint32,
	"Uint64":  TypeUint64,
	"Float32": TypeFloat32,
	"Float64": TypeFloat64,
}

var TypeToGroup map[int]int = map[int]int{
	TypeBool:    FieldTypeBasic1,
	TypeByte:    FieldTypeBasic1,
	TypeInt8:    FieldTypeBasic1,
	TypeInt16:   FieldTypeBasic2,
	TypeInt32:   FieldTypeBasic4,
	TypeInt64:   FieldTypeBasic8,
	TypeUint8:   FieldTypeBasic1,
	TypeUint16:  FieldTypeBasic2,
	TypeUint32:  FieldTypeBasic4,
	TypeUint64:  FieldTypeBasic8,
	TypeFloat32: FieldTypeBasic4,
	TypeFloat64: FieldTypeBasic8,
}

var TypeToSize map[int]int = map[int]int{
	TypeBool:    1,
	TypeByte:    1,
	TypeInt8:    1,
	TypeInt16:   2,
	TypeInt32:   4,
	TypeInt64:   8,
	TypeUint8:   1,
	TypeUint16:  2,
	TypeUint32:  4,
	TypeUint64:  8,
	TypeFloat32: 4,
	TypeFloat64: 8,
}

var IsStructName map[string]bool = map[string]bool{}

func SetNameIsStrunct(name string, enable bool) bool {
	if enable {
		IsStructName[name] = true
		return true
	}
	return false
}

func NameToTypeEnum(s string) int {
	if v, ok := NameToType[s]; ok {
		return v
	}
	return -1
}

func AtoiNoErr(i int, e error) int {
	if e != nil {
		Log(LOG_WARN, func() LogArgs {
			return F("AtoiNoErr(%d, e=%v) has error\n", i, e)
		})
		return 0
	}
	return i
}

const FieldTypeBasic = FieldTypeBasic1 | FieldTypeBasic2 | FieldTypeBasic4 | FieldTypeBasic8

func IsFieldStruct(i int) bool {
	return IsMatchBit(i, FieldTypeStruct)
}

func IsFieldUnion(i int) bool {
	return IsMatchBit(i, FieldTypeUnion)
}

func IsFieldBytes(i int) bool {
	return IsMatchBit(i, FieldTypeSlice) && IsMatchBit(i, FieldTypeBasic1)
}
func IsFieldSlice(i int) bool {
	return IsMatchBit(i, FieldTypeSlice)
}

func IsFieldTable(i int) bool {
	return IsMatchBit(i, FieldTypeTable)
}

func IsFieldBasicType(i int) bool {
	return IsMatchBit(i, FieldTypeBasic)
}

var UnionAlias map[string][]string = map[string][]string{}

func SetAlias(union, alias string) bool {

	if UnionAlias[union] == nil {
		UnionAlias[union] = []string{}
	}
	UnionAlias[union] = append(UnionAlias[union], alias)
	return true

}

func IsUnionName(s string) bool {
	_, ok := UnionAlias[s]
	return ok
}

func ToBool(s string) bool {
	return s == "true"
}

var All_IdxToType map[string](map[int]int) = map[string](map[int]int){}
var All_IdxToTypeGroup map[string](map[int]int) = map[string](map[int]int){}
var All_IdxToName map[string](map[int]string) = map[string](map[int]string){}
var All_NameToIdx map[string](map[string]int) = map[string](map[string]int){}

var IdxNames []string = []string{}

var Name2Idx map[string]int = map[string]int{}

func AddName(s string) {

	if _, ok := Name2Idx[s]; ok {
		return
	}

	IdxNames = append(IdxNames, s)
	idx, err := loncha.IndexOf(IdxNames, func(i int) bool { return IdxNames[i] == s })
	if err != nil {
		Log(LOG_WARN, func() LogArgs {
			return F("IdxNames: fail to add %s ", s)
		})
		return
	}
	Name2Idx[s] = idx
	return
}

func Symbol(s string) int {
	v, ok := Name2Idx[s]
	if !ok {
		return -1
	}
	return v
}

func SymbolToString(i int) string {
	if i >= len(IdxNames) {
		return ""
	}
	return IdxNames[i]
}

func NodeFieldSynmbol(node, field string) (uint32, error) {
	nidx := uint32(Symbol(node))
	fidx := uint32(Symbol(field))
	if nidx < 0 || fidx < 0 {
		return 0, ERR_INVALID_INDEX
	}
	nidx = nidx << 16
	nidx += fidx
	return nidx, nil
}

func SymbolToNames(idx uint32) (node, field string, e error) {

	nidx := idx >> 16
	fidx := idx - ((idx >> 16) << 16)

	node = SymbolToString(int(nidx))
	if node == "" {
		e = ERR_INVALID_INDEX
		return
	}
	field = SymbolToString(int(fidx))
	if field == "" {
		e = ERR_INVALID_INDEX
		return
	}

	return
}

func IsSameTable(i, j uint32) bool {
	return (i >> 16) == (j >> 16)
}

func IsSameField(i, j uint32) bool {
	return (i - ((i >> 16) << 16)) == (j - ((j >> 16) << 16))
}

func SetNameToIdx(name string, v map[string]int) {
	All_NameToIdx[name] = v
}

func SetIdxToName(name string, v map[int]string) {
	All_IdxToName[name] = v
}

func SetIdxToType(name string, v map[int]int) {
	All_IdxToType[name] = v
}
func SetdxToTypeGroup(name string, v map[int]int) {
	All_IdxToTypeGroup[name] = v
}

func GetTypeGroup(s string) (result int) {

	result = 0
	if enum, ok := NameToType[s]; ok {
		return TypeToGroup[enum]
	}

	if _, ok := UnionAlias[s]; ok {
		return FieldTypeUnion
	}
	if s[0:6] == "[]byte" {
		return FieldTypeSlice | FieldTypeBasic1
	}

	if s[0:2] == "[]" {
		result |= FieldTypeSlice
	}

	if enum, ok := NameToType[s[2:]]; ok {
		result |= TypeToGroup[enum]
		return
	}

	result |= FieldTypeTable
	return

}
