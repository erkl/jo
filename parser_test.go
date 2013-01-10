package jo

import (
	"fmt"
	"testing"
)

type event struct {
	where int
	what  Event
}

type parseTest struct {
	json   string
	events []event
}

var legalLiterals = []parseTest{
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

var illegalLiterals = []parseTest{
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

var legalObjects = []parseTest{
	{
		`{}`,
		[]event{
			{1, ObjectStart},
			{1, ObjectEnd},
			{0, None},
		},
	},
	{
		`{"foo":"bar"}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{4, KeyEnd},
			{2, StringStart},
			{4, StringEnd},
			{1, ObjectEnd},
			{0, None},
		},
	},
	{
		`{"1":1,"2":2}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{2, KeyEnd},
			{2, NumberStart},
			{0, NumberEnd},
			{2, KeyStart},
			{2, KeyEnd},
			{2, NumberStart},
			{0, NumberEnd},
			{1, ObjectEnd},
			{0, None},
		},
	},
	{
		`{"\u1234\t\n\b\u8BbF\"":0}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{21, KeyEnd},
			{2, NumberStart},
			{0, NumberEnd},
			{1, ObjectEnd},
			{0, None},
		},
	},
	{
		`{"{":"}"}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{2, KeyEnd},
			{2, StringStart},
			{2, StringEnd},
			{1, ObjectEnd},
			{0, None},
		},
	},
	{
		`{"foo":{"bar":{"baz":{}}}}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{4, KeyEnd},
			{2, ObjectStart},
			{1, KeyStart},
			{4, KeyEnd},
			{2, ObjectStart},
			{1, KeyStart},
			{4, KeyEnd},
			{2, ObjectStart},
			{1, ObjectEnd},
			{1, ObjectEnd},
			{1, ObjectEnd},
			{1, ObjectEnd},
			{0, None},
		},
	},
	{
		`{"true":true,"false":false,"null":null}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{5, KeyEnd},
			{2, BoolStart},
			{3, BoolEnd},
			{2, KeyStart},
			{6, KeyEnd},
			{2, BoolStart},
			{4, BoolEnd},
			{2, KeyStart},
			{5, KeyEnd},
			{2, NullStart},
			{3, NullEnd},
			{1, ObjectEnd},
			{0, None},
		},
	},
}

