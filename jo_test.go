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

type parseTest struct {
	json   string
	events []event
}

var parseTests = []parseTest{
	{`""`, []event{{1, StringStart}, {1, StringEnd}, {0, None}}},
	{`"abc"`, []event{{1, StringStart}, {4, StringEnd}, {0, None}}},
	{`123`, []event{{1, NumberStart}, {2, None}, {0, NumberEnd}}},
	{`true`, []event{{1, BoolStart}, {3, BoolEnd}, {0, None}}},
	{`false`, []event{{1, BoolStart}, {4, BoolEnd}, {0, None}}},
	{`null`, []event{{1, NullStart}, {3, NullEnd}, {0, None}}},
}

func TestParsing(t *testing.T) {
	for _, test := range parseTests {
		input := []byte(test.json)
		pos := 0

		// instantiate a new Parser and reset the log output
		// before each test
		p := Parser{}
		log := []string{"p := Parser{}"}

		// evaluate all expected events
		for i := 0; i < len(test.events); i++ {
			event := test.events[i]

			var where int
			var what Event

			// the final event is special -- it must be triggered
			// by Parser.Eof()
			if i < len(test.events)-1 {
				line := fmt.Sprintf("  .Parse(%#q)", input[pos:])
				log = append(log, line)

				where, what = p.Parse(input[pos:])
			} else {
				log = append(log, "  .Eof()")
				where, what = 0, p.Eof()
			}

			if where != event.where || what != event.what {
				// dump the log output we've accumulated
				for _, line := range log {
					t.Logf(line)
				}

				t.Fatalf("want %s at index %d, got %s at index %d",
					reverse[event.what], event.where, reverse[what], where)
			}

			// skip the bytes consumed during the last call
			// to parser.Parse()
			pos += where
		}
	}
}
