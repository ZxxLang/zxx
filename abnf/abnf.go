// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// 本包实现 ABNF 规范中的规则模式算法.
// 参见: http://tools.ietf.org/html/rfc5234
//
// 使用其它 Token, fork 本包并修改 token 包路径, 设定 eof.
// 而且 Token 需实现 Has 方法, 参见 zxx/token#Token 实现.
//
package abnf

import . "github.com/ZxxLang/zxx/token"

const eof = EOF

type Rule interface {
	// Match 返回匹配 tok 的状态标记. 实现必须遵守以下约定:
	//
	// 返回值为下列之一:
	//
	//	0
	// 	Matched
	// 	Matched|Standing
	// 	Matched|Finished
	// 	Finished
	//
	// 自动重置:
	//
	// 规则状态为 0, Finished, Matched|Finished 时自动重置, 可接受新的匹配.
	//
	// EOF 重置: 当最终状态为
	//
	//	Matched 最终状态是不确定完整匹配.
	//	Matched|Standing 最终状态是完整匹配.
	//
	// 时使用 Match(EOF) 重置规则并返回 Finished, 其它状态不应该使用 EOF 来重置.
	//
	// 末尾完整测试:
	//
	// 类似 Seq(Term(XX),Option(YY),Option(ZZ)) 规则, 单个 XX 也是合法的,
	// 但是由于 Option 的原因, 匹配单个 XX 的状态为 Matched,
	// 因此再匹配一个不可能出现的 Token, 可以测试规则是否完整.
	//
	Match(tok Token) Flag

	// Bind 应只在递归嵌套规则变量时使用, 先声明规则变量后绑定规则.
	// 其他情况都不应该使用 Bind.
	Bind(Rule)

	// Clone 返回克隆规则, 这是深度克隆, 但不含递归.
	// 递归规则在 Match 中通过判断 Handing 标记及时建立的.
	Clone() Rule

	// IsOption 返回该规则是否为可选规则.
	// 事实上除了 Option 是明确的可选规则外, 其它组合可能产生事实上的可选规则.
	IsOption() bool
}

// Flag 表示规则匹配 Token 后的状态
type Flag int

const (
	Matched  Flag = 1 << iota // 匹配成功, 消耗掉一个 Token
	Standing                  // 规则成立, 可以继续匹配 Token
	Finished                  // 规则完整, 不再接受匹配 Token

	// 下列标记由算法维护, Match 返回值不应该包含这些标记.
	Handing // 正在进行处理中(Match), 可用来检查发生递归
	Cloning // 正在进行克隆
	Custom  // 通用标记位, 包括更高的位都由 Match 自定义用途,

	// 短标记
	M = Matched
	S = Standing
	F = Finished
	H = Handing
	C = Custom
)

// Has 返回 f&flag != 0
func (f Flag) Has(flag Flag) bool {
	return f&flag != 0
}

// Must 返回 f&flag == flag
func (f Flag) Must(flag Flag) bool {
	return f&flag == flag
}

func (f Flag) String() (s string) {

	const flagString = "MSFHC"

	if f == 0 {
		return "0"
	}

	for i := 0; i < len(flagString); i++ {
		if f&Flag(1<<uint(i)) != 0 {
			s += string(flagString[i])
		}
	}

	return
}

type term struct {
	toks []Token
}

// Term 产生任一 Token 匹配规则. Match 方法返回值为:
//
//	0
//	Matched | Finished
//	Finished	当 EOF 重置或者 tok == nil
//
func Term(tok ...Token) Rule {
	return term{tok}
}

func (r term) Bind(Rule)   { panic("not support") }
func (r term) Clone() Rule { return r }
func (r term) Match(tok Token) Flag {
	if tok == eof || r.toks == nil {
		return Finished
	}
	for _, t := range r.toks {
		if t.Has(tok) {
			return Matched | Finished
		}
	}
	return 0
}

func (r term) IsOption() bool {
	return r.toks == nil
}

// base 维护 Flag 和递归 next.
type base struct {
	Flag
	next Rule
}

func (b *base) match(r Rule, tok Token) Flag {
	if tok == eof && b.next == nil {
		b.Flag = 0
		return Finished
	}
	if b.next == nil {
		b.next = r.Clone()
	}
	return b.next.Match(tok)
}

func (b *base) eof() {
	b.Flag = 0
	if b.next != nil {
		b.next.Match(eof)
	}
	return
}

type option struct {
	rule Rule
	base
}

// Option 产生零次或一次匹配规则.
func Option(rule Rule) Rule {
	return &option{rule: rule}
}

func (r *option) IsOption() bool {
	return true
}

