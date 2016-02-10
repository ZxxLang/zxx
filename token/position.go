// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import "fmt"

const (
	// 换行风格
	LF   = "\n"
	CR   = "\r"
	CRLF = "\r\n"
)

// Position 表示源码中的一个抽象位置.
type Position struct {
	Offset int // 字节偏移量, 从 0 开始
	Line   int // 行号, 从 1 开始
	Column int // 列号, 从 1 开始, 以字节为单位
}

// IsValid 返回位置是否有效. 通常意味着是否明确了行号.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

// String 返回该位置的字符串描述. 格式为:
//
//	file:line:column    三段式, 如果有文件名
//	line:column         两段式, 如果无文件名
//	file                有文件名且位置无效
//	-                   无文件名且位置无效
func (pos Position) String(filename string) string {
	if pos.IsValid() {
		if filename != "" {
			filename += ":"
		}
		filename += fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	}
	if filename == "" {
		filename = "-"
	}
	return filename
}

// Block 记录包粒度代码块中的文件(包括 hack 覆盖)关系.
type Block struct {
	files []File
	uri   string // use 语句中的 block 完全路径
}

// File 记录文件级别的 block 信息
type File struct {
	path  string // 抽象的文件路径
	nl    string // 换行符
	index int    // 在 Block.files 中的序号
}

// Code 记录代码块级别的 block 信息.
// 代码块是个抽象 block, 具体粒度根据需求变化.
type Code struct {
	pos, end Position // 开始和结束位置
	index    int      // 所属 Block.files 中的序号
}
