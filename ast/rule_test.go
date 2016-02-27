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

	z_mnl := Zero{Term{token.NL, token.EMPTYLINE}}
	t_comma := Term{token.NL, token.COMMA}
	t_semi := Term{token.NL, token.COMMA}

	r_type := Seq(
		Term{token.TYPE},
		Term{token.IDENT},
	)

	r_decl_use = Seq(
		Term{token.USE},
		Any(
			Seq(
				Zero{Term{token.IDENT}},
				Term{token.VALSTRING},
			),
			Seq(
				Term{token.LEFT},
				More(Seq(
					z_mnl,
					Zero{Term{token.IDENT}},
					Term{token.VALSTRING},
					Zero{t_comma},
				)),
				Term{token.RIGHT},
			),
		),
	)

	r_decl_type = Seq(
		Zero{Term{token.PUB}},
		Any(
			r_type,
			Seq(
				Term{token.LEFT},
				More(Seq(z_mnl, r_type, t_comma)),
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
		nsr(4, `use a 'b' 'd'`),
		nsr(5, `use a 'b' c 'd'`),
		nsr(4, `use ('b')`),
		nsr(5, `use (a 'b')`),
		nsr(5, `use ( 'b',)`),
		nsr(7, `use (a 'b' c 'd')`),
		nsr(7, "use (\na 'b'\n)"),
		nsr(8, "use (\n\ta 'b',\n\t)"),
		nsr(7, `use ( 'b', c 'd')`),
		nsr(7, "use ( 'b'\nc 'd')"),
		nsr(8, "use ( 'b',\n\tc 'd')"),
		nsr(10, "use (\n\n\na 'b',\n\n\tc 'd/e'\n\n\t)"),
	}

	bad = []ns{
		nsr(1, `use`, r_decl_use),
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
			t.Fatal(i, n, ns.s, tok)
		}
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
	parser.Fast(
		[]byte(src),
		func(_ scanner.Pos, t token.Token, _ string) error {
			tok = t

			if tok == token.EOF || tok > token.IDENT {
				return nil
			}

			eat, clear := r.Eat(tok)

			if eat == 0 && !clear {
				return io.ErrNoProgress
			}
			n++
			return nil
		})

	return
}
