package base

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kazu/loncha"

	log "github.com/kazu/fbshelper/query/log"
)

const (
	// DEFAULT_BUF_CAP .. cap size of Base's buffef([]byte)
	DEFAULT_BUF_CAP = 512
	// DEFAULT_NODE_BUF_CAP ... cap size of Base's buffer via CreateNode
	DEFAULT_NODE_BUF_CAP = 64
)

var (
	// tag for Marshal/Unmarshal
	FBS_TAG_NAME string = "fbs"
)

var (
	ERR_MUST_POINTER    error = log.ERR_MUST_POINTER
	ERR_INVALID_TYPE    error = log.ERR_INVALID_TYPE
	ERR_NOT_FOUND       error = log.ERR_NOT_FOUND
	ERR_READ_BUFFER     error = log.ERR_READ_BUFFER
	ERR_MORE_BUFFER     error = log.ERR_MORE_BUFFER
	ERR_NO_SUPPORT      error = log.ERR_NO_SUPPORT
	ERR_INVALID_INDEX   error = log.ERR_INVALID_INDEX
	ERR_NO_INCLUDE_ROOT error = log.ERR_NO_INCLUDE_ROOT
)

type LogLevel = log.LogLevel

var CurrentLogLevel LogLevel
var LogW io.Writer = os.Stderr

const (
	LOG_ERROR LogLevel = iota
	LOG_WARN
	LOG_DEBUG
)

type LogArgs struct {
	Fmt  string
	Infs []interface{}
}

type LogFn func() LogArgs

func SetLogLevel(l LogLevel) {
	CurrentLogLevel = l
}

func F(s string, v ...interface{}) LogArgs {
	return LogArgs{Fmt: s, Infs: v}
}

// if no output , not eval args
//  Log(LOG_DEBUG, func() LogArgs { return F("test %d \n", 1) })
func Log(l LogLevel, fn LogFn) {

	if CurrentLogLevel < l {
		return
	}

	var b strings.Builder
	switch l {
	case LOG_DEBUG:
		b.WriteString("D: ")
	case LOG_WARN:
		b.WriteString("W: ")
	case LOG_ERROR:
		b.WriteString("E: ")
	default:
		b.WriteString(" ")
	}
	args := fn()
	fmt.Fprintf(&b, args.Fmt, args.Infs...)
	io.WriteString(LogW, b.String())

	return

}

type Base interface {
	Next(skip int) Base
	HasIoReader() bool
	R(off int) []byte
	C(off, size int, src []byte) error
	Copy(src Base, srcOff, size, dstOff, extend int)
	D(off, size int) *Diff
	U(off, size int) []byte
	LenBuf() int
	Merge()
	Dedup()
	ClearValueInfoOnDirty(node *NodeList)
	insertBuf(pos, size int) Base
	insertSpace(pos, size int, isCreate bool) Base
	AddDirty(Dirty)
	GetDiffs() []Diff
	SetDiffs([]Diff)
}

// BaseImpl ... Base Object of byte buffer for flatbuffers
// read from r and store bytes.
// Diffs has jounals for writing
type BaseImpl struct {
	r       io.Reader
	bytes   []byte
	RDiffs  []Diff
	Diffs   []Diff
	dirties []Dirty
}

type Dirty struct {
	Pos int
}

// Diff is journal for create/update
type Diff struct {
	Offset int
	bytes  []byte
}

// Root manage top node of data.
type Root struct {
	*Node
}

func IsMatchBit(i, j int) bool {
	if (i & j) > 0 {
		return true
	}
	return false
}

// NewBase initialize Base struct via buffer(buf)
func NewBase(buf []byte) *BaseImpl {
	return &BaseImpl{bytes: buf}
}

// NewBaseByIO initialize , make Base
//   this dosent use []byte, use io.Reader
func NewBaseByIO(rio io.Reader, cap int) *BaseImpl {
	b := &BaseImpl{r: rio, bytes: make([]byte, 0, cap)}
	return b
}

// NextBase provide next root flatbuffers
// this is mainly for streaming data.
func (b *BaseImpl) Next(skip int) Base {
	newBase := &BaseImpl{
		r:     b.r,
		bytes: b.bytes,
	}
	newBase.bytes = newBase.bytes[skip:]
	if cap(b.bytes) > skip {
		b.bytes = b.bytes[:skip]
	}

	if len(b.bytes) <= skip && len(b.RDiffs) > 0 {
		//FIXME: impelement
		Log(LOG_WARN, func() LogArgs {
			return F("NextBase(): require RDiff. but not implemented\n")
		})
	}

	return newBase
}

// HasIoReader ... Base buffer is reading io.Reader or not.
func (b *BaseImpl) HasIoReader() bool {
	return b.r != nil
}

// R ... is access buffer data
// Base cannot access byte buffer directly.
func (b *BaseImpl) R(off int) []byte {
	if len(b.Diffs) == 0 {
		return b.readerR(off)
	}

	n, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})
	if e == nil && n >= 0 {
		return b.Diffs[n].bytes[(off - b.Diffs[n].Offset):]
	}

	if len(b.bytes) < off {
		Log(LOG_WARN, func() LogArgs {
			return F("base.R(): remain offset=%d lenBuf()=%d \n",
				off, b.LenBuf())
		})
	}

	return b.bytes[off:]
}

func (b *BaseImpl) readerR(off int) []byte {

	n, e := loncha.LastIndexOf(b.RDiffs, func(i int) bool {
		return b.RDiffs[i].Offset <= off && off < (b.RDiffs[i].Offset+len(b.RDiffs[i].bytes))
	})
	if e == nil && n >= 0 {
		return b.RDiffs[n].bytes[(off - b.RDiffs[n].Offset):]
	}

	if off+32 >= len(b.bytes) {
		b.expandBuf(off - len(b.bytes) + 32)
	}

	return b.bytes[off:]
}

func (b *BaseImpl) C(off, size int, src []byte) error {

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})
	if e == nil && sn >= 0 && b.Diffs[sn].Offset == off && len(b.Diffs[sn].bytes) == size {
		b.Diffs[sn].bytes = src[:size]
		return nil
	}

	diff := b.D(off, size)
	diff.bytes = src[:size]
	return nil
}

func (b *BaseImpl) Copy(osrc Base, srcOff, size, dstOff, extend int) {

	src, ok := osrc.(*BaseImpl)
	if !ok {
		log.Log(LOG_WARN, log.Printf("BaseImpl.Copy() src is only BaseImpl"))

		return
	}

	if len(b.bytes) > dstOff {
		diff := Diff{Offset: dstOff, bytes: b.bytes[dstOff:]}
		b.Diffs = append(b.Diffs, diff)
		//b.bytes = b.bytes[:dstOff]
	}

	for i, diff := range b.Diffs {
		if diff.Offset >= dstOff {
			diff.Offset += extend
		}
		b.Diffs[i] = diff
	}

	if len(src.bytes) > srcOff {
		nSize := len(src.bytes[srcOff:])
		if nSize > size {
			nSize = size
		}

		diff := Diff{Offset: dstOff, bytes: src.bytes[srcOff : srcOff+nSize]}
		b.Diffs = append(b.Diffs, diff)
	}
	for _, diff := range src.Diffs {
		if diff.Offset >= srcOff {
			nDiff := diff
			nDiff.Offset -= srcOff
			nDiff.Offset += dstOff
			b.Diffs = append(b.Diffs, nDiff)
		}
	}
	return
}

func (b *BaseImpl) D(off, size int) *Diff {

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})

	if e == nil && sn >= 0 {
		// if b.Diffs[sn].Offset+len(b.Diffs[sn].bytes) < off+size {
		// 	Log(LOG_ERROR, func() LogArgs {
		// 		return F("D(%d,%d) Invalid diff=%d\n",
		// 			off, size, )
		// 	})
		// }else
		if b.Diffs[sn].Offset+len(b.Diffs[sn].bytes) <= off+size {
			diff := Diff{Offset: off}
			b.Diffs = append(b.Diffs, diff)
			idx := len(b.Diffs) - 1
			return &b.Diffs[idx]
		}

		diffbefore := b.Diffs[sn]
		//MENTION: should increase cap ?
		diff := Diff{Offset: off}
		diffafter := Diff{Offset: off + size, bytes: diffbefore.bytes[off-diffbefore.Offset+size:]}
		diffbefore.bytes = diffbefore.bytes[:off+size-diffbefore.Offset]
		b.Diffs[sn] = diffbefore
		b.Diffs = append(b.Diffs[:sn+1],
			append([]Diff{diffafter}, b.Diffs[sn+1:]...)...)
		b.Diffs = append(b.Diffs, diff)
		return &b.Diffs[len(b.Diffs)-1]
	}

	if len(b.bytes) > off && cap(b.bytes) > off+size {
		diff := Diff{Offset: off, bytes: b.bytes[off : off+size]}
		b.Diffs = append(b.Diffs, diff)
		//b.bytes = b.bytes[:off]
	}

	//MENTION: should increase cap ?
	diff := Diff{Offset: off}
	b.Diffs = append(b.Diffs, diff)
	idx := len(b.Diffs) - 1

	return &b.Diffs[idx]
}

func (b *BaseImpl) U(off, size int) []byte {

	diff := b.D(off, size)
	diff.bytes = make([]byte, size)
	return diff.bytes
}

// FIXME:  support base.RDiff
func (b *BaseImpl) expandBuf(plus int) error {
	if !b.HasIoReader() && cap(b.bytes)-len(b.bytes) < plus {
		return nil
	}

	if !b.HasIoReader() {
		return nil
	}
	if cap(b.bytes) < len(b.bytes)+plus {
		diff := Diff{Offset: len(b.bytes), bytes: make([]byte, 0, DEFAULT_BUF_CAP)}
		n, err := io.ReadAtLeast(b.r, diff.bytes, plus)
		if n < plus || err != nil {
			diff.bytes = diff.bytes[:n]
			if n > 0 {
				b.RDiffs = append(b.RDiffs, diff)
			}
			return ERR_READ_BUFFER
		}
		b.RDiffs = append(b.RDiffs, diff)
	}

	l := len(b.bytes)
	b.bytes = b.bytes[:l+plus]

	n, err := io.ReadAtLeast(b.r, b.bytes[l:], plus)
	if n < plus || err != nil {
		b.bytes = b.bytes[:l+n]
		return ERR_READ_BUFFER
	}
	return nil
}

