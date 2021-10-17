package base_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	query2 "github.com/kazu/fbshelper/example/query2"
	"github.com/kazu/fbshelper/example/vfs_schema"
	"github.com/kazu/fbshelper/query/base"
	log "github.com/kazu/fbshelper/query/log"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	base.SetL2Current(log.LOG_ERROR, base.FLT_NORMAL)

	base.SetDefaultBase("")(&base.DefaultOption)
	m.Run()
	base.SetDefaultBase("")(&base.DefaultOption)

}

func TestHoge(t *testing.T) {

	//base.SetDefaultBase("")(&base.DefaultOption)
	buf := MakeRootWithRecord(111, 666, 12, 34)

	base.SetDefaultBase("")(&base.DefaultOption)
	root := query2.OpenByBuf(buf)

	base.SetDefaultBase("NoLayer")(&base.DefaultOption)
	rootnl := query2.OpenByBuf(buf)

	base.SetDefaultBase("DobuleLayer")(&base.DefaultOption)
	rootdl := query2.OpenByBuf(buf)

	assert.True(t, root.Equal(rootnl))
	assert.True(t, root.Equal(rootdl))

}

func ToBaseImpl(c *base.CommonNode) *base.BaseImpl {

	switch v := c.Base.(type) {
	case base.NoLayer:
		return v.BaseImpl
	case base.DoubleLayer:
		return v.BaseImpl
	case *base.BaseImpl:
		return v
	}
	return nil

}

func ToBytes(c *base.CommonNode) []byte {
	if impl := ToBaseImpl(c); impl != nil {
		return impl.Bytes()
	}
	return nil
}

func HexDump(s string, b []byte) {
	fmt.Fprintf(DefaultWriter, "--hexdump--start---\n%s\n", s)
	stdoutDumper := hex.Dumper(DefaultWriter)
	stdoutDumper.Write(b)
	fmt.Fprint(DefaultWriter, "\n--hexdump--end-----\n")
}
func Test_AddDiffInNoLayer(t *testing.T) {

	withoutE := func(q *query2.InvertedMapNum, e error) *query2.InvertedMapNum {
		return q
	}
	_ = withoutE

	obuf := MakeRootIndexNum(func() *query2.IndexNum {
		return MakeIndexNum("IndexNum R test", 11, 21, 31)
	})

	buf := make([]byte, len(obuf), cap(obuf))
	copy(buf, obuf)
	base.SetDefaultBase("NoLayer")(&base.DefaultOption)
	nolayer := query2.OpenByBuf(buf)

	idxNum := nolayer.Index().IndexNum()
	invlist := idxNum.Maps()

	invlist.Impl().Diffs = append(invlist.Impl().Diffs, base.Diff{Offset: 1000})

	assert.True(t, invlist.Impl().Equal(nolayer.Impl()))
}

func DupBytes(obuf []byte) []byte {

	buf := make([]byte, len(obuf), cap(obuf))
	copy(buf, obuf)
	return buf

}

