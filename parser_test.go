package jo

import (
	"fmt"
	"testing"
)

func ExampleParser() {
	input := []byte(`{ "foo": 10 }`)
	pos := 0

	p := Parser{}

	for pos < len(input) {
		n, ev := p.Parse(input[pos:])
		pos += n

		fmt.Printf("at %2d -> %s\n", pos, ev)
	}

	fmt.Printf("at %2d -> %s\n", pos, p.End())

	// Output:
	// at  1 -> ObjectStart
	// at  3 -> KeyStart
	// at  7 -> KeyEnd
	// at 10 -> NumberStart
	// at 11 -> NumberEnd
	// at 13 -> ObjectEnd
	// at 13 -> Done
}

type event struct {
	where int
	what  Event
}

var parserTests = []struct {
	json   string
	events []event
}{
	{
		`""`,
		[]event{
			{1, StringStart},
			{1, StringEnd},
			{0, Done},
		},
	},
	{
		`"abc"`,
		[]event{
			{1, StringStart},
			{4, StringEnd},
			{0, Done},
		},
	},
	{
		`"\u8bA0"`,
		[]event{
			{1, StringStart},
			{7, StringEnd},
			{0, Done},
		},
	},
	{
		`"\u12345"`,
		[]event{
			{1, StringStart},
			{8, StringEnd},
			{0, Done},
		},
	},
	{
		"\"\\b\\f\\n\\r\\t\\\\\"",
		[]event{
			{1, StringStart},
			{13, StringEnd},
			{0, Done},
		},
	},
	{
		`0`,
		[]event{
			{1, NumberStart},
			{0, NumberEnd},
		},
	},
	{
		`123`,
		[]event{
			{1, NumberStart},
			{2, Continue},
			{0, NumberEnd},
		},
	},
	{
		`-10`,
		[]event{
			{1, NumberStart},
			{2, Continue},
			{0, NumberEnd},
		},
	},
	{
		`0.10`,
		[]event{
			{1, NumberStart},
			{3, Continue},
			{0, NumberEnd},
		},
	},
	{
		`12.03e+1`,
		[]event{
			{1, NumberStart},
			{7, Continue},
			{0, NumberEnd},
		},
	},
	{
		`0.1e2`,
		[]event{
			{1, NumberStart},
			{4, Continue},
			{0, NumberEnd},
		},
	},
	{
		`1e1`,
		[]event{
			{1, NumberStart},
			{2, Continue},
			{0, NumberEnd},
		},
	},
	{
		`3.141569`,
		[]event{
			{1, NumberStart},
			{7, Continue},
			{0, NumberEnd},
		},
	},
	{
		`10000000000000e-10`,
		[]event{
			{1, NumberStart},
			{17, Continue},
			{0, NumberEnd},
		},
	},
	{
		`9223372036854775808`,
		[]event{
			{1, NumberStart},
			{18, Continue},
			{0, NumberEnd},
		},
	},
	{
		`6E-06`,
		[]event{
			{1, NumberStart},
			{4, Continue},
			{0, NumberEnd},
		},
	},
	{
		`1E-06`,
		[]event{
			{1, NumberStart},
			{4, Continue},
			{0, NumberEnd},
		},
	},
	{
		`false`,
		[]event{
			{1, BoolStart},
			{4, BoolEnd},
			{0, Done},
		},
	},
	{
		`true`,
		[]event{
			{1, BoolStart},
			{3, BoolEnd},
			{0, Done},
		},
	},
	{
		`null`,
		[]event{
			{1, NullStart},
			{3, NullEnd},
			{0, Done},
		},
	},
	{
		`"`,
		[]event{
			{1, StringStart},
			{0, SyntaxError},
		},
	},
	{
		`"foo`,
		[]event{
			{1, StringStart},
			{3, Continue},
			{0, SyntaxError},
		},
	},
	{
		`'single'`,
		[]event{
			{0, SyntaxError},
		},
	},
	{
		`"\u12g8"`,
		[]event{
			{1, StringStart},
			{4, SyntaxError},
		},
	},
	{
		`"\u"`,
		[]event{
			{1, StringStart},
			{2, SyntaxError},
		},
	},
	{
		`"you can\'t do this"`,
		[]event{
			{1, StringStart},
			{8, SyntaxError},
		},
	},
	{
		`-`,
		[]event{
			{1, NumberStart},
			{0, SyntaxError},
		},
	},
	{
		`0.`,
		[]event{
			{1, NumberStart},
			{1, Continue},
			{0, SyntaxError},
		},
	},
	{
		`123.456.789`,
		[]event{
			{1, NumberStart},
			{6, NumberEnd},
			{0, SyntaxError},
		},
	},
	{
		`10e`,
		[]event{
			{1, NumberStart},
			{2, Continue},
			{0, SyntaxError},
		},
	},
	{
		`10e+`,
		[]event{
			{1, NumberStart},
			{3, Continue},
			{0, SyntaxError},
		},
	},
	{
		`10e-`,
		[]event{
			{1, NumberStart},
			{3, Continue},
			{0, SyntaxError},
		},
	},
	{
		`0e1x`,
		[]event{
			{1, NumberStart},
			{2, NumberEnd},
			{0, SyntaxError},
		},
	},
	{
		`0e+13.`,
		[]event{
			{1, NumberStart},
			{4, NumberEnd},
			{0, SyntaxError},
		},
	},
	{
		`0e+-0`,
		[]event{
			{1, NumberStart},
			{2, SyntaxError},
		},
	},
	{
		`tr`,
		[]event{
			{1, BoolStart},
			{1, Continue},
			{0, SyntaxError},
		},
	},
	{
		`truE`,
		[]event{
			{1, BoolStart},
			{2, SyntaxError},
		},
	},
	{
		`fals`,
		[]event{
			{1, BoolStart},
			{3, Continue},
			{0, SyntaxError},
		},
	},
	{
		`fALSE`,
		[]event{
			{1, BoolStart},
			{0, SyntaxError},
		},
	},
	{
		`n`,
		[]event{
			{1, NullStart},
			{0, SyntaxError},
		},
	},
	{
		`NULL`,
		[]event{
			{0, SyntaxError},
		},
	},
	{
		`{}`,
		[]event{
			{1, ObjectStart},
			{1, ObjectEnd},
			{0, Done},
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
			{0, Done},
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
			{0, Done},
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
			{0, Done},
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
			{0, Done},
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
			{0, Done},
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
			{0, Done},
		},
	},
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
			{2, Continue},
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

	{
		`[]`,
		[]event{
			{1, ArrayStart},
			{1, ArrayEnd},
			{0, Done},
		},
	},
	{
		`[1]`,
		[]event{
			{1, ArrayStart},
			{1, NumberStart},
			{0, NumberEnd},
			{1, ArrayEnd},
			{0, Done},
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
			{0, Done},
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
			{0, Done},
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
			{0, Done},
		},
	},
	{
		`[10.3]`,
		[]event{
			{1, ArrayStart},
			{1, NumberStart},
			{3, NumberEnd},
			{1, ArrayEnd},
			{0, Done},
		},
	},
	{
		`["][\"]"]`,
		[]event{
			{1, ArrayStart},
			{1, StringStart},
			{6, StringEnd},
			{1, ArrayEnd},
			{0, Done},
		},
	},
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
	{
		` 17`,
		[]event{
			{2, NumberStart},
			{1, Continue},
			{0, NumberEnd},
		},
	},
	{`38 `,
		[]event{
			{1, NumberStart},
			{1, NumberEnd},
			{1, Continue},
			{0, Done},
		},
	},
	{
		`  " what ? "  `,
		[]event{
			{3, StringStart},
			{9, StringEnd},
			{2, Continue},
			{0, Done},
		},
	},
	{
		"\nnull",
		[]event{
			{2, NullStart},
			{3, NullEnd},
			{0, Done},
		},
	},
	{
		"\n\r\t true \r\n\t",
		[]event{
			{5, BoolStart},
			{3, BoolEnd},
			{4, Continue},
			{0, Done},
		},
	},
	{
		" { \"foo\": \t\"bar\" } ",
		[]event{
			{2, ObjectStart},
			{2, KeyStart},
			{4, KeyEnd},
			{4, StringStart},
			{4, StringEnd},
			{2, ObjectEnd},
			{1, Continue},
			{0, Done},
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
			{0, Done},
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
			{1, Continue},
			{0, Done},
		},
	},
}

