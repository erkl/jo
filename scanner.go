package jo

import (
	"fmt"
)

// A Scanner is a state machine which emits a series of Events when fed JSON
// input.
type Scanner struct {
	// Current state.
	scan func(*Scanner, byte) Event

	// Scheduled state.
	stack []func(*Scanner, byte) Event

	// Used when delaying end events.
	end Event

	// Persisted syntax error.
	err error
}

// NewScanner initializes a new Scanner.
func NewScanner() *Scanner {
	s := &Scanner{stack: make([]func(*Scanner, byte) Event, 0, 4)}
	s.Reset()
	return s
}

// Reset restores a Scanner to its initial state.
func (s *Scanner) Reset() {
	s.scan = scanValue
	s.stack = append(s.stack[:0], scanEnd)
	s.err = nil
}

// Scan accepts a byte of input and returns an Event.
func (s *Scanner) Scan(c byte) Event {
	return s.scan(s, c)
}

// End signals the Scanner that the end of input has been reached. It returns
// an event just as Scan does.
func (s *Scanner) End() Event {
	// Feeding the scan function whitespace will for NumberEnd events.
	// Note the bitwise operation to filter out the Space bit.
	ev := s.scan(s, ' ') & (^Space)

	if s.err != nil {
		return Error
	}
	if len(s.stack) > 0 {
		return s.errorf("TODO")
	}

	return ev
}

// LastError returns a syntax error description after either Scan or End has
// returned an Error event.
func (s *Scanner) LastError() error {
	return nil
}

// errorf generates and persists an error.
func (s *Scanner) errorf(str string, args ...interface{}) Event {
	s.scan = scanError
	s.err = fmt.Errorf(str, args...)
	return Error
}

// Push another scan function onto the stack.
func (s *Scanner) push(fn func(*Scanner, byte) Event) {
	s.stack = append(s.stack, fn)
}

// Move the top scan function to s.scan.
func (s *Scanner) pop() {
	n := len(s.stack) - 1
	s.scan = s.stack[n]
	s.stack = s.stack[:n]
}

func scanValue(s *Scanner, c byte) Event {
	if c <= '9' {
		if c >= '1' {
			s.scan = scanDigit
			return NumberStart
		} else if isSpace(c) {
			return Space
		} else if c == '"' {
			s.scan = scanInString
			return StringStart
		} else if c == '-' {
			s.scan = scanNeg
			return NumberStart
		} else if c == '0' {
			s.scan = scanZero
			return NumberStart
		}
	} else if c == '{' {
		// TODO
	} else if c == '[' {
		s.scan = scanArray
		return ArrayStart
	} else if c == 't' {
		s.scan = scanT
		return BoolStart
	} else if c == 'f' {
		s.scan = scanF
		return BoolStart
	} else if c == 'n' {
		s.scan = scanN
		return NullStart
	}

	return s.errorf("TODO")
}

func scanArray(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ']' {
		s.scan = scanDelay
		s.end = ArrayEnd
		return None
	}

	s.push(scanElement)
	return scanValue(s, c)
}

func scanElement(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ',' {
		s.scan = scanValue
		s.push(scanElement)
		return None
	} else if c == ']' {
		s.scan = scanDelay
		s.end = ArrayEnd
		return None
	}

	return s.errorf("TODO")
}

