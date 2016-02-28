// Copyright 2016 The Zxx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import "unicode/utf8"

// scanner 每次扫描一个字符.
type scanner struct {
	src []byte // source

	size   int
	offset int // 当前扫描所处 src 的字节偏移量
	nl     uint16
}

type Pos int

// Offset 返回 pos 位置偏移 offset 个字节的 Pos 值.
// 注意参数 offset 是有方向
func (pos Pos) Offset(offset int) Pos {
	return pos + Pos(offset)
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

func (s *scanner) Pos() Pos {
	return Pos(s.offset)
}

// IsEOF 返回当前扫描位置是否已经到 EOF
// 不恰当的使用 IsEOF 的话
// Rune()   得不到  0, 0
// Symbol() 得不到 "", true
func (s *scanner) IsEOF() bool {
	return s.offset >= s.size
}

// Eol 返回 uint16 表示的换行符
//	0   未确定
//	10  LF   风格 "\n"   0xa
//	13  CR   风格 "\r"   0xd
//	218 CRLF 风格 "\r\n" 0xda
func (s *scanner) Eol() uint16 {
	return s.nl
}

// Rune 返回当前位置的 rune 和 UTF-8 编码长度, 并前进 size 个字节.
// 如果 r == 0 && size == 0 , 表示扫描结束
func (s *scanner) Rune() (r rune, size int) {
	if s.offset >= s.size {
		return
	}

	r, size = utf8.DecodeRune(s.src[s.offset:])
	if r == utf8.RuneError {
		return
	}
	s.offset += size
	return
}

// Symbol 返回易于识别的字符串, 并前进.
// 如果 ok == false, 表示编码错误
//
// symbol 值:
//
//	"" 表示扫描结束
//	连续的			' ', '\t', '/', '-', '+', '\n', '\r', '\r\n'
//	两个字符		运算符, 操作符
//	单个字符		运算符, 定界符, 单双引号
//  多字节字符		直到行尾
//	连续的字符		直到空白, 运算符, 操作符, 定界符, 换行
//  连续的成员		a.b.c
func (s *scanner) Symbol() (symbol string, ok bool) {
	offset := s.offset
	r, size := s.Rune()

	if size == 0 {
		ok = r == 0
		return
	}
	ok = true

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

	switch c {
	case '\r', '\n': // 连续的换行
		nc := s.src[s.offset]
		// 识别换行风格
		if s.nl == 0 {
			s.nl = uint16(c)
			if nc != c && (nc == '\r' || nc == '\n') {
				s.nl = s.nl<<8 + uint16(nc)
			}
		}

		for s.offset != s.size && (nc == '\r' || nc == '\n') {
			nc = s.src[s.offset]
			s.offset++
		}

		symbol = string(s.src[offset:s.offset])

	case ' ', '\t': // 连续的

		for s.offset != s.size && s.src[s.offset] == c {
			s.offset++
		}

		symbol = string(s.src[offset:s.offset])

	case '-', '/': // 可后跟 '=' 或者连续多个的
		if s.src[s.offset] == '=' {
			s.offset++
		} else {
			for s.offset != s.size && s.src[s.offset] == c {
				s.offset++
			}
		}
		symbol = string(s.src[offset:s.offset])

	case '&', '|', '!', '~', '*', '=': // 可后跟 '='
		if s.src[s.offset] == '=' {
			s.offset++
		}
		symbol = string(s.src[offset:s.offset])

	case '>', '<', '+': // 可后跟 '=', 或者重复一个
		if s.src[s.offset] == '=' || s.src[s.offset] == c {
			s.offset++
		}
		symbol = string(s.src[offset:s.offset])

	case ',', '"', '\'', '{', '}', '(', ')', '[', ']': // 单个
		symbol = string(s.src[offset:s.offset])
	default:
		// 不严格的判断 integer, float, datetime, 标识符
		var num byte
		if c >= '0' && c <= '9' {
			num = 3 // 三种可能 integer, float, datetime
		} else if c >= 'a' && c <= 'z' || c >= 'Z' && c <= 'Z' || c == '_' {
			num = 1 // 标识符
		}

		for s.offset != s.size {

			switch s.src[s.offset] {
			case '.':
				if num == 3 {
					num = 'f' // float
					s.offset++
					continue
				}
				if num != 1 {
					break
				}
				// 成员写法
				s.offset++
				if s.offset != s.size {
					c = s.src[s.offset]
					if !(c >= 'a' && c <= 'z' || c >= 'Z' && c <= 'Z' || c == '_') {
						s.offset--
						break
					}
				}
				continue
			case ':':
				if num == 3 {
					num = 'd' // datetime
					s.offset++
					continue
				}
			case '+':
				if num == 'd' || num == 'f' {
					s.offset++
					continue
				}
			default:
				c = s.src[s.offset]
				if c < '0' ||
					c > '9' && c < 'A' ||
					c > 'Z' && c < 'a' ||
					c > 'z' {
					break
				}
				s.offset++
				continue
			}
			break
		}
		symbol = string(s.src[offset:s.offset])
	}

	return
}

// Tail 返回当前位置到行尾的字符串.
// 如果参数 nl 为 true, 返回字符讲包含连续的换行符.
// 如果当前位置已经是换行, 那么会返回 "".
// 该方法不检查非法 UTF-8 编码, 下一个符号将是换行符或 EOF.
func (s *scanner) Tail(nl bool) string {
	offset := s.offset

	for s.offset < s.size && s.src[s.offset] != '\n' && s.src[s.offset] != '\r' {
		s.offset++
	}

	for nl && (s.src[s.offset] == '\n' || s.src[s.offset] == '\r') {
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
