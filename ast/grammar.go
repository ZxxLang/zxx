package ast

import "github.com/ZxxLang/zxx/token"

// 词法轻规则, 配合些许 Resolve 代码完成构建

type rule interface {
	pass(token.Token) (eat int, pass bool)
}

// term 包装 Token, Alternative 规则
type term []token.Token

func (r term) pass(tok token.Token) (int, bool) {
	for _, t := range r {
		if token.Token(t).As(tok) {
			return 1, true
		}
	}
	return 0, false
}

// 零次或一次
type zero term

func (r zero) pass(tok token.Token) (int, bool) {
	eat, pass := rule(r).pass(tok)
	if pass && eat != 0 {
		return 1, true
	}
	return 0, true
}

// 一次或多次
type more term

func (r more) pass(tok token.Token) (int, bool) {
	_, pass := rule(r).pass(tok)
	if !pass {
		return 0, false
	}
	return -1, true
}

// queue 执行 Concatenation 规则, 顺序通过匹配
type queue struct {
	rules []rule
	at    int
}

func (q *queue) pass(tok token.Token) (int, bool) {
	if q.at >= len(q.rules) {
		q.at = 0
	}

	for q.at < len(q.rules) {
		eat, pass := q.rules[q.at].pass(tok)
		if !pass {
			return 0, false
		}
		switch eat {
		case 0:
			q.at++
			continue
		case 1:
			q.at++
			return 1, true
		}
	}
	return 0, true
}
