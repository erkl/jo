package jo

import (
	"fmt"
)

// A Scanner is a state machine which emits a series of Events when fed JSON
// input.
type Scanner struct {
	// Current state.
	state func(*Scanner, byte) Event

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
	s.state = scanValue
	s.stack = append(s.stack[:0], scanEnd)
	s.err = nil
}

// Scan accepts a byte of input and returns an Event.
func (s *Scanner) Scan(c byte) Event {
	return s.state(s, c)
}

// End signals the Scanner that the end of input has been reached. It returns
// an event just as Scan does.
func (s *Scanner) End() Event {
	// Feeding the state function whitespace will for NumberEnd events.
	// Note the bitwise operation to filter out the Space bit.
	ev := s.state(s, ' ') & (^Space)

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
	s.state = scanError
	s.err = fmt.Errorf(str, args...)
	return Error
}

// Push another state function onto the stack.
func (s *Scanner) push(fn func(*Scanner, byte) Event) {
	s.stack = append(s.stack, fn)
}

// next pops the next state function off the stack and invokes it.
func (s *Scanner) next(c byte) Event {
	n := len(s.stack) - 1
	s.state = s.stack[n]
	s.stack = s.stack[:n]

	return s.state(s, c)
}

// delay schedules an end event to be returned for the next byte of input.
func (s *Scanner) delay(ev Event) Event {
	s.state = scanDelay
	s.end = ev
	return None
}

func scanValue(s *Scanner, c byte) Event {
	if c <= '9' {
		if c >= '1' {
			s.state = scanDigit
			return NumberStart
		} else if isSpace(c) {
			return Space
		} else if c == '"' {
			s.state = scanInString
			s.end = StringEnd
			return StringStart
		} else if c == '-' {
			s.state = scanNeg
			return NumberStart
		} else if c == '0' {
			s.state = scanZero
			return NumberStart
		}
	} else if c == '{' {
		s.state = scanObject
		return ObjectStart
	} else if c == '[' {
		s.state = scanArray
		return ArrayStart
	} else if c == 't' {
		s.state = scanT
		return BoolStart
	} else if c == 'f' {
		s.state = scanF
		return BoolStart
	} else if c == 'n' {
		s.state = scanN
		return NullStart
	}

	return s.errorf("TODO")
}

func scanObject(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == '"' {
		s.state = scanInString
		s.end = KeyEnd
		s.push(scanKey)
		return KeyStart
	} else if c == '}' {
		return s.delay(ObjectEnd)
	}

	return s.errorf("TODO")
}

func scanKey(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ':' {
		s.state = scanValue
		s.push(scanProperty)
		return None
	}

	return s.errorf("TODO")
}

func scanProperty(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ',' {
		s.state = scanInString
		s.end = KeyEnd
		s.push(scanKey)
		return KeyStart
	} else if c == '}' {
		return s.delay(ObjectEnd)
	}

	return s.errorf("TODO")
}

func scanArray(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ']' {
		return s.delay(ArrayEnd)
	}

	s.push(scanElement)
	return scanValue(s, c)
}

func scanElement(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ',' {
		s.state = scanValue
		s.push(scanElement)
		return None
	} else if c == ']' {
		return s.delay(ArrayEnd)
	}

	return s.errorf("TODO")
}

func scanInString(s *Scanner, c byte) Event {
	if c == '"' {
		s.state = scanDelay
		return None
	} else if c == '\\' {
		s.state = scanInStringEsc
		return None
	} else if c >= 0x20 {
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEsc(s *Scanner, c byte) Event {
	if isEsc(c) {
		s.state = scanInString
		return None
	} else if c == 'u' {
		s.state = scanInStringEscU
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEscU(s *Scanner, c byte) Event {
	if isHex(c) {
		s.state = scanInStringEscU1
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEscU1(s *Scanner, c byte) Event {
	if isHex(c) {
		s.state = scanInStringEscU12
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEscU12(s *Scanner, c byte) Event {
	if isHex(c) {
		s.state = scanInStringEscU123
		return None
	}

	return s.errorf("TODO")
}

func scanInStringEscU123(s *Scanner, c byte) Event {
	if isHex(c) {
		s.state = scanInString
		return None
	}

	return s.errorf("TODO")
}

func scanNeg(s *Scanner, c byte) Event {
	if c == '0' {
		s.state = scanZero
		return None
	} else if '1' <= c && c <= '9' {
		s.state = scanDigit
		return None
	}

	return s.errorf("TODO")
}

func scanZero(s *Scanner, c byte) Event {
	if c == '.' {
		s.state = scanDot
		return None
	} else if c == 'e' || c == 'E' {
		s.state = scanE
		return None
	}

	return s.next(c) | NumberEnd
}

func scanDigit(s *Scanner, c byte) Event {
	if isDigit(c) {
		return None
	}

	return scanZero(s, c)
}

func scanDot(s *Scanner, c byte) Event {
	if isDigit(c) {
		s.state = scanDotDigit
		return None
	}

	return s.errorf("TODO")
}

func scanDotDigit(s *Scanner, c byte) Event {
	if isDigit(c) {
		return None
	} else if c == 'e' || c == 'E' {
		s.state = scanE
		return None
	}

	return s.next(c) | NumberEnd
}

func scanE(s *Scanner, c byte) Event {
	if isDigit(c) {
		s.state = scanEDigit
		return None
	} else if c == '-' || c == '+' {
		s.state = scanESign
		return None
	}

	return s.errorf("TODO")
}

func scanESign(s *Scanner, c byte) Event {
	if isDigit(c) {
		s.state = scanEDigit
		return None
	}

	return s.errorf("TODO")
}

func scanEDigit(s *Scanner, c byte) Event {
	if isDigit(c) {
		return None
	}

	return s.next(c) | NumberEnd
}

func scanT(s *Scanner, c byte) Event {
	if c == 'r' {
		s.state = scanTr
		return None
	}

	return s.errorf("TODO")
}

func scanTr(s *Scanner, c byte) Event {
	if c == 'u' {
		s.state = scanTru
		return None
	}

	return s.errorf("TODO")
}

func scanTru(s *Scanner, c byte) Event {
	if c == 'e' {
		return s.delay(BoolEnd)
	}

	return s.errorf("TODO")
}

func scanF(s *Scanner, c byte) Event {
	if c == 'a' {
		s.state = scanFa
		return None
	}

	return s.errorf("TODO")
}

func scanFa(s *Scanner, c byte) Event {
	if c == 'l' {
		s.state = scanFal
		return None
	}

	return s.errorf("TODO")
}

func scanFal(s *Scanner, c byte) Event {
	if c == 's' {
		s.state = scanFals
		return None
	}

	return s.errorf("TODO")
}

func scanFals(s *Scanner, c byte) Event {
	if c == 'e' {
		return s.delay(BoolEnd)
	}

	return s.errorf("TODO")
}

func scanN(s *Scanner, c byte) Event {
	if c == 'u' {
		s.state = scanNu
		return None
	}

	return s.errorf("TODO")
}

func scanNu(s *Scanner, c byte) Event {
	if c == 'l' {
		s.state = scanNul
		return None
	}

	return s.errorf("TODO")
}

func scanNul(s *Scanner, c byte) Event {
	if c == 'l' {
		return s.delay(NullEnd)
	}

	return s.errorf("TODO")
}

func scanDelay(s *Scanner, c byte) Event {
	return s.next(c) | s.end
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
