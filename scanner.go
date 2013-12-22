package jo

import (
	"fmt"
)

// Valid scanner states.
const (
	sClean = iota
	sEOF

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
	switch s.state {
	case sClean:
		switch c {
		case ' ', '\t', '\n', '\r':
			return OpSpace, 1
		case '{':
			// @todo
		case '[':
			// @todo
		case '"':
			// @todo
		case '-':
			// @todo
		case '0':
			// @todo
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// @todo
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
	return OpEOF
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
