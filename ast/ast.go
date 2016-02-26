// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Zxx 的 AST 是 flat-AST, 节点依次保存为数组, 节点关系在节点 Flag 中.
// 构建过程通过 token, scanner, parser, ast 多包配合逐步分析转换 Token 的过程,
// 该过程不是严格的 Token 筛选, 最终执行节点的 Final() 做最后的合法性验证.
//
package ast

import (
	"errors"

	"github.com/ZxxLang/zxx/scanner"
	"github.com/ZxxLang/zxx/token"
)

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

// Flag 表示语义节点具有的属性
type Flag uint

const (
	// 节点标记, 每个节点必须包含一项
	FBlock       Flag = 1 << iota // Block 节点
	FFile                         // File 节点
	FDeclaration                  // Decl 节点
	FChunk                        // Chunk 节点
	FStatement                    // Stmt 节点
	FExpression                   // Expr 节点
	FText                         // Text 节点

	// 节点状态
	FFinal // 是否已经完整并闭合

	// 语法风格标记
)

// 风格 mask 继承
const StyleMask = 0xFFFFFF00

type (
	Node interface {

		// Kind 返回该节点的语义标识
		Kind(mask Flag) Flag

		// Id 返回该节点在根节点中的序号, 从 0 开始.
		Id() int

		// Prev 返回该节点的上个非占位节点
		// 如果到了顶层, 返回 nil
		Prev() Node

		// Final 设置 FFinal 状态
		Final() error

		// 节点对应的 Token 来自 Resolve 的参数, 但有可能被转化.
		// 如果该值为 token.EOF, 表示该节点是个节点容器, 比如 File, Block.
		// 没有真的 EOF 节点.
		Token() token.Token

		// Text 返回该节点源代码.
		Text() string

		resolve(*Base)
	}

	// Base 是所有 Node 的共性结构
	Base struct {
		Flag

		Tok token.Token

		// Index 是所属上层节点子节点序号
		Index int

		// Source 是节点源码
		Source string

		Pos scanner.Pos

		// 上个节点序号
		prev int

		all *File
	}

	// Block 对应 Zxx 块, 类似其它语言的 package, 由 File 组成.
	// 事实上 Block 就是 File 容器和作用域
	Block struct {
		Base
		Files []File
	}

	// File 对应 Zxx 源码文件. 属性 Nodes 只包括 Text 占位节点, Decl 声明节点.
	// File 负责解决 SPACES,TABS,NL,EOF
	File struct {
		Base
		Nodes  []Node
		Active Node // 活动节点
		Last   Node // 最后的节点
		expect rule
	}

	// Decl 可包括声明语句 IsDeclare.
	// 可后跟: IsDeclare, LEFT, IDENT
	// 特别处理: USE
	Decl struct {
		Base
	}

	// Chunk 视觉上的块, Token 只能是 Left 或者 Right
	Chunk struct {
		Base
	}

	// Stmt 表示非 Decl 外的语句
	// 可包括: BREAK .. SWITCH,
	// 可后跟: Text, Expr(经 Text 转换而来)
	Stmt struct {
		Base
	}

	// Expr 表示表达式, 表达式总是从 Text 转换而来.
	// 一个独立的 Text 也是 Expr. 转换由 File.Add 完成
	Expr struct {
		Base
	}

	// Text
	Text struct {
		Base
	}
)

// ------------------- Base -------------------

func (b Base) Kind(mask Flag) Flag {
	if mask == 0 {
		return b.Flag
	}
	return b.Flag & mask
}

func (b Base) Id() int            { return b.Index }
func (b Base) Token() token.Token { return b.Tok }
func (b Base) Text() string       { return b.Source }

func (b Base) Prev() Node {
	if b.Index == 0 {
		return nil
	}

	n := b.all.Nodes[b.Index-1]
	tok := n.Token()
	for n.Id() > 0 && tok > token.IDENT {
		n = b.all.Nodes[n.Id()-1]
	}
	if n.Id() == 0 {
		return nil
	}

	return n
}

// ------------------- File -------------------

