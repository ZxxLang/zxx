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

	// processing 重置状态标记以便可以开始新的匹配
	Reset()

	// Set 用来设置规则, 只有 More, Seq 规则可用.
	// 有了该方法可解决递归规则.
	Set(Rule)

	Clone() Rule

	Name() string
	SetName(string) Rule
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

func (r Term) Clone() Rule {
	return r
}

// option 重复 {0,1}
type option struct {
	Rule
	n          int
	name       string
	processing bool
	next       *option
}

// Option  产生重复匹配 r{0,1}
func Option(r Rule) Rule {
	return &option{Rule: r}
}

func (r *option) Clone() Rule {
	return &option{Rule: r.Rule.Clone(), name: r.name}
}

func (r *option) Eat(tok token.Token) (n int, ok bool) {
	if r.processing {
		if r.next == nil {
			r.next = &option{Rule: r.Rule.Clone()}
		}
		n, ok = r.next.Eat(tok)
		if !ok && n == 0 {
			r.next.Rule.Reset()
		}
		return
	}

	r.processing = true

	n, ok = r.Rule.Eat(tok)

	r.n += n
	if n == 0 {
		ok = ok || r.n == 0
		r.n = 0
		r.Rule.Reset()
	}

	r.processing = false
	return
}

// more 的实现方法可能会改变, 不导出
type more struct {
	Rule
	sep        Term
	next       *more
	n          int
	name       string
	more       bool // 期望逗号
	processing bool
}

// More  产生重复匹配 r{1,} Rule
// sep 是多次重复匹配的分隔 token.Token.
// More 要求规则只是吃掉一个 Token 才会成功.
func More(sep Term, r Rule) Rule {
	return &more{Rule: r, sep: sep}
}

type onceSep struct {
	name string
	rule Rule
	once Rule
	did  bool
}

// Once 产生当 rule 吃不掉 Token 时尝试一次 once 的规则.
// 如果 once 成功, 返回 1, false, 会继续使用 rule.
// 无论成功失败 once 只执行一次.
func Once(once, rule Rule) Rule {
	return &onceSep{"", rule, once, false}
}

func (r *onceSep) Clone() Rule {
	return &onceSep{r.name, r.rule.Clone(), r.once.Clone(), false}
}

func (r *onceSep) Eat(tok token.Token) (n int, ok bool) {
	// 效率稍低
	n, ok = r.rule.Eat(tok)
	if n == 0 && !r.did {
		r.did = true
		n, _ = r.once.Eat(tok)
		if n == 1 {
			ok = false
		}
	}
	return
}

func (r *more) Clone() Rule {
	return &more{Rule: r.Rule.Clone(), sep: r.sep, name: r.name}
}

func (r *more) Eat(tok token.Token) (n int, ok bool) {
	if r.processing {
		if r.next == nil {
			r.next = &more{Rule: r.Rule.Clone(), sep: r.sep}
		}
		n, ok = r.next.Eat(tok)
		if !ok && n == 0 {
			r.next.Reset()
		}
		return
	}

	r.processing = true

	if r.more {
		r.more = false
		n, _ = r.sep.Eat(tok)
		if n == 0 {
			r.n = 0 // end
			ok = true
		}
	} else {
		n, ok = r.Rule.Eat(tok)

		if ok {
			r.n = 1
		}
		r.more = r.sep != nil && ok && n == 1 // 下次尝试 sep
		if n == 0 {

			if r.n != 0 { // 立即尝试 sep
				r.more = false

				if r.sep != nil {
					n, _ = r.sep.Eat(tok)
					ok = n == 0
					r.Rule.Reset()
				} else if ok && tok != token.EOF {
					// 再次尝试
					n, ok = r.Rule.Eat(tok)
				} else {
					ok = true // end
				}

				if n == 0 {
					r.n = 0 // end
				}

			} else {
				ok = r.sep == nil && tok == token.EOF
				r.n = 0 // end
				r.more = false
			}

		} else {
			ok = false
		}
	}

	r.processing = false
	return
}

// alternative 规则, 最多吃掉一个 Token
type alternative struct {
	name       string
	Rules      []Rule
	pos        int
	next       *alternative
	must       Rule
	processing bool
	cloning    bool
}

// Any 产生任一匹配 Rule. 这是唯一支持递归的规则.
func Any(rule ...Rule) Rule {
	return &alternative{Rules: rule, pos: 0}
}

func (r *alternative) Clone() Rule {
	return r.clone()
}

func (r *alternative) clone() *alternative {
	if r.cloning {
		return r
	}
	r.cloning = true
	nr := &alternative{}
	nr.name = r.name
	nr.Rules = make([]Rule, len(r.Rules))
	for i := 0; i < len(nr.Rules); i++ {
		nr.Rules[i] = r.Rules[i].Clone()
	}
	r.cloning = false
	return nr
}

// concatenation
type concatenation struct {
	name       string
	Rules      []Rule
	pos        int
	next       *concatenation
	processing bool
	cloning    bool
}

