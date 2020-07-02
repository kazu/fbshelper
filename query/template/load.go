package query 

	func SturctTeamplate() string {
		return `{{- $ := . }}
{{- $SName := $.Name}}
package {{$.PkgName}}

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type Fbs{{$.Name}} struct {
	*base.Node
}

{{ $IsUnion := $.IsUnion}}

{{- range $i, $v := .Fields}}
    {{- if eq (isUnion $IsUnion $v.Name) true }}
type FbsRoot{{$.Name}}{{$v.Name}} struct {
    *base.Node
}
    {{- else if eq $v.Type "[]byte" }}
    {{- else if eq (isSlice $v.Type) true }}
type Fbs{{$.Name}}{{$v.Name}} struct {
    *base.Node
    VPos   uint32
	VLen   uint32
	VStart uint32
}
    {{- end }}
{{- end }}

{{- if eq (isRoot $SName) true }}
func OpenByBuf(buf []byte) Fbs{{$.Name}} {
	return FbsRoot{
		Node: base.NewNode(&base.Base{Bytes: buf}, int(flatbuffers.GetUOffsetT(buf))),
	}
}
{{- end }}

{{- range $i, $v := .Fields}}
    {{- if eq (isUnion $IsUnion $v.Type) true }}
func (node Fbs{{$SName}}) {{$v.Name}}() Fbs{{$v.Type}} {
    {{- else if eq $v.Type "[]byte" }}
func (node Fbs{{$SName}}) {{$v.Name}}() {{$v.Type}} {
    {{- else if eq (isSlice $v.Type) true }}
//func (node Fbs{{$SName}}) {{$v.Name}}() {
func (node Fbs{{$SName}}) {{$v.Name}}() Fbs{{$.Name}}{{$v.Name}} {
    {{- else if eq (isMessage $v.Type) true }}
func (node Fbs{{$SName}}) {{$v.Name}}() Fbs{{$v.Type}} {
    {{- else }}
func (node Fbs{{$SName}}) {{$v.Name}}() {{$v.Type}} {
    {{- end }}
	        if node.VTable[{{$i}}] == 0 {
                    {{- if eq $v.Type "[]byte" }}
                return nil
                    {{- else if eq (isSlice $v.Type) true }}
                return Fbs{{$.Name}}{{$v.Name}}{}
                    {{- else if eq (isMessage $v.Type) true }}
                return Fbs{{$v.Type}}{}        
                    {{- else }}    
		        return {{$v.Type}}(0)
                    {{- end }}
	        }
    {{- if eq $v.Type "[]byte" }}
            buf := node.Bytes
	        pos := uint32(node.Pos + int(node.VTable[{{$i}}]))
	        sLenOff := flatbuffers.GetUint32(buf[pos:])
	        sLen := flatbuffers.GetUint32(buf[pos+sLenOff:])
	        start := pos + sLenOff + flatbuffers.SizeUOffsetT

            return buf[start:start+sLen]

    {{- else if eq (isUnion $IsUnion $v.Type) true }}
            pos := node.Pos + int(node.VTable[{{$i}}])
            return Fbs{{$v.Type}}{Node: base.NewNode(node.Base, int(flatbuffers.GetUint32(node.Bytes[pos:]))+pos)}
    {{- else if eq (isSlice $v.Type) true }}
            buf := node.Bytes
            vPos := uint32(node.Pos + int(node.VTable[{{$i}}]))
            vLenOff := flatbuffers.GetUint32(buf[vPos:])
            vLen := flatbuffers.GetUint32(buf[vPos+vLenOff:])
            start := vPos + vLenOff + flatbuffers.SizeUOffsetT

            return Fbs{{$.Name}}{{$v.Name}}{
                Node: base.NewNode(node.Base, node.Pos),
                VPos:   vPos,
		        VLen:   vLen,
		        VStart: start,
            }

    {{- else if eq (isMessage $v.Type) true }}
            pos := node.Pos + int(node.VTable[{{$i}}])
            return Fbs{{$v.Type}}{Node: base.NewNode(node.Base, int(flatbuffers.GetUint32(node.Bytes[pos:]))+pos)}
    {{- else }}
            pos := node.Pos + int(node.VTable[{{$i}}])
            return {{$v.Type}}(flatbuffers.Get{{(toCamel $v.Type)}}(node.Bytes[pos:]))              
    {{- end }}
}
{{- end }}

{{- range $i, $v := .Fields}}
    {{- if eq $v.Type "[]byte" }}
    {{- else if eq (isSlice $v.Type) true }}

{{ $SingleName := (toBareType $v.Type) }}

        {{- if eq (isSlice $SingleName) true }}
type Fbs{{ (toBareType $SingleName) }}s {{$SingleName}}    
            {{ $SingleName = (toBareType $SingleName) }}
            {{ $SingleName = printf "%ss" $SingleName }}
        {{- end }}


func (node Fbs{{$.Name}}{{$v.Name}}) At(i int) Fbs{{ $SingleName }} {
    if i > int(node.VLen) || i < 0 {
		return Fbs{{ $SingleName }}{}
	}

	buf := node.Bytes
	ptr := node.VStart + uint32(i-1)*4
{{- if eq $SingleName "bytes"}}
    return Fbs{{$SingleName}}(base.FbsString(base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:])))))
{{- else }}
	return Fbs{{$SingleName}}{Node: base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:])))}
{{- end }}
}


func (node Fbs{{$.Name}}{{$v.Name}}) First() Fbs{{ $SingleName }} {
	return node.At(0)
}


func (node Fbs{{$.Name}}{{$v.Name}}) Last() Fbs{{ $SingleName }} {
	return node.At(int(node.VLen))
}

func (node Fbs{{$.Name}}{{$v.Name}}) Select(fn func(m Fbs{{ $SingleName }}) bool) []Fbs{{ $SingleName }} {

	result := make([]Fbs{{ $SingleName }}, 0, int(node.VLen))
	for i := 0; i < int(node.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node Fbs{{$.Name}}{{$v.Name}}) Find(fn func(m Fbs{{ $SingleName }}) bool) Fbs{{ $SingleName }}{

	for i := 0; i < int(node.VLen); i++ {
		if m := node.At(i); fn(m) {
			return m
		}
	}
	return Fbs{{ $SingleName }}{}
}

func (node Fbs{{$.Name}}{{$v.Name}}) All() []Fbs{{ $SingleName }} {
	return node.Select(func(m Fbs{{ $SingleName }}) bool { return true })
}

func (node Fbs{{$.Name}}{{$v.Name}}) Count() int {
	return int(node.VLen)
}


    {{- end}}

{{- end }}`
	}

	func UnionTeamplate() string {
		return `{{- $ := . }}
{{- $Name := $.Name}}
package {{.PkgName}}

import (
	//flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

type Fbs{{.Name}} struct {
    *base.Node
}
{{- range $i, $v := .Aliases}}
func(node Fbs{{$Name}}) {{$v}}() Fbs{{$v}} {

    return Fbs{{$v}}{Node: node.Node}
}
{{- end }}` 
	}