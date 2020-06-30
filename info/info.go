package info

import (
	"fmt"
	"log"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/davecgh/go-spew/spew"
)

type FieldType byte

var ENABLE_LOG_DEBUG bool = false

func Debugf(f string, args ...interface{}) {
	if !ENABLE_LOG_DEBUG {
		return
	}
	log.Printf(fmt.Sprintf("D: %s", f), args...)

}

const (
	FIELD_ONE FieldType = iota
	FIELD_TWO
	FIELD_FOUR
	FIELD_EIGHT
	FIELD_STRING
	FIELD_NEST
)

type Fbs struct {
	Top     uint32
	VOffset uint32
	VPos    uint32
	Length  uint32

	SOffset uint32
	SPos    uint32
	SLen    uint32
	VLen    uint16
	TLen    uint16
	VTable  []uint16
	Table   []uint64
	Nest    map[int]uint32
	Childs  map[int]*Fbs
}

type OptionType struct {
	Key  string
	Size int
	Nest Option
}

type Option struct {
	Maps []OptionType
}

func GetFieldType(otype OptionType) FieldType {
	switch otype.Size {
	case 1:
		return FIELD_ONE
	case 2:
		return FIELD_TWO
	case 4:
		return FIELD_FOUR
	case 8:
		return FIELD_EIGHT
	default:
		if len(otype.Nest.Maps) > 0 {
			return FIELD_NEST
		}
	}
	return FIELD_STRING
}

func GetFbsRootInfo(buf []byte, opt Option) *Fbs {
	return GetFbsInfo(buf, uint32(flatbuffers.GetUOffsetT(buf)), opt)
}

func MaxLen(x, y uint32) uint32 {
	if y > x {
		return y
	}
	return x
}

func GetFbsInfo(buf []byte, top uint32, opt Option) *Fbs {

	info := &Fbs{Nest: make(map[int]uint32)}

	info.Top = top
	info.Length = top
	info.VOffset = uint32(flatbuffers.GetUOffsetT(buf[info.Top:]))
	info.VPos = info.Top - info.VOffset
	fieldStart := 4
	_ = fieldStart

	info.SOffset = uint32(flatbuffers.GetUOffsetT(buf[info.Top+4:]))
	if int(info.SOffset) < len(buf) {
		info.SPos = info.Top + 4 + info.SOffset
		info.SLen = uint32(flatbuffers.GetUOffsetT(buf[info.SOffset:]))
	} else {
		info.SOffset = 0
	}
	info.VLen = uint16(flatbuffers.GetVOffsetT(buf[info.VPos:]))
	info.TLen = uint16(flatbuffers.GetVOffsetT(buf[info.VPos+2:]))

	info.VTable = make([]uint16, 0, info.VLen)
	for cur := info.VPos + 4; cur < info.VPos+uint32(info.VLen); cur += 2 {
		info.VTable = append(info.VTable, uint16(flatbuffers.GetVOffsetT(buf[cur:])))
	}

	info.Table = make([]uint64, 0, info.TLen)
	for idx, _ := range info.VTable {

		pos := info.Top + uint32(info.VTable[idx])

		//size := opt.Maps[idx]

		switch GetFieldType(opt.Maps[idx]) {
		case FIELD_ONE:
			info.Table = append(info.Table, uint64(flatbuffers.GetByte(buf[pos:])))
			info.Length = MaxLen(info.Length, pos+1)
			Debugf("Table[%d] buf[%d:]=%+v\n", idx, pos, buf[pos:pos+1])
		case FIELD_TWO:
			info.Table = append(info.Table, uint64(flatbuffers.GetVOffsetT(buf[pos:])))
			info.Length = MaxLen(info.Length, pos+2)
			Debugf("Table[%d] buf[%d:]=%+v\n", idx, pos, buf[pos:pos+2])
		case FIELD_FOUR:
			//info.Table = append(info.Table, uint64(flatbuffers.GetUint32(buf[info.Top+info.VTable[idx]:])))
			info.Table = append(info.Table, uint64(flatbuffers.GetInt32(buf[pos:])))
			info.Length = MaxLen(info.Length, pos+4)
			Debugf("Table[%d] buf[%d:]=%+v\n", idx, pos, buf[pos:pos+4])
		case FIELD_EIGHT:
			info.Table = append(info.Table, uint64(flatbuffers.GetUint64(buf[pos:])))
			info.Length = MaxLen(info.Length, pos+8)
			Debugf("Table[%d] buf[%d:]=%+v\n", idx, pos, buf[pos:pos+8])
		case FIELD_NEST:
			val := uint32(flatbuffers.GetUint32(buf[pos:]) + pos)
			info.Nest[idx] = val
			info.Table = append(info.Table, 0)
		case FIELD_STRING:
			sLen := flatbuffers.GetUint32(buf[pos:])
			start := pos + sLen + flatbuffers.SizeUOffsetT
			info.Table = append(info.Table, uint64(sLen))
			info.Length = MaxLen(info.Length, start+sLen)
			Debugf("Table[%d] buf[%d:%d]='%s'\n", idx, start, start+sLen, buf[start:start+sLen])

		default:
			Debugf("unknow data: buf[%d:]=%+v\n", pos, buf[pos:])
			info.Table = append(info.Table, uint64(flatbuffers.GetUint32(buf[pos:])+info.Top))
		}
	}

	Debugf("dump VTable[%d:]\t= %+v\n", info.VPos, buf[info.VPos:info.VPos+uint32(info.VLen)])
	Debugf("dump Table[%d:]\t= %+v\n", info.Top, buf[info.Top:info.Top+uint32(info.TLen)])

	for i, pos := range info.Nest {
		Debugf("---Nested--%d---\n", pos)
		ninfo := GetFbsInfo(buf, pos, opt.Maps[i].Nest)
		Debugf("nested info %s\n", spew.Sdump(ninfo))
		info.Length = MaxLen(info.Length, ninfo.Length)
		Debugf("---Nested----\n")
		if info.Childs == nil {
			info.Childs = make(map[int]*Fbs)
		}
		info.Childs[i] = ninfo
	}

	return info
}
