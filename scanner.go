package jo

import (
	"fmt"
)

// Valid scanner states.
const (
	sClean = iota
	sEOF

	sObjectKeyOrBrace
	sObjectColon
	sObjectCommaOrBrace
	sObjectKey

	sArrayElementOrBracket
	sArrayCommaOrBracket

	sString
	sStringEsc
	sStringUnicode
	sStringUnicode1
	sStringUnicode12
	sStringUnicode123

	sNumberNeg
	sNumberZero
	sNumberDigit
	sNumberDot
	sNumberDotDigit
	sNumberExp
	sNumberExpSign
	sNumberExpDigit

	sBoolT
	sBoolTr
	sBoolTru

	sBoolF
	sBoolFa
	sBoolFal
	sBoolFals

	sNullN
	sNullNu
	sNullNul
)

// Scanner is a JSON scanning state machine. The zero value for Scanner is
// ready to be used without any further initialization.
type Scanner struct {
	state int
	err   error
	stack []int
	size  int
	isKey bool
}

// Scan feeds another byte to the scanner.
//
// The first return value, refered to as an opcode, tells the caller about
// significant parsing events like beginning and ending values, so that the
// caller can follow along if it wishes.
//
// The second return value is used to indicate whether or not, which is
// necessary because of number literals (is 123 a whole value, or just the
// beginning of 12345?). Even though it will only ever be either 1 or 0, it is
// an integer rather than a boolean to allow for these kinds of constructs:
//
//     for i := 0; i < len(input); {
//         _, n := s.Scan(input[i])
//         i += n
//     }
func (s *Scanner) Scan(c byte) (Op, int) {
rewind:
	switch s.state {
	case sClean:
		switch c {
		case ' ', '\t', '\n', '\r':
			return OpSpace, 1
		case '{':
			s.state = sObjectKeyOrBrace
			return OpObjectStart, 1
		case '[':
			s.state = sArrayElementOrBracket
			return OpArrayStart, 1
		case '"':
			s.state = sString
			s.isKey = false
			return OpStringStart, 1
		case '-':
			s.state = sNumberNeg
			return OpNumberStart, 1
		case '0':
			s.state = sNumberZero
			return OpNumberStart, 1
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s.state = sNumberDigit
			return OpNumberStart, 1
		case 't':
			s.state = sBoolT
			return OpBoolStart, 1
		case 'f':
			s.state = sBoolF
			return OpBoolStart, 1
		case 'n':
			s.state = sNullN
			return OpNullStart, 1
		}
		return s.errorf(`expected start of JSON value, found %q`, c)

	case sObjectKeyOrBrace:
		switch c {
		case ' ', '\t', '\r', '\n':
			return OpSpace, 1
		case '"':
			s.state = sString
			s.isKey = true
			s.push(sObjectColon)
			return OpObjectKeyStart, 1
		case '}':
			s.state = s.pop()
			return OpObjectEnd, 1
		}
		return s.errorf(`expected object key or '}', found %q`, c)

	case sObjectColon:
		switch c {
		case ' ', '\t', 'r', '\n':
			return OpSpace, 1
		case ':':
			s.state = sClean
			s.push(sObjectCommaOrBrace)
			return OpContinue, 1
		}
		return s.errorf(`expected ':' after object key, found %q`, c)

	case sObjectCommaOrBrace:
		switch c {
		case ' ', '\t', '\r', '\n':
			return OpSpace, 1
		case ',':
			s.state = sObjectKey
			return OpContinue, 1
		case '}':
			s.state = s.pop()
			return OpObjectEnd, 1
		}
		return s.errorf(`expected object key or '}' after object value, found %q`, c)

	case sObjectKey:
		switch c {
		case ' ', '\t', '\r', '\n':
			return OpSpace, 1
		case '"':
			s.state = sString
			s.isKey = true
			s.push(sObjectColon)
			return OpObjectKeyStart, 1
		}
		return s.errorf(`expected object key after ',', found %q`, c)

	case sArrayElementOrBracket:
		switch c {
		case ' ', '\t', '\r', '\n':
			return OpSpace, 1
		case ']':
			s.state = s.pop()
			return OpArrayEnd, 1
		}
		s.push(sArrayCommaOrBracket)
		s.state = sClean
		goto rewind

	case sArrayCommaOrBracket:
		switch c {
		case ' ', '\t', '\n', '\r':
			return OpSpace, 1
		case ',':
			s.state = sArrayElementOrBracket
			return OpContinue, 1
		case ']':
			s.state = s.pop()
			return OpArrayEnd, 1
		}
		return s.errorf(`expected ',' or ']' after array element, found %q`, c)

	case sString:
		if c == '"' {
			s.state = s.pop()
			if s.isKey {
				return OpObjectKeyEnd, 1
			}
			return OpStringEnd, 1
		}
		if c == '\\' {
			s.state = sStringEsc
			return OpContinue, 1
		}
		if c >= 0x20 {
			return OpContinue, 1
		}
		return s.errorf(`unexpected control character in string literal`)

	case sStringEsc:
		switch c {
		case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
			s.state = sString
			return OpContinue, 1
		case 'u':
			s.state = sStringUnicode
			return OpContinue, 1
		}
		return s.errorf(`illegal escape sequence in string literal`)

	case sStringUnicode:
		if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
			s.state = sStringUnicode1
			return OpContinue, 1
		}
		return s.errorf(`illegal unicode escape sequence in string literal`)

	case sStringUnicode1:
		if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
			s.state = sStringUnicode12
			return OpContinue, 1
		}
		return s.errorf(`illegal unicode escape sequence in string literal`)

	case sStringUnicode12:
		if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
			s.state = sStringUnicode123
			return OpContinue, 1
		}
		return s.errorf(`illegal unicode escape sequence in string literal`)

	case sStringUnicode123:
		if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
			s.state = sString
			return OpContinue, 1
		}
		return s.errorf(`illegal unicode escape sequence in string literal`)

	case sNumberNeg:
		if c == '0' {
			s.state = sNumberZero
			return OpContinue, 1
		}
		if '1' <= c && c <= '9' {
			s.state = sNumberDigit
			return OpContinue, 1
		}
		return s.errorf(`expected digit after "-", found %q`, c)

	case sNumberDigit:
		if '0' <= c && c <= '9' {
			return OpContinue, 1
		}
		fallthrough

	case sNumberZero:
		if c == '.' {
			s.state = sNumberDot
			return OpContinue, 1
		}
		if c == 'e' || c == 'E' {
			s.state = sNumberExp
			return OpContinue, 1
		}
		s.state = s.pop()
		return OpNumberEnd, 0

	case sNumberDot:
		if '0' <= c && c <= '9' {
			s.state = sNumberDotDigit
			return OpContinue, 1
		}
		return s.errorf(`expected digit after decimal point, found %q`, c)

	case sNumberDotDigit:
		if '0' <= c && c <= '9' {
			return OpContinue, 1
		}
		if c == 'e' || c == 'E' {
			s.state = sNumberExp
			return OpContinue, 1
		}
		s.state = s.pop()
		return OpNumberEnd, 0

	case sNumberExp:
		if '0' <= c && c <= '9' {
			s.state = sNumberExpDigit
			return OpContinue, 1
		}
		if c == '-' || c == '+' {
			s.state = sNumberExpSign
			return OpContinue, 1
		}
		return s.errorf(`expected sign or digit in exponent, found %q`, c)

	case sNumberExpSign:
		if '0' <= c && c <= '9' {
			s.state = sNumberExpDigit
			return OpContinue, 1
		}
		return s.errorf(`expected digit after exponent sign, found %q`)

	case sNumberExpDigit:
		if '0' <= c && c <= '9' {
			return OpContinue, 1
		}
		s.state = s.pop()
		return OpNumberEnd, 0

	case sBoolT:
		if c == 'r' {
			s.state = sBoolTr
			return OpContinue, 1
		}
		return s.errorf(`expected 'r' after "t", found %q`, c)

	case sBoolTr:
		if c == 'u' {
			s.state = sBoolTru
			return OpContinue, 1
		}
		return s.errorf(`expected 'u' after "tr", found %q`, c)

	case sBoolTru:
		if c == 'e' {
			s.state = s.pop()
			return OpBoolEnd, 1
		}
		return s.errorf(`expected 'e' after "tru", found %q`, c)

	case sBoolF:
		if c == 'a' {
			s.state = sBoolFa
			return OpContinue, 1
		}
		return s.errorf(`expected 'a' after "f", found %q`, c)

	case sBoolFa:
		if c == 'l' {
			s.state = sBoolFal
			return OpContinue, 1
		}
		return s.errorf(`expected 'l' after "fa", found %q`, c)

	case sBoolFal:
		if c == 's' {
			s.state = sBoolFals
			return OpContinue, 1
		}
		return s.errorf(`expected 's' after "fal", found %q`, c)

	case sBoolFals:
		if c == 'e' {
			s.state = s.pop()
			return OpBoolEnd, 1
		}
		return s.errorf(`expected 'e' after "fals", found %q`, c)

	case sNullN:
		if c == 'u' {
			s.state = sNullNu
			return OpContinue, 1
		}
		return s.errorf(`expected 'u' after "n", found %q`, c)

	case sNullNu:
		if c == 'l' {
			s.state = sNullNul
			return OpContinue, 1
		}
		return s.errorf(`expected 'l' after "nu", found %q`, c)

	case sNullNul:
		if c == 'l' {
			s.state = s.pop()
			return OpNullEnd, 1
		}
		return s.errorf(`expected 'l' after "nul", found %q`, c)

	case sEOF:
		if c <= ' ' && (c == ' ' || c == '\t' || c == '\n' || c == '\r') {
			return OpSpace, 1
		}
		return s.errorf(`unexpected %q after top-level value`, c)
	}

	return OpContinue, 1
}