var illegalObjects = []parseTest{
	{
		`{0:1}`,
		[]event{
			{1, ObjectStart},
			{0, SyntaxError},
		},
	},
	{
		`{"foo":"bar"`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{4, KeyEnd},
			{2, StringStart},
			{4, StringEnd},
			{0, SyntaxError},
		},
	},
	{
		`{{`,
		[]event{
			{1, ObjectStart},
			{0, SyntaxError},
		},
	},
	{
		`{"a":1,}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{2, KeyEnd},
			{2, NumberStart},
			{0, NumberEnd},
			{1, SyntaxError},
		},
	},
	{
		`{"a":1,,`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{2, KeyEnd},
			{2, NumberStart},
			{0, NumberEnd},
			{1, SyntaxError},
		},
	},
	{
		`{"a"}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{2, KeyEnd},
			{0, SyntaxError},
		},
	},
	{
		`{"a":"1}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{2, KeyEnd},
			{2, StringStart},
			{2, None},
			{0, SyntaxError},
		},
	},
	{
		`{"a":1"b":2}`,
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{2, KeyEnd},
			{2, NumberStart},
			{0, NumberEnd},
			{0, SyntaxError},
		},
	},
}

var legalArrays = []parseTest{
	{
		`[]`,
		[]event{
			{1, ArrayStart},
			{1, ArrayEnd},
			{0, None},
		},
	},
	{
		`[1]`,
		[]event{
			{1, ArrayStart},
			{1, NumberStart},
			{0, NumberEnd},
			{1, ArrayEnd},
		},
	},
	{
		`[1,2,3]`,
		[]event{
			{1, ArrayStart},
			{1, NumberStart},
			{0, NumberEnd},
			{2, NumberStart},
			{0, NumberEnd},
			{2, NumberStart},
			{0, NumberEnd},
			{1, ArrayEnd},
		},
	},
	{
		`["dude","what"]`,
		[]event{
			{1, ArrayStart},
			{1, StringStart},
			{5, StringEnd},
			{2, StringStart},
			{5, StringEnd},
			{1, ArrayEnd},
			{0, None},
		},
	},
	{
		`[[[],[]],[]]`,
		[]event{
			{1, ArrayStart},
			{1, ArrayStart},
			{1, ArrayStart},
			{1, ArrayEnd},
			{2, ArrayStart},
			{1, ArrayEnd},
			{1, ArrayEnd},
			{2, ArrayStart},
			{1, ArrayEnd},
			{1, ArrayEnd},
			{0, None},
		},
	},
	{
		`["][\"]"]`,
		[]event{
			{1, ArrayStart},
			{1, StringStart},
			{6, StringEnd},
			{1, ArrayEnd},
			{0, None},
		},
	},
}

var illegalArrays = []parseTest{
	{
		`[`,
		[]event{
			{1, ArrayStart},
			{0, SyntaxError},
		},
	},
	{
		`[,]`,
		[]event{
			{1, ArrayStart},
			{0, SyntaxError},
		},
	},
	{
		`[10,]`,
		[]event{
			{1, ArrayStart},
			{1, NumberStart},
			{1, NumberEnd},
			{1, SyntaxError},
		},
	},
	{
		`[}`,
		[]event{
			{1, ArrayStart},
			{0, SyntaxError},
		},
	},
}

var legalWhitespace = []parseTest{
	{` 17`, []event{{2, NumberStart}, {1, None}, {0, NumberEnd}}},
	{`38 `, []event{{1, NumberStart}, {1, NumberEnd}, {1, None}, {0, None}}},
	{`  " what ? "  `, []event{{3, StringStart}, {9, StringEnd}, {2, None}, {0, None}}},
	{"\nnull", []event{{2, NullStart}, {3, NullEnd}, {0, None}}},
	{"\n\r\t true \r\n\t", []event{{5, BoolStart}, {3, BoolEnd}, {4, None}, {0, None}}},

	{
		" { \"foo\": \t\"bar\" } ",
		[]event{
			{2, ObjectStart},
			{2, KeyStart},
			{4, KeyEnd},
			{4, StringStart},
			{4, StringEnd},
			{2, ObjectEnd},
			{1, None},
		},
	},
	{
		"\t[ 1 , 2\r, 3\n]",
		[]event{
			{2, ArrayStart},
			{2, NumberStart},
			{0, NumberEnd},
			{4, NumberStart},
			{0, NumberEnd},
			{4, NumberStart},
			{0, NumberEnd},
			{2, ArrayEnd},
			{0, None},
		},
	},
	{
		" \n{ \t\"foo\" : [ \"bar\", null ], \"what\\n\"\t: 10.3e1 } ",
		[]event{
			{3, ObjectStart},
			{3, KeyStart},
			{4, KeyEnd},
			{4, ArrayStart},
			{2, StringStart},
			{4, StringEnd},
			{3, NullStart},
			{3, NullEnd},
			{2, ArrayEnd},
			{3, KeyStart},
			{7, KeyEnd},
			{4, NumberStart},
			{5, NumberEnd},
			{2, ObjectEnd},
			{1, None},
			{0, None},
		},
	},
}

func TestParsing(t *testing.T) {
	tests := make([]parseTest, 0)
	tests = append(tests, legalLiterals...)
	tests = append(tests, legalObjects...)
	tests = append(tests, legalArrays...)
	tests = append(tests, illegalLiterals...)
	tests = append(tests, illegalObjects...)
	tests = append(tests, illegalArrays...)
	tests = append(tests, legalWhitespace...)

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
					event.what, event.where, what, where)
			}

			// skip the bytes consumed during the last call
			// to parser.Parse()
			pos += where
		}
	}
}
