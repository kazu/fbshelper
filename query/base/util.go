package base

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
	Typebyte
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUnt8
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

var NameToType map[string]int = map[string]int{
	"bool":    TypeBool,
	"byte":    Typebyte,
	"int8":    TypeInt8,
	"int16":   TypeInt16,
	"int32":   TypeInt32,
	"int64":   TypeInt64,
	"unt8":    TypeUnt8,
	"uint16":  TypeUint16,
	"uint32":  TypeUint32,
	"uint64":  TypeUint64,
	"float32": TypeFloat32,
	"float64": TypeFloat64,
}

var TypeToGroup map[int]int = map[int]int{
	TypeBool:    FieldTypeBasic1,
	Typebyte:    FieldTypeBasic1,
	TypeInt8:    FieldTypeBasic1,
	TypeInt16:   FieldTypeBasic2,
	TypeInt32:   FieldTypeBasic4,
	TypeInt64:   FieldTypeBasic8,
	TypeUnt8:    FieldTypeBasic1,
	TypeUint16:  FieldTypeBasic2,
	TypeUint32:  FieldTypeBasic4,
	TypeUint64:  FieldTypeBasic8,
	TypeFloat32: FieldTypeBasic4,
	TypeFloat64: FieldTypeBasic8,
}

var TypeToSize map[int]int = map[int]int{
	TypeBool:    1,
	Typebyte:    1,
	TypeInt8:    1,
	TypeInt16:   2,
	TypeInt32:   4,
	TypeInt64:   8,
	TypeUnt8:    1,
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

func NameToTypeEnum(string) int {
	return 0
}

func AtoiNoErr(i int, e error) int {
	if e != nil {
		return -1
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
