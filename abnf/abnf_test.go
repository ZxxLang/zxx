package abnf_test

import (
	"testing"

	. "github.com/ZxxLang/zxx/abnf"
	. "github.com/ZxxLang/zxx/token"
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

func ts(tok ...Token) []Token {
	return tok
}

func TestTerm(t *testing.T) {
	toks := ts(
		FUNC, IDENT, OUT, INT, BOOL,
	)

	i, tok, flag := test(toks, Term(toks...))
	if flag != M|F || tok != BOOL {
		t.Fatal(i, tok, flag)
	}
	// fatal 0

	i, tok, flag = test(ts(FUNC, NAN), Term(toks...))
	if flag != 0 || tok != NAN {
		t.Fatal(i, tok, flag)
	}
}

func TestSeq(t *testing.T) {
	toks := ts(
		FUNC, IDENT, OUT, INT, BOOL,
	)

	rule := Seq(
		Term(FUNC), Term(IDENT),
		Term(OUT), Term(INT), Term(BOOL),
	)

	i, tok, flag := test(toks, rule)
	if flag != M|F || tok != BOOL {
		t.Fatal(i, tok, flag)
	}

	out := Seq(
		Term(OUT),
		Term(Type),
	)

	i, tok, flag = test(ts(OUT, INT), out)
	if flag != M|F || tok != INT {
		t.Fatal(i, tok, flag)
	}

	// fatal test

	i, tok, flag = test(ts(FUNC, IDENT), rule)
	if flag != M || tok != IDENT {
		t.Fatal(i, tok, flag)
	}

	rule.Match(EOF)

	i, tok, flag = test(ts(FUNC, IDENT, NAN), rule)
	if flag != 0 || tok != NAN {
		t.Fatal(i, tok, flag)
	}

}

func TestMore(t *testing.T) {

	toks := ts(INT, COMMA, BOOL, COMMA)

	rules := []Rule{
		More(Term(Type, COMMA), nil),
		More(Term(Type), Term(COMMA)),
	}

	for k, rule := range rules {
		for j := 0; j < len(toks); j++ {
			i, tok, flag := test(toks[0:j+1], rule)
			if flag != M|S || tok != toks[j] {
				t.Fatal(j, i, tok, flag)
			}

			if flag = rule.Match(NAN); flag != F {
				t.Fatal(NAN, "want F", flag, k, toks[j])
			}

			if flag = rule.Match(NAN); flag != 0 {
				t.Fatal(NAN, "want 0", flag, k, toks[j])
			}
		}
	}

	rule := rules[1]
	toks = ts(INT, COMMA, BOOL, TRUE)
	i, tok, flag := test(toks, rule)

	if flag != F || tok != TRUE {
		t.Fatal(i, tok, flag)
	}

	rule = More(
		Seq(
			Term(VAR),
			Term(Type),
		),
		Term(COMMA),
	)

	toks = ts(VAR, INT, COMMA)
	i, tok, flag = test(toks, rule)

	if flag != M|S || tok != COMMA {
		t.Fatal(i, tok, flag)
	}

	rule = Seq(
		Term(OUT),
		More(Term(Type), nil),
	)
	i, tok, flag = test(ts(OUT, INT), rule)

	if flag != M|S || tok != INT {
		t.Fatal(i, tok, flag)
	}

	// fatal test
	rule.Match(EOF)
	i, tok, flag = test(ts(OUT, NAN), rule)

	if flag != 0 || tok != NAN {
		t.Fatal(i, tok, flag)
	}

}

func TestOnce(t *testing.T) {
	toks := ts(FUNC, IDENT, OUT, INT, BOOL)

	out := Once(Seq(
		Term(OUT),
		Term(Type),
	))

	i, tok, flag := test(ts(OUT, INT), out)
	if flag != M|F || tok != INT {
		t.Fatal(i, tok, flag)
	}

	if flag = out.Match(NAN); flag != F {
		t.Fatal("want F", flag)
	}

	out = Once(Seq(
		Term(OUT),
		More(Term(Type), nil),
	))

	i, tok, flag = test(ts(OUT, INT), out)
	if flag != M|S || tok != INT {
		t.Fatal(i, tok, flag)
	}

	rule := Seq(
		Term(FUNC), Term(IDENT),
		More(
			Any(
				Term(Type),
				out,
			),
			out,
		),
	)
	out.Match(EOF)

	i, tok, flag = test(toks, rule)
	if flag != M|S || tok != BOOL {
		t.Fatal(i, tok, flag)
	}

	rule.Match(EOF)

	toks = ts(FUNC, IDENT, INT, STRING, OUT, INT, BOOL)

	i, tok, flag = test(toks, rule)
	if flag != M|S || tok != BOOL {
		t.Fatal(i, tok, flag)
	}

	// fatal test

	rule.Match(EOF)

	toks = ts(FUNC, IDENT, OUT, INT, STRING, OUT, INT, BOOL)

	i, tok, flag = test(toks, rule)
	if i != 5 || flag != F || tok != OUT {
		t.Fatal(i, tok, flag)
	}

}
