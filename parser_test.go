package jo

import (
	"fmt"
	"testing"
)

func ExampleParser() {
	input := []byte(`{ "foo": 10 }`)
	parser := Parser{}
	parsed := 0

	for parsed < len(input) {
		n, event, _ := parser.Parse(input[parsed:])
		parsed += n

		fmt.Printf("at %2d -> %s\n", parsed, event)
	}

	event, _ := parser.End()
	fmt.Printf("at %2d -> %s\n", parsed, event)

	// Output:
	// at  1 -> ObjectStart
	// at  3 -> KeyStart
	// at  7 -> KeyEnd
	// at 10 -> NumberStart
	// at 11 -> NumberEnd
	// at 13 -> ObjectEnd
	// at 13 -> Done
}

// represents a step in a parser test case
type step func(*Parser, []byte) (int, bool, string)

// returns a function which checks an event returned by Parser.Parser()
func parse(offset int, want Event) step {
	return func(p *Parser, in []byte) (int, bool, string) {
		n, actual, err := p.Parse(in)
		log := fmt.Sprintf(".Parse(%#q) -> %d, %s, %#v", in, n, actual, err)

		if n != offset || actual != want {
			log = log + fmt.Sprintf(" (want %d, %s, <nil>)", offset, want)
			return n, false, log
		}

		return n, true, log
	}
}

// returns a function which checks the event returned by Parser.End()
func end(want Event) step {
	return func(p *Parser, in []byte) (int, bool, string) {
		actual, err := p.End()
		log := fmt.Sprintf(".End() -> %s, %#v", actual, err)

		if actual != want {
			log = log + fmt.Sprintf(" (want %s, <nil>)", actual)
			return 0, false, log
		}

		return 0, true, log
	}
}

// returns a function which invokes and checks Parser.Depth()
func depth(want int) step {
	return func(p *Parser, in []byte) (int, bool, string) {
		actual := p.Depth()
		log := fmt.Sprintf(".Depth() -> %d", actual)

		if actual != want {
			log = log + fmt.Sprintf(" (want %d)", want)
			return 0, false, log
		}

		return 0, true, log
	}
}

// returns a function which invokes Parser.Skip()
func skip(dead, empty int) step {
	return func(p *Parser, in []byte) (int, bool, string) {
		p.Skip(dead, empty)
		log := fmt.Sprintf(".Skip(%d, %d)", dead, empty)

		return 0, true, log
	}
}

// a parser test case
type parserTest struct {
	json  string
	steps []step
}

