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
			p.SetType("int8")
		case ruleAction11:
			p.SetType("int16")
		case ruleAction12:
			p.SetType("uint16")
		case ruleAction13:
			p.SetType("int32")
		case ruleAction14:
			p.SetType("uint32")
		case ruleAction15:
			p.SetType("int64")
		case ruleAction16:
			p.SetType("uint64")
		case ruleAction17:
			p.SetType("float32")
		case ruleAction18:
			p.SetType("float64")
		case ruleAction19:
			p.SetType("byte")
		case ruleAction20:
			p.SetType("ubyte")
		case ruleAction21:
			p.SetType("short")
		case ruleAction22:
			p.SetType("ushort")
		case ruleAction23:
			p.SetType("int")
		case ruleAction24:
			p.SetType("uint")
		case ruleAction25:
			p.SetType("float")
		case ruleAction26:
			p.SetType("long")
		case ruleAction27:
			p.SetType("ulong")
		case ruleAction28:
			p.SetType("double")
		case ruleAction29:
			p.SetRepeated("byte")
		case ruleAction30:
			p.SetType(text)
		case ruleAction31:
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
		/* 11 enum_fields <- <((spacing enum_field ',') / (spacing enum_field))> */
		func() bool {
			position84, tokenIndex84, depth84 := position, tokenIndex, depth
			{
				position85 := position
				depth++
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if !_rules[rulespacing]() {
						goto l87
					}
					if !_rules[ruleenum_field]() {
						goto l87
					}
					if buffer[position] != rune(',') {
						goto l87
					}
					position++
					goto l86
				l87:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
					if !_rules[rulespacing]() {
						goto l84
					}
					if !_rules[ruleenum_field]() {
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
		/* 12 union_decl <- <('u' 'n' 'i' 'o' 'n' spacing ident spacing metadata* '{' enum_fields+ '}' spacing Action6)> */
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
				if !_rules[ruleenum_fields]() {
					goto l88
				}
			l92:
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if !_rules[ruleenum_fields]() {
						goto l93
					}
					goto l92
				l93:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
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
			position94, tokenIndex94, depth94 := position, tokenIndex, depth
			{
				position95 := position
				depth++
				{
					position96, tokenIndex96, depth96 := position, tokenIndex, depth
					if !_rules[ruleenum_field_type]() {
						goto l97
					}
					if !_rules[rulespacing]() {
						goto l97
					}
					if !_rules[ruleAction7]() {
						goto l97
					}
					goto l96
				l97:
					position, tokenIndex, depth = position96, tokenIndex96, depth96
					if !_rules[ruleenum_field_type]() {
						goto l94
					}
					if !_rules[rulespacing]() {
						goto l94
					}
					if buffer[position] != rune('=') {
						goto l94
					}
					position++
					if !_rules[rulespacing]() {
						goto l94
					}
					if !_rules[ruleinteger_constant]() {
						goto l94
					}
					if !_rules[rulespacing]() {
						goto l94
					}
				}
			l96:
				depth--
				add(ruleenum_field, position95)
			}
			return true
		l94:
			position, tokenIndex, depth = position94, tokenIndex94, depth94
			return false
		},
		/* 14 enum_field_type <- <(ident Action8)> */
		func() bool {
			position98, tokenIndex98, depth98 := position, tokenIndex, depth
			{
				position99 := position
				depth++
				if !_rules[ruleident]() {
					goto l98
				}
				if !_rules[ruleAction8]() {
					goto l98
				}
				depth--
				add(ruleenum_field_type, position99)
			}
			return true
		l98:
			position, tokenIndex, depth = position98, tokenIndex98, depth98
			return false
		},
		/* 15 root_decl <- <('r' 'o' 'o' 't' '_' 't' 'y' 'p' 'e' spacing ident spacing ';' spacing)> */
		func() bool {
			position100, tokenIndex100, depth100 := position, tokenIndex, depth
			{
				position101 := position
				depth++
				if buffer[position] != rune('r') {
					goto l100
				}
				position++
				if buffer[position] != rune('o') {
					goto l100
				}
				position++
				if buffer[position] != rune('o') {
					goto l100
				}
				position++
				if buffer[position] != rune('t') {
					goto l100
				}
				position++
				if buffer[position] != rune('_') {
					goto l100
				}
				position++
				if buffer[position] != rune('t') {
					goto l100
				}
				position++
				if buffer[position] != rune('y') {
					goto l100
				}
				position++
				if buffer[position] != rune('p') {
					goto l100
				}
				position++
				if buffer[position] != rune('e') {
					goto l100
				}
				position++
				if !_rules[rulespacing]() {
					goto l100
				}
				if !_rules[ruleident]() {
					goto l100
				}
				if !_rules[rulespacing]() {
					goto l100
				}
				if buffer[position] != rune(';') {
					goto l100
				}
				position++
				if !_rules[rulespacing]() {
					goto l100
				}
				depth--
				add(ruleroot_decl, position101)
			}
			return true
		l100:
			position, tokenIndex, depth = position100, tokenIndex100, depth100
			return false
		},
		/* 16 file_extension_decl <- <('f' 'i' 'l' 'e' '_' 'e' 'x' 't' 'e' 'n' 's' 'i' 'o' 'n' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position102, tokenIndex102, depth102 := position, tokenIndex, depth
			{
				position103 := position
				depth++
				if buffer[position] != rune('f') {
					goto l102
				}
				position++
				if buffer[position] != rune('i') {
					goto l102
				}
				position++
				if buffer[position] != rune('l') {
					goto l102
				}
				position++
				if buffer[position] != rune('e') {
					goto l102
				}
				position++
				if buffer[position] != rune('_') {
					goto l102
				}
				position++
				if buffer[position] != rune('e') {
					goto l102
				}
				position++
				if buffer[position] != rune('x') {
					goto l102
				}
				position++
				if buffer[position] != rune('t') {
					goto l102
				}
				position++
				if buffer[position] != rune('e') {
					goto l102
				}
				position++
				if buffer[position] != rune('n') {
					goto l102
				}
				position++
				if buffer[position] != rune('s') {
					goto l102
				}
				position++
				if buffer[position] != rune('i') {
					goto l102
				}
				position++
				if buffer[position] != rune('o') {
					goto l102
				}
				position++
				if buffer[position] != rune('n') {
					goto l102
				}
				position++
				{
					position104 := position
					depth++
				l105:
					{
						position106, tokenIndex106, depth106 := position, tokenIndex, depth
						{
							position107, tokenIndex107, depth107 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l108
							}
							position++
							goto l107
						l108:
							position, tokenIndex, depth = position107, tokenIndex107, depth107
							if buffer[position] != rune('\t') {
								goto l106
							}
							position++
						}
					l107:
						goto l105
					l106:
						position, tokenIndex, depth = position106, tokenIndex106, depth106
					}
					depth--
					add(rulePegText, position104)
				}
				{
					position109 := position
					depth++
					{
						position112, tokenIndex112, depth112 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l112
						}
						position++
						goto l102
					l112:
						position, tokenIndex, depth = position112, tokenIndex112, depth112
					}
					if !matchDot() {
						goto l102
					}
				l110:
					{
						position111, tokenIndex111, depth111 := position, tokenIndex, depth
						{
							position113, tokenIndex113, depth113 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l113
							}
							position++
							goto l111
						l113:
							position, tokenIndex, depth = position113, tokenIndex113, depth113
						}
						if !matchDot() {
							goto l111
						}
						goto l110
					l111:
						position, tokenIndex, depth = position111, tokenIndex111, depth111
					}
					depth--
					add(rulePegText, position109)
				}
				if buffer[position] != rune(';') {
					goto l102
				}
				position++
				if !_rules[rulespacing]() {
					goto l102
				}
				depth--
				add(rulefile_extension_decl, position103)
			}
			return true
		l102:
			position, tokenIndex, depth = position102, tokenIndex102, depth102
			return false
		},
		/* 17 file_identifier_decl <- <('f' 'i' 'l' 'e' '_' 'i' 'd' 'e' 'n' 't' 'i' 'f' 'i' 'e' 'r' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position114, tokenIndex114, depth114 := position, tokenIndex, depth
			{
				position115 := position
				depth++
				if buffer[position] != rune('f') {
					goto l114
				}
				position++
				if buffer[position] != rune('i') {
					goto l114
				}
				position++
				if buffer[position] != rune('l') {
					goto l114
				}
				position++
				if buffer[position] != rune('e') {
					goto l114
				}
				position++
				if buffer[position] != rune('_') {
					goto l114
				}
				position++
				if buffer[position] != rune('i') {
					goto l114
				}
				position++
				if buffer[position] != rune('d') {
					goto l114
				}
				position++
				if buffer[position] != rune('e') {
					goto l114
				}
				position++
				if buffer[position] != rune('n') {
					goto l114
				}
				position++
				if buffer[position] != rune('t') {
					goto l114
				}
				position++
				if buffer[position] != rune('i') {
					goto l114
				}
				position++
				if buffer[position] != rune('f') {
					goto l114
				}
				position++
				if buffer[position] != rune('i') {
					goto l114
				}
				position++
				if buffer[position] != rune('e') {
					goto l114
				}
				position++
				if buffer[position] != rune('r') {
					goto l114
				}
				position++
				{
					position116 := position
					depth++
				l117:
					{
						position118, tokenIndex118, depth118 := position, tokenIndex, depth
						{
							position119, tokenIndex119, depth119 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l120
							}
							position++
							goto l119
						l120:
							position, tokenIndex, depth = position119, tokenIndex119, depth119
							if buffer[position] != rune('\t') {
								goto l118
							}
							position++
						}
					l119:
						goto l117
					l118:
						position, tokenIndex, depth = position118, tokenIndex118, depth118
					}
					depth--
					add(rulePegText, position116)
				}
				{
					position121 := position
					depth++
					{
						position124, tokenIndex124, depth124 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l124
						}
						position++
						goto l114
					l124:
						position, tokenIndex, depth = position124, tokenIndex124, depth124
					}
					if !matchDot() {
						goto l114
					}
				l122:
					{
						position123, tokenIndex123, depth123 := position, tokenIndex, depth
						{
							position125, tokenIndex125, depth125 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l125
							}
							position++
							goto l123
						l125:
							position, tokenIndex, depth = position125, tokenIndex125, depth125
						}
						if !matchDot() {
							goto l123
						}
						goto l122
					l123:
						position, tokenIndex, depth = position123, tokenIndex123, depth123
					}
					depth--
					add(rulePegText, position121)
				}
				if buffer[position] != rune(';') {
					goto l114
				}
				position++
				if !_rules[rulespacing]() {
					goto l114
				}
				depth--
				add(rulefile_identifier_decl, position115)
			}
			return true
		l114:
			position, tokenIndex, depth = position114, tokenIndex114, depth114
			return false
		},
		/* 18 attribute_decl <- <('a' 't' 't' 'r' 'i' 'b' 'u' 't' 'e' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position126, tokenIndex126, depth126 := position, tokenIndex, depth
			{
				position127 := position
				depth++
				if buffer[position] != rune('a') {
					goto l126
				}
				position++
				if buffer[position] != rune('t') {
					goto l126
				}
				position++
				if buffer[position] != rune('t') {
					goto l126
				}
				position++
				if buffer[position] != rune('r') {
					goto l126
				}
				position++
				if buffer[position] != rune('i') {
					goto l126
				}
				position++
				if buffer[position] != rune('b') {
					goto l126
				}
				position++
				if buffer[position] != rune('u') {
					goto l126
				}
				position++
				if buffer[position] != rune('t') {
					goto l126
				}
				position++
				if buffer[position] != rune('e') {
					goto l126
				}
				position++
				{
					position128 := position
					depth++
				l129:
					{
						position130, tokenIndex130, depth130 := position, tokenIndex, depth
						{
							position131, tokenIndex131, depth131 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l132
							}
							position++
							goto l131
						l132:
							position, tokenIndex, depth = position131, tokenIndex131, depth131
							if buffer[position] != rune('\t') {
								goto l130
							}
							position++
						}
					l131:
						goto l129
					l130:
						position, tokenIndex, depth = position130, tokenIndex130, depth130
					}
					depth--
					add(rulePegText, position128)
				}
				{
					position133 := position
					depth++
					{
						position136, tokenIndex136, depth136 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l136
						}
						position++
						goto l126
					l136:
						position, tokenIndex, depth = position136, tokenIndex136, depth136
					}
					if !matchDot() {
						goto l126
					}
				l134:
					{
						position135, tokenIndex135, depth135 := position, tokenIndex, depth
						{
							position137, tokenIndex137, depth137 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l137
							}
							position++
							goto l135
						l137:
							position, tokenIndex, depth = position137, tokenIndex137, depth137
						}
						if !matchDot() {
							goto l135
						}
						goto l134
					l135:
						position, tokenIndex, depth = position135, tokenIndex135, depth135
					}
					depth--
					add(rulePegText, position133)
				}
				if buffer[position] != rune(';') {
					goto l126
				}
				position++
				if !_rules[rulespacing]() {
					goto l126
				}
				depth--
				add(ruleattribute_decl, position127)
			}
			return true
		l126:
			position, tokenIndex, depth = position126, tokenIndex126, depth126
			return false
		},
		/* 19 rpc_decl <- <('r' 'p' 'c' '_' 's' 'e' 'r' 'v' 'i' 'c' 'e' <(' ' / '\t')*> ident '{' <(!'}' .)+> '}' spacing)> */
		func() bool {
			position138, tokenIndex138, depth138 := position, tokenIndex, depth
			{
				position139 := position
				depth++
				if buffer[position] != rune('r') {
					goto l138
				}
				position++
				if buffer[position] != rune('p') {
					goto l138
				}
				position++
				if buffer[position] != rune('c') {
					goto l138
				}
				position++
				if buffer[position] != rune('_') {
					goto l138
				}
				position++
				if buffer[position] != rune('s') {
					goto l138
				}
				position++
				if buffer[position] != rune('e') {
					goto l138
				}
				position++
				if buffer[position] != rune('r') {
					goto l138
				}
				position++
				if buffer[position] != rune('v') {
					goto l138
				}
				position++
				if buffer[position] != rune('i') {
					goto l138
				}
				position++
				if buffer[position] != rune('c') {
					goto l138
				}
				position++
				if buffer[position] != rune('e') {
					goto l138
				}
				position++
				{
					position140 := position
					depth++
				l141:
					{
						position142, tokenIndex142, depth142 := position, tokenIndex, depth
						{
							position143, tokenIndex143, depth143 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l144
							}
							position++
							goto l143
						l144:
							position, tokenIndex, depth = position143, tokenIndex143, depth143
							if buffer[position] != rune('\t') {
								goto l142
							}
							position++
						}
					l143:
						goto l141
					l142:
						position, tokenIndex, depth = position142, tokenIndex142, depth142
					}
					depth--
					add(rulePegText, position140)
				}
				if !_rules[ruleident]() {
					goto l138
				}
				if buffer[position] != rune('{') {
					goto l138
				}
				position++
				{
					position145 := position
					depth++
					{
						position148, tokenIndex148, depth148 := position, tokenIndex, depth
						if buffer[position] != rune('}') {
							goto l148
						}
						position++
						goto l138
					l148:
						position, tokenIndex, depth = position148, tokenIndex148, depth148
					}
					if !matchDot() {
						goto l138
					}
				l146:
					{
						position147, tokenIndex147, depth147 := position, tokenIndex, depth
						{
							position149, tokenIndex149, depth149 := position, tokenIndex, depth
							if buffer[position] != rune('}') {
								goto l149
							}
							position++
							goto l147
						l149:
							position, tokenIndex, depth = position149, tokenIndex149, depth149
						}
						if !matchDot() {
							goto l147
						}
						goto l146
					l147:
						position, tokenIndex, depth = position147, tokenIndex147, depth147
					}
					depth--
					add(rulePegText, position145)
				}
				if buffer[position] != rune('}') {
					goto l138
				}
				position++
				if !_rules[rulespacing]() {
					goto l138
				}
				depth--
				add(rulerpc_decl, position139)
			}
			return true
		l138:
			position, tokenIndex, depth = position138, tokenIndex138, depth138
			return false
		},
		/* 20 type <- <(('b' 'o' 'o' 'l' spacing Action9) / ('i' 'n' 't' '8' spacing Action10) / ('i' 'n' 't' '1' '6' spacing Action11) / ('u' 'i' 'n' 't' '1' '6' spacing Action12) / ('i' 'n' 't' '3' '2' spacing Action13) / ('u' 'i' 'n' 't' '3' '2' spacing Action14) / ('i' 'n' 't' '6' '4' spacing Action15) / ('u' 'i' 'n' 't' '6' '4' spacing Action16) / ('f' 'l' 'o' 'a' 't' '3' '2' spacing Action17) / ('f' 'l' 'o' 'a' 't' '6' '4' spacing Action18) / ('b' 'y' 't' 'e' spacing Action19) / ('u' 'b' 'y' 't' 'e' spacing Action20) / ('s' 'h' 'o' 'r' 't' spacing Action21) / ('u' 's' 'h' 'o' 'r' 't' spacing Action22) / ('i' 'n' 't' spacing Action23) / ('u' 'i' 'n' 't' spacing Action24) / ('f' 'l' 'o' 'a' 't' spacing Action25) / ('l' 'o' 'n' 'g' spacing Action26) / ('u' 'l' 'o' 'n' 'g' spacing Action27) / ('d' 'o' 'u' 'b' 'l' 'e' spacing Action28) / ('s' 't' 'r' 'i' 'n' 'g' spacing Action29) / (ident spacing Action30) / ('[' type ']' spacing Action31))> */
		func() bool {
			position150, tokenIndex150, depth150 := position, tokenIndex, depth
			{
				position151 := position
				depth++
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l153
					}
					position++
					if buffer[position] != rune('o') {
						goto l153
					}
					position++
					if buffer[position] != rune('o') {
						goto l153
					}
					position++
					if buffer[position] != rune('l') {
						goto l153
					}
					position++
					if !_rules[rulespacing]() {
						goto l153
					}
					if !_rules[ruleAction9]() {
						goto l153
					}
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('i') {
						goto l154
					}
					position++
					if buffer[position] != rune('n') {
						goto l154
					}
					position++
					if buffer[position] != rune('t') {
						goto l154
					}
					position++
					if buffer[position] != rune('8') {
						goto l154
					}
					position++
					if !_rules[rulespacing]() {
						goto l154
					}
					if !_rules[ruleAction10]() {
						goto l154
					}
					goto l152
				l154:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('i') {
						goto l155
					}
					position++
					if buffer[position] != rune('n') {
						goto l155
					}
					position++
					if buffer[position] != rune('t') {
						goto l155
					}
					position++
					if buffer[position] != rune('1') {
						goto l155
					}
					position++
					if buffer[position] != rune('6') {
						goto l155
					}
					position++
					if !_rules[rulespacing]() {
						goto l155
					}
					if !_rules[ruleAction11]() {
						goto l155
					}
					goto l152
				l155:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('u') {
						goto l156
					}
					position++
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
					if buffer[position] != rune('1') {
						goto l156
					}
					position++
					if buffer[position] != rune('6') {
						goto l156
					}
					position++
					if !_rules[rulespacing]() {
						goto l156
					}
					if !_rules[ruleAction12]() {
						goto l156
					}
					goto l152
				l156:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
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
					if buffer[position] != rune('3') {
						goto l157
					}
					position++
					if buffer[position] != rune('2') {
						goto l157
					}
					position++
					if !_rules[rulespacing]() {
						goto l157
					}
					if !_rules[ruleAction13]() {
						goto l157
					}
					goto l152
				l157:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('u') {
						goto l158
					}
					position++
					if buffer[position] != rune('i') {
						goto l158
					}
					position++
					if buffer[position] != rune('n') {
						goto l158
					}
					position++
					if buffer[position] != rune('t') {
						goto l158
					}
					position++
					if buffer[position] != rune('3') {
						goto l158
					}
					position++
					if buffer[position] != rune('2') {
						goto l158
					}
					position++
					if !_rules[rulespacing]() {
						goto l158
					}
					if !_rules[ruleAction14]() {
						goto l158
					}
					goto l152
				l158:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('i') {
						goto l159
					}
					position++
					if buffer[position] != rune('n') {
						goto l159
					}
					position++
					if buffer[position] != rune('t') {
						goto l159
					}
					position++
					if buffer[position] != rune('6') {
						goto l159
					}
					position++
					if buffer[position] != rune('4') {
						goto l159
					}
					position++
					if !_rules[rulespacing]() {
						goto l159
					}
					if !_rules[ruleAction15]() {
						goto l159
					}
					goto l152
				l159:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('u') {
						goto l160
					}
					position++
					if buffer[position] != rune('i') {
						goto l160
					}
					position++
					if buffer[position] != rune('n') {
						goto l160
					}
					position++
					if buffer[position] != rune('t') {
						goto l160
					}
					position++
					if buffer[position] != rune('6') {
						goto l160
					}
					position++
					if buffer[position] != rune('4') {
						goto l160
					}
					position++
					if !_rules[rulespacing]() {
						goto l160
					}
					if !_rules[ruleAction16]() {
						goto l160
					}
					goto l152
				l160:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('f') {
						goto l161
					}
					position++
					if buffer[position] != rune('l') {
						goto l161
					}
					position++
					if buffer[position] != rune('o') {
						goto l161
					}
					position++
					if buffer[position] != rune('a') {
						goto l161
					}
					position++
					if buffer[position] != rune('t') {
						goto l161
					}
					position++
					if buffer[position] != rune('3') {
						goto l161
					}
					position++
					if buffer[position] != rune('2') {
						goto l161
					}
					position++
					if !_rules[rulespacing]() {
						goto l161
					}
					if !_rules[ruleAction17]() {
						goto l161
					}
					goto l152
				l161:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('f') {
						goto l162
					}
					position++
					if buffer[position] != rune('l') {
						goto l162
					}
					position++
					if buffer[position] != rune('o') {
						goto l162
					}
					position++
					if buffer[position] != rune('a') {
						goto l162
					}
					position++
					if buffer[position] != rune('t') {
						goto l162
					}
					position++
					if buffer[position] != rune('6') {
						goto l162
					}
					position++
					if buffer[position] != rune('4') {
						goto l162
					}
					position++
					if !_rules[rulespacing]() {
						goto l162
					}
					if !_rules[ruleAction18]() {
						goto l162
					}
					goto l152
				l162:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('b') {
						goto l163
					}
					position++
					if buffer[position] != rune('y') {
						goto l163
					}
					position++
					if buffer[position] != rune('t') {
						goto l163
					}
					position++
					if buffer[position] != rune('e') {
						goto l163
					}
					position++
					if !_rules[rulespacing]() {
						goto l163
					}
					if !_rules[ruleAction19]() {
						goto l163
					}
					goto l152
				l163:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('u') {
						goto l164
					}
					position++
					if buffer[position] != rune('b') {
						goto l164
					}
					position++
					if buffer[position] != rune('y') {
						goto l164
					}
					position++
					if buffer[position] != rune('t') {
						goto l164
					}
					position++
					if buffer[position] != rune('e') {
						goto l164
					}
					position++
					if !_rules[rulespacing]() {
						goto l164
					}
					if !_rules[ruleAction20]() {
						goto l164
					}
					goto l152
				l164:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('s') {
						goto l165
					}
					position++
					if buffer[position] != rune('h') {
						goto l165
					}
					position++
					if buffer[position] != rune('o') {
						goto l165
					}
					position++
					if buffer[position] != rune('r') {
						goto l165
					}
					position++
					if buffer[position] != rune('t') {
						goto l165
					}
					position++
					if !_rules[rulespacing]() {
						goto l165
					}
					if !_rules[ruleAction21]() {
						goto l165
					}
					goto l152
				l165:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('u') {
						goto l166
					}
					position++
					if buffer[position] != rune('s') {
						goto l166
					}
					position++
					if buffer[position] != rune('h') {
						goto l166
					}
					position++
					if buffer[position] != rune('o') {
						goto l166
					}
					position++
					if buffer[position] != rune('r') {
						goto l166
					}
					position++
					if buffer[position] != rune('t') {
						goto l166
					}
					position++
					if !_rules[rulespacing]() {
						goto l166
					}
					if !_rules[ruleAction22]() {
						goto l166
					}
					goto l152
				l166:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
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
					if !_rules[rulespacing]() {
						goto l167
					}
					if !_rules[ruleAction23]() {
						goto l167
					}
					goto l152
				l167:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
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
					if !_rules[rulespacing]() {
						goto l168
					}
					if !_rules[ruleAction24]() {
						goto l168
					}
					goto l152
				l168:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
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
					if !_rules[rulespacing]() {
						goto l169
					}
					if !_rules[ruleAction25]() {
						goto l169
					}
					goto l152
				l169:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('l') {
						goto l170
					}
					position++
					if buffer[position] != rune('o') {
						goto l170
					}
					position++
					if buffer[position] != rune('n') {
						goto l170
					}
					position++
					if buffer[position] != rune('g') {
						goto l170
					}
					position++
					if !_rules[rulespacing]() {
						goto l170
					}
					if !_rules[ruleAction26]() {
						goto l170
					}
					goto l152
				l170:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('u') {
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
					if !_rules[ruleAction27]() {
						goto l171
					}
					goto l152
				l171:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('d') {
						goto l172
					}
					position++
					if buffer[position] != rune('o') {
						goto l172
					}
					position++
					if buffer[position] != rune('u') {
						goto l172
					}
					position++
					if buffer[position] != rune('b') {
						goto l172
					}
					position++
					if buffer[position] != rune('l') {
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
					if !_rules[ruleAction28]() {
						goto l172
					}
					goto l152
				l172:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('s') {
						goto l173
					}
					position++
					if buffer[position] != rune('t') {
						goto l173
					}
					position++
					if buffer[position] != rune('r') {
						goto l173
					}
					position++
					if buffer[position] != rune('i') {
						goto l173
					}
					position++
					if buffer[position] != rune('n') {
						goto l173
					}
					position++
					if buffer[position] != rune('g') {
						goto l173
					}
					position++
					if !_rules[rulespacing]() {
						goto l173
					}
					if !_rules[ruleAction29]() {
						goto l173
					}
					goto l152
				l173:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if !_rules[ruleident]() {
						goto l174
					}
					if !_rules[rulespacing]() {
						goto l174
					}
					if !_rules[ruleAction30]() {
						goto l174
					}
					goto l152
				l174:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('[') {
						goto l150
					}
					position++
					if !_rules[ruletype]() {
						goto l150
					}
					if buffer[position] != rune(']') {
						goto l150
					}
					position++
					if !_rules[rulespacing]() {
						goto l150
					}
					if !_rules[ruleAction31]() {
						goto l150
					}
				}
			l152:
				depth--
				add(ruletype, position151)
			}
			return true
		l150:
			position, tokenIndex, depth = position150, tokenIndex150, depth150
			return false
		},
		/* 21 scalar <- <(integer_constant / float_constant)> */
		func() bool {
			position175, tokenIndex175, depth175 := position, tokenIndex, depth
			{
				position176 := position
				depth++
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if !_rules[ruleinteger_constant]() {
						goto l178
					}
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if !_rules[rulefloat_constant]() {
						goto l175
					}
				}
			l177:
				depth--
				add(rulescalar, position176)
			}
			return true
		l175:
			position, tokenIndex, depth = position175, tokenIndex175, depth175
			return false
		},
		/* 22 integer_constant <- <(<[0-9]+> / ('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		func() bool {
			position179, tokenIndex179, depth179 := position, tokenIndex, depth
			{
				position180 := position
				depth++
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					{
						position183 := position
						depth++
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l182
						}
						position++
					l184:
						{
							position185, tokenIndex185, depth185 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l185
							}
							position++
							goto l184
						l185:
							position, tokenIndex, depth = position185, tokenIndex185, depth185
						}
						depth--
						add(rulePegText, position183)
					}
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('t') {
						goto l186
					}
					position++
					if buffer[position] != rune('r') {
						goto l186
					}
					position++
					if buffer[position] != rune('u') {
						goto l186
					}
					position++
					if buffer[position] != rune('e') {
						goto l186
					}
					position++
					goto l181
				l186:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('f') {
						goto l179
					}
					position++
					if buffer[position] != rune('a') {
						goto l179
					}
					position++
					if buffer[position] != rune('l') {
						goto l179
					}
					position++
					if buffer[position] != rune('s') {
						goto l179
					}
					position++
					if buffer[position] != rune('e') {
						goto l179
					}
					position++
				}
			l181:
				depth--
				add(ruleinteger_constant, position180)
			}
			return true
		l179:
			position, tokenIndex, depth = position179, tokenIndex179, depth179
			return false
		},
		/* 23 float_constant <- <(<('-'* [0-9]+ . [0-9])> / float_constant_exp)> */
		func() bool {
			position187, tokenIndex187, depth187 := position, tokenIndex, depth
			{
				position188 := position
				depth++
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					{
						position191 := position
						depth++
					l192:
						{
							position193, tokenIndex193, depth193 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l193
							}
							position++
							goto l192
						l193:
							position, tokenIndex, depth = position193, tokenIndex193, depth193
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l190
						}
						position++
					l194:
						{
							position195, tokenIndex195, depth195 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l195
							}
							position++
							goto l194
						l195:
							position, tokenIndex, depth = position195, tokenIndex195, depth195
						}
						if !matchDot() {
							goto l190
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l190
						}
						position++
						depth--
						add(rulePegText, position191)
					}
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if !_rules[rulefloat_constant_exp]() {
						goto l187
					}
				}
			l189:
				depth--
				add(rulefloat_constant, position188)
			}
			return true
		l187:
			position, tokenIndex, depth = position187, tokenIndex187, depth187
			return false
		},
		/* 24 float_constant_exp <- <(<('-'* [0-9]+ . [0-9]+)> <('e' / 'E')> <([+-]] / '>' / ' ' / '<' / '[' / [0-9])+>)> */
		func() bool {
			position196, tokenIndex196, depth196 := position, tokenIndex, depth
			{
				position197 := position
				depth++
				{
					position198 := position
					depth++
				l199:
					{
						position200, tokenIndex200, depth200 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l200
						}
						position++
						goto l199
					l200:
						position, tokenIndex, depth = position200, tokenIndex200, depth200
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l196
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
					if !matchDot() {
						goto l196
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l196
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
					depth--
					add(rulePegText, position198)
				}
				{
					position205 := position
					depth++
					{
						position206, tokenIndex206, depth206 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l207
						}
						position++
						goto l206
					l207:
						position, tokenIndex, depth = position206, tokenIndex206, depth206
						if buffer[position] != rune('E') {
							goto l196
						}
						position++
					}
				l206:
					depth--
					add(rulePegText, position205)
				}
				{
					position208 := position
					depth++
					{
						position211, tokenIndex211, depth211 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('+') || c > rune(']') {
							goto l212
						}
						position++
						goto l211
					l212:
						position, tokenIndex, depth = position211, tokenIndex211, depth211
						if buffer[position] != rune('>') {
							goto l213
						}
						position++
						goto l211
					l213:
						position, tokenIndex, depth = position211, tokenIndex211, depth211
						if buffer[position] != rune(' ') {
							goto l214
						}
						position++
						goto l211
					l214:
						position, tokenIndex, depth = position211, tokenIndex211, depth211
						if buffer[position] != rune('<') {
							goto l215
						}
						position++
						goto l211
					l215:
						position, tokenIndex, depth = position211, tokenIndex211, depth211
						if buffer[position] != rune('[') {
							goto l216
						}
						position++
						goto l211
					l216:
						position, tokenIndex, depth = position211, tokenIndex211, depth211
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l196
						}
						position++
					}
				l211:
				l209:
					{
						position210, tokenIndex210, depth210 := position, tokenIndex, depth
						{
							position217, tokenIndex217, depth217 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('+') || c > rune(']') {
								goto l218
							}
							position++
							goto l217
						l218:
							position, tokenIndex, depth = position217, tokenIndex217, depth217
							if buffer[position] != rune('>') {
								goto l219
							}
							position++
							goto l217
						l219:
							position, tokenIndex, depth = position217, tokenIndex217, depth217
							if buffer[position] != rune(' ') {
								goto l220
							}
							position++
							goto l217
						l220:
							position, tokenIndex, depth = position217, tokenIndex217, depth217
							if buffer[position] != rune('<') {
								goto l221
							}
							position++
							goto l217
						l221:
							position, tokenIndex, depth = position217, tokenIndex217, depth217
							if buffer[position] != rune('[') {
								goto l222
							}
							position++
							goto l217
						l222:
							position, tokenIndex, depth = position217, tokenIndex217, depth217
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l210
							}
							position++
						}
					l217:
						goto l209
					l210:
						position, tokenIndex, depth = position210, tokenIndex210, depth210
					}
					depth--
					add(rulePegText, position208)
				}
				depth--
				add(rulefloat_constant_exp, position197)
			}
			return true
		l196:
			position, tokenIndex, depth = position196, tokenIndex196, depth196
			return false
		},
		/* 25 ident <- <<(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position223, tokenIndex223, depth223 := position, tokenIndex, depth
			{
				position224 := position
				depth++
				{
					position225 := position
					depth++
					{
						position226, tokenIndex226, depth226 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l227
						}
						position++
						goto l226
					l227:
						position, tokenIndex, depth = position226, tokenIndex226, depth226
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l228
						}
						position++
						goto l226
					l228:
						position, tokenIndex, depth = position226, tokenIndex226, depth226
						if buffer[position] != rune('_') {
							goto l223
						}
						position++
					}
				l226:
				l229:
					{
						position230, tokenIndex230, depth230 := position, tokenIndex, depth
						{
							position231, tokenIndex231, depth231 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l232
							}
							position++
							goto l231
						l232:
							position, tokenIndex, depth = position231, tokenIndex231, depth231
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l233
							}
							position++
							goto l231
						l233:
							position, tokenIndex, depth = position231, tokenIndex231, depth231
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l234
							}
							position++
							goto l231
						l234:
							position, tokenIndex, depth = position231, tokenIndex231, depth231
							if buffer[position] != rune('_') {
								goto l230
							}
							position++
						}
					l231:
						goto l229
					l230:
						position, tokenIndex, depth = position230, tokenIndex230, depth230
					}
					depth--
					add(rulePegText, position225)
				}
				depth--
				add(ruleident, position224)
			}
			return true
		l223:
			position, tokenIndex, depth = position223, tokenIndex223, depth223
			return false
		},
		/* 26 only_comment <- <(spacing ';')> */
		func() bool {
			position235, tokenIndex235, depth235 := position, tokenIndex, depth
			{
				position236 := position
				depth++
				if !_rules[rulespacing]() {
					goto l235
				}
				if buffer[position] != rune(';') {
					goto l235
				}
				position++
				depth--
				add(ruleonly_comment, position236)
			}
			return true
		l235:
			position, tokenIndex, depth = position235, tokenIndex235, depth235
			return false
		},
		/* 27 spacing <- <space_comment*> */
		func() bool {
			{
				position238 := position
				depth++
			l239:
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if !_rules[rulespace_comment]() {
						goto l240
					}
					goto l239
				l240:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
				}
				depth--
				add(rulespacing, position238)
			}
			return true
		},
		/* 28 space_comment <- <(space / comment)> */
		func() bool {
			position241, tokenIndex241, depth241 := position, tokenIndex, depth
			{
				position242 := position
				depth++
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if !_rules[rulespace]() {
						goto l244
					}
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if !_rules[rulecomment]() {
						goto l241
					}
				}
			l243:
				depth--
				add(rulespace_comment, position242)
			}
			return true
		l241:
			position, tokenIndex, depth = position241, tokenIndex241, depth241
			return false
		},
		/* 29 comment <- <('/' '/' (!end_of_line .)* end_of_line)> */
		func() bool {
			position245, tokenIndex245, depth245 := position, tokenIndex, depth
			{
				position246 := position
				depth++
				if buffer[position] != rune('/') {
					goto l245
				}
				position++
				if buffer[position] != rune('/') {
					goto l245
				}
				position++
			l247:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					{
						position249, tokenIndex249, depth249 := position, tokenIndex, depth
						if !_rules[ruleend_of_line]() {
							goto l249
						}
						goto l248
					l249:
						position, tokenIndex, depth = position249, tokenIndex249, depth249
					}
					if !matchDot() {
						goto l248
					}
					goto l247
				l248:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
				}
				if !_rules[ruleend_of_line]() {
					goto l245
				}
				depth--
				add(rulecomment, position246)
			}
			return true
		l245:
			position, tokenIndex, depth = position245, tokenIndex245, depth245
			return false
		},
		/* 30 space <- <(' ' / '\t' / end_of_line)> */
		func() bool {
			position250, tokenIndex250, depth250 := position, tokenIndex, depth
			{
				position251 := position
				depth++
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l253
					}
					position++
					goto l252
				l253:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if buffer[position] != rune('\t') {
						goto l254
					}
					position++
					goto l252
				l254:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if !_rules[ruleend_of_line]() {
						goto l250
					}
				}
			l252:
				depth--
				add(rulespace, position251)
			}
			return true
		l250:
			position, tokenIndex, depth = position250, tokenIndex250, depth250
			return false
		},
		/* 31 end_of_line <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position255, tokenIndex255, depth255 := position, tokenIndex, depth
			{
				position256 := position
				depth++
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l258
					}
					position++
					if buffer[position] != rune('\n') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('\n') {
						goto l259
					}
					position++
					goto l257
				l259:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('\r') {
						goto l255
					}
					position++
				}
			l257:
				depth--
				add(ruleend_of_line, position256)
			}
			return true
		l255:
			position, tokenIndex, depth = position255, tokenIndex255, depth255
			return false
		},
		/* 32 end_of_file <- <!.> */
		func() bool {
			position260, tokenIndex260, depth260 := position, tokenIndex, depth
			{
				position261 := position
				depth++
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if !matchDot() {
						goto l262
					}
					goto l260
				l262:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
				}
				depth--
				add(ruleend_of_file, position261)
			}
			return true
		l260:
			position, tokenIndex, depth = position260, tokenIndex260, depth260
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
		/* 45 Action10 <- <{p.SetType("int8")}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 46 Action11 <- <{p.SetType("int16")}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 47 Action12 <- <{p.SetType("uint16")}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 48 Action13 <- <{p.SetType("int32")}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 49 Action14 <- <{p.SetType("uint32")}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 50 Action15 <- <{p.SetType("int64")}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 51 Action16 <- <{p.SetType("uint64")}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 52 Action17 <- <{p.SetType("float32")}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 53 Action18 <- <{p.SetType("float64")}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 54 Action19 <- <{p.SetType("byte")}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 55 Action20 <- <{p.SetType("ubyte")}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 56 Action21 <- <{p.SetType("short")}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 57 Action22 <- <{p.SetType("ushort")}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 58 Action23 <- <{p.SetType("int")}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 59 Action24 <- <{p.SetType("uint")}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 60 Action25 <- <{p.SetType("float")}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 61 Action26 <- <{p.SetType("long")}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 62 Action27 <- <{p.SetType("ulong")}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 63 Action28 <- <{p.SetType("double")}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 64 Action29 <- <{p.SetRepeated("byte")}> */
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
		/* 66 Action31 <- <{p.SetRepeated("") }> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
	}
	p.rules = _rules
}
