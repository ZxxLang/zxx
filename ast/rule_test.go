package ast_test

import (
	"io"
	"testing"

	. "github.com/ZxxLang/zxx/ast"
	"github.com/ZxxLang/zxx/parser"
	"github.com/ZxxLang/zxx/scanner"
	"github.com/ZxxLang/zxx/token"
)

var (
	r_decl_use, r_decl_type,
	r_ Rule
)

func init() {
	// r_pub    Rule = Term{token.PUB}
	// r_static      = Term{token.STATIC}
	// r_var         = Term{token.VAR}
	// r_type        = Term{token.TYPE}
	// r_const       = Term{token.CONST}
	// r_func        = Term{token.FUNC}
	// r_proc        = Term{token.PROC}
	// r_modif       = Zero{Any{
	// 	Seq(r_pub, r_static),
	// 	Seq(r_static),
	// }}
	// r_left  = Term{token.LEFT}
	// r_right = Term{token.RIGHT}

	z_nl := Zero{Term{token.NL}}
	t_comma := Term{token.NL, token.COMMA}
	t_semi := Term{token.NL, token.SEMICOLON}
	t_ident := Term{token.IDENT}

	r_type := Seq(
		Term{token.TYPE},
		Term{token.IDENT},
	)

	r_decl_use = Seq(
		Term{token.USE},
		Any(
			Seq(
				Term{token.LEFT},
				More(Seq(
					z_nl,
					Zero{t_ident},
					Term{token.VALSTRING},
					Zero{t_comma},
				)),
				Term{token.RIGHT},
			),
			Seq(
				Zero{t_ident},
				Term{token.VALSTRING},
			),
		),
	)

	r_decl_type = Seq(
		Zero{Term{token.PUB}},
		Any(
			r_type,
			Seq(
				Term{token.LEFT},
				More(Seq(z_nl, r_type, t_comma)),
				Term{token.RIGHT},
			),
			t_semi,
		),
	)
}

type ns struct {
	n int
	s string
	r Rule
}

func nsr(n int, s string, r ...Rule) ns {
	if r == nil {
		return ns{n, s, nil}
	}
	return ns{n, s, r[0]}
}

var good, bad []ns

func init() {
	good = []ns{
		nsr(2, `use 'b'`, r_decl_use),
		nsr(3, `use a 'b'`),
		nsr(4, `use ('b')`),
		nsr(5, `use (a 'b')`),
		nsr(5, `use ( 'b',)`),
		nsr(8, `use (a 'b', c 'd')`),
		nsr(7, "use (\na 'b'\n)"),
		nsr(8, "use (\n\ta 'b',\n\t)"),
		nsr(7, `use ( 'b', c 'd')`),
		nsr(8, "use ( 'b',\nc 'd')"),
		nsr(8, "use ( 'b',\n\tc 'd')"),
		nsr(11, "use (世界\n\n你好\na 'b',\n\n\tc 'd/e'\n\n\t)", r_decl_use),
	}

	bad = []ns{
		nsr(1, `use`, r_decl_use),
		nsr(2, `use a`),
		nsr(2, `use (a,'d')`),
		nsr(5, `use a 'b' c 'd'`),
		nsr(6, `use a 'b', c 'd'`),
		nsr(7, "use ( 'b'\n, c 'd')"),
	}
}

func Test_good(t *testing.T) {
	var r Rule
	for i, ns := range good {
		if ns.r != nil {
			r = ns.r
		}
		n, tok := testParse(ns.s, r)
		if n != ns.n || tok != token.EOF || !r.Ok() {
			t.Fatal(i, n, ns.s, tok, r.Ok())
		}
		r.Reset()
	}
}

func Test_bad(t *testing.T) {
	var r Rule
	for i, ns := range bad {
		if ns.r != nil {
			r = ns.r
		}
		n, tok := testParse(ns.s, r)
		if r.Ok() {
			t.Fatal(i, n, ns.s, tok, r.Ok())
		}
		r.Reset()
	}
}

func testParse(src string, r Rule) (n int, tok token.Token) {
	nl := false
	parser.Fast(
		[]byte(src),
		func(_ scanner.Pos, t token.Token, _ string) error {
			tok = t

			// 去占位节点
			if tok == token.EOF || tok > token.IDENT {
				return nil
			}

			if nl && tok == token.NL { // 合并连续的换行
				return nil
			}
			nl = tok == token.NL

			eat, _ := r.Eat(tok)

			if eat != 1 {
				return io.ErrNoProgress
			}
			n++
			return nil
		})

	return
}