func Test_RootIndexNum(t *testing.T) {

	// withoutE := func(q *query2.InvertedMapNum, e error) *query2.InvertedMapNum {
	// 	return q
	// }

	base.SetDefaultBase("")(&base.DefaultOption)
	obuf := MakeRootIndexNum(func() *query2.IndexNum {
		return MakeIndexNum("IndexNum R test", 11, 21, 31)
	})

	buf := DupBytes(obuf)
	base.SetDefaultBase("")(&base.DefaultOption)
	expect := query2.OpenByBuf(buf)
	expectCnt := expect.Index().IndexNum().Maps().Count()

	buf = DupBytes(obuf)
	base.SetDefaultBase("NoLayer")(&base.DefaultOption)
	nolayer := query2.OpenByBuf(buf)
	nolayerCnt := nolayer.Index().IndexNum().Maps().Count()

	buf = DupBytes(obuf)
	base.SetDefaultBase("DoubleLayer")(&base.DefaultOption)
	doubleLayer := query2.OpenByBuf(buf)
	doubleLayerCnt := doubleLayer.Index().IndexNum().Maps().Count()

	assert.True(t, expect.Equal(nolayer))
	assert.True(t, expect.Equal(doubleLayer))

	assert.Equal(t, expectCnt, nolayerCnt)
	assert.Equal(t, expectCnt, doubleLayerCnt)

	editFn := func(root *query2.Root, isDump bool) {
		idxNum := root.Index().IndexNum()
		invlist := idxNum.Maps()
		if isDump {
			HexDump(
				fmt.Sprintf("invlist.NodeList.Pos-4(from vector len)=%d after getting invlis",
					invlist.NodeList.ValueInfo.Pos-4),
				invlist.Impl().R(invlist.NodeList.ValueInfo.Pos-4))
		}

		inv := query2.NewInvertedMapNum()

		rec := query2.NewRecord()
		rec.SetFileId(query2.FromUint64(98 + 1))
		rec.SetOffset(query2.FromInt64(98 + 2))

		inv.SetKey(query2.FromInt64(98))
		inv.SetValue(rec)

		if isDump {
			HexDump(
				fmt.Sprintf("invlist.NodeList.Pos-4(from vector len)=%d before add to invlis",
					invlist.NodeList.ValueInfo.Pos-4),
				invlist.Impl().R(invlist.NodeList.ValueInfo.Pos-4))
		}

		invlist.SetAt(invlist.Count(), inv)
		if isDump {
			HexDump(
				fmt.Sprintf("invlist.NodeList.Pos-4(from vector len)=%d after add to invlis",
					invlist.NodeList.ValueInfo.Pos-4),
				invlist.Impl().R(invlist.NodeList.ValueInfo.Pos-4))
		}

		//assert.True(t, inv.Equal(withoutE(invlist.Last())))
		inv.Equal(*invlist.AtWihoutError(invlist.Count() - 1))

		if isDump {
			HexDump(
				fmt.Sprintf("invlist.NodeList.Pos-4(from vector len)=%d before set invlist to expect",
					invlist.NodeList.ValueInfo.Pos-4),
				invlist.Impl().R(invlist.NodeList.ValueInfo.Pos-4))
		}
		idxNum.SetMaps(invlist)
		if isDump {
			HexDump(
				fmt.Sprintf("idxNum.Maps().NodeList.Pos-4(from vector len)=%d before set invlist to expect",
					idxNum.Maps().NodeList.ValueInfo.Pos-4),
				idxNum.Maps().Impl().R(invlist.NodeList.ValueInfo.Pos-4))
		}
		root.SetIndex(&query2.Index{CommonNode: idxNum.CommonNode})
	}

	base.SetDefaultBase("")(&base.DefaultOption)
	HexDump(
		fmt.Sprintf("expect.Index().IndexNum().Maps().NodeList.Pos-4(from vector len)=%d before editting expect",
			expect.Index().IndexNum().Maps().NodeList.ValueInfo.Pos-4),
		expect.Impl().R(expect.Index().IndexNum().Maps().NodeList.ValueInfo.Pos-4))
	editFn(&expect, false)
	hoge := expect.Index().IndexNum()
	maps := hoge.Maps()
	_ = maps
	HexDump(
		fmt.Sprintf("expect.Index().IndexNum().Maps().NodeList.Pos-4(from vector len)=%d after editting expect",
			expect.Index().IndexNum().Maps().NodeList.ValueInfo.Pos-4),
		expect.Impl().R(expect.Index().IndexNum().Maps().NodeList.ValueInfo.Pos-4))

	base.SetDefaultBase("NoLayer")(&base.DefaultOption)
	editFn(&nolayer, false)

	base.SetDefaultBase("DoubleLayer")(&base.DefaultOption)
	editFn(&doubleLayer, true)

	assert.True(t, expect.Equal(nolayer))
	assert.True(t, expect.Equal(doubleLayer))

}

func MakeRootIndexNum(fn func() *query2.IndexNum) []byte {

	base.SetDefaultBase("")(&base.DefaultOption)

	root := query2.NewRoot()
	root.WithHeader()
	root.SetVersion(query2.FromInt32(1))
	root.SetIndexType(query2.FromByte(byte(vfs_schema.IndexIndexNum)))
	root.SetIndex(&query2.Index{CommonNode: fn().CommonNode})

	rec := query2.NewRecord()
	rec.SetFileId(query2.FromUint64(21))
	rec.SetOffset(query2.FromInt64(22))
	rec.SetSize(query2.FromInt64(23))
	rec.SetOffsetOfValue(query2.FromInt32(0))
	rec.SetValueSize(query2.FromInt32(0))
	root.SetRecord(rec)
	root.Flatten()

	return root.R(0)
}

func MakeIndexNum(prefix string, cntOfMaps, prefixKey, offsetRec int) *query2.IndexNum {
	idxNum := query2.NewIndexNum()

	idxNum.SetSize(query2.FromInt32(12))

	invList := query2.NewInvertedMapNumList()

	for i := 0; i < cntOfMaps; i++ {
		inv := query2.NewInvertedMapNum()
		inv.SetKey(query2.FromInt64(int64(prefixKey + i)))

		rec := query2.NewRecord()
		rec.SetFileId(query2.FromUint64(uint64(offsetRec + i)))
		rec.SetOffset(query2.FromInt64(int64(offsetRec + i + 100)))

		inv.SetValue(rec)
		invList.SetAt(invList.Count(), inv)
	}
	idxNum.SetMaps(invList)
	return idxNum
}

