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
	ruleenum_field
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
	"enum_field",
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
	rules  [62]func() bool
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
			p.SetType("bool")
		case ruleAction7:
			p.SetType("byte")
		case ruleAction8:
			p.SetType("ubyte")
		case ruleAction9:
			p.SetType("short")
		case ruleAction10:
			p.SetType("ushort")
		case ruleAction11:
			p.SetType("int")
		case ruleAction12:
			p.SetType("uint")
		case ruleAction13:
			p.SetType("float")
		case ruleAction14:
			p.SetType("long")
		case ruleAction15:
			p.SetType("ulong")
		case ruleAction16:
			p.SetType("double")
		case ruleAction17:
			p.SetType("int8")
		case ruleAction18:
			p.SetType("int16")
		case ruleAction19:
			p.SetType("uint16")
		case ruleAction20:
			p.SetType("int32")
		case ruleAction21:
			p.SetType("uint32")
		case ruleAction22:
			p.SetType("int64")
		case ruleAction23:
			p.SetType("uint64")
		case ruleAction24:
			p.SetType("float32")
		case ruleAction25:
			p.SetType("float64")
		case ruleAction26:
			p.SetType("string")
		case ruleAction27:
			p.SetType(text)
		case ruleAction28:
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
		/* 1 statment_decl <- <(namespace_decl / type_decl / enum_decl / root_decl / file_extension_decl / file_identifier_decl / attribute_decl / rpc_decl / only_comment)> */
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
					if !_rules[ruletype_decl]() {
						goto l13
					}
					goto l11
				l13:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[ruleenum_decl]() {
						goto l14
					}
					goto l11
				l14:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[ruleroot_decl]() {
						goto l15
					}
					goto l11
				l15:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulefile_extension_decl]() {
						goto l16
					}
					goto l11
				l16:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulefile_identifier_decl]() {
						goto l17
					}
					goto l11
				l17:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[ruleattribute_decl]() {
						goto l18
					}
					goto l11
				l18:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulerpc_decl]() {
						goto l19
					}
					goto l11
				l19:
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
			position20, tokenIndex20, depth20 := position, tokenIndex, depth
			{
				position21 := position
				depth++
				if buffer[position] != rune('n') {
					goto l20
				}
				position++
				if buffer[position] != rune('a') {
					goto l20
				}
				position++
				if buffer[position] != rune('m') {
					goto l20
				}
				position++
				if buffer[position] != rune('e') {
					goto l20
				}
				position++
				if buffer[position] != rune('s') {
					goto l20
				}
				position++
				if buffer[position] != rune('p') {
					goto l20
				}
				position++
				if buffer[position] != rune('a') {
					goto l20
				}
				position++
				if buffer[position] != rune('c') {
					goto l20
				}
				position++
				if buffer[position] != rune('e') {
					goto l20
				}
				position++
				if !_rules[rulespacing]() {
					goto l20
				}
				{
					position22 := position
					depth++
					{
						position25, tokenIndex25, depth25 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('A') || c > rune('z') {
							goto l26
						}
						position++
						goto l25
					l26:
						position, tokenIndex, depth = position25, tokenIndex25, depth25
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l27
						}
						position++
						goto l25
					l27:
						position, tokenIndex, depth = position25, tokenIndex25, depth25
						if buffer[position] != rune('_') {
							goto l28
						}
						position++
						goto l25
					l28:
						position, tokenIndex, depth = position25, tokenIndex25, depth25
						if buffer[position] != rune('.') {
							goto l29
						}
						position++
						goto l25
					l29:
						position, tokenIndex, depth = position25, tokenIndex25, depth25
						if buffer[position] != rune('-') {
							goto l20
						}
						position++
					}
				l25:
				l23:
					{
						position24, tokenIndex24, depth24 := position, tokenIndex, depth
						{
							position30, tokenIndex30, depth30 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('A') || c > rune('z') {
								goto l31
							}
							position++
							goto l30
						l31:
							position, tokenIndex, depth = position30, tokenIndex30, depth30
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l32
							}
							position++
							goto l30
						l32:
							position, tokenIndex, depth = position30, tokenIndex30, depth30
							if buffer[position] != rune('_') {
								goto l33
							}
							position++
							goto l30
						l33:
							position, tokenIndex, depth = position30, tokenIndex30, depth30
							if buffer[position] != rune('.') {
								goto l34
							}
							position++
							goto l30
						l34:
							position, tokenIndex, depth = position30, tokenIndex30, depth30
							if buffer[position] != rune('-') {
								goto l24
							}
							position++
						}
					l30:
						goto l23
					l24:
						position, tokenIndex, depth = position24, tokenIndex24, depth24
					}
					depth--
					add(rulePegText, position22)
				}
				if !_rules[ruleAction0]() {
					goto l20
				}
				if buffer[position] != rune(';') {
					goto l20
				}
				position++
				if !_rules[rulespacing]() {
					goto l20
				}
				depth--
				add(rulenamespace_decl, position21)
			}
			return true
		l20:
			position, tokenIndex, depth = position20, tokenIndex20, depth20
			return false
		},
		/* 3 include <- <('i' 'n' 'c' 'l' 'u' 'd' 'e' spacing ident comment ';' spacing)> */
		func() bool {
			position35, tokenIndex35, depth35 := position, tokenIndex, depth
			{
				position36 := position
				depth++
				if buffer[position] != rune('i') {
					goto l35
				}
				position++
				if buffer[position] != rune('n') {
					goto l35
				}
				position++
				if buffer[position] != rune('c') {
					goto l35
				}
				position++
				if buffer[position] != rune('l') {
					goto l35
				}
				position++
				if buffer[position] != rune('u') {
					goto l35
				}
				position++
				if buffer[position] != rune('d') {
					goto l35
				}
				position++
				if buffer[position] != rune('e') {
					goto l35
				}
				position++
				if !_rules[rulespacing]() {
					goto l35
				}
				if !_rules[ruleident]() {
					goto l35
				}
				if !_rules[rulecomment]() {
					goto l35
				}
				if buffer[position] != rune(';') {
					goto l35
				}
				position++
				if !_rules[rulespacing]() {
					goto l35
				}
				depth--
				add(ruleinclude, position36)
			}
			return true
		l35:
			position, tokenIndex, depth = position35, tokenIndex35, depth35
			return false
		},
		/* 4 type_decl <- <(type_label spacing typename spacing metadata* '{' field_decl+ '}' spacing Action1)> */
		func() bool {
			position37, tokenIndex37, depth37 := position, tokenIndex, depth
			{
				position38 := position
				depth++
				if !_rules[ruletype_label]() {
					goto l37
				}
				if !_rules[rulespacing]() {
					goto l37
				}
				if !_rules[ruletypename]() {
					goto l37
				}
				if !_rules[rulespacing]() {
					goto l37
				}
			l39:
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if !_rules[rulemetadata]() {
						goto l40
					}
					goto l39
				l40:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
				}
				if buffer[position] != rune('{') {
					goto l37
				}
				position++
				if !_rules[rulefield_decl]() {
					goto l37
				}
			l41:
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if !_rules[rulefield_decl]() {
						goto l42
					}
					goto l41
				l42:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
				}
				if buffer[position] != rune('}') {
					goto l37
				}
				position++
				if !_rules[rulespacing]() {
					goto l37
				}
				if !_rules[ruleAction1]() {
					goto l37
				}
				depth--
				add(ruletype_decl, position38)
			}
			return true
		l37:
			position, tokenIndex, depth = position37, tokenIndex37, depth37
			return false
		},
		/* 5 type_label <- <(('t' 'a' 'b' 'l' 'e') / ('s' 't' 'r' 'u' 'c' 't'))> */
		func() bool {
			position43, tokenIndex43, depth43 := position, tokenIndex, depth
			{
				position44 := position
				depth++
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l46
					}
					position++
					if buffer[position] != rune('a') {
						goto l46
					}
					position++
					if buffer[position] != rune('b') {
						goto l46
					}
					position++
					if buffer[position] != rune('l') {
						goto l46
					}
					position++
					if buffer[position] != rune('e') {
						goto l46
					}
					position++
					goto l45
				l46:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if buffer[position] != rune('s') {
						goto l43
					}
					position++
					if buffer[position] != rune('t') {
						goto l43
					}
					position++
					if buffer[position] != rune('r') {
						goto l43
					}
					position++
					if buffer[position] != rune('u') {
						goto l43
					}
					position++
					if buffer[position] != rune('c') {
						goto l43
					}
					position++
					if buffer[position] != rune('t') {
						goto l43
					}
					position++
				}
			l45:
				depth--
				add(ruletype_label, position44)
			}
			return true
		l43:
			position, tokenIndex, depth = position43, tokenIndex43, depth43
			return false
		},
		/* 6 typename <- <(ident Action2)> */
		func() bool {
			position47, tokenIndex47, depth47 := position, tokenIndex, depth
			{
				position48 := position
				depth++
				if !_rules[ruleident]() {
					goto l47
				}
				if !_rules[ruleAction2]() {
					goto l47
				}
				depth--
				add(ruletypename, position48)
			}
			return true
		l47:
			position, tokenIndex, depth = position47, tokenIndex47, depth47
			return false
		},
		/* 7 metadata <- <('(' <(!')' .)*> ')')> */
		func() bool {
			position49, tokenIndex49, depth49 := position, tokenIndex, depth
			{
				position50 := position
				depth++
				if buffer[position] != rune('(') {
					goto l49
				}
				position++
				{
					position51 := position
					depth++
				l52:
					{
						position53, tokenIndex53, depth53 := position, tokenIndex, depth
						{
							position54, tokenIndex54, depth54 := position, tokenIndex, depth
							if buffer[position] != rune(')') {
								goto l54
							}
							position++
							goto l53
						l54:
							position, tokenIndex, depth = position54, tokenIndex54, depth54
						}
						if !matchDot() {
							goto l53
						}
						goto l52
					l53:
						position, tokenIndex, depth = position53, tokenIndex53, depth53
					}
					depth--
					add(rulePegText, position51)
				}
				if buffer[position] != rune(')') {
					goto l49
				}
				position++
				depth--
				add(rulemetadata, position50)
			}
			return true
		l49:
			position, tokenIndex, depth = position49, tokenIndex49, depth49
			return false
		},
		/* 8 field_decl <- <((spacing field_type ':' type metadata* ';' spacing Action3) / (spacing field_type ':' type <(' ' / '\t')*> '=' <(' ' / '\t')*> scalar metadata* ';' spacing Action4))> */
		func() bool {
			position55, tokenIndex55, depth55 := position, tokenIndex, depth
			{
				position56 := position
				depth++
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					if !_rules[rulespacing]() {
						goto l58
					}
					if !_rules[rulefield_type]() {
						goto l58
					}
					if buffer[position] != rune(':') {
						goto l58
					}
					position++
					if !_rules[ruletype]() {
						goto l58
					}
				l59:
					{
						position60, tokenIndex60, depth60 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l60
						}
						goto l59
					l60:
						position, tokenIndex, depth = position60, tokenIndex60, depth60
					}
					if buffer[position] != rune(';') {
						goto l58
					}
					position++
					if !_rules[rulespacing]() {
						goto l58
					}
					if !_rules[ruleAction3]() {
						goto l58
					}
					goto l57
				l58:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
					if !_rules[rulespacing]() {
						goto l55
					}
					if !_rules[rulefield_type]() {
						goto l55
					}
					if buffer[position] != rune(':') {
						goto l55
					}
					position++
					if !_rules[ruletype]() {
						goto l55
					}
					{
						position61 := position
						depth++
					l62:
						{
							position63, tokenIndex63, depth63 := position, tokenIndex, depth
							{
								position64, tokenIndex64, depth64 := position, tokenIndex, depth
								if buffer[position] != rune(' ') {
									goto l65
								}
								position++
								goto l64
							l65:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('\t') {
									goto l63
								}
								position++
							}
						l64:
							goto l62
						l63:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
						}
						depth--
						add(rulePegText, position61)
					}
					if buffer[position] != rune('=') {
						goto l55
					}
					position++
					{
						position66 := position
						depth++
					l67:
						{
							position68, tokenIndex68, depth68 := position, tokenIndex, depth
							{
								position69, tokenIndex69, depth69 := position, tokenIndex, depth
								if buffer[position] != rune(' ') {
									goto l70
								}
								position++
								goto l69
							l70:
								position, tokenIndex, depth = position69, tokenIndex69, depth69
								if buffer[position] != rune('\t') {
									goto l68
								}
								position++
							}
						l69:
							goto l67
						l68:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
						}
						depth--
						add(rulePegText, position66)
					}
					if !_rules[rulescalar]() {
						goto l55
					}
				l71:
					{
						position72, tokenIndex72, depth72 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l72
						}
						goto l71
					l72:
						position, tokenIndex, depth = position72, tokenIndex72, depth72
					}
					if buffer[position] != rune(';') {
						goto l55
					}
					position++
					if !_rules[rulespacing]() {
						goto l55
					}
					if !_rules[ruleAction4]() {
						goto l55
					}
				}
			l57:
				depth--
				add(rulefield_decl, position56)
			}
			return true
		l55:
			position, tokenIndex, depth = position55, tokenIndex55, depth55
			return false
		},
		/* 9 field_type <- <(ident Action5)> */
		func() bool {
			position73, tokenIndex73, depth73 := position, tokenIndex, depth
			{
				position74 := position
				depth++
				if !_rules[ruleident]() {
					goto l73
				}
				if !_rules[ruleAction5]() {
					goto l73
				}
				depth--
				add(rulefield_type, position74)
			}
			return true
		l73:
			position, tokenIndex, depth = position73, tokenIndex73, depth73
			return false
		},
		/* 10 enum_decl <- <(('e' 'n' 'u' 'm' spacing ident spacing metadata* '{' enum_fields '}' spacing) / ('e' 'n' 'u' 'm' spacing ident ':' type spacing metadata* '{' enum_fields '}' spacing) / ('u' 'n' 'i' 'o' 'n' spacing ident metadata* '{' enum_fields '}' spacing))> */
		func() bool {
			position75, tokenIndex75, depth75 := position, tokenIndex, depth
			{
				position76 := position
				depth++
				{
					position77, tokenIndex77, depth77 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l78
					}
					position++
					if buffer[position] != rune('n') {
						goto l78
					}
					position++
					if buffer[position] != rune('u') {
						goto l78
					}
					position++
					if buffer[position] != rune('m') {
						goto l78
					}
					position++
					if !_rules[rulespacing]() {
						goto l78
					}
					if !_rules[ruleident]() {
						goto l78
					}
					if !_rules[rulespacing]() {
						goto l78
					}
				l79:
					{
						position80, tokenIndex80, depth80 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l80
						}
						goto l79
					l80:
						position, tokenIndex, depth = position80, tokenIndex80, depth80
					}
					if buffer[position] != rune('{') {
						goto l78
					}
					position++
					if !_rules[ruleenum_fields]() {
						goto l78
					}
					if buffer[position] != rune('}') {
						goto l78
					}
					position++
					if !_rules[rulespacing]() {
						goto l78
					}
					goto l77
				l78:
					position, tokenIndex, depth = position77, tokenIndex77, depth77
					if buffer[position] != rune('e') {
						goto l81
					}
					position++
					if buffer[position] != rune('n') {
						goto l81
					}
					position++
					if buffer[position] != rune('u') {
						goto l81
					}
					position++
					if buffer[position] != rune('m') {
						goto l81
					}
					position++
					if !_rules[rulespacing]() {
						goto l81
					}
					if !_rules[ruleident]() {
						goto l81
					}
					if buffer[position] != rune(':') {
						goto l81
					}
					position++
					if !_rules[ruletype]() {
						goto l81
					}
					if !_rules[rulespacing]() {
						goto l81
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
						goto l81
					}
					position++
					if !_rules[ruleenum_fields]() {
						goto l81
					}
					if buffer[position] != rune('}') {
						goto l81
					}
					position++
					if !_rules[rulespacing]() {
						goto l81
					}
					goto l77
				l81:
					position, tokenIndex, depth = position77, tokenIndex77, depth77
					if buffer[position] != rune('u') {
						goto l75
					}
					position++
					if buffer[position] != rune('n') {
						goto l75
					}
					position++
					if buffer[position] != rune('i') {
						goto l75
					}
					position++
					if buffer[position] != rune('o') {
						goto l75
					}
					position++
					if buffer[position] != rune('n') {
						goto l75
					}
					position++
					if !_rules[rulespacing]() {
						goto l75
					}
					if !_rules[ruleident]() {
						goto l75
					}
				l84:
					{
						position85, tokenIndex85, depth85 := position, tokenIndex, depth
						if !_rules[rulemetadata]() {
							goto l85
						}
						goto l84
					l85:
						position, tokenIndex, depth = position85, tokenIndex85, depth85
					}
					if buffer[position] != rune('{') {
						goto l75
					}
					position++
					if !_rules[ruleenum_fields]() {
						goto l75
					}
					if buffer[position] != rune('}') {
						goto l75
					}
					position++
					if !_rules[rulespacing]() {
						goto l75
					}
				}
			l77:
				depth--
				add(ruleenum_decl, position76)
			}
			return true
		l75:
			position, tokenIndex, depth = position75, tokenIndex75, depth75
			return false
		},
		/* 11 enum_fields <- <(enum_field / (enum_field ',' enum_fields))> */
		func() bool {
			position86, tokenIndex86, depth86 := position, tokenIndex, depth
			{
				position87 := position
				depth++
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if !_rules[ruleenum_field]() {
						goto l89
					}
					goto l88
				l89:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if !_rules[ruleenum_field]() {
						goto l86
					}
					if buffer[position] != rune(',') {
						goto l86
					}
					position++
					if !_rules[ruleenum_fields]() {
						goto l86
					}
				}
			l88:
				depth--
				add(ruleenum_fields, position87)
			}
			return true
		l86:
			position, tokenIndex, depth = position86, tokenIndex86, depth86
			return false
		},
		/* 12 enum_field <- <(ident / (ident spacing '=' spacing integer_constant))> */
		func() bool {
			position90, tokenIndex90, depth90 := position, tokenIndex, depth
			{
				position91 := position
				depth++
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if !_rules[ruleident]() {
						goto l93
					}
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if !_rules[ruleident]() {
						goto l90
					}
					if !_rules[rulespacing]() {
						goto l90
					}
					if buffer[position] != rune('=') {
						goto l90
					}
					position++
					if !_rules[rulespacing]() {
						goto l90
					}
					if !_rules[ruleinteger_constant]() {
						goto l90
					}
				}
			l92:
				depth--
				add(ruleenum_field, position91)
			}
			return true
		l90:
			position, tokenIndex, depth = position90, tokenIndex90, depth90
			return false
		},
		/* 13 root_decl <- <('r' 'o' 'o' 't' '_' 't' 'y' 'p' 'e' spacing ident spacing ';' spacing)> */
		func() bool {
			position94, tokenIndex94, depth94 := position, tokenIndex, depth
			{
				position95 := position
				depth++
				if buffer[position] != rune('r') {
					goto l94
				}
				position++
				if buffer[position] != rune('o') {
					goto l94
				}
				position++
				if buffer[position] != rune('o') {
					goto l94
				}
				position++
				if buffer[position] != rune('t') {
					goto l94
				}
				position++
				if buffer[position] != rune('_') {
					goto l94
				}
				position++
				if buffer[position] != rune('t') {
					goto l94
				}
				position++
				if buffer[position] != rune('y') {
					goto l94
				}
				position++
				if buffer[position] != rune('p') {
					goto l94
				}
				position++
				if buffer[position] != rune('e') {
					goto l94
				}
				position++
				if !_rules[rulespacing]() {
					goto l94
				}
				if !_rules[ruleident]() {
					goto l94
				}
				if !_rules[rulespacing]() {
					goto l94
				}
				if buffer[position] != rune(';') {
					goto l94
				}
				position++
				if !_rules[rulespacing]() {
					goto l94
				}
				depth--
				add(ruleroot_decl, position95)
			}
			return true
		l94:
			position, tokenIndex, depth = position94, tokenIndex94, depth94
			return false
		},
		/* 14 file_extension_decl <- <('f' 'i' 'l' 'e' '_' 'e' 'x' 't' 'e' 'n' 's' 'i' 'o' 'n' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position96, tokenIndex96, depth96 := position, tokenIndex, depth
			{
				position97 := position
				depth++
				if buffer[position] != rune('f') {
					goto l96
				}
				position++
				if buffer[position] != rune('i') {
					goto l96
				}
				position++
				if buffer[position] != rune('l') {
					goto l96
				}
				position++
				if buffer[position] != rune('e') {
					goto l96
				}
				position++
				if buffer[position] != rune('_') {
					goto l96
				}
				position++
				if buffer[position] != rune('e') {
					goto l96
				}
				position++
				if buffer[position] != rune('x') {
					goto l96
				}
				position++
				if buffer[position] != rune('t') {
					goto l96
				}
				position++
				if buffer[position] != rune('e') {
					goto l96
				}
				position++
				if buffer[position] != rune('n') {
					goto l96
				}
				position++
				if buffer[position] != rune('s') {
					goto l96
				}
				position++
				if buffer[position] != rune('i') {
					goto l96
				}
				position++
				if buffer[position] != rune('o') {
					goto l96
				}
				position++
				if buffer[position] != rune('n') {
					goto l96
				}
				position++
				{
					position98 := position
					depth++
				l99:
					{
						position100, tokenIndex100, depth100 := position, tokenIndex, depth
						{
							position101, tokenIndex101, depth101 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l102
							}
							position++
							goto l101
						l102:
							position, tokenIndex, depth = position101, tokenIndex101, depth101
							if buffer[position] != rune('\t') {
								goto l100
							}
							position++
						}
					l101:
						goto l99
					l100:
						position, tokenIndex, depth = position100, tokenIndex100, depth100
					}
					depth--
					add(rulePegText, position98)
				}
				{
					position103 := position
					depth++
					{
						position106, tokenIndex106, depth106 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l106
						}
						position++
						goto l96
					l106:
						position, tokenIndex, depth = position106, tokenIndex106, depth106
					}
					if !matchDot() {
						goto l96
					}
				l104:
					{
						position105, tokenIndex105, depth105 := position, tokenIndex, depth
						{
							position107, tokenIndex107, depth107 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l107
							}
							position++
							goto l105
						l107:
							position, tokenIndex, depth = position107, tokenIndex107, depth107
						}
						if !matchDot() {
							goto l105
						}
						goto l104
					l105:
						position, tokenIndex, depth = position105, tokenIndex105, depth105
					}
					depth--
					add(rulePegText, position103)
				}
				if buffer[position] != rune(';') {
					goto l96
				}
				position++
				if !_rules[rulespacing]() {
					goto l96
				}
				depth--
				add(rulefile_extension_decl, position97)
			}
			return true
		l96:
			position, tokenIndex, depth = position96, tokenIndex96, depth96
			return false
		},
		/* 15 file_identifier_decl <- <('f' 'i' 'l' 'e' '_' 'i' 'd' 'e' 'n' 't' 'i' 'f' 'i' 'e' 'r' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position108, tokenIndex108, depth108 := position, tokenIndex, depth
			{
				position109 := position
				depth++
				if buffer[position] != rune('f') {
					goto l108
				}
				position++
				if buffer[position] != rune('i') {
					goto l108
				}
				position++
				if buffer[position] != rune('l') {
					goto l108
				}
				position++
				if buffer[position] != rune('e') {
					goto l108
				}
				position++
				if buffer[position] != rune('_') {
					goto l108
				}
				position++
				if buffer[position] != rune('i') {
					goto l108
				}
				position++
				if buffer[position] != rune('d') {
					goto l108
				}
				position++
				if buffer[position] != rune('e') {
					goto l108
				}
				position++
				if buffer[position] != rune('n') {
					goto l108
				}
				position++
				if buffer[position] != rune('t') {
					goto l108
				}
				position++
				if buffer[position] != rune('i') {
					goto l108
				}
				position++
				if buffer[position] != rune('f') {
					goto l108
				}
				position++
				if buffer[position] != rune('i') {
					goto l108
				}
				position++
				if buffer[position] != rune('e') {
					goto l108
				}
				position++
				if buffer[position] != rune('r') {
					goto l108
				}
				position++
				{
					position110 := position
					depth++
				l111:
					{
						position112, tokenIndex112, depth112 := position, tokenIndex, depth
						{
							position113, tokenIndex113, depth113 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l114
							}
							position++
							goto l113
						l114:
							position, tokenIndex, depth = position113, tokenIndex113, depth113
							if buffer[position] != rune('\t') {
								goto l112
							}
							position++
						}
					l113:
						goto l111
					l112:
						position, tokenIndex, depth = position112, tokenIndex112, depth112
					}
					depth--
					add(rulePegText, position110)
				}
				{
					position115 := position
					depth++
					{
						position118, tokenIndex118, depth118 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l118
						}
						position++
						goto l108
					l118:
						position, tokenIndex, depth = position118, tokenIndex118, depth118
					}
					if !matchDot() {
						goto l108
					}
				l116:
					{
						position117, tokenIndex117, depth117 := position, tokenIndex, depth
						{
							position119, tokenIndex119, depth119 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l119
							}
							position++
							goto l117
						l119:
							position, tokenIndex, depth = position119, tokenIndex119, depth119
						}
						if !matchDot() {
							goto l117
						}
						goto l116
					l117:
						position, tokenIndex, depth = position117, tokenIndex117, depth117
					}
					depth--
					add(rulePegText, position115)
				}
				if buffer[position] != rune(';') {
					goto l108
				}
				position++
				if !_rules[rulespacing]() {
					goto l108
				}
				depth--
				add(rulefile_identifier_decl, position109)
			}
			return true
		l108:
			position, tokenIndex, depth = position108, tokenIndex108, depth108
			return false
		},
		/* 16 attribute_decl <- <('a' 't' 't' 'r' 'i' 'b' 'u' 't' 'e' <(' ' / '\t')*> <(!';' .)+> ';' spacing)> */
		func() bool {
			position120, tokenIndex120, depth120 := position, tokenIndex, depth
			{
				position121 := position
				depth++
				if buffer[position] != rune('a') {
					goto l120
				}
				position++
				if buffer[position] != rune('t') {
					goto l120
				}
				position++
				if buffer[position] != rune('t') {
					goto l120
				}
				position++
				if buffer[position] != rune('r') {
					goto l120
				}
				position++
				if buffer[position] != rune('i') {
					goto l120
				}
				position++
				if buffer[position] != rune('b') {
					goto l120
				}
				position++
				if buffer[position] != rune('u') {
					goto l120
				}
				position++
				if buffer[position] != rune('t') {
					goto l120
				}
				position++
				if buffer[position] != rune('e') {
					goto l120
				}
				position++
				{
					position122 := position
					depth++
				l123:
					{
						position124, tokenIndex124, depth124 := position, tokenIndex, depth
						{
							position125, tokenIndex125, depth125 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l126
							}
							position++
							goto l125
						l126:
							position, tokenIndex, depth = position125, tokenIndex125, depth125
							if buffer[position] != rune('\t') {
								goto l124
							}
							position++
						}
					l125:
						goto l123
					l124:
						position, tokenIndex, depth = position124, tokenIndex124, depth124
					}
					depth--
					add(rulePegText, position122)
				}
				{
					position127 := position
					depth++
					{
						position130, tokenIndex130, depth130 := position, tokenIndex, depth
						if buffer[position] != rune(';') {
							goto l130
						}
						position++
						goto l120
					l130:
						position, tokenIndex, depth = position130, tokenIndex130, depth130
					}
					if !matchDot() {
						goto l120
					}
				l128:
					{
						position129, tokenIndex129, depth129 := position, tokenIndex, depth
						{
							position131, tokenIndex131, depth131 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l131
							}
							position++
							goto l129
						l131:
							position, tokenIndex, depth = position131, tokenIndex131, depth131
						}
						if !matchDot() {
							goto l129
						}
						goto l128
					l129:
						position, tokenIndex, depth = position129, tokenIndex129, depth129
					}
					depth--
					add(rulePegText, position127)
				}
				if buffer[position] != rune(';') {
					goto l120
				}
				position++
				if !_rules[rulespacing]() {
					goto l120
				}
				depth--
				add(ruleattribute_decl, position121)
			}
			return true
		l120:
			position, tokenIndex, depth = position120, tokenIndex120, depth120
			return false
		},
		/* 17 rpc_decl <- <('r' 'p' 'c' '_' 's' 'e' 'r' 'v' 'i' 'c' 'e' <(' ' / '\t')*> ident '{' <(!'}' .)+> '}' spacing)> */
		func() bool {
			position132, tokenIndex132, depth132 := position, tokenIndex, depth
			{
				position133 := position
				depth++
				if buffer[position] != rune('r') {
					goto l132
				}
				position++
				if buffer[position] != rune('p') {
					goto l132
				}
				position++
				if buffer[position] != rune('c') {
					goto l132
				}
				position++
				if buffer[position] != rune('_') {
					goto l132
				}
				position++
				if buffer[position] != rune('s') {
					goto l132
				}
				position++
				if buffer[position] != rune('e') {
					goto l132
				}
				position++
				if buffer[position] != rune('r') {
					goto l132
				}
				position++
				if buffer[position] != rune('v') {
					goto l132
				}
				position++
				if buffer[position] != rune('i') {
					goto l132
				}
				position++
				if buffer[position] != rune('c') {
					goto l132
				}
				position++
				if buffer[position] != rune('e') {
					goto l132
				}
				position++
				{
					position134 := position
					depth++
				l135:
					{
						position136, tokenIndex136, depth136 := position, tokenIndex, depth
						{
							position137, tokenIndex137, depth137 := position, tokenIndex, depth
							if buffer[position] != rune(' ') {
								goto l138
							}
							position++
							goto l137
						l138:
							position, tokenIndex, depth = position137, tokenIndex137, depth137
							if buffer[position] != rune('\t') {
								goto l136
							}
							position++
						}
					l137:
						goto l135
					l136:
						position, tokenIndex, depth = position136, tokenIndex136, depth136
					}
					depth--
					add(rulePegText, position134)
				}
				if !_rules[ruleident]() {
					goto l132
				}
				if buffer[position] != rune('{') {
					goto l132
				}
				position++
				{
					position139 := position
					depth++
					{
						position142, tokenIndex142, depth142 := position, tokenIndex, depth
						if buffer[position] != rune('}') {
							goto l142
						}
						position++
						goto l132
					l142:
						position, tokenIndex, depth = position142, tokenIndex142, depth142
					}
					if !matchDot() {
						goto l132
					}
				l140:
					{
						position141, tokenIndex141, depth141 := position, tokenIndex, depth
						{
							position143, tokenIndex143, depth143 := position, tokenIndex, depth
							if buffer[position] != rune('}') {
								goto l143
							}
							position++
							goto l141
						l143:
							position, tokenIndex, depth = position143, tokenIndex143, depth143
						}
						if !matchDot() {
							goto l141
						}
						goto l140
					l141:
						position, tokenIndex, depth = position141, tokenIndex141, depth141
					}
					depth--
					add(rulePegText, position139)
				}
				if buffer[position] != rune('}') {
					goto l132
				}
				position++
				if !_rules[rulespacing]() {
					goto l132
				}
				depth--
				add(rulerpc_decl, position133)
			}
			return true
		l132:
			position, tokenIndex, depth = position132, tokenIndex132, depth132
			return false
		},
		/* 18 type <- <(('b' 'o' 'o' 'l' spacing Action6) / ('b' 'y' 't' 'e' spacing Action7) / ('u' 'b' 'y' 't' 'e' spacing Action8) / ('s' 'h' 'o' 'r' 't' spacing Action9) / ('u' 's' 'h' 'o' 'r' 't' spacing Action10) / ('i' 'n' 't' spacing Action11) / ('u' 'i' 'n' 't' spacing Action12) / ('f' 'l' 'o' 'a' 't' spacing Action13) / ('l' 'o' 'n' 'g' spacing Action14) / ('u' 'l' 'o' 'n' 'g' spacing Action15) / ('d' 'o' 'u' 'b' 'l' 'e' spacing Action16) / ('i' 'n' 't' '8' spacing Action17) / ('i' 'n' 't' '1' '6' spacing Action18) / ('u' 'i' 'n' 't' '1' '6' spacing Action19) / ('i' 'n' 't' '3' '2' spacing Action20) / ('u' 'i' 'n' 't' '3' '2' spacing Action21) / ('i' 'n' 't' '6' '4' spacing Action22) / ('u' 'i' 'n' 't' '6' '4' spacing Action23) / ('f' 'l' 'o' 'a' 't' '3' '2' spacing Action24) / ('f' 'l' 'o' 'a' 't' '6' '4' spacing Action25) / ('s' 't' 'r' 'i' 'n' 'g' spacing Action26) / (ident spacing Action27) / ('[' type ']' spacing Action28))> */
		func() bool {
			position144, tokenIndex144, depth144 := position, tokenIndex, depth
			{
				position145 := position
				depth++
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l147
					}
					position++
					if buffer[position] != rune('o') {
						goto l147
					}
					position++
					if buffer[position] != rune('o') {
						goto l147
					}
					position++
					if buffer[position] != rune('l') {
						goto l147
					}
					position++
					if !_rules[rulespacing]() {
						goto l147
					}
					if !_rules[ruleAction6]() {
						goto l147
					}
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('b') {
						goto l148
					}
					position++
					if buffer[position] != rune('y') {
						goto l148
					}
					position++
					if buffer[position] != rune('t') {
						goto l148
					}
					position++
					if buffer[position] != rune('e') {
						goto l148
					}
					position++
					if !_rules[rulespacing]() {
						goto l148
					}
					if !_rules[ruleAction7]() {
						goto l148
					}
					goto l146
				l148:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('u') {
						goto l149
					}
					position++
					if buffer[position] != rune('b') {
						goto l149
					}
					position++
					if buffer[position] != rune('y') {
						goto l149
					}
					position++
					if buffer[position] != rune('t') {
						goto l149
					}
					position++
					if buffer[position] != rune('e') {
						goto l149
					}
					position++
					if !_rules[rulespacing]() {
						goto l149
					}
					if !_rules[ruleAction8]() {
						goto l149
					}
					goto l146
				l149:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('s') {
						goto l150
					}
					position++
					if buffer[position] != rune('h') {
						goto l150
					}
					position++
					if buffer[position] != rune('o') {
						goto l150
					}
					position++
					if buffer[position] != rune('r') {
						goto l150
					}
					position++
					if buffer[position] != rune('t') {
						goto l150
					}
					position++
					if !_rules[rulespacing]() {
						goto l150
					}
					if !_rules[ruleAction9]() {
						goto l150
					}
					goto l146
				l150:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('u') {
						goto l151
					}
					position++
					if buffer[position] != rune('s') {
						goto l151
					}
					position++
					if buffer[position] != rune('h') {
						goto l151
					}
					position++
					if buffer[position] != rune('o') {
						goto l151
					}
					position++
					if buffer[position] != rune('r') {
						goto l151
					}
					position++
					if buffer[position] != rune('t') {
						goto l151
					}
					position++
					if !_rules[rulespacing]() {
						goto l151
					}
					if !_rules[ruleAction10]() {
						goto l151
					}
					goto l146
				l151:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('i') {
						goto l152
					}
					position++
					if buffer[position] != rune('n') {
						goto l152
					}
					position++
					if buffer[position] != rune('t') {
						goto l152
					}
					position++
					if !_rules[rulespacing]() {
						goto l152
					}
					if !_rules[ruleAction11]() {
						goto l152
					}
					goto l146
				l152:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('u') {
						goto l153
					}
					position++
					if buffer[position] != rune('i') {
						goto l153
					}
					position++
					if buffer[position] != rune('n') {
						goto l153
					}
					position++
					if buffer[position] != rune('t') {
						goto l153
					}
					position++
					if !_rules[rulespacing]() {
						goto l153
					}
					if !_rules[ruleAction12]() {
						goto l153
					}
					goto l146
				l153:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('f') {
						goto l154
					}
					position++
					if buffer[position] != rune('l') {
						goto l154
					}
					position++
					if buffer[position] != rune('o') {
						goto l154
					}
					position++
					if buffer[position] != rune('a') {
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
					if !_rules[ruleAction13]() {
						goto l154
					}
					goto l146
				l154:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('l') {
						goto l155
					}
					position++
					if buffer[position] != rune('o') {
						goto l155
					}
					position++
					if buffer[position] != rune('n') {
						goto l155
					}
					position++
					if buffer[position] != rune('g') {
						goto l155
					}
					position++
					if !_rules[rulespacing]() {
						goto l155
					}
					if !_rules[ruleAction14]() {
						goto l155
					}
					goto l146
				l155:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('u') {
						goto l156
					}
					position++
					if buffer[position] != rune('l') {
						goto l156
					}
					position++
					if buffer[position] != rune('o') {
						goto l156
					}
					position++
					if buffer[position] != rune('n') {
						goto l156
					}
					position++
					if buffer[position] != rune('g') {
						goto l156
					}
					position++
					if !_rules[rulespacing]() {
						goto l156
					}
					if !_rules[ruleAction15]() {
						goto l156
					}
					goto l146
				l156:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('d') {
						goto l157
					}
					position++
					if buffer[position] != rune('o') {
						goto l157
					}
					position++
					if buffer[position] != rune('u') {
						goto l157
					}
					position++
					if buffer[position] != rune('b') {
						goto l157
					}
					position++
					if buffer[position] != rune('l') {
						goto l157
					}
					position++
					if buffer[position] != rune('e') {
						goto l157
					}
					position++
					if !_rules[rulespacing]() {
						goto l157
					}
					if !_rules[ruleAction16]() {
						goto l157
					}
					goto l146
				l157:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
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
					if buffer[position] != rune('8') {
						goto l158
					}
					position++
					if !_rules[rulespacing]() {
						goto l158
					}
					if !_rules[ruleAction17]() {
						goto l158
					}
					goto l146
				l158:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
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
					if buffer[position] != rune('1') {
						goto l159
					}
					position++
					if buffer[position] != rune('6') {
						goto l159
					}
					position++
					if !_rules[rulespacing]() {
						goto l159
					}
					if !_rules[ruleAction18]() {
						goto l159
					}
					goto l146
				l159:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
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
					if buffer[position] != rune('1') {
						goto l160
					}
					position++
					if buffer[position] != rune('6') {
						goto l160
					}
					position++
					if !_rules[rulespacing]() {
						goto l160
					}
					if !_rules[ruleAction19]() {
						goto l160
					}
					goto l146
				l160:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('i') {
						goto l161
					}
					position++
					if buffer[position] != rune('n') {
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
					if !_rules[ruleAction20]() {
						goto l161
					}
					goto l146
				l161:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('u') {
						goto l162
					}
					position++
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
					if buffer[position] != rune('3') {
						goto l162
					}
					position++
					if buffer[position] != rune('2') {
						goto l162
					}
					position++
					if !_rules[rulespacing]() {
						goto l162
					}
					if !_rules[ruleAction21]() {
						goto l162
					}
					goto l146
				l162:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
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
					if buffer[position] != rune('6') {
						goto l163
					}
					position++
					if buffer[position] != rune('4') {
						goto l163
					}
					position++
					if !_rules[rulespacing]() {
						goto l163
					}
					if !_rules[ruleAction22]() {
						goto l163
					}
					goto l146
				l163:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
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
					if buffer[position] != rune('6') {
						goto l164
					}
					position++
					if buffer[position] != rune('4') {
						goto l164
					}
					position++
					if !_rules[rulespacing]() {
						goto l164
					}
					if !_rules[ruleAction23]() {
						goto l164
					}
					goto l146
				l164:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('f') {
						goto l165
					}
					position++
					if buffer[position] != rune('l') {
						goto l165
					}
					position++
					if buffer[position] != rune('o') {
						goto l165
					}
					position++
					if buffer[position] != rune('a') {
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
					if !_rules[ruleAction24]() {
						goto l165
					}
					goto l146
				l165:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('f') {
						goto l166
					}
					position++
					if buffer[position] != rune('l') {
						goto l166
					}
					position++
					if buffer[position] != rune('o') {
						goto l166
					}
					position++
					if buffer[position] != rune('a') {
						goto l166
					}
					position++
					if buffer[position] != rune('t') {
						goto l166
					}
					position++
					if buffer[position] != rune('6') {
						goto l166
					}
					position++
					if buffer[position] != rune('4') {
						goto l166
					}
					position++
					if !_rules[rulespacing]() {
						goto l166
					}
					if !_rules[ruleAction25]() {
						goto l166
					}
					goto l146
				l166:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('s') {
						goto l167
					}
					position++
					if buffer[position] != rune('t') {
						goto l167
					}
					position++
					if buffer[position] != rune('r') {
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
					if buffer[position] != rune('g') {
						goto l167
					}
					position++
					if !_rules[rulespacing]() {
						goto l167
					}
					if !_rules[ruleAction26]() {
						goto l167
					}
					goto l146
				l167:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if !_rules[ruleident]() {
						goto l168
					}
					if !_rules[rulespacing]() {
						goto l168
					}
					if !_rules[ruleAction27]() {
						goto l168
					}
					goto l146
				l168:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('[') {
						goto l144
					}
					position++
					if !_rules[ruletype]() {
						goto l144
					}
					if buffer[position] != rune(']') {
						goto l144
					}
					position++
					if !_rules[rulespacing]() {
						goto l144
					}
					if !_rules[ruleAction28]() {
						goto l144
					}
				}
			l146:
				depth--
				add(ruletype, position145)
			}
			return true
		l144:
			position, tokenIndex, depth = position144, tokenIndex144, depth144
			return false
		},
		/* 19 scalar <- <(integer_constant / float_constant)> */
		func() bool {
			position169, tokenIndex169, depth169 := position, tokenIndex, depth
			{
				position170 := position
				depth++
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if !_rules[ruleinteger_constant]() {
						goto l172
					}
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if !_rules[rulefloat_constant]() {
						goto l169
					}
				}
			l171:
				depth--
				add(rulescalar, position170)
			}
			return true
		l169:
			position, tokenIndex, depth = position169, tokenIndex169, depth169
			return false
		},
		/* 20 integer_constant <- <(<[0-9]+> / ('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> */
		func() bool {
			position173, tokenIndex173, depth173 := position, tokenIndex, depth
			{
				position174 := position
				depth++
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					{
						position177 := position
						depth++
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l176
						}
						position++
					l178:
						{
							position179, tokenIndex179, depth179 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l179
							}
							position++
							goto l178
						l179:
							position, tokenIndex, depth = position179, tokenIndex179, depth179
						}
						depth--
						add(rulePegText, position177)
					}
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('t') {
						goto l180
					}
					position++
					if buffer[position] != rune('r') {
						goto l180
					}
					position++
					if buffer[position] != rune('u') {
						goto l180
					}
					position++
					if buffer[position] != rune('e') {
						goto l180
					}
					position++
					goto l175
				l180:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('f') {
						goto l173
					}
					position++
					if buffer[position] != rune('a') {
						goto l173
					}
					position++
					if buffer[position] != rune('l') {
						goto l173
					}
					position++
					if buffer[position] != rune('s') {
						goto l173
					}
					position++
					if buffer[position] != rune('e') {
						goto l173
					}
					position++
				}
			l175:
				depth--
				add(ruleinteger_constant, position174)
			}
			return true
		l173:
			position, tokenIndex, depth = position173, tokenIndex173, depth173
			return false
		},
		/* 21 float_constant <- <(<('-'* [0-9]+ . [0-9])> / float_constant_exp)> */
		func() bool {
			position181, tokenIndex181, depth181 := position, tokenIndex, depth
			{
				position182 := position
				depth++
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					{
						position185 := position
						depth++
					l186:
						{
							position187, tokenIndex187, depth187 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l187
							}
							position++
							goto l186
						l187:
							position, tokenIndex, depth = position187, tokenIndex187, depth187
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l184
						}
						position++
					l188:
						{
							position189, tokenIndex189, depth189 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l189
							}
							position++
							goto l188
						l189:
							position, tokenIndex, depth = position189, tokenIndex189, depth189
						}
						if !matchDot() {
							goto l184
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l184
						}
						position++
						depth--
						add(rulePegText, position185)
					}
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					if !_rules[rulefloat_constant_exp]() {
						goto l181
					}
				}
			l183:
				depth--
				add(rulefloat_constant, position182)
			}
			return true
		l181:
			position, tokenIndex, depth = position181, tokenIndex181, depth181
			return false
		},
		/* 22 float_constant_exp <- <(<('-'* [0-9]+ . [0-9]+)> <('e' / 'E')> <([+-]] / '>' / ' ' / '<' / '[' / [0-9])+>)> */
		func() bool {
			position190, tokenIndex190, depth190 := position, tokenIndex, depth
			{
				position191 := position
				depth++
				{
					position192 := position
					depth++
				l193:
					{
						position194, tokenIndex194, depth194 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l194
						}
						position++
						goto l193
					l194:
						position, tokenIndex, depth = position194, tokenIndex194, depth194
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l190
					}
					position++
				l195:
					{
						position196, tokenIndex196, depth196 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l196
						}
						position++
						goto l195
					l196:
						position, tokenIndex, depth = position196, tokenIndex196, depth196
					}
					if !matchDot() {
						goto l190
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l190
					}
					position++
				l197:
					{
						position198, tokenIndex198, depth198 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l198
						}
						position++
						goto l197
					l198:
						position, tokenIndex, depth = position198, tokenIndex198, depth198
					}
					depth--
					add(rulePegText, position192)
				}
				{
					position199 := position
					depth++
					{
						position200, tokenIndex200, depth200 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l201
						}
						position++
						goto l200
					l201:
						position, tokenIndex, depth = position200, tokenIndex200, depth200
						if buffer[position] != rune('E') {
							goto l190
						}
						position++
					}
				l200:
					depth--
					add(rulePegText, position199)
				}
				{
					position202 := position
					depth++
					{
						position205, tokenIndex205, depth205 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('+') || c > rune(']') {
							goto l206
						}
						position++
						goto l205
					l206:
						position, tokenIndex, depth = position205, tokenIndex205, depth205
						if buffer[position] != rune('>') {
							goto l207
						}
						position++
						goto l205
					l207:
						position, tokenIndex, depth = position205, tokenIndex205, depth205
						if buffer[position] != rune(' ') {
							goto l208
						}
						position++
						goto l205
					l208:
						position, tokenIndex, depth = position205, tokenIndex205, depth205
						if buffer[position] != rune('<') {
							goto l209
						}
						position++
						goto l205
					l209:
						position, tokenIndex, depth = position205, tokenIndex205, depth205
						if buffer[position] != rune('[') {
							goto l210
						}
						position++
						goto l205
					l210:
						position, tokenIndex, depth = position205, tokenIndex205, depth205
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l190
						}
						position++
					}
				l205:
				l203:
					{
						position204, tokenIndex204, depth204 := position, tokenIndex, depth
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
								goto l204
							}
							position++
						}
					l211:
						goto l203
					l204:
						position, tokenIndex, depth = position204, tokenIndex204, depth204
					}
					depth--
					add(rulePegText, position202)
				}
				depth--
				add(rulefloat_constant_exp, position191)
			}
			return true
		l190:
			position, tokenIndex, depth = position190, tokenIndex190, depth190
			return false
		},
		/* 23 ident <- <<(([a-z] / [A-Z] / '_') ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position217, tokenIndex217, depth217 := position, tokenIndex, depth
			{
				position218 := position
				depth++
				{
					position219 := position
					depth++
					{
						position220, tokenIndex220, depth220 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l221
						}
						position++
						goto l220
					l221:
						position, tokenIndex, depth = position220, tokenIndex220, depth220
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l222
						}
						position++
						goto l220
					l222:
						position, tokenIndex, depth = position220, tokenIndex220, depth220
						if buffer[position] != rune('_') {
							goto l217
						}
						position++
					}
				l220:
				l223:
					{
						position224, tokenIndex224, depth224 := position, tokenIndex, depth
						{
							position225, tokenIndex225, depth225 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l226
							}
							position++
							goto l225
						l226:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l227
							}
							position++
							goto l225
						l227:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l228
							}
							position++
							goto l225
						l228:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							if buffer[position] != rune('_') {
								goto l224
							}
							position++
						}
					l225:
						goto l223
					l224:
						position, tokenIndex, depth = position224, tokenIndex224, depth224
					}
					depth--
					add(rulePegText, position219)
				}
				depth--
				add(ruleident, position218)
			}
			return true
		l217:
			position, tokenIndex, depth = position217, tokenIndex217, depth217
			return false
		},
		/* 24 only_comment <- <(spacing ';')> */
		func() bool {
			position229, tokenIndex229, depth229 := position, tokenIndex, depth
			{
				position230 := position
				depth++
				if !_rules[rulespacing]() {
					goto l229
				}
				if buffer[position] != rune(';') {
					goto l229
				}
				position++
				depth--
				add(ruleonly_comment, position230)
			}
			return true
		l229:
			position, tokenIndex, depth = position229, tokenIndex229, depth229
			return false
		},
		/* 25 spacing <- <space_comment*> */
		func() bool {
			{
				position232 := position
				depth++
			l233:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if !_rules[rulespace_comment]() {
						goto l234
					}
					goto l233
				l234:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
				}
				depth--
				add(rulespacing, position232)
			}
			return true
		},
		/* 26 space_comment <- <(space / comment)> */
		func() bool {
			position235, tokenIndex235, depth235 := position, tokenIndex, depth
			{
				position236 := position
				depth++
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if !_rules[rulespace]() {
						goto l238
					}
					goto l237
				l238:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
					if !_rules[rulecomment]() {
						goto l235
					}
				}
			l237:
				depth--
				add(rulespace_comment, position236)
			}
			return true
		l235:
			position, tokenIndex, depth = position235, tokenIndex235, depth235
			return false
		},
		/* 27 comment <- <('/' '/' (!end_of_line .)* end_of_line)> */
		func() bool {
			position239, tokenIndex239, depth239 := position, tokenIndex, depth
			{
				position240 := position
				depth++
				if buffer[position] != rune('/') {
					goto l239
				}
				position++
				if buffer[position] != rune('/') {
					goto l239
				}
				position++
			l241:
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					{
						position243, tokenIndex243, depth243 := position, tokenIndex, depth
						if !_rules[ruleend_of_line]() {
							goto l243
						}
						goto l242
					l243:
						position, tokenIndex, depth = position243, tokenIndex243, depth243
					}
					if !matchDot() {
						goto l242
					}
					goto l241
				l242:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
				}
				if !_rules[ruleend_of_line]() {
					goto l239
				}
				depth--
				add(rulecomment, position240)
			}
			return true
		l239:
			position, tokenIndex, depth = position239, tokenIndex239, depth239
			return false
		},
		/* 28 space <- <(' ' / '\t' / end_of_line)> */
		func() bool {
			position244, tokenIndex244, depth244 := position, tokenIndex, depth
			{
				position245 := position
				depth++
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if buffer[position] != rune('\t') {
						goto l248
					}
					position++
					goto l246
				l248:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if !_rules[ruleend_of_line]() {
						goto l244
					}
				}
			l246:
				depth--
				add(rulespace, position245)
			}
			return true
		l244:
			position, tokenIndex, depth = position244, tokenIndex244, depth244
			return false
		},
		/* 29 end_of_line <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position249, tokenIndex249, depth249 := position, tokenIndex, depth
			{
				position250 := position
				depth++
				{
					position251, tokenIndex251, depth251 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l252
					}
					position++
					if buffer[position] != rune('\n') {
						goto l252
					}
					position++
					goto l251
				l252:
					position, tokenIndex, depth = position251, tokenIndex251, depth251
					if buffer[position] != rune('\n') {
						goto l253
					}
					position++
					goto l251
				l253:
					position, tokenIndex, depth = position251, tokenIndex251, depth251
					if buffer[position] != rune('\r') {
						goto l249
					}
					position++
				}
			l251:
				depth--
				add(ruleend_of_line, position250)
			}
			return true
		l249:
			position, tokenIndex, depth = position249, tokenIndex249, depth249
			return false
		},
		/* 30 end_of_file <- <!.> */
		func() bool {
			position254, tokenIndex254, depth254 := position, tokenIndex, depth
			{
				position255 := position
				depth++
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if !matchDot() {
						goto l256
					}
					goto l254
				l256:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
				}
				depth--
				add(ruleend_of_file, position255)
			}
			return true
		l254:
			position, tokenIndex, depth = position254, tokenIndex254, depth254
			return false
		},
		nil,
		/* 33 Action0 <- <{p.SetNameSpace(text)}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 34 Action1 <- <{p.ExtractStruct()}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 35 Action2 <- <{p.SetTypeName(text)}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 36 Action3 <- <{p.NewExtractField()}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 37 Action4 <- <{p.NewExtractFieldWithValue()}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 38 Action5 <- <{p.FieldNaame(text)}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 39 Action6 <- <{p.SetType("bool")}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 40 Action7 <- <{p.SetType("byte")}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 41 Action8 <- <{p.SetType("ubyte")}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 42 Action9 <- <{p.SetType("short")}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 43 Action10 <- <{p.SetType("ushort")}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 44 Action11 <- <{p.SetType("int")}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 45 Action12 <- <{p.SetType("uint")}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 46 Action13 <- <{p.SetType("float")}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 47 Action14 <- <{p.SetType("long")}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 48 Action15 <- <{p.SetType("ulong")}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 49 Action16 <- <{p.SetType("double")}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 50 Action17 <- <{p.SetType("int8")}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 51 Action18 <- <{p.SetType("int16")}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 52 Action19 <- <{p.SetType("uint16")}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 53 Action20 <- <{p.SetType("int32")}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 54 Action21 <- <{p.SetType("uint32")}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 55 Action22 <- <{p.SetType("int64")}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 56 Action23 <- <{p.SetType("uint64")}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 57 Action24 <- <{p.SetType("float32")}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 58 Action25 <- <{p.SetType("float64")}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 59 Action26 <- <{p.SetType("string")}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 60 Action27 <- <{p.SetType(text)}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 61 Action28 <- <{p.SetRepeated(text) }> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
	}
	p.rules = _rules
}
