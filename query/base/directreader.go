package base

import (
	"io"

	log "github.com/kazu/fbshelper/query/log"
	"github.com/kazu/loncha"
)

// DirectReader ... is Base without buffer in data area
type DirectReader struct {
	IO
	r io.ReaderAt
}

// NewDirectReader ... return new DirectReader instance.
func NewDirectReader(b IO, r io.ReaderAt) DirectReader {
	return DirectReader{IO: b, r: r}
}

// Type ... type of Base
func (b DirectReader) Type() uint8 { return BASE_DIRECT_READER }

// R ... read data.
func (b DirectReader) R(offset int, opts ...OptIO) (result []byte) {
	if b.IO.LenBuf() > offset {
		return b.IO.R(offset, opts...)
	}
	bLen := b.IO.LenBuf()
	result = make([]byte, 128)
	n, err := b.r.ReadAt(result, int64(offset-bLen))
	result = result[:n]
	if err != nil {
		//Log(LOG_WARN, "P(%d) cannot read \n", offset)
		if n > 0 {
			return
		}
		log.Log(LOG_WARN, log.Printf("P(%d) cannot read \n", offset))
		return nil
	}
	return
}

// ShouldCheckBound ... DirectReader is not checking boundary.
func (b DirectReader) ShouldCheckBound() bool {
	return false
}

// NoLayer is BaseImpl without journal infomation for perfomance.
type NoLayer struct {
	*BaseImpl
}

// NewNoLayer ... return new NoLayer instnace. b must be BaseImpl/NoLayer
func NewNoLayer(b IO) NoLayer {
	impl, ok := b.(*BaseImpl)
	if ok {
		return NoLayer{BaseImpl: impl}
	}
	if _, already := b.(NoLayer); already {
		return b.(NoLayer)
	}
	if b.Impl() != nil {
		return NoLayer{BaseImpl: b.Impl()}
	}

	return NoLayer{}
}

// FIXME: should impleent WriteAt for performance
// func (b NoLayer) WriteAt(p []byte, oi int64) (n int, err error) {
//}

// Dup ... return copied Base
func (b NoLayer) Dup() (dst IO) {

	return NoLayer{BaseImpl: b.BaseImpl.Dup().Impl()}
}

// Type ... type of Base interface
func (b NoLayer) Type() uint8 { return BASE_NO_LAYER }

func (b NoLayer) insertBuf(pos, size int) IO {
	return b.insertSpace(pos, size, true)
}

func (b NoLayer) insertSpace(pos, size int, isCreate bool) IO {
	nImple := b.BaseImpl.insertSpace(pos, size, isCreate).(*BaseImpl)

	if nImple.LenBuf() >= pos+size {
		//b.BaseImpl = nImple
		return NoLayer{BaseImpl: nImple}
	} else {
		//b.BaseImpl = b.BaseImpl.insertSpace(pos, size, true).(*BaseImpl)
		return NoLayer{BaseImpl: b.BaseImpl.insertSpace(pos, size, true).(*BaseImpl)}
	}

	return b
}

// R ... read buffer
func (b NoLayer) R(off int, opts ...OptIO) []byte {
	p := NewParamOpt(opts...)

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off+MinInt(1, p.req) <= (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})
	if e == nil && sn >= 0 {
		return b.Diffs[sn].bytes[off-b.Diffs[sn].Offset:]
	}

	return b.Impl().R(off, opts...)

}

// D ... return Diff for write
func (b NoLayer) D(off, size int) *Diff {

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+cap(b.Diffs[i].bytes))
	})

	if e == nil && sn >= 0 {
		if b.Diffs[sn].Offset == off {
			if e := b.Impl().moveLastInDiff(sn); e != nil {
				panic(e.Error())
			}
			return &b.Diffs[len(b.Diffs)-1]
		}
		diff := b.Diffs[sn]
		offDiff := off - diff.Offset
		b.Diffs[sn].bytes = b.Diffs[sn].bytes[:offDiff+size]
		return &Diff{Offset: off, bytes: b.Diffs[sn].bytes[offDiff : offDiff+size]}
	}

	if off+size <= len(b.bytes) {
		return &Diff{Offset: off, bytes: b.bytes[off:]}
	}

	if off+size <= cap(b.bytes) {
		b.bytes = b.bytes[:off+size]

		return &Diff{Offset: off, bytes: b.bytes[off:]}
	}

	newDiff := Diff{Offset: off, bytes: make([]byte, 0, 512)}
	newDiff.bytes = newDiff.bytes[:size]
	b.Diffs = append(b.Diffs, newDiff)

	return &newDiff
}

// U ... return buffer for update
func (b NoLayer) U(off, size int) []byte {

	diff := b.D(off, size)
	if cap(diff.bytes) < size {
		diff.bytes = make([]byte, size)
	}
	if len(diff.bytes) < size {
		diff.bytes = diff.bytes[:size:size]
	}
	return diff.bytes
}

// Copy ... copy buffer from Base
func (b NoLayer) Copy(osrc IO, srcOff, size, dstOff, extend int) {

	b.Impl().Copy(osrc, srcOff, size, dstOff, extend)
}

