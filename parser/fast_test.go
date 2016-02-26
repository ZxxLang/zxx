package parser_test

import (
	"fmt"
	"testing"

	"github.com/ZxxLang/zxx/parser"
)

func Test_Fast(t *testing.T) {
	nodes, err := parser.Fast([]byte(code))
	if err != nil {
		t.Fatal(err)
	}

	for _, n := range nodes {
		fmt.Println(n.Pos, n.Tok, n.Source)
	}
}

const code = `
this is test

pub
	type Block
		use Base
		[File] Files


`
