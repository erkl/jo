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
		return s.errorf(`unexpected end of JSON input`)
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

// push pushes a state function onto the stack.
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
		} else if table[c]&isSpace != 0 {
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

	return s.errorf(`invalid character %q in place of value start`)
}

func beforeFirstObjectKey(s *Scanner, c byte) Event {
	if table[c]&isSpace != 0 {
		return Space
	} else if c == '"' {
		s.state = afterQuote
		s.end = KeyEnd
		s.push(afterObjectKey)
		return KeyStart
	} else if c == '}' {
		return s.delay(ObjectEnd)
	}

	return s.errorf(`invalid character %q in object`, c)
}

func afterObjectKey(s *Scanner, c byte) Event {
	if table[c]&isSpace != 0 {
		return Space
	} else if c == ':' {
		s.state = beforeValue
		s.push(afterObjectValue)
		return None
	}

	return s.errorf(`invalid character %q after object key`, c)
}

func afterObjectValue(s *Scanner, c byte) Event {
	if table[c]&isSpace != 0 {
		return Space
	} else if c == ',' {
		s.state = afterObjectComma
		return None
	} else if c == '}' {
		return s.delay(ObjectEnd)
	}

	return s.errorf(`invalid character %q after object value`, c)
}

func afterObjectComma(s *Scanner, c byte) Event {
	if table[c]&isSpace != 0 {
		return Space
	} else if c == '"' {
		s.state = afterQuote
		s.end = KeyEnd
		s.push(afterObjectKey)
		return KeyStart
	}

	return s.errorf(`invalid character %q in place of object key`, c)
}

func beforeFirstArrayElement(s *Scanner, c byte) Event {
	if table[c]&isSpace != 0 {
		return Space
	} else if c == ']' {
		return s.delay(ArrayEnd)
	}

	s.push(afterArrayElement)
	return beforeValue(s, c)
}

func afterArrayElement(s *Scanner, c byte) Event {
	if table[c]&isSpace != 0 {
		return Space
	} else if c == ',' {
		s.state = beforeValue
		s.push(afterArrayElement)
		return None
	} else if c == ']' {
		return s.delay(ArrayEnd)
	}

	return s.errorf(`invalid character %q after array element`, c)
}

func afterQuote(s *Scanner, c byte) Event {
	if c == '"' {
		// At this point, s.end has already been set to either StringEnd or
		// KeyEnd depending on the previous state function.
		s.state = delayed
		return None
	} else if c == '\\' {
		s.state = afterEsc
		return None
	} else if c >= 0x20 {
		return None
	}

	return s.errorf(`invalid character %q in string literal`, c)
}

func afterEsc(s *Scanner, c byte) Event {
	if table[c]&isEsc != 0 {
		s.state = afterQuote
		return None
	} else if c == 'u' {
		s.state = afterEscU
		return None
	}

	return s.errorf(`invalid character %q in character escape`, c)
}

func afterEscU(s *Scanner, c byte) Event {
	if table[c]&isHex != 0 {
		s.state = afterEscU1
		return None
	}

	return s.errorf(`invalid character %q in hexadecimal character escape`, c)
}

func afterEscU1(s *Scanner, c byte) Event {
	if table[c]&isHex != 0 {
		s.state = afterEscU12
		return None
	}

	return s.errorf(`invalid character %q in hexadecimal character escape`, c)
}

func afterEscU12(s *Scanner, c byte) Event {
	if table[c]&isHex != 0 {
		s.state = afterEscU123
		return None
	}

	return s.errorf(`invalid character %q in hexadecimal character escape`, c)
}

func afterEscU123(s *Scanner, c byte) Event {
	if table[c]&isHex != 0 {
		s.state = afterQuote
		return None
	}

	return s.errorf(`invalid character %q in hexadecimal character escape`, c)
}

func afterMinus(s *Scanner, c byte) Event {
	if c == '0' {
		s.state = afterZero
		return None
	} else if '1' <= c && c <= '9' {
		s.state = afterDigit
		return None
	}

	return s.errorf(`invalid character %q after "-"`, c)
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
	if table[c]&isDigit != 0 {
		return None
	}

	return afterZero(s, c)
}

