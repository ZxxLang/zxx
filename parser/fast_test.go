package parser_test

import (
	"testing"

	"github.com/ZxxLang/zxx/parser"
)

var good [][]string = [][]string{
	[]string{
		`use a 'b'`,
		`use`, `a`, `'b'`,
	},
	[]string{
		`use (a 'b')`,
		`use`, `(`, `a`, `'b'`, `)`,
	},
	[]string{
		`use( a 'b') `,
		`use`, `(`, `a`, `'b'`, `)`,
	},
	[]string{
		`use( a'b' )`,
		`use`, `(`, `a`, `'b'`, `)`,
	},
}

func Test_eq(t *testing.T) {
	for _, ss := range good {
		nodes, err := parser.Fast([]byte(ss[0]), nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(nodes) != len(ss)-1 {
			t.Fatal(nodes)
		}
		ss = ss[1:]

		for i, n := range nodes {
			if ss[i] != n.Source {
				t.Fatal(n)
			}
		}
	}
}

const code = `
this is test

pub
	type Block
		use Base
		[File] Files


`
