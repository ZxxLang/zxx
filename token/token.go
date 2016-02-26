// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import "strconv"

type Token int

const (
	EOF Token = iota

	Operator // Token 分类标记
	// 运算符

	//DOLLAR // $

	ANTI // ~

	BITAND // &
	BITOR  // |
	XOR    // xor

	MUL     // mul
	MULSIGN // *
	DIV     // div
	DIVSIGN // /
	MOD     // mod
	REM     // %
	SHL     // shl
	SHLSIGN // <<
	SHR     // shr
	SHRSIGN // >>

	ADD  // add
	PLUS // +
	SUB  // -

	DOTDOT // ..

	EQL // ==
	NEQ // !=
	LEQ // <=
	GEQ // >=
	LSS // <
	GTR // >

	IS    // is
	ISNOT // isnot
	HAS   // has

	NOT // not

	AND // and

	OR // or

	Assign // 分类标记
	// 赋值语句
	ASSIGN // =
	INC    // ++
	DEC    // --

	Declare // 分类标记
	// 声明
	USE
	PUB
	CONST
	TYPE
	STATIC
	VAR
	FUNC
	PROC

	Statement // 分类标记
	// 语句
	BREAK
	CASE
	CONTINUE
	DEFAULT
	DEFER
	ELSE
	FOR
	GO
	GOTO
	IF
	SWITCH
	OUT // 具有二义性

	Divide // 分类标记
	// 分界
	COLON     // :
	COMMA     // ,
	SEMICOLON // ;

	DOT // .

	Type // 分类标记
	// 预定义类型
	BYTE
	STRING
	UINT
	INT
	U8
	U16
	U32
	U64
	I8
	I16
	I32
	I64
	F32
	F64
	F128
	DATETIME
	BOOL
	MAP
	ARRAY

	Literal // 分类标记
	// 字面值
	NULL
	// 风格值, 由解析器分析 PLACEHOLDER 转义得到
	VALSTRING
	VALINTEGER
	VALFLOAT
	VALDATETIME
	VALBOOL

	Alone // 分类标记
	// 下列每个都是单独的

	// 成对符号
	LEFT  // [{(
	RIGHT // ]})

	// 由解析器分析转义得到
	IDENT
	EMPTYLINE
	INDENTATION
	MEMBER  // IDENT 包含一个 '.'
	MEMBERS // IDENT 包含多个 '.'
	SUGAR

	// 注释
	COMMENT // '//'

	// 换行
	NL

	// 占位
	PLACEHOLDER

	// AST 中是没有下列 Token, 他们被转义了

	NAN
	INFINITE
	TRUE
	FALSE

	SPACES   // 连续的空格
	TABS     // 连续的制表符
	COMMENTS // '---'
)

var tokens = [...]string{
	EOF:         "EOF",
	IDENT:       "IDENT",
	LEFT:        "LEFT",
	NL:          "NEWLINE",
	EMPTYLINE:   "EMPTYLINE",
	PLACEHOLDER: "PLACEHOLDER",
	RIGHT:       "RIGHT",
	SPACES:      "SPACES",
	TABS:        "TABS",

	COMMENT:  "//",
	COMMENTS: "---",

	//DOLLAR:  "$",
	ANTI:    "~",
	BITAND:  "&",
	BITOR:   "|",
	XOR:     "xor",
	MUL:     "mul",
	MULSIGN: "*",
	DIV:     "div",
	DIVSIGN: "/",
	MOD:     "mod",
	REM:     "rem",
	SHL:     "shl",
	SHLSIGN: "<<",
	SHR:     "shr",
	SHRSIGN: ">>",
	ADD:     "add",
	PLUS:    "+",
	SUB:     "-",
	DOTDOT:  "..",
	EQL:     "==",
	NEQ:     "!=",
	LEQ:     "<=",
	GEQ:     ">=",
	LSS:     "<",
	GTR:     ">",
	IS:      "is",
	ISNOT:   "isnot",
	HAS:     "has",
	NOT:     "not",
	AND:     "and",
	OR:      "or",

	ASSIGN: "=",
	INC:    "++",
	DEC:    "--",

	COLON:     ":",
	COMMA:     ",",
	SEMICOLON: ";",
	DOT:       ".",

	BREAK:    "break",
	CASE:     "case",
	CONST:    "const",
	CONTINUE: "continue",
	DEFAULT:  "default",
	DEFER:    "defer",
	ELSE:     "else",
	FOR:      "for",
	FUNC:     "func",
	GO:       "go",
	GOTO:     "goto",
	IF:       "if",
	OUT:      "out",
	PROC:     "proc",
	PUB:      "pub",
	STATIC:   "static",
	SWITCH:   "switch",
	TYPE:     "type",
	USE:      "use",
	VAR:      "var",

	NULL:     "null",
	NAN:      "nan",
	INFINITE: "infinite",
	FALSE:    "false",
	TRUE:     "true",

	BYTE:     "byte",
	STRING:   "string",
	UINT:     "uint",
	INT:      "int",
	U8:       "u8",
	U16:      "u16",
	U32:      "u32",
	U64:      "u64",
	I8:       "i8",
	I16:      "i16",
	I32:      "i32",
	I64:      "i64",
	F32:      "f32",
	F64:      "f64",
	F128:     "f128",
	DATETIME: "datetime",
	BOOL:     "bool",
	MAP:      "map",
	ARRAY:    "array",
}

