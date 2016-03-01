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

	o_nl := Option(Term{token.NL})
	a_ident := Term{token.IDENT}

	a_type := Term{token.IDENT, token.Type}

	expr := Seq()
	expr0 := More(Term{token.COMMA, token.Operator}, Any(
		Term{token.Literal, token.IDENT},
		expr,
	))
	expr = (Any(
		Option(Seq(
			Term{token.Literal, token.IDENT},
			Option(Any(
				Seq(
					Term{token.Operator},
					expr0,
				),
				Seq(Term{token.LEFT}, expr0, Term{token.RIGHT}),
			)),
		)),
		Seq(
			Term{token.NOT, token.ANTI},
			expr0,
		),
		Seq(Term{token.LEFT}, expr0, Term{token.RIGHT}),
	))

	s_assign := Seq(
		Term{token.ASSIGN, token.COLON},
		expr,
	)

	r_decl_use = Seq(
		Term{token.USE},
		Any(
			Seq(
				Term{token.LEFT},
				More(Term{token.COMMA}, Seq(
					o_nl,
					Option(Term{token.IDENT}),
					Term{token.VALSTRING},
					o_nl,
				)),
				Term{token.RIGHT},
			),
			Seq(Option(a_ident), Term{token.VALSTRING}),
		),
	)

	r_decl_type = Seq(
		Term{token.TYPE}, Term{token.IDENT},
		Any(
			Term{token.IDENT},
			Seq(
				Term{token.LEFT},
				Any(
					Term{token.RIGHT},
					Seq(
						More(Term{token.COMMA}, Seq(
							o_nl,
							a_type, More(
								Term{token.COMMA},
								Seq(
									Term{token.IDENT},
									Option(s_assign),
								),
							),
						)),
						Term{token.RIGHT},
					),
				),
			),
		),
	)

	// r_ 临时测试用
	r_ = Seq(
		Term{token.USE},
		Term{token.LEFT},
		More(Term{token.COMMA}, Seq(
			Term{token.IDENT},
			Option(s_assign),
		)),
		Term{token.RIGHT},
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
		nsr(8, `use (a, b = "")`, r_),
		nsr(10, `use (a, b = [true])`),
		nsr(10, `use (a, b = [11])`),
		nsr(12, `use (a, b = [true,true])`),
		nsr(12, `use (a, b = [true,"1"])`),
		nsr(13, `use (a, b = not [true,"1"])`),
		// nsr(2, `use 'b'`, r_decl_use),
		// nsr(3, `use a 'b'`),
		// nsr(4, `use ('b')`),
		// nsr(5, `use ( 'b',)`),
		// nsr(6, `use ('b','d')`, r_decl_use),
		// nsr(5, `use (a 'b')`),
		// nsr(8, `use (a 'b', c 'd')`),
		// nsr(9, `use (a 'b', c 'd',)`),
		// nsr(7, "use (\na 'b'\n)"),
		// nsr(8, "use (\n\ta 'b',\n\t)"),
		// nsr(7, `use ( 'b', c 'd')`),
		// nsr(8, "use ( 'b',\nc 'd')"),
		// nsr(8, "use ( 'b',\n\tc 'd')"),
		// nsr(11, "use (世界\n\n你好\na 'b',\n\n\tc 'd/e'\n\n\t)"),

		// nsr(8, "use ( 'b'\n, c 'd')"),

		// nsr(4, `type a{}`, r_decl_type),
		// nsr(6, `type a{int b}`),
		// nsr(8, `type a{int b,c}`),
	}

	bad = []ns{
		nsr(1, `use`, r_decl_use),
		nsr(2, `use a`),
		nsr(3, `use 'd')`),
		nsr(3, `use ('d'`),
		nsr(2, `use (a,'d')`),
		nsr(5, `use a 'b' c 'd'`),
		nsr(6, `use a 'b', c 'd'`),
		nsr(6, `use ('b' c 'd')`),
		nsr(6, `use (a 'b' 'd')`),
		nsr(7, `use (a 'b' c 'd')`),
	}
}

func Test_good(t *testing.T) {

	var r Rule
	for i, ns := range good {
		if ns.r != nil {
			r = ns.r
		}
		n, tok := testParse(ns.s, r)
		if n != ns.n || tok != token.EOF {
			t.Fatal(i, n, tok, ns.n, ":", ns.s)
		}
		r.Reset()
	}
}

func Test_bad(t *testing.T) {
	return
	var r Rule
	for i, ns := range bad {
		if ns.r != nil {
			r = ns.r
		}
		n, tok := testParse(ns.s, r)
		if tok == token.EOF {
			t.Fatal(i, n, ns.n, ":", ns.s)
		}
		r.Reset()
	}
}

func testParse(src string, r Rule) (n int, tok token.Token) {
	eat := 0
	nl, ok := false, false
	parser.Fast(
		[]byte(src),
		func(_ scanner.Pos, t token.Token, _ string) error {

			// 去占位, 合并连续的换行
			if t > token.IDENT || nl && t == token.NL {
				return nil
			}
			//fmt.Println(n, t)
			if t == token.EOF {
				if !ok { // 让 Option, More 结束
					eat, ok = r.Eat(t)
				}
				if ok {
					tok = token.EOF
				}
				return nil
			}

			tok = t
			nl = t == token.NL
			eat, ok = r.Eat(t)

			if eat == 0 {
				return io.ErrNoProgress
			}
			n++
			return nil
		})

	return
}
