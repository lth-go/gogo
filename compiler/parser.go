//line parser.go.y:2
package compiler

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"github.com/lth-go/gogogogo/vm"
	"strconv"
)

//line parser.go.y:10
type yySymType struct {
	yys            int
	parameter_list []*Parameter
	argument_list  []Expression

	statement      Statement
	statement_list []Statement

	expression      Expression
	expression_list []Expression

	block     *Block
	elif_list []*Elif

	basic_type_specifier *TypeSpecifier
	type_specifier       *TypeSpecifier

	array_dimension      *ArrayDimension
	array_dimension_list []*ArrayDimension

	package_name []string
	require_list []*Require

	extends_list        []*Extend
	member_declaration  []MemberDeclaration
	function_definition *FunctionDefinition

	class_name []string

	tok Token
}

const IF = 57346
const ELSE = 57347
const ELIF = 57348
const FOR = 57349
const RETURN_T = 57350
const BREAK = 57351
const CONTINUE = 57352
const LP = 57353
const RP = 57354
const LC = 57355
const RC = 57356
const LB = 57357
const RB = 57358
const SEMICOLON = 57359
const COMMA = 57360
const COLON = 57361
const ASSIGN_T = 57362
const LOGICAL_AND = 57363
const LOGICAL_OR = 57364
const EQ = 57365
const NE = 57366
const GT = 57367
const GE = 57368
const LT = 57369
const LE = 57370
const ADD = 57371
const SUB = 57372
const MUL = 57373
const DIV = 57374
const INT_LITERAL = 57375
const DOUBLE_LITERAL = 57376
const STRING_LITERAL = 57377
const TRUE_T = 57378
const FALSE_T = 57379
const NULL_T = 57380
const IDENTIFIER = 57381
const EXCLAMATION = 57382
const DOT = 57383
const VOID_T = 57384
const BOOLEAN_T = 57385
const INT_T = 57386
const DOUBLE_T = 57387
const STRING_T = 57388
const NEW = 57389
const REQUIRE = 57390
const CLASS_T = 57391
const THIS_T = 57392

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IF",
	"ELSE",
	"ELIF",
	"FOR",
	"RETURN_T",
	"BREAK",
	"CONTINUE",
	"LP",
	"RP",
	"LC",
	"RC",
	"LB",
	"RB",
	"SEMICOLON",
	"COMMA",
	"COLON",
	"ASSIGN_T",
	"LOGICAL_AND",
	"LOGICAL_OR",
	"EQ",
	"NE",
	"GT",
	"GE",
	"LT",
	"LE",
	"ADD",
	"SUB",
	"MUL",
	"DIV",
	"INT_LITERAL",
	"DOUBLE_LITERAL",
	"STRING_LITERAL",
	"TRUE_T",
	"FALSE_T",
	"NULL_T",
	"IDENTIFIER",
	"EXCLAMATION",
	"DOT",
	"VOID_T",
	"BOOLEAN_T",
	"INT_T",
	"DOUBLE_T",
	"STRING_T",
	"NEW",
	"REQUIRE",
	"CLASS_T",
	"THIS_T",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:709

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 33,
	39, 18,
	-2, 63,
	-1, 88,
	15, 18,
	-2, 80,
	-1, 152,
	14, 122,
	-2, 120,
}

const yyPrivate = 57344

const yyLast = 465