func Test_R(t *testing.T) {

	buf := MakeRootIndexNum(func() *query2.IndexNum {
		return MakeIndexNum("IndexNum R test", 11, 21, 31)
	})

	base.SetDefaultBase("")(&base.DefaultOption)
	expect := query2.OpenByBuf(DupBytes(buf))

	base.SetDefaultBase("NoLayer")(&base.DefaultOption)
	nolayer := query2.OpenByBuf(DupBytes(buf))

	base.SetDefaultBase("DoubleLayer")(&base.DefaultOption)
	doubleLayer := query2.OpenByBuf(DupBytes(buf))

	prepareFn := func(root *query2.Root, isTarget bool, idxOfTarget int) {
		if !isTarget {
			base.SetDefaultBase("")(&base.DefaultOption)
			v := query2.OpenByBuf(DupBytes(buf))
			root = &v
		} else {
			if idxOfTarget == 0 {
				base.SetDefaultBase("NoLayer")(&base.DefaultOption)
				v := query2.OpenByBuf(DupBytes(buf))
				root = &v
			} else {
				base.SetDefaultBase("DoubleLayer")(&base.DefaultOption)
				v := query2.OpenByBuf(DupBytes(buf))
				root = &v
			}
		}
	}
	withoutE := func(q *query2.InvertedMapNum, e error) *query2.InvertedMapNum {
		return q
	}

	tests := []struct {
		name    string
		expect  *query2.Root
		targets []*query2.Root
		prepare func(*query2.Root, bool, int)
		off     func(*query2.Root) int
	}{
		{
			name:    "Get root.Index().IndexNum().Size()",
			expect:  &expect,
			targets: []*query2.Root{&nolayer, &doubleLayer},
			prepare: prepareFn,
			off: func(root *query2.Root) int {
				return root.Index().IndexNum().Size().Node.Pos
			},
		},
		{
			name:    "Get root.Index().IndexNum().Size()",
			expect:  &expect,
			targets: []*query2.Root{&nolayer, &doubleLayer},
			prepare: prepareFn,
			off: func(root *query2.Root) int {
				idxNum := root.Index().IndexNum()
				invlist := idxNum.Maps()
				cnt := invlist.Count()
				inv := query2.NewInvertedMapNum()

				rec := query2.NewRecord()
				rec.SetFileId(query2.FromUint64(98 + 1))
				rec.SetOffset(query2.FromInt64(98 + 2))

				inv.SetKey(query2.FromInt64(98))
				inv.SetValue(rec)

				invlist.SetAt(invlist.Count(), inv)
				idxNum.SetMaps(invlist)
				root.SetIndex(&query2.Index{CommonNode: idxNum.CommonNode})

				return withoutE(root.Index().IndexNum().Maps().At(cnt - 1)).Key().Node.Pos
			},
		},
		// {
		// 	name:    "Get root.Index().IndexNum().Size()",
		// 	expect:  &expect,
		// 	targets: []*query2.Root{&nolayer, &doubleLayer},
		// 	prepare: func(root *query2.Root, isTarget bool, i int) {
		// 		idxNum := root.Index().IndexNum()
		// 		invlist := idxNum.Maps()

		// 		inv := query2.NewInvertedMapNum()

		// 		rec := query2.NewRecord()
		// 		rec.SetFileId(query2.FromUint64(198 + 1))
		// 		rec.SetOffset(query2.FromInt64(198 + 2))

		// 		inv.SetKey(query2.FromInt64(98))
		// 		inv.SetValue(rec)

		// 		invlist.SetAt(invlist.Count(), inv)

		// 		if isTarget {
		// 			idxNum.SetMaps(invlist)
		// 			root.SetIndex(&query2.Index{CommonNode: idxNum.CommonNode})
		// 		}
		// 	},
		// 	off: func(root *query2.Root) int {
		// 		return root.Index().IndexNum().Maps()
		// 	},
		// },

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(&expect, false, 0)
			eOff := tt.off(&expect)

			for i, target := range tt.targets {
				tt.prepare(target, true, i)
				off := tt.off(target)
				testR(expect.Base, target.Base, eOff, off, t)
			}
		})
	}

}

func testR(expect base.Base, target base.Base, eoff, off int, t *testing.T) {
	assert.Equal(t, expect.R(eoff), target.R(off)[:len(expect.R(eoff))])
}
