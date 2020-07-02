package fbsparser

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleschema
	rulestatment_decl
	rulenamespace_decl
	ruleinclude
	ruletype_decl
	rulestruct_decl
	ruletable_decl
	ruletypename
	rulemetadata
	rulefield_decl
	rulefield_type
	ruleenum_decl
	ruleenum_fields
	ruleunion_decl
	ruleunion_name
	ruleenum_field
	ruleenum_field_type
	ruleroot_decl
	rulefile_extension_decl
	rulefile_identifier_decl
	ruleattribute_decl
	rulerpc_decl
	ruletype
	rulescalar
	ruleinteger_constant
	rulefloat_constant
	rulefloat_constant_exp
	ruleident
	ruleonly_comment
	rulespacing
	rulespace_comment
	rulecomment
	rulespace
	ruleend_of_line
	ruleend_of_file
	rulePegText
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	ruleAction31
	ruleAction32
	ruleAction33
	ruleAction34
	ruleAction35

	rulePre
	ruleIn
	ruleSuf
)

var rul3s = [...]string{
	"Unknown",
	"schema",
	"statment_decl",
	"namespace_decl",
	"include",
	"type_decl",
	"struct_decl",
	"table_decl",
	"typename",
	"metadata",
	"field_decl",
	"field_type",
	"enum_decl",
	"enum_fields",
	"union_decl",
	"union_name",
	"enum_field",
	"enum_field_type",
	"root_decl",
	"file_extension_decl",
	"file_identifier_decl",
	"attribute_decl",
	"rpc_decl",
	"type",
	"scalar",
	"integer_constant",
	"float_constant",
	"float_constant_exp",
	"ident",
	"only_comment",
	"spacing",
	"space_comment",
	"comment",
	"space",
	"end_of_line",
	"end_of_file",
	"PegText",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",

	"Pre_",
	"_In_",
	"_Suf",
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (node *node32) Print(buffer string) {
	node.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: ruleIn, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: ruleSuf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens32) Expand(index int) {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
}

type Parser struct {
	Fbs

	Buffer string
	buffer []rune
	rules  [73]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokens32
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Parser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *Parser) PrintSyntaxTree() {
	p.tokens32.PrintSyntaxTree(p.Buffer)
}

func (p *Parser) Highlighter() {
	p.PrintSyntax()
}

func (p *Parser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.SetNameSpace(text)
		case ruleAction1:
			p.ExtractStruct(false)
		case ruleAction2:
			p.ExtractStruct(true)
		case ruleAction3:
			p.SetTypeName(text)
		case ruleAction4:
			p.NewExtractField()
		case ruleAction5:
			p.NewExtractFieldWithValue()
		case ruleAction6:
			p.FieldName(text)
		case ruleAction7:
			p.NewUnion()
		case ruleAction8:
			p.UnionName(text)
		case ruleAction9:
			p.NewExtractField()
		case ruleAction10:
			p.EnumName(text)
		case ruleAction11:
			p.SetRootType(text)
		case ruleAction12:
			p.SetType("bool")
		case ruleAction13:
			p.SetType("int8")
		case ruleAction14:
			p.SetType("uint8")
		case ruleAction15:
			p.SetType("int16")
		case ruleAction16:
			p.SetType("uint16")
		case ruleAction17:
			p.SetType("int32")
		case ruleAction18:
			p.SetType("uint32")
		case ruleAction19:
			p.SetType("int64")
		case ruleAction20:
			p.SetType("uint64")
		case ruleAction21:
			p.SetType("float32")
		case ruleAction22:
			p.SetType("float64")
		case ruleAction23:
			p.SetType("uint8")
		case ruleAction24:
			p.SetType("uint8")
		case ruleAction25:
			p.SetType("short")
		case ruleAction26:
			p.SetType("ushort")
		case ruleAction27:
			p.SetType("int32")
		case ruleAction28:
			p.SetType("uint32")
		case ruleAction29:
			p.SetType("float")
		case ruleAction30:
			p.SetType("long")
		case ruleAction31:
			p.SetType("ulong")
		case ruleAction32:
			p.SetType("double")
		case ruleAction33:
			p.SetRepeated("byte")
		case ruleAction34:
			p.SetType(text)
		case ruleAction35:
			p.SetRepeated("")

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *Parser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
		p.buffer = append(p.buffer, endSymbol)
	}

	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		tree.Expand(tokenIndex)
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 schema <- <((include end_of_file) / (include statment_decl+ end_of_file) / (statment_decl+ end_of_file))> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !_rules[ruleinclude]() {
						goto l3
					}
					if !_rules[ruleend_of_file]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleinclude]() {
						goto l4
					}
					if !_rules[rulestatment_decl]() {
						goto l4
					}
				l5:
					{
						position6, tokenIndex6, depth6 := position, tokenIndex, depth
						if !_rules[rulestatment_decl]() {
							goto l6
						}
						goto l5
					l6:
						position, tokenIndex, depth = position6, tokenIndex6, depth6
					}
					if !_rules[ruleend_of_file]() {
						goto l4
					}
					goto l2
				l4:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[rulestatment_decl]() {
						goto l0
					}
				l7:
					{
						position8, tokenIndex8, depth8 := position, tokenIndex, depth
						if !_rules[rulestatment_decl]() {
							goto l8
						}
						goto l7
					l8:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
					}
					if !_rules[ruleend_of_file]() {
						goto l0
					}
				}
			l2:
				depth--
				add(ruleschema, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 statment_decl <- <(namespace_decl / union_decl / type_decl / enum_decl / root_decl / file_extension_decl / file_identifier_decl / attribute_decl / rpc_decl / only_comment)> */
		func() bool {
			position9, tokenIndex9, depth9 := position, tokenIndex, depth
			{
				position10 := position
				depth++
				{
					position11, tokenIndex11, depth11 := position, tokenIndex, depth
					if !_rules[rulenamespace_decl]() {
						goto l12
					}
					goto l11
				l12:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[ruleunion_decl]() {
						goto l13
					}
					goto l11
				l13:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[ruletype_decl]() {
						goto l14
					}
					goto l11
				l14:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[ruleenum_decl]() {
						goto l15
					}
					goto l11
				l15:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[ruleroot_decl]() {
						goto l16
					}
					goto l11
				l16:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulefile_extension_decl]() {
						goto l17
					}
					goto l11
				l17:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulefile_identifier_decl]() {
						goto l18
					}
					goto l11
				l18:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[ruleattribute_decl]() {
						goto l19
					}
					goto l11
				l19:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulerpc_decl]() {
						goto l20
					}
					goto l11
				l20:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[ruleonly_comment]() {
						goto l9
					}
				}
			l11:
				depth--
				add(rulestatment_decl, position10)
			}
			return true
		l9:
			position, tokenIndex, depth = position9, tokenIndex9, depth9
			return false
		},
		/* 2 namespace_decl <- <('n' 'a' 'm' 'e' 's' 'p' 'a' 'c' 'e' spacing <([A-z] / [0-9] / '_' / '.' / '-')+> Action0 ';' spacing)> */
		func() bool {
			position21, tokenIndex21, depth21 := position, tokenIndex, depth
			{
				position22 := position
				depth++
				if buffer[position] != rune('n') {
					goto l21
				}
				position++
				if buffer[position] != rune('a') {
					goto l21
				}
				position++
				if buffer[position] != rune('m') {
					goto l21
				}
				position++
				if buffer[position] != rune('e') {
					goto l21
				}
				position++
				if buffer[position] != rune('s') {
					goto l21
				}
				position++
				if buffer[position] != rune('p') {
					goto l21
				}
				position++
				if buffer[position] != rune('a') {
					goto l21
				}
				position++
				if buffer[position] != rune('c') {
					goto l21
				}
				position++
				if buffer[position] != rune('e') {
					goto l21
				}
				position++
				if !_rules[rulespacing]() {
					goto l21
				}
				{
					position23 := position
					depth++
					{
						position26, tokenIndex26, depth26 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('A') || c > rune('z') {
							goto l27
						}
						position++
						goto l26
					l27:
						position, tokenIndex, depth = position26, tokenIndex26, depth26
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l28
						}
						position++
						goto l26
					l28:
						position, tokenIndex, depth = position26, tokenIndex26, depth26
						if buffer[position] != rune('_') {
							goto l29
						}
						position++
						goto l26
					l29:
						position, tokenIndex, depth = position26, tokenIndex26, depth26
						if buffer[position] != rune('.') {
							goto l30
						}
						position++
						goto l26
					l30:
						position, tokenIndex, depth = position26, tokenIndex26, depth26
						if buffer[position] != rune('-') {
							goto l21
						}
						position++
					}
				l26:
				l24:
					{
						position25, tokenIndex25, depth25 := position, tokenIndex, depth
						{
							position31, tokenIndex31, depth31 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('A') || c > rune('z') {
								goto l32
							}
							position++
							goto l31
						l32:
							position, tokenIndex, depth = position31, tokenIndex31, depth31
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l33
							}
							position++
							goto l31
						l33:
							position, tokenIndex, depth = position31, tokenIndex31, depth31
							if buffer[position] != rune('_') {
								goto l34
							}
							position++
							goto l31
						l34:
							position, tokenIndex, depth = position31, tokenIndex31, depth31
							if buffer[position] != rune('.') {
								goto l35
							}
							position++
							goto l31
						l35:
							position, tokenIndex, depth = position31, tokenIndex31, depth31
							if buffer[position] != rune('-') {
								goto l25
							}
							position++
						}
					l31:
						goto l24
					l25:
						position, tokenIndex, depth = position25, tokenIndex25, depth25
					}
					depth--
					add(rulePegText, position23)
				}
				if !_rules[ruleAction0]() {
					goto l21
				}
				if buffer[position] != rune(';') {
					goto l21
				}
				position++
				if !_rules[rulespacing]() {
					goto l21
				}
				depth--
				add(rulenamespace_decl, position22)
			}
			return true
		l21:
			position, tokenIndex, depth = position21, tokenIndex21, depth21
			return false
		},
		/* 3 include <- <('i' 'n' 'c' 'l' 'u' 'd' 'e' spacing ident comment ';' spacing)> */
		func() bool {
			position36, tokenIndex36, depth36 := position, tokenIndex, depth
			{
				position37 := position
				depth++
				if buffer[position] != rune('i') {
					goto l36
				}
				position++
				if buffer[position] != rune('n') {
					goto l36
				}
				position++
				if buffer[position] != rune('c') {
					goto l36
				}
				position++
				if buffer[position] != rune('l') {
					goto l36
				}
				position++
				if buffer[position] != rune('u') {
					goto l36
				}
				position++
				if buffer[position] != rune('d') {
					goto l36
				}
				position++
				if buffer[position] != rune('e') {
					goto l36
				}
				position++
				if !_rules[rulespacing]() {
					goto l36
				}
				if !_rules[ruleident]() {
					goto l36
				}
				if !_rules[rulecomment]() {
					goto l36
				}
				if buffer[position] != rune(';') {
					goto l36
				}
				position++
				if !_rules[rulespacing]() {
					goto l36
				}
				depth--
				add(ruleinclude, position37)
			}
			return true
		l36:
			position, tokenIndex, depth = position36, tokenIndex36, depth36
			return false
		},
		/* 4 type_decl <- <(struct_decl / table_decl)> */
		func() bool {
			position38, tokenIndex38, depth38 := position, tokenIndex, depth
			{
				position39 := position
				depth++
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if !_rules[rulestruct_decl]() {
						goto l41
					}
					goto l40
				l41:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
					if !_rules[ruletable_decl]() {
						goto l38
					}
				}
			l40:
				depth--
				add(ruletype_decl, position39)
			}
			return true
		l38:
			position, tokenIndex, depth = position38, tokenIndex38, depth38
			return false
		},
		/* 5 struct_decl <- <('s' 't' 'r' 'u' 'c' 't' spacing typename spacing metadata* '{' field_decl+ '}' spacing Action1)> */
		func() bool {
			position42, tokenIndex42, depth42 := position, tokenIndex, depth
			{
				position43 := position
				depth++
				if buffer[position] != rune('s') {
					goto l42
				}
				position++
				if buffer[position] != rune('t') {
					goto l42
				}
				position++
				if buffer[position] != rune('r') {
					goto l42
				}
				position++
				if buffer[position] != rune('u') {
					goto l42
				}
				position++
				if buffer[position] != rune('c') {
					goto l42
				}
				position++
				if buffer[position] != rune('t') {
					goto l42
				}
				position++
				if !_rules[rulespacing]() {
					goto l42
				}
				if !_rules[ruletypename]() {
					goto l42
				}
				if !_rules[rulespacing]() {
					goto l42
				}
			l44:
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					if !_rules[rulemetadata]() {
						goto l45
					}
					goto l44
				l45:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
				}
				if buffer[position] != rune('{') {
					goto l42
				}
				position++
				if !_rules[rulefield_decl]() {
					goto l42
				}
			l46:
				{
					position47, tokenIndex47, depth47 := position, tokenIndex, depth
					if !_rules[rulefield_decl]() {
						goto l47
					}
					goto l46
				l47:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
				}
				if buffer[position] != rune('}') {
					goto l42
				}
				position++
				if !_rules[rulespacing]() {
					goto l42
				}
				if !_rules[ruleAction1]() {
					goto l42
				}
				depth--
				add(rulestruct_decl, position43)
			}
			return true
		l42:
			position, tokenIndex, depth = position42, tokenIndex42, depth42
			return false
		},
		/* 6 table_decl <- <('t' 'a' 'b' 'l' 'e' spacing typename spacing metadata* '{' field_decl+ '}' spacing Action2)> */
		func() bool {
			position48, tokenIndex48, depth48 := position, tokenIndex, depth
			{
				position49 := position
				depth++
				if buffer[position] != rune('t') {
					goto l48
				}
				position++
				if buffer[position] != rune('a') {
					goto l48
				}
				position++
				if buffer[position] != rune('b') {
					goto l48
				}
				position++
				if buffer[position] != rune('l') {
					goto l48
				}
				position++
				if buffer[position] != rune('e') {
					goto l48
				}
				position++
				if !_rules[rulespacing]() {
					goto l48
				}
				if !_rules[ruletypename]() {
					goto l48
				}
				if !_rules[rulespacing]() {
					goto l48
				}
			l50:
				{
					position51, tokenIndex51, depth51 := position, tokenIndex, depth
					if !_rules[rulemetadata]() {
						goto l51
					}
					goto l50
				l51:
					position, tokenIndex, depth = position51, tokenIndex51, depth51
				}
				if buffer[position] != rune('{') {
					goto l48
				}
				position++
				if !_rules[rulefield_decl]() {
					goto l48
				}
			l52:
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if !_rules[rulefield_decl]() {
						goto l53
					}
					goto l52
				l53:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
				}
				if buffer[position] != rune('}') {
					goto l48
				}
				position++
				if !_rules[rulespacing]() {
					goto l48
				}
				if !_rules[ruleAction2]() {
					goto l48
				}
				depth--
				add(ruletable_decl, position49)
			}
			return true
		l48:
			position, tokenIndex, depth = position48, tokenIndex48, depth48
			return false
		},
		/* 7 typename <- <(ident Action3)> */
		func() bool {
			position54, tokenIndex54, depth54 := position, tokenIndex, depth
			{
				position55 := position
				depth++
				if !_rules[ruleident]() {
					goto l54
				}
				if !_rules[ruleAction3]() {
					goto l54
				}
				depth--
				add(ruletypename, position55)
			}
			return true
		l54:
			position, tokenIndex, depth = position54, tokenIndex54, depth54
			return false
		},
		/* 8 metadata <- <('(' <(!')' .)*> ')')> */
		func() bool {
			position56, tokenIndex56, depth56 := position, tokenIndex, depth
			{
				position57 := position
				depth++
				if buffer[position] != rune('(') {
					goto l56
				}
				position++
				{
					position58 := position
					depth++
				l59:
					{
						position60, tokenIndex60, depth60 := position, tokenIndex, depth
						{
							position61, tokenIndex61, depth61 := position, tokenIndex, depth
							if buffer[position] != rune(')') {
								goto l61
							}
							position++
							goto l60
						l61:
							position, tokenIndex, depth = position61, tokenIndex61, depth61
						}
						if !matchDot() {
							goto l60
						}
						goto l59
					l60:
						position, tokenIndex, depth = position60, tokenIndex60, depth60
					}
					depth--
					add(rulePegText, position58)
				}
				if buffer[position] != rune(')') {
					goto l56
				}
				position++
				depth--
				add(rulemetadata, position57)
			}
			return true
		l56:
			position, tokenIndex, depth = position56, tokenIndex56, depth56
			return false
		},
		/* 9 field_decl <- <((spacing field_type ':' type metadata* ';' spacing Action4) / (spacing field_type ':' type <(' ' / '\t')*> '=' <(' ' / '\t')*> scalar metadata* ';' spacing Action5))> */
		func() bool {
			position62, tokenIndex62, depth62 := position, tokenIndex, depth
			{
				position63 := position
				depth++
				{
					position64, tokenIndex64, depth64 := position, tokenIndex, depth
					if !_rules[rulespacing]() {
						goto l65
					}
					if !_rules[rulefield_type]() {
						goto l65
					}
					if buffer[position] != rune(':') {
						goto l65
					}
					position++
					if !_rules[ruletype]() {
						goto l65
					}
				l66:
					{
						position67, tokenIndex67, depth67 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l67
						}
						goto l66
					l67:
						position, tokenIndex, depth = position67, tokenIndex67, depth67
					}
					if buffer[position] != rune(';') {
						goto l65
					}
					position++
					if !_rules[rulespacing]() {
						goto l65
					}
					if !_rules[ruleAction4]() {
						goto l65
					}
					goto l64
				l65:
					position, tokenIndex, depth = position64, tokenIndex64, depth64
					if !_rules[rulespacing]() {
						goto l62
					}
					if !_rules[rulefield_type]() {
						goto l62
					}
					if buffer[position] != rune(':') {
						goto l62
					}
					position++
					if !_rules[ruletype]() {
						goto l62
					}
					{
						position68 := position
						depth++
					l69:
						{
							position70, tokenIndex70, depth70 := position, tokenIndex, depth
							{
								position71, tokenIndex71, depth71 := position, tokenIndex, depth
								if buffer[position] != rune(' ') {
									goto l72
								}
								position++
								goto l71
							l72:
								position, tokenIndex, depth = position71, tokenIndex71, depth71
								if buffer[position] != rune('\t') {
									goto l70
								}
								position++
							}
						l71:
							goto l69
						l70:
							position, tokenIndex, depth = position70, tokenIndex70, depth70
						}
						depth--
						add(rulePegText, position68)
					}
					if buffer[position] != rune('=') {
						goto l62
					}
					position++
					{
						position73 := position
						depth++
					l74:
						{
							position75, tokenIndex75, depth75 := position, tokenIndex, depth
							{
								position76, tokenIndex76, depth76 := position, tokenIndex, depth
								if buffer[position] != rune(' ') {
									goto l77
								}
								position++
								goto l76
							l77:
								position, tokenIndex, depth = position76, tokenIndex76, depth76
								if buffer[position] != rune('\t') {
									goto l75
								}
								position++
							}
						l76:
							goto l74
						l75:
							position, tokenIndex, depth = position75, tokenIndex75, depth75
						}
						depth--
						add(rulePegText, position73)
					}
					if !_rules[rulescalar]() {
						goto l62
					}
				l78:
					{
						position79, tokenIndex79, depth79 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l79
						}
						goto l78
					l79:
						position, tokenIndex, depth = position79, tokenIndex79, depth79
					}
					if buffer[position] != rune(';') {
						goto l62
					}
					position++
					if !_rules[rulespacing]() {
						goto l62
					}
					if !_rules[ruleAction5]() {
						goto l62
					}
				}
			l64:
				depth--
				add(rulefield_decl, position63)
			}
			return true
		l62:
			position, tokenIndex, depth = position62, tokenIndex62, depth62
			return false
		},
		/* 10 field_type <- <(ident Action6)> */
		func() bool {
			position80, tokenIndex80, depth80 := position, tokenIndex, depth
			{
				position81 := position
				depth++
				if !_rules[ruleident]() {
					goto l80
				}
				if !_rules[ruleAction6]() {
					goto l80
				}
				depth--
				add(rulefield_type, position81)
			}
			return true
		l80:
			position, tokenIndex, depth = position80, tokenIndex80, depth80
			return false
		},
		/* 11 enum_decl <- <(('e' 'n' 'u' 'm' spacing ident spacing metadata* '{' enum_fields '}' spacing) / ('e' 'n' 'u' 'm' spacing ident ':' type spacing metadata* '{' enum_fields '}' spacing))> */
		func() bool {
			position82, tokenIndex82, depth82 := position, tokenIndex, depth
			{
				position83 := position
				depth++
				{
					position84, tokenIndex84, depth84 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l85
					}
					position++
					if buffer[position] != rune('n') {
						goto l85
					}
					position++
					if buffer[position] != rune('u') {
						goto l85
					}
					position++
					if buffer[position] != rune('m') {
						goto l85
					}
					position++
					if !_rules[rulespacing]() {
						goto l85
					}
					if !_rules[ruleident]() {
						goto l85
					}
					if !_rules[rulespacing]() {
						goto l85
					}
				l86:
					{
						position87, tokenIndex87, depth87 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l87
						}
						goto l86
					l87:
						position, tokenIndex, depth = position87, tokenIndex87, depth87
					}
					if buffer[position] != rune('{') {
						goto l85
					}
					position++
					if !_rules[ruleenum_fields]() {
						goto l85
					}
					if buffer[position] != rune('}') {
						goto l85
					}
					position++
					if !_rules[rulespacing]() {
						goto l85
					}
					goto l84
				l85:
					position, tokenIndex, depth = position84, tokenIndex84, depth84
					if buffer[position] != rune('e') {
						goto l82
					}
					position++
					if buffer[position] != rune('n') {
						goto l82
					}
					position++
					if buffer[position] != rune('u') {
						goto l82
					}
					position++
					if buffer[position] != rune('m') {
						goto l82
					}
					position++
					if !_rules[rulespacing]() {
						goto l82
					}
					if !_rules[ruleident]() {
						goto l82
					}
					if buffer[position] != rune(':') {
						goto l82
					}
					position++
					if !_rules[ruletype]() {
						goto l82
					}
					if !_rules[rulespacing]() {
						goto l82
					}
				l88:
					{
						position89, tokenIndex89, depth89 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l89
						}
						goto l88
					l89:
						position, tokenIndex, depth = position89, tokenIndex89, depth89
					}
					if buffer[position] != rune('{') {
						goto l82
					}
					position++
					if !_rules[ruleenum_fields]() {
						goto l82
					}
					if buffer[position] != rune('}') {
						goto l82
					}
					position++
					if !_rules[rulespacing]() {
						goto l82
					}
				}
			l84:
				depth--
				add(ruleenum_decl, position83)
			}
			return true
		l82:
			position, tokenIndex, depth = position82, tokenIndex82, depth82
			return false
		},
		/* 12 enum_fields <- <((spacing enum_field ',') / (spacing enum_field))> */
		func() bool {
			position90, tokenIndex90, depth90 := position, tokenIndex, depth
			{
				position91 := position
				depth++
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if !_rules[rulespacing]() {
						goto l93
					}
					if !_rules[ruleenum_field]() {
						goto l93
					}
					if buffer[position] != rune(',') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if !_rules[rulespacing]() {
						goto l90
					}
					if !_rules[ruleenum_field]() {
						goto l90
					}
				}
			l92:
				depth--
				add(ruleenum_fields, position91)
			}
			return true
		l90:
			position, tokenIndex, depth = position90, tokenIndex90, depth90
			return false
		},
		/* 13 union_decl <- <('u' 'n' 'i' 'o' 'n' spacing union_name spacing metadata* '{' enum_fields+ '}' spacing Action7)> */
		func() bool {
			position94, tokenIndex94, depth94 := position, tokenIndex, depth
			{
				position95 := position
				depth++
				if buffer[position] != rune('u') {
					goto l94
				}
				position++
				if buffer[position] != rune('n') {
					goto l94
				}
				position++
				if buffer[position] != rune('i') {
					goto l94
				}
				position++
				if buffer[position] != rune('o') {
					goto l94
				}
				position++
				if buffer[position] != rune('n') {
					goto l94
				}
				position++
				if !_rules[rulespacing]() {
					goto l94
				}
				if !_rules[ruleunion_name]() {
					goto l94
				}
				if !_rules[rulespacing]() {
					goto l94
				}
			l96:
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					if !_rules[rulemetadata]() {
						goto l97
					}
					goto l96
				l97:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
				}
				if buffer[position] != rune('{') {
					goto l94
				}
				position++
				if !_rules[ruleenum_fields]() {
					goto l94
				}
			l98:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if !_rules[ruleenum_fields]() {
						goto l99
					}
					goto l98
				l99:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
				}
				if buffer[position] != rune('}') {
					goto l94
				}
				position++
				if !_rules[rulespacing]() {
					goto l94
				}
				if !_rules[ruleAction7]() {
					goto l94
				}
				depth--
				add(ruleunion_decl, position95)
			}
			return true
		l94:
			position, tokenIndex, depth = position94, tokenIndex94, depth94
			return false
		},
		/* 14 union_name <- <(ident Action8)> */
		func() bool {
			position100, tokenIndex100, depth100 := position, tokenIndex, depth
			{
				position101 := position
				depth++
				if !_rules[ruleident]() {
					goto l100
				}
				if !_rules[ruleAction8]() {
					goto l100
				}
				depth--
				add(ruleunion_name, position101)
			}
			return true
		l100:
			position, tokenIndex, depth = position100, tokenIndex100, depth100
			return false
		},
		/* 15 enum_field <- <((enum_field_type spacing Action9) / (enum_field_type spacing '=' spacing integer_constant spacing))> */
		func() bool {
			position102, tokenIndex102, depth102 := position, tokenIndex, depth
			{
				position103 := position
				depth++
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if !_rules[ruleenum_field_type]() {
						goto l105
					}
					if !_rules[rulespacing]() {
						goto l105
					}
					if !_rules[ruleAction9]() {
						goto l105
					}
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if !_rules[ruleenum_field_type]() {
						goto l102
					}
					if !_rules[rulespacing]() {
						goto l102
					}
					if buffer[position] != rune('=') {
						goto l102
					}
					position++
					if !_rules[rulespacing]() {
						goto l102
					}
					if !_rules[ruleinteger_constant]() {
						goto l102
					}
					if !_rules[rulespacing]() {
						goto l102
					}
				}
			l104:
				depth--
				add(ruleenum_field, position103)
			}
			return true
		l102:
			position, tokenIndex, depth = position102, tokenIndex102, depth102
			return false
		},
		/* 16 enum_field_type <- <(ident Action10)> */
		func() bool {
			position106, tokenIndex106, depth106 := position, tokenIndex, depth
			{
				position107 := position
				depth++
				if !_rules[ruleident]() {
					goto l106
				}
				if !_rules[ruleAction10]() {
					goto l106
				}
				depth--
				add(ruleenum_field_type, position107)
			}
			return true
		l106:
			position, tokenIndex, depth = position106, tokenIndex106, depth106
			return false
		},
		/* 17 root_decl <- <('r' 'o' 'o' 't' '_' 't' 'y' 'p' 'e' spacing <([A-z] / [0-9] / '_' / '.' / '-')+> Action11 ';' spacing)> */
		func() bool {
			position108, tokenIndex108, depth108 := position, tokenIndex, depth
			{
				position109 := position
				depth++
				if buffer[position] != rune('r') {
					goto l108
				}
				position++
				if buffer[position] != rune('o') {
					goto l108
				}
				position++
				if buffer[position] != rune('o') {
					goto l108
				}
				position++
				if buffer[position] != rune('t') {
					goto l108
				}
				position++
				if buffer[position] != rune('_') {
					goto l108
				}
				position++
				if buffer[position] != rune('t') {
					goto l108
				}
				position++
				if buffer[position] != rune('y') {
					goto l108
				}
				position++
				if buffer[position] != rune('p') {
					goto l108
				}
				position++
				if buffer[position] != rune('e') {
					goto l108
				}
				position++
				if !_rules[rulespacing]() {
					goto l108
				}
				{
					position110 := position
					depth++
					{
						position113, tokenIndex113, depth113 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('A') || c > rune('z') {
							goto l114
						}
						position++
						goto l113
					l114:
						position, tokenIndex, depth = position113, tokenIndex113, depth113
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l115
						}
						position++
						goto l113
					l115:
						position, tokenIndex, depth = position113, tokenIndex113, depth113
						if buffer[position] != rune('_') {
							goto l116
						}
						position++
						goto l113
					l116:
						position, tokenIndex, depth = position113, tokenIndex113, depth113
						if buffer[position] != rune('.') {
							goto l117
						}
						position++
						goto l113
					l117:
						position, tokenIndex, depth = position113, tokenIndex113, depth113
						if buffer[position] != rune('-') {
							goto l108
						}
						position++
					}
				l113:
				l111:
					{
						position112, tokenIndex112, depth112 := position, tokenIndex, depth
						{
							position118, tokenIndex118, depth118 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('A') || c > rune('z') {
								goto l119
							}
							position++
							goto l118
						l119:
							position, tokenIndex, depth = position118, tokenIndex118, depth118
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l120
							}
							position++
							goto l118
						l120:
							position, tokenIndex, depth = position118, tokenIndex118, depth118
							if buffer[position] != rune('_') {
								goto l121
							}
							position++
							goto l118
						l121:
							position, tokenIndex, depth = position118, tokenIndex118, depth118
							if buffer[position] != rune('.') {
								goto l122
							}
							position++
							goto l118
						l122:
							position, tokenIndex, depth = position118, tokenIndex118, depth118
							if buffer[position] != rune('-') {
								goto l112
							}
							position++
						}
					l118:
						goto l111
					l112:
						position, tokenIndex, depth = position112, tokenIndex112, depth112
					}
					depth--
					add(rulePegText, position110)
				}
				if !_rules[ruleAction11]() {
					goto l108
				}
				if buffer[position] != rune(';') {
					goto l108
				}
				position++
				if !_rules[rulespacing]() {
					goto l108
				}
				depth--
				add(ruleroot_decl, position109)
			}
			return true
		l108:
			position, tokenIndex, depth = position108, tokenIndex108, depth108
			return false
		},
		/* 18 file_extension_decl <- <('f' 'i' 'l' 'e' '_' 'e' 'x' 't' 'e' 'n' 's' 'i' 'o' 'n' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position123, tokenIndex123, depth123 := position, tokenIndex, depth
			{
				position124 := position
				depth++
				if buffer[position] != rune('f') {
					goto l123
				}
				position++
				if buffer[position] != rune('i') {
					goto l123
				}
				position++
				if buffer[position] != rune('l') {
					goto l123
				}
				position++
				if buffer[position] != rune('e') {
					goto l123
				}
				position++
				if buffer[position] != rune('_') {
					goto l123
				}
				position++
				if buffer[position] != rune('e') {
					goto l123
				}
				position++
				if buffer[position] != rune('x') {
					goto l123
				}
				position++
				if buffer[position] != rune('t') {
					goto l123
				}
				position++
				if buffer[position] != rune('e') {
					goto l123
				}
				position++
				if buffer[position] != rune('n') {
					goto l123
				}
				position++
				if buffer[position] != rune('s') {
					goto l123
				}
				position++
				if buffer[position] != rune('i') {
					goto l123
				}
				position++
				if buffer[position] != rune('o') {
					goto l123
				}
				position++
				if buffer[position] != rune('n') {
					goto l123
				}
				position++
				{
					position125 := position
					depth++
				l126:
					{
						position127, tokenIndex127, depth127 := position, tokenIndex, depth
						{
							position128, tokenIndex128, depth128 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l129
							}
							position++
							goto l128
						l129:
							position, tokenIndex, depth = position128, tokenIndex128, depth128
							if buffer[position] != rune('\t') {
								goto l127
							}
							position++
						}
					l128:
						goto l126
					l127:
						position, tokenIndex, depth = position127, tokenIndex127, depth127
					}
					depth--
					add(rulePegText, position125)
				}
				{
					position130 := position
					depth++
					{
						position133, tokenIndex133, depth133 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l133
						}
						position++
						goto l123
					l133:
						position, tokenIndex, depth = position133, tokenIndex133, depth133
					}
					if !matchDot() {
						goto l123
					}
				l131:
					{
						position132, tokenIndex132, depth132 := position, tokenIndex, depth
						{
							position134, tokenIndex134, depth134 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l134
							}
							position++
							goto l132
						l134:
							position, tokenIndex, depth = position134, tokenIndex134, depth134
						}
						if !matchDot() {
							goto l132
						}
						goto l131
					l132:
						position, tokenIndex, depth = position132, tokenIndex132, depth132
					}
					depth--
					add(rulePegText, position130)
				}
				if buffer[position] != rune(';') {
					goto l123
				}
				position++
				if !_rules[rulespacing]() {
					goto l123
				}
				depth--
				add(rulefile_extension_decl, position124)
			}
			return true
		l123:
			position, tokenIndex, depth = position123, tokenIndex123, depth123
			return false
		},
		/* 19 file_identifier_decl <- <('f' 'i' 'l' 'e' '_' 'i' 'd' 'e' 'n' 't' 'i' 'f' 'i' 'e' 'r' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position135, tokenIndex135, depth135 := position, tokenIndex, depth
			{
				position136 := position
				depth++
				if buffer[position] != rune('f') {
					goto l135
				}
				position++
				if buffer[position] != rune('i') {
					goto l135
				}
				position++
				if buffer[position] != rune('l') {
					goto l135
				}
				position++
				if buffer[position] != rune('e') {
					goto l135
				}
				position++
				if buffer[position] != rune('_') {
					goto l135
				}
				position++
				if buffer[position] != rune('i') {
					goto l135
				}
				position++
				if buffer[position] != rune('d') {
					goto l135
				}
				position++
				if buffer[position] != rune('e') {
					goto l135
				}
				position++
				if buffer[position] != rune('n') {
					goto l135
				}
				position++
				if buffer[position] != rune('t') {
					goto l135
				}
				position++
				if buffer[position] != rune('i') {
					goto l135
				}
				position++
				if buffer[position] != rune('f') {
					goto l135
				}
				position++
				if buffer[position] != rune('i') {
					goto l135
				}
				position++
				if buffer[position] != rune('e') {
					goto l135
				}
				position++
				if buffer[position] != rune('r') {
					goto l135
				}
				position++
				{
					position137 := position
					depth++
				l138:
					{
						position139, tokenIndex139, depth139 := position, tokenIndex, depth
						{
							position140, tokenIndex140, depth140 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l141
							}
							position++
							goto l140
						l141:
							position, tokenIndex, depth = position140, tokenIndex140, depth140
							if buffer[position] != rune('\t') {
								goto l139
							}
							position++
						}
					l140:
						goto l138
					l139:
						position, tokenIndex, depth = position139, tokenIndex139, depth139
					}
					depth--
					add(rulePegText, position137)
				}
				{
					position142 := position
					depth++
					{
						position145, tokenIndex145, depth145 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l145
						}
						position++
						goto l135
					l145:
						position, tokenIndex, depth = position145, tokenIndex145, depth145
					}
					if !matchDot() {
						goto l135
					}
				l143:
					{
						position144, tokenIndex144, depth144 := position, tokenIndex, depth
						{
							position146, tokenIndex146, depth146 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l146
							}
							position++
							goto l144
						l146:
							position, tokenIndex, depth = position146, tokenIndex146, depth146
						}
						if !matchDot() {
							goto l144
						}
						goto l143
					l144:
						position, tokenIndex, depth = position144, tokenIndex144, depth144
					}
					depth--
					add(rulePegText, position142)
				}
				if buffer[position] != rune(';') {
					goto l135
				}
				position++
				if !_rules[rulespacing]() {
					goto l135
				}
				depth--
				add(rulefile_identifier_decl, position136)
			}
			return true
		l135:
			position, tokenIndex, depth = position135, tokenIndex135, depth135
			return false
		},
		/* 20 attribute_decl <- <('a' 't' 't' 'r' 'i' 'b' 'u' 't' 'e' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position147, tokenIndex147, depth147 := position, tokenIndex, depth
			{
				position148 := position
				depth++
				if buffer[position] != rune('a') {
					goto l147
				}
				position++
				if buffer[position] != rune('t') {
					goto l147
				}
				position++
				if buffer[position] != rune('t') {
					goto l147
				}
				position++
				if buffer[position] != rune('r') {
					goto l147
				}
				position++
				if buffer[position] != rune('i') {
					goto l147
				}
				position++
				if buffer[position] != rune('b') {
					goto l147
				}
				position++
				if buffer[position] != rune('u') {
					goto l147
				}
				position++
				if buffer[position] != rune('t') {
					goto l147
				}
				position++
				if buffer[position] != rune('e') {
					goto l147
				}
				position++
				{
					position149 := position
					depth++
				l150:
					{
						position151, tokenIndex151, depth151 := position, tokenIndex, depth
						{
							position152, tokenIndex152, depth152 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l153
							}
							position++
							goto l152
						l153:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('\t') {
								goto l151
							}
							position++
						}
					l152:
						goto l150
					l151:
						position, tokenIndex, depth = position151, tokenIndex151, depth151
					}
					depth--
					add(rulePegText, position149)
				}
				{
					position154 := position
					depth++
					{
						position157, tokenIndex157, depth157 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l157
						}
						position++
						goto l147
					l157:
						position, tokenIndex, depth = position157, tokenIndex157, depth157
					}
					if !matchDot() {
						goto l147
					}
				l155:
					{
						position156, tokenIndex156, depth156 := position, tokenIndex, depth
						{
							position158, tokenIndex158, depth158 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l158
							}
							position++
							goto l156
						l158:
							position, tokenIndex, depth = position158, tokenIndex158, depth158
						}
						if !matchDot() {
							goto l156
						}
						goto l155
					l156:
						position, tokenIndex, depth = position156, tokenIndex156, depth156
					}
					depth--
					add(rulePegText, position154)
				}
				if buffer[position] != rune(';') {
					goto l147
				}
				position++
				if !_rules[rulespacing]() {
					goto l147
				}
				depth--
				add(ruleattribute_decl, position148)
			}
			return true
		l147:
			position, tokenIndex, depth = position147, tokenIndex147, depth147
			return false
		},
		/* 21 rpc_decl <- <('r' 'p' 'c' '_' 's' 'e' 'r' 'v' 'i' 'c' 'e' <(' ' / '\t')*> ident '{' <(!'}' .)+> '}' spacing)> */
		func() bool {
			position159, tokenIndex159, depth159 := position, tokenIndex, depth
			{
				position160 := position
				depth++
				if buffer[position] != rune('r') {
					goto l159
				}
				position++
				if buffer[position] != rune('p') {
					goto l159
				}
				position++
				if buffer[position] != rune('c') {
					goto l159
				}
				position++
				if buffer[position] != rune('_') {
					goto l159
				}
				position++
				if buffer[position] != rune('s') {
					goto l159
				}
				position++
				if buffer[position] != rune('e') {
					goto l159
				}
				position++
				if buffer[position] != rune('r') {
					goto l159
				}
				position++
				if buffer[position] != rune('v') {
					goto l159
				}
				position++
				if buffer[position] != rune('i') {
					goto l159
				}
				position++
				if buffer[position] != rune('c') {
					goto l159
				}
				position++
				if buffer[position] != rune('e') {
					goto l159
				}
				position++
				{
					position161 := position
					depth++
				l162:
					{
						position163, tokenIndex163, depth163 := position, tokenIndex, depth
						{
							position164, tokenIndex164, depth164 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l165
							}
							position++
							goto l164
						l165:
							position, tokenIndex, depth = position164, tokenIndex164, depth164
							if buffer[position] != rune('\t') {
								goto l163
							}
							position++
						}
					l164:
						goto l162
					l163:
						position, tokenIndex, depth = position163, tokenIndex163, depth163
					}
					depth--
					add(rulePegText, position161)
				}
				if !_rules[ruleident]() {
					goto l159
				}
				if buffer[position] != rune('{') {
					goto l159
				}
				position++
				{
					position166 := position
					depth++
					{
						position169, tokenIndex169, depth169 := position, tokenIndex, depth
						if buffer[position] != rune('}') {
							goto l169
						}
						position++
						goto l159
					l169:
						position, tokenIndex, depth = position169, tokenIndex169, depth169
					}
					if !matchDot() {
						goto l159
					}
				l167:
					{
						position168, tokenIndex168, depth168 := position, tokenIndex, depth
						{
							position170, tokenIndex170, depth170 := position, tokenIndex, depth
							if buffer[position] != rune('}') {
								goto l170
							}
							position++
							goto l168
						l170:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
						}
						if !matchDot() {
							goto l168
						}
						goto l167
					l168:
						position, tokenIndex, depth = position168, tokenIndex168, depth168
					}
					depth--
					add(rulePegText, position166)
				}
				if buffer[position] != rune('}') {
					goto l159
				}
				position++
				if !_rules[rulespacing]() {
					goto l159
				}
				depth--
				add(rulerpc_decl, position160)
			}
			return true
		l159:
			position, tokenIndex, depth = position159, tokenIndex159, depth159
			return false
		},
		/* 22 type <- <(('b' 'o' 'o' 'l' spacing Action12) / ('i' 'n' 't' '8' spacing Action13) / ('u' 'i' 'n' 't' '8' spacing Action14) / ('i' 'n' 't' '1' '6' spacing Action15) / ('u' 'i' 'n' 't' '1' '6' spacing Action16) / ('i' 'n' 't' '3' '2' spacing Action17) / ('u' 'i' 'n' 't' '3' '2' spacing Action18) / ('i' 'n' 't' '6' '4' spacing Action19) / ('u' 'i' 'n' 't' '6' '4' spacing Action20) / ('f' 'l' 'o' 'a' 't' '3' '2' spacing Action21) / ('f' 'l' 'o' 'a' 't' '6' '4' spacing Action22) / ('b' 'y' 't' 'e' spacing Action23) / ('u' 'b' 'y' 't' 'e' spacing Action24) / ('s' 'h' 'o' 'r' 't' spacing Action25) / ('u' 's' 'h' 'o' 'r' 't' spacing Action26) / ('i' 'n' 't' spacing Action27) / ('u' 'i' 'n' 't' spacing Action28) / ('f' 'l' 'o' 'a' 't' spacing Action29) / ('l' 'o' 'n' 'g' spacing Action30) / ('u' 'l' 'o' 'n' 'g' spacing Action31) / ('d' 'o' 'u' 'b' 'l' 'e' spacing Action32) / ('s' 't' 'r' 'i' 'n' 'g' spacing Action33) / (ident spacing Action34) / ('[' type ']' spacing Action35))> */
		func() bool {
			position171, tokenIndex171, depth171 := position, tokenIndex, depth
			{
				position172 := position
				depth++
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l174
					}
					position++
					if buffer[position] != rune('o') {
						goto l174
					}
					position++
					if buffer[position] != rune('o') {
						goto l174
					}
					position++
					if buffer[position] != rune('l') {
						goto l174
					}
					position++
					if !_rules[rulespacing]() {
						goto l174
					}
					if !_rules[ruleAction12]() {
						goto l174
					}
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('i') {
						goto l175
					}
					position++
					if buffer[position] != rune('n') {
						goto l175
					}
					position++
					if buffer[position] != rune('t') {
						goto l175
					}
					position++
					if buffer[position] != rune('8') {
						goto l175
					}
					position++
					if !_rules[rulespacing]() {
						goto l175
					}
					if !_rules[ruleAction13]() {
						goto l175
					}
					goto l173
				l175:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('u') {
						goto l176
					}
					position++
					if buffer[position] != rune('i') {
						goto l176
					}
					position++
					if buffer[position] != rune('n') {
						goto l176
					}
					position++
					if buffer[position] != rune('t') {
						goto l176
					}
					position++
					if buffer[position] != rune('8') {
						goto l176
					}
					position++
					if !_rules[rulespacing]() {
						goto l176
					}
					if !_rules[ruleAction14]() {
						goto l176
					}
					goto l173
				l176:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('i') {
						goto l177
					}
					position++
					if buffer[position] != rune('n') {
						goto l177
					}
					position++
					if buffer[position] != rune('t') {
						goto l177
					}
					position++
					if buffer[position] != rune('1') {
						goto l177
					}
					position++
					if buffer[position] != rune('6') {
						goto l177
					}
					position++
					if !_rules[rulespacing]() {
						goto l177
					}
					if !_rules[ruleAction15]() {
						goto l177
					}
					goto l173
				l177:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('u') {
						goto l178
					}
					position++
					if buffer[position] != rune('i') {
						goto l178
					}
					position++
					if buffer[position] != rune('n') {
						goto l178
					}
					position++
					if buffer[position] != rune('t') {
						goto l178
					}
					position++
					if buffer[position] != rune('1') {
						goto l178
					}
					position++
					if buffer[position] != rune('6') {
						goto l178
					}
					position++
					if !_rules[rulespacing]() {
						goto l178
					}
					if !_rules[ruleAction16]() {
						goto l178
					}
					goto l173
				l178:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('i') {
						goto l179
					}
					position++
					if buffer[position] != rune('n') {
						goto l179
					}
					position++
					if buffer[position] != rune('t') {
						goto l179
					}
					position++
					if buffer[position] != rune('3') {
						goto l179
					}
					position++
					if buffer[position] != rune('2') {
						goto l179
					}
					position++
					if !_rules[rulespacing]() {
						goto l179
					}
					if !_rules[ruleAction17]() {
						goto l179
					}
					goto l173
				l179:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('u') {
						goto l180
					}
					position++
					if buffer[position] != rune('i') {
						goto l180
					}
					position++
					if buffer[position] != rune('n') {
						goto l180
					}
					position++
					if buffer[position] != rune('t') {
						goto l180
					}
					position++
					if buffer[position] != rune('3') {
						goto l180
					}
					position++
					if buffer[position] != rune('2') {
						goto l180
					}
					position++
					if !_rules[rulespacing]() {
						goto l180
					}
					if !_rules[ruleAction18]() {
						goto l180
					}
					goto l173
				l180:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('i') {
						goto l181
					}
					position++
					if buffer[position] != rune('n') {
						goto l181
					}
					position++
					if buffer[position] != rune('t') {
						goto l181
					}
					position++
					if buffer[position] != rune('6') {
						goto l181
					}
					position++
					if buffer[position] != rune('4') {
						goto l181
					}
					position++
					if !_rules[rulespacing]() {
						goto l181
					}
					if !_rules[ruleAction19]() {
						goto l181
					}
					goto l173
				l181:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('u') {
						goto l182
					}
					position++
					if buffer[position] != rune('i') {
						goto l182
					}
					position++
					if buffer[position] != rune('n') {
						goto l182
					}
					position++
					if buffer[position] != rune('t') {
						goto l182
					}
					position++
					if buffer[position] != rune('6') {
						goto l182
					}
					position++
					if buffer[position] != rune('4') {
						goto l182
					}
					position++
					if !_rules[rulespacing]() {
						goto l182
					}
					if !_rules[ruleAction20]() {
						goto l182
					}
					goto l173
				l182:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('f') {
						goto l183
					}
					position++
					if buffer[position] != rune('l') {
						goto l183
					}
					position++
					if buffer[position] != rune('o') {
						goto l183
					}
					position++
					if buffer[position] != rune('a') {
						goto l183
					}
					position++
					if buffer[position] != rune('t') {
						goto l183
					}
					position++
					if buffer[position] != rune('3') {
						goto l183
					}
					position++
					if buffer[position] != rune('2') {
						goto l183
					}
					position++
					if !_rules[rulespacing]() {
						goto l183
					}
					if !_rules[ruleAction21]() {
						goto l183
					}
					goto l173
				l183:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('f') {
						goto l184
					}
					position++
					if buffer[position] != rune('l') {
						goto l184
					}
					position++
					if buffer[position] != rune('o') {
						goto l184
					}
					position++
					if buffer[position] != rune('a') {
						goto l184
					}
					position++
					if buffer[position] != rune('t') {
						goto l184
					}
					position++
					if buffer[position] != rune('6') {
						goto l184
					}
					position++
					if buffer[position] != rune('4') {
						goto l184
					}
					position++
					if !_rules[rulespacing]() {
						goto l184
					}
					if !_rules[ruleAction22]() {
						goto l184
					}
					goto l173
				l184:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('b') {
						goto l185
					}
					position++
					if buffer[position] != rune('y') {
						goto l185
					}
					position++
					if buffer[position] != rune('t') {
						goto l185
					}
					position++
					if buffer[position] != rune('e') {
						goto l185
					}
					position++
					if !_rules[rulespacing]() {
						goto l185
					}
					if !_rules[ruleAction23]() {
						goto l185
					}
					goto l173
				l185:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('u') {
						goto l186
					}
					position++
					if buffer[position] != rune('b') {
						goto l186
					}
					position++
					if buffer[position] != rune('y') {
						goto l186
					}
					position++
					if buffer[position] != rune('t') {
						goto l186
					}
					position++
					if buffer[position] != rune('e') {
						goto l186
					}
					position++
					if !_rules[rulespacing]() {
						goto l186
					}
					if !_rules[ruleAction24]() {
						goto l186
					}
					goto l173
				l186:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('s') {
						goto l187
					}
					position++
					if buffer[position] != rune('h') {
						goto l187
					}
					position++
					if buffer[position] != rune('o') {
						goto l187
					}
					position++
					if buffer[position] != rune('r') {
						goto l187
					}
					position++
					if buffer[position] != rune('t') {
						goto l187
					}
					position++
					if !_rules[rulespacing]() {
						goto l187
					}
					if !_rules[ruleAction25]() {
						goto l187
					}
					goto l173
				l187:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('u') {
						goto l188
					}
					position++
					if buffer[position] != rune('s') {
						goto l188
					}
					position++
					if buffer[position] != rune('h') {
						goto l188
					}
					position++
					if buffer[position] != rune('o') {
						goto l188
					}
					position++
					if buffer[position] != rune('r') {
						goto l188
					}
					position++
					if buffer[position] != rune('t') {
						goto l188
					}
					position++
					if !_rules[rulespacing]() {
						goto l188
					}
					if !_rules[ruleAction26]() {
						goto l188
					}
					goto l173
				l188:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('i') {
						goto l189
					}
					position++
					if buffer[position] != rune('n') {
						goto l189
					}
					position++
					if buffer[position] != rune('t') {
						goto l189
					}
					position++
					if !_rules[rulespacing]() {
						goto l189
					}
					if !_rules[ruleAction27]() {
						goto l189
					}
					goto l173
				l189:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('u') {
						goto l190
					}
					position++
					if buffer[position] != rune('i') {
						goto l190
					}
					position++
					if buffer[position] != rune('n') {
						goto l190
					}
					position++
					if buffer[position] != rune('t') {
						goto l190
					}
					position++
					if !_rules[rulespacing]() {
						goto l190
					}
					if !_rules[ruleAction28]() {
						goto l190
					}
					goto l173
				l190:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('f') {
						goto l191
					}
					position++
					if buffer[position] != rune('l') {
						goto l191
					}
					position++
					if buffer[position] != rune('o') {
						goto l191
					}
					position++
					if buffer[position] != rune('a') {
						goto l191
					}
					position++
					if buffer[position] != rune('t') {
						goto l191
					}
					position++
					if !_rules[rulespacing]() {
						goto l191
					}
					if !_rules[ruleAction29]() {
						goto l191
					}
					goto l173
				l191:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('l') {
						goto l192
					}
					position++
					if buffer[position] != rune('o') {
						goto l192
					}
					position++
					if buffer[position] != rune('n') {
						goto l192
					}
					position++
					if buffer[position] != rune('g') {
						goto l192
					}
					position++
					if !_rules[rulespacing]() {
						goto l192
					}
					if !_rules[ruleAction30]() {
						goto l192
					}
					goto l173
				l192:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('u') {
						goto l193
					}
					position++
					if buffer[position] != rune('l') {
						goto l193
					}
					position++
					if buffer[position] != rune('o') {
						goto l193
					}
					position++
					if buffer[position] != rune('n') {
						goto l193
					}
					position++
					if buffer[position] != rune('g') {
						goto l193
					}
					position++
					if !_rules[rulespacing]() {
						goto l193
					}
					if !_rules[ruleAction31]() {
						goto l193
					}
					goto l173
				l193:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('d') {
						goto l194
					}
					position++
					if buffer[position] != rune('o') {
						goto l194
					}
					position++
					if buffer[position] != rune('u') {
						goto l194
					}
					position++
					if buffer[position] != rune('b') {
						goto l194
					}
					position++
					if buffer[position] != rune('l') {
						goto l194
					}
					position++
					if buffer[position] != rune('e') {
						goto l194
					}
					position++
					if !_rules[rulespacing]() {
						goto l194
					}
					if !_rules[ruleAction32]() {
						goto l194
					}
					goto l173
				l194:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('s') {
						goto l195
					}
					position++
					if buffer[position] != rune('t') {
						goto l195
					}
					position++
					if buffer[position] != rune('r') {
						goto l195
					}
					position++
					if buffer[position] != rune('i') {
						goto l195
					}
					position++
					if buffer[position] != rune('n') {
						goto l195
					}
					position++
					if buffer[position] != rune('g') {
						goto l195
					}
					position++
					if !_rules[rulespacing]() {
						goto l195
					}
					if !_rules[ruleAction33]() {
						goto l195
					}
					goto l173
				l195:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if !_rules[ruleident]() {
						goto l196
					}
					if !_rules[rulespacing]() {
						goto l196
					}
					if !_rules[ruleAction34]() {
						goto l196
					}
					goto l173
				l196:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('[') {
						goto l171
					}
					position++
					if !_rules[ruletype]() {
						goto l171
					}
					if buffer[position] != rune(']') {
						goto l171
					}
					position++
					if !_rules[rulespacing]() {
						goto l171
					}
					if !_rules[ruleAction35]() {
						goto l171
					}
				}
			l173:
				depth--
				add(ruletype, position172)
			}
			return true
		l171:
			position, tokenIndex, depth = position171, tokenIndex171, depth171
			return false
		},
		/* 23 scalar <- <(integer_constant / float_constant)> */
		func() bool {
			position197, tokenIndex197, depth197 := position, tokenIndex, depth
			{
				position198 := position
				depth++
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if !_rules[ruleinteger_constant]() {
						goto l200
					}
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if !_rules[rulefloat_constant]() {
						goto l197
					}
				}
			l199:
				depth--
				add(rulescalar, position198)
			}
			return true
		l197:
			position, tokenIndex, depth = position197, tokenIndex197, depth197
			return false
		},
		/* 24 integer_constant <- <(<[0-9]+> / ('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		func() bool {
			position201, tokenIndex201, depth201 := position, tokenIndex, depth
			{
				position202 := position
				depth++
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					{
						position205 := position
						depth++
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l204
						}
						position++
					l206:
						{
							position207, tokenIndex207, depth207 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l207
							}
							position++
							goto l206
						l207:
							position, tokenIndex, depth = position207, tokenIndex207, depth207
						}
						depth--
						add(rulePegText, position205)
					}
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('t') {
						goto l208
					}
					position++
					if buffer[position] != rune('r') {
						goto l208
					}
					position++
					if buffer[position] != rune('u') {
						goto l208
					}
					position++
					if buffer[position] != rune('e') {
						goto l208
					}
					position++
					goto l203
				l208:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('f') {
						goto l201
					}
					position++
					if buffer[position] != rune('a') {
						goto l201
					}
					position++
					if buffer[position] != rune('l') {
						goto l201
					}
					position++
					if buffer[position] != rune('s') {
						goto l201
					}
					position++
					if buffer[position] != rune('e') {
						goto l201
					}
					position++
				}
			l203:
				depth--
				add(ruleinteger_constant, position202)
			}
			return true
		l201:
			position, tokenIndex, depth = position201, tokenIndex201, depth201
			return false
		},
		/* 25 float_constant <- <(<('-'* [0-9]+ . [0-9])> / float_constant_exp)> */
		func() bool {
			position209, tokenIndex209, depth209 := position, tokenIndex, depth
			{
				position210 := position
				depth++
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					{
						position213 := position
						depth++
					l214:
						{
							position215, tokenIndex215, depth215 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l215
							}
							position++
							goto l214
						l215:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l212
						}
						position++
					l216:
						{
							position217, tokenIndex217, depth217 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l217
							}
							position++
							goto l216
						l217:
							position, tokenIndex, depth = position217, tokenIndex217, depth217
						}
						if !matchDot() {
							goto l212
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l212
						}
						position++
						depth--
						add(rulePegText, position213)
					}
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if !_rules[rulefloat_constant_exp]() {
						goto l209
					}
				}
			l211:
				depth--
				add(rulefloat_constant, position210)
			}
			return true
		l209:
			position, tokenIndex, depth = position209, tokenIndex209, depth209
			return false
		},
		/* 26 float_constant_exp <- <(<('-'* [0-9]+ . [0-9]+)> <('e' / 'E')> <([+-]] / '>' / ' ' / '<' / '[' / [0-9])+>)> */
		func() bool {
			position218, tokenIndex218, depth218 := position, tokenIndex, depth
			{
				position219 := position
				depth++
				{
					position220 := position
					depth++
				l221:
					{
						position222, tokenIndex222, depth222 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l222
						}
						position++
						goto l221
					l222:
						position, tokenIndex, depth = position222, tokenIndex222, depth222
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l218
					}
					position++
				l223:
					{
						position224, tokenIndex224, depth224 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l224
						}
						position++
						goto l223
					l224:
						position, tokenIndex, depth = position224, tokenIndex224, depth224
					}
					if !matchDot() {
						goto l218
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l218
					}
					position++
				l225:
					{
						position226, tokenIndex226, depth226 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l226
						}
						position++
						goto l225
					l226:
						position, tokenIndex, depth = position226, tokenIndex226, depth226
					}
					depth--
					add(rulePegText, position220)
				}
				{
					position227 := position
					depth++
					{
						position228, tokenIndex228, depth228 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l229
						}
						position++
						goto l228
					l229:
						position, tokenIndex, depth = position228, tokenIndex228, depth228
						if buffer[position] != rune('E') {
							goto l218
						}
						position++
					}
				l228:
					depth--
					add(rulePegText, position227)
				}
				{
					position230 := position
					depth++
					{
						position233, tokenIndex233, depth233 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('+') || c > rune(']') {
							goto l234
						}
						position++
						goto l233
					l234:
						position, tokenIndex, depth = position233, tokenIndex233, depth233
						if buffer[position] != rune('>') {
							goto l235
						}
						position++
						goto l233
					l235:
						position, tokenIndex, depth = position233, tokenIndex233, depth233
						if buffer[position] != rune(' ') {
							goto l236
						}
						position++
						goto l233
					l236:
						position, tokenIndex, depth = position233, tokenIndex233, depth233
						if buffer[position] != rune('<') {
							goto l237
						}
						position++
						goto l233
					l237:
						position, tokenIndex, depth = position233, tokenIndex233, depth233
						if buffer[position] != rune('[') {
							goto l238
						}
						position++
						goto l233
					l238:
						position, tokenIndex, depth = position233, tokenIndex233, depth233
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l218
						}
						position++
					}
				l233:
				l231:
					{
						position232, tokenIndex232, depth232 := position, tokenIndex, depth
						{
							position239, tokenIndex239, depth239 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('+') || c > rune(']') {
								goto l240
							}
							position++
							goto l239
						l240:
							position, tokenIndex, depth = position239, tokenIndex239, depth239
							if buffer[position] != rune('>') {
								goto l241
							}
							position++
							goto l239
						l241:
							position, tokenIndex, depth = position239, tokenIndex239, depth239
							if buffer[position] != rune(' ') {
								goto l242
							}
							position++
							goto l239
						l242:
							position, tokenIndex, depth = position239, tokenIndex239, depth239
							if buffer[position] != rune('<') {
								goto l243
							}
							position++
							goto l239
						l243:
							position, tokenIndex, depth = position239, tokenIndex239, depth239
							if buffer[position] != rune('[') {
								goto l244
							}
							position++
							goto l239
						l244:
							position, tokenIndex, depth = position239, tokenIndex239, depth239
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l232
							}
							position++
						}
					l239:
						goto l231
					l232:
						position, tokenIndex, depth = position232, tokenIndex232, depth232
					}
					depth--
					add(rulePegText, position230)
				}
				depth--
				add(rulefloat_constant_exp, position219)
			}
			return true
		l218:
			position, tokenIndex, depth = position218, tokenIndex218, depth218
			return false
		},
		/* 27 ident <- <<(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position245, tokenIndex245, depth245 := position, tokenIndex, depth
			{
				position246 := position
				depth++
				{
					position247 := position
					depth++
					{
						position248, tokenIndex248, depth248 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l249
						}
						position++
						goto l248
					l249:
						position, tokenIndex, depth = position248, tokenIndex248, depth248
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l250
						}
						position++
						goto l248
					l250:
						position, tokenIndex, depth = position248, tokenIndex248, depth248
						if buffer[position] != rune('_') {
							goto l245
						}
						position++
					}
				l248:
				l251:
					{
						position252, tokenIndex252, depth252 := position, tokenIndex, depth
						{
							position253, tokenIndex253, depth253 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l254
							}
							position++
							goto l253
						l254:
							position, tokenIndex, depth = position253, tokenIndex253, depth253
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l255
							}
							position++
							goto l253
						l255:
							position, tokenIndex, depth = position253, tokenIndex253, depth253
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l256
							}
							position++
							goto l253
						l256:
							position, tokenIndex, depth = position253, tokenIndex253, depth253
							if buffer[position] != rune('_') {
								goto l252
							}
							position++
						}
					l253:
						goto l251
					l252:
						position, tokenIndex, depth = position252, tokenIndex252, depth252
					}
					depth--
					add(rulePegText, position247)
				}
				depth--
				add(ruleident, position246)
			}
			return true
		l245:
			position, tokenIndex, depth = position245, tokenIndex245, depth245
			return false
		},
		/* 28 only_comment <- <(spacing ';')> */
		func() bool {
			position257, tokenIndex257, depth257 := position, tokenIndex, depth
			{
				position258 := position
				depth++
				if !_rules[rulespacing]() {
					goto l257
				}
				if buffer[position] != rune(';') {
					goto l257
				}
				position++
				depth--
				add(ruleonly_comment, position258)
			}
			return true
		l257:
			position, tokenIndex, depth = position257, tokenIndex257, depth257
			return false
		},
		/* 29 spacing <- <space_comment*> */
		func() bool {
			{
				position260 := position
				depth++
			l261:
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if !_rules[rulespace_comment]() {
						goto l262
					}
					goto l261
				l262:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
				}
				depth--
				add(rulespacing, position260)
			}
			return true
		},
		/* 30 space_comment <- <(space / comment)> */
		func() bool {
			position263, tokenIndex263, depth263 := position, tokenIndex, depth
			{
				position264 := position
				depth++
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if !_rules[rulespace]() {
						goto l266
					}
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if !_rules[rulecomment]() {
						goto l263
					}
				}
			l265:
				depth--
				add(rulespace_comment, position264)
			}
			return true
		l263:
			position, tokenIndex, depth = position263, tokenIndex263, depth263
			return false
		},
		/* 31 comment <- <('/' '/' (!end_of_line .)* end_of_line)> */
		func() bool {
			position267, tokenIndex267, depth267 := position, tokenIndex, depth
			{
				position268 := position
				depth++
				if buffer[position] != rune('/') {
					goto l267
				}
				position++
				if buffer[position] != rune('/') {
					goto l267
				}
				position++
			l269:
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					{
						position271, tokenIndex271, depth271 := position, tokenIndex, depth
						if !_rules[ruleend_of_line]() {
							goto l271
						}
						goto l270
					l271:
						position, tokenIndex, depth = position271, tokenIndex271, depth271
					}
					if !matchDot() {
						goto l270
					}
					goto l269
				l270:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
				}
				if !_rules[ruleend_of_line]() {
					goto l267
				}
				depth--
				add(rulecomment, position268)
			}
			return true
		l267:
			position, tokenIndex, depth = position267, tokenIndex267, depth267
			return false
		},
		/* 32 space <- <(' ' / '\t' / end_of_line)> */
		func() bool {
			position272, tokenIndex272, depth272 := position, tokenIndex, depth
			{
				position273 := position
				depth++
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l275
					}
					position++
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if buffer[position] != rune('\t') {
						goto l276
					}
					position++
					goto l274
				l276:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if !_rules[ruleend_of_line]() {
						goto l272
					}
				}
			l274:
				depth--
				add(rulespace, position273)
			}
			return true
		l272:
			position, tokenIndex, depth = position272, tokenIndex272, depth272
			return false
		},
		/* 33 end_of_line <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position277, tokenIndex277, depth277 := position, tokenIndex, depth
			{
				position278 := position
				depth++
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l280
					}
					position++
					if buffer[position] != rune('\n') {
						goto l280
					}
					position++
					goto l279
				l280:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
					if buffer[position] != rune('\n') {
						goto l281
					}
					position++
					goto l279
				l281:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
					if buffer[position] != rune('\r') {
						goto l277
					}
					position++
				}
			l279:
				depth--
				add(ruleend_of_line, position278)
			}
			return true
		l277:
			position, tokenIndex, depth = position277, tokenIndex277, depth277
			return false
		},
		/* 34 end_of_file <- <!.> */
		func() bool {
			position282, tokenIndex282, depth282 := position, tokenIndex, depth
			{
				position283 := position
				depth++
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if !matchDot() {
						goto l284
					}
					goto l282
				l284:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
				}
				depth--
				add(ruleend_of_file, position283)
			}
			return true
		l282:
			position, tokenIndex, depth = position282, tokenIndex282, depth282
			return false
		},
		nil,
		/* 37 Action0 <- <{p.SetNameSpace(text)}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 38 Action1 <- <{p.ExtractStruct(false)}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 39 Action2 <- <{p.ExtractStruct(true)}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 40 Action3 <- <{p.SetTypeName(text)}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 41 Action4 <- <{p.NewExtractField()}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 42 Action5 <- <{p.NewExtractFieldWithValue()}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 43 Action6 <- <{p.FieldName(text)}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 44 Action7 <- <{p.NewUnion()}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 45 Action8 <- <{p.UnionName(text)}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 46 Action9 <- <{p.NewExtractField()}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 47 Action10 <- <{p.EnumName(text)}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 48 Action11 <- <{p.SetRootType(text)}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 49 Action12 <- <{p.SetType("bool")}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 50 Action13 <- <{p.SetType("int8")}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 51 Action14 <- <{p.SetType("uint8")}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 52 Action15 <- <{p.SetType("int16")}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 53 Action16 <- <{p.SetType("uint16")}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 54 Action17 <- <{p.SetType("int32")}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 55 Action18 <- <{p.SetType("uint32")}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 56 Action19 <- <{p.SetType("int64")}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 57 Action20 <- <{p.SetType("uint64")}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 58 Action21 <- <{p.SetType("float32")}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 59 Action22 <- <{p.SetType("float64")}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 60 Action23 <- <{p.SetType("uint8")}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 61 Action24 <- <{p.SetType("uint8")}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 62 Action25 <- <{p.SetType("short")}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 63 Action26 <- <{p.SetType("ushort")}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 64 Action27 <- <{p.SetType("int32")}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 65 Action28 <- <{p.SetType("uint32")}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 66 Action29 <- <{p.SetType("float")}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 67 Action30 <- <{p.SetType("long")}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 68 Action31 <- <{p.SetType("ulong")}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 69 Action32 <- <{p.SetType("double")}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 70 Action33 <- <{p.SetRepeated("byte")}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 71 Action34 <- <{p.SetType(text)}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 72 Action35 <- <{p.SetRepeated("") }> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
	}
	p.rules = _rules
}