func scanInString(s *Scanner, c byte) Event {
	if c == '"' {
		s.scan = scanDelay
		s.end = StringEnd
		return None
	} else if c == '\\' {
		s.scan = scanInStringEsc
		return None
	} else if c >= 0x20 {
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEsc(s *Scanner, c byte) Event {
	if isEsc(c) {
		s.scan = scanInString
		return None
	} else if c == 'u' {
		s.scan = scanInStringEscU
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEscU(s *Scanner, c byte) Event {
	if isHex(c) {
		s.scan = scanInStringEscU1
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEscU1(s *Scanner, c byte) Event {
	if isHex(c) {
		s.scan = scanInStringEscU12
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEscU12(s *Scanner, c byte) Event {
	if isHex(c) {
		s.scan = scanInStringEscU123
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEscU123(s *Scanner, c byte) Event {
	if isHex(c) {
		s.scan = scanInString
		return None
	}

	return s.errorf("TODO")
}

func scanNeg(s *Scanner, c byte) Event {
	if c == '0' {
		s.scan = scanZero
		return None
	} else if '1' <= c && c <= '9' {
		s.scan = scanDigit
		return None
	}

	return s.errorf("TODO")
}

func scanZero(s *Scanner, c byte) Event {
	if c == '.' {
		s.scan = scanDot
		return None
	} else if c == 'e' || c == 'E' {
		s.scan = scanE
		return None
	}

	s.pop()
	return s.scan(s, c) | NumberEnd
}

func scanDigit(s *Scanner, c byte) Event {
	if isDigit(c) {
		return None
	}

	return scanZero(s, c)
}

func scanDot(s *Scanner, c byte) Event {
	if isDigit(c) {
		s.scan = scanDotDigit
		return None
	}

	return s.errorf("TODO")
}

func scanDotDigit(s *Scanner, c byte) Event {
	if isDigit(c) {
		return None
	} else if c == 'e' || c == 'E' {
		s.scan = scanE
		return None
	}

	s.pop()
	return s.scan(s, c) | NumberEnd
}

func scanE(s *Scanner, c byte) Event {
	if isDigit(c) {
		s.scan = scanEDigit
		return None
	} else if c == '-' || c == '+' {
		s.scan = scanESign
		return None
	}

	return s.errorf("TODO")
}

func scanESign(s *Scanner, c byte) Event {
	if isDigit(c) {
		s.scan = scanEDigit
		return None
	}

	return s.errorf("TODO")
}

func scanEDigit(s *Scanner, c byte) Event {
	if isDigit(c) {
		return None
	}

	s.pop()
	return s.scan(s, c) | NumberEnd
}

func scanT(s *Scanner, c byte) Event {
	if c == 'r' {
		s.scan = scanTr
		return None
	}

	return s.errorf("TODO")
}

func scanTr(s *Scanner, c byte) Event {
	if c == 'u' {
		s.scan = scanTru
		return None
	}

	return s.errorf("TODO")
}

func scanTru(s *Scanner, c byte) Event {
	if c == 'e' {
		s.scan = scanDelay
		s.end = BoolEnd
		return None
	}

	return s.errorf("TODO")
}

func scanF(s *Scanner, c byte) Event {
	if c == 'a' {
		s.scan = scanFa
		return None
	}

	return s.errorf("TODO")
}

func scanFa(s *Scanner, c byte) Event {
	if c == 'l' {
		s.scan = scanFal
		return None
	}

	return s.errorf("TODO")
}

func scanFal(s *Scanner, c byte) Event {
	if c == 's' {
		s.scan = scanFals
		return None
	}

	return s.errorf("TODO")
}

func scanFals(s *Scanner, c byte) Event {
	if c == 'e' {
		s.scan = scanDelay
		s.end = BoolEnd
		return None
	}

	return s.errorf("TODO")
}

func scanN(s *Scanner, c byte) Event {
	if c == 'u' {
		s.scan = scanNu
		return None
	}

	return s.errorf("TODO")
}

func scanNu(s *Scanner, c byte) Event {
	if c == 'l' {
		s.scan = scanNul
		return None
	}

	return s.errorf("TODO")
}

func scanNul(s *Scanner, c byte) Event {
	if c == 'l' {
		s.scan = scanDelay
		s.end = NullEnd
		return None
	}

	return s.errorf("TODO")
}

func scanDelay(s *Scanner, c byte) Event {
	s.pop()
	return s.scan(s, c) | s.end
}

func scanEnd(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	}

	return s.errorf("TODO")
}

func scanError(s *Scanner, c byte) Event {
	return Error
}

// isSpace returns true if c is a whitespace character.
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// isDigit returns true if c is a valid decimal digit.
func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// isHex returns true if c is a valid hexadecimal digit.
func isHex(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

// isEsc returns true if `\` + c is a valid escape sequence.
func isEsc(c byte) bool {
	return c == 'b' || c == 'f' || c == 'n' || c == 'r' || c == 't' ||
		c == '\\' || c == '/' || c == '"'
}
