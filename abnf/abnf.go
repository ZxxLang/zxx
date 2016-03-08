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
	//  0
	// 	Matched
	// 	Matched|Standing
	// 	Matched|Finished
	// 	Finished
	//
	// 自动重置:
	//
	// 当状态为 0, Finished, Matched|Finished 的规则自动重置, 可接受新的匹配.
	//
	// EOF 重置: 当最终状态为
	//
	//	Matched 最终状态是不完整匹配.
	//	Matched|Standing 最终状态是完整匹配.
	//
	// 时使用 Match(EOF) 重置规则并返回 Finished, 其它状态不应该使用 EOF 来重置.
	//
	Match(tok Token) Flag

	// Bind 应只在递归嵌套规则变量时使用, 先声明规则变量后绑定规则.
	// 其他情况都不应该使用 Bind.
	Bind(Rule)

	// Clone 返回克隆规则, 这是深度克隆, 但不含递归.
	// 递归规则在 Match 中通过判断 Handing 标记及时建立的.
	Clone() Rule
}

// Flag 表示规则匹配 Token 后的状态
type Flag int

const (
	Matched  Flag = 1 << iota // 匹配成功, 消耗掉一个 Token
	Standing                  // 规则成立, 可以继续匹配 Token
	Finished                  // 规则完整, 不再接受匹配 Token

	// 下列标记由算法维护, Match 返回值不应该包含这些标记.
	Handing // 正在进行处理中(Match), 可用来检查发生递归
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

// end 在结束匹配时调用, 综合计算 f 和 flag, 返回值为:
//
//	0
//	Matched
//	Finished
//	Matched | Finished
//
func (f Flag) end(flag Flag) Flag {
	flag = f&6 | flag&7
	if flag.Has(Standing) {
		flag |= Finished
		flag ^= Standing
	}
	return flag
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
//
func Term(tok ...Token) Rule {
	return term{tok}
}

func (r term) Bind(Rule)   { panic("not support") }
func (r term) Clone() Rule { return r }
func (r term) Match(tok Token) Flag {
	for _, t := range r.toks {
		if t.Has(tok) {
			return Matched | Finished
		}
	}
	return 0
}

// base 维护 Flag 和递归 next.
type base struct {
	Flag
	next Rule
}

func (b *base) match(r Rule, tok Token) Flag {
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

func (r *option) Bind(rule Rule) { r.rule = rule }
func (r *option) Clone() (c Rule) {
	c = Option(r.rule.Clone())
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
		flag = r.Flag.end(0)
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

// more 的实现方法可能会改变, 不导出
type more struct {
	rule Rule
	base
	sep Rule
}

// More 产生重复匹配规则.
// rule 必须初次匹配成功, 然后当 rule 匹配结果为 0, Finished 时尝试 sep 匹配,
// 如果 sep 匹配成功可继续匹配 rule.
// 如果 sep 为 nil, 可重复匹配 rule.
//
func More(rule, sep Rule) Rule {
	return &more{rule: rule, sep: sep}
}

func (r *more) Bind(rule Rule) { r.rule = rule }
func (r *more) Clone() (c Rule) {
	if r.sep != nil {
		c = More(r.rule.Clone(), r.sep.Clone())
	} else {
		c = More(r.rule.Clone(), nil)
	}
	return
}

// matchSep 匹配 sep, 调用条件 r.has(Finished)
func (r *more) matchSep(tok Token) (flag Flag) {
	const (
		sepMatched  = Custom
		sepStanding = Custom << 1
	)
	if r.sep == nil {
		flag = Finished
	} else {
		flag = r.sep.Match(tok)
	}

	if tok == eof {
		r.eof()
		return
	}

	switch flag {
	case 0:
		if !r.Has(sepMatched) {
			flag = Finished
		}
		r.Flag -= Finished

	case Matched: // 需要更多 Token
		r.Flag |= sepMatched

	case Matched | Standing:
		r.Flag |= sepStanding
		if r.Has(sepMatched) {
			r.Flag -= sepMatched
		}

	case Matched | Finished:
		r.Flag -= Finished
		flag = Matched | Standing
		if r.Has(sepStanding) {
			r.Flag -= sepStanding
		}

	case Finished:
		r.Flag -= Finished
		if r.Has(sepStanding) {
			r.Flag -= sepStanding
		}
		if r.Has(sepMatched) {
			r.Flag -= sepMatched
		}
	}

	return
}

func (r *more) Match(tok Token) (flag Flag) {

	if r.Has(Handing) {
		return r.match(r, tok)
	}
	r.Flag |= Handing

	if r.Has(Finished) {
		flag = r.matchSep(tok)
		if tok == eof {
			r.eof()
			return
		}
		if flag != Finished {
			if r.Flag == Handing || flag.Has(Matched) {
				r.Flag -= Handing
				return
			}
		}
	}

	flag = r.rule.Match(tok)

	if tok == eof {
		r.matchSep(eof)
		r.eof()
		return
	}

	switch flag {
	case 0:
		if !r.Has(Standing) {
			r.Flag = Handing
			break
		}
		r.Flag |= Finished

		flag = r.matchSep(tok)
		if !flag.Has(Matched) {
			r.Flag = Handing
			flag = Finished
		}

	case Matched, Matched | Standing:
		r.Flag = r.Flag&0xf8 | flag
	case Matched | Finished:
		r.Flag = r.Flag&0xf8 | Finished | Standing
		flag = Matched | Standing
	case Finished:
		r.Flag = r.Flag&0xf8 | Finished | Standing
		flag = r.matchSep(tok)

		if !flag.Has(Matched) {
			r.Flag = Handing
			flag = Finished
		}
	}
	r.Flag -= Handing
	return
}

// alternative 规则, 最多吃掉一个 Token
type alternative struct {
	rules []Rule
	base
	pos int
}

// Any 产生任一匹配规则.
func Any(rule ...Rule) Rule {
	return &alternative{rules: rule}
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
	r := &alternative{}
	r.rules = make([]Rule, len(c.rules))
	for i := 0; i < len(r.rules); i++ {
		r.rules[i] = c.rules[i].Clone()
	}
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
				r.pos = len(r.rules)
			}
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

// concatenation
type concatenation struct {
	rules []Rule
	base
	pos int
}

// Seq 产生顺序匹配规则
func Seq(rule ...Rule) Rule {
	return &concatenation{rules: rule}
}

func (r *concatenation) Bind(rule Rule) {
	c, ok := rule.(*concatenation)
	if !ok || r.Has(Handing) {
		panic("not support")
	}
	r.rules = make([]Rule, len(c.rules))
	for i := 0; i < len(r.rules); i++ {
		r.rules[i] = c.rules[i]
	}
}

func (c *concatenation) Clone() Rule {
	r := &concatenation{}
	r.rules = make([]Rule, len(c.rules))
	for i := 0; i < len(r.rules); i++ {
		r.rules[i] = c.rules[i].Clone()
	}
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
			r.pos = 0
			r.eof()
			return
		}

		switch flag {
		case 0:
			r.pos = 0
		case Matched:
		case Matched | Standing:
			if r.pos+1 != len(r.rules) {
				flag = Matched
			}
			r.Flag = Handing | Standing
		case Matched | Finished:
			r.pos++
			if r.pos != len(r.rules) {
				flag = Matched
			}
			r.Flag = Handing
		case Finished:
			r.pos++
			if r.Has(Standing) {
				r.Flag -= Standing
				continue
			}
		}
		break
	}

	r.Flag -= Handing

	return
}
