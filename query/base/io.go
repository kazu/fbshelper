package base

import (
	"errors"
	"fmt"
	"io"
)

var (
	AlreadReadError = errors.New("read already")
)

type ReaderWithAt interface {
	io.Reader
	io.ReaderAt
}

type ReaderAt struct {
	r   io.Reader
	cur int
}

type ErrorMustRead struct {
	Off  int64
	Size int
}

func (e ErrorMustRead) Error() string {
	return fmt.Sprintf("must not read off=%d, size=%d", e.Off, e.Size)
}

func (e ErrorMustRead) ToParam() (int, int) {

	return int(e.Off), e.Size
}

func NewReaderAt(r io.Reader) ReaderAt {
	return ReaderAt{r: r, cur: 0}
}

// ReadAt ... warpiing ReadAt with Read
func (r ReaderAt) ReadAt(buf []byte, off int64) (n int, err error) {

	if nr, ok := r.r.(io.ReaderAt); ok {
		return nr.ReadAt(buf, off)
	}

	if int64(r.cur) > off {
		return 0, AlreadReadError
	}
	if int64(r.cur) < off {
		return 0, ErrorMustRead{Off: int64(r.cur), Size: int(off) - r.cur}
	}

	n, err = io.ReadAtLeast(r.r, buf, len(buf))
	if n > 0 {
		r.cur += n
	}

	return
}

func (r ReaderAt) Read(buf []byte) (n int, err error) {

	n, err = r.r.Read(buf)
	if n > 0 {
		r.cur += n
	}
	return

}
