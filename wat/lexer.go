package wat

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenStream interface {
	HasNext() bool
	Next() Token
}

type Lexer struct {
	input []byte
	pos   Position
	next  Token

	scratchBytes   [1024]byte
	scratchStrings [64]string
}

func NewLexer(input []byte) *Lexer {
	return &Lexer{input: input}
}

func (lexer *Lexer) HasNext() bool {
	if lexer.next.IsTerminal() {
		return false
	}

	lexer.next = Token{}
	lexer.next.Type = InvalidToken
	lexer.next.Span = Span{Begin: lexer.pos, End: lexer.pos}

	if len(lexer.input) <= 0 {
		lexer.next.Type = AcceptToken
		return true
	}

	lexer.lex()
	return true
}

func (lexer *Lexer) Next() Token {
	return lexer.next
}

func (lexer *Lexer) lex() {
	ch, ok := lexer.peekRune()
	switch {
	case !ok:
		return
	case ch < 0:
		lexer.accept()
	case ch == '\t':
		lexer.lexSpace(HT)
	case ch == '\n':
		lexer.lexSpace(LF)
	case ch == '\r':
		lexer.lexSpace(CR)
	case ch == ' ':
		lexer.lexSpace(SP)
	case ch >= '0' && ch <= '9':
		lexer.lexNumber(0)
	case ch >= 'A' && ch <= 'Z':
		lexer.lexKeyword()
	case ch >= 'a' && ch <= 'z':
		lexer.lexKeyword()
	default:
		lexer.lexOther()
	}
}

func (lexer *Lexer) lexSpace(t SpaceType) {
	n := uint(0)
	r := t.Rune()
	m := lexer.createMark()
	ch, ok := lexer.readRune()
	for ok && ch == r {
		n++
		m = lexer.createMark()
		ch, ok = lexer.readRune()
	}
	if t == CR && n == 1 && ok && ch == '\n' {
		t = CRLF
		for {
			m = lexer.createMark()
			ch, ok = lexer.readRune()
			if !ok || ch != '\r' {
				break
			}

			m = lexer.createMark()
			ch, ok = lexer.readRune()
			if !ok || ch != '\n' {
				break
			}

			n++
		}
	}
	m.rewind(lexer)
	lexer.done(SpaceToken, Space{Type: t, Count: n})
}

func (lexer *Lexer) lexNumber(flags NumFlags) {
	var num Num
	num.Flags = flags

	m := lexer.createMark()
	keyword := lexer.takeWhile(nil, isIdentifier)
	if keywordToNumber(&num, keyword) {
		lexer.done(NumberToken, num)
		return
	}
	m.rewind(lexer)

	hex := lexer.match('0', 'x')
	exponentRune := rune('e')
	if hex {
		num.Flags |= FlagHex
		exponentRune = 'p'
	}

	num.Integer = lexer.takeDigits(nil, hex)
	if num.Integer == "" {
		num.Integer = "0"
	}

	if lexer.match('.') {
		num.Flags |= FlagFloat
		num.Fraction = lexer.takeDigits(nil, hex)
		if num.Fraction == "" {
			num.Fraction = "0"
		}
	}

	if lexer.match(exponentRune) {
		num.Flags |= FlagFloat | lexer.takeSign(FlagExpSign, FlagExpNeg)
		num.Exponent = lexer.takeDigits(nil, false)
		if num.Exponent == "" {
			ch, ok := lexer.peekRune()
			if ok {
				lexer.rejectUnexpected(ch, `at least one exponent digit`)
			}
			return
		}
	}

	lexer.done(NumberToken, num)
}

func (lexer *Lexer) lexKeyword() {
	keyword := lexer.takeWhile(nil, isIdentifier)
	var num Num
	if keywordToNumber(&num, keyword) {
		lexer.done(NumberToken, num)
		return
	}
	lexer.done(KeywordToken, keyword)
	return
}

