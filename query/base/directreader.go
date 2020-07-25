package base

import "io"

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
		return nil
	}
	return
}

func (b DirectReader) ShouldCheckBound() bool {
	return false
}
