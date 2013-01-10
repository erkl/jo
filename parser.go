package jo

const (
	_StateValue = iota
	_StateDone

	_StateObjectKeyOrBrace   // {
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
	_StateStringEscaped  // "\

	_StateKeyUnicode  // "\u
	_StateKeyUnicode2 // "\u1
	_StateKeyUnicode3 // "\u12
	_StateKeyUnicode4 // "\u123
	_StateKey         // "
	_StateKeyEscaped  // "\

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

// Parser state machine.
type Parser struct {
	state int
	queue []int
	err   error
}

// Our own little implementation of the `error` interface.
type err string

func (e err) Error() string {
	return string(e)
}

// Returns the syntax error stored in the parser, if any.
func (p *Parser) LastError() error {
	return p.err
}

// Parses a byte slice containing JSON data. Returns the number of bytes
// read and an appropriate Event.
func (p *Parser) Parse(input []byte) (int, Event) {
	for i := 0; i < len(input); i++ {
		var event Event = Continue
		var b byte = input[i]

		if p.state < _IgnoreSpace && isSpace(b) {
			continue
		}

		switch p.state {
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
				event = p.error(`expecting beginning of JSON value`)
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
				p.state = _StateKey
			} else {
				event = p.error(`expected object key`)
			}

		case _StateObjectColon:
			if b == ':' {
				p.state = _StateValue
				p.push(_StateObjectCommaOrBrace)
			} else {
				event = p.error(`expected ':' after object key`)
			}

		case _StateObjectCommaOrBrace:
			if b == ',' {
				p.state = _StateObjectKey
			} else if b == '}' {
				event = ObjectEnd
				p.state = p.next()
			} else {
				event = p.error(`expected ',' or '}' after object value`)
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
				event = p.error(`expected ',' or ']' after array value`)
			}

		case _StateStringUnicode, _StateKeyUnicode,
			_StateStringUnicode2, _StateKeyUnicode2,
			_StateStringUnicode3, _StateKeyUnicode3,
			_StateStringUnicode4, _StateKeyUnicode4:
			if isHex(b) {
				// move on to the next unicode byte state, or back to
				// `_State{String/Key}` if this was the fourth hexadecimal
				// character after "\u"
				p.state++
			} else {
				event = p.error(`expected four hexadecimal chars after "\u"`)
			}

		case _StateString, _StateKey:
			if b == '"' {
				if p.state == _StateKey {
					event = KeyEnd
					p.state = _StateObjectColon
				} else {
					event = StringEnd
					p.state = p.next()
				}
			} else if b == '\\' {
				// continue with `_State{String,Key}Escaped`
				p.state++
			} else if b < 0x20 {
				event = p.error(`expected valid string character`)
			}

		case _StateStringEscaped, _StateKeyEscaped:
			switch b {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				// jump back to _State{String,Key}
				p.state--
			case 'u':
				// continue to `_State{String,Key}Unicode`
				p.state -= 5
			default:
				event = p.error(`expected valid escape sequence after '\'`)
			}

		case _StateNumberNegative:
			if b == '0' {
				p.state = _StateNumberZero
			} else if '1' <= b && b <= '9' {
				p.state = _StateNumber
			} else {
				event = p.error(`digit after '-'`)
			}

		case _StateNumber:
			if isDigit(b) {
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
			if isDigit(b) {
				p.state = _StateNumberDotDigit
			} else {
				event = p.error(`expected digit after dot in number`)
			}

		case _StateNumberDotDigit:
			if b == 'e' || b == 'E' {
				p.state = _StateNumberExpSign
			} else if !isDigit(b) {
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
			if !isDigit(b) {
				event = p.error(`expected digit after exponent in number`)
			} else {
				p.state++
			}

		case _StateNumberExpDigit:
			if !isDigit(b) {
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
				event = p.error(`expected 'r' in literal true`)
			}

		case _StateTrue2:
			if b == 'u' {
				p.state = _StateTrue3
			} else {
				return i, p.error(`expected 'u' in literal true`)
			}

		case _StateTrue3:
			if b == 'e' {
				event = BoolEnd
				p.state = p.next()
			} else {
				return i, p.error(`expected 'e' in literal true`)
			}

		case _StateFalse:
			if b == 'a' {
				p.state = _StateFalse2
			} else {
				event = p.error(`expected 'a' in literal false`)
			}

		case _StateFalse2:
			if b == 'l' {
				p.state = _StateFalse3
			} else {
				event = p.error(`expected 'l' in literal false`)
			}

		case _StateFalse3:
			if b == 's' {
				p.state = _StateFalse4
			} else {
				event = p.error(`expected 's' in literal false`)
			}

		case _StateFalse4:
			if b == 'e' {
				event = BoolEnd
				p.state = p.next()
			} else {
				event = p.error(`expected 'e' in literal false`)
			}

		case _StateNull:
			if b == 'u' {
				p.state = _StateNull2
			} else {
				event = p.error(`expected 'u' in literal false`)
			}

		case _StateNull2:
			if b == 'l' {
				p.state = _StateNull3
			} else {
				event = p.error(`expected 'l' in literal false`)
			}

		case _StateNull3:
			if b == 'l' {
				event = NullEnd
				p.state = p.next()
			} else {
				event = p.error(`expected 'l' in literal false`)
			}

		case _StateDone:
			return i, p.error(`expected nothing after top-level value`)

		default:
			panic(`invalid state`)
		}

		if event == SyntaxError {
			// don't consume the byte that caused the error
			return i, SyntaxError
		}

		if event != Continue {
			return i + 1, event
		}
	}

	return len(input), Continue
}

// Informs the parser not to expect any further input. Will return pending
// NumberEnd events (if there are any), or a SyntaxError event if EOF was
// not expected -- otherwise Done.
func (p *Parser) Eof() Event {
	switch p.state {
	case _StateNumberZero,
		_StateNumber,
		_StateNumberDotDigit,
		_StateNumberExpDigit:
		p.state = _StateDone
		return NumberEnd
	case _StateDone:
		return Done
	}
	return p.error(`.Eof(): @todo`)
}

// Pops the next state off the parser struct's queue.
func (p *Parser) next() int {
	length := len(p.queue)

	// with the "state queue" empty, we can only wait for EOF
	if length == 0 {
		return _StateDone
	}

	state := p.queue[length-1]
	p.queue = p.queue[:length-1]

	return state
}

// Insert a new state at the top of the queue.
func (p *Parser) push(state int) {
	p.queue = append(p.queue, state)
}

// Registers a syntax error. Always returns a SyntaxError event.
func (p *Parser) error(message string) Event {
	p.err = err(message)
	return SyntaxError
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func isHex(b byte) bool {
	return isDigit(b) || ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F')
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}
