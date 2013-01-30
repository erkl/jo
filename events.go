package jo

const (
	Primitive Event = 1 << iota
	Composite
	Start
	End

	shift
)

const (
	Continue Event = (iota << uint(shift))
	Done           = (iota << uint(shift))

	ObjectStart = (iota << uint(shift)) | Composite | Start
	ObjectEnd   = (iota << uint(shift)) | Composite | End
	ArrayStart  = (iota << uint(shift)) | Composite | Start
	ArrayEnd    = (iota << uint(shift)) | Composite | End

	KeyStart = (iota << uint(shift)) | Start
	KeyEnd   = (iota << uint(shift)) | End

	StringStart = (iota << uint(shift)) | Primitive | Start
	StringEnd   = (iota << uint(shift)) | Primitive | End
	NumberStart = (iota << uint(shift)) | Primitive | Start
	NumberEnd   = (iota << uint(shift)) | Primitive | End
	BoolStart   = (iota << uint(shift)) | Primitive | Start
	BoolEnd     = (iota << uint(shift)) | Primitive | End
	NullStart   = (iota << uint(shift)) | Primitive | Start
	NullEnd     = (iota << uint(shift)) | Primitive | End

	SyntaxError = (iota << uint(shift))
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
//
//   BoolStart.String()  // -> "BoolStart"
//   Event(-1).String()  // -> "INVALID"
func (e Event) String() string {
	name, ok := names[e]
	if !ok {
		return "INVALID"
	}

	return name
}
