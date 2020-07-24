package base

import (
	"fmt"
	"io"

	"github.com/kazu/fbshelper/query/log"
	"github.com/kazu/loncha"
)

type CommonList struct {
	*CommonNode
	dataW io.Writer
	dCur  int
	dLen  int
}

func (l *CommonList) SetDataWriter(w io.Writer) {
	l.dCur = 0
	l.dataW = w
	l.dCur = 0
}

func (l *CommonList) WriteDataAll() (e error) {

	if l.Count() == 0 {
		return
	}
	l.Merge()

	first, e := l.CommonNode.At(0)
	if e != nil {
		return e
	}

	vSize := first.CountOfField()*2 + 4
	first.Node.Size = first.Info().Size

	g := GetTypeGroup(first.Name)
	if !IsFieldTable(g) {
		return log.ERR_NO_SUPPORT
	}

	for i, elm := range l.All() {
		pos := elm.Node.Pos - vSize
		elm.Node.Size = elm.Info().Size
		size := elm.Node.Size + vSize
		l.dataW.Write(l.R(pos)[:size])
		l.dCur = i
		l.dLen += size
	}
	last := l.LenBuf()
	_ = last

	//remove writed data area
	bytes := l.R(0)
	bytes = bytes[0 : first.Node.Pos-vSize : first.Node.Pos-vSize]
	l.Base = NewBase(bytes)

	return nil

}

func (l *CommonList) At(i int) (*CommonNode, error) {
	if l.dataW == nil {
		return l.CommonNode.At(i)
	}
	// FIXME
	return nil, log.ERR_NO_SUPPORT
}

func (l *CommonList) SetAt(i int, elm *CommonNode) error {
	if l.dataW == nil {
		return l.CommonNode.SetAt(i, elm)
	}
	if l.Count() != i {
		return log.ERR_NO_SUPPORT
	}

	vSize := elm.CountOfField()*2 + 4
	elm.Node.Size = elm.Info().Size
	pos := elm.Node.Pos - vSize
	size := elm.Node.Size + vSize
	//l.dataW.Write(elm.R(pos)[:size])

	//elm.Merge()

	l.CommonNode.SetAt(i, elm)
	last := l.LenBuf()
	_ = last
	l.WriteElm(elm, pos, size)
	l.dCur = i

	return nil
}

func (l *CommonList) dataStart() int {

	return l.NodeList.ValueInfo.Pos + int(l.VLen())*4

}

func (l *CommonList) WriteElm(elm *CommonNode, pos, size int) {

	defer func() {
		l.dLen += size
	}()
	if len(l.GetDiffs()) == 0 {
		l.dataW.Write(elm.R(pos)[:size])
		//remove writed data area
		bytes := l.R(0)
		if len(bytes) >= pos+size {
			bytes = bytes[0 : pos+size : pos+size]
			l.Base = NewBase(bytes)
		}
		return
	}
	wStart := l.dataStart() + l.dLen

	if w, ok := l.dataW.(io.WriterAt); ok {

		//w.WriteAt(elm.R(0), int64(l.dLen+pos))
		cur := l.dLen //wStart

		w.WriteAt(elm.R(0), int64(cur))

		for i := range elm.GetDiffs() {
			bytes := elm.GetDiffs()[i].bytes
			offset := elm.GetDiffs()[i].Offset
			if len(bytes)+offset > size {
				fmt.Printf("????\n")
			}
			w.WriteAt(bytes, int64(cur+pos+offset))
		}
	}
	// remove written data
	diffs := l.GetDiffs()
	loncha.Delete(&diffs, func(i int) bool {
		return wStart+pos+size <= diffs[i].Offset+len(diffs[i].bytes) || wStart <= diffs[i].Offset
	})
	l.SetDiffs(diffs)
	return
}
