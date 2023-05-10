package wat

import (
	"io/fs"
	"path"
	"reflect"
	"testing"
)

func L(v ...string) []string {
	return v
}

func S(t SpaceType, n uint) Space {
	return Space{Type: t, Count: n}
}

func TestLexer(t *testing.T) {
	type TestCase struct {
		Name   string
		Expect []Token
	}

	testCases := [...]TestCase{
		{
			Name: "file1.wat",
			Expect: []Token{
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "module"},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: AcceptToken},
			},
		},
		{
			Name: "file2.wat",
			Expect: []Token{
				Token{Type: LineCommentToken, Value: ` one line comment`},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: LineCommentToken, Value: ` another one line comment`},
				Token{Type: SpaceToken, Value: S(LF, 2)},
				Token{Type: BlockCommentToken, Value: L(``, ` multiline`, ` block`, ` comment`, ` `)},
				Token{Type: SpaceToken, Value: S(LF, 2)},
				Token{Type: BlockCommentToken, Value: L(` block comment with embedded block comment (; whee i can nest ;) trailing text `)},
				Token{Type: SpaceToken, Value: S(LF, 2)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "module"},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: AcceptToken},
			},
		},
		{
			Name: "file3.wat",
			Expect: []Token{
				Token{Type: LineCommentToken, Value: ` Copied from https://github.com/bytecodealliance/wasmtime/blob/main/docs/WASI-tutorial.md`},
				Token{Type: SpaceToken, Value: S(LF, 2)},

				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "module"},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: LineCommentToken, Value: ` Import the required fd_write WASI function which will write the given io vectors to stdout`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: LineCommentToken, Value: ` The function signature for fd_write is:`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: LineCommentToken, Value: ` (File Descriptor, *iovs, iovs_len, nwritten) -> Returns number of bytes written`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "import"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: StringToken, Value: "wasi_unstable"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: StringToken, Value: "fd_write"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "func"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: IdentifierToken, Value: "$fd_write"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "param"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: KeywordToken, Value: "i32"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: KeywordToken, Value: "i32"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: KeywordToken, Value: "i32"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: KeywordToken, Value: "i32"},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "result"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: KeywordToken, Value: "i32"},
				Token{Type: CloseParenToken},
				Token{Type: CloseParenToken},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 2)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "memory"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "1"}},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "export"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: StringToken, Value: "memory"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "memory"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "0"}},
				Token{Type: CloseParenToken},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 2)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: LineCommentToken, Value: ` Write 'hello world\n' to memory at an offset of 8 bytes`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: LineCommentToken, Value: ` Note the trailing newline which is required for the text to appear`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "data"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.const"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "8"}},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: StringToken, Value: "hello world\n"},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 2)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "func"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: IdentifierToken, Value: "$main"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "export"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: StringToken, Value: "_start"},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 8)},
				Token{Type: LineCommentToken, Value: ` Creating a new io vector within linear memory`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 8)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.store"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.const"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "0"}},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.const"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "8"}},
				Token{Type: CloseParenToken},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 2)},
				Token{Type: LineCommentToken, Value: ` iov.iov_base - This is a pointer to the start of the 'hello world\n' string`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 8)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.store"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.const"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "4"}},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.const"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "12"}},
				Token{Type: CloseParenToken},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 2)},
				Token{Type: LineCommentToken, Value: ` iov.iov_len - The length of the 'hello world\n' string`},
				Token{Type: SpaceToken, Value: S(LF, 2)},

				Token{Type: SpaceToken, Value: S(SP, 8)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "call"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: IdentifierToken, Value: "$fd_write"},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 12)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.const"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "1"}},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: LineCommentToken, Value: ` file_descriptor - 1 for stdout`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 12)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.const"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "0"}},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: LineCommentToken, Value: ` *iovs - The pointer to the iov array, which is stored at memory location 0`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 12)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.const"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "1"}},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: LineCommentToken, Value: ` iovs_len - We're printing 1 string stored in an iov - so one.`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 12)},
				Token{Type: OpenParenToken},
				Token{Type: KeywordToken, Value: "i32.const"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "20"}},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: LineCommentToken, Value: ` nwritten - A place in memory to store the number of bytes written`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 8)},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 8)},
				Token{Type: KeywordToken, Value: "drop"},
				Token{Type: SpaceToken, Value: S(SP, 1)},
				Token{Type: LineCommentToken, Value: ` Discard the number of bytes written from the top of the stack`},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: SpaceToken, Value: S(SP, 4)},
				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 1)},

				Token{Type: CloseParenToken},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: AcceptToken},
			},
		},
		{
			Name: "strings.wat",
			Expect: []Token{
				Token{Type: StringToken, Value: "this is a string"},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: StringToken, Value: "tab: \t"},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: StringToken, Value: "newline: \n"},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: StringToken, Value: "crlf: \r\n"},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: StringToken, Value: `backslash: \`},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: StringToken, Value: `single quote: '`},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: StringToken, Value: `double quote: "`},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: StringToken, Value: "ESC[0m: \x1b[0m"},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: StringToken, Value: "smiley: \u263a\ufe0f"},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: StringToken, Value: "smiley: \u263a\ufe0f"},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: AcceptToken},
			},
		},
		{
			Name: "numbers.wat",
			Expect: []Token{
				Token{Type: NumberToken, Value: Num{Integer: "0"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "123"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagSign, Integer: "0"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagSign, Integer: "123"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagSign | FlagNeg, Integer: "0"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagSign | FlagNeg, Integer: "123"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat, Integer: "0", Fraction: "0"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat, Integer: "0", Fraction: "123"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagSign, Integer: "0", Fraction: "0"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagSign, Integer: "0", Fraction: "123"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagSign | FlagNeg, Integer: "0", Fraction: "0"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagSign | FlagNeg, Integer: "0", Fraction: "123"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagHex, Integer: "0"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagHex, Integer: "a"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagHex | FlagSign, Integer: "0"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagHex | FlagSign, Integer: "a"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagHex | FlagSign | FlagNeg, Integer: "0"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagHex | FlagSign | FlagNeg, Integer: "a"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagHex, Integer: "a", Fraction: "b"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Integer: "123"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat, Integer: "123", Fraction: "456"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagHex, Integer: "123", Fraction: "456"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagInf}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagInf | FlagSign}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagInf | FlagSign | FlagNeg}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagNaN}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagNaN | FlagSign}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagNaN | FlagSign | FlagNeg}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagNaN | FlagAcanonical | FlagHex, Integer: "deadbeef"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagNaN | FlagAcanonical | FlagHex | FlagSign, Integer: "deadbeef"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagNaN | FlagAcanonical | FlagHex | FlagSign | FlagNeg, Integer: "deadbeef"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat, Integer: "0", Fraction: "5", Exponent: "20"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagExpSign, Integer: "0", Fraction: "5", Exponent: "20"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagExpSign | FlagExpNeg, Integer: "0", Fraction: "5", Exponent: "20"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagHex, Integer: "0", Fraction: "8", Exponent: "12"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagHex | FlagExpSign, Integer: "0", Fraction: "8", Exponent: "12"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: NumberToken, Value: Num{Flags: FlagFloat | FlagHex | FlagExpSign | FlagExpNeg, Integer: "0", Fraction: "8", Exponent: "12"}},
				Token{Type: SpaceToken, Value: S(LF, 1)},
				Token{Type: AcceptToken},
			},
		},
	}

	for _, row := range testCases {
		t.Run(row.Name, func(t *testing.T) {
			testDataPath := path.Join("testdata", row.Name)

			raw, err := fs.ReadFile(testDataFS, testDataPath)
			if err != nil {
				t.Errorf("failed to read %q: %v", testDataPath, err)
				return
			}

			expect := row.Expect
			expectLen := uint(len(expect))
			actual := make([]Token, 0, expectLen)

			lexer := NewLexer(raw)
			for lexer.HasNext() {
				token := lexer.Next()
				actual = append(actual, token)
			}
			actualLen := uint(len(actual))

			for i := uint(0); i < expectLen && i < actualLen; i++ {
				a := expect[i]
				b := actual[i]
				if a.Type != b.Type || !reflect.DeepEqual(a.Value, b.Value) {
					t.Errorf("token #%d: mismatch\n\texpect: %v\n\tactual: %v", i, a, b)
				}
			}
			if expectLen > actualLen {
				t.Errorf("token stream ends %d elements earlier than expected", expectLen-actualLen)
			}
			if actualLen > expectLen {
				t.Errorf("token stream ends %d elements later than expected", actualLen-expectLen)
				for i := expectLen; i < actualLen; i++ {
					x := actual[i]
					t.Logf("token #%d: %v", i, x)
				}
			}
		})
	}
}
