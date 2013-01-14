package jo

const (
	_Literal = 1 << 0
	_Start   = 1 << 2
	_End     = 1 << 3

	Continue = (iota << 4)
	Done     = (iota << 4)

	ObjectStart = (iota << 4) | _Start
	ObjectEnd   = (iota << 4) | _End
	KeyStart    = (iota << 4) | _Start
	KeyEnd      = (iota << 4) | _End
	ArrayStart  = (iota << 4) | _Start
	ArrayEnd    = (iota << 4) | _End

	StringStart = (iota << 4) | _Literal | _Start
	StringEnd   = (iota << 4) | _Literal | _End
	NumberStart = (iota << 4) | _Literal | _Start
	NumberEnd   = (iota << 4) | _Literal | _End
	BoolStart   = (iota << 4) | _Literal | _Start
	BoolEnd     = (iota << 4) | _Literal | _End
	NullStart   = (iota << 4) | _Literal | _Start
	NullEnd     = (iota << 4) | _Literal | _End

	SyntaxError = (iota << 4)
)

// Represents a change in parser context.
type Event int

var names = map[Event]string{
	Continue:    "Continue",
	Done:        "Done",
	ObjectStart: "ObjectStart",
	ObjectEnd:   "ObjectEnd",
	KeyStart:    "KeyStart",
	KeyEnd:      "KeyEnd",
	ArrayStart:  "ArrayStart",
	ArrayEnd:    "ArrayEnd",
	StringStart: "StringStart",
	StringEnd:   "StringEnd",
	NumberStart: "NumberStart",
	NumberEnd:   "NumberEnd",
	BoolStart:   "BoolStart",
	BoolEnd:     "BoolEnd",
	NullStart:   "NullStart",
	NullEnd:     "NullEnd",
	SyntaxError: "SyntaxError",
}

// Returns the event's name.
func (e Event) String() string {
	name, ok := names[e]
	if !ok {
		return "INVALID"
	}

	return name
}

// Returns true if the event regards a literal value.
func (e Event) IsLiteral() bool {
	return e&_Literal != 0
}

// Returns true if the event marks the start of a value.
func (e Event) IsStart() bool {
	return e&_Start != 0
}

// Returns true if the event marks the end of a value.
func (e Event) IsEnd() bool {
	return e&_End != 0
}
