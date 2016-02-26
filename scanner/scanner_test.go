package scanner_test

import (
	"fmt"
	"testing"

	"github.com/ZxxLang/zxx/scanner"
)

const code = `
a.b.c
1+2
a+b*c/d*=*e
a().b.c
`

func TestEcho(t *testing.T) {
	scan := scanner.New([]byte(code))
	code, ok := scan.Symbol()
	for ; ok; code, ok = scan.Symbol() {
		fmt.Println(code)
	}
}