// Tests basic JSON parsing.
func TestParsing(t *testing.T) {
	for _, test := range parserTests {
		h := helper{t: t, in: []byte(test.json)}

		for i := 0; ; i++ {
			want := test.events[i]
			n, ev, eof := h.next()

			if n != want.where || ev != want.what {
				h.fail("wanted %s at offset %d", want.what, want.where)
				break
			}

			if eof {
				break
			}
		}
	}
}

// Tests the parser's Depth() method.
func TestDepth(t *testing.T) {
	for _, test := range parserTests {
		h := helper{t: t, in: []byte(test.json)}
		d := 0

		for {
			_, ev, eof := h.next()

			switch ev {
			case ObjectStart, ArrayStart:
				d++
			case ObjectEnd, ArrayEnd:
				d--
			}

			if got := h.Depth(); got != d {
				h.fail("depth should be %d, was %d", d, got)
				break
			}

			if eof {
				break
			}
		}
	}
}

type escape struct {
	when, depth int
}

var escapeTests = []struct {
	json    string
	escapes []escape
	events  []event
}{
	{
		`[]`,
		[]escape{{1, 1}},
		[]event{
			{1, ArrayStart},
			// p.Escape(1)
			{1, ArrayEnd},
			{0, Done},
		},
	},
	{
		`{"foo":"bar"}`,
		[]escape{{1, 1}},
		[]event{
			{1, ObjectStart},
			// p.Escape(1)
			{12, ObjectEnd},
			{0, Done},
		},
	},
	{
		`[[{},2,3]]`,
		[]escape{{2, 1}},
		[]event{
			{1, ArrayStart},
			{1, ArrayStart},
			// p.Escape(1)
			{7, ArrayEnd},
			{1, ArrayEnd},
			{0, Done},
		},
	},
	{
		`[[[[[],1],2],3],4]`,
		[]escape{{5, 5}},
		[]event{
			{1, ArrayStart},
			{1, ArrayStart},
			{1, ArrayStart},
			{1, ArrayStart},
			{1, ArrayStart},
			// p.Escape(5)
			{1, ArrayEnd},
			{3, ArrayEnd},
			{3, ArrayEnd},
			{3, ArrayEnd},
			{3, ArrayEnd},
			{0, Done},
		},
	},
	{
		`[{"foo":"bar", "num": 1}, ["deeper", ["and deeper", 1, 2]]]`,
		[]escape{{2, 2}},
		[]event{
			{1, ArrayStart},
			{1, ObjectStart},
			// p.Escape(2)
			{22, ObjectEnd},
			{35, ArrayEnd},
			{0, Done},
		},
	},
	{
		`{"key":123}`,
		[]escape{{3, 1}},
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{4, KeyEnd},
			// p.Escape(1)
			{5, ObjectEnd},
			{0, Done},
		},
	},
	{
		`null`,
		[]escape{{1, 1}},
		[]event{
			{1, NullStart},
			// p.Escape(1)
			{3, Continue},
			{0, Done},
		},
	},
	{
		`[{"foo":"bar"}]`,
		[]escape{{5, 2}},
		[]event{
			{1, ArrayStart},
			{1, ObjectStart},
			{1, KeyStart},
			{4, KeyEnd},
			{2, StringStart},
			// p.Escape(2)
			{5, ObjectEnd},
			{1, ArrayEnd},
			{0, Done},
		},
	},
	{
		`[{"foo":"bar"}]`,
		[]escape{{5, 2}},
		[]event{
			{1, ArrayStart},
			{1, ObjectStart},
			{1, KeyStart},
			{4, KeyEnd},
			{2, StringStart},
			// p.Escape(2)
			{5, ObjectEnd},
			{1, ArrayEnd},
			{0, Done},
		},
	},
	{
		`[["skip"],[]]`,
		[]escape{{3, 1}, {4, 1}},
		[]event{
			{1, ArrayStart},
			{1, ArrayStart},
			{1, StringStart},
			// p.Escape(1)
			{6, ArrayEnd},
			// p.Escape(1)
			{4, ArrayEnd},
			{0, Done},
		},
	},
	{
		`{"a":"wub","wrong"}`,
		[]escape{{6, 1}},
		[]event{
			{1, ObjectStart},
			{1, KeyStart},
			{2, KeyEnd},
			{2, StringStart},
			{4, StringEnd},
			{2, KeyStart},
			// p.Escape(1)
			{6, SyntaxError},
		},
	},
	{
		`[[],[]]`,
		[]escape{{2, 1}},
		[]event{
			{1, ArrayStart},
			{1, ArrayStart},
			// p.Escape(1)
			{1, ArrayEnd},
			{2, ArrayStart},
			{1, ArrayEnd},
			{1, ArrayEnd},
			{0, Done},
		},
	},
}

