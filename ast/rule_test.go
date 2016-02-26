package ast

import (
	"testing"

	"github.com/ZxxLang/zxx/scanner"
	"github.com/ZxxLang/zxx/token"
)

var qPub queue

func init() {
	qPub = queue{
		rules: []rule{
			term{token.PUB},
			term{token.VAR},
			term{token.Type, token.IDENT},
			term{token.IDENT, token.MEMBER},
		},
	}
}

func TestRule(t *testing.T) {
	tok := testParse("pub var int a.b", &qPub)
	if tok != token.EOF {
		t.Fatal(qPub.at, tok)
	}
}

func testParse(src string, q rule) (tok token.Token) {
	var code string
	scan := scanner.New([]byte(src))
	for {
		code, _ = scan.Symbol()
		tok = token.Lookup(code)
		if tok == token.SPACES {
			continue
		}

		if tok == token.EOF {
			return
		}
		if tok == token.PLACEHOLDER {
			tok = lookup(code)
		}

		_, ok := q.pass(tok)
		if !ok {
			return
		}
	}
}

func lookup(code string) (tok token.Token) {
	tok = token.PLACEHOLDER
	// 整数, 浮点数, datetime
	// ??? 缺少严格检查
	if code[0] >= '0' && code[0] <= '9' {
		tok = token.VALINTEGER
		if code[0] == '0' && len(code) > 2 && (code[1] == 'x' || code[1] == 'b') {
		} else {
			for _, c := range code {
				if c == '.' || c == 'e' {
					tok = token.VALFLOAT
				} else if c == 'T' || c == ':' || c == 'Z' {
					tok = token.VALDATETIME
				} else if (c < '0' || c > '9') && c != '+' && c != '-' && c != '_' {
					tok = token.PLACEHOLDER
					break
				}
			}
		}
	} else {
		// 标识符, 成员
		tok = token.IDENT
		dot := 0
		for _, c := range code {
			if c == '.' {
				dot++
				continue
			}

			if c != '_' && !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
				tok = token.PLACEHOLDER
				break
			}
		}
		if tok == token.IDENT && dot != 0 {
			if dot == 1 {
				tok = token.MEMBER
			} else {
				tok = token.MEMBERS
			}
		}
	}
	return
}