// LenBuf ... size of buffer
func (b *BaseImpl) LenBuf() int {

	if len(b.Diffs) < 1 {
		return len(b.bytes)
	}

	max := len(b.bytes)
	for i := range b.Diffs {
		if b.Diffs[i].Offset+len(b.Diffs[i].bytes) > max {
			max = b.Diffs[i].Offset + len(b.Diffs[i].bytes)
		}
	}

	return max

}
func (b *BaseImpl) Merge() {

	nbytes := make([]byte, b.LenBuf())

	copy(nbytes, b.bytes)
	for i := range b.Diffs {
		copy(nbytes[b.Diffs[i].Offset:], b.Diffs[i].bytes)
	}
	b.bytes = nbytes
	b.Diffs = []Diff{}

}

func (b *BaseImpl) Dedup() {

	loncha.Delete(&b.Diffs, func(i int) bool {
		n, err := loncha.LastIndexOf(b.Diffs, func(j int) bool {
			return i < j &&
				b.Diffs[j].Offset <= b.Diffs[i].Offset &&
				b.Diffs[j].Offset+len(b.Diffs[j].bytes) >= b.Diffs[i].Offset+len(b.Diffs[i].bytes)
		})
		if err == nil && n > 0 {
			ok := true
			_ = ok
		}
		return err == nil && n > 0
	})
}

// RawBufInfo ... capacity/length infomation of Base's buffer
type RawBufInfo struct {
	Len int
	Cap int
}

// BufInfo is buffer detail infomation
// information of Base.bytes is stored to 0 indexed
// information of Base.Diffs is stored to 1 indexed
type BufInfo [2]RawBufInfo

// BufInfo ... return infos(buffer detail information)
func (b *BaseImpl) BufInfo() (infos BufInfo) {

	infos[0].Len = len(b.bytes)
	infos[0].Cap = cap(b.bytes)

	infos[1].Len = infos[0].Len
	infos[1].Cap = infos[0].Cap
	for _, diff := range b.Diffs {
		if len(diff.bytes)+diff.Offset > infos[1].Len {
			infos[1].Len = diff.Offset + len(diff.bytes)
		}
		if cap(diff.bytes)+diff.Offset > infos[1].Cap {
			infos[1].Cap = diff.Offset + cap(diff.bytes)
		}
	}
	return
}

func (b *BaseImpl) ClearValueInfoOnDirty(node *NodeList) {

	err := loncha.Delete(&b.dirties, func(i int) bool {
		return (b.dirties[i].Pos == node.Node.Pos) ||
			(node.ValueInfo.Pos > 0 && b.dirties[i].Pos == node.ValueInfo.Pos)
	})
	if err != nil {
		return
	}
}

func (b *BaseImpl) GetDiffs() []Diff {
	return b.Diffs
}

func (b *BaseImpl) SetDiffs(d []Diff) {
	b.Diffs = d
}

func (b *BaseImpl) AddDirty(d Dirty) {

	b.dirties = append(b.dirties, d)

}
func (b *BaseImpl) insertBuf(pos, size int) Base {
	return b.insertSpace(pos, size, true)
}
func (b *BaseImpl) insertSpace(pos, size int, isCreate bool) Base {

	newBase := &BaseImpl{
		r:     b.r,
		bytes: b.bytes,
	}

	newBase.Diffs = make([]Diff, len(b.Diffs), cap(b.Diffs))

	copy(newBase.Diffs, b.Diffs)

	for i, diff := range newBase.Diffs {
		if diff.Offset < pos && diff.Offset+len(diff.bytes) > pos {
			newBase.Diffs = append(newBase.Diffs[:i],
				append([]Diff{
					Diff{Offset: diff.Offset, bytes: diff.bytes[:pos-diff.Offset]},
					Diff{Offset: pos, bytes: diff.bytes[pos-diff.Offset:]}},
					b.Diffs[i+1:]...)...)
		}
	}

	for i := range newBase.Diffs {
		if newBase.Diffs[i].Offset >= pos {
			newBase.Diffs[i].Offset += size
		}
	}

	if len(newBase.bytes) > pos {
		newBase.Diffs = append(newBase.Diffs,
			Diff{Offset: pos + size, bytes: newBase.bytes[pos:]})
		newBase.bytes = newBase.bytes[:pos]
	}
	if isCreate {
		newBase.Diffs = append(newBase.Diffs,
			Diff{Offset: pos, bytes: make([]byte, size)})
	}
	return newBase
}
