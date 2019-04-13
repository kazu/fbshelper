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
	ruletype_label
	ruletypename
	rulemetadata
	rulefield_decl
	rulefield_type
	ruleenum_decl
	ruleenum_fields
	ruleunion_decl
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
	"type_label",
	"typename",
	"metadata",
	"field_decl",
	"field_type",
	"enum_decl",
	"enum_fields",
	"union_decl",
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
	rules  [67]func() bool
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
			p.ExtractStruct()
		case ruleAction2:
			p.SetTypeName(text)
		case ruleAction3:
			p.NewExtractField()
		case ruleAction4:
			p.NewExtractFieldWithValue()
		case ruleAction5:
			p.FieldNaame(text)
		case ruleAction6:
			p.NewUnion(text)
		case ruleAction7:
			p.NewExtractField()
		case ruleAction8:
			p.FieldNaame(text)
		case ruleAction9:
			p.SetType("bool")
		case ruleAction10:
			p.SetType("byte")
		case ruleAction11:
			p.SetType("ubyte")
		case ruleAction12:
			p.SetType("short")
		case ruleAction13:
			p.SetType("ushort")
		case ruleAction14:
			p.SetType("int")
		case ruleAction15:
			p.SetType("uint")
		case ruleAction16:
			p.SetType("float")
		case ruleAction17:
			p.SetType("long")
		case ruleAction18:
			p.SetType("ulong")
		case ruleAction19:
			p.SetType("double")
		case ruleAction20:
			p.SetType("int8")
		case ruleAction21:
			p.SetType("int16")
		case ruleAction22:
			p.SetType("uint16")
		case ruleAction23:
			p.SetType("int32")
		case ruleAction24:
			p.SetType("uint32")
		case ruleAction25:
			p.SetType("int64")
		case ruleAction26:
			p.SetType("uint64")
		case ruleAction27:
			p.SetType("float32")
		case ruleAction28:
			p.SetType("float64")
		case ruleAction29:
			p.SetType("string")
		case ruleAction30:
			p.SetType(text)
		case ruleAction31:
			p.SetRepeated(text)

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
		/* 4 type_decl <- <(type_label spacing typename spacing metadata* '{' field_decl+ '}' spacing Action1)> */
		func() bool {
			position38, tokenIndex38, depth38 := position, tokenIndex, depth
			{
				position39 := position
				depth++
				if !_rules[ruletype_label]() {
					goto l38
				}
				if !_rules[rulespacing]() {
					goto l38
				}
				if !_rules[ruletypename]() {
					goto l38
				}
				if !_rules[rulespacing]() {
					goto l38
				}
			l40:
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if !_rules[rulemetadata]() {
						goto l41
					}
					goto l40
				l41:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
				}
				if buffer[position] != rune('{') {
					goto l38
				}
				position++
				if !_rules[rulefield_decl]() {
					goto l38
				}
			l42:
				{
					position43, tokenIndex43, depth43 := position, tokenIndex, depth
					if !_rules[rulefield_decl]() {
						goto l43
					}
					goto l42
				l43:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
				}
				if buffer[position] != rune('}') {
					goto l38
				}
				position++
				if !_rules[rulespacing]() {
					goto l38
				}
				if !_rules[ruleAction1]() {
					goto l38
				}
				depth--
				add(ruletype_decl, position39)
			}
			return true
		l38:
			position, tokenIndex, depth = position38, tokenIndex38, depth38
			return false
		},
		/* 5 type_label <- <(('t' 'a' 'b' 'l' 'e') / ('s' 't' 'r' 'u' 'c' 't'))> */
		func() bool {
			position44, tokenIndex44, depth44 := position, tokenIndex, depth
			{
				position45 := position
				depth++
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l47
					}
					position++
					if buffer[position] != rune('a') {
						goto l47
					}
					position++
					if buffer[position] != rune('b') {
						goto l47
					}
					position++
					if buffer[position] != rune('l') {
						goto l47
					}
					position++
					if buffer[position] != rune('e') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
					if buffer[position] != rune('s') {
						goto l44
					}
					position++
					if buffer[position] != rune('t') {
						goto l44
					}
					position++
					if buffer[position] != rune('r') {
						goto l44
					}
					position++
					if buffer[position] != rune('u') {
						goto l44
					}
					position++
					if buffer[position] != rune('c') {
						goto l44
					}
					position++
					if buffer[position] != rune('t') {
						goto l44
					}
					position++
				}
			l46:
				depth--
				add(ruletype_label, position45)
			}
			return true
		l44:
			position, tokenIndex, depth = position44, tokenIndex44, depth44
			return false
		},
		/* 6 typename <- <(ident Action2)> */
		func() bool {
			position48, tokenIndex48, depth48 := position, tokenIndex, depth
			{
				position49 := position
				depth++
				if !_rules[ruleident]() {
					goto l48
				}
				if !_rules[ruleAction2]() {
					goto l48
				}
				depth--
				add(ruletypename, position49)
			}
			return true
		l48:
			position, tokenIndex, depth = position48, tokenIndex48, depth48
			return false
		},
		/* 7 metadata <- <('(' <(!')' .)*> ')')> */
		func() bool {
			position50, tokenIndex50, depth50 := position, tokenIndex, depth
			{
				position51 := position
				depth++
				if buffer[position] != rune('(') {
					goto l50
				}
				position++
				{
					position52 := position
					depth++
				l53:
					{
						position54, tokenIndex54, depth54 := position, tokenIndex, depth
						{
							position55, tokenIndex55, depth55 := position, tokenIndex, depth
							if buffer[position] != rune(')') {
								goto l55
							}
							position++
							goto l54
						l55:
							position, tokenIndex, depth = position55, tokenIndex55, depth55
						}
						if !matchDot() {
							goto l54
						}
						goto l53
					l54:
						position, tokenIndex, depth = position54, tokenIndex54, depth54
					}
					depth--
					add(rulePegText, position52)
				}
				if buffer[position] != rune(')') {
					goto l50
				}
				position++
				depth--
				add(rulemetadata, position51)
			}
			return true
		l50:
			position, tokenIndex, depth = position50, tokenIndex50, depth50
			return false
		},
		/* 8 field_decl <- <((spacing field_type ':' type metadata* ';' spacing Action3) / (spacing field_type ':' type <(' ' / '\t')*> '=' <(' ' / '\t')*> scalar metadata* ';' spacing Action4))> */
		func() bool {
			position56, tokenIndex56, depth56 := position, tokenIndex, depth
			{
				position57 := position
				depth++
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if !_rules[rulespacing]() {
						goto l59
					}
					if !_rules[rulefield_type]() {
						goto l59
					}
					if buffer[position] != rune(':') {
						goto l59
					}
					position++
					if !_rules[ruletype]() {
						goto l59
					}
				l60:
					{
						position61, tokenIndex61, depth61 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l61
						}
						goto l60
					l61:
						position, tokenIndex, depth = position61, tokenIndex61, depth61
					}
					if buffer[position] != rune(';') {
						goto l59
					}
					position++
					if !_rules[rulespacing]() {
						goto l59
					}
					if !_rules[ruleAction3]() {
						goto l59
					}
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if !_rules[rulespacing]() {
						goto l56
					}
					if !_rules[rulefield_type]() {
						goto l56
					}
					if buffer[position] != rune(':') {
						goto l56
					}
					position++
					if !_rules[ruletype]() {
						goto l56
					}
					{
						position62 := position
						depth++
					l63:
						{
							position64, tokenIndex64, depth64 := position, tokenIndex, depth
							{
								position65, tokenIndex65, depth65 := position, tokenIndex, depth
								if buffer[position] != rune(' ') {
									goto l66
								}
								position++
								goto l65
							l66:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('\t') {
									goto l64
								}
								position++
							}
						l65:
							goto l63
						l64:
							position, tokenIndex, depth = position64, tokenIndex64, depth64
						}
						depth--
						add(rulePegText, position62)
					}
					if buffer[position] != rune('=') {
						goto l56
					}
					position++
					{
						position67 := position
						depth++
					l68:
						{
							position69, tokenIndex69, depth69 := position, tokenIndex, depth
							{
								position70, tokenIndex70, depth70 := position, tokenIndex, depth
								if buffer[position] != rune(' ') {
									goto l71
								}
								position++
								goto l70
							l71:
								position, tokenIndex, depth = position70, tokenIndex70, depth70
								if buffer[position] != rune('\t') {
									goto l69
								}
								position++
							}
						l70:
							goto l68
						l69:
							position, tokenIndex, depth = position69, tokenIndex69, depth69
						}
						depth--
						add(rulePegText, position67)
					}
					if !_rules[rulescalar]() {
						goto l56
					}
				l72:
					{
						position73, tokenIndex73, depth73 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l73
						}
						goto l72
					l73:
						position, tokenIndex, depth = position73, tokenIndex73, depth73
					}
					if buffer[position] != rune(';') {
						goto l56
					}
					position++
					if !_rules[rulespacing]() {
						goto l56
					}
					if !_rules[ruleAction4]() {
						goto l56
					}
				}
			l58:
				depth--
				add(rulefield_decl, position57)
			}
			return true
		l56:
			position, tokenIndex, depth = position56, tokenIndex56, depth56
			return false
		},
		/* 9 field_type <- <(ident Action5)> */
		func() bool {
			position74, tokenIndex74, depth74 := position, tokenIndex, depth
			{
				position75 := position
				depth++
				if !_rules[ruleident]() {
					goto l74
				}
				if !_rules[ruleAction5]() {
					goto l74
				}
				depth--
				add(rulefield_type, position75)
			}
			return true
		l74:
			position, tokenIndex, depth = position74, tokenIndex74, depth74
			return false
		},
		/* 10 enum_decl <- <(('e' 'n' 'u' 'm' spacing ident spacing metadata* '{' enum_fields '}' spacing) / ('e' 'n' 'u' 'm' spacing ident ':' type spacing metadata* '{' enum_fields '}' spacing))> */
		func() bool {
			position76, tokenIndex76, depth76 := position, tokenIndex, depth
			{
				position77 := position
				depth++
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l79
					}
					position++
					if buffer[position] != rune('n') {
						goto l79
					}
					position++
					if buffer[position] != rune('u') {
						goto l79
					}
					position++
					if buffer[position] != rune('m') {
						goto l79
					}
					position++
					if !_rules[rulespacing]() {
						goto l79
					}
					if !_rules[ruleident]() {
						goto l79
					}
					if !_rules[rulespacing]() {
						goto l79
					}
				l80:
					{
						position81, tokenIndex81, depth81 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l81
						}
						goto l80
					l81:
						position, tokenIndex, depth = position81, tokenIndex81, depth81
					}
					if buffer[position] != rune('{') {
						goto l79
					}
					position++
					if !_rules[ruleenum_fields]() {
						goto l79
					}
					if buffer[position] != rune('}') {
						goto l79
					}
					position++
					if !_rules[rulespacing]() {
						goto l79
					}
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if buffer[position] != rune('e') {
						goto l76
					}
					position++
					if buffer[position] != rune('n') {
						goto l76
					}
					position++
					if buffer[position] != rune('u') {
						goto l76
					}
					position++
					if buffer[position] != rune('m') {
						goto l76
					}
					position++
					if !_rules[rulespacing]() {
						goto l76
					}
					if !_rules[ruleident]() {
						goto l76
					}
					if buffer[position] != rune(':') {
						goto l76
					}
					position++
					if !_rules[ruletype]() {
						goto l76
					}
					if !_rules[rulespacing]() {
						goto l76
					}
				l82:
					{
						position83, tokenIndex83, depth83 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l83
						}
						goto l82
					l83:
						position, tokenIndex, depth = position83, tokenIndex83, depth83
					}
					if buffer[position] != rune('{') {
						goto l76
					}
					position++
					if !_rules[ruleenum_fields]() {
						goto l76
					}
					if buffer[position] != rune('}') {
						goto l76
					}
					position++
					if !_rules[rulespacing]() {
						goto l76
					}
				}
			l78:
				depth--
				add(ruleenum_decl, position77)
			}
			return true
		l76:
			position, tokenIndex, depth = position76, tokenIndex76, depth76
			return false
		},
		/* 11 enum_fields <- <(enum_field / (enum_field ',' enum_fields))> */
		func() bool {
			position84, tokenIndex84, depth84 := position, tokenIndex, depth
			{
				position85 := position
				depth++
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if !_rules[ruleenum_field]() {
						goto l87
					}
					goto l86
				l87:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
					if !_rules[ruleenum_field]() {
						goto l84
					}
					if buffer[position] != rune(',') {
						goto l84
					}
					position++
					if !_rules[ruleenum_fields]() {
						goto l84
					}
				}
			l86:
				depth--
				add(ruleenum_fields, position85)
			}
			return true
		l84:
			position, tokenIndex, depth = position84, tokenIndex84, depth84
			return false
		},
		/* 12 union_decl <- <('u' 'n' 'i' 'o' 'n' spacing ident spacing metadata* '{' spacing enum_fields '}' spacing Action6)> */
		func() bool {
			position88, tokenIndex88, depth88 := position, tokenIndex, depth
			{
				position89 := position
				depth++
				if buffer[position] != rune('u') {
					goto l88
				}
				position++
				if buffer[position] != rune('n') {
					goto l88
				}
				position++
				if buffer[position] != rune('i') {
					goto l88
				}
				position++
				if buffer[position] != rune('o') {
					goto l88
				}
				position++
				if buffer[position] != rune('n') {
					goto l88
				}
				position++
				if !_rules[rulespacing]() {
					goto l88
				}
				if !_rules[ruleident]() {
					goto l88
				}
				if !_rules[rulespacing]() {
					goto l88
				}
			l90:
				{
					position91, tokenIndex91, depth91 := position, tokenIndex, depth
					if !_rules[rulemetadata]() {
						goto l91
					}
					goto l90
				l91:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
				}
				if buffer[position] != rune('{') {
					goto l88
				}
				position++
				if !_rules[rulespacing]() {
					goto l88
				}
				if !_rules[ruleenum_fields]() {
					goto l88
				}
				if buffer[position] != rune('}') {
					goto l88
				}
				position++
				if !_rules[rulespacing]() {
					goto l88
				}
				if !_rules[ruleAction6]() {
					goto l88
				}
				depth--
				add(ruleunion_decl, position89)
			}
			return true
		l88:
			position, tokenIndex, depth = position88, tokenIndex88, depth88
			return false
		},
		/* 13 enum_field <- <((enum_field_type spacing Action7) / (enum_field_type spacing '=' spacing integer_constant spacing))> */
		func() bool {
			position92, tokenIndex92, depth92 := position, tokenIndex, depth
			{
				position93 := position
				depth++
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					if !_rules[ruleenum_field_type]() {
						goto l95
					}
					if !_rules[rulespacing]() {
						goto l95
					}
					if !_rules[ruleAction7]() {
						goto l95
					}
					goto l94
				l95:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
					if !_rules[ruleenum_field_type]() {
						goto l92
					}
					if !_rules[rulespacing]() {
						goto l92
					}
					if buffer[position] != rune('=') {
						goto l92
					}
					position++
					if !_rules[rulespacing]() {
						goto l92
					}
					if !_rules[ruleinteger_constant]() {
						goto l92
					}
					if !_rules[rulespacing]() {
						goto l92
					}
				}
			l94:
				depth--
				add(ruleenum_field, position93)
			}
			return true
		l92:
			position, tokenIndex, depth = position92, tokenIndex92, depth92
			return false
		},
		/* 14 enum_field_type <- <(ident Action8)> */
		func() bool {
			position96, tokenIndex96, depth96 := position, tokenIndex, depth
			{
				position97 := position
				depth++
				if !_rules[ruleident]() {
					goto l96
				}
				if !_rules[ruleAction8]() {
					goto l96
				}
				depth--
				add(ruleenum_field_type, position97)
			}
			return true
		l96:
			position, tokenIndex, depth = position96, tokenIndex96, depth96
			return false
		},
		/* 15 root_decl <- <('r' 'o' 'o' 't' '_' 't' 'y' 'p' 'e' spacing ident spacing ';' spacing)> */
		func() bool {
			position98, tokenIndex98, depth98 := position, tokenIndex, depth
			{
				position99 := position
				depth++
				if buffer[position] != rune('r') {
					goto l98
				}
				position++
				if buffer[position] != rune('o') {
					goto l98
				}
				position++
				if buffer[position] != rune('o') {
					goto l98
				}
				position++
				if buffer[position] != rune('t') {
					goto l98
				}
				position++
				if buffer[position] != rune('_') {
					goto l98
				}
				position++
				if buffer[position] != rune('t') {
					goto l98
				}
				position++
				if buffer[position] != rune('y') {
					goto l98
				}
				position++
				if buffer[position] != rune('p') {
					goto l98
				}
				position++
				if buffer[position] != rune('e') {
					goto l98
				}
				position++
				if !_rules[rulespacing]() {
					goto l98
				}
				if !_rules[ruleident]() {
					goto l98
				}
				if !_rules[rulespacing]() {
					goto l98
				}
				if buffer[position] != rune(';') {
					goto l98
				}
				position++
				if !_rules[rulespacing]() {
					goto l98
				}
				depth--
				add(ruleroot_decl, position99)
			}
			return true
		l98:
			position, tokenIndex, depth = position98, tokenIndex98, depth98
			return false
		},
		/* 16 file_extension_decl <- <('f' 'i' 'l' 'e' '_' 'e' 'x' 't' 'e' 'n' 's' 'i' 'o' 'n' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position100, tokenIndex100, depth100 := position, tokenIndex, depth
			{
				position101 := position
				depth++
				if buffer[position] != rune('f') {
					goto l100
				}
				position++
				if buffer[position] != rune('i') {
					goto l100
				}
				position++
				if buffer[position] != rune('l') {
					goto l100
				}
				position++
				if buffer[position] != rune('e') {
					goto l100
				}
				position++
				if buffer[position] != rune('_') {
					goto l100
				}
				position++
				if buffer[position] != rune('e') {
					goto l100
				}
				position++
				if buffer[position] != rune('x') {
					goto l100
				}
				position++
				if buffer[position] != rune('t') {
					goto l100
				}
				position++
				if buffer[position] != rune('e') {
					goto l100
				}
				position++
				if buffer[position] != rune('n') {
					goto l100
				}
				position++
				if buffer[position] != rune('s') {
					goto l100
				}
				position++
				if buffer[position] != rune('i') {
					goto l100
				}
				position++
				if buffer[position] != rune('o') {
					goto l100
				}
				position++
				if buffer[position] != rune('n') {
					goto l100
				}
				position++
				{
					position102 := position
					depth++
				l103:
					{
						position104, tokenIndex104, depth104 := position, tokenIndex, depth
						{
							position105, tokenIndex105, depth105 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l106
							}
							position++
							goto l105
						l106:
							position, tokenIndex, depth = position105, tokenIndex105, depth105
							if buffer[position] != rune('\t') {
								goto l104
							}
							position++
						}
					l105:
						goto l103
					l104:
						position, tokenIndex, depth = position104, tokenIndex104, depth104
					}
					depth--
					add(rulePegText, position102)
				}
				{
					position107 := position
					depth++
					{
						position110, tokenIndex110, depth110 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l110
						}
						position++
						goto l100
					l110:
						position, tokenIndex, depth = position110, tokenIndex110, depth110
					}
					if !matchDot() {
						goto l100
					}
				l108:
					{
						position109, tokenIndex109, depth109 := position, tokenIndex, depth
						{
							position111, tokenIndex111, depth111 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l111
							}
							position++
							goto l109
						l111:
							position, tokenIndex, depth = position111, tokenIndex111, depth111
						}
						if !matchDot() {
							goto l109
						}
						goto l108
					l109:
						position, tokenIndex, depth = position109, tokenIndex109, depth109
					}
					depth--
					add(rulePegText, position107)
				}
				if buffer[position] != rune(';') {
					goto l100
				}
				position++
				if !_rules[rulespacing]() {
					goto l100
				}
				depth--
				add(rulefile_extension_decl, position101)
			}
			return true
		l100:
			position, tokenIndex, depth = position100, tokenIndex100, depth100
			return false
		},
		/* 17 file_identifier_decl <- <('f' 'i' 'l' 'e' '_' 'i' 'd' 'e' 'n' 't' 'i' 'f' 'i' 'e' 'r' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position112, tokenIndex112, depth112 := position, tokenIndex, depth
			{
				position113 := position
				depth++
				if buffer[position] != rune('f') {
					goto l112
				}
				position++
				if buffer[position] != rune('i') {
					goto l112
				}
				position++
				if buffer[position] != rune('l') {
					goto l112
				}
				position++
				if buffer[position] != rune('e') {
					goto l112
				}
				position++
				if buffer[position] != rune('_') {
					goto l112
				}
				position++
				if buffer[position] != rune('i') {
					goto l112
				}
				position++
				if buffer[position] != rune('d') {
					goto l112
				}
				position++
				if buffer[position] != rune('e') {
					goto l112
				}
				position++
				if buffer[position] != rune('n') {
					goto l112
				}
				position++
				if buffer[position] != rune('t') {
					goto l112
				}
				position++
				if buffer[position] != rune('i') {
					goto l112
				}
				position++
				if buffer[position] != rune('f') {
					goto l112
				}
				position++
				if buffer[position] != rune('i') {
					goto l112
				}
				position++
				if buffer[position] != rune('e') {
					goto l112
				}
				position++
				if buffer[position] != rune('r') {
					goto l112
				}
				position++
				{
					position114 := position
					depth++
				l115:
					{
						position116, tokenIndex116, depth116 := position, tokenIndex, depth
						{
							position117, tokenIndex117, depth117 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l118
							}
							position++
							goto l117
						l118:
							position, tokenIndex, depth = position117, tokenIndex117, depth117
							if buffer[position] != rune('\t') {
								goto l116
							}
							position++
						}
					l117:
						goto l115
					l116:
						position, tokenIndex, depth = position116, tokenIndex116, depth116
					}
					depth--
					add(rulePegText, position114)
				}
				{
					position119 := position
					depth++
					{
						position122, tokenIndex122, depth122 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l122
						}
						position++
						goto l112
					l122:
						position, tokenIndex, depth = position122, tokenIndex122, depth122
					}
					if !matchDot() {
						goto l112
					}
				l120:
					{
						position121, tokenIndex121, depth121 := position, tokenIndex, depth
						{
							position123, tokenIndex123, depth123 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l123
							}
							position++
							goto l121
						l123:
							position, tokenIndex, depth = position123, tokenIndex123, depth123
						}
						if !matchDot() {
							goto l121
						}
						goto l120
					l121:
						position, tokenIndex, depth = position121, tokenIndex121, depth121
					}
					depth--
					add(rulePegText, position119)
				}
				if buffer[position] != rune(';') {
					goto l112
				}
				position++
				if !_rules[rulespacing]() {
					goto l112
				}
				depth--
				add(rulefile_identifier_decl, position113)
			}
			return true
		l112:
			position, tokenIndex, depth = position112, tokenIndex112, depth112
			return false
		},
		/* 18 attribute_decl <- <('a' 't' 't' 'r' 'i' 'b' 'u' 't' 'e' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position124, tokenIndex124, depth124 := position, tokenIndex, depth
			{
				position125 := position
				depth++
				if buffer[position] != rune('a') {
					goto l124
				}
				position++
				if buffer[position] != rune('t') {
					goto l124
				}
				position++
				if buffer[position] != rune('t') {
					goto l124
				}
				position++
				if buffer[position] != rune('r') {
					goto l124
				}
				position++
				if buffer[position] != rune('i') {
					goto l124
				}
				position++
				if buffer[position] != rune('b') {
					goto l124
				}
				position++
				if buffer[position] != rune('u') {
					goto l124
				}
				position++
				if buffer[position] != rune('t') {
					goto l124
				}
				position++
				if buffer[position] != rune('e') {
					goto l124
				}
				position++
				{
					position126 := position
					depth++
				l127:
					{
						position128, tokenIndex128, depth128 := position, tokenIndex, depth
						{
							position129, tokenIndex129, depth129 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l130
							}
							position++
							goto l129
						l130:
							position, tokenIndex, depth = position129, tokenIndex129, depth129
							if buffer[position] != rune('\t') {
								goto l128
							}
							position++
						}
					l129:
						goto l127
					l128:
						position, tokenIndex, depth = position128, tokenIndex128, depth128
					}
					depth--
					add(rulePegText, position126)
				}
				{
					position131 := position
					depth++
					{
						position134, tokenIndex134, depth134 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l134
						}
						position++
						goto l124
					l134:
						position, tokenIndex, depth = position134, tokenIndex134, depth134
					}
					if !matchDot() {
						goto l124
					}
				l132:
					{
						position133, tokenIndex133, depth133 := position, tokenIndex, depth
						{
							position135, tokenIndex135, depth135 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l135
							}
							position++
							goto l133
						l135:
							position, tokenIndex, depth = position135, tokenIndex135, depth135
						}
						if !matchDot() {
							goto l133
						}
						goto l132
					l133:
						position, tokenIndex, depth = position133, tokenIndex133, depth133
					}
					depth--
					add(rulePegText, position131)
				}
				if buffer[position] != rune(';') {
					goto l124
				}
				position++
				if !_rules[rulespacing]() {
					goto l124
				}
				depth--
				add(ruleattribute_decl, position125)
			}
			return true
		l124:
			position, tokenIndex, depth = position124, tokenIndex124, depth124
			return false
		},
		/* 19 rpc_decl <- <('r' 'p' 'c' '_' 's' 'e' 'r' 'v' 'i' 'c' 'e' <(' ' / '\t')*> ident '{' <(!'}' .)+> '}' spacing)> */
		func() bool {
			position136, tokenIndex136, depth136 := position, tokenIndex, depth
			{
				position137 := position
				depth++
				if buffer[position] != rune('r') {
					goto l136
				}
				position++
				if buffer[position] != rune('p') {
					goto l136
				}
				position++
				if buffer[position] != rune('c') {
					goto l136
				}
				position++
				if buffer[position] != rune('_') {
					goto l136
				}
				position++
				if buffer[position] != rune('s') {
					goto l136
				}
				position++
				if buffer[position] != rune('e') {
					goto l136
				}
				position++
				if buffer[position] != rune('r') {
					goto l136
				}
				position++
				if buffer[position] != rune('v') {
					goto l136
				}
				position++
				if buffer[position] != rune('i') {
					goto l136
				}
				position++
				if buffer[position] != rune('c') {
					goto l136
				}
				position++
				if buffer[position] != rune('e') {
					goto l136
				}
				position++
				{
					position138 := position
					depth++
				l139:
					{
						position140, tokenIndex140, depth140 := position, tokenIndex, depth
						{
							position141, tokenIndex141, depth141 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l142
							}
							position++
							goto l141
						l142:
							position, tokenIndex, depth = position141, tokenIndex141, depth141
							if buffer[position] != rune('\t') {
								goto l140
							}
							position++
						}
					l141:
						goto l139
					l140:
						position, tokenIndex, depth = position140, tokenIndex140, depth140
					}
					depth--
					add(rulePegText, position138)
				}
				if !_rules[ruleident]() {
					goto l136
				}
				if buffer[position] != rune('{') {
					goto l136
				}
				position++
				{
					position143 := position
					depth++
					{
						position146, tokenIndex146, depth146 := position, tokenIndex, depth
						if buffer[position] != rune('}') {
							goto l146
						}
						position++
						goto l136
					l146:
						position, tokenIndex, depth = position146, tokenIndex146, depth146
					}
					if !matchDot() {
						goto l136
					}
				l144:
					{
						position145, tokenIndex145, depth145 := position, tokenIndex, depth
						{
							position147, tokenIndex147, depth147 := position, tokenIndex, depth
							if buffer[position] != rune('}') {
								goto l147
							}
							position++
							goto l145
						l147:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
						}
						if !matchDot() {
							goto l145
						}
						goto l144
					l145:
						position, tokenIndex, depth = position145, tokenIndex145, depth145
					}
					depth--
					add(rulePegText, position143)
				}
				if buffer[position] != rune('}') {
					goto l136
				}
				position++
				if !_rules[rulespacing]() {
					goto l136
				}
				depth--
				add(rulerpc_decl, position137)
			}
			return true
		l136:
			position, tokenIndex, depth = position136, tokenIndex136, depth136
			return false
		},
		/* 20 type <- <(('b' 'o' 'o' 'l' spacing Action9) / ('b' 'y' 't' 'e' spacing Action10) / ('u' 'b' 'y' 't' 'e' spacing Action11) / ('s' 'h' 'o' 'r' 't' spacing Action12) / ('u' 's' 'h' 'o' 'r' 't' spacing Action13) / ('i' 'n' 't' spacing Action14) / ('u' 'i' 'n' 't' spacing Action15) / ('f' 'l' 'o' 'a' 't' spacing Action16) / ('l' 'o' 'n' 'g' spacing Action17) / ('u' 'l' 'o' 'n' 'g' spacing Action18) / ('d' 'o' 'u' 'b' 'l' 'e' spacing Action19) / ('i' 'n' 't' '8' spacing Action20) / ('i' 'n' 't' '1' '6' spacing Action21) / ('u' 'i' 'n' 't' '1' '6' spacing Action22) / ('i' 'n' 't' '3' '2' spacing Action23) / ('u' 'i' 'n' 't' '3' '2' spacing Action24) / ('i' 'n' 't' '6' '4' spacing Action25) / ('u' 'i' 'n' 't' '6' '4' spacing Action26) / ('f' 'l' 'o' 'a' 't' '3' '2' spacing Action27) / ('f' 'l' 'o' 'a' 't' '6' '4' spacing Action28) / ('s' 't' 'r' 'i' 'n' 'g' spacing Action29) / (ident spacing Action30) / ('[' type ']' spacing Action31))> */
		func() bool {
			position148, tokenIndex148, depth148 := position, tokenIndex, depth
			{
				position149 := position
				depth++
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l151
					}
					position++
					if buffer[position] != rune('o') {
						goto l151
					}
					position++
					if buffer[position] != rune('o') {
						goto l151
					}
					position++
					if buffer[position] != rune('l') {
						goto l151
					}
					position++
					if !_rules[rulespacing]() {
						goto l151
					}
					if !_rules[ruleAction9]() {
						goto l151
					}
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('b') {
						goto l152
					}
					position++
					if buffer[position] != rune('y') {
						goto l152
					}
					position++
					if buffer[position] != rune('t') {
						goto l152
					}
					position++
					if buffer[position] != rune('e') {
						goto l152
					}
					position++
					if !_rules[rulespacing]() {
						goto l152
					}
					if !_rules[ruleAction10]() {
						goto l152
					}
					goto l150
				l152:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('u') {
						goto l153
					}
					position++
					if buffer[position] != rune('b') {
						goto l153
					}
					position++
					if buffer[position] != rune('y') {
						goto l153
					}
					position++
					if buffer[position] != rune('t') {
						goto l153
					}
					position++
					if buffer[position] != rune('e') {
						goto l153
					}
					position++
					if !_rules[rulespacing]() {
						goto l153
					}
					if !_rules[ruleAction11]() {
						goto l153
					}
					goto l150
				l153:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('s') {
						goto l154
					}
					position++
					if buffer[position] != rune('h') {
						goto l154
					}
					position++
					if buffer[position] != rune('o') {
						goto l154
					}
					position++
					if buffer[position] != rune('r') {
						goto l154
					}
					position++
					if buffer[position] != rune('t') {
						goto l154
					}
					position++
					if !_rules[rulespacing]() {
						goto l154
					}
					if !_rules[ruleAction12]() {
						goto l154
					}
					goto l150
				l154:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('u') {
						goto l155
					}
					position++
					if buffer[position] != rune('s') {
						goto l155
					}
					position++
					if buffer[position] != rune('h') {
						goto l155
					}
					position++
					if buffer[position] != rune('o') {
						goto l155
					}
					position++
					if buffer[position] != rune('r') {
						goto l155
					}
					position++
					if buffer[position] != rune('t') {
						goto l155
					}
					position++
					if !_rules[rulespacing]() {
						goto l155
					}
					if !_rules[ruleAction13]() {
						goto l155
					}
					goto l150
				l155:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('i') {
						goto l156
					}
					position++
					if buffer[position] != rune('n') {
						goto l156
					}
					position++
					if buffer[position] != rune('t') {
						goto l156
					}
					position++
					if !_rules[rulespacing]() {
						goto l156
					}
					if !_rules[ruleAction14]() {
						goto l156
					}
					goto l150
				l156:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('u') {
						goto l157
					}
					position++
					if buffer[position] != rune('i') {
						goto l157
					}
					position++
					if buffer[position] != rune('n') {
						goto l157
					}
					position++
					if buffer[position] != rune('t') {
						goto l157
					}
					position++
					if !_rules[rulespacing]() {
						goto l157
					}
					if !_rules[ruleAction15]() {
						goto l157
					}
					goto l150
				l157:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('f') {
						goto l158
					}
					position++
					if buffer[position] != rune('l') {
						goto l158
					}
					position++
					if buffer[position] != rune('o') {
						goto l158
					}
					position++
					if buffer[position] != rune('a') {
						goto l158
					}
					position++
					if buffer[position] != rune('t') {
						goto l158
					}
					position++
					if !_rules[rulespacing]() {
						goto l158
					}
					if !_rules[ruleAction16]() {
						goto l158
					}
					goto l150
				l158:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('l') {
						goto l159
					}
					position++
					if buffer[position] != rune('o') {
						goto l159
					}
					position++
					if buffer[position] != rune('n') {
						goto l159
					}
					position++
					if buffer[position] != rune('g') {
						goto l159
					}
					position++
					if !_rules[rulespacing]() {
						goto l159
					}
					if !_rules[ruleAction17]() {
						goto l159
					}
					goto l150
				l159:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('u') {
						goto l160
					}
					position++
					if buffer[position] != rune('l') {
						goto l160
					}
					position++
					if buffer[position] != rune('o') {
						goto l160
					}
					position++
					if buffer[position] != rune('n') {
						goto l160
					}
					position++
					if buffer[position] != rune('g') {
						goto l160
					}
					position++
					if !_rules[rulespacing]() {
						goto l160
					}
					if !_rules[ruleAction18]() {
						goto l160
					}
					goto l150
				l160:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('d') {
						goto l161
					}
					position++
					if buffer[position] != rune('o') {
						goto l161
					}
					position++
					if buffer[position] != rune('u') {
						goto l161
					}
					position++
					if buffer[position] != rune('b') {
						goto l161
					}
					position++
					if buffer[position] != rune('l') {
						goto l161
					}
					position++
					if buffer[position] != rune('e') {
						goto l161
					}
					position++
					if !_rules[rulespacing]() {
						goto l161
					}
					if !_rules[ruleAction19]() {
						goto l161
					}
					goto l150
				l161:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
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
					if !_rules[ruleAction20]() {
						goto l162
					}
					goto l150
				l162:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
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
					if buffer[position] != rune('1') {
						goto l163
					}
					position++
					if buffer[position] != rune('6') {
						goto l163
					}
					position++
					if !_rules[rulespacing]() {
						goto l163
					}
					if !_rules[ruleAction21]() {
						goto l163
					}
					goto l150
				l163:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('u') {
						goto l164
					}
					position++
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
					if !_rules[ruleAction22]() {
						goto l164
					}
					goto l150
				l164:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
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
					if buffer[position] != rune('3') {
						goto l165
					}
					position++
					if buffer[position] != rune('2') {
						goto l165
					}
					position++
					if !_rules[rulespacing]() {
						goto l165
					}
					if !_rules[ruleAction23]() {
						goto l165
					}
					goto l150
				l165:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('u') {
						goto l166
					}
					position++
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
					if !_rules[ruleAction24]() {
						goto l166
					}
					goto l150
				l166:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
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
					if buffer[position] != rune('6') {
						goto l167
					}
					position++
					if buffer[position] != rune('4') {
						goto l167
					}
					position++
					if !_rules[rulespacing]() {
						goto l167
					}
					if !_rules[ruleAction25]() {
						goto l167
					}
					goto l150
				l167:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('u') {
						goto l168
					}
					position++
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
					if !_rules[ruleAction26]() {
						goto l168
					}
					goto l150
				l168:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('f') {
						goto l169
					}
					position++
					if buffer[position] != rune('l') {
						goto l169
					}
					position++
					if buffer[position] != rune('o') {
						goto l169
					}
					position++
					if buffer[position] != rune('a') {
						goto l169
					}
					position++
					if buffer[position] != rune('t') {
						goto l169
					}
					position++
					if buffer[position] != rune('3') {
						goto l169
					}
					position++
					if buffer[position] != rune('2') {
						goto l169
					}
					position++
					if !_rules[rulespacing]() {
						goto l169
					}
					if !_rules[ruleAction27]() {
						goto l169
					}
					goto l150
				l169:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
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
					if buffer[position] != rune('6') {
						goto l170
					}
					position++
					if buffer[position] != rune('4') {
						goto l170
					}
					position++
					if !_rules[rulespacing]() {
						goto l170
					}
					if !_rules[ruleAction28]() {
						goto l170
					}
					goto l150
				l170:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('s') {
						goto l171
					}
					position++
					if buffer[position] != rune('t') {
						goto l171
					}
					position++
					if buffer[position] != rune('r') {
						goto l171
					}
					position++
					if buffer[position] != rune('i') {
						goto l171
					}
					position++
					if buffer[position] != rune('n') {
						goto l171
					}
					position++
					if buffer[position] != rune('g') {
						goto l171
					}
					position++
					if !_rules[rulespacing]() {
						goto l171
					}
					if !_rules[ruleAction29]() {
						goto l171
					}
					goto l150
				l171:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if !_rules[ruleident]() {
						goto l172
					}
					if !_rules[rulespacing]() {
						goto l172
					}
					if !_rules[ruleAction30]() {
						goto l172
					}
					goto l150
				l172:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('[') {
						goto l148
					}
					position++
					if !_rules[ruletype]() {
						goto l148
					}
					if buffer[position] != rune(']') {
						goto l148
					}
					position++
					if !_rules[rulespacing]() {
						goto l148
					}
					if !_rules[ruleAction31]() {
						goto l148
					}
				}
			l150:
				depth--
				add(ruletype, position149)
			}
			return true
		l148:
			position, tokenIndex, depth = position148, tokenIndex148, depth148
			return false
		},
		/* 21 scalar <- <(integer_constant / float_constant)> */
		func() bool {
			position173, tokenIndex173, depth173 := position, tokenIndex, depth
			{
				position174 := position
				depth++
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if !_rules[ruleinteger_constant]() {
						goto l176
					}
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if !_rules[rulefloat_constant]() {
						goto l173
					}
				}
			l175:
				depth--
				add(rulescalar, position174)
			}
			return true
		l173:
			position, tokenIndex, depth = position173, tokenIndex173, depth173
			return false
		},
		/* 22 integer_constant <- <(<[0-9]+> / ('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		func() bool {
			position177, tokenIndex177, depth177 := position, tokenIndex, depth
			{
				position178 := position
				depth++
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					{
						position181 := position
						depth++
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l180
						}
						position++
					l182:
						{
							position183, tokenIndex183, depth183 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l183
							}
							position++
							goto l182
						l183:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
						}
						depth--
						add(rulePegText, position181)
					}
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('t') {
						goto l184
					}
					position++
					if buffer[position] != rune('r') {
						goto l184
					}
					position++
					if buffer[position] != rune('u') {
						goto l184
					}
					position++
					if buffer[position] != rune('e') {
						goto l184
					}
					position++
					goto l179
				l184:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('f') {
						goto l177
					}
					position++
					if buffer[position] != rune('a') {
						goto l177
					}
					position++
					if buffer[position] != rune('l') {
						goto l177
					}
					position++
					if buffer[position] != rune('s') {
						goto l177
					}
					position++
					if buffer[position] != rune('e') {
						goto l177
					}
					position++
				}
			l179:
				depth--
				add(ruleinteger_constant, position178)
			}
			return true
		l177:
			position, tokenIndex, depth = position177, tokenIndex177, depth177
			return false
		},
		/* 23 float_constant <- <(<('-'* [0-9]+ . [0-9])> / float_constant_exp)> */
		func() bool {
			position185, tokenIndex185, depth185 := position, tokenIndex, depth
			{
				position186 := position
				depth++
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					{
						position189 := position
						depth++
					l190:
						{
							position191, tokenIndex191, depth191 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l191
							}
							position++
							goto l190
						l191:
							position, tokenIndex, depth = position191, tokenIndex191, depth191
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l188
						}
						position++
					l192:
						{
							position193, tokenIndex193, depth193 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l193
							}
							position++
							goto l192
						l193:
							position, tokenIndex, depth = position193, tokenIndex193, depth193
						}
						if !matchDot() {
							goto l188
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l188
						}
						position++
						depth--
						add(rulePegText, position189)
					}
					goto l187
				l188:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
					if !_rules[rulefloat_constant_exp]() {
						goto l185
					}
				}
			l187:
				depth--
				add(rulefloat_constant, position186)
			}
			return true
		l185:
			position, tokenIndex, depth = position185, tokenIndex185, depth185
			return false
		},
		/* 24 float_constant_exp <- <(<('-'* [0-9]+ . [0-9]+)> <('e' / 'E')> <([+-]] / '>' / ' ' / '<' / '[' / [0-9])+>)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
				{
					position196 := position
					depth++
				l197:
					{
						position198, tokenIndex198, depth198 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l198
						}
						position++
						goto l197
					l198:
						position, tokenIndex, depth = position198, tokenIndex198, depth198
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l194
					}
					position++
				l199:
					{
						position200, tokenIndex200, depth200 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l200
						}
						position++
						goto l199
					l200:
						position, tokenIndex, depth = position200, tokenIndex200, depth200
					}
					if !matchDot() {
						goto l194
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l194
					}
					position++
				l201:
					{
						position202, tokenIndex202, depth202 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l202
						}
						position++
						goto l201
					l202:
						position, tokenIndex, depth = position202, tokenIndex202, depth202
					}
					depth--
					add(rulePegText, position196)
				}
				{
					position203 := position
					depth++
					{
						position204, tokenIndex204, depth204 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l205
						}
						position++
						goto l204
					l205:
						position, tokenIndex, depth = position204, tokenIndex204, depth204
						if buffer[position] != rune('E') {
							goto l194
						}
						position++
					}
				l204:
					depth--
					add(rulePegText, position203)
				}
				{
					position206 := position
					depth++
					{
						position209, tokenIndex209, depth209 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('+') || c > rune(']') {
							goto l210
						}
						position++
						goto l209
					l210:
						position, tokenIndex, depth = position209, tokenIndex209, depth209
						if buffer[position] != rune('>') {
							goto l211
						}
						position++
						goto l209
					l211:
						position, tokenIndex, depth = position209, tokenIndex209, depth209
						if buffer[position] != rune(' ') {
							goto l212
						}
						position++
						goto l209
					l212:
						position, tokenIndex, depth = position209, tokenIndex209, depth209
						if buffer[position] != rune('<') {
							goto l213
						}
						position++
						goto l209
					l213:
						position, tokenIndex, depth = position209, tokenIndex209, depth209
						if buffer[position] != rune('[') {
							goto l214
						}
						position++
						goto l209
					l214:
						position, tokenIndex, depth = position209, tokenIndex209, depth209
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l194
						}
						position++
					}
				l209:
				l207:
					{
						position208, tokenIndex208, depth208 := position, tokenIndex, depth
						{
							position215, tokenIndex215, depth215 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('+') || c > rune(']') {
								goto l216
							}
							position++
							goto l215
						l216:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
							if buffer[position] != rune('>') {
								goto l217
							}
							position++
							goto l215
						l217:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
							if buffer[position] != rune(' ') {
								goto l218
							}
							position++
							goto l215
						l218:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
							if buffer[position] != rune('<') {
								goto l219
							}
							position++
							goto l215
						l219:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
							if buffer[position] != rune('[') {
								goto l220
							}
							position++
							goto l215
						l220:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l208
							}
							position++
						}
					l215:
						goto l207
					l208:
						position, tokenIndex, depth = position208, tokenIndex208, depth208
					}
					depth--
					add(rulePegText, position206)
				}
				depth--
				add(rulefloat_constant_exp, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 25 ident <- <<(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position221, tokenIndex221, depth221 := position, tokenIndex, depth
			{
				position222 := position
				depth++
				{
					position223 := position
					depth++
					{
						position224, tokenIndex224, depth224 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l225
						}
						position++
						goto l224
					l225:
						position, tokenIndex, depth = position224, tokenIndex224, depth224
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l226
						}
						position++
						goto l224
					l226:
						position, tokenIndex, depth = position224, tokenIndex224, depth224
						if buffer[position] != rune('_') {
							goto l221
						}
						position++
					}
				l224:
				l227:
					{
						position228, tokenIndex228, depth228 := position, tokenIndex, depth
						{
							position229, tokenIndex229, depth229 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l230
							}
							position++
							goto l229
						l230:
							position, tokenIndex, depth = position229, tokenIndex229, depth229
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l231
							}
							position++
							goto l229
						l231:
							position, tokenIndex, depth = position229, tokenIndex229, depth229
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l232
							}
							position++
							goto l229
						l232:
							position, tokenIndex, depth = position229, tokenIndex229, depth229
							if buffer[position] != rune('_') {
								goto l228
							}
							position++
						}
					l229:
						goto l227
					l228:
						position, tokenIndex, depth = position228, tokenIndex228, depth228
					}
					depth--
					add(rulePegText, position223)
				}
				depth--
				add(ruleident, position222)
			}
			return true
		l221:
			position, tokenIndex, depth = position221, tokenIndex221, depth221
			return false
		},
		/* 26 only_comment <- <(spacing ';')> */
		func() bool {
			position233, tokenIndex233, depth233 := position, tokenIndex, depth
			{
				position234 := position
				depth++
				if !_rules[rulespacing]() {
					goto l233
				}
				if buffer[position] != rune(';') {
					goto l233
				}
				position++
				depth--
				add(ruleonly_comment, position234)
			}
			return true
		l233:
			position, tokenIndex, depth = position233, tokenIndex233, depth233
			return false
		},
		/* 27 spacing <- <space_comment*> */
		func() bool {
			{
				position236 := position
				depth++
			l237:
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if !_rules[rulespace_comment]() {
						goto l238
					}
					goto l237
				l238:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
				}
				depth--
				add(rulespacing, position236)
			}
			return true
		},
		/* 28 space_comment <- <(space / comment)> */
		func() bool {
			position239, tokenIndex239, depth239 := position, tokenIndex, depth
			{
				position240 := position
				depth++
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if !_rules[rulespace]() {
						goto l242
					}
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if !_rules[rulecomment]() {
						goto l239
					}
				}
			l241:
				depth--
				add(rulespace_comment, position240)
			}
			return true
		l239:
			position, tokenIndex, depth = position239, tokenIndex239, depth239
			return false
		},
		/* 29 comment <- <('/' '/' (!end_of_line .)* end_of_line)> */
		func() bool {
			position243, tokenIndex243, depth243 := position, tokenIndex, depth
			{
				position244 := position
				depth++
				if buffer[position] != rune('/') {
					goto l243
				}
				position++
				if buffer[position] != rune('/') {
					goto l243
				}
				position++
			l245:
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					{
						position247, tokenIndex247, depth247 := position, tokenIndex, depth
						if !_rules[ruleend_of_line]() {
							goto l247
						}
						goto l246
					l247:
						position, tokenIndex, depth = position247, tokenIndex247, depth247
					}
					if !matchDot() {
						goto l246
					}
					goto l245
				l246:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
				}
				if !_rules[ruleend_of_line]() {
					goto l243
				}
				depth--
				add(rulecomment, position244)
			}
			return true
		l243:
			position, tokenIndex, depth = position243, tokenIndex243, depth243
			return false
		},
		/* 30 space <- <(' ' / '\t' / end_of_line)> */
		func() bool {
			position248, tokenIndex248, depth248 := position, tokenIndex, depth
			{
				position249 := position
				depth++
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if buffer[position] != rune('\t') {
						goto l252
					}
					position++
					goto l250
				l252:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if !_rules[ruleend_of_line]() {
						goto l248
					}
				}
			l250:
				depth--
				add(rulespace, position249)
			}
			return true
		l248:
			position, tokenIndex, depth = position248, tokenIndex248, depth248
			return false
		},
		/* 31 end_of_line <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position253, tokenIndex253, depth253 := position, tokenIndex, depth
			{
				position254 := position
				depth++
				{
					position255, tokenIndex255, depth255 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l256
					}
					position++
					if buffer[position] != rune('\n') {
						goto l256
					}
					position++
					goto l255
				l256:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
					if buffer[position] != rune('\n') {
						goto l257
					}
					position++
					goto l255
				l257:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
					if buffer[position] != rune('\r') {
						goto l253
					}
					position++
				}
			l255:
				depth--
				add(ruleend_of_line, position254)
			}
			return true
		l253:
			position, tokenIndex, depth = position253, tokenIndex253, depth253
			return false
		},
		/* 32 end_of_file <- <!.> */
		func() bool {
			position258, tokenIndex258, depth258 := position, tokenIndex, depth
			{
				position259 := position
				depth++
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if !matchDot() {
						goto l260
					}
					goto l258
				l260:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
				}
				depth--
				add(ruleend_of_file, position259)
			}
			return true
		l258:
			position, tokenIndex, depth = position258, tokenIndex258, depth258
			return false
		},
		nil,
		/* 35 Action0 <- <{p.SetNameSpace(text)}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 36 Action1 <- <{p.ExtractStruct()}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 37 Action2 <- <{p.SetTypeName(text)}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 38 Action3 <- <{p.NewExtractField()}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 39 Action4 <- <{p.NewExtractFieldWithValue()}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 40 Action5 <- <{p.FieldNaame(text)}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 41 Action6 <- <{p.NewUnion(text)}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 42 Action7 <- <{p.NewExtractField()}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 43 Action8 <- <{p.FieldNaame(text)}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 44 Action9 <- <{p.SetType("bool")}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 45 Action10 <- <{p.SetType("byte")}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 46 Action11 <- <{p.SetType("ubyte")}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 47 Action12 <- <{p.SetType("short")}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 48 Action13 <- <{p.SetType("ushort")}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 49 Action14 <- <{p.SetType("int")}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 50 Action15 <- <{p.SetType("uint")}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 51 Action16 <- <{p.SetType("float")}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 52 Action17 <- <{p.SetType("long")}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 53 Action18 <- <{p.SetType("ulong")}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 54 Action19 <- <{p.SetType("double")}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 55 Action20 <- <{p.SetType("int8")}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 56 Action21 <- <{p.SetType("int16")}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 57 Action22 <- <{p.SetType("uint16")}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 58 Action23 <- <{p.SetType("int32")}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 59 Action24 <- <{p.SetType("uint32")}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 60 Action25 <- <{p.SetType("int64")}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 61 Action26 <- <{p.SetType("uint64")}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 62 Action27 <- <{p.SetType("float32")}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 63 Action28 <- <{p.SetType("float64")}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 64 Action29 <- <{p.SetType("string")}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 65 Action30 <- <{p.SetType(text)}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 66 Action31 <- <{p.SetRepeated(text) }> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
	}
	p.rules = _rules
}
