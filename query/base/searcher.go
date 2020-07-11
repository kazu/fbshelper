package base

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
	TraverseInfo(pos int, fn TraverseRec, condFn TraverseCond)
}

type UnionNoder interface {
	Info(int) Info
	Member(int) Noder
}

type TraverseCond func(int, int, int) bool
type TraverseRec func(*CommonNode, int, int, int, int)
