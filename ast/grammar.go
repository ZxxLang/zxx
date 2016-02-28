package ast

import (
	"fmt"

	"github.com/ZxxLang/zxx/token"
)

// 词法轻规则, 配合些许 Resolve 代码完成构建
// 所有的规则必须能吃掉一个 Token 才能继续
// http://tools.ietf.org/html/rfc5234

type Rule interface {
	// Eat 返回 Rule 匹配 tok 的情况:
	//
	// 0, false 匹配失败
	// 1, false 匹配成功, 但不完整, 还需要 tok
	// 0, true  匹配完成, 但需要其它 Rule 匹配 tok
	// 1, true	匹配完成
	Eat(tok token.Token) (n int, ok bool)

	// Reset 重置状态标记以便可以开始新的匹配
	Reset()
}

// Term 规则任一 token.Token 匹配
type Term []token.Token

func (r Term) Eat(tok token.Token) (int, bool) {
	for _, t := range r {
		if t.Has(tok) {
			return 1, true
		}
	}
	return 0, false
}

// option 重复 {0,1}
type option struct {
	Rule
}

// Option  产生重复匹配 r{0,1}
func Option(r Rule) Rule {
	return option{r}
}

func (r option) Eat(tok token.Token) (int, bool) {
	n, ok := r.Rule.Eat(tok)

	if !ok {
		if n == 1 {
			// 继续
			return n, false
		}
		r.Rule.Reset()
		return n, true
	}

	return n, true
}

// more 的实现方法可能会改变, 不导出
type more struct {
	Rule
	sep token.Token
	n   int // 成功的次数
}

// More  产生重复匹配 r{1,} Rule
// sep 是多次重复匹配的分隔 token.Token.
func More(sep token.Token, r Rule) Rule {
	return &more{r, sep, 0}
}

func (r *more) Eat(tok token.Token) (n int, ok bool) {

	n, ok = r.Rule.Eat(tok)

	if n == 0 && (ok || r.n == 1) && r.sep != token.EOF {
		if r.sep.Has(tok) {
			n = 1
			r.n = 2
			ok = false
			r.Rule.Reset()
			return
		}
	}

	if ok {
		r.n = 1
	}

	if n == 1 { // 持续吃
		ok = false
	} else {
		ok = r.n != 0
		r.n = 0
		r.Rule.Reset()
	}
	return
}

// Alternative 规则, 最多吃掉一个 Token
type Alternative struct {
	Rules []Rule
	pos   int
}

// Any 产生任一匹配 Rule
func Any(rule ...Rule) Rule {
	return &Alternative{rule, 0}
}

// Concatenation
type Concatenation struct {
	Rules []Rule
	pos   int
	zero  int
}

// Seq 产生序列匹配 Rule
func Seq(rule ...Rule) *Concatenation {
	return &Concatenation{rule, 0, 0}
}

func (q *Alternative) Eat(tok token.Token) (int, bool) {
	if q.pos == len(q.Rules) {
		q.pos = 0
	}
	for ; q.pos < len(q.Rules); q.pos++ {

		n, ok := q.Rules[q.pos].Eat(tok)

		if n == 1 {
			if ok {
				q.pos = len(q.Rules)
			}
			return n, ok
		}
		q.Rules[q.pos].Reset()
	}

	q.pos = 0
	return 0, false
}

func (q *Concatenation) Eat(tok token.Token) (n int, ok bool) {
	if q.pos == len(q.Rules) {
		q.pos = 0
	}

	for ; q.pos < len(q.Rules); q.pos++ {
		n, ok = q.Rules[q.pos].Eat(tok)
		if q.zero == 1 {
			fmt.Println(n, ok, tok, q.pos)
		}
		if !ok {
			if n == 0 {
				q.pos = 0
				q.Rules[q.pos].Reset()
			}
			return
		}

		if n == 0 {
			continue
		}
		q.pos++
		break
	}

	ok = q.pos == len(q.Rules)
	return
}

func (r Term) Reset()  {}
func (r *more) Reset() { r.n = 0 }

func (q *Alternative) Reset() {
	for _, r := range q.Rules {
		r.Reset()
	}
	q.pos = 0
}

func (q *Concatenation) Reset() {
	for _, r := range q.Rules {
		r.Reset()
	}
	q.pos = q.zero
}