// Tests the parser's Escape() method.
func TestEscape(t *testing.T) {
	for _, test := range escapeTests {
		h := helper{t: t, in: []byte(test.json)}
		n := 0

		for i := 0; ; i++ {
			// see if the test case wants us to invoke Escape() here
			if n < len(test.escapes) {
				esc := test.escapes[n]

				if esc.when == i {
					h.logf("p.Escape(%d)", esc.depth)
					h.Escape(esc.depth)

					n++
				}
			}

			want := test.events[i]
			n, ev, eof := h.next()

			if n != want.where || ev != want.what {
				h.fail("wanted %s at offset %d", want.what, want.where)
				break
			}

			if eof {
				break
			}
		}
	}
}

// Test case helper which wraps around a Parser struct.
type helper struct {
	Parser
	t   *testing.T
	in  []byte
	log []string
}

// Feeds what's left of the JSON input through the parser, then returns the
// outcome. Calls Parser.End() automatically when all input has been parsed.
func (h *helper) next() (int, Event, bool) {
	total := 0

	for len(h.in) > 0 {
		n, ev := h.Parse(h.in[:1])
		h.logf("p.Parse(%#q) -> %d, %s", h.in, n, ev)

		h.in = h.in[n:]
		total += n

		if ev == Continue && len(h.in) > 0 {
			continue
		}

		return total, ev, ev == SyntaxError
	}

	ev := h.End()
	h.logf("p.End() -> %s", ev)

	return 0, ev, true
}

// Logs anything specific to this test case.
func (h *helper) logf(format string, args ...interface{}) {
	h.log = append(h.log, fmt.Sprintf(format, args...))
}

// Reports an error in the current test case.
func (h *helper) fail(format string, args ...interface{}) {
	h.t.Log("p := Parser{}")
	for _, s := range h.log {
		h.t.Log(s)
	}
	h.t.Errorf(format, args...)
}