func (lexer *Lexer) lexOther() {
	ch, ok := lexer.readRune()
	if !ok {
		return
	}

	switch ch {
	case ';':
		ch, ok = lexer.readRune()
		if !ok {
			return
		}
		if ch != ';' {
			lexer.rejectUnexpected(ch, `start of single-line comment ';;'`)
			return
		}
		lexer.done(LineCommentToken, lexer.takeWhile(nil, isLineComment))

	case '(':
		if lexer.match(';') {
			lexer.lexBlockComment()
			return
		}
		lexer.done(OpenParenToken, nil)

	case ')':
		lexer.done(CloseParenToken, nil)

	case '$':
		partial := lexer.scratchBytes[:0]
		partial = append(partial, '$')
		lexer.done(IdentifierToken, lexer.takeWhile(partial, isIdentifier))

	case '+':
		lexer.lexNumber(FlagSign)

	case '-':
		lexer.lexNumber(FlagSign | FlagNeg)

	case '"':
		lexer.lexString()

	default:
		lexer.rejectUnexpected(ch, `start of token`)
	}
}

func (lexer *Lexer) lexBlockComment() {
	const (
		stateReady = iota
		stateOpenParen
		stateSemicolon
		stateCR
	)

	lines := lexer.scratchStrings[:0]
	partial := lexer.scratchBytes[:0]
	counter := 1
	state := stateReady

	flush := func() {
		lines = append(lines, string(partial))
		partial = lexer.scratchBytes[:0]
	}

	for {
		ch, ok := lexer.readRune()
		if !ok {
			return
		}

		if ch < 0 {
			lexer.rejectUnexpected(ch, `block comment terminator ';)'`)
			return
		}

		if state == stateSemicolon && ch == ')' {
			counter--
			if counter <= 0 {
				flush()
				out := make([]string, len(lines))
				copy(out, lines)
				lexer.done(BlockCommentToken, out)
				return
			}
		}

		if state == stateSemicolon {
			partial = append(partial, ';')
			state = stateReady
		}

		switch ch {
		case '\r':
			flush()
			state = stateCR

		case '\n':
			if state != stateCR {
				flush()
			}
			state = stateReady

		case '(':
			partial = append(partial, '(')
			state = stateOpenParen

		case ';':
			switch state {
			case stateOpenParen:
				partial = append(partial, ';')
				counter++
				state = stateReady

			default:
				state = stateSemicolon
			}

		default:
			partial = utf8.AppendRune(partial, ch)
			state = stateReady
		}
	}
}

func (lexer *Lexer) lexString() {
	partial := lexer.scratchBytes[:0]
	for {
		ch, ok := lexer.readRune()
		if !ok {
			return
		}
		if ch < 0 {
			lexer.rejectUnexpected(ch, `string terminator '"'`)
			return
		}

		if ch == '"' {
			lexer.done(StringToken, string(partial))
			return
		}

		if ch == '\\' {
			var ok bool
			partial, ok = lexer.appendEscape(partial)
			if !ok {
				return
			}
			continue
		}

		partial = utf8.AppendRune(partial, ch)
	}
}

func (lexer *Lexer) appendEscape(partial []byte) ([]byte, bool) {
	ch, ok := lexer.readRune()
	if !ok {
		return nil, false
	}

	if isSimpleEscape(ch) {
		partial = append(partial, simpleEscapeMap[ch])
		return partial, true
	}

	if isHexDigit(ch) {
		var tmp [2]byte
		tmp[0] = byte(ch)

		ch, ok = lexer.readRune()
		if !ok {
			return nil, false
		}
		if !isHexDigit(ch) {
			lexer.rejectUnexpected(ch, "hex digit")
			return nil, false
		}
		tmp[1] = byte(ch)
		hex := string(tmp[:])

		u64, err := strconv.ParseUint(hex, 16, 8)
		if err != nil {
			panic(err)
		}

		ch = rune(u64)
		partial = utf8.AppendRune(partial, ch)
		return partial, true
	}

	if ch == 'u' {
		ch, ok = lexer.readRune()
		if !ok {
			return nil, false
		}
		if ch != '{' {
			lexer.rejectUnexpected(ch, `'{' as next character in Unicode escape '\u{H+}'`)
			return nil, false
		}

		var tmp [8]byte
		hex := lexer.takeWhile(tmp[:0], isHexDigit)
		if hex == "" {
			ch, ok = lexer.peekRune()
			if ok {
				lexer.rejectUnexpected(ch, `at least 1 hex digit`)
			}
			return nil, false
		}

		ch, ok = lexer.readRune()
		if ch != '}' {
			if ok {
				lexer.rejectUnexpected(ch, `'}' as next character in Unicode escape '\u{H+}'`)
			}
			return nil, false
		}

		u64, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			lexer.rejectf("hex escape \\u{%s} is out of range for a 32-bit unsigned integer", hex)
			return nil, false
		}

		ch = rune(u64)
		if ch < 0 || !utf8.ValidRune(ch) {
			lexer.rejectf("hex escape \\u{%s} encodes an invalid Unicode rune, U+%04x", hex, u64)
			return nil, false
		}

		partial = utf8.AppendRune(partial, ch)
		return partial, true
	}

	lexer.rejectUnexpected(ch, `escape sequence (one of '\\', '\"', '\'', '\t', '\n', '\r', '\HH', or '\u{H+}')`)
	return nil, false
}

