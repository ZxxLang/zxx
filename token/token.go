// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import "strconv"

type Token int

const (
	EOF Token = iota
	// 运算符
	DOLLAR // $

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

	// 关键字

	USE
	PUB
	CONST
	TYPE
	STATIC
	VAR
	FUNC
	PROC

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

	MAP

	SWITCH

	OUT

	// 定界
	COLON     // :
	COMMA     // ,
	SEMICOLON // ;

	DOT // .

	// 自成语句的符号
	ASSIGN // =
	INC    // ++
	DEC    // --

	// 成对符号
	LEFT  // [{(
	RIGHT // ]})

	SPACES // 连续的空格
	TABS   // 连续的制表符

	// 需要进行语义识别, 转化

	LITERAL // 字面值
	IDENT   // 标识符

	PLACEHOLDER // 占位

	// 注释
	COMMENT  // '//'
	COMMENTS // '---'

	// 换行
	NL
)

var tokens = [...]string{
	EOF:         "EOF",
	LITERAL:     "LITERAL",
	PLACEHOLDER: "PLACEHOLDER",
	TABS:        "TABS",
	SPACES:      "SPACES",
	COMMENT:     "COMMENT",
	COMMENTS:    "COMMENTS",
	NL:          "NEWLINE",
	DOLLAR:      "$",
	ANTI:        "~",
	BITAND:      "&",
	BITOR:       "|",
	XOR:         "xor",
	MUL:         "mul",
	MULSIGN:     "*",
	DIV:         "div",
	DIVSIGN:     "/",
	MOD:         "mod",
	REM:         "rem",
	SHL:         "shl",
	SHLSIGN:     "<<",
	SHR:         "shr",
	SHRSIGN:     ">>",
	ADD:         "add",
	PLUS:        "+",
	SUB:         "-",
	DOTDOT:      "..",
	EQL:         "==",
	NEQ:         "!=",
	LEQ:         "<=",
	GEQ:         ">=",
	LSS:         "<",
	GTR:         ">",
	IS:          "is",
	ISNOT:       "isnot",
	HAS:         "has",
	NOT:         "not",
	AND:         "and",
	OR:          "or",
	ASSIGN:      "=",
	INC:         "++",
	DEC:         "--",
	LEFT:        "LEFT",
	RIGHT:       "RIGHT",
	COLON:       ":",
	COMMA:       ",",
	SEMICOLON:   ";",
	DOT:         ".",
	BREAK:       "break",
	CASE:        "case",
	CONST:       "const",
	CONTINUE:    "continue",
	DEFAULT:     "default",
	DEFER:       "defer",
	ELSE:        "else",
	FOR:         "for",
	FUNC:        "func",
	GO:          "go",
	GOTO:        "goto",
	IF:          "if",
	MAP:         "map",
	OUT:         "out",
	PROC:        "proc",
	PUB:         "pub",
	SWITCH:      "switch",
	STATIC:      "static",
	TYPE:        "type",
	USE:         "use",
	VAR:         "var",
	IDENT:       "IDENT",
}

// String 返回 tok 名称或者字面值.
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
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
	case DOLLAR:
		return 11
	}

	return 0
}

var letters map[string]Token // fixed

func init() {
	letters = make(map[string]Token)
	for i := EOF + 1; i < LEFT; i++ {
		letters[tokens[i]] = i
	}
	letters["["] = LEFT
	letters["{"] = LEFT
	letters["("] = LEFT
	letters["]"] = RIGHT
	letters[")"] = RIGHT
	letters["}"] = RIGHT
}

// Lookup 识别字面值 letter 对应的 Token.
// 适当的调用 Lookup 才有意义.
func Lookup(letter string) Token {
	if letter == "" {
		return EOF
	}

	if letter[0] == ' ' {
		return SPACES
	}

	if letter[0] == '\t' {
		return TABS
	}

	if tok, is := letters[letter]; is {
		return tok
	}

	if letter[0] == '\n' || letter[1] == '\r' {
		return NL
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

	return LITERAL
}

// IsAssign 返回 tok 是否为赋值语句的符号 '=','++','--'
func (tok Token) IsAssign() bool {
	return tok == ASSIGN || tok == INC || tok == DEC
}

// IsOperator 返回 tok 是否为操作符
func (tok Token) IsOperator() bool { return tok >= DOLLAR && tok <= OR }

// IsKeyword 返回 tok 是否为关键字
func (tok Token) IsKeyword() bool { return tok >= USE && tok <= OUT }

// IsDeclare 返回 tok 是否为声明
func (tok Token) IsDeclare() bool { return tok >= USE && tok <= PROC }
