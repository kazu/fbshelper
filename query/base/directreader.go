package base

import (
	"io"

	log "github.com/kazu/fbshelper/query/log"
	"github.com/kazu/loncha"
)

type DirectReader struct {
	Base
	r io.ReaderAt
}

func NewDirectReader(b Base, r io.ReaderAt) DirectReader {
	return DirectReader{Base: b, r: r}
}

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

func (b DirectReader) ShouldCheckBound() bool {
	return false
}

type NoLayer struct {
	*BaseImpl
}

func NewNoLayer(b Base) NoLayer {
	impl := b.(*BaseImpl)
	return NoLayer{BaseImpl: impl}
}

func (b NoLayer) insertBuf(pos, size int) Base {
	return b.insertSpace(pos, size, true)
}

func (b NoLayer) insertSpace(pos, size int, isCreate bool) Base {
	nImple := b.BaseImpl.insertSpace(pos, size, isCreate).(*BaseImpl)
	if nImple.LenBuf() >= pos+size {
		b.BaseImpl = nImple
	} else {
		b.BaseImpl = b.BaseImpl.insertSpace(pos, size, true).(*BaseImpl)
	}

	return b
}

func (b NoLayer) R(off int) []byte {

	if off < len(b.bytes) {
		return b.bytes[off:]
	}

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+len(b.Diffs[i].bytes))
	})
	if e == nil && sn >= 0 {
		return b.Diffs[sn].bytes[off-b.Diffs[sn].Offset:]
	}

	if off < cap(b.bytes) {
		if b.expandBuf(off-len(b.bytes)) == nil {
			return b.bytes[off:]
		}
	}

	log.Log(LOG_WARN, log.Printf("NoyLayer.R(%d) invalid len(bytes)=%d cap(bytes)=%d\n",
		off, len(b.bytes), cap(b.bytes)))

	newDiff := Diff{Offset: cap(b.bytes), bytes: make([]byte, 0, 512)}
	b.Diffs = append(b.Diffs, newDiff)

	return newDiff.bytes[off-cap(b.bytes):]

}

func (b NoLayer) D(off, size int) *Diff {

	sn, e := loncha.LastIndexOf(b.Diffs, func(i int) bool {
		return b.Diffs[i].Offset <= off && off < (b.Diffs[i].Offset+cap(b.Diffs[i].bytes))
	})

	if e == nil && sn >= 0 {
		if b.Diffs[sn].Offset == off {
			return &b.Diffs[sn]
		}
		off_diff := off - b.Diffs[sn].Offset
		return &Diff{Offset: off, bytes: b.Diffs[sn].bytes[off_diff:]}
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

func (b NoLayer) U(off, size int) []byte {

	diff := b.D(off, size)
	if cap(b.bytes) < size {
		diff.bytes = make([]byte, size)
	}
	return diff.bytes
}

func (b NoLayer) Copy(osrc Base, srcOff, size, dstOff, extend int) {

	// src, ok := osrc.(*BaseImpl)
	// _ = src
	// if !ok {
	// 	log.Log(LOG_WARN, log.Printf("BaseImpl.Copy() src is only BaseImpl"))

	// 	return
	// }
	if cap(b.bytes) > dstOff {
		diff := Diff{Offset: dstOff, bytes: b.bytes[dstOff:]}
		b.Diffs = append(b.Diffs, diff)
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
		return srcDiffs[i].Offset > srcOff+size || srcDiffs[i].Offset+len(srcDiffs[i].bytes) < srcOff
	})

	for _, diff := range srcDiffs {
		diff.Offset += dstOff - srcOff
		loncha.Delete(&b.Diffs, func(i int) bool {
			return diff.Offset <= b.Diffs[i].Offset &&
				b.Diffs[i].Offset+len(b.Diffs[i].bytes) <= diff.Offset+len(diff.bytes)
		})
		if diff.Offset > 0 {
			b.Diffs = append(b.Diffs, diff)
		}
	}

}
func (b NoLayer) New(bytes []byte) Base {

	return NewNoLayer(NewBase(bytes))

}
