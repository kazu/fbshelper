package info

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	flatbuffers "github.com/google/flatbuffers/go"
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
	FIELD_VECTOR
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

	Vector map[int]bool // ?
}

type OptionType struct {
	Key      string
	Size     int
	IsVector bool
	Nest     Option
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
		if otype.IsVector {
			return FIELD_VECTOR
		}
		if len(otype.Nest.Maps) > 0 {
			return FIELD_NEST
		}
	}
	return FIELD_STRING
}

func GetFbsRootInfo(buf []byte, opt Option) *Fbs {
	return GetFbsInfo(buf, uint32(flatbuffers.GetUOffsetT(buf)), opt)
}

func GetFbsInfo(buf []byte, top uint32, opt Option) *Fbs {

	info := &Fbs{Nest: make(map[int]uint32)}
	info.Top = top
	info.Length = top

	info.FetchVtable(buf, opt)
	info.FetchTable(buf, opt)

	return info
}

func MaxLen(x, y uint32) uint32 {
	if y > x {
		return y
	}
	return x
}

func (info *Fbs) FetchVtable(buf []byte, opt Option) error {

	if len(info.VTable) > 0 {
		return nil
	}

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

	return nil
}

func (info *Fbs) FetchTable(buf []byte, opt Option) error {

	if len(info.Table) > 0 {
		return nil
	}

	info.Table = make([]uint64, 0, info.TLen)
	for idx, _ := range info.VTable {

		pos := info.Top + uint32(info.VTable[idx])

		switch GetFieldType(opt.Maps[idx]) {
		case FIELD_ONE:
			info.Length = MaxLen(info.Length, pos+1)
		case FIELD_TWO:
			info.Length = MaxLen(info.Length, pos+2)
		case FIELD_FOUR:
			info.Length = MaxLen(info.Length, pos+4)
		case FIELD_EIGHT:
			info.Length = MaxLen(info.Length, pos+8)
		case FIELD_VECTOR:
			vLenOff := flatbuffers.GetUint32(buf[pos:])
			vLen := flatbuffers.GetUint32(buf[pos+vLenOff:])
			start := pos + vLenOff + flatbuffers.SizeUOffsetT

			// get last value of vector
			ptrLast := start + (vLen-1)*4
			vBuf := buf[start : start+64]
			_ = vBuf

			last := ptrLast + flatbuffers.GetUint32(buf[ptrLast:])
			info.Nest[idx] = last

		case FIELD_NEST:
			val := uint32(flatbuffers.GetUint32(buf[pos:]) + pos)
			info.Nest[idx] = val
		case FIELD_STRING:
			//sLen := flatbuffers.GetUint32(buf[pos:])
			sLenOff := flatbuffers.GetUint32(buf[pos:])
			sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])
			//align := (sLen + 1) % flatbuffers.SizeUOffsetT

			//align++

			start := pos + sLenOff + flatbuffers.SizeUOffsetT
			//align := uint32(stringAlign(int(sLen) + 1))

			info.Table = append(info.Table, uint64(sLen))
			info.Length = MaxLen(info.Length, start+sLen)

		default:
			log.Printf("W: unknow data: buf[%d:]=%+v\n", pos, buf[pos:])
		}
	}

	for i, pos := range info.Nest {
		ninfo := &Fbs{Nest: make(map[int]uint32)}
		ninfo.Top = pos
		ninfo.Length = info.Length

		ninfo.FetchVtable(buf, opt.Maps[i].Nest)
		ninfo.FetchTable(buf, opt.Maps[i].Nest)

		info.Length = MaxLen(info.Length, ninfo.Length)
		if info.Childs == nil {
			info.Childs = make(map[int]*Fbs)
		}
		info.Childs[i] = ninfo
	}

	return nil

}

func stringAlign(size int) int {

	alignSize := (^0 + size) + 1
	alignSize &= (flatbuffers.SizeUOffsetT - 1)
	return alignSize
}

func (info *Fbs) FetchAll(buf []byte, opt Option) error {

	if len(info.Table) > 0 {
		return nil
	}

	info.Table = make([]uint64, 0, info.TLen)
	for idx, _ := range info.VTable {

		pos := info.Top + uint32(info.VTable[idx])

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
		case FIELD_VECTOR:
			vLenOff := flatbuffers.GetUint32(buf[pos:])
			vLen := flatbuffers.GetUint32(buf[pos+vLenOff:])
			start := pos + vLenOff + flatbuffers.SizeUOffsetT

			ptrLast := start + (vLen-1)*4
			last := ptrLast + flatbuffers.GetUint32(buf[ptrLast:])
			info.Nest[idx] = last

		case FIELD_NEST:
			val := uint32(flatbuffers.GetUint32(buf[pos:]) + pos)
			info.Nest[idx] = val
			info.Table = append(info.Table, 0)
		case FIELD_STRING:
			bb := buf[pos:]
			_ = bb
			sLenOff := flatbuffers.GetUint32(buf[pos:])
			sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])

			start := pos + sLenOff + flatbuffers.SizeUOffsetT
			//align := (sLen + 1) % flatbuffers.SizeUOffsetT
			//align++

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

	info.Childs = make(map[int]*Fbs)
	for i, pos := range info.Nest {
		ninfo := &Fbs{Nest: make(map[int]uint32)}
		ninfo.Top = pos
		ninfo.Length = info.Length

		ninfo.FetchVtable(buf, opt.Maps[i].Nest)
		ninfo.FetchAll(buf, opt.Maps[i].Nest)
		Debugf("---Nested--%d---\n", pos)
		Debugf("nested info %s\n", spew.Sdump(ninfo))
		Debugf("---Nested----\n")

		info.Length = MaxLen(info.Length, ninfo.Length)
		if info.Childs == nil {
			info.Childs = make(map[int]*Fbs)
		}
		info.Childs[i] = ninfo

	}

	return nil
}
