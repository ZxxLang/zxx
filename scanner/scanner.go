// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import "unicode/utf8"

// scanner 每次扫描一个字符.
type scanner struct {
	src []byte // source

	size   int
	offset int  // 当前扫描所处 src 的字节偏移量
	nl     byte // 换行风格 0 未知, '\n' LF, ' \r' CR, 其他 CRLF
}

// New 返回一个 scanner 并进行一些前期处理.
// 如果有 BOM 头则移动当前位置到 BOM 之后, scanner.Pos() 一定为 3.
func New(source []byte) (scan *scanner) {
	scan = &scanner{src: source, size: len(source)}

	// BOM 0xFEFF
	if len(source) > 2 && source[0] == 0xef && source[1] == 0xbb && source[2] == 0xbf {
		scan.offset = 3
	}
	return scan
}

func (s *scanner) Pos() int {
	return s.offset
}

// Rune 返回当前位置的 rune 和 UTF-8 编码长度, 并前进 size 个字节.
// 如果 size  < 0, 表示编码错误
// 如果 size == 0, 表示扫描结束
func (s *scanner) Rune() (r rune, size int) {
	if s.offset >= s.size {
		return
	}

	r, size = utf8.DecodeRune(s.src[s.offset:])
	if r == utf8.RuneError {
		size = -1
		return
	}
	s.size += size
	return
}

// Symbol 返回易于解析的 zxx 语义符号, 并前进.
// 如果 ok == false, 表示编码错误
// symbol:
//	'' 表示扫描结束
//	换行符 '\n', '\r', '\r\n'
//	两个以上的 ' ', '\t', '/','--','++'
//	[a-zA-Z_]+[a-zA-Z_0-9]*
//	[0-9]+[a-zA-Z_0-9\.]*
//	两字符的 Zxx 运算符, 操作符
//	单字符的 Zxx 运算符, 定界符, 单双引号
//  多字节字符直到行尾
//
// 注意 symbol 未考虑上下文结合性. 合理的使用 symbol 才有意义.
func (s *scanner) Symbol() (symbol string, ok bool) {
	offset := s.offset
	r, size := s.Rune()

	ok = size >= 0
	if !ok || r == 0 {
		return
	}

	if s.offset >= s.size {
		symbol = string(r)
		return
	}

	// 多字节, 直接到行尾
	if r > 127 {
		for s.offset != s.size && s.src[s.offset] != '\n' && s.src[s.offset] != '\r' {
			s.offset++
		}
		symbol = string(s.src[offset:s.offset])
		return
	}

	c := byte(r)

	nc := s.src[s.offset]

	switch c {
	case '\r', '\n': // 识别换行风格
		if s.nl == 0 {
			s.nl = nc
		}
		if c != nc && s.nl == nc {
			s.offset++
			symbol = string(s.src[offset:s.offset])
		} else {
			symbol = string(c)
		}

	case ' ', '\t': // 连续的

		for s.offset != s.size && s.src[s.offset] == c {
			s.offset++
		}

		symbol = string(s.src[offset:s.offset])

	case '-', '/': // 可后跟 '=' 或者连续多个的
		if nc == '=' {
			s.offset++
		} else {
			for s.offset != s.size && s.src[s.offset] == c {
				s.offset++
			}
		}
		symbol = string(s.src[offset:s.offset])

	case '&', '|', '!', '~', '=': // 可后跟 '='
		if nc == '=' {
			s.offset++
		}
		symbol = string(s.src[offset:s.offset])

	case '>', '<', '+': // 可后跟 '=', 或者重复一个
		if nc == '=' || nc == c {
			s.offset++
		}
		symbol = string(s.src[offset:s.offset])

	case '"', '\'', '{', '}', '(', ')', '[', ']': // 单个
		symbol = string(c)
	default: // 直到非 [a-zA-Z_0-9] 结束
		// [a-zA-Z_]+[a-zA-Z_0-9]*
		// [0-9]+[a-zA-Z_0-9\.]*
		for s.offset != s.size &&
			(nc >= 'a' && nc <= 'z' || nc >= 'A' && nc <= 'Z' ||
				nc >= '0' && nc <= '9' || nc == '_' || nc == '.') {

			s.offset++
			nc = s.src[s.offset]
		}
		symbol = string(s.src[offset:s.offset])
	}

	return
}

// Tail 返回当前位置到行尾的字符串.
// 如果返回值为 "" 表示扫描结束
// 该方法不检查非法 UTF-8 编码, 下一个符号将是换行符或 EOF.
func (s *scanner) Tail() string {
	offset := s.offset

	for s.offset < s.size && s.src[s.offset] != '\n' && s.src[s.offset] != '\r' {
		s.offset++
	}
	return string(s.src[offset:s.offset])
}

// EndString 返回当前位置到 escape 指示的 Zxx 的字符串结尾.
// escape 为 true 表示支持 "\" 逃逸的双引号结尾的字符串.
// 否则表示单引号结尾字符串.
// 目前只支持简单的单个字符逃逸
func (s *scanner) EndString(escape bool) string {
	offset := s.offset
	for ; s.offset < s.size; s.offset++ {
		if escape && s.src[s.offset] == '\\' {
			s.offset++
			continue
		}
		if escape && s.src[s.offset] == '"' ||
			!escape && s.src[s.offset] == '\'' {
			s.offset++
			break
		}
	}
	return string(s.src[offset:s.offset])
}
