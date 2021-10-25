package base

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
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
	ERR_MUST_POINTER       error = log.ERR_MUST_POINTER
	ERR_INVALID_TYPE       error = log.ERR_INVALID_TYPE
	ERR_NOT_FOUND          error = log.ERR_NOT_FOUND
	ERR_READ_BUFFER        error = log.ERR_READ_BUFFER
	ERR_MORE_BUFFER        error = log.ERR_MORE_BUFFER
	ERR_NO_SUPPORT         error = log.ERR_NO_SUPPORT
	ERR_INVALID_INDEX      error = log.ERR_INVALID_INDEX
	ERR_NO_INCLUDE_ROOT    error = log.ERR_NO_INCLUDE_ROOT
	ERR_INVLIAD_WRITE_SIZE error = log.ERR_INVLIAD_WRITE_SIZE
	ERR_OVERWARP_DIFF      error = errors.New("old diff is overwarap current dfff")
)

type LogLevel = log.LogLevel

var CurrentLogLevel LogLevel
var LogW io.Writer = os.Stderr

const (
	LOG_ERROR LogLevel = iota
	LOG_WARN
	LOG_DEBUG
)

const (
	BASE_IMPL = iota
	BASE_NO_LAYER
	BASE_DOUBLE_LAYER
	BASE_DIRECT_READER
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

type ParamIO struct {
	req int
}

type OptIO func(*ParamIO)

func (p *ParamIO) merge(opts ...OptIO) {

	for _, optFn := range opts {
		optFn(p)
	}
}

func Size(require int) OptIO {

	return func(p *ParamIO) {
		p.req = require
	}

}

func NewParamOpt(opts ...OptIO) (p *ParamIO) {

	p = &ParamIO{req: 0}
	p.merge(opts...)

	return
}

// IO ... low level buffer
type IO interface {
	Next(skip int) IO
	HasIoReader() bool
	R(off int, opts ...OptIO) []byte
	D(off, size int) *Diff
	C(off, size int, src []byte) error
	Copy(src IO, srcOff, size, dstOff, extend int)
	U(off, size int) []byte
	LenBuf() int
	Merge() // Deprecated: should use Flatten()
	Flatten()
	Dedup()
	ClearValueInfoOnDirty(node *NodeList)
	insertBuf(pos, size int) IO
	insertSpace(pos, size int, isCreate bool) IO
	AddDirty(Dirty)
	GetDiffs() []Diff
	SetDiffs([]Diff)
	ShouldCheckBound() bool
	New(IO) IO
	NewFromBytes([]byte) IO
	Impl() *BaseImpl
	Type() uint8
	Dump(int, ...DumpOptFn) string
	Dup() IO
	// for io interface
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
	Write(p []byte) (n int, err error)
	WriteAt([]byte, int64) (int, error)
}

// BaseImpl ... Base Object of byte buffer for flatbuffers
// read from r and store bytes.
// Diffs has jounals for writing
type BaseImpl struct {
	r       ReaderWithAt
	bytes   []byte
	RDiffs  []Diff
	Diffs   []Diff
	dirties []Dirty

	// for io interface
	seekCur int
}

// Dirty ... Dirty
// Deprecated: not use.
type Dirty struct {
	Pos int
}

// Diff is journal for create/update
type Diff struct {
	Offset int
	bytes  []byte
}

// NewDiff .. create Diff from outside package
func NewDiff(o int, data []byte) Diff {
	return Diff{Offset: o, bytes: data}
}

// Include .. check pos with diff range.
func (d Diff) Include(pos int) bool {
	return d.Offset <= pos && pos <= d.Offset+len(d.bytes)-1
}

func (d Diff) Included(pos, size int) bool {
	return pos <= d.Offset && d.Offset <= pos+size-1
}

// Inner undocumented
func (d Diff) Inner(pos, size int) bool {
	if !d.Include(pos) {
		return false
	}
	return pos+size <= d.Offset+len(d.bytes)
}

func (d Diff) Innerd(pos, size int) bool {
	if pos > d.Offset {
		return false
	}

	if d.Offset+len(d.bytes) < pos+size {
		return true
	}
	return false
}

func (d *Diff) Merge(s *Diff) {

	var nbytes []byte
	nlen := -1

	if d.Inner(s.Offset, len(s.bytes)) {
		goto COPY
	}

	if d.Include(s.Offset) && d.Offset+cap(d.bytes) > s.Offset+len(s.bytes) {
		goto COPY
	}
	nlen = MaxInt(d.Offset+len(d.bytes), s.Offset+len(s.bytes)) - MinInt(d.Offset, s.Offset)
	nbytes = make([]byte, nlen, nlen*2)

	if d.Include(s.Offset) {
		copy(nbytes, d.bytes)
		d.bytes = nbytes
		goto COPY
	}

	if d.Offset > s.Offset {
		copy(nbytes[d.Offset-s.Offset:], d.bytes)
		copy(nbytes, s.bytes)
		d.Offset = s.Offset
		return
	}

	// d.Offset + len < s.Offset
	copy(nbytes, d.bytes)
	copy(nbytes[s.Offset-d.Offset:], s.bytes)
	d.bytes = nbytes
	return

COPY:
	copy(d.bytes[s.Offset-d.Offset:], s.bytes)

	return
}

// Equal ... return true if equal
func (d *Diff) Equal(s *Diff) bool {
	if d.Offset != s.Offset {
		return false
	}

	if !bytes.Equal(d.bytes, s.bytes) {
		return false
	}
	return true
}

// Root manage top node of data.
type Root struct {
	*Node
}

// IsMatchBit ... is used bit field . mainly Field Type
func IsMatchBit(i, j int) bool {
	if (i & j) > 0 {
		return true
	}
	return false
}

type DumpOptFn func(opt *DumpOpt)
type DumpOpt struct {
	size int
	out  io.Writer
}

func (opt *DumpOpt) merge(fns ...DumpOptFn) {

	for _, fn := range fns {
		fn(opt)
	}
}

func OptDumpSize(size int) DumpOptFn {

	return func(opt *DumpOpt) {
		opt.size = size
	}
}

func OptDumpOut(w io.Writer) DumpOptFn {

	return func(opt *DumpOpt) {
		opt.out = w
	}
}

func (b *BaseImpl) Dump(pos int, opts ...DumpOptFn) (out string) {

	opt := DumpOpt{size: 0}
	opt.merge(opts...)
	var builder strings.Builder

	if opt.out == nil {
		opt.out = &builder
		defer func() { out = builder.String() }()

	}
	w := opt.out

	stdoutDumper := hex.Dumper(w)

	if opt.size == 0 || len(b.R(pos, Size(opt.size))) > opt.size {
		stdoutDumper.Write(b.R(pos, Size(opt.size)))
		return
	}

	d := &BaseImpl{
		r:       b.r,
		bytes:   append([]byte{}, b.bytes...),
		Diffs:   append([]Diff{}, b.Diffs...),
		seekCur: 0,
	}
	d.Flatten()
	stdoutDumper.Write(d.R(pos))
	return
}

func (b *BaseImpl) Read(p []byte) (n int, e error) {

	n, e = b.ReadAt(p, int64(b.seekCur))
	b.seekCur += n
	return
}

func (b *BaseImpl) ReadAt(p []byte, i int64) (int, error) {

	bytes := b.R(int(i), Size(len(p)))
	l := len(bytes)

	copy(p, b.R(int(i), Size(l))[:l])

	return l, nil
}

func (b *BaseImpl) Write(p []byte) (n int, err error) {

	n, err = b.WriteAt(p, int64(b.seekCur))
	b.seekCur += n
	return
}

func (b *BaseImpl) WriteAt(p []byte, oi int64) (n int, err error) {
	i := int(oi)

	diff := Diff{Offset: i, bytes: make([]byte, len(p))}
	copy(diff.bytes, p)
	n = len(diff.bytes)

	b.Diffs = append(b.Diffs, diff)
	return
}

// NewBaseImpl ... initialize BaseImpl struct via buffer(buf)
func NewBaseImpl(buf []byte) *BaseImpl {
	return &BaseImpl{bytes: buf, seekCur: 0}
}

// NewBaseImplByIO ... return new BaseImpl instance with io.Reader
func NewBaseImplByIO(rio io.Reader, cap int) *BaseImpl {
	b := &BaseImpl{bytes: make([]byte, 0, cap), seekCur: 0}
	if origReader, ok := rio.(ReaderAt); ok {
		newReader := NewReaderAt(rio)
		newReader.orig = origReader.orig
		newReader.offFromOrig = origReader.offFromOrig
		b.r = newReader
		return b
	}
	b.r = NewReaderAt(rio)

	return b
}

// Dup ... return copied Base.
func (b *BaseImpl) Dup() (dst IO) {

	dst = NewBaseImplByIO(b.r, DEFAULT_BUF_CAP)
	dbytes := make([]byte, len(b.bytes))
	copy(dbytes, b.bytes)

	diffs := make([]Diff, 0, cap(b.RDiffs)+cap(b.GetDiffs()))

	for _, diff := range b.GetDiffs() {
		dbytes := make([]byte, len(diff.bytes), cap(diff.bytes))
		copy(dbytes, diff.bytes)
		diffs = append(diffs, Diff{Offset: diff.Offset, bytes: dbytes})
	}
	dst.SetDiffs(diffs)

	rdiffs := make([]Diff, 0, cap(b.RDiffs))

	for _, diff := range b.RDiffs {
		dbytes := make([]byte, len(diff.bytes), cap(diff.bytes))
		copy(dbytes, diff.bytes)
		diffs = append(rdiffs, Diff{Offset: diff.Offset, bytes: dbytes})
	}
	dst.Impl().RDiffs = rdiffs

	return dst
}

func (b *BaseImpl) Type() uint8 { return BASE_IMPL }

func (dst *BaseImpl) overwrite(src *BaseImpl) {
	dst.r = src.r
	dst.bytes = src.bytes
	dst.RDiffs = src.RDiffs
	dst.Diffs = src.Diffs
	dst.dirties = src.dirties
}

// NewFromBytes ... return new BaseImpl instance with byte buffer
func (b *BaseImpl) NewFromBytes(bytes []byte) IO {
	return NewBaseImpl(bytes)
}

// New ... return new Base Interface (instance is BaseImpl)
func (b *BaseImpl) New(n IO) IO {
	return n
}

// BaseImpl undocumented
func (b *BaseImpl) Impl() *BaseImpl {
	return b
}

// Bytes undocumented
func (b *BaseImpl) Bytes() []byte {
	return b.bytes
}

// Next ... provide next root flatbuffers
// this is mainly for streaming data.
func (b *BaseImpl) Next(skip int) IO {
	newBase := &BaseImpl{
		r:       b.r,
		bytes:   b.bytes,
		seekCur: 0,
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
func (b *BaseImpl) R(off int, opts ...OptIO) (result []byte) {

	result = b.writerR(off, opts...)
	if len(result) == 0 {
		result = b.readerR(off, opts...)
	}
	return
}

func (b *BaseImpl) writerR(off int, opts ...OptIO) []byte {
	//	param := NewParamOpt(opts...)

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
		return nil
	}

	if len(b.Diffs) == 0 {
		return b.bytes[off:]
	}

	minoff := Diffs(b.Diffs).minOffset(
		func(i int) bool {

			return b.Diffs[i].Offset > off
		},
	)
	if len(b.Diffs) == 0 || minoff < 0 {
		return b.bytes[off:]
	}
	len := MinInt(off+len(b.bytes[off:]), minoff)

	return b.bytes[off:len]
}

func (b *BaseImpl) parent() *BaseImpl {

	if b.r == nil {
		return nil
	}
	reader, ok := b.r.(ReaderAt)
	if !ok {
		return nil
	}
	if reader.orig == nil {
		return nil
	}

	return reader.orig
}

func (b *BaseImpl) readerR(off int, opts ...OptIO) []byte {

	if parent := b.parent(); parent != nil {
		return parent.R(off, opts...)
	}

	p := NewParamOpt(opts...)

RETRY:
	result := Diffs(b.RDiffs).findByIndex(loncha.LastIndexOf(b.RDiffs, func(i int) bool {
		return b.RDiffs[i].Offset <= off && off+p.req <= (b.RDiffs[i].Offset+len(b.RDiffs[i].bytes))
	}))

	if result.diff != nil {
		return result.diff.bytes[(off - result.diff.Offset):]
	}

	if off+p.req > len(b.bytes) {
		err := b.loadBuf(off, p.req)
		if off+MaxInt(1, p.req) <= len(b.bytes) {
			goto RETURN_BYTES
		}
		if err != nil {
			return nil
		}
		goto REGET_DIFF

	}
RETURN_BYTES:

	return b.bytes[off:]

REGET_DIFF:
	result = Diffs(b.RDiffs).findByIndex(loncha.LastIndexOf(b.RDiffs, func(i int) bool {
		return b.RDiffs[i].Offset <= off && b.RDiffs[i].Offset+len(b.RDiffs[i].bytes) > off
	}))
	if result.diff == nil {
		panic("infinit loop")
		goto RETRY
	}
	return result.diff.bytes[(off - result.diff.Offset):]

}

// C ... copy buffer
// Deprecated: should use Copy
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

// Copy ... Copy buffer from src to b as BaseImplt
func (b *BaseImpl) Copy(src IO, srcOff, size, dstOff, extend int) {

	if len(b.bytes) > dstOff {
		diff := Diff{Offset: dstOff, bytes: b.bytes[dstOff:]}
		b.Diffs = append([]Diff{diff}, b.Diffs...)
	}

	for i, diff := range b.Diffs {
		if diff.Offset >= dstOff {
			diff.Offset += extend
		}
		b.Diffs[i] = diff
	}

	for srcPtr := srcOff; srcPtr < srcOff+size; {
		data := src.R(srcPtr, Size(size-(srcPtr-srcOff)))
		if len(data) == 0 {
			panic("hoge")
		}

		b.Diffs = append(b.Diffs, Diff{Offset: (srcPtr - srcOff) + dstOff, bytes: data})
		srcPtr += len(data)
	}
	return
}

// D ... return new Diff of buffer for updating
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

// U ... return buffer for updating
func (b *BaseImpl) U(off, size int) []byte {

	diff := b.D(off, size)
	diff.bytes = make([]byte, size)
	return diff.bytes
}

func calcBufSize(req int) int {

	i := req / 4096
	if i%4096 != 0 {
		i++
	}
	if i == 0 {
		i++
	}

	return i * 4096

}

func (b *BaseImpl) loadBuf(offset, size int) error {
	if !b.HasIoReader() {
		return ERR_READ_BUFFER
	}

	if !b.HasIoReader() && cap(b.bytes)-len(b.bytes) < size {
		return ERR_READ_BUFFER
	}

	// expand b.bytes and return
	if cap(b.bytes) >= offset+size && len(b.bytes) < offset+size {
		blen := len(b.bytes)
		b.bytes = b.bytes[:offset+size]
		n, err := b.r.ReadAt(b.bytes[blen:offset+size], int64(blen))
		if n > 0 {
			b.bytes = b.bytes[:blen+n]
		}
		return err

	}

	// not expand b.bytes
	b.bytes = b.bytes[:len(b.bytes):len(b.bytes)]

	// if offset + size innner diff , already return readerR()
	smaller := Diffs(b.RDiffs).findByIndex(
		loncha.LastIndexOf(b.RDiffs, func(i int) bool {
			return b.RDiffs[i].Offset <= offset
		}))

	requireSize := size
	capLimit := 0
	if smaller.diff != nil && smaller.n+1 < len(b.RDiffs) {
		requireSize = MinInt(size, b.RDiffs[smaller.n+1].Offset-smaller.diff.Offset)
		capLimit = b.RDiffs[smaller.n+1].Offset - offset
	}
	if smaller.diff == nil && len(b.RDiffs) > 0 {
		requireSize = MinInt(size, b.RDiffs[0].Offset-offset)
	}

	if smaller.diff != nil && smaller.diff.Offset+cap(smaller.diff.bytes) >= offset+requireSize {
		goto EXPAND_AND_READ_DIFF
	}

	// shrink diff.bytes
	if smaller.diff != nil && smaller.diff.Offset+cap(smaller.diff.bytes) >= offset && smaller.diff.Offset+len(smaller.diff.bytes) <= offset {
		smaller.diff.bytes = smaller.diff.bytes[: len(smaller.diff.bytes) : offset-smaller.diff.Offset]
	}

	// make new diff and load
	if true {
		ncap := calcBufSize(MaxInt(size, DEFAULT_BUF_CAP))
		if size != requireSize {
			ncap = requireSize
		}
		if capLimit > 0 && capLimit < ncap {
			ncap = capLimit
		}
		diff := Diff{Offset: offset, bytes: make([]byte, requireSize, ncap)}
		var n int
		var err error
	LOAD_CAP:
		n, err = b.r.ReadAt(diff.bytes, int64(offset))
		if errMust, ok := err.(ErrorMustRead); ok {
			b.loadBuf(errMust.ToParam())
			goto LOAD_CAP
		}
		if n > 0 {
			diff.bytes = diff.bytes[:n]
		}
		b.AddRDiff(diff)
		return err
	}

EXPAND_AND_READ_DIFF:

	diff := smaller.diff

	// expand diff and load
	start := len(diff.bytes) + diff.Offset
	sizeForRead := (offset + requireSize) - (len(diff.bytes) + diff.Offset)

	diff.bytes = diff.bytes[:len(diff.bytes)+sizeForRead]

	n, e := b.r.ReadAt(diff.bytes[start-diff.Offset:], int64(start))
	if errMust, ok := e.(ErrorMustRead); ok {
		b.loadBuf(errMust.ToParam())
	}

	if n < len(diff.bytes[start-diff.Offset:]) && n > 0 {
		diff.bytes = diff.bytes[:start-diff.Offset+n]
	}
	return e

}

// LenBuf ... size of buffer
func (b *BaseImpl) LenBuf() int {

	if len(b.Diffs) < 1 && len(b.RDiffs) < 1 {
		return len(b.bytes)
	}

	max := len(b.bytes)
	for i := range b.RDiffs {
		if b.RDiffs[i].Offset+len(b.RDiffs[i].bytes) > max {
			max = b.RDiffs[i].Offset + len(b.RDiffs[i].bytes)
		}
	}

	for i := range b.Diffs {
		if b.Diffs[i].Offset+len(b.Diffs[i].bytes) > max {
			max = b.Diffs[i].Offset + len(b.Diffs[i].bytes)
		}
	}

	return max

}

// Merge ... join Diffs buffer to bytes
// Deprecated: Use Flatten()
func (b *BaseImpl) Merge() {
	b.Flatten()
}

// Flatten ... Diffs buffer join to bytes
func (b *BaseImpl) Flatten() {
	b.FlattenWithLen(-1)
}

// FlattenWithLen ... Diffs buffer join to bytes
func (b *BaseImpl) FlattenWithLen(l int) {

	nbytes := make([]byte, b.LenBuf())

	copy(nbytes, b.bytes)
	for i := range b.Diffs {
		if l > 0 && b.Diffs[i].Offset+len(b.Diffs[i].bytes) > l {
			_ = nbytes
		}
		copy(nbytes[b.Diffs[i].Offset:], b.Diffs[i].bytes)
	}
	b.bytes = nbytes
	b.Diffs = []Diff{}

}

// Dedup ... dedup in b.Diffs
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

// ClearValueInfoOnDirty ... notify to be changed ValueInfo.
func (b *BaseImpl) ClearValueInfoOnDirty(node *NodeList) {

	err := loncha.Delete(&b.dirties, func(i int) bool {
		return (b.dirties[i].Pos == node.Node.Pos) ||
			(node.ValueInfo.Pos > 0 && b.dirties[i].Pos == node.ValueInfo.Pos)
	})
	if err != nil {
		return
	}
}

// GetDiffs ... return Diffs of buffer
func (b *BaseImpl) GetDiffs() []Diff {
	return b.Diffs
}

// SetDiffs ... set Diffs of buffer
func (b *BaseImpl) SetDiffs(d []Diff) {
	b.Diffs = d
}

// AddDirty ... add dirty of ValueInfo.
func (b *BaseImpl) AddDirty(d Dirty) {

	b.dirties = append(b.dirties, d)

}
func (b *BaseImpl) insertBuf(pos, size int) IO {
	return b.insertSpace(pos, size, true)
}
func (b *BaseImpl) insertSpace(pos, size int, isCreate bool) IO {

	newBase := &BaseImpl{
		r:       b.r,
		bytes:   b.bytes,
		seekCur: 0,
	}

	newBase.Diffs = make([]Diff, len(b.Diffs), cap(b.Diffs))

	copy(newBase.Diffs, b.Diffs)

	for i, diff := range newBase.Diffs {
		//if diff.Offset < pos && diff.Offset+len(diff.bytes) > pos {
		if diff.Offset < pos && diff.Include(pos) {
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

// ShouldCheckBound .. which check bounding of buffer.
func (b *BaseImpl) ShouldCheckBound() bool {
	return true
}

func (b *BaseImpl) Equal(c *BaseImpl) bool {
	// if len(b.bytes) != len(c.bytes) {
	// 	return false
	// }
	// if cap(b.bytes) != cap(c.bytes) {
	// 	return false
	// }

	// if len(b.bytes) > 0 && &b.bytes[0] != &c.bytes[0] {
	// 	return false
	// }

	if !bytes.Equal(b.bytes, c.bytes) {
		return false
	}

	if len(b.Diffs) != len(c.Diffs) {
		return false
	}

	for i, _ := range b.Diffs {
		if !b.Diffs[i].Equal(&c.Diffs[i]) {
			return false
		}
	}
	return true
}
func (b *BaseImpl) shrink(start, size int) error {

	//bsize := size

	if len(b.bytes) < start {
		b.bytes = []byte{}
		goto SETUP_DIFF
	}

	if start+size < len(b.bytes) {
		b.bytes = b.bytes[:start+size]
	}

	if len(b.bytes) >= start {
		b.bytes = b.bytes[start:]
	}

SETUP_DIFF:

	loncha.Delete(&b.Diffs, func(i int) bool {
		if b.Diffs[i].Include(start) || b.Diffs[i].Include(start+size) || b.Diffs[i].Innerd(start, size) {
			return false
		}
		return true
	})

	for i := range b.Diffs {

		diff := &b.Diffs[i]
		diff.Offset -= start
		if diff.Offset < 0 {
			diff.bytes = diff.bytes[-diff.Offset:]
			diff.Offset = 0
		}
		if diff.Offset+len(diff.bytes) > start+size {

			len := start + size - diff.Offset
			diff.bytes = diff.bytes[:len]
		}
	}

	return nil
}

func (b *BaseImpl) moveLastInDiff(idx int) error {

	if len(b.Diffs)-1 == idx {
		return nil
	}
	if len(b.Diffs) <= idx {
		return errors.New("moveDiff invalid idx")
	}

	swapFn := reflect.Swapper(b.Diffs)

	for i := idx; i < len(b.Diffs)-1; i++ {
		swapFn(i, i+1)
	}

	return nil
}

func (b *BaseImpl) findByIndexResult(diffs []Diff, n int, err error) (diff *Diff) {

	if err != nil {
		return nil
	}

	if n < 0 {
		return nil
	}

	if n >= len(diffs) {
		return nil
	}
	return &diffs[n]
}

type DiffsIndexResult struct {
	n    int
	e    error
	diff *Diff
}

type Diffs []Diff

func (diffs Diffs) findByIndex(n int, e error) (result *DiffsIndexResult) {

	result = &DiffsIndexResult{n: n, e: e, diff: nil}
	if result.e != nil {
		return
	}
	if n < 0 {
		return
	}

	if n >= len(diffs) {
		return
	}

	result.diff = &diffs[n]
	return
}

func (diffs Diffs) minOffset(fn func(i int) bool) (off int) {

	off = -1

	for i := range diffs {
		if !fn(i) {
			continue
		}

		if off < 0 || diffs[i].Offset < off {
			off = diffs[i].Offset
		}
	}

	return off
}

func (b *BaseImpl) AddRDiff(diff Diff) error {

	result := Diffs(b.RDiffs).findByIndex(loncha.LastIndexOf(b.RDiffs, func(i int) bool {
		return b.RDiffs[i].Offset < diff.Offset
	}))

	if result.diff != nil && result.diff.Offset+len(result.diff.bytes) > diff.Offset {
		return ERR_OVERWARP_DIFF
	}

	b.RDiffs = append(b.RDiffs, diff)

	sort.Slice(b.RDiffs, func(i, j int) bool {
		return b.RDiffs[i].Offset < b.RDiffs[j].Offset
	})

	return nil
}
