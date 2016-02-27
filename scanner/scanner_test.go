package scanner_test

import (
	"testing"

	"github.com/ZxxLang/zxx/scanner"
)

type seq []string

var good []seq = []seq{
	seq{
		`use a 'b'`,
		`use`, ` `, `a`, ` `, `'`, `b`, `'`,
	},
}

func Test_eq(t *testing.T) {
	for _, ss := range good {
		scan := scanner.New([]byte(ss[0]))
		for _, s := range ss[1:] {
			code, ok := scan.Symbol()
			if !ok || s != code {
				t.Fatal(ok, s, code)
			}
		}
	}
}