func (r *option) Bind(rule Rule) { r.rule = rule }
func (r *option) Clone() (c Rule) {
	if r.Has(Cloning) {
		return r
	}
	r.Flag += Cloning
	c = Option(r.rule.Clone())
	r.Flag -= Cloning
	return
}

func (r *option) Match(tok Token) (flag Flag) {
	if r.Has(Handing) {
		return r.match(r, tok)
	}

	r.Flag += Handing
	flag = r.rule.Match(tok)
	r.Flag -= Handing

	switch flag {
	case 0:
		if r.Flag == 0 || r.Flag.Has(Standing) {
			flag = Finished
		}
		r.Flag = 0
	case Matched, Matched | Standing:
		r.Flag = flag
	case Matched | Finished, Finished:
		r.Flag = 0
	}

	return
}

type once struct {
	rule Rule
	base
}

// Once 产生当 Match 返回 Finished 后, 再次 Match 总是返回 Finished 的规则.
//
// 该规则需要进行 EOF 重置才能开始新的匹配.
//
// 该规则不支持递归, Clone() 返回它本身, 而不是产生副本.
//
func Once(rule Rule) Rule {
	return &once{rule: rule}
}

func (r *once) IsOption() bool {
	return r.rule.IsOption()
}

func (r *once) Bind(rule Rule) { r.rule = rule }
func (r *once) Clone() Rule    { return r }

func (r *once) Match(tok Token) (flag Flag) {
	if r.Has(Handing) {
		// 现实中不应该出现递归
		panic("Once does not support recursion")
	}
	if r.Has(Finished) {
		if tok == eof {
			r.rule.Match(eof)
			r.eof()
		}
		flag = Finished
		return
	}

	r.Flag += Handing
	flag = r.rule.Match(tok)
	r.Flag -= Handing

	if tok == eof {
		r.eof()
	} else if flag.Has(Finished) {
		r.Flag = Finished
	}

	return
}

type more struct {
	rule Rule
	base
	sep Rule
}

// More 产生重复匹配规则.
// rule 必须初次匹配成功, 然后当 rule 匹配结果为 0, Finished 时尝试 sep 匹配,
// 如果 sep 匹配成功则继续匹配 rule.
//
func More(rule, sep Rule) Rule {
	if sep == nil {
		sep = Term()
	}
	return &more{rule: rule, sep: sep}
}

func (r *more) IsOption() bool {
	return false
}

func (r *more) Bind(rule Rule) { r.rule = rule }
func (r *more) Clone() (c Rule) {
	if r.Has(Cloning) {
		return r
	}
	r.Flag += Cloning
	c = More(r.rule.Clone(), r.sep.Clone())
	r.Flag -= Cloning
	return c
}

// matchSep 匹配 sep, 调用条件 r.has(Finished)
func (r *more) matchSep(tok Token) (flag Flag) {
	const (
		sepMatched = Custom
	)
	flag = r.sep.Match(tok)

	if tok == eof {
		r.eof()
		return
	}

	switch flag {
	case 0:
		if r.Has(Standing) {
			flag = Finished
		}
	case Matched: // 需要更多 Token
		r.Flag = Finished | sepMatched

	case Matched | Standing:
		flag = Matched
		r.Flag = Finished | sepMatched
	case Matched | Finished:
		flag = Matched
		r.Flag = Standing

	case Finished:
	}

	return
}

func (r *more) Match(tok Token) (flag Flag) {
	var sep bool
	if r.Has(Handing) {
		return r.match(r, tok)
	}
	r.Flag |= Handing

	if r.Has(Custom) { // Custom 标记表示 sep 匹配中
		sep = true
		flag = r.matchSep(tok)
		if flag == 0 || tok == eof {
			r.Flag = 0
			return
		}
		if flag != Finished {
			return
		}
	}

	flag = r.rule.Match(tok)

	if tok == eof {
		r.matchSep(eof)
		return
	}

	// 从 rule 直接进入 sep 匹配必须具一定条件
	if !sep && (flag == 0 && r.Has(Standing) || flag == Finished) {
		sep = true
		flag = r.matchSep(tok)
		if flag == 0 { // 结束
			if r.Has(Standing) {
				flag = Finished
			}
			r.Flag = 0
			return
		}
		if flag != Finished { // sep 需要继续匹配
			return
		}
		// 无 sep , 纯 rule 重复, rule 应该已经被重置
		flag = r.rule.Match(tok)
	}

	// 这里全都是 rule 的结果
	switch flag {
	case 0:
		if sep || r.Has(Standing|Finished) {
			flag = Finished
		}
		r.Flag = 0
		r.rule.Match(eof)

	case Matched, Matched | Standing:
		r.Flag = r.Flag&0xf8 | flag&0x6

	case Matched | Finished:
		r.Flag = r.Flag&0xf8 | Standing
		flag = Matched | Standing

	case Finished:
		r.Flag = 0
	}

	if r.Has(Handing) {
		r.Flag -= Handing
	}
	return
}

