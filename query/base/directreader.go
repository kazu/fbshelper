package base

import (
	"io"

	log "github.com/kazu/fbshelper/query/log"
	"github.com/kazu/loncha"
)

// DirectReader ... is Base without buffer in data area
type DirectReader struct {
	Base
	r io.ReaderAt
}

// NewDirectReader ... return new DirectReader instance.
func NewDirectReader(b Base, r io.ReaderAt) DirectReader {
	return DirectReader{Base: b, r: r}
}

// Type ... type of Base
func (b DirectReader) Type() uint8 { return BASE_DIRECT_READER }

// R ... read data.
func (b DirectReader) R(offset int) (result []byte) {
	if b.Base.LenBuf() > offset {
		return b.Base.R(offset)
	}
	bLen := b.Base.LenBuf()
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
func NewNoLayer(b Base) NoLayer {
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

// Type ... type of Base interface
func (b NoLayer) Type() uint8 { return BASE_NO_LAYER }

func (b NoLayer) insertBuf(pos, size int) Base {
	return b.insertSpace(pos, size, true)
}

func (b NoLayer) insertSpace(pos, size int, isCreate bool) Base {
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
func (b NoLayer) R(off int) []byte {

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})
	if e == nil && sn >= 0 {
		return b.Diffs[sn].bytes[off-b.Diffs[sn].Offset:]
	}

	if off < len(b.bytes) {
		return b.bytes[off:]
	}

	// MENTION: should check off +32 ?
	if off+8 < cap(b.bytes) {
		if b.expandBuf(off-len(b.bytes)+32) == nil || off < len(b.bytes) {
			return b.bytes[off:]
		}
	}

	log.Log(LOG_WARN, log.Printf("NoyLayer.R(%d) invalid len(bytes)=%d cap(bytes)=%d\n",
		off, len(b.bytes), cap(b.bytes)))

	newDiff := Diff{Offset: cap(b.bytes), bytes: make([]byte, 0, 512)}
	b.Diffs = append(b.Diffs, newDiff)

	if off-cap(b.bytes) < 0 || len(newDiff.bytes) < off-cap(b.bytes) {
		//panic("hoge")
		return nil
	}

	return newDiff.bytes[off-cap(b.bytes):]

}

// D ... return Diff for write
func (b NoLayer) D(off, size int) *Diff {

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+cap(b.Diffs[i].bytes))
	})

	if e == nil && sn >= 0 {
		if b.Diffs[sn].Offset == off {
			return &b.Diffs[sn]
		}
		off_diff := off - b.Diffs[sn].Offset
		return &Diff{Offset: off, bytes: b.Diffs[sn].bytes[off_diff : off_diff+size]}
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
func (b NoLayer) Copy(osrc Base, srcOff, size, dstOff, extend int) {

	if cap(b.bytes) > dstOff {
		if len(b.bytes) > dstOff {
			diff := Diff{Offset: dstOff, bytes: b.bytes[dstOff:]}
			b.Diffs = append(b.Diffs, diff)
		}
		b.bytes = b.bytes[:dstOff:dstOff]
	}

	for i, diff := range b.Diffs {
		if diff.Offset >= dstOff {
			diff.Offset += extend
		}
		b.Diffs[i] = diff
	}
	srcDiffs := append([]Diff{Diff{Offset: 0, bytes: osrc.R(0)}},
		osrc.GetDiffs()...)

	loncha.Delete(&srcDiffs, func(i int) bool {

		diffComp := srcDiffs[i].Inner(srcOff, size) || srcDiffs[i].Include(srcOff) || srcDiffs[i].Included(srcOff, size) || srcDiffs[i].Innerd(srcOff, size)
		diffComp = !diffComp

		result := srcDiffs[i].Offset > srcOff+size || srcDiffs[i].Offset+len(srcDiffs[i].bytes) <= srcOff
		_ = result
		// if result != diffComp {
		// 	panic("invalid")
		// }

		//return diffComp
		return result
	})

	for _, diff := range srcDiffs {
		diff.Offset -= srcOff
		if diff.Offset < 0 && len(diff.bytes)+diff.Offset > 0 {
			diff.bytes = diff.bytes[-diff.Offset:]
			diff.Offset = 0
			if len(diff.bytes) > size {
				diff.bytes = diff.bytes[:size:size]
			}
		}
		diff.Offset += dstOff
		//diff.Offset = 0
		// diff.Offset += srcOff - dstOff
		// diff.Offset += dstOff
		// diff.bytes = diff.bytes[srcOff:]
		// if len(diff.bytes) > size {
		// 	diff.bytes = diff.bytes[:size:size]
		// }

		loncha.Delete(&b.Diffs, func(i int) bool {
			return diff.Offset <= b.Diffs[i].Offset &&
				b.Diffs[i].Offset+len(b.Diffs[i].bytes) <= diff.Offset+len(diff.bytes)
		})
		if diff.Offset >= 0 {
			b.Diffs = append(b.Diffs, diff)
		}
	}

}

// NewFromBytes ... return new NoLayer instance with byte buffer.
func (b NoLayer) NewFromBytes(bytes []byte) Base {
	return NewNoLayer(NewBaseImpl(bytes))
}

// New ... return new NoLayer instance
func (b NoLayer) New(n Base) Base {
	return NewNoLayer(n)
}

// DoubleLayer is BaseImpl with Single Diff
type DoubleLayer struct {
	*BaseImpl
}

// NewDoubleLayer ... return DobuleLayer
func NewDoubleLayer(b Base) DoubleLayer {
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

// Type ... type of Base interface
func (b DoubleLayer) Type() uint8 { return BASE_DOUBLE_LAYER }

func (b DoubleLayer) mergeDiffs() {

	if len(b.Diffs) < 2 {
		return
	}

	for i := 1; i < len(b.Diffs); i++ {
		b.Diffs[0].Merge(&b.Diffs[i])
	}
	return
}

func (b DoubleLayer) insertBuf(pos, size int) Base {
	return b.insertSpace(pos, size, true)
}

func (b DoubleLayer) insertSpace(pos, size int, isCreate bool) Base {

	b.mergeDiffs()

	newBase := &BaseImpl{
		r:     b.r,
		bytes: b.bytes,
	}
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
func (b DoubleLayer) R(off int) []byte {

	b.mergeDiffs()

	return NoLayer{BaseImpl: b.BaseImpl}.R(off)
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

func (b DoubleLayer) Copy(osrc Base, srcOff, size, dstOff, extend int) {

	NoLayer{BaseImpl: b.BaseImpl}.Copy(osrc, srcOff, size, dstOff, extend)
	b.mergeDiffs()
	return
}

// NewFromBytes ... return new NoLayer instance with byte buffer.
func (b DoubleLayer) NewFromBytes(bytes []byte) Base {
	return NewDoubleLayer(NewBaseImpl(bytes))
}

// New ... return new NoLayer instance
func (b DoubleLayer) New(n Base) Base {
	return NewDoubleLayer(n)
}
