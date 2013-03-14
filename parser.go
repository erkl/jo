package jo

import (
	"fmt"
)

// Events signal a change in context, for example the start of
// a string literal.
type Event int

const (
	Primitive Event = 1 << (5 + iota)
	Composite
	Start
	End
)

const (
	Continue    Event = iota
	SyntaxError       = iota
	Done              = iota

	ObjectStart = iota | Composite | Start
	ObjectEnd   = iota | Composite | End
	ArrayStart  = iota | Composite | Start
	ArrayEnd    = iota | Composite | End

	KeyStart = iota | Start
	KeyEnd   = iota | End

	StringStart = iota | Primitive | Start
	StringEnd   = iota | Primitive | End
	NumberStart = iota | Primitive | Start
	NumberEnd   = iota | Primitive | End
	BoolStart   = iota | Primitive | Start
	BoolEnd     = iota | Primitive | End
	NullStart   = iota | Primitive | Start
	NullEnd     = iota | Primitive | End
)

// String returns the name of the event.
func (e Event) String() string {
	switch e {
	case Continue:
		return "Continue"
	case SyntaxError:
		return "SyntaxError"
	case Done:
		return "Done"
	case ObjectStart:
		return "ObjectStart"
	case ObjectEnd:
		return "ObjectEnd"
	case KeyStart:
		return "KeyStart"
	case KeyEnd:
		return "KeyEnd"
	case ArrayStart:
		return "ArrayStart"
	case ArrayEnd:
		return "ArrayEnd"
	case StringStart:
		return "StringStart"
	case StringEnd:
		return "StringEnd"
	case NumberStart:
		return "NumberStart"
	case NumberEnd:
		return "NumberEnd"
	case BoolStart:
		return "BoolStart"
	case BoolEnd:
		return "BoolEnd"
	case NullStart:
		return "NullStart"
	case NullEnd:
		return "NullEnd"
	}
	return "<unknown event>"
}

// Parser is the state machine used while parsing a stream of JSON data.
// It requires no initialization before use.
type Parser struct {
	state int
	stack []int
}

const (
	stateValue int = iota
	stateDone

	stateObjectKeyOrBrace   // {
	stateObjectKeyDone      // {"foo
	stateObjectColon        // {"foo"
	stateObjectCommaOrBrace // {"foo":"bar"
	stateObjectKey          // {"foo":"bar",

	stateArrayValueOrBracket // [
	stateArrayCommaOrBracket // ["any value"

	// leading whitespace must be be consumed before any of
	// the states listed above are processed
	__CONSUME_SPACE__

	stateStringEscU    // "\u
	stateStringEscU1   // "\u1
	stateStringEscU12  // "\u12
	stateStringEscU123 // "\u123
	stateString        // "
	stateStringDone    // "foo
	stateStringEsc     // "\

	stateNumberNegative      // -
	stateNumberZero          // 0
	stateNumber              // 123
	stateNumberDotFirstDigit // 123.
	stateNumberDotDigit      // 123.4
	stateNumberExpSign       // 123e
	stateNumberExpFirstDigit // 123e+
	stateNumberExpDigit      // 123e+1

	stateT   // t
	stateTr  // tr
	stateTru // tru

	stateF    // f
	stateFa   // fa
	stateFal  // fal
	stateFals // fals

	stateN   // n
	stateNu  // nu
	stateNul // nul
)

var expected = map[int]string{
	stateValue:               "start of JSON value",
	stateDone:                "end of input",
	stateObjectKeyOrBrace:    "object key or '}'",
	stateObjectColon:         "':'",
	stateObjectCommaOrBrace:  "',' or '}'",
	stateObjectKey:           "object key",
	stateArrayValueOrBracket: "array element or ']'",
	stateArrayCommaOrBracket: "',' or ']'",
	stateStringEscU:          "hexadecimal digit",
	stateStringEscU1:         "hexadecimal digit",
	stateStringEscU12:        "hexadecimal digit",
	stateStringEscU123:       "hexadecimal digit",
	stateString:              "valid string character or '\"'",
	stateStringEsc:           "'b', 'f', 'n', 'r', 't', 'u', '\\', '/' or '\"'",
	stateNumberNegative:      "digit",
	stateNumberZero:          "'.', 'e' or 'E'",
	stateNumberDotFirstDigit: "digit",
	stateNumberExpSign:       "'-', '+' or digit",
	stateNumberExpFirstDigit: "digit",
	stateT:                   "'r' in literal true",
	stateTr:                  "'u' in literal true",
	stateTru:                 "'e' in literal true",
	stateF:                   "'a' in literal false",
	stateFa:                  "'l' in literal false",
	stateFal:                 "'s' in literal false",
	stateFals:                "'e' in literal false",
	stateN:                   "'u' in literal null",
	stateNu:                  "'l' in literal null",
	stateNul:                 "'l' in literal null",
}

