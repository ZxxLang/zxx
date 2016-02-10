// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import "strconv"

type Token int

const (
	ILLEGAL Token = iota
	EOF

	SPACES
	INDENT

	PLACEHOLDER
	COMMENT

	LITERAL // 字面值

	operator_beg
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
	operator_end

	// 自成语句的符号
	ASSIGN // =
	INC    // ++
	DEC    // --

	// 成对符号
	LEFT  // [{('"
	RIGHT // ]})'"

	// 定界
	COLON     // :
	COMMA     // ,
	SEMICOLON // ;

	DOT // .

	keyword_beg
	// 关键字
	BREAK
	CASE
	CONST
	CONTINUE

	DEFAULT
	DEFER
	ELSE
	FOR

	FUNC
	GO
	GOTO
	IF

	MAP

	OUT
	PROC
	PUB
	SWITCH
	STATIC
	TYPE
	USE
	VAR

	keyword_end

	// 标识符
	IDENT
)

var tokens = [...]string{
	ILLEGAL:     "ILLEGAL",
	EOF:         "EOF",
	PLACEHOLDER: "PLACEHOLDER",
	INDENT:      "INDENT",
	SPACES:      "SPACES",
	COMMENT:     "COMMENT",
	LITERAL:     "LITERAL",
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

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup 识别 ident 对应的 Token 值
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

// IsLiteral 返回 tok 是否为字面值标记
func (tok Token) IsLiteral() bool { return LITERAL == tok }

// IsOperator 返回 tok 是否为操作符
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword 返回 tok 是否为关键字
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }
