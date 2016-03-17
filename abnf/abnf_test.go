package abnf_test

import (
	"testing"

	. "github.com/ZxxLang/zxx/abnf"
	. "github.com/ZxxLang/zxx/token"
)

const (
	MS = M | S
	MF = M | F
)

func test(toks []Token, rule Rule) (i int, tok Token, flag Flag) {
	for i, tok = range toks {
		flag = rule.Match(tok)
		if !flag.Has(Matched) {
			return
		}
	}
	return
}

type step struct {
	tok  Token
	flag Flag
}

func tokens(s ...Token) []Token {
	return s
}

func flags(s ...Flag) []Flag {
	return s
}

func ts(toks []Token, flag []Flag) []step {
	steps := make([]step, len(toks))
	for i, tok := range toks {
		steps[i] = step{tok, flag[i]}
	}
	return steps
}

func TestTerm(t *testing.T) {
	toks := tokens(FUNC, IDENT, OUT, INT, BOOL)

	i, tok, flag := test(toks, Term(toks...))
	if flag != M|F || tok != BOOL {
		t.Fatal(i, tok, flag)
	}
	// fatal 0

	i, tok, flag = test(tokens(FUNC, NAN), Term(toks...))
	if flag != 0 || tok != NAN {
		t.Fatal(i, tok, flag)
	}

	flag = Term().Match(NAN) | Term().Match(EOF)
	if flag != F {
		t.Fatal(`want F`, flag)
	}
}

func TestSeq(t *testing.T) {
	rule := Seq(
		Term(FUNC), Term(IDENT),
		Term(OUT), Term(INT), Term(BOOL),
	)

	walk := ts(
		tokens(FUNC, IDENT, OUT, INT, BOOL, NAN, EOF),
		flags(M, M, M, M, MF, 0, F),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

	rule = Seq(
		Term(FUNC), Term(IDENT), Option(Term(OUT)), Option(Term(INT)),
	)

	walk = ts(
		tokens(FUNC, IDENT, SPACES),
		flags(M, MS, F),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

}

func TestMore(t *testing.T) {

	rule := More(Term(Type), Term(COMMA))

	walk := ts(
		tokens(INT, COMMA, BOOL, COMMA, NAN, NAN, EOF),
		flags(MS, M, MS, M, F, 0, F),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

	rule = More(
		Seq(
			Term(VAR),
			Term(Type),
		),
		Term(SEMICOLON),
	)

	walk = ts(
		tokens(VAR, INT, SEMICOLON, VAR, INT, SEMICOLON, NAN, NAN, EOF),
		flags(M, MS, M, M, MS, M, F, 0, F),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

	rule = Seq(
		Term(Type), Term(IDENT),
		Option(More(Seq(Term(COMMA), Term(IDENT)), nil)),
	)

	walk = ts(
		tokens(
			INT, IDENT, NAN, EOF,
		),
		flags(
			M, MS, F, F,
		),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}
}

func TestOnce(t *testing.T) {

	rule := Once(Seq(
		Term(Type),
		Term(COMMA),
	))

	walk := ts(
		tokens(INT, COMMA, BOOL, EOF, BOOL, COMMA, NAN, NAN, EOF),
		flags(M, MF, F, F, M, MF, F, F, F),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

}

func TestOption(t *testing.T) {
	rule := Option(Term(IDENT))
	walk := ts(
		tokens(
			USE, IDENT, NAN, EOF,
			IDENT, IDENT, NAN, EOF,
		),
		flags(
			F, MF, F, F,
			MF, MF, F, F,
		),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

	rule = Seq(
		Term(USE), Option(Term(IDENT)),
		Term(VALSTRING),
	)

	walk = ts(
		tokens(
			USE, IDENT, VALSTRING, NAN, EOF,
			USE, VALSTRING, NAN, EOF,
		),
		flags(
			M, M, MF, 0, F,
			M, MF, 0, F,
		),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

	rule = More(Seq(
		Term(Type), Term(IDENT),
		Option(More(Seq(Term(COMMA), Term(IDENT)), nil)),
	), nil)

	walk = ts(
		tokens(
			INT, IDENT, NAN, EOF,
			INT, IDENT, COMMA, IDENT, NAN, EOF,
			INT, IDENT, COMMA, IDENT,
			INT, IDENT, COMMA, IDENT, NAN, EOF,
		),
		flags(
			M, MS, F, F,
			M, MS, M, MS, F, F,
			M, MS, M, MS,
			M, MS, M, MS, F, F,
		),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

	rule = More(
		Seq(Term(Type), Term(IDENT)),
		Option(More(Seq(Term(COMMA), Term(IDENT)), nil)),
	)

	walk = ts(
		tokens(
			INT, IDENT, NAN, EOF,
			INT, IDENT, COMMA, IDENT, NAN, EOF,
			INT, IDENT, COMMA, IDENT,
			INT, IDENT, COMMA, IDENT, NAN, EOF,
		),
		flags(
			M, MS, F, F,
			M, MS, M, M, F, F,
			M, MS, M, M,
			M, MS, M, M, F, F,
		),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

}

func TestAny(t *testing.T) {
	tmap := Seq()
	tfunc := Seq()
	tarray := Seq()

	length := Seq(Term(COMMA), Term(VALINTEGER, Type, IDENT, MEMBER))

	types := Any(
		Term(Type, IDENT, MEMBER),
		tmap,
		tfunc,
		tarray,
	)

	tarray.Bind(Seq(Term(ARRAY), Term(LEFT), types, Option(length), Term(RIGHT)))

	tmap.Bind(Seq(Term(MAP), Term(LEFT), types, types, Option(length), Term(RIGHT)))

	params := More(
		Seq(
			Option(Once(Term(OUT))),
			types,
			Option(Seq(
				Term(IDENT),
				Option(More(Seq(Term(COMMA), Term(IDENT)), nil)),
			)),
		),
		Term(SEMICOLON),
	)

	tfunc.Bind(Seq(
		Term(FUNC),
		Option(Seq(Term(LEFT), params, Term(RIGHT))),
	))

	rule := params

	walk := ts(
		tokens(
			INT, IDENT, SEMICOLON,
			BOOL, IDENT, COMMA, IDENT, SEMICOLON,
		),
		flags(
			MS, MS, M,
			MS, MS, M, MS, M,
		),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

	rule = params

	walk = ts(
		tokens(
			INT, IDENT, SEMICOLON,
			BOOL, IDENT, COMMA, IDENT, NAN,

			FUNC, LEFT, IDENT, SEMICOLON,
			BOOL, IDENT, SEMICOLON,
			OUT, F32, SEMICOLON, F64, SEMICOLON, RIGHT, IDENT,
		),
		flags(
			MS, MS, M,
			MS, MS, M, MS, F,
			MS, M, M, M,
			M, M, M,
			M, M, M, M, M, MS, MS,
		),
	)

	for i, w := range walk {
		if flag := rule.Match(w.tok); flag != w.flag {
			t.Fatal(i, flag, w.tok, w.flag)
		}
	}

	rule = Seq(
		Term(PROC), Term(IDENT, MEMBER),
		Option(Any(
			Seq(
				Term(OUT),
				More(Seq(types, Option(Seq(Term(IDENT), Option(More(Seq(Term(COMMA), Term(IDENT)), nil))))),
					nil),
			),
			More(Seq(types, Term(IDENT), Option(More(Seq(Term(COMMA), Term(IDENT)), nil))),
				nil),
		)),
	)
}
