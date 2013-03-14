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
	_Value int = iota
	_Done

	_ObjectKeyOrBrace   // {
	_ObjectKeyDone      // {"foo
	_ObjectColon        // {"foo"
	_ObjectCommaOrBrace // {"foo":"bar"
	_ObjectKey          // {"foo":"bar",

	_ArrayValueOrBracket // [
	_ArrayCommaOrBracket // ["any value"

	// leading whitespace must be be consumed before any of
	// the states listed above are processed
	__CONSUME_SPACE__

	_StringUnicode  // "\u
	_StringUnicode2 // "\u1
	_StringUnicode3 // "\u12
	_StringUnicode4 // "\u123
	_String         // "
	_StringDone     // "foo
	_StringEscaped  // "\

	_NumberNegative      // -
	_NumberZero          // 0
	_Number              // 123
	_NumberDotFirstDigit // 123.
	_NumberDotDigit      // 123.4
	_NumberExpSign       // 123e
	_NumberExpFirstDigit // 123e+
	_NumberExpDigit      // 123e+1

	_True  // t
	_True2 // tr
	_True3 // tru

	_False  // f
	_False2 // fa
	_False3 // fal
	_False4 // fals

	_Null  // n
	_Null2 // nu
	_Null3 // nul
)

var expected = map[int]string{
	_Value:               "start of JSON value",
	_Done:                "end of input",
	_ObjectKeyOrBrace:    "object key or '}'",
	_ObjectColon:         "':'",
	_ObjectCommaOrBrace:  "',' or '}'",
	_ObjectKey:           "object key",
	_ArrayValueOrBracket: "array element or ']'",
	_ArrayCommaOrBracket: "',' or ']'",
	_StringUnicode:       "hexadecimal digit",
	_StringUnicode2:      "hexadecimal digit",
	_StringUnicode3:      "hexadecimal digit",
	_StringUnicode4:      "hexadecimal digit",
	_String:              "valid string character or '\"'",
	_StringEscaped:       "'b', 'f', 'n', 'r', 't', 'u', '\\', '/' or '\"'",
	_NumberNegative:      "digit",
	_NumberZero:          "'.', 'e' or 'E'",
	_NumberDotFirstDigit: "digit",
	_NumberExpSign:       "'-', '+' or digit",
	_NumberExpFirstDigit: "digit",
	_True:                "'r' in literal true",
	_True2:               "'u' in literal true",
	_True3:               "'e' in literal true",
	_False:               "'a' in literal false",
	_False2:              "'l' in literal false",
	_False3:              "'s' in literal false",
	_False4:              "'e' in literal false",
	_Null:                "'u' in literal null",
	_Null2:               "'l' in literal null",
	_Null3:               "'l' in literal null",
}

