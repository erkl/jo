package jo

// Events signal changes in scanning state.
type Event int

const (
	// Nothing of interest, continue scanning.
	None Event = 0

	// Same as None, but specifically for whitespace.
	Space = (1 << iota)

	// Start events.
	ObjectStart = (1 << iota)
	KeyStart
	ArrayStart
	StringStart
	NumberStart
	BoolStart
	NullStart

	// End events.
	ObjectEnd = (1 << iota)
	KeyEnd
	ArrayEnd
	StringEnd
	NumberEnd
	BoolEnd
	NullEnd

	// Syntax error.
	Error = (1 << iota)
)
