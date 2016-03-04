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

	expr0 := More(Term{token.COMMA, token.Operator}, nil)
	expr := Any(
		Seq(
			Term{token.NOT, token.ANTI},
			expr0,
		),
		Seq(
			Term{token.Literal, token.IDENT},
			Option(expr0),
		),
		Seq(
			Term{token.LEFT},
			Any(
				Term{token.RIGHT},
				Seq(expr0, Term{token.RIGHT}),
			),
		),
	)

	expr0.Set(expr)

	r_decl_use = Seq(
		Term{token.USE},
		Any(
			Term{token.VALSTRING},
			Seq(Term{token.IDENT}, Term{token.VALSTRING}),
			Seq(
				Term{token.LEFT}, Option(Term{token.NL}),
				More(Term{token.COMMA, token.NL},
					Seq(
						Option(Term{token.NL}),
						Any(
							Term{token.VALSTRING},
							Seq(Term{token.IDENT}, Term{token.VALSTRING}),
						),
					),
				),
				Option(Term{token.NL}),
				Term{token.RIGHT},
			),
		),
	)

	style_type_body := More(Term{token.SEMICOLON, token.NL}, Seq(
		Option(Term{token.PUB}), Option(Term{token.STATIC}),
		Term{token.Type, token.IDENT},
		More(
			Term{token.COMMA},
			Seq(
				Term{token.IDENT},
				Option(Seq(
					Term{token.ASSIGN, token.COLON},
					expr,
				)),
			),
		),
	))

	r_decl_type = Seq(
		Term{token.TYPE}, Term{token.IDENT},
		Any(
			Term{token.IDENT, token.Type}, // 别名
			Seq(
				Term{token.LEFT},
				Option(style_type_body),
				Term{token.RIGHT},
			),
		),
	)

	// r_ 临时测试用
	r_ = Seq(
		Term{token.TYPE}, Term{token.IDENT},
		Any(
			Term{token.IDENT, token.Type}, // 别名
			Seq(
				Term{token.LEFT},
				Option(style_type_body),
				Term{token.RIGHT},
			),
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
		nsr(8, `type a{int b =1}`, r_),
		nsr(2, `use 'b'`, r_decl_use),
		nsr(3, `use a 'b'`),
		nsr(4, `use ('b')`),
		nsr(5, `use ( 'b',)`),
		nsr(5, `use (a 'b')`),
		nsr(8, `use (a 'b', c 'd')`),
		nsr(9, `use (a 'b', c 'd',)`),
		nsr(7, "use (\na 'b'\n)"),
		nsr(8, "use (\n\ta 'b',\n\t)"),
		nsr(7, `use ( 'b', c 'd')`),
		nsr(8, "use ( 'b',\nc 'd')"),
		nsr(8, "use ( 'b',\n\tc 'd')"),
		nsr(11, "use (世界\n\n你好\na 'b',\n\n\tc 'd/e'\n\n\t)"),

		nsr(3, `type a int`, r_decl_type),
		nsr(4, `type a{}`),
		nsr(6, `type a{int b}`),
		nsr(6, `type a{i8 b}`),
		nsr(8, `type a{i b,c}`),

		nsr(9, "type a{int b\nint c}"),
		nsr(9, "type a{int b;int c}"),
		nsr(8, `type a{int b =1}`),
		nsr(10, `type a{int b,c=1}`),
		nsr(8, `type a{int b =true}`),
		nsr(9, `type a{int b =[]}`),
		nsr(10, `type a{int b =[1]}`),
		nsr(12, `type a{int b =[1,true]}`),
		nsr(13, `type a{c b = not [true,"1"]}`),
		nsr(11, `type a{int b =[[]]}`),
		nsr(10, "type a{pub int b;int c}"),
		nsr(11, "type a{pub int b;static int c}"),
		nsr(11, "type a{int b;pub static int c}"),
	}

	bad = []ns{
		nsr(1, `use`, r_decl_use),
		nsr(2, `use a`),
		nsr(3, `use 'd')`),
		nsr(3, `use ('d'`),
		nsr(6, `use (a,'d')`),
		nsr(5, `use a 'b' c 'd'`),
		nsr(6, `use a 'b', c 'd'`),
		nsr(6, `use ('b' c 'd')`),
		nsr(6, `use (a 'b' 'd')`),
		nsr(7, `use (a 'b' c 'd')`),
		// 未来可能会支持下列
		nsr(8, "use ( 'b'\n, c 'd')"),
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
				if !ok {
					// 让 More 结束
					_, ok = r.Eat(t)
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