// Next parses a slice of JSON data, signalling the first change in context by
// returning the number of bytes read and an appropriate event. If parsing
// concludes without any change in parser state, the Continue psuedo-event
// will be returned instead.
//
// When the event returned is SyntaxError, the error return value will describe
// why parsing failed.
func (p *Parser) Next(data []byte) (int, Event, error) {
	for i := 0; i < len(data); i++ {
		ev := Continue
		b := data[i]

		// trim insignificant whitespace
		if p.state < __CONSUME_SPACE__ && isSpace(b) {
			continue
		}

		switch p.state {
		case _Value:
			if b == '{' {
				ev = ObjectStart
				p.state = _ObjectKeyOrBrace
			} else if b == '[' {
				ev = ArrayStart
				p.state = _ArrayValueOrBracket
			} else if b == '"' {
				ev = StringStart
				p.state = _String
				p.push(_StringDone)
			} else if b == '-' {
				ev = NumberStart
				p.state = _NumberNegative
			} else if b == '0' {
				ev = NumberStart
				p.state = _NumberZero
			} else if '1' <= b && b <= '9' {
				ev = NumberStart
				p.state = _Number
			} else if b == 't' {
				ev = BoolStart
				p.state = _True
			} else if b == 'f' {
				ev = BoolStart
				p.state = _False
			} else if b == 'n' {
				ev = NullStart
				p.state = _Null
			} else {
				goto abort
			}

		case _ObjectKeyOrBrace:
			if b == '}' {
				ev = ObjectEnd
				p.state = p.pop()
				break
			}

			// if it's not a brace, it must be a key
			p.state = _ObjectKey
			fallthrough

		case _ObjectKey:
			if b == '"' {
				ev = KeyStart
				p.state = _String
				p.push(_ObjectKeyDone)
			} else {
				goto abort
			}

		case _ObjectKeyDone:
			// we wouldn't be here unless b == '"', so we can avoid
			// checking it again
			ev = KeyEnd
			p.state = _ObjectColon

		case _ObjectColon:
			if b == ':' {
				p.state = _Value
				p.push(_ObjectCommaOrBrace)
			} else {
				goto abort
			}

		case _ObjectCommaOrBrace:
			if b == ',' {
				p.state = _ObjectKey
			} else if b == '}' {
				ev = ObjectEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case _ArrayValueOrBracket:
			if b == ']' {
				ev = ArrayEnd
				p.state = p.pop()
			} else {
				p.state = _Value
				p.push(_ArrayCommaOrBracket)

				// rewind and let _Value parse this byte for us
				i--
			}

		case _ArrayCommaOrBracket:
			if b == ',' {
				p.state = _Value
				p.push(_ArrayCommaOrBracket)
			} else if b == ']' {
				ev = ArrayEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case _StringUnicode,
			_StringUnicode2,
			_StringUnicode3,
			_StringUnicode4:
			if isHex(b) {
				// move on to the next unicode byte state, or back to
				// `_String` if this was the fourth hexadecimal
				// character after "\u"
				p.state++
			} else {
				goto abort
			}

		case _String:
			if b == '"' {
				// forget we saw the double quote, let the next state
				// "discover" it instead
				i--
				p.state = p.pop()
			} else if b == '\\' {
				p.state = _StringEscaped
			} else if b < 0x20 {
				goto abort
			}

		case _StringDone:
			// we wouldn't be here unless b == '"', so we can avoid
			// checking it again
			ev = StringEnd
			p.state = p.pop()

		case _StringEscaped:
			switch b {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				p.state = _String
			case 'u':
				p.state = _StringUnicode
			default:
				goto abort
			}

		case _NumberNegative:
			if b == '0' {
				p.state = _NumberZero
			} else if '1' <= b && b <= '9' {
				p.state = _Number
			} else {
				goto abort
			}

		case _Number:
			if b >= '0' && b <= '9' {
				break
			}
			fallthrough

		case _NumberZero:
			if b == '.' {
				p.state = _NumberDotFirstDigit
			} else if b == 'e' || b == 'E' {
				p.state = _NumberExpSign
			} else {
				ev = NumberEnd
				p.state = p.pop()

				// rewind a byte, because the character we encountered was
				// not part of the number
				i--
			}

		case _NumberDotFirstDigit:
			if b >= '0' && b <= '9' {
				p.state = _NumberDotDigit
			} else {
				goto abort
			}

		case _NumberDotDigit:
			if b == 'e' || b == 'E' {
				p.state = _NumberExpSign
			} else if b < '0' || b > '9' {
				ev = NumberEnd
				p.state = p.pop()

				// rewind a byte, because the character we encountered was
				// not part of the number
				i--
			}

		case _NumberExpSign:
			p.state = _NumberExpFirstDigit
			if b == '+' || b == '-' {
				break
			}
			fallthrough

		case _NumberExpFirstDigit:
			if b < '0' || b > '9' {
				goto abort
			} else {
				p.state++
			}

		case _NumberExpDigit:
			if b < '0' || b > '9' {
				ev = NumberEnd
				p.state = p.pop()

				// rewind a byte, because the character we encountered was
				// not part of the number
				i--
			}

		case _True:
			if b == 'r' {
				p.state = _True2
			} else {
				goto abort
			}

		case _True2:
			if b == 'u' {
				p.state = _True3
			} else {
				goto abort
			}

		case _True3:
			if b == 'e' {
				ev = BoolEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case _False:
			if b == 'a' {
				p.state = _False2
			} else {
				goto abort
			}

		case _False2:
			if b == 'l' {
				p.state = _False3
			} else {
				goto abort
			}

		case _False3:
			if b == 's' {
				p.state = _False4
			} else {
				goto abort
			}

		case _False4:
			if b == 'e' {
				ev = BoolEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case _Null:
			if b == 'u' {
				p.state = _Null2
			} else {
				goto abort
			}

		case _Null2:
			if b == 'l' {
				p.state = _Null3
			} else {
				goto abort
			}

		case _Null3:
			if b == 'l' {
				ev = NullEnd
				p.state = p.pop()
			} else {
				goto abort
			}

		case _Done:
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

		return i + 1, ev, nil

	abort:
		return i, SyntaxError, fmt.Errorf("expected %s, found %q",
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
		return _Done
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
	case _Done:
		return Done, nil
	case _NumberZero, _Number, _NumberDotDigit, _NumberExpDigit:
		if len(p.stack) == 0 {
			p.state = _Done
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
