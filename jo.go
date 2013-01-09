package jo

// Parsing events.
type Event int

const (
	Continue = iota

	StringStart
	StringEnd
	BoolStart
	BoolEnd
	NullStart
	NullEnd

	SyntaxError
)

// Parser states.
const (
	_StateValue = iota

	_StateStringUnicode  // \u
	_StateStringUnicode2 // \u1
	_StateStringUnicode3 // \u12
	_StateStringUnicode4 // \u123
	_StateString         // "
	_StateStringEscaped  // \

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

	_StateDone
	_StateSyntaxError
)

// Our own little implementation of the `error` interface.
type syntaxError string

func (e syntaxError) Error() string {
	return string(e)
}

// Parser state machine.
type Parser struct {
	state int
	queue []int
	err   error
}

// Parses a byte slice containing JSON data. Returns the number of bytes
// read and an appropriate Event.
func (p *Parser) Parse(input []byte) (int, Event) {
	for i, b := range input {
		switch p.state {
		case _StateValue:
			switch b {
			case '"':
				p.state = _StateString
				return i + 1, StringStart
			case 't':
				p.state = _StateTrue
				return i + 1, BoolStart
			case 'f':
				p.state = _StateFalse
				return i + 1, BoolStart
			case 'n':
				p.state = _StateNull
				return i + 1, NullStart
			default:
				return i, p.error(`_StateValue: @todo`)
			}

		case _StateStringUnicode, _StateStringUnicode2,
			_StateStringUnicode3, _StateStringUnicode4:
			switch {
			case '0' <= b && b <= '9':
			case 'a' <= b && b <= 'f':
			case 'A' <= b && b <= 'F':
			default:
				return i, p.error(`_StateStringUnicodeX: @todo`)
			}

			p.state++ // note that `_StateString == (_StateStringUnicode4 + 1)`

		case _StateString:
			switch {
			case b == '"':
				p.state = p.next()
				return i + 1, StringEnd
			case b == '\\':
				p.state = _StateStringEscaped
			case b < 0x20:
				return i, p.error(`_StateString: @todo`)
			}

		case _StateStringEscaped:
			switch b {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				p.state = _StateString
			case 'u':
				p.state = _StateStringUnicode
			default:
				return i, p.error(`_StateStringEscaped: @todo`)
			}

		case _StateTrue:
			if b != 'r' {
				return i, p.error(`_StateTrue: @todo`)
			}
			p.state++

		case _StateTrue2:
			if b != 'u' {
				return i, p.error(`_StateTrue2: @todo`)
			}
			p.state++

		case _StateTrue3:
			if b != 'e' {
				return i, p.error(`_StateTrue3: @todo`)
			}
			p.state = p.next()

			return i + 1, BoolEnd

		case _StateFalse:
			if b != 'a' {
				return i, p.error(`_StateFalse: @todo`)
			}
			p.state++

		case _StateFalse2:
			if b != 'l' {
				return i, p.error(`_StateFalse2: @todo`)
			}
			p.state++

		case _StateFalse3:
			if b != 's' {
				return i, p.error(`_StateFalse3: @todo`)
			}
			p.state++

		case _StateFalse4:
			if b != 'e' {
				return i, p.error(`_StateFalse4: @todo`)
			}
			p.state = p.next()

			return i + 1, BoolEnd

		case _StateNull:
			if b != 'u' {
				return i, p.error(`_StateNull: @todo`)
			}
			p.state++

		case _StateNull2:
			if b != 'l' {
				return i, p.error(`_StateNull2: @todo`)
			}
			p.state++

		case _StateNull3:
			if b != 'l' {
				return i, p.error(`_StateNull3: @todo`)
			}
			p.state = p.next()

			return i + 1, NullEnd

		case _StateDone:
			return i, p.error(`_StateDone: @todo`)

		default:
			panic(`invalid state`)
		}
	}

	return len(input) - 1, Continue
}

// Pops the next state off the parser struct's queue.
func (p *Parser) next() int {
	length := len(p.queue)

	// with the "state queue" empty, we can only wait for EOF
	if length == 0 {
		return _StateDone
	}

	state := p.queue[length]
	p.queue = p.queue[:length-1]

	return state
}

// Insert a new state at the top of the queue.
func (p *Parser) push(state int) {
	p.queue = append(p.queue, state)
}

// Registers a syntax error. Always returns a SyntaxError event.
func (p *Parser) error(message string) Event {
	p.err = syntaxError(message)
	return SyntaxError
}
