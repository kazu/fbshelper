package info

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"sort"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

type Fbs struct {
	Top     uint32
	VOffset uint32
	VPos    uint32

	SOffset uint32
	SPos    uint32
	SLen    uint32
	VLen    uint16
	TLen    uint16
	VTable  []uint16
	Table   []uint64
	Nest    map[int]uint32
}

type OptionType struct {
	Key string
	Size int
	Nest Option
}

type Option struct {
	Maps []OptionType
}


func GetFbsRootInfo(buf []byte, opt Option) *Fbs {
	return GetFbsInfo(buf, uint32(flatbuffers.GetUOffsetT(buf)), opt)
}

func GetFbsInfo(buf []byte, top uint32, opt Option) *Fbs {

	info := &Fbs{Nest: make(map[int]uint32)}

	info.Top = top
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

	sort.Slice(info.VTable, func(i, j int) bool { return info.VTable[i] < info.VTable[j] })

	info.Table = make([]uint64, 0, info.TLen)
	for idx, _ := range info.VTable {

		diff := info.TLen - info.VTable[idx]

		if idx < len(info.VTable)-1 {
			diff = info.VTable[idx+1] - info.VTable[idx]
		}

		pos := info.Top + uint32(info.VTable[idx])

		switch diff {
		case uint16(1):
			info.Table = append(info.Table, uint64(flatbuffers.GetByte(buf[pos:])))
			fmt.Printf("Table[%d] buf[%d:]=%+v\n", idx, pos, buf[pos:pos+1])
		case uint16(2):
			info.Table = append(info.Table, uint64(flatbuffers.GetVOffsetT(buf[pos:])))
			fmt.Printf("Table[%d] buf[%d:]=%+v\n", idx, pos, buf[pos:pos+2])
		case uint16(4):
			//info.Table = append(info.Table, uint64(flatbuffers.GetUint32(buf[info.Top+info.VTable[idx]:])))
			info.Table = append(info.Table, uint64(flatbuffers.GetUint32(buf[pos:])))
			fmt.Printf("Table[%d] buf[%d:]=%+v\n", idx, pos, buf[pos:pos+4])
		case uint16(8):
			info.Table = append(info.Table, uint64(flatbuffers.GetUint64(buf[pos:])))
			fmt.Printf("Table[%d] buf[%d:]=%+v\n", idx, pos, buf[pos:pos+8])
			//info.Table = append(info.Table, uint64(flatbuffers.GetUint64(buf[info.Top+info.VTable[idx]:])))
		//case uint16(12):

		default:
			fmt.Printf("unknow data: buf[%d:%d]=%+v\n", pos, pos+uint32(diff), buf[pos:pos+uint32(diff)])
			//info.Table = append(info.Table, uint64(flatbuffers.GetUint32(buf[pos:])+info.Top))
			if uint32(diff) == flatbuffers.GetUint32(buf[pos:])+4 {
				fmt.Printf("string data: buf[%d:%d]=%s\n", pos+4, pos+4+uint32(diff), string(buf[pos+4:pos+uint32(diff)]))
			}

			if diff == uint16(5) {
				val := uint32(flatbuffers.GetUint32(buf[pos:]) + pos)
				info.Nest[idx] = val
				info.Table = append(info.Table, 0)
			} else {
				info.Table = append(info.Table, uint64(flatbuffers.GetUint32(buf[pos:])+info.Top))
			}
		}
	}

	sort.Slice(info.VTable, func(i, j int) bool { return i > j })
	sort.Slice(info.Table, func(i, j int) bool { return i > j })

	fmt.Printf("dump VTable\t= %+v\n", buf[info.VPos:info.VPos+uint32(info.VLen)])
	fmt.Printf("dump Table\t= %+v\n", buf[info.Top:info.Top+uint32(info.TLen)])

	for _, pos := range info.Nest {
		fmt.Printf("---Nested--%d---\n", pos)
		ninfo := GetFbsInfo(buf, pos, opt)

		spew.Dump(ninfo)
	}

	return info
}
