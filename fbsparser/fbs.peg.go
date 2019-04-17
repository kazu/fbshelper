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
	rules  [72]func() bool
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
			p.FieldNaame(text)
		case ruleAction7:
			p.NewUnion()
		case ruleAction8:
			p.UnionName(text)
		case ruleAction9:
			p.NewExtractField()
		case ruleAction10:
			p.FieldNaame(text)
		case ruleAction11:
			p.SetType("bool")
		case ruleAction12:
			p.SetType("int8")
		case ruleAction13:
			p.SetType("uint8")
		case ruleAction14:
			p.SetType("int16")
		case ruleAction15:
			p.SetType("uint16")
		case ruleAction16:
			p.SetType("int32")
		case ruleAction17:
			p.SetType("uint32")
		case ruleAction18:
			p.SetType("int64")
		case ruleAction19:
			p.SetType("uint64")
		case ruleAction20:
			p.SetType("float32")
		case ruleAction21:
			p.SetType("float64")
		case ruleAction22:
			p.SetType("uint8")
		case ruleAction23:
			p.SetType("uint8")
		case ruleAction24:
			p.SetType("short")
		case ruleAction25:
			p.SetType("ushort")
		case ruleAction26:
			p.SetType("int")
		case ruleAction27:
			p.SetType("uint")
		case ruleAction28:
			p.SetType("float")
		case ruleAction29:
			p.SetType("long")
		case ruleAction30:
			p.SetType("ulong")
		case ruleAction31:
			p.SetType("double")
		case ruleAction32:
			p.SetRepeated("byte")
		case ruleAction33:
			p.SetType(text)
		case ruleAction34:
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
		/* 17 root_decl <- <('r' 'o' 'o' 't' '_' 't' 'y' 'p' 'e' spacing ident spacing ';' spacing)> */
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
				if !_rules[ruleident]() {
					goto l108
				}
				if !_rules[rulespacing]() {
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
			position110, tokenIndex110, depth110 := position, tokenIndex, depth
			{
				position111 := position
				depth++
				if buffer[position] != rune('f') {
					goto l110
				}
				position++
				if buffer[position] != rune('i') {
					goto l110
				}
				position++
				if buffer[position] != rune('l') {
					goto l110
				}
				position++
				if buffer[position] != rune('e') {
					goto l110
				}
				position++
				if buffer[position] != rune('_') {
					goto l110
				}
				position++
				if buffer[position] != rune('e') {
					goto l110
				}
				position++
				if buffer[position] != rune('x') {
					goto l110
				}
				position++
				if buffer[position] != rune('t') {
					goto l110
				}
				position++
				if buffer[position] != rune('e') {
					goto l110
				}
				position++
				if buffer[position] != rune('n') {
					goto l110
				}
				position++
				if buffer[position] != rune('s') {
					goto l110
				}
				position++
				if buffer[position] != rune('i') {
					goto l110
				}
				position++
				if buffer[position] != rune('o') {
					goto l110
				}
				position++
				if buffer[position] != rune('n') {
					goto l110
				}
				position++
				{
					position112 := position
					depth++
				l113:
					{
						position114, tokenIndex114, depth114 := position, tokenIndex, depth
						{
							position115, tokenIndex115, depth115 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l116
							}
							position++
							goto l115
						l116:
							position, tokenIndex, depth = position115, tokenIndex115, depth115
							if buffer[position] != rune('\t') {
								goto l114
							}
							position++
						}
					l115:
						goto l113
					l114:
						position, tokenIndex, depth = position114, tokenIndex114, depth114
					}
					depth--
					add(rulePegText, position112)
				}
				{
					position117 := position
					depth++
					{
						position120, tokenIndex120, depth120 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l120
						}
						position++
						goto l110
					l120:
						position, tokenIndex, depth = position120, tokenIndex120, depth120
					}
					if !matchDot() {
						goto l110
					}
				l118:
					{
						position119, tokenIndex119, depth119 := position, tokenIndex, depth
						{
							position121, tokenIndex121, depth121 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l121
							}
							position++
							goto l119
						l121:
							position, tokenIndex, depth = position121, tokenIndex121, depth121
						}
						if !matchDot() {
							goto l119
						}
						goto l118
					l119:
						position, tokenIndex, depth = position119, tokenIndex119, depth119
					}
					depth--
					add(rulePegText, position117)
				}
				if buffer[position] != rune(';') {
					goto l110
				}
				position++
				if !_rules[rulespacing]() {
					goto l110
				}
				depth--
				add(rulefile_extension_decl, position111)
			}
			return true
		l110:
			position, tokenIndex, depth = position110, tokenIndex110, depth110
			return false
		},
		/* 19 file_identifier_decl <- <('f' 'i' 'l' 'e' '_' 'i' 'd' 'e' 'n' 't' 'i' 'f' 'i' 'e' 'r' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position122, tokenIndex122, depth122 := position, tokenIndex, depth
			{
				position123 := position
				depth++
				if buffer[position] != rune('f') {
					goto l122
				}
				position++
				if buffer[position] != rune('i') {
					goto l122
				}
				position++
				if buffer[position] != rune('l') {
					goto l122
				}
				position++
				if buffer[position] != rune('e') {
					goto l122
				}
				position++
				if buffer[position] != rune('_') {
					goto l122
				}
				position++
				if buffer[position] != rune('i') {
					goto l122
				}
				position++
				if buffer[position] != rune('d') {
					goto l122
				}
				position++
				if buffer[position] != rune('e') {
					goto l122
				}
				position++
				if buffer[position] != rune('n') {
					goto l122
				}
				position++
				if buffer[position] != rune('t') {
					goto l122
				}
				position++
				if buffer[position] != rune('i') {
					goto l122
				}
				position++
				if buffer[position] != rune('f') {
					goto l122
				}
				position++
				if buffer[position] != rune('i') {
					goto l122
				}
				position++
				if buffer[position] != rune('e') {
					goto l122
				}
				position++
				if buffer[position] != rune('r') {
					goto l122
				}
				position++
				{
					position124 := position
					depth++
				l125:
					{
						position126, tokenIndex126, depth126 := position, tokenIndex, depth
						{
							position127, tokenIndex127, depth127 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l128
							}
							position++
							goto l127
						l128:
							position, tokenIndex, depth = position127, tokenIndex127, depth127
							if buffer[position] != rune('\t') {
								goto l126
							}
							position++
						}
					l127:
						goto l125
					l126:
						position, tokenIndex, depth = position126, tokenIndex126, depth126
					}
					depth--
					add(rulePegText, position124)
				}
				{
					position129 := position
					depth++
					{
						position132, tokenIndex132, depth132 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l132
						}
						position++
						goto l122
					l132:
						position, tokenIndex, depth = position132, tokenIndex132, depth132
					}
					if !matchDot() {
						goto l122
					}
				l130:
					{
						position131, tokenIndex131, depth131 := position, tokenIndex, depth
						{
							position133, tokenIndex133, depth133 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l133
							}
							position++
							goto l131
						l133:
							position, tokenIndex, depth = position133, tokenIndex133, depth133
						}
						if !matchDot() {
							goto l131
						}
						goto l130
					l131:
						position, tokenIndex, depth = position131, tokenIndex131, depth131
					}
					depth--
					add(rulePegText, position129)
				}
				if buffer[position] != rune(';') {
					goto l122
				}
				position++
				if !_rules[rulespacing]() {
					goto l122
				}
				depth--
				add(rulefile_identifier_decl, position123)
			}
			return true
		l122:
			position, tokenIndex, depth = position122, tokenIndex122, depth122
			return false
		},
		/* 20 attribute_decl <- <('a' 't' 't' 'r' 'i' 'b' 'u' 't' 'e' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position134, tokenIndex134, depth134 := position, tokenIndex, depth
			{
				position135 := position
				depth++
				if buffer[position] != rune('a') {
					goto l134
				}
				position++
				if buffer[position] != rune('t') {
					goto l134
				}
				position++
				if buffer[position] != rune('t') {
					goto l134
				}
				position++
				if buffer[position] != rune('r') {
					goto l134
				}
				position++
				if buffer[position] != rune('i') {
					goto l134
				}
				position++
				if buffer[position] != rune('b') {
					goto l134
				}
				position++
				if buffer[position] != rune('u') {
					goto l134
				}
				position++
				if buffer[position] != rune('t') {
					goto l134
				}
				position++
				if buffer[position] != rune('e') {
					goto l134
				}
				position++
				{
					position136 := position
					depth++
				l137:
					{
						position138, tokenIndex138, depth138 := position, tokenIndex, depth
						{
							position139, tokenIndex139, depth139 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l140
							}
							position++
							goto l139
						l140:
							position, tokenIndex, depth = position139, tokenIndex139, depth139
							if buffer[position] != rune('\t') {
								goto l138
							}
							position++
						}
					l139:
						goto l137
					l138:
						position, tokenIndex, depth = position138, tokenIndex138, depth138
					}
					depth--
					add(rulePegText, position136)
				}
				{
					position141 := position
					depth++
					{
						position144, tokenIndex144, depth144 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l144
						}
						position++
						goto l134
					l144:
						position, tokenIndex, depth = position144, tokenIndex144, depth144
					}
					if !matchDot() {
						goto l134
					}
				l142:
					{
						position143, tokenIndex143, depth143 := position, tokenIndex, depth
						{
							position145, tokenIndex145, depth145 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l145
							}
							position++
							goto l143
						l145:
							position, tokenIndex, depth = position145, tokenIndex145, depth145
						}
						if !matchDot() {
							goto l143
						}
						goto l142
					l143:
						position, tokenIndex, depth = position143, tokenIndex143, depth143
					}
					depth--
					add(rulePegText, position141)
				}
				if buffer[position] != rune(';') {
					goto l134
				}
				position++
				if !_rules[rulespacing]() {
					goto l134
				}
				depth--
				add(ruleattribute_decl, position135)
			}
			return true
		l134:
			position, tokenIndex, depth = position134, tokenIndex134, depth134
			return false
		},
		/* 21 rpc_decl <- <('r' 'p' 'c' '_' 's' 'e' 'r' 'v' 'i' 'c' 'e' <(' ' / '\t')*> ident '{' <(!'}' .)+> '}' spacing)> */
		func() bool {
			position146, tokenIndex146, depth146 := position, tokenIndex, depth
			{
				position147 := position
				depth++
				if buffer[position] != rune('r') {
					goto l146
				}
				position++
				if buffer[position] != rune('p') {
					goto l146
				}
				position++
				if buffer[position] != rune('c') {
					goto l146
				}
				position++
				if buffer[position] != rune('_') {
					goto l146
				}
				position++
				if buffer[position] != rune('s') {
					goto l146
				}
				position++
				if buffer[position] != rune('e') {
					goto l146
				}
				position++
				if buffer[position] != rune('r') {
					goto l146
				}
				position++
				if buffer[position] != rune('v') {
					goto l146
				}
				position++
				if buffer[position] != rune('i') {
					goto l146
				}
				position++
				if buffer[position] != rune('c') {
					goto l146
				}
				position++
				if buffer[position] != rune('e') {
					goto l146
				}
				position++
				{
					position148 := position
					depth++
				l149:
					{
						position150, tokenIndex150, depth150 := position, tokenIndex, depth
						{
							position151, tokenIndex151, depth151 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l152
							}
							position++
							goto l151
						l152:
							position, tokenIndex, depth = position151, tokenIndex151, depth151
							if buffer[position] != rune('\t') {
								goto l150
							}
							position++
						}
					l151:
						goto l149
					l150:
						position, tokenIndex, depth = position150, tokenIndex150, depth150
					}
					depth--
					add(rulePegText, position148)
				}
				if !_rules[ruleident]() {
					goto l146
				}
				if buffer[position] != rune('{') {
					goto l146
				}
				position++
				{
					position153 := position
					depth++
					{
						position156, tokenIndex156, depth156 := position, tokenIndex, depth
						if buffer[position] != rune('}') {
							goto l156
						}
						position++
						goto l146
					l156:
						position, tokenIndex, depth = position156, tokenIndex156, depth156
					}
					if !matchDot() {
						goto l146
					}
				l154:
					{
						position155, tokenIndex155, depth155 := position, tokenIndex, depth
						{
							position157, tokenIndex157, depth157 := position, tokenIndex, depth
							if buffer[position] != rune('}') {
								goto l157
							}
							position++
							goto l155
						l157:
							position, tokenIndex, depth = position157, tokenIndex157, depth157
						}
						if !matchDot() {
							goto l155
						}
						goto l154
					l155:
						position, tokenIndex, depth = position155, tokenIndex155, depth155
					}
					depth--
					add(rulePegText, position153)
				}
				if buffer[position] != rune('}') {
					goto l146
				}
				position++
				if !_rules[rulespacing]() {
					goto l146
				}
				depth--
				add(rulerpc_decl, position147)
			}
			return true
		l146:
			position, tokenIndex, depth = position146, tokenIndex146, depth146
			return false
		},
		/* 22 type <- <(('b' 'o' 'o' 'l' spacing Action11) / ('i' 'n' 't' '8' spacing Action12) / ('u' 'i' 'n' 't' '8' spacing Action13) / ('i' 'n' 't' '1' '6' spacing Action14) / ('u' 'i' 'n' 't' '1' '6' spacing Action15) / ('i' 'n' 't' '3' '2' spacing Action16) / ('u' 'i' 'n' 't' '3' '2' spacing Action17) / ('i' 'n' 't' '6' '4' spacing Action18) / ('u' 'i' 'n' 't' '6' '4' spacing Action19) / ('f' 'l' 'o' 'a' 't' '3' '2' spacing Action20) / ('f' 'l' 'o' 'a' 't' '6' '4' spacing Action21) / ('b' 'y' 't' 'e' spacing Action22) / ('u' 'b' 'y' 't' 'e' spacing Action23) / ('s' 'h' 'o' 'r' 't' spacing Action24) / ('u' 's' 'h' 'o' 'r' 't' spacing Action25) / ('i' 'n' 't' spacing Action26) / ('u' 'i' 'n' 't' spacing Action27) / ('f' 'l' 'o' 'a' 't' spacing Action28) / ('l' 'o' 'n' 'g' spacing Action29) / ('u' 'l' 'o' 'n' 'g' spacing Action30) / ('d' 'o' 'u' 'b' 'l' 'e' spacing Action31) / ('s' 't' 'r' 'i' 'n' 'g' spacing Action32) / (ident spacing Action33) / ('[' type ']' spacing Action34))> */
		func() bool {
			position158, tokenIndex158, depth158 := position, tokenIndex, depth
			{
				position159 := position
				depth++
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l161
					}
					position++
					if buffer[position] != rune('o') {
						goto l161
					}
					position++
					if buffer[position] != rune('o') {
						goto l161
					}
					position++
					if buffer[position] != rune('l') {
						goto l161
					}
					position++
					if !_rules[rulespacing]() {
						goto l161
					}
					if !_rules[ruleAction11]() {
						goto l161
					}
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('i') {
						goto l162
					}
					position++
					if buffer[position] != rune('n') {
						goto l162
					}
					position++
					if buffer[position] != rune('t') {
						goto l162
					}
					position++
					if buffer[position] != rune('8') {
						goto l162
					}
					position++
					if !_rules[rulespacing]() {
						goto l162
					}
					if !_rules[ruleAction12]() {
						goto l162
					}
					goto l160
				l162:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('u') {
						goto l163
					}
					position++
					if buffer[position] != rune('i') {
						goto l163
					}
					position++
					if buffer[position] != rune('n') {
						goto l163
					}
					position++
					if buffer[position] != rune('t') {
						goto l163
					}
					position++
					if buffer[position] != rune('8') {
						goto l163
					}
					position++
					if !_rules[rulespacing]() {
						goto l163
					}
					if !_rules[ruleAction13]() {
						goto l163
					}
					goto l160
				l163:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('i') {
						goto l164
					}
					position++
					if buffer[position] != rune('n') {
						goto l164
					}
					position++
					if buffer[position] != rune('t') {
						goto l164
					}
					position++
					if buffer[position] != rune('1') {
						goto l164
					}
					position++
					if buffer[position] != rune('6') {
						goto l164
					}
					position++
					if !_rules[rulespacing]() {
						goto l164
					}
					if !_rules[ruleAction14]() {
						goto l164
					}
					goto l160
				l164:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('u') {
						goto l165
					}
					position++
					if buffer[position] != rune('i') {
						goto l165
					}
					position++
					if buffer[position] != rune('n') {
						goto l165
					}
					position++
					if buffer[position] != rune('t') {
						goto l165
					}
					position++
					if buffer[position] != rune('1') {
						goto l165
					}
					position++
					if buffer[position] != rune('6') {
						goto l165
					}
					position++
					if !_rules[rulespacing]() {
						goto l165
					}
					if !_rules[ruleAction15]() {
						goto l165
					}
					goto l160
				l165:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('i') {
						goto l166
					}
					position++
					if buffer[position] != rune('n') {
						goto l166
					}
					position++
					if buffer[position] != rune('t') {
						goto l166
					}
					position++
					if buffer[position] != rune('3') {
						goto l166
					}
					position++
					if buffer[position] != rune('2') {
						goto l166
					}
					position++
					if !_rules[rulespacing]() {
						goto l166
					}
					if !_rules[ruleAction16]() {
						goto l166
					}
					goto l160
				l166:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('u') {
						goto l167
					}
					position++
					if buffer[position] != rune('i') {
						goto l167
					}
					position++
					if buffer[position] != rune('n') {
						goto l167
					}
					position++
					if buffer[position] != rune('t') {
						goto l167
					}
					position++
					if buffer[position] != rune('3') {
						goto l167
					}
					position++
					if buffer[position] != rune('2') {
						goto l167
					}
					position++
					if !_rules[rulespacing]() {
						goto l167
					}
					if !_rules[ruleAction17]() {
						goto l167
					}
					goto l160
				l167:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('i') {
						goto l168
					}
					position++
					if buffer[position] != rune('n') {
						goto l168
					}
					position++
					if buffer[position] != rune('t') {
						goto l168
					}
					position++
					if buffer[position] != rune('6') {
						goto l168
					}
					position++
					if buffer[position] != rune('4') {
						goto l168
					}
					position++
					if !_rules[rulespacing]() {
						goto l168
					}
					if !_rules[ruleAction18]() {
						goto l168
					}
					goto l160
				l168:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('u') {
						goto l169
					}
					position++
					if buffer[position] != rune('i') {
						goto l169
					}
					position++
					if buffer[position] != rune('n') {
						goto l169
					}
					position++
					if buffer[position] != rune('t') {
						goto l169
					}
					position++
					if buffer[position] != rune('6') {
						goto l169
					}
					position++
					if buffer[position] != rune('4') {
						goto l169
					}
					position++
					if !_rules[rulespacing]() {
						goto l169
					}
					if !_rules[ruleAction19]() {
						goto l169
					}
					goto l160
				l169:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('f') {
						goto l170
					}
					position++
					if buffer[position] != rune('l') {
						goto l170
					}
					position++
					if buffer[position] != rune('o') {
						goto l170
					}
					position++
					if buffer[position] != rune('a') {
						goto l170
					}
					position++
					if buffer[position] != rune('t') {
						goto l170
					}
					position++
					if buffer[position] != rune('3') {
						goto l170
					}
					position++
					if buffer[position] != rune('2') {
						goto l170
					}
					position++
					if !_rules[rulespacing]() {
						goto l170
					}
					if !_rules[ruleAction20]() {
						goto l170
					}
					goto l160
				l170:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('f') {
						goto l171
					}
					position++
					if buffer[position] != rune('l') {
						goto l171
					}
					position++
					if buffer[position] != rune('o') {
						goto l171
					}
					position++
					if buffer[position] != rune('a') {
						goto l171
					}
					position++
					if buffer[position] != rune('t') {
						goto l171
					}
					position++
					if buffer[position] != rune('6') {
						goto l171
					}
					position++
					if buffer[position] != rune('4') {
						goto l171
					}
					position++
					if !_rules[rulespacing]() {
						goto l171
					}
					if !_rules[ruleAction21]() {
						goto l171
					}
					goto l160
				l171:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('b') {
						goto l172
					}
					position++
					if buffer[position] != rune('y') {
						goto l172
					}
					position++
					if buffer[position] != rune('t') {
						goto l172
					}
					position++
					if buffer[position] != rune('e') {
						goto l172
					}
					position++
					if !_rules[rulespacing]() {
						goto l172
					}
					if !_rules[ruleAction22]() {
						goto l172
					}
					goto l160
				l172:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('u') {
						goto l173
					}
					position++
					if buffer[position] != rune('b') {
						goto l173
					}
					position++
					if buffer[position] != rune('y') {
						goto l173
					}
					position++
					if buffer[position] != rune('t') {
						goto l173
					}
					position++
					if buffer[position] != rune('e') {
						goto l173
					}
					position++
					if !_rules[rulespacing]() {
						goto l173
					}
					if !_rules[ruleAction23]() {
						goto l173
					}
					goto l160
				l173:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('s') {
						goto l174
					}
					position++
					if buffer[position] != rune('h') {
						goto l174
					}
					position++
					if buffer[position] != rune('o') {
						goto l174
					}
					position++
					if buffer[position] != rune('r') {
						goto l174
					}
					position++
					if buffer[position] != rune('t') {
						goto l174
					}
					position++
					if !_rules[rulespacing]() {
						goto l174
					}
					if !_rules[ruleAction24]() {
						goto l174
					}
					goto l160
				l174:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('u') {
						goto l175
					}
					position++
					if buffer[position] != rune('s') {
						goto l175
					}
					position++
					if buffer[position] != rune('h') {
						goto l175
					}
					position++
					if buffer[position] != rune('o') {
						goto l175
					}
					position++
					if buffer[position] != rune('r') {
						goto l175
					}
					position++
					if buffer[position] != rune('t') {
						goto l175
					}
					position++
					if !_rules[rulespacing]() {
						goto l175
					}
					if !_rules[ruleAction25]() {
						goto l175
					}
					goto l160
				l175:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
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
					if !_rules[rulespacing]() {
						goto l176
					}
					if !_rules[ruleAction26]() {
						goto l176
					}
					goto l160
				l176:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('u') {
						goto l177
					}
					position++
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
					if !_rules[rulespacing]() {
						goto l177
					}
					if !_rules[ruleAction27]() {
						goto l177
					}
					goto l160
				l177:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('f') {
						goto l178
					}
					position++
					if buffer[position] != rune('l') {
						goto l178
					}
					position++
					if buffer[position] != rune('o') {
						goto l178
					}
					position++
					if buffer[position] != rune('a') {
						goto l178
					}
					position++
					if buffer[position] != rune('t') {
						goto l178
					}
					position++
					if !_rules[rulespacing]() {
						goto l178
					}
					if !_rules[ruleAction28]() {
						goto l178
					}
					goto l160
				l178:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('l') {
						goto l179
					}
					position++
					if buffer[position] != rune('o') {
						goto l179
					}
					position++
					if buffer[position] != rune('n') {
						goto l179
					}
					position++
					if buffer[position] != rune('g') {
						goto l179
					}
					position++
					if !_rules[rulespacing]() {
						goto l179
					}
					if !_rules[ruleAction29]() {
						goto l179
					}
					goto l160
				l179:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('u') {
						goto l180
					}
					position++
					if buffer[position] != rune('l') {
						goto l180
					}
					position++
					if buffer[position] != rune('o') {
						goto l180
					}
					position++
					if buffer[position] != rune('n') {
						goto l180
					}
					position++
					if buffer[position] != rune('g') {
						goto l180
					}
					position++
					if !_rules[rulespacing]() {
						goto l180
					}
					if !_rules[ruleAction30]() {
						goto l180
					}
					goto l160
				l180:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('d') {
						goto l181
					}
					position++
					if buffer[position] != rune('o') {
						goto l181
					}
					position++
					if buffer[position] != rune('u') {
						goto l181
					}
					position++
					if buffer[position] != rune('b') {
						goto l181
					}
					position++
					if buffer[position] != rune('l') {
						goto l181
					}
					position++
					if buffer[position] != rune('e') {
						goto l181
					}
					position++
					if !_rules[rulespacing]() {
						goto l181
					}
					if !_rules[ruleAction31]() {
						goto l181
					}
					goto l160
				l181:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('s') {
						goto l182
					}
					position++
					if buffer[position] != rune('t') {
						goto l182
					}
					position++
					if buffer[position] != rune('r') {
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
					if buffer[position] != rune('g') {
						goto l182
					}
					position++
					if !_rules[rulespacing]() {
						goto l182
					}
					if !_rules[ruleAction32]() {
						goto l182
					}
					goto l160
				l182:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if !_rules[ruleident]() {
						goto l183
					}
					if !_rules[rulespacing]() {
						goto l183
					}
					if !_rules[ruleAction33]() {
						goto l183
					}
					goto l160
				l183:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('[') {
						goto l158
					}
					position++
					if !_rules[ruletype]() {
						goto l158
					}
					if buffer[position] != rune(']') {
						goto l158
					}
					position++
					if !_rules[rulespacing]() {
						goto l158
					}
					if !_rules[ruleAction34]() {
						goto l158
					}
				}
			l160:
				depth--
				add(ruletype, position159)
			}
			return true
		l158:
			position, tokenIndex, depth = position158, tokenIndex158, depth158
			return false
		},
		/* 23 scalar <- <(integer_constant / float_constant)> */
		func() bool {
			position184, tokenIndex184, depth184 := position, tokenIndex, depth
			{
				position185 := position
				depth++
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if !_rules[ruleinteger_constant]() {
						goto l187
					}
					goto l186
				l187:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
					if !_rules[rulefloat_constant]() {
						goto l184
					}
				}
			l186:
				depth--
				add(rulescalar, position185)
			}
			return true
		l184:
			position, tokenIndex, depth = position184, tokenIndex184, depth184
			return false
		},
		/* 24 integer_constant <- <(<[0-9]+> / ('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		func() bool {
			position188, tokenIndex188, depth188 := position, tokenIndex, depth
			{
				position189 := position
				depth++
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					{
						position192 := position
						depth++
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l191
						}
						position++
					l193:
						{
							position194, tokenIndex194, depth194 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l194
							}
							position++
							goto l193
						l194:
							position, tokenIndex, depth = position194, tokenIndex194, depth194
						}
						depth--
						add(rulePegText, position192)
					}
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('t') {
						goto l195
					}
					position++
					if buffer[position] != rune('r') {
						goto l195
					}
					position++
					if buffer[position] != rune('u') {
						goto l195
					}
					position++
					if buffer[position] != rune('e') {
						goto l195
					}
					position++
					goto l190
				l195:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('f') {
						goto l188
					}
					position++
					if buffer[position] != rune('a') {
						goto l188
					}
					position++
					if buffer[position] != rune('l') {
						goto l188
					}
					position++
					if buffer[position] != rune('s') {
						goto l188
					}
					position++
					if buffer[position] != rune('e') {
						goto l188
					}
					position++
				}
			l190:
				depth--
				add(ruleinteger_constant, position189)
			}
			return true
		l188:
			position, tokenIndex, depth = position188, tokenIndex188, depth188
			return false
		},
		/* 25 float_constant <- <(<('-'* [0-9]+ . [0-9])> / float_constant_exp)> */
		func() bool {
			position196, tokenIndex196, depth196 := position, tokenIndex, depth
			{
				position197 := position
				depth++
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					{
						position200 := position
						depth++
					l201:
						{
							position202, tokenIndex202, depth202 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l202
							}
							position++
							goto l201
						l202:
							position, tokenIndex, depth = position202, tokenIndex202, depth202
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l199
						}
						position++
					l203:
						{
							position204, tokenIndex204, depth204 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l204
							}
							position++
							goto l203
						l204:
							position, tokenIndex, depth = position204, tokenIndex204, depth204
						}
						if !matchDot() {
							goto l199
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l199
						}
						position++
						depth--
						add(rulePegText, position200)
					}
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if !_rules[rulefloat_constant_exp]() {
						goto l196
					}
				}
			l198:
				depth--
				add(rulefloat_constant, position197)
			}
			return true
		l196:
			position, tokenIndex, depth = position196, tokenIndex196, depth196
			return false
		},
		/* 26 float_constant_exp <- <(<('-'* [0-9]+ . [0-9]+)> <('e' / 'E')> <([+-]] / '>' / ' ' / '<' / '[' / [0-9])+>)> */
		func() bool {
			position205, tokenIndex205, depth205 := position, tokenIndex, depth
			{
				position206 := position
				depth++
				{
					position207 := position
					depth++
				l208:
					{
						position209, tokenIndex209, depth209 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l209
						}
						position++
						goto l208
					l209:
						position, tokenIndex, depth = position209, tokenIndex209, depth209
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l205
					}
					position++
				l210:
					{
						position211, tokenIndex211, depth211 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l211
						}
						position++
						goto l210
					l211:
						position, tokenIndex, depth = position211, tokenIndex211, depth211
					}
					if !matchDot() {
						goto l205
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l205
					}
					position++
				l212:
					{
						position213, tokenIndex213, depth213 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l213
						}
						position++
						goto l212
					l213:
						position, tokenIndex, depth = position213, tokenIndex213, depth213
					}
					depth--
					add(rulePegText, position207)
				}
				{
					position214 := position
					depth++
					{
						position215, tokenIndex215, depth215 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l216
						}
						position++
						goto l215
					l216:
						position, tokenIndex, depth = position215, tokenIndex215, depth215
						if buffer[position] != rune('E') {
							goto l205
						}
						position++
					}
				l215:
					depth--
					add(rulePegText, position214)
				}
				{
					position217 := position
					depth++
					{
						position220, tokenIndex220, depth220 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('+') || c > rune(']') {
							goto l221
						}
						position++
						goto l220
					l221:
						position, tokenIndex, depth = position220, tokenIndex220, depth220
						if buffer[position] != rune('>') {
							goto l222
						}
						position++
						goto l220
					l222:
						position, tokenIndex, depth = position220, tokenIndex220, depth220
						if buffer[position] != rune(' ') {
							goto l223
						}
						position++
						goto l220
					l223:
						position, tokenIndex, depth = position220, tokenIndex220, depth220
						if buffer[position] != rune('<') {
							goto l224
						}
						position++
						goto l220
					l224:
						position, tokenIndex, depth = position220, tokenIndex220, depth220
						if buffer[position] != rune('[') {
							goto l225
						}
						position++
						goto l220
					l225:
						position, tokenIndex, depth = position220, tokenIndex220, depth220
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l205
						}
						position++
					}
				l220:
				l218:
					{
						position219, tokenIndex219, depth219 := position, tokenIndex, depth
						{
							position226, tokenIndex226, depth226 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('+') || c > rune(']') {
								goto l227
							}
							position++
							goto l226
						l227:
							position, tokenIndex, depth = position226, tokenIndex226, depth226
							if buffer[position] != rune('>') {
								goto l228
							}
							position++
							goto l226
						l228:
							position, tokenIndex, depth = position226, tokenIndex226, depth226
							if buffer[position] != rune(' ') {
								goto l229
							}
							position++
							goto l226
						l229:
							position, tokenIndex, depth = position226, tokenIndex226, depth226
							if buffer[position] != rune('<') {
								goto l230
							}
							position++
							goto l226
						l230:
							position, tokenIndex, depth = position226, tokenIndex226, depth226
							if buffer[position] != rune('[') {
								goto l231
							}
							position++
							goto l226
						l231:
							position, tokenIndex, depth = position226, tokenIndex226, depth226
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l219
							}
							position++
						}
					l226:
						goto l218
					l219:
						position, tokenIndex, depth = position219, tokenIndex219, depth219
					}
					depth--
					add(rulePegText, position217)
				}
				depth--
				add(rulefloat_constant_exp, position206)
			}
			return true
		l205:
			position, tokenIndex, depth = position205, tokenIndex205, depth205
			return false
		},
		/* 27 ident <- <<(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position232, tokenIndex232, depth232 := position, tokenIndex, depth
			{
				position233 := position
				depth++
				{
					position234 := position
					depth++
					{
						position235, tokenIndex235, depth235 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l236
						}
						position++
						goto l235
					l236:
						position, tokenIndex, depth = position235, tokenIndex235, depth235
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l237
						}
						position++
						goto l235
					l237:
						position, tokenIndex, depth = position235, tokenIndex235, depth235
						if buffer[position] != rune('_') {
							goto l232
						}
						position++
					}
				l235:
				l238:
					{
						position239, tokenIndex239, depth239 := position, tokenIndex, depth
						{
							position240, tokenIndex240, depth240 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l241
							}
							position++
							goto l240
						l241:
							position, tokenIndex, depth = position240, tokenIndex240, depth240
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l242
							}
							position++
							goto l240
						l242:
							position, tokenIndex, depth = position240, tokenIndex240, depth240
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l243
							}
							position++
							goto l240
						l243:
							position, tokenIndex, depth = position240, tokenIndex240, depth240
							if buffer[position] != rune('_') {
								goto l239
							}
							position++
						}
					l240:
						goto l238
					l239:
						position, tokenIndex, depth = position239, tokenIndex239, depth239
					}
					depth--
					add(rulePegText, position234)
				}
				depth--
				add(ruleident, position233)
			}
			return true
		l232:
			position, tokenIndex, depth = position232, tokenIndex232, depth232
			return false
		},
		/* 28 only_comment <- <(spacing ';')> */
		func() bool {
			position244, tokenIndex244, depth244 := position, tokenIndex, depth
			{
				position245 := position
				depth++
				if !_rules[rulespacing]() {
					goto l244
				}
				if buffer[position] != rune(';') {
					goto l244
				}
				position++
				depth--
				add(ruleonly_comment, position245)
			}
			return true
		l244:
			position, tokenIndex, depth = position244, tokenIndex244, depth244
			return false
		},
		/* 29 spacing <- <space_comment*> */
		func() bool {
			{
				position247 := position
				depth++
			l248:
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if !_rules[rulespace_comment]() {
						goto l249
					}
					goto l248
				l249:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
				}
				depth--
				add(rulespacing, position247)
			}
			return true
		},
		/* 30 space_comment <- <(space / comment)> */
		func() bool {
			position250, tokenIndex250, depth250 := position, tokenIndex, depth
			{
				position251 := position
				depth++
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if !_rules[rulespace]() {
						goto l253
					}
					goto l252
				l253:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if !_rules[rulecomment]() {
						goto l250
					}
				}
			l252:
				depth--
				add(rulespace_comment, position251)
			}
			return true
		l250:
			position, tokenIndex, depth = position250, tokenIndex250, depth250
			return false
		},
		/* 31 comment <- <('/' '/' (!end_of_line .)* end_of_line)> */
		func() bool {
			position254, tokenIndex254, depth254 := position, tokenIndex, depth
			{
				position255 := position
				depth++
				if buffer[position] != rune('/') {
					goto l254
				}
				position++
				if buffer[position] != rune('/') {
					goto l254
				}
				position++
			l256:
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					{
						position258, tokenIndex258, depth258 := position, tokenIndex, depth
						if !_rules[ruleend_of_line]() {
							goto l258
						}
						goto l257
					l258:
						position, tokenIndex, depth = position258, tokenIndex258, depth258
					}
					if !matchDot() {
						goto l257
					}
					goto l256
				l257:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
				}
				if !_rules[ruleend_of_line]() {
					goto l254
				}
				depth--
				add(rulecomment, position255)
			}
			return true
		l254:
			position, tokenIndex, depth = position254, tokenIndex254, depth254
			return false
		},
		/* 32 space <- <(' ' / '\t' / end_of_line)> */
		func() bool {
			position259, tokenIndex259, depth259 := position, tokenIndex, depth
			{
				position260 := position
				depth++
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l262
					}
					position++
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if buffer[position] != rune('\t') {
						goto l263
					}
					position++
					goto l261
				l263:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if !_rules[ruleend_of_line]() {
						goto l259
					}
				}
			l261:
				depth--
				add(rulespace, position260)
			}
			return true
		l259:
			position, tokenIndex, depth = position259, tokenIndex259, depth259
			return false
		},
		/* 33 end_of_line <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position264, tokenIndex264, depth264 := position, tokenIndex, depth
			{
				position265 := position
				depth++
				{
					position266, tokenIndex266, depth266 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l267
					}
					position++
					if buffer[position] != rune('\n') {
						goto l267
					}
					position++
					goto l266
				l267:
					position, tokenIndex, depth = position266, tokenIndex266, depth266
					if buffer[position] != rune('\n') {
						goto l268
					}
					position++
					goto l266
				l268:
					position, tokenIndex, depth = position266, tokenIndex266, depth266
					if buffer[position] != rune('\r') {
						goto l264
					}
					position++
				}
			l266:
				depth--
				add(ruleend_of_line, position265)
			}
			return true
		l264:
			position, tokenIndex, depth = position264, tokenIndex264, depth264
			return false
		},
		/* 34 end_of_file <- <!.> */
		func() bool {
			position269, tokenIndex269, depth269 := position, tokenIndex, depth
			{
				position270 := position
				depth++
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if !matchDot() {
						goto l271
					}
					goto l269
				l271:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
				}
				depth--
				add(ruleend_of_file, position270)
			}
			return true
		l269:
			position, tokenIndex, depth = position269, tokenIndex269, depth269
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
		/* 43 Action6 <- <{p.FieldNaame(text)}> */
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
		/* 47 Action10 <- <{p.FieldNaame(text)}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 48 Action11 <- <{p.SetType("bool")}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 49 Action12 <- <{p.SetType("int8")}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 50 Action13 <- <{p.SetType("uint8")}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 51 Action14 <- <{p.SetType("int16")}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 52 Action15 <- <{p.SetType("uint16")}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 53 Action16 <- <{p.SetType("int32")}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 54 Action17 <- <{p.SetType("uint32")}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 55 Action18 <- <{p.SetType("int64")}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 56 Action19 <- <{p.SetType("uint64")}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 57 Action20 <- <{p.SetType("float32")}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 58 Action21 <- <{p.SetType("float64")}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 59 Action22 <- <{p.SetType("uint8")}> */
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
		/* 61 Action24 <- <{p.SetType("short")}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 62 Action25 <- <{p.SetType("ushort")}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 63 Action26 <- <{p.SetType("int")}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 64 Action27 <- <{p.SetType("uint")}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 65 Action28 <- <{p.SetType("float")}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 66 Action29 <- <{p.SetType("long")}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 67 Action30 <- <{p.SetType("ulong")}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 68 Action31 <- <{p.SetType("double")}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 69 Action32 <- <{p.SetRepeated("byte")}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 70 Action33 <- <{p.SetType(text)}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 71 Action34 <- <{p.SetRepeated("") }> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
	}
	p.rules = _rules
}
