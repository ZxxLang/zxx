// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"errors"

	"github.com/ZxxLang/zxx/ast"
	"github.com/ZxxLang/zxx/scanner"
	"github.com/ZxxLang/zxx/token"
)

// Parse 解析 src 到 ast
func Parse(src []byte) (blk *ast.Block, err error) {
	var node ast.Node

	scan := scanner.New(src)

	List := make([]ast.Node, 0, 1024)

	for {
		sym, ok := scan.Symbol()

		if !ok {
			err = errors.New("invalid UTF-8 encode")
			return
		}

		if sym == "" {
			break
		}

		tok := token.Lookup(sym)

		switch {

		case tok == token.SPACES:

		case tok == token.TABS:

		case tok == token.NL: // new line

		case tok == token.COMMENT: // 注释

		case tok == token.COMMENTS:

		case tok == token.LITERAL: // 字面值

		case tok.IsDeclare():

		case tok.IsKeyword(): // 声明之外的保留字

		case tok.IsOperator():

		case tok.IsAssign(): // '=','++','--'

		}
	}
	return
}