// Eof signals the scanner that the end of input has been reached. It returns
// an opcode, just as s.Scan does, which will be either OpEOF, OpNumberEnd or
// OpSyntaxError.
func (s *Scanner) Eof() Op {
	switch s.state {
	case sNumberZero, sNumberDigit, sNumberDotDigit, sNumberExpDigit:
		if len(s.stack) == 0 {
			s.state = sEOF
			return OpNumberEnd
		}
	case sEOF:
		return OpEOF
	}

	op, _ := s.errorf(`unexpected end of JSON input`)
	return op
}

// LastError returns the last error raised by the Scan or Eof methods.
func (s *Scanner) LastError() error {
	return s.err
}

// Reset resets the scanner's internal state.
func (s *Scanner) Reset() {
	s.state = 0
	s.err = nil
	s.stack = s.stack[:0]
	s.size = 0
}

// push adds a new state to the stack.
func (s *Scanner) push(x int) {
	if len(s.stack) == s.size {
		s.stack = append(s.stack, x)
	} else {
		s.stack[s.size] = x
	}
	s.size++
}

// pop takes the top state from the stack.
func (s *Scanner) pop() int {
	if s.size == 0 {
		return sEOF
	}

	s.size--
	return s.stack[s.size]
}

// errorf is a convenience function which both sets the scanner's error
// field and returns an OpSyntaxError opcode.
func (s *Scanner) errorf(format string, args ...interface{}) (Op, int) {
	s.err = fmt.Errorf(format, args...)
	return OpSyntaxError, 0
}