func NewFile() (file *File) {
	file = new(File)
	file.Flag = FFile
	file.Tok = token.EOF
	file.Nodes = make([]Node, 0, 1024)
	file.Nodes = append(file.Nodes, file)
	file.Last = file
	file.Active = file
	return
}

func (b *File) Len() int { return len(b.Nodes) }

// file.add 做最后的检查, 并根据 Token, Flag 设置 Last, Active.
func (b *File) add(base Base) (err error) {
	var n, part Node

	// NL 自动 Final
	if base.Tok == token.NL {
		last := b.Active.Token()
		if last == token.RIGHT || last == token.IDENT ||
			last.As(token.Literal) {

			err = b.Active.Final()
			if err != nil {
				return err
			}
			base.Flag |= FFinal
		}
	}

	base.Index = b.Len()
	base.all = b
	base.prev = b.Active.Id()
	base.Flag |= b.Active.Kind(StyleMask)

	switch base.Flag & 0x7F {
	case FDeclaration:
		n = &Decl{base}
	case FChunk:
		n = &Chunk{base}
	case FStatement:
		n = &Stmt{base}
	case FExpression:
		n = &Expr{base}
	case FText:
		n = &Text{base}
	}

	if n == nil || n.Token() > token.PLACEHOLDER {
		err = errors.New("ast: Oop! invalid Base")
		return
	}

	b.Last = n
	b.Nodes = append(b.Nodes, n)

	// 分组自动结束
	if base.Tok == token.RIGHT {
		for part = b.Prev(); part != nil; part = part.Prev() {
			if part.Token() == token.LEFT && part.Kind(FFinal) == 0 {
				if err = n.Final(); err == nil {
					err = part.Final()
				}
				if err != nil {
					return
				}
				break
			}
		}

		if part.Kind(FFinal) == 0 {
			return errors.New("ast: Oop! Unpaired LEFT and RIGHT")
		}

		if b.Active = part.Prev(); b.Active == nil {
			b.Active = b
			return
		}
		// 需要检查 =LEFT 的右侧
		return
	}

	return
}

// ------------------- Final ------------------

func (b *File) Final() error {
	b.Flag |= FFinal
	return nil
}

func (b *Decl) Final() error {
	b.Flag |= FFinal
	return nil
}

func (b *Chunk) Final() error {
	b.Flag |= FFinal
	return nil
}

func (b *Text) Final() error {
	b.Flag |= FFinal
	return nil
}

func (b *Stmt) Final() error {
	b.Flag |= FFinal
	return nil
}

func (b *Expr) Final() error {
	b.Flag |= FFinal
	return nil
}

// 解决所有的 PLACEHOLDER, INDENTATION, COMMENT, NL, EOF.
// 并合并多个空行为 EMPTYLINE

// File.Push 接收扫描到的 Token,
func (b *File) Push(pos scanner.Pos, tok token.Token, code string) error {
	var flag Flag
	switch tok {
	case token.NL:
		flag = FText

		last := b.Last.Token()

		if last == token.NL {
			tok = token.EMPTYLINE // 转
			break
		}

		if last == token.INDENTATION ||
			last == token.EMPTYLINE {

			text := b.Last.(*Text)
			text.Source += code
			text.Tok = token.EMPTYLINE // 合并
			return nil
		}

		last = b.Active.Token()
		if last == token.RIGHT || last == token.IDENT ||
			last.As(token.Literal) {

			err := b.Active.Final()
			if err != nil {
				return err
			}
			flag |= FFinal
		}

		// 识别 Python 缩进风格

	case token.RIGHT:
	case token.COMMENT, token.COMMENTS:
		flag = FText

		if b.Last.Token() == token.EMPTYLINE ||
			b.Last.Token() == token.PLACEHOLDER {
			tok = token.PLACEHOLDER // 转
			break
		}
		// 先不使用合并
		tok = token.COMMENT

	case token.PLACEHOLDER, token.INDENTATION:
		// 统一处理右括号闭合
		flag = FText
	default:
		if tok > token.PLACEHOLDER || b.Active.Kind(FFile|FBlock|FText) != 0 {
			return errors.New("ast: Oop! invalid " + tok.String())
		}

		if b.expect != nil {
			// expect 尝试生成 Text 节点, 并且 tok 不变
			if _, pass := b.expect.pass(tok); pass {
				flag = FText
				b.expect = nil
				break
			}
			b.expect = nil
		}

		base := Base{
			Tok:    tok,
			Pos:    pos,
			Source: code,
		}

		b.Active.resolve(&base)
		if base.Flag == 0 {
			return errors.New("ast: Oop! invalid " + tok.String())
		}
		return b.add(base)

	case token.EOF:
		// 最后的闭合检查
	}

	if flag == 0 {
		return errors.New("ast: Oop! invalid " + tok.String())
	}

	return b.add(Base{
		Flag:   flag,
		Tok:    tok,
		Pos:    pos,
		Source: code,
	})

}

