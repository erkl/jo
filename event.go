package jo

import (
	"strings"
)

// Events signal changes in scanning state.
type Event int

const (
	// Syntax error.
	Error Event = -1

	// Nothing of interest, continue scanning.
	None = 0

	// Same as None, but specifically for whitespace.
	Space = (1 << iota)

	// Start and end events.
	ObjectStart = (1 << iota)
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

	// Start and end bitsets.
	Start = ObjectStart | KeyStart | ArrayStart | StringStart | NumberStart | BoolStart | NullStart
	End   = ObjectEnd | KeyEnd | ArrayEnd | StringEnd | NumberEnd | BoolEnd | NullEnd
)

// String returns a string representation of the Event.
func (ev Event) String() string {
	if ev == None {
		return "None"
	}
	if ev == Error {
		return "Error"
	}

	// Make sure no unknown bits are set.
	if ev&^(Space|Start|End) != 0 {
		return "INVALID"
	}

	var parts []string

	if ev&ObjectEnd != 0 {
		parts = append(parts, "ObjectEnd")
	}
	if ev&KeyEnd != 0 {
		parts = append(parts, "KeyEnd")
	}
	if ev&ArrayEnd != 0 {
		parts = append(parts, "ArrayEnd")
	}
	if ev&StringEnd != 0 {
		parts = append(parts, "StringEnd")
	}
	if ev&NumberEnd != 0 {
		parts = append(parts, "NumberEnd")
	}
	if ev&BoolEnd != 0 {
		parts = append(parts, "BoolEnd")
	}
	if ev&NullEnd != 0 {
		parts = append(parts, "NullEnd")
	}

	if ev&ObjectStart != 0 {
		parts = append(parts, "ObjectStart")
	}
	if ev&KeyStart != 0 {
		parts = append(parts, "KeyStart")
	}
	if ev&ArrayStart != 0 {
		parts = append(parts, "ArrayStart")
	}
	if ev&StringStart != 0 {
		parts = append(parts, "StringStart")
	}
	if ev&NumberStart != 0 {
		parts = append(parts, "NumberStart")
	}
	if ev&BoolStart != 0 {
		parts = append(parts, "BoolStart")
	}
	if ev&NullStart != 0 {
		parts = append(parts, "NullStart")
	}

	if ev&Space != 0 {
		parts = append(parts, "Space")
	}

	return strings.Join(parts, " | ")
}