// Seq 产生序列匹配 Rule
func Seq(rule ...Rule) Rule {
	return &concatenation{Rules: rule, pos: 0}
}

func (r *concatenation) Clone() Rule {
	return r.clone()
}

func (r *concatenation) clone() *concatenation {
	if r.cloning {
		return r
	}
	r.cloning = true
	nr := &concatenation{}
	nr.name = r.name
	nr.Rules = make([]Rule, len(r.Rules))
	for i := 0; i < len(nr.Rules); i++ {
		nr.Rules[i] = r.Rules[i].Clone()
	}
	r.cloning = false
	return nr
}

func (q *alternative) Eat(tok token.Token) (n int, ok bool) {
	if q.processing {

		if q.next == nil {
			q.next = q.clone()
		}
		n, ok = q.next.Eat(tok)
		if !ok && n == 0 {
			ok = true
			q.next.Reset()
		}
		return
	}

	q.processing = true

	// 继续上次没有进行完的 Rule
	if q.must != nil {
		n, ok = q.must.Eat(tok)
		if ok || n == 0 {
			q.must.Reset()
			q.must = nil
			q.pos = len(q.Rules)
		}
		q.processing = false
		return
	}

	if q.pos == len(q.Rules) {
		q.pos = 0
	}
	for ; q.pos < len(q.Rules); q.pos++ {

		n, ok = q.Rules[q.pos].Eat(tok)

		if n == 1 {
			if ok {
				q.pos = len(q.Rules)
			} else {
				q.must = q.Rules[q.pos]
			}
			break
		}

		if !ok {
			q.Rules[q.pos].Reset()
			continue
		}

		q.pos = len(q.Rules)
		break
	}
	q.processing = false
	return
}

func (q *concatenation) Eat(tok token.Token) (n int, ok bool) {

	if q.processing {
		if q.next == nil {
			q.next = q.clone()
		}
		n, ok = q.next.Eat(tok)
		if !ok && n == 0 {
			ok = true
			q.next.Reset()
		}
		return
	}

	// 上次成功过, 这次如果吃不掉, 那么返回 0,true
	prev := false
	q.processing = true
	if q.pos == len(q.Rules) {
		q.pos = 0
		prev = true
	}

	for ; q.pos < len(q.Rules); q.pos++ {
		n, ok = q.Rules[q.pos].Eat(tok)

		if n == 0 {
			if ok {
				continue
			}
			q.Rules[q.pos].Reset()
			q.pos = 0
			break
		}

		if ok {
			q.pos++
		}
		break
	}

	// 保持上次的成功状态
	if prev && !ok && n == 0 {
		q.pos = len(q.Rules)
	}
	ok = q.pos == len(q.Rules)

	q.processing = false
	return
}

func (r Term) Reset() {}

func (r *option) Reset() {
	if !r.processing {
		r.processing = true
		r.n = 0
		r.next = nil
		r.Rule.Reset()
		r.processing = false
	}
}

func (r *more) Reset() {
	if !r.processing {
		r.processing = true
		r.n = 0
		r.more = false
		r.next = nil
		r.Rule.Reset()
		r.processing = false
	}
}

func (r *onceSep) Reset() {
	r.did = false
	r.rule.Reset()
	r.once.Reset()
}

func (q *alternative) Reset() {
	if !q.processing {
		q.processing = true
		for _, r := range q.Rules {
			r.Reset()
		}
		q.pos = 0
		q.must = nil
		q.next = nil
		q.processing = false
	}
}

func (q *concatenation) Reset() {
	if !q.processing {
		q.processing = true
		for _, r := range q.Rules {
			r.Reset()
		}
		q.pos = 0
		q.next = nil
		q.processing = false
	}
}

func (r *onceSep) Set(rule Rule) {
	panic("not support")
}

func (r *more) Set(rule Rule) {
	if r.processing {
		panic("not support")
	}
	r.Rule = rule
}

func (r Term) Set(Rule)    { panic("not support") }
func (r *option) Set(Rule) { panic("not support") }
func (q *alternative) Set(Rule) {
	panic("not support")
}

func (q *concatenation) Set(rule Rule) {
	c, ok := rule.(*concatenation)
	if !ok || q.processing {
		panic("not support")
	}
	q.Rules = make([]Rule, len(c.Rules))
	for i := 0; i < len(q.Rules); i++ {
		q.Rules[i] = c.Rules[i]
	}
}

func (r Term) Name() string        { return "" }
func (r Term) SetName(string) Rule { return r }

func (r *option) Name() string { return r.name }
func (r *option) SetName(name string) Rule {
	r.name = name
	return r
}

func (r *more) Name() string { return r.name }
func (r *more) SetName(name string) Rule {
	r.name = name
	return r
}

func (r *onceSep) Name() string { return r.name }
func (r *onceSep) SetName(name string) Rule {
	r.name = name
	return r
}

func (r *alternative) Name() string { return r.name }
func (r *alternative) SetName(name string) Rule {
	r.name = name
	return r
}

func (r *concatenation) Name() string { return r.name }
func (r *concatenation) SetName(name string) Rule {
	r.name = name
	return r
}
