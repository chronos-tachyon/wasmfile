package wat

import (
	"io/fs"
	"path"
	"testing"
)

func XN(v ...*Node) *Node {
	return &Node{Type: ExprNode, Value: v}
}

func KN(str string) *Node {
	return &Node{Type: KeywordNode, Value: str}
}

func IN(str string) *Node {
	return &Node{Type: IdentifierNode, Value: str}
}

func SVN(str string) *Node {
	return &Node{Type: StringNode, Value: str}
}

func NVN(bits NumFlags, v ...string) *Node {
	return &Node{Type: NumberNode, Value: N(bits, v...)}
}

func TestParse(t *testing.T) {
	type TestCase struct {
		Name   string
		Expect *Node
	}

	testCases := [...]TestCase{
		{
			Name: "file1.wat",
			Expect: XN(
				XN(KN("module")),
			),
		},
		{
			Name: "file2.wat",
			Expect: XN(
				XN(KN("module")),
			),
		},
		{
			Name: "file3.wat",
			Expect: XN(
				XN(
					KN("module"),
					XN(
						KN("import"),
						SVN("wasi_unstable"),
						SVN("fd_write"),
						XN(
							KN("func"),
							IN("$fd_write"),
							XN(KN("param"), KN("i32"), KN("i32"), KN("i32"), KN("i32")),
							XN(KN("result"), KN("i32")),
						),
					),
					XN(
						KN("memory"),
						NVN(0, "1"),
					),
					XN(
						KN("export"),
						SVN("memory"),
						XN(KN("memory"), NVN(0, "0")),
					),
					XN(
						KN("data"),
						XN(KN("i32.const"), NVN(0, "8")),
						SVN("hello world\n"),
					),
					XN(
						KN("func"),
						IN("$main"),
						XN(KN("export"), SVN("_start")),
						XN(
							KN("i32.store"),
							XN(KN("i32.const"), NVN(0, "0")),
							XN(KN("i32.const"), NVN(0, "8")),
						),
						XN(
							KN("i32.store"),
							XN(KN("i32.const"), NVN(0, "4")),
							XN(KN("i32.const"), NVN(0, "12")),
						),
						XN(
							KN("call"),
							IN("$fd_write"),
							XN(KN("i32.const"), NVN(0, "1")),
							XN(KN("i32.const"), NVN(0, "0")),
							XN(KN("i32.const"), NVN(0, "1")),
							XN(KN("i32.const"), NVN(0, "20")),
						),
						KN("drop"),
					),
				),
			),
		},
		{
			Name: "strings.wat",
			Expect: XN(
				SVN("this is a string"),
				SVN("tab: \t"),
				SVN("newline: \n"),
				SVN("crlf: \r\n"),
				SVN(`backslash: \`),
				SVN(`single quote: '`),
				SVN(`double quote: "`),
				SVN("ESC[0m: \x1b[0m"),
				SVN("smiley: \u263a\ufe0f"),
				SVN("smiley: \u263a\ufe0f"),
			),
		},
		{
			Name: "numbers.wat",
			Expect: XN(
				NVN(0, "0"),
				NVN(0, "123"),
				NVN(FlagSign, "0"),
				NVN(FlagSign, "123"),
				NVN(FlagSign|FlagNeg, "0"),
				NVN(FlagSign|FlagNeg, "123"),
				NVN(FlagFloat, "0", "0"),
				NVN(FlagFloat, "0", "123"),
				NVN(FlagFloat|FlagSign, "0", "0"),
				NVN(FlagFloat|FlagSign, "0", "123"),
				NVN(FlagFloat|FlagSign|FlagNeg, "0", "0"),
				NVN(FlagFloat|FlagSign|FlagNeg, "0", "123"),
				NVN(FlagHex, "0"),
				NVN(FlagHex, "a"),
				NVN(FlagHex|FlagSign, "0"),
				NVN(FlagHex|FlagSign, "a"),
				NVN(FlagHex|FlagSign|FlagNeg, "0"),
				NVN(FlagHex|FlagSign|FlagNeg, "a"),
				NVN(FlagFloat|FlagHex, "a", "b"),
				NVN(0, "123"),
				NVN(FlagFloat, "123", "456"),
				NVN(FlagFloat|FlagHex, "123", "456"),
				NVN(FlagFloat|FlagInf),
				NVN(FlagFloat|FlagInf|FlagSign),
				NVN(FlagFloat|FlagInf|FlagSign|FlagNeg),
				NVN(FlagFloat|FlagNaN),
				NVN(FlagFloat|FlagNaN|FlagSign),
				NVN(FlagFloat|FlagNaN|FlagSign|FlagNeg),
				NVN(FlagFloat|FlagNaN|FlagAcanonical|FlagHex, "deadbeef"),
				NVN(FlagFloat|FlagNaN|FlagAcanonical|FlagHex|FlagSign, "deadbeef"),
				NVN(FlagFloat|FlagNaN|FlagAcanonical|FlagHex|FlagSign|FlagNeg, "deadbeef"),
				NVN(FlagFloat, "0", "5", "20"),
				NVN(FlagFloat|FlagExpSign, "0", "5", "20"),
				NVN(FlagFloat|FlagExpSign|FlagExpNeg, "0", "5", "20"),
				NVN(FlagFloat|FlagHex, "0", "8", "12"),
				NVN(FlagFloat|FlagHex|FlagExpSign, "0", "8", "12"),
				NVN(FlagFloat|FlagHex|FlagExpSign|FlagExpNeg, "0", "8", "12"),
			),
		},
	}

	var p Parser
	for _, row := range testCases {
		t.Run(row.Name, func(t *testing.T) {
			testDataPath := path.Join("testdata", row.Name)

			raw, err := fs.ReadFile(testDataFS, testDataPath)
			if err != nil {
				t.Errorf("failed to read %q: %v", testDataPath, err)
				return
			}

			node, err := p.Parse(NewLexer(raw))
			if err != nil {
				t.Errorf("parse failed: %v", err)
			}

			if !row.Expect.Equals(node) {
				t.Errorf("parse gave wrong result:\n\texpect: %v\n\tactual: %v", row.Expect, node)
			}
		})
	}
}
