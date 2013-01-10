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

var validPrimitives = []parseTest{
	{`""`, []event{{1, StringStart}, {1, StringEnd}}},
	{`"abc"`, []event{{1, StringStart}, {4, StringEnd}}},
	{`"\u8bA0"`, []event{{1, StringStart}, {7, StringEnd}}},
	{`"\u12345"`, []event{{1, StringStart}, {8, StringEnd}}},
	{"\"\\b\\f\\n\\r\\t\\\\\"", []event{{1, StringStart}, {13, StringEnd}}},

	{`0`, []event{{1, NumberStart}, {0, NumberEnd}}},
	{`123`, []event{{1, NumberStart}, {2, None}, {0, NumberEnd}}},
	{`-10`, []event{{1, NumberStart}, {2, None}, {0, NumberEnd}}},
	{`0.10`, []event{{1, NumberStart}, {3, None}, {0, NumberEnd}}},
	{`12.03e+1`, []event{{1, NumberStart}, {7, None}, {0, NumberEnd}}},
	{`0.1e2`, []event{{1, NumberStart}, {4, None}, {0, NumberEnd}}},
	{`1e1`, []event{{1, NumberStart}, {2, None}, {0, NumberEnd}}},
	{`3.141569`, []event{{1, NumberStart}, {7, None}, {0, NumberEnd}}},
	{`10000000000000e-10`, []event{{1, NumberStart}, {17, None}, {0, NumberEnd}}},
	{`9223372036854775808`, []event{{1, NumberStart}, {18, None}, {0, NumberEnd}}},
	{`6E-06`, []event{{1, NumberStart}, {4, None}, {0, NumberEnd}}},
	{`1E-06`, []event{{1, NumberStart}, {4, None}, {0, NumberEnd}}},

	{`false`, []event{{1, BoolStart}, {4, BoolEnd}}},
	{`true`, []event{{1, BoolStart}, {3, BoolEnd}}},
	{`null`, []event{{1, NullStart}, {3, NullEnd}}},
}

var invalidPrimitives = []parseTest{
	{`"`, []event{{1, StringStart}, {0, SyntaxError}}},
	{`"foo`, []event{{1, StringStart}, {3, None}, {0, SyntaxError}}},
	{`'single'`, []event{{0, SyntaxError}}},
	{`"\u12g8"`, []event{{1, StringStart}, {4, SyntaxError}}},
	{`"\u"`, []event{{1, StringStart}, {2, SyntaxError}}},
	{`"you can\'t do this"`, []event{{1, StringStart}, {8, SyntaxError}}},

	{`-`, []event{{1, NumberStart}, {0, SyntaxError}}},
	{`0.`, []event{{1, NumberStart}, {1, None}, {0, SyntaxError}}},
	{`123.456.789`, []event{{1, NumberStart}, {6, SyntaxError}}},
	{`10e`, []event{{1, NumberStart}, {2, None}, {0, SyntaxError}}},
	{`10e+`, []event{{1, NumberStart}, {3, None}, {0, SyntaxError}}},
	{`10e-`, []event{{1, NumberStart}, {3, None}, {0, SyntaxError}}},
	{`0e1x`, []event{{1, NumberStart}, {2, NumberEnd}, {0, SyntaxError}}},
	{`0e+13.`, []event{{1, NumberStart}, {4, NumberEnd}, {0, SyntaxError}}},
	{`0e+-0`, []event{{1, NumberStart}, {2, SyntaxError}}},

	{`tr`, []event{{1, BoolStart}, {1, None}, {0, SyntaxError}}},
	{`truE`, []event{{1, BoolStart}, {2, SyntaxError}}},
	{`fals`, []event{{1, BoolStart}, {3, None}, {0, SyntaxError}}},
	{`fALSE`, []event{{1, BoolStart}, {0, SyntaxError}}},
	{`n`, []event{{1, NullStart}, {0, SyntaxError}}},
	{`NULL`, []event{{0, SyntaxError}}},
}

func TestParsing(t *testing.T) {
	tests := make([]parseTest, 0)
	tests = append(tests, validPrimitives...)
	tests = append(tests, invalidPrimitives...)

	for _, test := range tests {
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

			if len(input)-pos == 0 {
				log = append(log, "  .Eof()")
				where, what = 0, p.Eof()
			} else {
				line := fmt.Sprintf("  .Parse(%#q)", input[pos:])
				log = append(log, line)

				where, what = p.Parse(input[pos:])
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
