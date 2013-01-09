package jo

import (
	"fmt"
	"testing"
)

// reverse lookups of event names
var reverse = map[Event]string{
	None:        "None",
	SyntaxError: "SyntaxError",
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
}

type event struct {
	where int
	what  Event
}

type test struct {
	json   string
	events []event
}

var tests = []test{
	{`""`, []event{{1, StringStart}, {1, StringEnd}, {0, None}}},
	{`"abc"`, []event{{1, StringStart}, {4, StringEnd}, {0, None}}},
	{`123`, []event{{1, NumberStart}, {2, None}, {0, NumberEnd}}},
	{`true`, []event{{1, BoolStart}, {3, BoolEnd}, {0, None}}},
	{`false`, []event{{1, BoolStart}, {4, BoolEnd}, {0, None}}},
	{`null`, []event{{1, NullStart}, {3, NullEnd}, {0, None}}},
}

func TestEverything(t *testing.T) {
	for _, test := range tests {
		p := Parser{}

		input := []byte(test.json)
		pos := 0

		// only save log output for this particular test
		log := []string{"p := Parser{}"}

		// run through all expected events
		for i := 0; i < len(test.events); i++ {
			event := test.events[i]

			var n int
			var e Event

			// the last event must come from `p.Eof()
			if i == len(test.events)-1 {
				log = append(log, "  .Eof()")
				n, e = 0, p.Eof()
			} else {
				log = append(log, fmt.Sprintf("  .Parse(%#q)", input[pos:]))
				n, e = p.Parse(input[pos:])
			}

			// are we happy with the outcome?
			if n != event.where || e != event.what {
				for _, line := range log {
					t.Logf(line)
				}

				t.Fatalf("want %s at index %d, got %s at index %d",
					reverse[event.what], event.where, reverse[e], n)
			}

			// skip the consumed bytes
			pos += n
		}
	}
}
