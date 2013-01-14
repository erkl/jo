package jo

const (
	Continue Event = iota
	Done

	ObjectStart
	ObjectEnd
	KeyStart
	KeyEnd
	ArrayStart
	ArrayEnd

	StringStart
	StringEnd
	NumberStart
	NumberEnd
	BoolStart
	BoolEnd
	NullStart
	NullEnd

	SyntaxError
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
