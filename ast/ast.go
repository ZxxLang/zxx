// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

type (
	// Stmt 表示语句
	Stmt struct {
	}

	// Text 表示文本字符串字面值
	Text struct {
	}

	// Decl 表示声明语句
	Decl struct {
	}

	// Exec 表示执行语句(过程调用)
	Exec struct {
	}

	// Pred 表示预定义语句
	Pred struct {
	}

	// Expr 表示表达式
	Expr struct {
	}

	// 缩进
	Ident struct {
		Name string
	}

	Node interface {
	}

	Block struct {
		List []Node
	}
)
