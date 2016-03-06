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
	r_decl_use, r_decl_const, r_decl_var, r_decl_type, r_decl_func, r_decl_proc,
	r_stmt,
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

	style_more_option_assign := More(
		Term{token.COMMA},
		Seq(
			Term{token.IDENT},
			Option(Seq(Term{token.ASSIGN}, expr)),
		),
	)

	r_decl_use = Seq(
		Term{token.USE},
		Any(
			Term{token.VALSTRING},
			Seq(Term{token.IDENT}, Term{token.VALSTRING}),
			Seq(
				Term{token.LEFT},
				More(Term{token.COMMA},
					Seq(
						Any(
							Term{token.VALSTRING},
							Seq(Term{token.IDENT}, Term{token.VALSTRING}),
						),
					),
				),
				Term{token.RIGHT},
			),
		),
	)

	r_decl_const = Seq(
		Term{token.CONST},
		Any(
			Seq(Term{token.IDENT}, Term{token.ASSIGN}, expr),
			Seq(
				Term{token.LEFT},
				More(Term{token.SEMICOLON},
					Seq(Term{token.IDENT}, Term{token.ASSIGN}, expr),
				),
				Term{token.RIGHT},
			),
		),
	)

	array_type := Seq().SetName("array")
	map_type := Seq()

	// 特别注意 token.Type 包括 map, array
	a_type := Any(
		Seq(Term{token.MAP}, map_type),
		Seq(Term{token.ARRAY}, array_type),
		Term{token.IDENT, token.MEMBER, token.Type},
	).SetName("T")

	array_type.Set(Seq(
		Term{token.LEFT},
		Option(Seq(
			Term{token.VALINTEGER, token.IDENT, token.MEMBER},
			Term{token.COMMA},
		)),
		a_type,
		Term{token.RIGHT},
	))

	map_type.Set(Seq(
		Term{token.LEFT},
		a_type, Term{token.COMMA}, a_type,
		Term{token.RIGHT},
	))

	r_decl_var = Seq(
		Term{token.VAR},
		Any(
			Seq(
				a_type,
				style_more_option_assign,
			),
			Seq(
				Term{token.LEFT},
				More(Term{token.SEMICOLON},
					Seq(
						Term{token.IDENT, token.Type},
						style_more_option_assign,
					),
				),
				Term{token.RIGHT},
			),
		),
	).SetName("var")

	style_type_body := More(Term{token.SEMICOLON},
		Seq(
			Option(Term{token.PUB}), Option(Term{token.STATIC}),
			Term{token.IDENT, token.MEMBER, token.Type},
			style_more_option_assign,
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

	r_decl_func = Seq(
		Term{token.FUNC}, Term{token.IDENT, token.MEMBER},
		Option(
			Once(Term{token.OUT},
				More(nil, Term{token.IDENT, token.MEMBER, token.Type})),
		),
	)

	r_decl_proc = Seq(
		Term{token.PROC}, Term{token.IDENT, token.MEMBER},
		Option(
			Once(Term{token.OUT},
				More(nil,
					Any(
						Seq(
							Term{token.IDENT, token.MEMBER, token.Type},
							More(Term{token.COMMA}, Term{token.IDENT}),
						),
						Seq(
							Term{token.FUNC}, Term{token.IDENT},
							Option(Once(Term{token.OUT},
								More(nil, Term{token.IDENT, token.MEMBER, token.Type})),
							),
						),
					),
				),
			),
		),
	)

	// r_ 临时测试用
	r_ = Seq(
		Term{token.USE},
		Seq(
			Term{token.LEFT},
			More(Term{token.COMMA, token.NL},
				Seq(
					Any(
						Term{token.VALSTRING},
						Seq(Term{token.IDENT}, Term{token.VALSTRING}),
					),
				),
			),
			Term{token.RIGHT},
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
		nsr(5, `use ( 'b',)`),
		nsr(5, `use (a 'b')`),
		nsr(8, `use (a 'b', c 'd')`),
		nsr(9, `use (a 'b', c 'd',)`),
		nsr(5, "use (\na 'b'\n)"),
		nsr(6, "use (\n\ta 'b',\n\t)"),
		nsr(7, `use ( 'b', c 'd')`),
		nsr(7, "use ( 'b',\nc 'd')"),
		nsr(7, "use ( 'b'\n, c 'd')"),
		nsr(7, "use ( 'b',\n\tc 'd')"),
		nsr(8, "use (世界\n\n你好\na 'b',\n\n\tc 'd/e'\n\n\t)"),

		nsr(4, `const a ='b'`, r_decl_const),
		nsr(6, `const {a ='b'}`),
		nsr(11, "const {a ='b'\n c=[]}"),
		nsr(11, "const {a ='b';c=[]}"),

		nsr(3, `var string a`, r_decl_var),
		nsr(5, `var string a=''`),
		nsr(7, `var string a,b=''`),
		nsr(13, "var {string a ='b';int c=[]}"),

		nsr(3, `type a int`, r_decl_type),
		nsr(4, `type a{}`),
		nsr(6, `type a{int b}`),
		nsr(6, `type a{i8 b}`),
		nsr(8, `type a{i b,c}`),

		nsr(9, "type a{int b\nint c}"),
		nsr(9, "type a{int b;int c}"),
		nsr(10, "type a{pub int b;int c}"),
		nsr(8, `type a{int b =1}`),
		nsr(10, `type a{int b,c=1}`),
		nsr(8, `type a{int b =true}`),
		nsr(9, `type a{int b =[]}`),
		nsr(10, `type a{int b =[1]}`),
		nsr(12, `type a{int b =[1,true]}`),
		nsr(13, `type a{c b = not [true,"1"]}`),
		nsr(11, `type a{int b =[[]]}`),
		nsr(11, "type a{pub int b;static int c}"),
		nsr(11, "type a{int b;pub static int c}"),

		nsr(2, `func a`, r_decl_func),
		nsr(2, `func a.b`),
		nsr(3, `func a int`),
		nsr(5, `func a int string bool`),
		nsr(5, `func a b out int`),
		nsr(6, `func a.b b out int c.d`),

		nsr(2, `proc a`, r_decl_proc),
		nsr(4, `proc a int b`),
		nsr(6, `proc a int a,b`),
		nsr(6, `proc a int a bool b`),
		nsr(10, `proc a int b,c bool d,e`),
		nsr(11, `proc a int b,c out bool d,e`),
		nsr(13, `proc a int b,c out bool d,e int i`),
		nsr(8, `proc a int b,c func call`),
		nsr(12, `proc a int b,c func call bool int out int`),
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
	ok := false
	parser.Fast(
		[]byte(src),
		func(_ scanner.Pos, t token.Token, _ string) error {

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

			// 去占位,
			if t > token.IDENT {
				return nil
			}
			if t == token.NL {
				// 合并连续的换行, 清理排版换行
				switch tok {
				case token.NL, token.EOF, token.COMMA, token.SEMICOLON,
					token.LEFT, token.RIGHT:
					return nil
				}
				// 推延换行
				tok = token.NL
				return nil
			}

			if tok == token.NL {
				switch t {
				case token.COMMA, token.SEMICOLON, token.RIGHT:
				default:
					// 换行改 ";"
					eat, ok = r.Eat(token.SEMICOLON)
					if eat == 0 {
						return io.ErrNoProgress
					}
					n++
				}
			}

			tok = t
			eat, ok = r.Eat(t)

			if eat == 0 {
				return io.ErrNoProgress
			}
			n++
			return nil
		})

	return
}