// String 返回 tok 名称或者字面值.
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "Token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// Precedence 返回运算符 op 的优先级
func (op Token) Precedence() int {
	switch op {
	case OR:
		return 1
	case AND:
		return 2
	case NOT:
		return 3
	case IS, ISNOT, HAS:
		return 4
	case EQL, NEQ, LEQ, GEQ, LSS, GTR:
		return 5
	case DOTDOT:
		return 6
	case ADD, PLUS, SUB:
		return 7
	case MUL, MULSIGN, DIV, DIVSIGN, MOD, REM, SHL, SHLSIGN, SHR, SHRSIGN:
		return 8
	case BITAND, BITOR, XOR:
		return 9
	case ANTI:
		return 10
	}
	//case DOLLAR:
	//	return 11

	return 0
}

var letters map[string]Token // fixed

func init() {
	letters = map[string]Token{
		"[": LEFT,
		"{": LEFT,
		"(": LEFT,
		"]": RIGHT,
		")": RIGHT,
		"}": RIGHT,
	}

	for i := ANTI; i < Alone; i++ {
		if i == Assign || i == Declare || i == Statement || i == Literal || i == Divide {
			continue
		}
		letters[tokens[i]] = i
	}
}

// Lookup 分析 letter 的前缀字符串, 猜测它相应的 Token.
// 难于猜测的 letter 被判定为 PLACEHOLDER.
// 返回值包括:
// 	'$' 之外的运算符, 保留字
// 	EOF,SPACES,TABS,NL,COMMENT,COMMENTS
// 	LEFT,RIGHT,
// 	COLON,COMMA,SEMICOLON,DOT,ASSIGN,INC,DEC
// 	PLACEHOLDER
//
// 特别的, 多字节开头被当做 COMMENT
func Lookup(letter string) Token {
	if letter == "" {
		return EOF
	}
	// 多字节
	if letter[0] > 127 {
		return COMMENT
	}

	switch letter[0] {
	case ' ':
		return SPACES
	case '\t':
		return TABS
	case '\n', '\r':
		return NL
	}

	if tok, is := letters[letter]; is {
		return tok
	}

	if len(letter) > 1 {
		if letter[0] == '/' && letter[1] == '/' {
			return COMMENT
		}
	}

	if len(letter) > 2 {
		if letter[0] == '-' && letter[1] == '-' && letter[2] == '-' {
			return COMMENTS
		}
	}

	return PLACEHOLDER
}

// As 返回 token 属于 tok
func (token Token) As(tok Token) bool {
	return tok.Has(token)
}

// Has 返回 token 包括 tok
func (token Token) Has(tok Token) bool {
	if tok == token {
		return true
	}
	switch token {
	case Operator:
		return tok > Operator && tok < Assign
	case Assign:
		return tok > Assign && tok < Declare
	case Declare:
		return tok > Declare && tok < Statement
	case Statement:
		return tok > Statement && tok < Divide
	case Divide:
		return tok > Divide && tok < Type
	case Type:
		return tok > Type && tok < Literal
	case Literal:
		return tok > Literal && tok < Alone
	case Alone:
		return tok > Alone
	}
	return false
}