func (lexer *Lexer) takeDigits(partial []byte, hex bool) string {
	if partial == nil {
		partial = lexer.scratchBytes[:0]
	}

	isDigit := isDecDigit
	if hex {
		isDigit = isHexDigit
	}

	m := lexer.createMark()
	ch, ok := lexer.readRune()
	if !ok || !isDigit(ch) {
		m.rewind(lexer)
		return ""
	}

	partial = utf8.AppendRune(partial, ch)
	for {
		m = lexer.createMark()
		ch, ok = lexer.readRune()
		if !ok {
			break
		}
		if ch == '_' {
			continue
		}
		if !isDigit(ch) {
			break
		}
		partial = utf8.AppendRune(partial, ch)
	}
	m.rewind(lexer)
	return string(partial)
}

func (lexer *Lexer) takeSign(signed, neg NumFlags) NumFlags {
	m := lexer.createMark()
	ch, _ := lexer.readRune()
	switch ch {
	case '+':
		return signed
	case '-':
		return signed | neg
	default:
		m.rewind(lexer)
		return 0
	}
}

func (lexer *Lexer) takeWhile(partial []byte, pred func(rune) bool) string {
	if partial == nil {
		partial = lexer.scratchBytes[:0]
	}

	m := lexer.createMark()
	ch, ok := lexer.readRune()
	for ok && pred(ch) {
		partial = utf8.AppendRune(partial, ch)
		m = lexer.createMark()
		ch, ok = lexer.readRune()
	}
	m.rewind(lexer)
	return string(partial)
}

func (lexer *Lexer) match(v ...rune) bool {
	m := lexer.createMark()
	for len(v) > 0 {
		ch, ok := lexer.readRune()
		if !ok || ch != v[0] {
			m.rewind(lexer)
			return false
		}
		v = v[1:]
	}
	return true
}

func (lexer *Lexer) peekRune() (ch rune, ok bool) {
	m := lexer.createMark()
	ch, ok = lexer.readRune()
	if ok {
		m.rewind(lexer)
	}
	return
}

func (lexer *Lexer) readRune() (rune, bool) {
	if len(lexer.input) <= 0 {
		return -1, true
	}

	ch, size := utf8.DecodeRune(lexer.input)
	if size < 1 || (size == 1 && ch == utf8.RuneError) {
		tmp := hexBytes(lexer.input)
		suffix := ""
		if len(tmp) > 8 {
			tmp = tmp[:8]
			suffix = "..."
		}
		lexer.rejectf("UTF-8 decode error: %v%s", tmp, suffix)
		return -1, false
	}

	if ch < 0 || !utf8.ValidRune(ch) {
		lexer.rejectf("invalid Unicode character %q U+%04x", ch, ch)
		return -1, false
	}

	if unicode.IsControl(ch) && !isSpace(ch) {
		lexer.rejectf("unexpected Unicode control character %q U+%04x", ch, ch)
		return -1, false
	}

	lexer.input = lexer.input[size:]
	lexer.pos.Advance(ch, size)
	return ch, true
}

