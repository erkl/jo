// Light-weight, event driven JSON parser.
package jo

import (
	"errors"
)

const (
	_StateValue = iota
	_StateDone

	_StateObjectKeyOrBrace   // {
	_StateObjectKeyDone      // {"foo
	_StateObjectColon        // {"foo"
	_StateObjectCommaOrBrace // {"foo":"bar"
	_StateObjectKey          // {"foo":"bar",

	_StateArrayValueOrBracket // [
	_StateArrayCommaOrBracket // ["any value"

	// leading whitespace must be be consumed before any of the above
	// states are processed
	_IgnoreSpace

	_StateStringUnicode  // "\u
	_StateStringUnicode2 // "\u1
	_StateStringUnicode3 // "\u12
	_StateStringUnicode4 // "\u123
	_StateString         // "
	_StateStringDone     // "foo
	_StateStringEscaped  // "\

	_StateNumberNegative      // -
	_StateNumberZero          // 0
	_StateNumber              // 123
	_StateNumberDotFirstDigit // 123.
	_StateNumberDotDigit      // 123.4
	_StateNumberExpSign       // 123e
	_StateNumberExpFirstDigit // 123e+
	_StateNumberExpDigit      // 123e+1

	_StateTrue  // t
	_StateTrue2 // tr
	_StateTrue3 // tru

	_StateFalse  // f
	_StateFalse2 // fa
	_StateFalse3 // fal
	_StateFalse4 // fals

	_StateNull  // n
	_StateNull2 // nu
	_StateNull3 // nul
)

// Parser state machine. Requires no initialization before use.
type Parser struct {
	state int
	queue []int

	depth int
	limit int
	drop  int
	empty int

	property bool
}

