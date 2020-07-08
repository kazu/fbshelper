package base

//"io"
//"reflect"

//"github.com/kazu/loncha"

//flatbuffers "github.com/google/flatbuffers/go"

type CondFn func(int, Info) bool
type RecFn func(NodePath, Info)

type Noder interface {
	IsLeafAt(int) bool
	Info() Info
	ValueInfo(int) ValueInfo
	//FieldAt(int) Noder
	SearchInfo(int, RecFn, CondFn)
}

type Searcher interface {
	SearchInfo(int, RecFn, CondFn)
}

type UnionNoder interface {
	Info(int) Info
	Member(int) Noder
}