// basic parsing tests
var basicTests = []parserTest{
	{
		`""`,
		[]step{
			parse(1, StringStart),
			parse(1, StringEnd),
			end(Done),
		},
	},
	{
		`"abc"`,
		[]step{
			parse(1, StringStart),
			parse(4, StringEnd),
			end(Done),
		},
	},
	{
		`"\u8bA0"`,
		[]step{
			parse(1, StringStart),
			parse(7, StringEnd),
			end(Done),
		},
	},
	{
		`"\u12345"`,
		[]step{
			parse(1, StringStart),
			parse(8, StringEnd),
			end(Done),
		},
	},
	{
		"\"\\b\\f\\n\\r\\t\\\\\"",
		[]step{
			parse(1, StringStart),
			parse(13, StringEnd),
			end(Done),
		},
	},
	{
		`0`,
		[]step{
			parse(1, NumberStart),
			end(NumberEnd),
		},
	},
	{
		`123`,
		[]step{
			parse(1, NumberStart),
			parse(2, Continue),
			end(NumberEnd),
		},
	},
	{
		`-10`,
		[]step{
			parse(1, NumberStart),
			parse(2, Continue),
			end(NumberEnd),
		},
	},
	{
		`0.10`,
		[]step{
			parse(1, NumberStart),
			parse(3, Continue),
			end(NumberEnd),
		},
	},
	{
		`12.03e+1`,
		[]step{
			parse(1, NumberStart),
			parse(7, Continue),
			end(NumberEnd),
		},
	},
	{
		`0.1e2`,
		[]step{
			parse(1, NumberStart),
			parse(4, Continue),
			end(NumberEnd),
		},
	},
	{
		`1e1`,
		[]step{
			parse(1, NumberStart),
			parse(2, Continue),
			end(NumberEnd),
		},
	},
	{
		`3.141569`,
		[]step{
			parse(1, NumberStart),
			parse(7, Continue),
			end(NumberEnd),
		},
	},
	{
		`10000000000000e-10`,
		[]step{
			parse(1, NumberStart),
			parse(17, Continue),
			end(NumberEnd),
		},
	},
	{
		`9223372036854775808`,
		[]step{
			parse(1, NumberStart),
			parse(18, Continue),
			end(NumberEnd),
		},
	},
	{
		`6E-06`,
		[]step{
			parse(1, NumberStart),
			parse(4, Continue),
			end(NumberEnd),
		},
	},
	{
		`1E-06`,
		[]step{
			parse(1, NumberStart),
			parse(4, Continue),
			end(NumberEnd),
		},
	},
	{
		`false`,
		[]step{
			parse(1, BoolStart),
			parse(4, BoolEnd),
			end(Done),
		},
	},
	{
		`true`,
		[]step{
			parse(1, BoolStart),
			parse(3, BoolEnd),
			end(Done),
		},
	},
	{
		`null`,
		[]step{
			parse(1, NullStart),
			parse(3, NullEnd),
			end(Done),
		},
	},
	{
		`"`,
		[]step{
			parse(1, StringStart),
			end(SyntaxError),
		},
	},
	{
		`"foo`,
		[]step{
			parse(1, StringStart),
			parse(3, Continue),
			end(SyntaxError),
		},
	},
	{
		`'single'`,
		[]step{
			parse(0, SyntaxError),
		},
	},
	{
		`"\u12g8"`,
		[]step{
			parse(1, StringStart),
			parse(4, SyntaxError),
		},
	},
	{
		`"\u"`,
		[]step{
			parse(1, StringStart),
			parse(2, SyntaxError),
		},
	},
	{
		`"you can\'t do this"`,
		[]step{
			parse(1, StringStart),
			parse(8, SyntaxError),
		},
	},
	{
		`-`,
		[]step{
			parse(1, NumberStart),
			end(SyntaxError),
		},
	},
	{
		`0.`,
		[]step{
			parse(1, NumberStart),
			parse(1, Continue),
			end(SyntaxError),
		},
	},
	{
		`123.456.789`,
		[]step{
			parse(1, NumberStart),
			parse(6, NumberEnd),
			parse(0, SyntaxError),
		},
	},
	{
		`10e`,
		[]step{
			parse(1, NumberStart),
			parse(2, Continue),
			end(SyntaxError),
		},
	},
	{
		`10e+`,
		[]step{
			parse(1, NumberStart),
			parse(3, Continue),
			end(SyntaxError),
		},
	},
	{
		`10e-`,
		[]step{
			parse(1, NumberStart),
			parse(3, Continue),
			end(SyntaxError),
		},
	},
	{
		`0e1x`,
		[]step{
			parse(1, NumberStart),
			parse(2, NumberEnd),
			parse(0, SyntaxError),
		},
	},
	{
		`0e+13.`,
		[]step{
			parse(1, NumberStart),
			parse(4, NumberEnd),
			parse(0, SyntaxError),
		},
	},
	{
		`0e+-0`,
		[]step{
			parse(1, NumberStart),
			parse(2, SyntaxError),
		},
	},
	{
		`tr`,
		[]step{
			parse(1, BoolStart),
			parse(1, Continue),
			end(SyntaxError),
		},
	},
	{
		`truE`,
		[]step{
			parse(1, BoolStart),
			parse(2, SyntaxError),
		},
	},
	{
		`fals`,
		[]step{
			parse(1, BoolStart),
			parse(3, Continue),
			end(SyntaxError),
		},
	},
	{
		`fALSE`,
		[]step{
			parse(1, BoolStart),
			parse(0, SyntaxError),
		},
	},
	{
		`n`,
		[]step{
			parse(1, NullStart),
			end(SyntaxError),
		},
	},
	{
		`NULL`,
		[]step{
			parse(0, SyntaxError),
		},
	},
	{
		`{}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"foo":"bar"}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(4, KeyEnd),
			parse(2, StringStart),
			parse(4, StringEnd),
			parse(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"1":1,"2":2}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(2, KeyEnd),
			parse(2, NumberStart),
			parse(0, NumberEnd),
			parse(2, KeyStart),
			parse(2, KeyEnd),
			parse(2, NumberStart),
			parse(0, NumberEnd),
			parse(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"\u1234\t\n\b\u8BbF\"":0}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(21, KeyEnd),
			parse(2, NumberStart),
			parse(0, NumberEnd),
			parse(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"{":"}"}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(2, KeyEnd),
			parse(2, StringStart),
			parse(2, StringEnd),
			parse(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"foo":{"bar":{"baz":{}}}}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(4, KeyEnd),
			parse(2, ObjectStart),
			parse(1, KeyStart),
			parse(4, KeyEnd),
			parse(2, ObjectStart),
			parse(1, KeyStart),
			parse(4, KeyEnd),
			parse(2, ObjectStart),
			parse(1, ObjectEnd),
			parse(1, ObjectEnd),
			parse(1, ObjectEnd),
			parse(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"true":true,"false":false,"null":null}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(5, KeyEnd),
			parse(2, BoolStart),
			parse(3, BoolEnd),
			parse(2, KeyStart),
			parse(6, KeyEnd),
			parse(2, BoolStart),
			parse(4, BoolEnd),
			parse(2, KeyStart),
			parse(5, KeyEnd),
			parse(2, NullStart),
			parse(3, NullEnd),
			parse(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{0:1}`,
		[]step{
			parse(1, ObjectStart),
			parse(0, SyntaxError),
		},
	},
	{
		`{"foo":"bar"`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(4, KeyEnd),
			parse(2, StringStart),
			parse(4, StringEnd),
			end(SyntaxError),
		},
	},
	{
		`{{`,
		[]step{
			parse(1, ObjectStart),
			parse(0, SyntaxError),
		},
	},
	{
		`{"a":1,}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(2, KeyEnd),
			parse(2, NumberStart),
			parse(0, NumberEnd),
			parse(1, SyntaxError),
		},
	},
	{
		`{"a":1,,`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(2, KeyEnd),
			parse(2, NumberStart),
			parse(0, NumberEnd),
			parse(1, SyntaxError),
		},
	},
	{
		`{"a"}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(2, KeyEnd),
			parse(0, SyntaxError),
		},
	},
	{
		`{"a":"1}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(2, KeyEnd),
			parse(2, StringStart),
			parse(2, Continue),
			end(SyntaxError),
		},
	},
	{
		`{"a":1"b":2}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(2, KeyEnd),
			parse(2, NumberStart),
			parse(0, NumberEnd),
			parse(0, SyntaxError),
		},
	},

	{
		`[]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[1]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, NumberStart),
			parse(0, NumberEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[1,2,3]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, NumberStart),
			parse(0, NumberEnd),
			parse(2, NumberStart),
			parse(0, NumberEnd),
			parse(2, NumberStart),
			parse(0, NumberEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`["dude","what"]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, StringStart),
			parse(5, StringEnd),
			parse(2, StringStart),
			parse(5, StringEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[[[],[]],[]]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			parse(1, ArrayEnd),
			parse(2, ArrayStart),
			parse(1, ArrayEnd),
			parse(1, ArrayEnd),
			parse(2, ArrayStart),
			parse(1, ArrayEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[10.3]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, NumberStart),
			parse(3, NumberEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`["][\"]"]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, StringStart),
			parse(6, StringEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[`,
		[]step{
			parse(1, ArrayStart),
			end(SyntaxError),
		},
	},
	{
		`[,]`,
		[]step{
			parse(1, ArrayStart),
			parse(0, SyntaxError),
		},
	},
	{
		`[10,]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, NumberStart),
			parse(1, NumberEnd),
			parse(1, SyntaxError),
		},
	},
	{
		`[}`,
		[]step{
			parse(1, ArrayStart),
			parse(0, SyntaxError),
		},
	},
	{
		` 17`,
		[]step{
			parse(2, NumberStart),
			parse(1, Continue),
			end(NumberEnd),
		},
	},
	{`38 `,
		[]step{
			parse(1, NumberStart),
			parse(1, NumberEnd),
			parse(1, Continue),
			end(Done),
		},
	},
	{
		`  " what ? "  `,
		[]step{
			parse(3, StringStart),
			parse(9, StringEnd),
			parse(2, Continue),
			end(Done),
		},
	},
	{
		"\nnull",
		[]step{
			parse(2, NullStart),
			parse(3, NullEnd),
			end(Done),
		},
	},
	{
		"\n\r\t true \r\n\t",
		[]step{
			parse(5, BoolStart),
			parse(3, BoolEnd),
			parse(4, Continue),
			end(Done),
		},
	},
	{
		" { \"foo\": \t\"bar\" } ",
		[]step{
			parse(2, ObjectStart),
			parse(2, KeyStart),
			parse(4, KeyEnd),
			parse(4, StringStart),
			parse(4, StringEnd),
			parse(2, ObjectEnd),
			parse(1, Continue),
			end(Done),
		},
	},
	{
		"\t[ 1 , 2\r, 3\n]",
		[]step{
			parse(2, ArrayStart),
			parse(2, NumberStart),
			parse(0, NumberEnd),
			parse(4, NumberStart),
			parse(0, NumberEnd),
			parse(4, NumberStart),
			parse(0, NumberEnd),
			parse(2, ArrayEnd),
			end(Done),
		},
	},
	{
		" \n{ \t\"foo\" : [ \"bar\", null ], \"what\\n\"\t: 10.3e1 } ",
		[]step{
			parse(3, ObjectStart),
			parse(3, KeyStart),
			parse(4, KeyEnd),
			parse(4, ArrayStart),
			parse(2, StringStart),
			parse(4, StringEnd),
			parse(3, NullStart),
			parse(3, NullEnd),
			parse(2, ArrayEnd),
			parse(3, KeyStart),
			parse(7, KeyEnd),
			parse(4, NumberStart),
			parse(5, NumberEnd),
			parse(2, ObjectEnd),
			parse(1, Continue),
			end(Done),
		},
	},
}

// tests involving Parser.Depth()
var depthTests = []parserTest{
	{
		`"hello"`,
		[]step{
			depth(0),
			parse(1, StringStart),
			depth(1),
			parse(6, StringEnd),
			end(Done),
		},
	},
	{
		`{"what":false}`,
		[]step{
			parse(1, ObjectStart),
			depth(1),
			parse(1, KeyStart),
			depth(2),
			parse(5, KeyEnd),
			depth(2),
			parse(2, BoolStart),
			parse(4, BoolEnd),
			depth(1),
			parse(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`[[[]]]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			depth(2),
			parse(1, ArrayStart),
			parse(1, ArrayEnd),
			parse(1, ArrayEnd),
			depth(1),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
}

// tests involving Parser.Skip()
var skipTests = []parserTest{
	{
		`[]`,
		[]step{
			parse(1, ArrayStart),
			skip(1, 0),
			parse(1, Continue),
			end(Done),
		},
	},
	{
		`[]`,
		[]step{
			parse(1, ArrayStart),
			skip(0, 1),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`{"foo":"bar"}`,
		[]step{
			parse(1, ObjectStart),
			skip(0, 1),
			parse(12, ObjectEnd),
			end(Done),
		},
	},
	{
		`[[{},2,3]]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			skip(0, 2),
			parse(7, ArrayEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[[[[[],1],2],3],4]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			skip(0, 5),
			parse(1, ArrayEnd),
			parse(3, ArrayEnd),
			parse(3, ArrayEnd),
			parse(3, ArrayEnd),
			parse(3, ArrayEnd),
			end(Done),
		},
	},
	{
		`[{"foo":"bar", "num": 1}, ["deeper", ["and deeper", 1, 2]]]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ObjectStart),
			skip(1, 1),
			parse(57, ArrayEnd),
			end(Done),
		},
	},
	{
		`{"key":123}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(4, KeyEnd),
			skip(0, 1),
			parse(5, ObjectEnd),
			end(Done),
		},
	},
	{
		`null`,
		[]step{
			parse(1, NullStart),
			skip(0, 1),
			parse(3, NullEnd),
			end(Done),
		},
	},
	{
		`[{"foo":"bar"},10]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(4, KeyEnd),
			parse(2, StringStart),
			skip(0, 3),
			parse(4, StringEnd),
			parse(1, ObjectEnd),
			parse(4, ArrayEnd),
			end(Done),
		},
	},
	{
		`[{"foo":"bar"}]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(4, KeyEnd),
			parse(2, StringStart),
			skip(1, 0),
			parse(5, ObjectEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[["skip"],[]]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			parse(1, StringStart),
			skip(1, 0),
			parse(6, ArrayEnd),
			skip(0, 1),
			parse(4, ArrayEnd),
			end(Done),
		},
	},
	{
		`{"a":"wub","wrong"}`,
		[]step{
			parse(1, ObjectStart),
			parse(1, KeyStart),
			parse(2, KeyEnd),
			parse(2, StringStart),
			parse(4, StringEnd),
			parse(2, KeyStart),
			skip(2, 0),
			parse(6, SyntaxError),
		},
	},
	{
		`[[],[]]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			skip(0, 1),
			parse(1, ArrayEnd),
			parse(2, ArrayStart),
			parse(1, ArrayEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[[],[]]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, ArrayStart),
			skip(1, 0),
			parse(3, ArrayStart),
			parse(1, ArrayEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`"hello"`,
		[]step{
			parse(1, StringStart),
			skip(1, 0),
			parse(6, Continue),
			end(Done),
		},
	},
	{
		`["a","b"]`,
		[]step{
			parse(1, ArrayStart),
			parse(1, StringStart),
			skip(1, 0),
			parse(4, StringStart),
			parse(2, StringEnd),
			parse(1, ArrayEnd),
			end(Done),
		},
	},
}

func TestParser(t *testing.T) {
	tests := make([]parserTest, 0)
	tests = append(tests, basicTests...)
	tests = append(tests, depthTests...)
	tests = append(tests, skipTests...)

	for _, test := range tests {
		b := []byte(test.json)
		o := "\np := Parser{}"
		p := &Parser{}

		for i := 0; i < len(test.steps); i++ {
			n, ok, log := test.steps[i](p, b)
			o = o + "\np" + log

			if !ok {
				t.Errorf(o)
				break
			}

			b = b[n:]
		}
	}
}