// Parses a byte slice containing JSON data. Returns the number of bytes
// read and an appropriate Event.
//
// The third return value will be nil unless the second is a SyntaxError
// event, in which case it's an error describing the syntax error
// encountered.
func (p *Parser) Parse(input []byte) (int, Event, error) {
	for i := 0; i < len(input); i++ {
		event := Continue
		err := error(nil)

		s := p.state
		b := input[i]

		if s < _IgnoreSpace && isSpace(b) {
			continue
		}

		switch s {
		case _StateValue:
			if b == '{' {
				event = ObjectStart
				p.state = _StateObjectKeyOrBrace
			} else if b == '[' {
				event = ArrayStart
				p.state = _StateArrayValueOrBracket
			} else if b == '"' {
				event = StringStart
				p.state = _StateString
				p.push(_StateStringDone)
			} else if b == '-' {
				event = NumberStart
				p.state = _StateNumberNegative
			} else if b == '0' {
				event = NumberStart
				p.state = _StateNumberZero
			} else if '1' <= b && b <= '9' {
				event = NumberStart
				p.state = _StateNumber
			} else if b == 't' {
				event = BoolStart
				p.state = _StateTrue
			} else if b == 'f' {
				event = BoolStart
				p.state = _StateFalse
			} else if b == 'n' {
				event = NullStart
				p.state = _StateNull
			} else {
				event = SyntaxError
				err = errors.New(`expected beginning of JSON value`)
			}

		case _StateObjectKeyOrBrace:
			if b == '}' {
				event = ObjectEnd
				p.state = p.next()
				break
			}

			// if it's not a brace, it must be a key
			p.state = _StateObjectKey
			fallthrough

		case _StateObjectKey:
			if b == '"' {
				event = KeyStart
				p.state = _StateString
				p.push(_StateObjectKeyDone)
			} else {
				event = SyntaxError
				err = errors.New(`expected object key`)
			}

		case _StateObjectKeyDone:
			// we wouldn't be here unless b == '"', so we can avoid
			// checking it again
			event = KeyEnd
			p.state = _StateObjectColon

		case _StateObjectColon:
			if b == ':' {
				p.state = _StateValue
				p.push(_StateObjectCommaOrBrace)
			} else {
				event = SyntaxError
				err = errors.New(`expected ':' after object key`)
			}

		case _StateObjectCommaOrBrace:
			if b == ',' {
				p.state = _StateObjectKey
			} else if b == '}' {
				event = ObjectEnd
				p.state = p.next()
			} else {
				event = SyntaxError
				err = errors.New(`expected ',' or '}' after object value`)
			}

		case _StateArrayValueOrBracket:
			if b == ']' {
				event = ArrayEnd
				p.state = p.next()
			} else {
				p.state = _StateValue
				p.push(_StateArrayCommaOrBracket)

				// rewind and let _StateValue parse this byte for us
				i--
			}

		case _StateArrayCommaOrBracket:
			if b == ',' {
				p.state = _StateValue
				p.push(_StateArrayCommaOrBracket)
			} else if b == ']' {
				event = ArrayEnd
				p.state = p.next()
			} else {
				event = SyntaxError
				err = errors.New(`expected ',' or ']' after array value`)
			}

		case _StateStringUnicode,
			_StateStringUnicode2,
			_StateStringUnicode3,
			_StateStringUnicode4:
			if isHex(b) {
				// move on to the next unicode byte state, or back to
				// `_StateString` if this was the fourth hexadecimal
				// character after "\u"
				p.state++
			} else {
				event = SyntaxError
				err = errors.New(`expected four hexadecimal chars after "\u"`)
			}

		case _StateString:
			if b == '"' {
				// forget we saw the double quote, let the next state
				// "discover" it instead
				i--
				p.state = p.next()
			} else if b == '\\' {
				p.state = _StateStringEscaped
			} else if b < 0x20 {
				event = SyntaxError
				err = errors.New(`expected valid string character`)
			}

		case _StateStringDone:
			// we wouldn't be here unless b == '"', so we can avoid
			// checking it again
			event = StringEnd
			p.state = p.next()

		case _StateStringEscaped:
			switch b {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				p.state = _StateString
			case 'u':
				p.state = _StateStringUnicode
			default:
				event = SyntaxError
				err = errors.New(`expected valid escape sequence after '\'`)
			}

		case _StateNumberNegative:
			if b == '0' {
				p.state = _StateNumberZero
			} else if '1' <= b && b <= '9' {
				p.state = _StateNumber
			} else {
				event = SyntaxError
				err = errors.New(`digit after '-'`)
			}

		case _StateNumber:
			if isDecimal(b) {
				break
			}

			// the same limits apply here as in _StateNumberZero
			fallthrough

		case _StateNumberZero:
			if b == '.' {
				p.state = _StateNumberDotFirstDigit
			} else if b == 'e' || b == 'E' {
				p.state = _StateNumberExpSign
			} else {
				event = NumberEnd
				p.state = p.next()

				// rewind a byte, because the character we encountered was
				// not part of the number
				i--
			}

		case _StateNumberDotFirstDigit:
			if isDecimal(b) {
				p.state = _StateNumberDotDigit
			} else {
				event = SyntaxError
				err = errors.New(`expected digit after dot in number`)
			}

		case _StateNumberDotDigit:
			if b == 'e' || b == 'E' {
				p.state = _StateNumberExpSign
			} else if !isDecimal(b) {
				event = NumberEnd
				p.state = p.next()

				// rewind a byte, because the character we encountered was
				// not part of the number
				i--
			}

		case _StateNumberExpSign:
			p.state = _StateNumberExpFirstDigit
			if b == '+' || b == '-' {
				break
			}
			fallthrough

		case _StateNumberExpFirstDigit:
			if !isDecimal(b) {
				event = SyntaxError
				err = errors.New(`expected digit after exponent in number`)
			} else {
				p.state++
			}

		case _StateNumberExpDigit:
			if !isDecimal(b) {
				event = NumberEnd
				p.state = p.next()

				// rewind a byte, because the character we encountered was
				// not part of the number
				i--
			}

		case _StateTrue:
			if b == 'r' {
				p.state = _StateTrue2
			} else {
				event = SyntaxError
				err = errors.New(`expected 'r' in literal true`)
			}

		case _StateTrue2:
			if b == 'u' {
				p.state = _StateTrue3
			} else {
				event = SyntaxError
				err = errors.New(`expected 'u' in literal true`)
			}

		case _StateTrue3:
			if b == 'e' {
				event = BoolEnd
				p.state = p.next()
			} else {
				event = SyntaxError
				err = errors.New(`expected 'e' in literal true`)
			}

		case _StateFalse:
			if b == 'a' {
				p.state = _StateFalse2
			} else {
				event = SyntaxError
				err = errors.New(`expected 'a' in literal false`)
			}

		case _StateFalse2:
			if b == 'l' {
				p.state = _StateFalse3
			} else {
				event = SyntaxError
				err = errors.New(`expected 'l' in literal false`)
			}

		case _StateFalse3:
			if b == 's' {
				p.state = _StateFalse4
			} else {
				event = SyntaxError
				err = errors.New(`expected 's' in literal false`)
			}

		case _StateFalse4:
			if b == 'e' {
				event = BoolEnd
				p.state = p.next()
			} else {
				event = SyntaxError
				err = errors.New(`expected 'e' in literal false`)
			}

		case _StateNull:
			if b == 'u' {
				p.state = _StateNull2
			} else {
				event = SyntaxError
				err = errors.New(`expected 'u' in literal false`)
			}

		case _StateNull2:
			if b == 'l' {
				p.state = _StateNull3
			} else {
				event = SyntaxError
				err = errors.New(`expected 'l' in literal false`)
			}

		case _StateNull3:
			if b == 'l' {
				event = NullEnd
				p.state = p.next()
			} else {
				event = SyntaxError
				err = errors.New(`expected 'l' in literal false`)
			}

		case _StateDone:
			event = SyntaxError
			err = errors.New(`expected nothing after top-level value`)

		default:
			panic(`invalid state`)
		}

		// if this byte didn't yield an event, try the next
		if event == Continue {
			continue
		}

		// in the case of a syntax error, don't consume the offending byte
		if event == SyntaxError {
			return i, SyntaxError, err
		}

		// make sure p.depth is accurate
		switch {
		case event == KeyStart:
			p.depth++
			p.property = true
		case event&(Composite|Primitive) != 0 && event&Start != 0:
			if !p.property {
				p.depth++
			} else {
				p.property = false
			}
		case event&(Composite|Primitive) != 0 && event&End != 0:
			p.depth--
		}

		// determine if we should skip this event
		if p.drop != 0 || p.empty != 0 {
			// silence all events for values below the depth limit
			if p.depth > p.limit {
				continue
			}

			p.limit--

			if p.drop > 0 {
				p.drop--
				continue
			} else {
				p.empty--
				if event&End == 0 {
					continue
				}
			}
		}

		return i + 1, event, nil
	}

	return len(input), Continue, nil
}

