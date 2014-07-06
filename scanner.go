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
	s.state = beforeValue
	s.stack = append(s.stack[:0], afterTopValue)
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
	return s.err
}

// errorf generates and persists an error.
func (s *Scanner) errorf(str string, args ...interface{}) Event {
	s.state = afterError
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
	s.state = delayed
	s.end = ev
	return None
}

func beforeValue(s *Scanner, c byte) Event {
	if c <= '9' {
		if c >= '1' {
			s.state = afterDigit
			return NumberStart
		} else if isSpace(c) {
			return Space
		} else if c == '"' {
			s.state = afterQuote
			s.end = StringEnd
			return StringStart
		} else if c == '-' {
			s.state = afterMinus
			return NumberStart
		} else if c == '0' {
			s.state = afterZero
			return NumberStart
		}
	} else if c == '{' {
		s.state = beforeFirstObjectKey
		return ObjectStart
	} else if c == '[' {
		s.state = beforeFirstArrayElement
		return ArrayStart
	} else if c == 't' {
		s.state = afterT
		return BoolStart
	} else if c == 'f' {
		s.state = afterF
		return BoolStart
	} else if c == 'n' {
		s.state = afterN
		return NullStart
	}

	return s.errorf("TODO")
}

func beforeFirstObjectKey(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == '"' {
		s.state = afterQuote
		s.end = KeyEnd
		s.push(afterObjectKey)
		return KeyStart
	} else if c == '}' {
		return s.delay(ObjectEnd)
	}

	return s.errorf("TODO")
}

func afterObjectKey(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ':' {
		s.state = beforeValue
		s.push(afterObjectValue)
		return None
	}

	return s.errorf("TODO")
}

func afterObjectValue(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ',' {
		s.state = afterObjectComma
		return None
	} else if c == '}' {
		return s.delay(ObjectEnd)
	}

	return s.errorf("TODO")
}

func afterObjectComma(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == '"' {
		s.state = afterQuote
		s.end = KeyEnd
		s.push(afterObjectKey)
		return KeyStart
	}

	return s.errorf("TODO")
}

func beforeFirstArrayElement(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ']' {
		return s.delay(ArrayEnd)
	}

	s.push(afterArrayElement)
	return beforeValue(s, c)
}

func afterArrayElement(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	} else if c == ',' {
		s.state = beforeValue
		s.push(afterArrayElement)
		return None
	} else if c == ']' {
		return s.delay(ArrayEnd)
	}

	return s.errorf("TODO")
}

func afterQuote(s *Scanner, c byte) Event {
	if c == '"' {
		// At thie point, s.end has already been set to either StringEnd or
		// KeyEnd depending on the previous state function.
		s.state = delayed
		return None
	} else if c == '\\' {
		s.state = afterEsc
		return None
	} else if c >= 0x20 {
		return None
	}

	return s.errorf("TODO")
}

func afterEsc(s *Scanner, c byte) Event {
	if isEsc(c) {
		s.state = afterQuote
		return None
	} else if c == 'u' {
		s.state = afterEscU
		return None
	}

	return s.errorf("TODO")
}

func afterEscU(s *Scanner, c byte) Event {
	if isHex(c) {
		s.state = afterEscU1
		return None
	}

	return s.errorf("TODO")
}

func afterEscU1(s *Scanner, c byte) Event {
	if isHex(c) {
		s.state = afterEscU12
		return None
	}

	return s.errorf("TODO")
}

func afterEscU12(s *Scanner, c byte) Event {
	if isHex(c) {
		s.state = afterEscU123
		return None
	}

	return s.errorf("TODO")
}

func afterEscU123(s *Scanner, c byte) Event {
	if isHex(c) {
		s.state = afterQuote
		return None
	}

	return s.errorf("TODO")
}

func afterMinus(s *Scanner, c byte) Event {
	if c == '0' {
		s.state = afterZero
		return None
	} else if '1' <= c && c <= '9' {
		s.state = afterDigit
		return None
	}

	return s.errorf("TODO")
}

func afterZero(s *Scanner, c byte) Event {
	if c == '.' {
		s.state = afterDot
		return None
	} else if c == 'e' || c == 'E' {
		s.state = afterE
		return None
	}

	return s.next(c) | NumberEnd
}

func afterDigit(s *Scanner, c byte) Event {
	if isDigit(c) {
		return None
	}

	return afterZero(s, c)
}

func afterDot(s *Scanner, c byte) Event {
	if isDigit(c) {
		s.state = afterDotDigit
		return None
	}

	return s.errorf("TODO")
}

func afterDotDigit(s *Scanner, c byte) Event {
	if isDigit(c) {
		return None
	} else if c == 'e' || c == 'E' {
		s.state = afterE
		return None
	}

	return s.next(c) | NumberEnd
}

func afterE(s *Scanner, c byte) Event {
	if isDigit(c) {
		s.state = afterEDigit
		return None
	} else if c == '-' || c == '+' {
		s.state = afterESign
		return None
	}

	return s.errorf("TODO")
}

func afterESign(s *Scanner, c byte) Event {
	if isDigit(c) {
		s.state = afterEDigit
		return None
	}

	return s.errorf("TODO")
}

func afterEDigit(s *Scanner, c byte) Event {
	if isDigit(c) {
		return None
	}

	return s.next(c) | NumberEnd
}

func afterT(s *Scanner, c byte) Event {
	if c == 'r' {
		s.state = afterTr
		return None
	}

	return s.errorf("TODO")
}

func afterTr(s *Scanner, c byte) Event {
	if c == 'u' {
		s.state = afterTru
		return None
	}

	return s.errorf("TODO")
}

func afterTru(s *Scanner, c byte) Event {
	if c == 'e' {
		return s.delay(BoolEnd)
	}

	return s.errorf("TODO")
}

func afterF(s *Scanner, c byte) Event {
	if c == 'a' {
		s.state = afterFa
		return None
	}

	return s.errorf("TODO")
}

func afterFa(s *Scanner, c byte) Event {
	if c == 'l' {
		s.state = afterFal
		return None
	}

	return s.errorf("TODO")
}

func afterFal(s *Scanner, c byte) Event {
	if c == 's' {
		s.state = afterFals
		return None
	}

	return s.errorf("TODO")
}

func afterFals(s *Scanner, c byte) Event {
	if c == 'e' {
		return s.delay(BoolEnd)
	}

	return s.errorf("TODO")
}

func afterN(s *Scanner, c byte) Event {
	if c == 'u' {
		s.state = afterNu
		return None
	}

	return s.errorf("TODO")
}

func afterNu(s *Scanner, c byte) Event {
	if c == 'l' {
		s.state = afterNul
		return None
	}

	return s.errorf("TODO")
}

func afterNul(s *Scanner, c byte) Event {
	if c == 'l' {
		return s.delay(NullEnd)
	}

	return s.errorf("TODO")
}

func delayed(s *Scanner, c byte) Event {
	return s.next(c) | s.end
}

func afterTopValue(s *Scanner, c byte) Event {
	if isSpace(c) {
		return Space
	}

	return s.errorf("TODO")
}

func afterError(s *Scanner, c byte) Event {
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