func afterDot(s *Scanner, c byte) Event {
	if table[c]&isDigit != 0 {
		s.state = afterDotDigit
		return None
	}

	return s.errorf(`invalid character %q after decimal point in numeric literal`, c)
}

func afterDotDigit(s *Scanner, c byte) Event {
	if table[c]&isDigit != 0 {
		return None
	} else if c == 'e' || c == 'E' {
		s.state = afterE
		return None
	}

	return s.next(c) | NumberEnd
}

func afterE(s *Scanner, c byte) Event {
	if table[c]&isDigit != 0 {
		s.state = afterEDigit
		return None
	} else if c == '-' || c == '+' {
		s.state = afterESign
		return None
	}

	return s.errorf(`invalid character %q in exponent of numeric literal`, c)
}

func afterESign(s *Scanner, c byte) Event {
	if table[c]&isDigit != 0 {
		s.state = afterEDigit
		return None
	}

	return s.errorf(`invalid character %q in exponent of numeric literal`, c)
}

func afterEDigit(s *Scanner, c byte) Event {
	if table[c]&isDigit != 0 {
		return None
	}

	return s.next(c) | NumberEnd
}

func afterT(s *Scanner, c byte) Event {
	if c == 'r' {
		s.state = afterTr
		return None
	}

	return s.errorf(`invalid character %q after "t"`, c)
}

func afterTr(s *Scanner, c byte) Event {
	if c == 'u' {
		s.state = afterTru
		return None
	}

	return s.errorf(`invalid character %q after "tr"`, c)
}

func afterTru(s *Scanner, c byte) Event {
	if c == 'e' {
		return s.delay(BoolEnd)
	}

	return s.errorf(`invalid character %q after "tru"`, c)
}

func afterF(s *Scanner, c byte) Event {
	if c == 'a' {
		s.state = afterFa
		return None
	}

	return s.errorf(`invalid character %q after "f"`, c)
}

func afterFa(s *Scanner, c byte) Event {
	if c == 'l' {
		s.state = afterFal
		return None
	}

	return s.errorf(`invalid character %q after "fa"`, c)
}

func afterFal(s *Scanner, c byte) Event {
	if c == 's' {
		s.state = afterFals
		return None
	}

	return s.errorf(`invalid character %q after "fal"`, c)
}

func afterFals(s *Scanner, c byte) Event {
	if c == 'e' {
		return s.delay(BoolEnd)
	}

	return s.errorf(`invalid character %q after "fals"`, c)
}

func afterN(s *Scanner, c byte) Event {
	if c == 'u' {
		s.state = afterNu
		return None
	}

	return s.errorf(`invalid character %q after "n"`, c)
}

func afterNu(s *Scanner, c byte) Event {
	if c == 'l' {
		s.state = afterNul
		return None
	}

	return s.errorf(`invalid character %q after "nu"`, c)
}

func afterNul(s *Scanner, c byte) Event {
	if c == 'l' {
		return s.delay(NullEnd)
	}

	return s.errorf(`invalid character %q after "nul"`, c)
}

func delayed(s *Scanner, c byte) Event {
	return s.next(c) | s.end
}

func afterTopValue(s *Scanner, c byte) Event {
	if table[c]&isSpace != 0 {
		return Space
	}

	return s.errorf(`invalid character %q after top-level value`, c)
}

func afterError(s *Scanner, c byte) Event {
	return Error
}

// Character type lookup table.
var table = [256]byte{}

const (
	isSpace = 1 << iota
	isDigit
	isHex
	isEsc
)

func init() {
	for i := 0; i < 256; i++ {
		c := byte(i)

		if c == ' ' || c == '\n' || c == '\t' || c == '\r' {
			table[i] |= isSpace
		}
		if '0' <= c && c <= '9' {
			table[i] |= isDigit
		}
		if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
			table[i] |= isHex
		}
		if c == 'b' || c == 'f' || c == 'n' || c == 'r' || c == 't' ||
			c == '\\' || c == '/' || c == '"' {
			table[i] |= isEsc
		}
	}
}
