package ast

import "github.com/ZxxLang/zxx/token"

// 词法轻规则, 配合些许 Resolve 代码完成构建
// 所有的规则必须能吃掉一个 Token 才能继续

// Factor 代表规则的因子
type Factor interface {
	Token() token.Token
}

type Rule interface {
	// Eat 返回 Rule 匹配的情况.
	//
	// n	表示前 n 个 token 被匹配
	// want	表示是否期望更多的 Factor
	Eat(tok Factor) (n int, want bool)

	// Ok 返回规则是否被完整匹配
	Ok() bool

	// Reset 重置状态标记以便可以开始新的匹配
	Reset()
}

// Term 包装 Factor, 使用 Alternative 规则, 最多吃掉一个 Token
type Term []Factor

func (r Term) Eat(tok Factor) (int, bool) {
	for _, t := range r {
		if t.Token().Has(tok.Token()) {
			return 1, false
		}
	}
	return 0, false
}

// Zero 吃掉零个或一个 Token
type Zero struct {
	Term
}

func (r Zero) Eat(tok Factor) (int, bool) {
	n, _ := r.Term.Eat(tok)
	return n, false
}

// more 的实现方法可能会改变, 不导出
type more struct {
	Rule
	n int
}

// More 重复匹配规则, 吃掉一个或者多个 Token
func More(r Rule) Rule {
	return &more{r, 0}
}

func (r *more) Eat(tok Factor) (int, bool) {
	n, _ := r.Rule.Eat(tok)
	r.n += n
	return n, n != 0
}

// Any 使用 Alternative 规则, 最多吃掉一个 Token
// Any 和 Term 只是元素类型的不同.
type alternative struct {
	Rules []Rule
	Pos   int // sequence 会维护 Pos 归零.
	//ok    bool
}

func Any(rule ...Rule) Rule {
	return &alternative{rule, 0}
}

// sequence
type sequence struct {
	Rules []Rule
	Pos   int // sequence 会维护 Pos 归零.
}

// Seq 执行 Concatenation 规则, 顺序匹配 rule
func Seq(rule ...Rule) Rule {
	return &sequence{rule, 0}
}

func (q *alternative) Eat(tok Factor) (int, bool) {
	if q.Pos == len(q.Rules) {
		q.Pos = 0
	}
	for q.Pos < len(q.Rules) {
		r := q.Rules[q.Pos]
		n, w := r.Eat(tok)

		if n == 1 {
			if !w {
				q.Pos = len(q.Rules)
			}
			return 1, w
		}

		if !w {
			q.Pos++
			continue
		}
		break
	}
	q.Pos = 0
	return 0, false
}

func (q *sequence) Eat(tok Factor) (int, bool) {
	if q.Pos == len(q.Rules) {
		q.Pos = 0
	}
	for q.Pos < len(q.Rules) {
		r := q.Rules[q.Pos]
		n, w := r.Eat(tok)

		if n == 1 {
			if !w {
				q.Pos++
				return 1, q.Pos < len(q.Rules)
			}
			return 1, w
		}

		if !w && r.Ok() { // Zero
			q.Pos++
			continue
		}
		break
	}
	q.Pos = 0
	return 0, false
}

func (r Term) Reset()  {}
func (r *more) Reset() { r.n = 0 }

func (r Term) Ok() bool  { return false }
func (r Zero) Ok() bool  { return true }
func (r *more) Ok() bool { return r.n != 0 }

func (q *alternative) Reset() {
	for _, r := range q.Rules {
		r.Reset()
	}
	q.Pos = 0
}

func (q *sequence) Reset() {
	for _, r := range q.Rules {
		r.Reset()
	}
	q.Pos = 0
}

func (q *alternative) Ok() bool {
	return q.Pos == len(q.Rules)
}

func (q *sequence) Ok() bool {
	for i := q.Pos; i < len(q.Rules); i++ {
		if !q.Rules[i].Ok() {
			return false
		}
	}
	return true
}