var yyAct = [...]int{

	112, 10, 71, 147, 12, 202, 168, 9, 133, 22,
	132, 35, 122, 50, 54, 21, 36, 23, 53, 39,
	24, 25, 26, 27, 40, 5, 51, 213, 68, 150,
	72, 52, 28, 29, 30, 31, 32, 79, 102, 220,
	19, 217, 214, 56, 208, 84, 41, 42, 43, 44,
	45, 46, 33, 57, 179, 28, 29, 30, 31, 32,
	49, 94, 103, 48, 167, 87, 154, 78, 100, 100,
	146, 99, 101, 115, 109, 72, 121, 63, 88, 62,
	118, 28, 29, 30, 31, 32, 126, 120, 100, 124,
	86, 61, 100, 119, 100, 100, 127, 128, 135, 227,
	125, 100, 100, 100, 100, 130, 149, 100, 100, 100,
	100, 151, 144, 145, 142, 143, 79, 97, 98, 118,
	80, 136, 137, 138, 139, 77, 150, 95, 96, 28,
	29, 30, 31, 32, 76, 131, 82, 83, 108, 171,
	124, 169, 172, 166, 169, 105, 78, 184, 106, 177,
	64, 174, 89, 90, 91, 92, 185, 197, 104, 64,
	191, 215, 188, 192, 105, 72, 190, 106, 113, 160,
	218, 229, 193, 64, 194, 171, 198, 176, 200, 175,
	162, 181, 64, 116, 206, 176, 164, 209, 64, 211,
	163, 191, 210, 65, 64, 150, 162, 212, 28, 29,
	30, 31, 32, 206, 161, 140, 64, 219, 113, 141,
	23, 221, 233, 24, 25, 26, 27, 40, 222, 51,
	72, 74, 73, 129, 216, 149, 228, 226, 230, 64,
	232, 117, 223, 113, 111, 148, 56, 231, 224, 41,
	42, 43, 44, 45, 46, 33, 57, 110, 28, 29,
	30, 31, 32, 49, 23, 11, 48, 24, 25, 26,
	27, 40, 150, 51, 195, 28, 29, 30, 31, 32,
	113, 113, 225, 207, 199, 178, 180, 170, 134, 114,
	56, 81, 75, 41, 42, 43, 44, 45, 46, 33,
	57, 67, 28, 29, 30, 31, 32, 49, 66, 40,
	48, 51, 70, 159, 196, 113, 152, 186, 187, 155,
	157, 4, 183, 6, 182, 59, 58, 40, 56, 51,
	173, 41, 42, 43, 44, 45, 46, 69, 57, 158,
	8, 7, 2, 1, 205, 49, 56, 204, 48, 41,
	42, 43, 44, 45, 46, 69, 57, 40, 165, 51,
	203, 201, 107, 49, 153, 20, 48, 156, 189, 18,
	17, 16, 15, 14, 13, 93, 56, 38, 47, 41,
	42, 43, 44, 45, 46, 69, 57, 40, 123, 51,
	37, 55, 34, 49, 3, 60, 48, 85, 0, 0,
	0, 0, 0, 0, 0, 0, 56, 0, 0, 41,
	42, 43, 44, 45, 46, 69, 57, 40, 0, 51,
	0, 0, 117, 49, 0, 0, 48, 0, 0, 0,
	0, 0, 0, 0, 0, 40, 56, 51, 0, 41,
	42, 43, 44, 45, 46, 69, 57, 0, 0, 0,
	0, 0, 0, 49, 56, 0, 48, 41, 42, 43,
	44, 45, 46, 69, 57, 0, 0, 0, 0, 0,
	0, 49, 0, 0, 48,
}
var yyPact = [...]int{

	-23, 206, 206, -23, -1000, 52, -1000, -1000, -1000, -1000,
	40, 38, 176, -1000, -1000, -1000, -1000, -1000, -1000, 283,
	276, -1000, -1000, 414, 291, 414, 205, 204, -1000, -1000,
	-1000, -1000, -1000, 267, 112, 105, 99, 266, -1000, 113,
	414, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 39,
	127, 414, 98, 86, -1000, -1000, 414, 414, -1000, -1000,
	21, -1000, 147, 119, 414, -1000, 231, 218, 155, 264,
	414, 166, 132, -1000, -1000, 396, 414, 414, 37, 366,
	414, 414, 414, 414, 211, 94, 263, 263, -1000, 414,
	414, 414, 414, 191, -1000, 414, 414, 414, 414, -1000,
	26, -1000, -1000, 31, 223, -1000, 414, 293, 27, -1000,
	-1000, -1000, 304, 289, 414, 152, -1000, -1000, 188, 99,
	-1000, -1000, 178, -1000, -1000, 113, 170, 127, 127, -1000,
	336, 25, 262, -1000, 414, 262, 98, 98, 98, 98,
	-1000, 306, 86, 86, -1000, -1000, -1000, 167, 258, 15,
	261, 164, -1000, 129, -1000, 292, 302, 414, 250, -1000,
	414, -1000, 414, -1000, -1000, -1000, 162, -1000, 249, -1000,
	288, 141, 249, -1000, -1000, 257, -10, -1000, -1000, -1000,
	215, -1000, -10, 259, 5, -1000, 292, 414, 155, 13,
	-1000, 3, 144, -1000, -1000, 208, -1000, -1000, -1000, -1000,
	2, 156, -1000, -1000, -1000, -1000, 0, -1000, -1000, -1000,
	155, -1000, -1000, -1000, 128, 414, -1000, -1000, -1000, -1000,
	221, -1000, 260, 87, -1000, 292, 159, 220, -1000, 195,
	-1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 387, 385, 384, 311, 4, 2, 9, 16, 382,
	19, 13, 31, 18, 14, 381, 11, 380, 368, 367,
	365, 7, 364, 363, 362, 361, 360, 359, 358, 3,
	12, 0, 357, 40, 1, 15, 355, 8, 10, 6,
	354, 352, 5, 351, 350, 337, 334, 333, 332, 313,
	331, 330, 329, 314, 312,
}
var yyR1 = [...]int{

	0, 47, 47, 48, 48, 3, 3, 4, 2, 2,
	49, 49, 49, 33, 33, 33, 33, 33, 35, 36,
	36, 36, 34, 34, 34, 50, 50, 50, 50, 29,
	29, 30, 30, 28, 28, 5, 5, 7, 7, 9,
	9, 8, 8, 10, 10, 10, 11, 11, 11, 11,
	11, 12, 12, 12, 13, 13, 13, 14, 14, 14,
	15, 16, 16, 16, 17, 17, 17, 17, 17, 17,
	17, 17, 17, 17, 17, 17, 17, 17, 17, 17,
	1, 1, 18, 18, 19, 19, 19, 19, 38, 38,
	37, 39, 39, 20, 20, 20, 21, 21, 21, 21,
	21, 21, 21, 22, 22, 22, 22, 32, 32, 23,
	6, 6, 24, 25, 26, 27, 27, 52, 31, 31,
	53, 51, 54, 51, 41, 41, 40, 40, 43, 43,
	42, 42, 44, 46, 46, 46, 46, 45,
}
var yyR2 = [...]int{

	0, 2, 2, 0, 1, 1, 2, 3, 1, 3,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 3,
	3, 3, 1, 1, 1, 6, 5, 6, 5, 2,
	4, 1, 3, 1, 2, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 3, 3, 1, 3, 3, 3,
	3, 1, 3, 3, 1, 3, 3, 1, 2, 2,
	1, 1, 1, 1, 4, 4, 3, 4, 3, 3,
	1, 1, 1, 1, 1, 1, 1, 1, 4, 5,
	1, 3, 3, 4, 3, 4, 3, 4, 1, 2,
	3, 2, 3, 0, 1, 3, 2, 1, 1, 1,
	1, 1, 1, 3, 5, 4, 6, 3, 4, 9,
	0, 1, 3, 2, 2, 3, 5, 0, 4, 2,
	0, 7, 0, 6, 0, 2, 1, 3, 1, 2,
	1, 1, 1, 6, 5, 6, 5, 3,
}
var yyChk = [...]int{

	-1000, -47, -48, -3, -4, 48, -49, -50, -51, -21,
	-34, 49, -5, -22, -23, -24, -25, -26, -27, -33,
	-36, -35, -7, 4, 7, 8, 9, 10, 42, 43,
	44, 45, 46, 39, -9, -16, -8, -17, -19, -10,
	11, 33, 34, 35, 36, 37, 38, -18, 50, 47,
	-11, 13, -12, -13, -14, -15, 30, 40, -49, -4,
	-2, 39, 39, 39, 18, 17, 15, 15, -5, 39,
	11, -6, -5, 17, 17, 15, 22, 20, 41, 11,
	21, 15, 23, 24, -5, -1, -33, -35, 39, 25,
	26, 27, 28, -20, -7, 29, 30, 31, 32, -14,
	-16, -14, 17, 41, 11, 17, 20, -41, 19, -7,
	16, 16, -31, 13, 15, -6, 17, 16, -5, -8,
	-7, 39, -30, 12, -7, -10, -5, -11, -11, 12,
	11, 41, -38, -37, 15, -38, -12, -12, -12, -12,
	14, 18, -13, -13, -14, -14, 39, -29, 12, -34,
	39, -5, 13, -40, 39, 5, -32, 6, -52, 14,
	17, 16, 18, 12, 16, 12, -30, 39, -39, -37,
	15, -5, -39, 14, -7, 12, 18, -31, 17, 39,
	15, 17, -53, -54, 18, -31, 5, 6, -5, -28,
	-21, -34, -6, -7, 12, 15, 16, 16, -31, 17,
	-34, -43, -42, -44, -45, -46, -34, 14, 39, -31,
	-5, -31, -21, 14, 39, 17, 16, 39, 14, -42,
	39, -31, -6, 11, 17, 12, -29, 12, -31, 12,
	-31, 17, -31, 17,
}
var yyDef = [...]int{

	3, -2, 0, 4, 5, 0, 2, 10, 11, 12,
	0, 0, 0, 97, 98, 99, 100, 101, 102, 22,
	23, 24, 35, 0, 0, 110, 0, 0, 13, 14,
	15, 16, 17, -2, 37, 60, 39, 61, 62, 41,
	0, 70, 71, 72, 73, 74, 75, 76, 77, 0,
	43, 93, 46, 51, 54, 57, 0, 0, 1, 6,
	0, 8, 0, 124, 0, 96, 0, 0, 0, 63,
	110, 0, 111, 113, 114, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, -2, 0,
	0, 0, 0, 0, 94, 0, 0, 0, 0, 58,
	60, 59, 7, 0, 0, 115, 0, 0, 0, 36,
	19, 21, 103, 117, 0, 0, 112, 20, 0, 40,
	38, 66, 0, 68, 31, 42, 0, 44, 45, 69,
	0, 0, 84, 88, 0, 86, 47, 48, 49, 50,
	82, 0, 52, 53, 55, 56, 9, 0, 0, 0,
	18, 0, -2, 125, 126, 0, 105, 0, 0, 119,
	110, 65, 0, 67, 64, 78, 0, 81, 85, 89,
	0, 0, 87, 83, 95, 0, 0, 26, 28, 29,
	0, 116, 0, 0, 0, 104, 0, 0, 0, 0,
	33, 0, 0, 32, 79, 0, 91, 90, 25, 27,
	0, 0, 128, 130, 131, 132, 0, 123, 127, 106,
	0, 107, 34, 118, 0, 110, 92, 30, 121, 129,
	0, 108, 0, 0, 137, 0, 0, 0, 109, 0,
	134, 136, 133, 135,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:97
		{
			setRequireList(nil)
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:101
		{
			setRequireList(yyDollar[1].require_list)
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:108
		{
			yyVAL.require_list = chainRequireList(yyDollar[1].require_list, yyDollar[2].require_list)
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:114
		{
			yyVAL.require_list = createRequireList(yyDollar[2].package_name)
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:120
		{
			yyVAL.package_name = createPackageName(yyDollar[1].tok.Lit)
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:124
		{
			yyVAL.package_name = chainPackageName(yyDollar[1].package_name, yyDollar[3].tok.Lit)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:132
		{
			l := yylex.(*Lexer)
			l.compiler.statementList = append(l.compiler.statementList, yyDollar[1].statement)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:139
		{
			yyVAL.type_specifier = createTypeSpecifier(vm.VoidType, yyDollar[1].tok.Position())
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:143
		{
			yyVAL.type_specifier = createTypeSpecifier(vm.BooleanType, yyDollar[1].tok.Position())
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:147
		{
			yyVAL.type_specifier = createTypeSpecifier(vm.IntType, yyDollar[1].tok.Position())
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:151
		{
			yyVAL.type_specifier = createTypeSpecifier(vm.DoubleType, yyDollar[1].tok.Position())
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:155
		{
			yyVAL.type_specifier = createTypeSpecifier(vm.StringType, yyDollar[1].tok.Position())
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:161
		{
			yyVAL.type_specifier = createClassTypeSpecifier(yyDollar[1].tok.Lit, yyDollar[1].tok.Position())
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:167
		{
			yyVAL.type_specifier = createArrayTypeSpecifier(yyDollar[1].type_specifier)
			yyVAL.type_specifier.SetPosition(yyDollar[1].type_specifier.Position())
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:172
		{
			class_type := createClassTypeSpecifier(yyDollar[1].tok.Lit, yyDollar[1].tok.Position())
			yyVAL.type_specifier = createArrayTypeSpecifier(class_type)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:177
		{
			yyVAL.type_specifier = createArrayTypeSpecifier(yyDollar[1].type_specifier)
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:183
		{
			yyVAL.type_specifier = yyDollar[1].type_specifier
		}
	case 25:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:191
		{
			l := yylex.(*Lexer)
			l.compiler.functionDefine(yyDollar[1].type_specifier, yyDollar[2].tok.Lit, yyDollar[4].parameter_list, yyDollar[6].block)
		}
	case 26:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:196
		{
			l := yylex.(*Lexer)
			l.compiler.functionDefine(yyDollar[1].type_specifier, yyDollar[2].tok.Lit, []*Parameter{}, yyDollar[5].block)
		}
	case 27:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:201
		{
			l := yylex.(*Lexer)
			l.compiler.functionDefine(yyDollar[1].type_specifier, yyDollar[2].tok.Lit, yyDollar[4].parameter_list, nil)
		}
	case 28:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:206
		{
			l := yylex.(*Lexer)
			l.compiler.functionDefine(yyDollar[1].type_specifier, yyDollar[2].tok.Lit, []*Parameter{}, nil)
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:213
		{
			parameter := &Parameter{typeSpecifier: yyDollar[1].type_specifier, name: yyDollar[2].tok.Lit}
			yyVAL.parameter_list = []*Parameter{parameter}
		}
	case 30:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:218
		{
			yyVAL.parameter_list = append(yyDollar[1].parameter_list, &Parameter{typeSpecifier: yyDollar[3].type_specifier, name: yyDollar[4].tok.Lit})
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:224
		{
			yyVAL.argument_list = []Expression{yyDollar[1].expression}
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:228
		{
			yyVAL.argument_list = append(yyDollar[1].argument_list, yyDollar[3].expression)
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:234
		{
			yyVAL.statement_list = []Statement{yyDollar[1].statement}
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:238
		{
			yyVAL.statement_list = append(yyDollar[1].statement_list, yyDollar[2].statement)
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:245
		{
			yyVAL.expression = &CommaExpression{left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:253
		{
			yyVAL.expression = &AssignExpression{left: yyDollar[1].expression, operand: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:261
		{
			yyVAL.expression = &BinaryExpression{operator: LogicalOrOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:269
		{
			yyVAL.expression = &BinaryExpression{operator: LogicalAndOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:277
		{
			yyVAL.expression = &BinaryExpression{operator: EqOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:282
		{
			yyVAL.expression = &BinaryExpression{operator: NeOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:290
		{
			yyVAL.expression = &BinaryExpression{operator: GtOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:295
		{
			yyVAL.expression = &BinaryExpression{operator: GeOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:300
		{
			yyVAL.expression = &BinaryExpression{operator: LtOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:305
		{
			yyVAL.expression = &BinaryExpression{operator: LeOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:313
		{
			yyVAL.expression = &BinaryExpression{operator: AddOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:318
		{
			yyVAL.expression = &BinaryExpression{operator: SubOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:326
		{
			yyVAL.expression = &BinaryExpression{operator: MulOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:331
		{
			yyVAL.expression = &BinaryExpression{operator: DivOperator, left: yyDollar[1].expression, right: yyDollar[3].expression}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:339
		{
			yyVAL.expression = &MinusExpression{operand: yyDollar[2].expression}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:344
		{
			yyVAL.expression = &LogicalNotExpression{operand: yyDollar[2].expression}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:356
		{
			yyVAL.expression = createIdentifierExpression(yyDollar[1].tok.Lit, yyDollar[1].tok.Position())
		}
	case 64:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:362
		{
			yyVAL.expression = createIndexExpression(yyDollar[1].expression, yyDollar[3].expression, yyDollar[1].expression.Position())
		}
	case 65:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:366
		{
			identifier := createIdentifierExpression(yyDollar[1].tok.Lit, yyDollar[1].tok.Position())
			yyVAL.expression = createIndexExpression(identifier, yyDollar[3].expression, yyDollar[1].tok.Position())
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:371
		{
			yyVAL.expression = createMemberExpression(yyDollar[1].expression, yyDollar[3].tok.Lit)
		}
	case 67:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:375
		{
			yyVAL.expression = &FunctionCallExpression{function: yyDollar[1].expression, argumentList: yyDollar[3].argument_list}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:380
		{
			yyVAL.expression = &FunctionCallExpression{function: yyDollar[1].expression, argumentList: []Expression{}}
			yyVAL.expression.SetPosition(yyDollar[1].expression.Position())
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:385
		{
			yyVAL.expression = yyDollar[2].expression
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:389
		{
			value, _ := strconv.Atoi(yyDollar[1].tok.Lit)
			yyVAL.expression = &IntExpression{intValue: value}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:395
		{
			value, _ := strconv.ParseFloat(yyDollar[1].tok.Lit, 64)
			yyVAL.expression = &DoubleExpression{doubleValue: value}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:401
		{
			yyVAL.expression = &StringExpression{stringValue: yyDollar[1].tok.Lit}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:406
		{
			yyVAL.expression = &BooleanExpression{booleanValue: true}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:411
		{
			yyVAL.expression = &BooleanExpression{booleanValue: false}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:416
		{
			yyVAL.expression = &NullExpression{}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:422
		{
			yyVAL.expression = createThisExpression(yyDollar[1].tok.Position())
		}
	case 78:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:426
		{
			yyVAL.expression = createNewExpression(yyDollar[2].class_name, nil, yyDollar[1].tok.Position())
		}
	case 79:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:430
		{
			yyVAL.expression = createNewExpression(yyDollar[2].class_name, yyDollar[4].argument_list, yyDollar[1].tok.Position())
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:436
		{
			yyVAL.class_name = []string{yyDollar[1].tok.Lit}
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:440
		{
			yyVAL.class_name = append(yyDollar[1].class_name, yyDollar[3].tok.Lit)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:446
		{
			yyVAL.expression = &ArrayLiteralExpression{arrayLiteral: yyDollar[2].expression_list}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 83:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:451
		{
			yyVAL.expression = &ArrayLiteralExpression{arrayLiteral: yyDollar[2].expression_list}
			yyVAL.expression.SetPosition(yyDollar[1].tok.Position())
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:458
		{
			yyVAL.expression = createBasicArrayCreation(yyDollar[2].type_specifier, yyDollar[3].array_dimension_list, nil, yyDollar[1].tok.Position())
		}
	case 85:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:462
		{
			yyVAL.expression = createBasicArrayCreation(yyDollar[2].type_specifier, yyDollar[3].array_dimension_list, yyDollar[4].array_dimension_list, yyDollar[1].tok.Position())
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:466
		{
			yyVAL.expression = createClassArrayCreation(yyDollar[2].type_specifier, yyDollar[3].array_dimension_list, nil, yyDollar[1].tok.Position())
		}
	case 87:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:470
		{
			yyVAL.expression = createClassArrayCreation(yyDollar[2].type_specifier, yyDollar[3].array_dimension_list, yyDollar[4].array_dimension_list, yyDollar[1].tok.Position())
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:476
		{
			yyVAL.array_dimension_list = []*ArrayDimension{yyDollar[1].array_dimension}
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:480
		{
			yyVAL.array_dimension_list = append(yyDollar[1].array_dimension_list, yyDollar[2].array_dimension)
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:486
		{
			yyVAL.array_dimension = &ArrayDimension{expression: yyDollar[2].expression}
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:492
		{
			yyVAL.array_dimension_list = []*ArrayDimension{&ArrayDimension{}}
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:496
		{
			yyVAL.array_dimension_list = append(yyDollar[1].array_dimension_list, &ArrayDimension{})
		}
	case 93:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:502
		{
			yyVAL.expression_list = nil
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:506
		{
			yyVAL.expression_list = []Expression{yyDollar[1].expression}
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:510
		{
			yyVAL.expression_list = append(yyDollar[1].expression_list, yyDollar[3].expression)
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:516
		{
			yyVAL.statement = &ExpressionStatement{expression: yyDollar[1].expression}
			yyVAL.statement.SetPosition(yyDollar[1].expression.Position())
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:529
		{
			yyVAL.statement = &IfStatement{condition: yyDollar[2].expression, thenBlock: yyDollar[3].block, elifList: []*Elif{}, elseBlock: nil}
			yyVAL.statement.SetPosition(yyDollar[1].tok.Position())
		}
	case 104:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:534
		{
			yyVAL.statement = &IfStatement{condition: yyDollar[2].expression, thenBlock: yyDollar[3].block, elifList: []*Elif{}, elseBlock: yyDollar[5].block}
			yyVAL.statement.SetPosition(yyDollar[1].tok.Position())
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:539
		{
			yyVAL.statement = &IfStatement{condition: yyDollar[2].expression, thenBlock: yyDollar[3].block, elifList: yyDollar[4].elif_list, elseBlock: nil}
			yyVAL.statement.SetPosition(yyDollar[1].tok.Position())
		}
	case 106:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:544
		{
			yyVAL.statement = &IfStatement{condition: yyDollar[2].expression, thenBlock: yyDollar[3].block, elifList: yyDollar[4].elif_list, elseBlock: yyDollar[6].block}
			yyVAL.statement.SetPosition(yyDollar[1].tok.Position())
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:551
		{
			yyVAL.elif_list = []*Elif{&Elif{condition: yyDollar[2].expression, block: yyDollar[3].block}}
		}
	case 108:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:555
		{
			yyVAL.elif_list = append(yyDollar[1].elif_list, &Elif{condition: yyDollar[3].expression, block: yyDollar[4].block})
		}
	case 109:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.go.y:561
		{
			yyVAL.statement = &ForStatement{init: yyDollar[3].expression, condition: yyDollar[5].expression, post: yyDollar[7].expression, block: yyDollar[9].block}
			yyVAL.statement.SetPosition(yyDollar[1].tok.Position())
			yyDollar[9].block.parent = &StatementBlockInfo{statement: yyVAL.statement}
		}
	case 110:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:569
		{
			yyVAL.expression = nil
		}
	case 112:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:576
		{
			yyVAL.statement = &ReturnStatement{returnValue: yyDollar[2].expression}
			yyVAL.statement.SetPosition(yyDollar[1].tok.Position())
		}
	case 113:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:583
		{
			yyVAL.statement = &BreakStatement{}
			yyVAL.statement.SetPosition(yyDollar[1].tok.Position())
		}
	case 114:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:590
		{
			yyVAL.statement = &ContinueStatement{}
			yyVAL.statement.SetPosition(yyDollar[1].tok.Position())
		}
	case 115:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:597
		{
			yyVAL.statement = &Declaration{typeSpecifier: yyDollar[1].type_specifier, name: yyDollar[2].tok.Lit, variableIndex: -1}
			yyVAL.statement.SetPosition(yyDollar[1].type_specifier.Position())
		}
	case 116:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:602
		{
			yyVAL.statement = &Declaration{typeSpecifier: yyDollar[1].type_specifier, name: yyDollar[2].tok.Lit, initializer: yyDollar[4].expression, variableIndex: -1}
			yyVAL.statement.SetPosition(yyDollar[1].type_specifier.Position())
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:609
		{
			l := yylex.(*Lexer)
			l.compiler.currentBlock = &Block{outerBlock: l.compiler.currentBlock}
			yyVAL.block = l.compiler.currentBlock
		}
	case 118:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:615
		{
			currentBlock := yyDollar[2].block
			currentBlock.statementList = yyDollar[3].statement_list

			l := yylex.(*Lexer)

			yyVAL.block = l.compiler.currentBlock
			l.compiler.currentBlock = currentBlock.outerBlock
		}
	case 119:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:625
		{
			l := yylex.(*Lexer)
			yyVAL.block = &Block{outerBlock: l.compiler.currentBlock}
		}
	case 120:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:632
		{
			startClassDefine(yyDollar[2].tok.Lit, yyDollar[3].extends_list, yyDollar[1].tok.Position())
		}
	case 121:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:636
		{
			endClassDefine(yyDollar[6].member_declaration)
		}
	case 122:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:640
		{
			startClassDefine(yyDollar[2].tok.Lit, yyDollar[3].extends_list, yyDollar[1].tok.Position())
		}
	case 123:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:644
		{
			endClassDefine(nil)
		}
	case 124:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:650
		{
			yyVAL.extends_list = nil
		}
	case 125:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:654
		{
			yyVAL.extends_list = yyDollar[2].extends_list
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:660
		{
			yyVAL.extends_list = createExtendList(yyDollar[1].tok.Lit)
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:664
		{
			yyVAL.extends_list = chainExtendList(yyDollar[1].extends_list, yyDollar[3].tok.Lit)
		}
	case 129:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:671
		{
			yyVAL.member_declaration = chainMemberDeclaration(yyDollar[1].member_declaration, yyDollar[2].member_declaration)
		}
	case 132:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:681
		{
			yyVAL.member_declaration = createMethodMember(yyDollar[1].function_definition, yyDollar[1].function_definition.typeSpecifier.Position())
		}
	case 133:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:687
		{
			yyVAL.function_definition = methodFunctionDefine(yyDollar[1].type_specifier, yyDollar[2].tok.Lit, yyDollar[4].parameter_list, yyDollar[6].block)
		}
	case 134:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:691
		{
			yyVAL.function_definition = methodFunctionDefine(yyDollar[1].type_specifier, yyDollar[2].tok.Lit, nil, yyDollar[5].block)
		}
	case 135:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:695
		{
			yyVAL.function_definition = methodFunctionDefine(yyDollar[1].type_specifier, yyDollar[2].tok.Lit, yyDollar[4].parameter_list, nil)
		}
	case 136:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:699
		{
			yyVAL.function_definition = methodFunctionDefine(yyDollar[1].type_specifier, yyDollar[2].tok.Lit, nil, nil)
		}
	case 137:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:705
		{
			yyVAL.member_declaration = createFieldMember(yyDollar[1].type_specifier, yyDollar[2].tok.Lit, yyDollar[1].type_specifier.Position())
		}
	}
	goto yystack /* stack new state and value */
}
