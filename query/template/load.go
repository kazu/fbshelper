package query 

	func SturctTeamplate() string {
		return `{{- $ := . }}
{{- $SName := $.Name}}
{{- $IsTable := $.IsTable }}
package {{$.PkgName}}

import (
    flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)

const (
    DUMMY_{{$SName}} = flatbuffers.VtableMetadataFields
)

const (

    {{- range $i, $v := .Fields}}
        {{$SName}}_{{$v.Name}} =     {{$i}}
    {{- end}}
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
    *base.NodeList
}
    {{- end }}
{{- end }}

{{- if eq (isRoot $SName) true }}
func OpenByBuf(buf []byte) Fbs{{$.Name}} {
	return FbsRoot{
		Node: base.NewNode(&base.Base{Bytes: buf}, int(flatbuffers.GetUOffsetT(buf))),
	}
}
func (node Fbs{{$.Name}}) Len() int {
    info := node.Info()
    size := info.Pos + info.Size

    return size + (8 - (size % 8)) 
}


{{- end }}

{{- if eq $IsTable false }}
func (node Fbs{{$SName}}) Info() base.Info {
    info := base.Info{Pos: node.Pos, Size: -1}
    size := 0
    {{- range $i, $v := .Fields}}
        size += base.SizeOf{{$v.Type}}
    {{- end}}
    info.Size = size

    return info

}
{{- else }}
func (node Fbs{{$SName}}) Info() base.Info {

    info := base.Info{Pos: node.Pos, Size: -1}
    for i := 0; i < len(node.VTable); i++ {
        vInfo := node.ValueInfo(i)
        if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
            info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
        }
    }
    return info    
}

{{- end }}

func (node Fbs{{$SName}}) ValueInfo(i int) base.ValueInfo {

    switch i {
{{- range $i, $v := .Fields}}
    case {{$i}}:
    {{- if eq (isStruct $v.Type) true }}
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = node.{{$v.Name}}().Info().Size     

    {{- else if eq (isUnion $IsUnion $v.Type) true }}
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPosTable(i)
        }
        node.ValueInfos[i].Size = node.{{$v.Name}}().Info(i-1).Size
    {{- else if eq $v.Type "[]byte" }}
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPosBytes(i)
        }
    {{- else if eq (isSlice $v.Type) true }}
         if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPosList(i)
        }
        node.ValueInfos[i].Size = node.{{$v.Name}}().Info().Size
    {{- else if eq (isMessage $v.Type) true }}
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPosTable(i)
        }
        node.ValueInfos[i].Size = node.{{$v.Name}}().Info().Size       
    {{- else }}
        if node.ValueInfos[i].IsNotReady() {
            node.ValueInfoPos(i)
        }
        node.ValueInfos[i].Size = base.SizeOf{{$v.Type}}
    {{- end  }}
 {{- end }}
     }
     return node.ValueInfos[i]
}

func (node Fbs{{$SName}}) FieldAt(i int) interface{} {

    switch i {
{{- range $i, $v := .Fields}}
    case {{$i}}:
        return node.{{$v.Name}}()
 {{- end }}
     }
     return nil
}





{{ range $i, $v := .Fields}}

    {{- if eq $IsTable false }}
func (node Fbs{{$SName}}) {{$v.Name}}() {{$v.Type}} {
    buf := node.Bytes
    pos := node.Pos
        {{- range $ii, $vv := $.Fields}}
            {{- if lt $ii $i }}
                pos += base.SizeOf{{$vv.Type}}
            {{- end }}
        {{- end }}
    return {{$v.Type}}(flatbuffers.Get{{(toCamel $v.Type)}}(buf[pos:]))
}
    {{- else if eq (isUnion $IsUnion $v.Type) true }}

func (node Fbs{{$SName}}) {{$v.Name}}() Fbs{{$v.Type}} {
    if node.VTable[{{$i}}] == 0 {
        return Fbs{{$v.Type}}{}  
    }
    return Fbs{{$v.Type}}{Node: node.ValueTable({{$i}})}
}


    {{- else if eq $v.Type "[]byte" }}
func (node Fbs{{$SName}}) {{$v.Name}}() {{$v.Type}} {
    if node.VTable[{{$i}}] == 0 {
        return nil
    }
    return node.ValueBytes({{$i}})
}

    {{- else if eq (isSlice $v.Type) true }}
func (node Fbs{{$SName}}) {{$v.Name}}() Fbs{{$.Name}}{{$v.Name}} {
    if node.VTable[{{$i}}] == 0 {
        return Fbs{{$.Name}}{{$v.Name}}{}
    }
    nodelist :=  node.ValueList({{$i}})
    return Fbs{{$.Name}}{{$v.Name}}{
                NodeList: &nodelist,
    }
}

    {{- else if eq (isStruct $v.Type) true}}
func (node Fbs{{$SName}}) {{$v.Name}}() Fbs{{$v.Type}} {    
   if node.VTable[{{$i}}] == 0 {
        return Fbs{{$v.Type}}{}  
    }
    return Fbs{{$v.Type}}{Node: node.ValueStruct({{$i}})}

}

    {{- else if eq (isMessage $v.Type) true }}
func (node Fbs{{$SName}}) {{$v.Name}}() Fbs{{$v.Type}} {
    if node.VTable[{{$i}}] == 0 {
        return Fbs{{$v.Type}}{}  
    }
    return Fbs{{$v.Type}}{Node: node.ValueTable({{$i}})}
}

    {{- else }}
func (node Fbs{{$SName}}) {{$v.Name}}() {{$v.Type}} {
    if node.VTable[{{$i}}] == 0 {
        return {{$v.Type}}(0)
    }
    return {{$v.Type}}(flatbuffers.Get{{(toCamel $v.Type)}}(node.ValueNormal({{$i}})))
}

    {{- end}}

{{ end}}


func (node Fbs{{$.Name}}) CountOfField() int {
    return {{len $.Fields}}
}

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
    if i > int(node.ValueInfo.VLen) || i < 0 {
		return Fbs{{ $SingleName }}{}
	}

	buf := node.Bytes
	ptr := uint32(node.ValueInfo.Pos + (i-1)*4)
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
	return node.At(int(node.ValueInfo.VLen))
}

func (node Fbs{{$.Name}}{{$v.Name}}) Select(fn func(m Fbs{{ $SingleName }}) bool) []Fbs{{ $SingleName }} {

	result := make([]Fbs{{ $SingleName }}, 0, int(node.ValueInfo.VLen))
	for i := 0; i < int(node.ValueInfo.VLen); i++ {
		if m := node.At(i); fn(m) {
			result = append(result, m)
		}
	}
	return result
}

func (node Fbs{{$.Name}}{{$v.Name}}) Find(fn func(m Fbs{{ $SingleName }}) bool) Fbs{{ $SingleName }}{

	for i := 0; i < int(node.ValueInfo.VLen); i++ {
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
	return int(node.ValueInfo.VLen)
}

func (node Fbs{{$.Name}}{{$v.Name}}) Info() base.Info {

    info := base.Info{Pos: node.ValueInfo.Pos, Size: -1}
    
    {{- if eq $SingleName "bytes"}}
    buf := node.Bytes
	ptr := uint32(node.ValueInfo.Pos + (int(node.ValueInfo.VLen)-1)*4)

    vInfo := base.FbsStringInfo(base.NewNode(node.Base, int(ptr+flatbuffers.GetUint32(buf[ptr:]))))
    {{- else }}
    vInfo := node.Last().Info()
    {{- end }}



    if info.Pos + info.Size < vInfo.Pos + vInfo.Size {
        info.Size = (vInfo.Pos + vInfo.Size) - info.Pos
    }
    return info
}



    {{- end}}

{{- end }}
`
	}

	func UnionTeamplate() string {
		return `{{- $ := . }}
{{- $Name := $.Name}}
package {{.PkgName}}

import (
	//flatbuffers "github.com/google/flatbuffers/go"
    base "github.com/kazu/fbshelper/query/base"
)



{{- $Type := "byte" }}

{{- if le (len .Aliases) 8 -}}
{{- else if le (len .Aliases) 16 -}}
    {{- $Type = "uint16" }}
{{- else if le (len .Aliases) 32 -}}
    {{- $Type = "uint32" }}
{{- else -}}
    {{- $Type = "uint64" }}
{{- end }}

type Enum{{.Name}} {{$Type}}

const ( 
    {{$Name}}NONE                   Enum{{$Name}}   =  0
{{- range $i, $v := .Aliases}}
    {{$Name}}{{$v}}                 Enum{{$Name}}   = {{(add $i 1)}}
{{- end }}
)


type Fbs{{.Name}} struct {
    *base.Node
}
{{- range $i, $v := .Aliases}}
func(node Fbs{{$Name}}) {{$v}}() Fbs{{$v}} {

    return Fbs{{$v}}{Node: node.Node}
}
{{- end }}

func(node Fbs{{$Name}}) Info(i int) base.Info {
    switch i-1 {
{{- range $i, $v := .Aliases}}
    case {{$i}}:
        return node.{{$v}}().Info()
{{- end}}
    }

    return base.Info{}
}


` 
	}