// Informs the parser not to expect any further input (EOF).
//
// Returns a SyntaxError event and a descriptive error if invoked before the
// top-level value has been completely parsed. Otherwise returns dangling
// NumberEnd events, or Done.
func (p *Parser) End() (Event, error) {
	switch p.state {
	case _StateNumberZero,
		_StateNumber,
		_StateNumberDotDigit,
		_StateNumberExpDigit:
		p.state = _StateDone
		return NumberEnd, nil
	case _StateDone:
		return Done, nil
	}

	return SyntaxError, errors.New(`unexpected end of input`)
}

// Returns the current depth of nested objects and arrays. Will be 0 for
// top-level literal values.
func (p *Parser) Depth() int {
	return p.depth
}

// Skip provides advanced functionality to skip events, based on the depth of
// nested objects. It lets the user both drop all start and end events, or
// only allow the parser to emit end events, until parsing has reached a depth
// N levels up from where Skip was triggered.
//
// The number of levels to drop or make empty will be applied in order. So in
// the following example, events for the current level and the one above will
// be completely dropped. Then, the object/array above that won't emit any
// more properties/elements. After that object/array, parsing will resume as
// usual.
//
//   p.Skip(2, 1)
//
// Here's a more thourough example:
//
//   in := []byte(`[{"foo":"bar"},{"numbers":[1,2,3]}]`)
//   p := Parser{}
//    
//   p.Parse(in[0:])   // -> (1, ArrayStart, nil)
//   p.Parse(in[1:])   // -> (1, ObjectStart, nil)
//    
//   p.Skip(0, 1)      // we don't care about what's in this object;
//                     // skip all its properties
//    
//   p.Parse(in[2:])   // -> (12, ObjectEnd, nil)
//   p.Parse(in[14:])  // -> (2, ObjectStart, nil)
//   p.Parse(in[16:])  // -> (1, KeyStart, nil)
//    
//   p.Skip(2, 0)      // completely drop this key/value pair, and
//                     // whatever remains of the parent object
//    
//   p.Parse(in[17:])  // -> (18, ArrayEnd, nil)
//   p.End()           // -> (Done, nil)
//
// Panics when a) either drop or empty is negative, or b) drop + empty is
// greater than the current depth.
func (p *Parser) Skip(drop, empty int) {
	if drop < 0 || empty < 0 {
		panic(`both drop and empty must be positive`)
	}
	if drop+empty > p.depth {
		panic(`drop + empty must not be greater than the current depth`)
	}

	// Parser.Skip(1, 0) should be equal to Parser.Skip(0, 1) if invoked inside,
	// or just after, an object key; it wouldn't make sense to receive an end
	// event we haven't seen the start of
	if p.property && drop == 0 && empty > 0 {
		drop++
		empty--
	}

	p.drop = drop
	p.empty = empty

	p.limit = p.depth - 1
}

// Puts a new state at the top of the queue.
func (p *Parser) push(state int) {
	p.queue = append(p.queue, state)
}

// Fetches the next state in the queue.
func (p *Parser) next() int {
	length := len(p.queue)

	// if the state queue is empty, the top level value has ended
	if length == 0 {
		return _StateDone
	}

	state := p.queue[length-1]
	p.queue = p.queue[:length-1]

	return state
}

// Returns true if b is a whitespace character.
func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// Returns true if b is a hexadecimal character.
func isHex(b byte) bool {
	return isDecimal(b) || ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F')
}

// Returns true if b is a decimal digit.
func isDecimal(b byte) bool {
	return '0' <= b && b <= '9'
}