// NewFromBytes ... return new NoLayer instance with byte buffer.
func (b NoLayer) NewFromBytes(bytes []byte) IO {
	return NewNoLayer(NewBaseImpl(bytes))
}

// New ... return new NoLayer instance
func (b NoLayer) New(n IO) IO {
	return NewNoLayer(n)
}

// DoubleLayer is BaseImpl with Single Diff
type DoubleLayer struct {
	*BaseImpl
}

// NewDoubleLayer ... return DobuleLayer
func NewDoubleLayer(b IO) DoubleLayer {
	if _, already := b.(DoubleLayer); already {
		return b.(DoubleLayer)
	}
	if impl, ok := b.(*BaseImpl); ok {

		if len(impl.Diffs) > 0 {
			impl.Flatten()
		}
		return DoubleLayer{BaseImpl: impl}
	}
	return DoubleLayer{}

}

func (b DoubleLayer) Dup() (dst IO) {
	return DoubleLayer{BaseImpl: b.BaseImpl.Dup().Impl()}
}

// Type ... type of Base interface
func (b DoubleLayer) Type() uint8 { return BASE_DOUBLE_LAYER }

func (b DoubleLayer) mergeDiffs() {

	if len(b.Diffs) < 2 {
		return
	}

	for i := 0; i < len(b.RDiffs); i++ {
		b.Diffs[0].Merge(&b.RDiffs[i])
	}

	for i := 1; i < len(b.Diffs); i++ {
		b.Diffs[0].Merge(&b.Diffs[i])
	}
	return
}

func (b DoubleLayer) insertBuf(pos, size int) IO {
	return b.insertSpace(pos, size, true)
}

func (b DoubleLayer) insertSpace(pos, size int, isCreate bool) IO {

	b.mergeDiffs()
	newBase := NewBaseImpl(b.bytes)
	newBase.r = b.r
	newBase.checkBoundary = b.checkBoundary

	newBase.Diffs = make([]Diff, len(b.Diffs), cap(b.Diffs))
	copy(newBase.Diffs, b.Diffs)

	if len(newBase.Diffs) == 0 {
		blen := MaxInt(size, size+len(newBase.bytes)-pos)
		diff := Diff{Offset: pos, bytes: make([]byte, blen, blen)}
		defer func() {
			newBase.Diffs = append(newBase.Diffs, diff)
		}()

		if len(newBase.bytes) <= pos {
			return newBase
		}
		copy(diff.bytes[:size], newBase.bytes[pos:])

		return newBase
	}

	diff := &newBase.Diffs[0]

	if diff.Include(pos) {

		if cap(diff.bytes[:pos])-len(diff.bytes[:pos]) >= size {
			diff.bytes = diff.bytes[:len(diff.bytes)+size]
			copy(diff.bytes[pos-diff.Offset+size:], diff.bytes[pos-diff.Offset:])
			copy(diff.bytes[pos-diff.Offset:pos-diff.Offset+size], make([]byte, size, size))
		} else {
			nbytes := make([]byte, len(diff.bytes)+size, 2*len(diff.bytes)+size)
			copy(nbytes, diff.bytes[:pos-diff.Offset])
			copy(nbytes[pos-diff.Offset+size:], diff.bytes[pos-diff.Offset:])
			diff.bytes = nbytes
		}
		return newBase
	}

	if diff.Offset > pos {
		nbytes := make([]byte, len(diff.bytes)+size, 2*(len(diff.bytes)+size))
		copy(nbytes[size:], diff.bytes)
		diff.Offset = pos
		diff.bytes = nbytes
	} else {
		nbytes := make([]byte, (pos - diff.Offset + size), 2*(pos-diff.Offset+size))
		copy(nbytes, diff.bytes)
		diff.bytes = nbytes
	}

	return newBase

}

// R ... read buffer
func (b DoubleLayer) R(off int, opts ...OptIO) []byte {

	b.mergeDiffs()

	return NoLayer{BaseImpl: b.BaseImpl}.R(off, opts...)
}

// D ... return Diff for write
func (b DoubleLayer) D(off, size int) *Diff {

	b.mergeDiffs()

	if len(b.Diffs) == 0 {
		b.Diffs = append(b.Diffs, Diff{Offset: off, bytes: make([]byte, size, size*2)})
		return &b.Diffs[0]
	}

	return NoLayer{BaseImpl: b.BaseImpl}.D(off, size)

}

// U ... return buffer for update
func (b DoubleLayer) U(off, size int) []byte {
	b.mergeDiffs()

	return NoLayer{BaseImpl: b.BaseImpl}.U(off, size)
}

func (b DoubleLayer) Copy(osrc IO, srcOff, size, dstOff, extend int) {

	NoLayer{BaseImpl: b.BaseImpl}.Copy(osrc, srcOff, size, dstOff, extend)
	b.mergeDiffs()
	return
}

// NewFromBytes ... return new NoLayer instance with byte buffer.
func (b DoubleLayer) NewFromBytes(bytes []byte) IO {
	return NewDoubleLayer(NewBaseImpl(bytes))
}

// New ... return new NoLayer instance
func (b DoubleLayer) New(n IO) IO {
	return NewDoubleLayer(n)
}