// ------------------- Resolve -------------------

// File 解决顶层声明
func (b *File) resolve(base *Base) {
	var flag Flag
	switch base.Tok {
	case token.USE:
		flag = FDeclaration
	case token.FUNC:
		flag = FDeclaration
	case token.PROC:
		flag = FDeclaration
	case token.CONST:
		flag = FDeclaration
	case token.VAR:
		flag = FDeclaration
	case token.TYPE, token.STATIC:
		flag = FDeclaration
	case token.PUB:
		flag = FDeclaration
	}
	base.Flag = flag
	return
}

// Text 只可能是用于声明的 IDENT, Type, OUT
func (b *Text) resolve(base *Base) {
	return
}

// Decl.Resolve 是所有干净节点的入口.
func (b *Decl) resolve(base *Base) {

	var flag Flag
	tok := base.Tok

	switch b.Tok {
	case token.USE:
		switch tok {
		case token.LEFT:
			flag = FChunk
			b.all.expect = term{token.IDENT, token.VALSTRING}
		case token.IDENT:
			flag = FText
			b.all.expect = term{token.VALSTRING}
		case token.VALSTRING:
			flag = FText | FFinal
		}

	case token.CONST:
		switch tok {
		case token.LEFT:
			flag = FChunk
			b.all.expect = term{token.IDENT, token.Type}
		case token.IDENT, token.Type:
			flag = FText
			b.all.expect = term{token.Assign}
		}

	case token.VAR:
		switch tok {
		case token.LEFT:
			flag = FChunk
			b.all.expect = term{token.IDENT, token.MEMBER, token.Type}
		case token.IDENT, token.Type:
			flag = FText
			b.all.expect = term{token.IDENT}
		}

	case token.STATIC:
		switch tok {
		case token.LEFT:
			flag = FChunk
			b.all.expect = term{token.IDENT, token.MEMBER, token.Type}
		case token.FUNC, token.PROC:
			flag = FDeclaration
		case token.IDENT, token.Type:
			flag = FText
			b.all.expect = term{token.IDENT}
		}
	case token.TYPE:
		switch tok {
		case token.LEFT:
			flag = FChunk
			b.all.expect = term{token.IDENT}
		case token.IDENT:
			flag = FText
			b.all.expect = term{token.LEFT}
		}

	case token.PROC:
		switch {
		case tok == token.OUT: // 匿名过程
			flag = FText
			b.all.expect = term{token.IDENT, token.MEMBER, token.FUNC}
		case tok == token.IDENT || tok == token.MEMBER || tok.As(token.Type):
			flag = FText
			b.all.expect = term{token.IDENT}
		}

	case token.FUNC:
		switch {
		case tok == token.OUT: // 匿名过程
			flag = FText
			b.all.expect = term{token.IDENT}
		case tok == token.IDENT || tok == token.MEMBER || tok.As(token.Type):
			flag = FText
			b.all.expect = term{token.OUT, token.IDENT}
		}
	case token.PUB:
		// 暂时不支持分组
		if tok != token.PUB && tok.As(token.Declare) {
			flag = FDeclaration
			b.all.expect = term{token.TYPE, token.PROC, token.FUNC,
				token.VAR, token.CONST, token.STATIC}
		}
	}
	base.Flag = flag
	return
}

func (b *Chunk) resolve(base *Base) {
	return
}

func (b *Stmt) resolve(base *Base) {
	return
}

func (b *Expr) resolve(base *Base) {
	return
}