// Next parses a slice of JSON data, signalling the first change in context by
// returning the number of bytes read and an appropriate event. If parsing
// concludes without any change in parser state, the Continue psuedo-event
// will be returned instead.
//
// When the event returned is SyntaxError, the error return value will describe
// why parsing failed.
func (p *Parser) Next(data []byte) (int, Event, error) {
	var ev Event = Continue
	var b byte

	for r := 0; r < len(data); r++ {
		b = data[r]

		// trim insignificant whitespace
		if p.state < __CONSUME_SPACE__ && isSpace(b) {
			continue
		}

		switch p.state {
		case stateValue:
			if b == '"' {
				ev = StringStart
				p.state = stateString
				p.push(stateStringDone)
			} else if '1' <= b && b <= '9' {
				ev = NumberStart
				p.state = stateNumber
			} else if b == '{' {
				ev = ObjectStart
				p.state = stateObjectKeyOrBrace
			} else if b == '[' {
				ev = ArrayStart
				p.state = stateArrayValueOrBracket
			} else if b == '0' {
				ev = NumberStart
				p.state = stateNumberZero
			} else if b == '-' {
				ev = NumberStart
				p.state = stateNumberNegative
			} else if b == 't' {
				ev = BoolStart
				p.state = stateT
			} else if b == 'f' {
				ev = BoolStart
				p.state = stateF
			} else if b == 'n' {
				ev = NullStart
				p.state = stateN
			} else {
				goto abort
			}

		case stateObjectKeyOrBrace:
			if b == '}' {
				ev = ObjectEnd
				p.state = p.pop()
				break
			}

			// if it's not a brace, it must be a key
			p.state = stateObjectKey
			fallthrough

		case stateObjectKey:
			if b == '"' {
				ev = KeyStart
				p.state = stateString
				p.push(stateObjectKeyDone)
			} else {
				goto abort
			}

		case stateObjectKeyDone:
			// we wouldn't be here unless b == '"', so we can avoid
			// checking it again
			ev = KeyEnd
			p.state = stateObjectColon

		case stateObjectColon:
			if b == ':' {
				p.state = stateValue
				p.push(stateObjectCommaOrBrace)
			} else {
				goto abort
			}

		case stateObjectCommaOrBrace:
			if b == ',' {
				p.state = stateObjectKey
			} else if b == '}' {
				ev = ObjectEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case stateArrayValueOrBracket:
			if b == ']' {
				ev = ArrayEnd
				p.state = p.pop()
			} else {
				p.state = stateValue
				p.push(stateArrayCommaOrBracket)

				// rewind and let stateValue parse this byte for us
				r--
			}

		case stateArrayCommaOrBracket:
			if b == ',' {
				p.state = stateValue
				p.push(stateArrayCommaOrBracket)
			} else if b == ']' {
				ev = ArrayEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case stateStringEscU:
			if !isHex(b) {
				goto abort
			}
			p.state = stateStringEscU1

		case stateStringEscU1:
			if !isHex(b) {
				goto abort
			}
			p.state = stateStringEscU12

		case stateStringEscU12:
			if !isHex(b) {
				goto abort
			}
			p.state = stateStringEscU123

		case stateStringEscU123:
			if !isHex(b) {
				goto abort
			}
			p.state = stateString

		case stateString:
			if b == '"' {
				// forget we saw the double quote, let the next state
				// "discover" it instead
				r--
				p.state = p.pop()
			} else if b == '\\' {
				p.state = stateStringEsc
			} else if b < 0x20 {
				goto abort
			}

		case stateStringDone:
			// we wouldn't be here unless b == '"', so we can avoid
			// checking it again
			ev = StringEnd
			p.state = p.pop()

		case stateStringEsc:
			switch b {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				p.state = stateString
			case 'u':
				p.state = stateStringEscU
			default:
				goto abort
			}

		case stateNumberNegative:
			if b == '0' {
				p.state = stateNumberZero
			} else if '1' <= b && b <= '9' {
				p.state = stateNumber
			} else {
				goto abort
			}

		case stateNumber:
			if b >= '0' && b <= '9' {
				break
			}
			fallthrough

		case stateNumberZero:
			if b == '.' {
				p.state = stateNumberDotFirstDigit
			} else if b == 'e' || b == 'E' {
				p.state = stateNumberExpSign
			} else {
				ev = NumberEnd
				p.state = p.pop()

				// rewind a byte, because the character we encountered was
				// not part of the number
				r--
			}

		case stateNumberDotFirstDigit:
			if b >= '0' && b <= '9' {
				p.state = stateNumberDotDigit
			} else {
				goto abort
			}

		case stateNumberDotDigit:
			if b == 'e' || b == 'E' {
				p.state = stateNumberExpSign
			} else if b < '0' || b > '9' {
				ev = NumberEnd
				p.state = p.pop()

				// rewind a byte, because the character we encountered was
				// not part of the number
				r--
			}

		case stateNumberExpSign:
			p.state = stateNumberExpFirstDigit
			if b == '+' || b == '-' {
				break
			}
			fallthrough

		case stateNumberExpFirstDigit:
			if b < '0' || b > '9' {
				goto abort
			} else {
				p.state++
			}

		case stateNumberExpDigit:
			if b < '0' || b > '9' {
				ev = NumberEnd
				p.state = p.pop()

				// rewind a byte, because the character we encountered was
				// not part of the number
				r--
			}

		case stateT:
			if b == 'r' {
				p.state = stateTr
			} else {
				goto abort
			}

		case stateTr:
			if b == 'u' {
				p.state = stateTru
			} else {
				goto abort
			}

		case stateTru:
			if b == 'e' {
				ev = BoolEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case stateF:
			if b == 'a' {
				p.state = stateFa
			} else {
				goto abort
			}

		case stateFa:
			if b == 'l' {
				p.state = stateFal
			} else {
				goto abort
			}

		case stateFal:
			if b == 's' {
				p.state = stateFals
			} else {
				goto abort
			}

		case stateFals:
			if b == 'e' {
				ev = BoolEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case stateN:
			if b == 'u' {
				p.state = stateNu
			} else {
				goto abort
			}

		case stateNu:
			if b == 'l' {
				p.state = stateNul
			} else {
				goto abort
			}

		case stateNul:
			if b == 'l' {
				ev = NullEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case stateDone:
			// only whitespace characters are legal after the
			// top-level value
			goto abort

		default:
			panic("invalid state")
		}

		// if this byte didn't yield an event, continue with
		// the next one
		if ev == Continue {
			continue
		}

		return r + 1, ev, nil

	abort:
		return r, SyntaxError, fmt.Errorf("expected %s, found %q",
			expected[p.state], b)
	}

	return len(data), Continue, nil
}

// push puts a new state at the top of the stack.
func (p *Parser) push(state int) {
	p.stack = append(p.stack, state)
}

// pop retrieves the top state in the stack.
func (p *Parser) pop() int {
	length := len(p.stack)

	// if the state queue is empty, the top level value has ended
	if length == 0 {
		return stateDone
	}

	state := p.stack[length-1]
	p.stack = p.stack[:length-1]

	return state
}

// End informs the parser not to expect any further input (i.e. EOF).
//
// A SyntaxError event and a relevant error will be returned if the method
// was invoked before the top-level had been completely parsed. Returns either
// a NumberEnd or Done event otherwise.
func (p *Parser) End() (Event, error) {
	switch p.state {
	case stateDone:
		return Done, nil
	case stateNumberZero, stateNumber, stateNumberDotDigit, stateNumberExpDigit:
		if len(p.stack) == 0 {
			p.state = stateDone
			return NumberEnd, nil
		}
	}

	return SyntaxError, fmt.Errorf("expected %s, found end of input",
		expected[p.state])
}

// Reset resets the parser struct to its initial state. Convenient when parsing
// a stream of more than one JSON value (simply reset the parser after each Done
// event).
func (p *Parser) Reset() {
	*p = Parser{}
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func isHex(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'f') ||
		(b >= 'A' && b <= 'F')
}