type alternative struct {
	rules []Rule
	base
	pos int
}

// Any 产生任一匹配规则.
// 不要用 Any(rule, Term()) 替代 Option, 那会让 IsOption() 不可靠.
func Any(rule ...Rule) Rule {
	return &alternative{rules: rule}
}

func (r *alternative) IsOption() bool {
	for _, rule := range r.rules {
		if !rule.IsOption() {
			return false
		}
	}
	return true
}

func (r *alternative) Bind(rule Rule) {
	c, ok := rule.(*alternative)
	if !ok || r.Has(Handing) {
		panic("not support")
	}
	r.rules = make([]Rule, len(c.rules))
	for i := 0; i < len(r.rules); i++ {
		r.rules[i] = c.rules[i]
	}
}

func (c *alternative) Clone() Rule {
	if c.Has(Cloning) {
		return c
	}
	c.Flag += Cloning

	r := &alternative{}
	r.rules = make([]Rule, len(c.rules))
	for i := 0; i < len(r.rules); i++ {
		r.rules[i] = c.rules[i].Clone()
	}
	c.Flag -= Cloning
	return r
}

func (r *alternative) Match(tok Token) (flag Flag) {
	if r.Has(Handing) {
		return r.match(r, tok)
	}

	if r.pos >= len(r.rules) {
		r.pos = 0
	}
	r.Flag += Handing
	for r.pos < len(r.rules) {
		flag = r.rules[r.pos].Match(tok)
		if tok == eof {
			r.pos = 0
			r.eof()
			return
		}

		switch flag {
		case 0:
			if r.Has(Matched) {
				r.Flag = Handing
				r.pos = len(r.rules)
			}
		case Matched, Matched | Standing:
			r.Flag |= flag
		case Matched | Finished:
			r.Flag = Handing
			r.pos = len(r.rules)
		case Finished:
			if r.Has(Matched) {
				r.Flag = Handing
				break
			}
			r.pos++
			if r.pos == len(r.rules) {
				r.Flag = Handing
			}
			continue
			r.Flag = Handing
		}

		if r.Has(Matched) {
			break
		}
		r.pos++
	}
	r.Flag -= Handing
	return
}

type concatenation struct {
	rules []Rule
	base
	opt, pos int
}

// Seq 产生顺序匹配规则
func Seq(rule ...Rule) Rule {

	r := &concatenation{rules: rule}
	r.opt = len(r.rules)
	for i := len(r.rules) - 1; i >= 0; i-- {
		if r.rules[i].IsOption() {
			r.opt = i
		} else {
			break
		}
	}
	return r
}

func (r *concatenation) IsOption() bool {
	return false
}

func (r *concatenation) Bind(rule Rule) {
	c, ok := rule.(*concatenation)
	if !ok || r.Has(Handing) {
		panic("not support")
	}
	r.rules = make([]Rule, len(c.rules))
	opt := true
	r.opt = len(r.rules)
	for i := len(r.rules) - 1; i >= 0; i-- {
		r.rules[i] = c.rules[i]
		if opt && r.rules[i].IsOption() {
			r.opt = i
		} else {
			opt = false
		}
	}
}

func (c *concatenation) Clone() Rule {
	if c.Has(Cloning) {
		return c
	}
	c.Flag += Cloning
	r := &concatenation{opt: c.opt}
	r.rules = make([]Rule, len(c.rules))
	for i := 0; i < len(c.rules); i++ {
		r.rules[i] = c.rules[i].Clone()
	}
	c.Flag -= Cloning
	return r
}

func (r *concatenation) Match(tok Token) (flag Flag) {
	if r.Has(Handing) {
		return r.match(r, tok)
	}

	r.Flag += Handing

	if r.pos >= len(r.rules) {
		r.pos = 0
	}

	for r.pos < len(r.rules) {

		flag = r.rules[r.pos].Match(tok)

		if tok == eof {
			r.eof()
			r.pos = 0
			return
		}

		switch flag {
		case 0:
			r.pos = 0
			r.Flag = Handing
		case Matched:
		case Matched | Standing:
			if r.pos+1 < r.opt {
				flag = Matched
			}
			r.Flag = Handing | Standing
		case Matched | Finished:
			r.pos++
			if r.pos != len(r.rules) {
				if r.pos < r.opt {
					flag = Matched
				} else {
					flag = Matched | Standing
				}
			}
			r.Flag = Handing
		case Finished:
			r.pos++
			r.Flag = Handing
			continue
		}
		break
	}

	r.Flag -= Handing

	return
}
