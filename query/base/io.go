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
	r           io.Reader
	cur         int
	orig        *BaseImpl
	offFromOrig int
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
	return ReaderAt{r: r, cur: 0, orig: nil, offFromOrig: 0}
}

func (r ReaderAt) origReadAt(buf []byte, off int64) (n int, err error) {

	if r.orig == nil {
		return -1, ERR_NO_SUPPORT
	}

	data := r.orig.R(int(off)+r.offFromOrig, Size(len(buf)))
	if len(data) < len(buf) {
		buf = buf[:len(data)]
		copy(buf, data)
	}
	return len(data), nil

}

// ReadAt ... warpiing ReadAt with Read
func (r ReaderAt) ReadAt(buf []byte, off int64) (n int, err error) {

	if r.orig != nil {
		return r.origReadAt(buf, off)
	}

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

	if r.orig != nil {
		n, err = r.origReadAt(buf, int64(r.cur))
	} else {
		n, err = r.r.Read(buf)
	}

	if n > 0 {
		r.cur += n
	}
	return

}
