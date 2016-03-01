package ast

import "github.com/ZxxLang/zxx/token"

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

	// Set 用来设置规则, 只有 Any 可用
	Set(...Rule)
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
	reset bool
}

// Option  产生重复匹配 r{0,1}
func Option(r Rule) Rule {
	return option{Rule: r}
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
	sep   Term
	n     int // 成功的次数
	reset bool
}

// More  产生重复匹配 r{1,} Rule
// sep 是多次重复匹配的分隔 token.Token.
func More(sep Term, r Rule) Rule {
	return &more{Rule: r, sep: sep, n: 0}
}

func (r *more) Eat(tok token.Token) (n int, ok bool) {

	n, ok = r.Rule.Eat(tok)

	if n == 0 && (ok || r.n == 1) && r.sep != nil {
		if x, y := r.sep.Eat(tok); y && x == 1 {
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

// alternative 规则, 最多吃掉一个 Token
type alternative struct {
	Rules []Rule
	pos   int
	reset bool
	next  *alternative // 递归生成新的规则
}

// Any 产生任一匹配 Rule
func Any(rule ...Rule) Rule {
	return &alternative{Rules: rule, pos: 0}
}

// concatenation
type concatenation struct {
	Rules []Rule
	pos   int
	reset bool
}

// Seq 产生序列匹配 Rule
func Seq(rule ...Rule) Rule {
	return &concatenation{Rules: rule, pos: 0}
}

func (q *alternative) Eat(tok token.Token) (int, bool) {
	if q.reset {
		if q.next == nil {
			q.next = &alternative{q.Rules, 0, false, nil}
		}
		n, ok := q.next.Eat(tok)
		if ok {
			q.next = nil
		}
		return n, ok
	}
	q.reset = true
	if q.pos == len(q.Rules) {
		q.pos = 0
	}
	for ; q.pos < len(q.Rules); q.pos++ {

		n, ok := q.Rules[q.pos].Eat(tok)

		if n == 1 {
			if ok {
				q.pos = len(q.Rules)
			}
			q.reset = false
			return n, ok
		}
		q.Rules[q.pos].Reset()
	}

	q.pos = 0
	q.reset = false
	return 0, false
}

func (q *concatenation) Eat(tok token.Token) (n int, ok bool) {
	if q.pos == len(q.Rules) {
		q.Reset()
		q.pos = 0
	}

	for ; q.pos < len(q.Rules); q.pos++ {
		n, ok = q.Rules[q.pos].Eat(tok)

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

func (r Term) Reset() {}

func (r option) Reset() {
	if r.reset {
		return
	}
	r.reset = true
	r.Rule.Reset()
	r.reset = false
}

func (r *more) Reset() {
	if r.reset {
		return
	}
	r.reset = true
	r = More(r.sep, r.Rule).(*more)
	r.Rule.Reset()
	r.reset = false
}

func (q *alternative) Reset() {
	if q.reset {
		return
	}
	q.reset = true

	q = Any(q.Rules...).(*alternative)
	for _, r := range q.Rules {
		r.Reset()
	}
	q.reset = false
}

func (q *concatenation) Reset() {
	if q.reset {
		return
	}
	q.reset = true

	q = Seq(q.Rules...).(*concatenation)
	for _, r := range q.Rules {
		r.Reset()
	}
	q.reset = false
}

func (r Term) Set(...Rule)   {}
func (r option) Set(...Rule) {}
func (r *more) Set(...Rule)  {}
func (q *alternative) Set(r ...Rule) {
	q.Rules = r
}

func (q *concatenation) Set(r ...Rule) {
	// q.Rules = r
}