func (lexer *Lexer) rejectUnexpected(ch rune, expect string) {
	partial := lexer.scratchBytes[:0]
	partial = append(partial, "unexpected "...)
	str := "end of input"
	if ch >= 0 {
		str = fmt.Sprintf("character %q U+%04x", ch, ch)
	}
	partial = append(partial, str...)
	if expect != "" {
		partial = append(partial, ": expect "...)
		partial = append(partial, expect...)
	}
	lexer.reject(errors.New(string(partial)))
}

func (lexer *Lexer) rejectf(format string, v ...any) {
	lexer.reject(fmt.Errorf(format, v...))
}

func (lexer *Lexer) reject(err error) {
	lexer.done(RejectToken, err)
}

func (lexer *Lexer) accept() {
	lexer.done(AcceptToken, nil)
}

func (lexer *Lexer) done(tt TokenType, tv any) {
	lexer.next.Type = tt
	lexer.next.Value = tv
	lexer.next.Span.End = lexer.pos
	if err := lexer.next.Validate(); err != nil {
		panic(err)
	}
}

var _ TokenStream = (*Lexer)(nil)

type mark struct {
	input []byte
	pos   Position
	next  Token
}

func (lexer *Lexer) createMark() mark {
	return mark{lexer.input, lexer.pos, lexer.next}
}

func (m mark) rewind(lexer *Lexer) {
	lexer.input = m.input
	lexer.pos = m.pos
	lexer.next = m.next
}

func keywordToNumber(num *Num, keyword string) bool {
	if keyword == "inf" {
		num.Flags |= FlagFloat | FlagInf
		return true
	}
	if keyword == "nan" {
		num.Flags |= FlagFloat | FlagNaN
		return true
	}
	if strings.HasPrefix(keyword, "nan:0x") {
		num.Flags |= FlagFloat | FlagNaN | FlagHex | FlagAcanonical
		num.Integer = keyword[6:]
		return true
	}
	return false
}

func isLineComment(ch rune) bool {
	return ch != '\r' && ch != '\n'
}

func isRuneInTable(ch rune, tab [4]uint32) bool {
	var ok bool
	if ch >= 0x00 && ch < 0x80 {
		byteValue := byte(ch)
		i := (byteValue >> 5)
		j := (byteValue & 0x1f)
		row := tab[i]
		bit := uint32(1) << j
		ok = (row & bit) != 0
	}
	return ok
}

var spaceTable = [4]uint32{
	0x00002600, // 00 to 1f
	0x00000001, // 20 to 3f
	0x00000000, // 40 to 5f
	0x00000000, // 60 to 7f
}

func isSpace(ch rune) bool {
	return isRuneInTable(ch, spaceTable)
}

var decDigitTable = [4]uint32{
	0x00000000, // 00 to 1f
	0x03ff0000, // 20 to 3f
	0x00000000, // 40 to 5f
	0x00000000, // 60 to 7f
}

func isDecDigit(ch rune) bool {
	return isRuneInTable(ch, decDigitTable)
}

var hexDigitTable = [4]uint32{
	0x00000000, // 00 to 1f
	0x03ff0000, // 20 to 3f
	0x0000007e, // 40 to 5f
	0x0000007e, // 60 to 7f
}

func isHexDigit(ch rune) bool {
	return isRuneInTable(ch, hexDigitTable)
}

var identifierTable = [4]uint32{
	0x00000000, // 00 to 1f
	0xf7ffecea, // 20 to 3f
	0xd7ffffff, // 40 to 5f
	0x57ffffff, // 60 to 7f
}

func isIdentifier(ch rune) bool {
	return isRuneInTable(ch, identifierTable)
}

var simpleEscapeTable = [4]uint32{
	0x00000000, // 00 to 1f
	0x00000084, // 20 to 3f
	0x10000000, // 40 to 5f
	0x00144000, // 60 to 7f
}

var simpleEscapeMap = map[rune]byte{
	't':  0x09,
	'n':  0x0a,
	'r':  0x0d,
	'"':  0x22,
	'\'': 0x27,
	'\\': 0x5c,
}

func isSimpleEscape(ch rune) bool {
	return isRuneInTable(ch, simpleEscapeTable)
}

func trace(pos Position, ch rune) {
	fmt.Fprintf(os.Stderr, "trace: %v, %q U+%04x\n", pos, ch, ch)
}
