package jo

import (
	"fmt"
	"testing"
)

var eventStringTests = []struct {
	in  Event
	out string
}{
	{
		None,
		"None",
	},
	{
		Error,
		"Error",
	},
	{
		StringStart,
		"StringStart",
	},
	{
		NumberEnd | Space,
		"NumberEnd | Space",
	},
	{
		ObjectEnd | KeyEnd | ArrayEnd | StringEnd | NumberEnd | BoolEnd | NullEnd | ObjectStart | KeyStart | ArrayStart | StringStart | NumberStart | BoolStart | NullStart | Space,
		"ObjectEnd | KeyEnd | ArrayEnd | StringEnd | NumberEnd | BoolEnd | NullEnd | ObjectStart | KeyStart | ArrayStart | StringStart | NumberStart | BoolStart | NullStart | Space",
	},
	{
		Error - 1,
		"INVALID",
	},
}

func TestEventString(t *testing.T) {
	for _, test := range eventStringTests {
		s := test.in.String()
		if s != test.out {
			t.Errorf("(%s).String():", test.out)
			t.Errorf("  got  %q", s)
			t.Errorf("  want %q", test.out)
		}
	}
}

func ExampleScanner() {
	var s = NewScanner()

	for _, c := range `{ "foo": 123 }` {
		fmt.Printf("at %q: %s\n", c, s.Scan(byte(c)))
	}

	fmt.Printf("at EOF: %s\n", s.End())
	// Output:
	// at '{': ObjectStart
	// at ' ': Space
	// at '"': KeyStart
	// at 'f': None
	// at 'o': None
	// at 'o': None
	// at '"': None
	// at ':': KeyEnd
	// at ' ': Space
	// at '1': NumberStart
	// at '2': None
	// at '3': None
	// at ' ': NumberEnd | Space
	// at '}': None
	// at EOF: ObjectEnd
}

var scannerTests = []struct {
	in  string
	out []Event
}{
	{
		` { } `,
		[]Event{
			Space,             // ' '
			ObjectStart,       // '{'
			Space,             // ' '
			None,              // '}'
			ObjectEnd | Space, // ' '
			None,              // EOF
		},
	},
	{
		`{ "a": { "b":{} , "c" : 0} }`,
		[]Event{
			ObjectStart,       // '{'
			Space,             // ' '
			KeyStart,          // '"'
			None,              // 'a'
			None,              // '"'
			KeyEnd,            // ':'
			Space,             // ' '
			ObjectStart,       // '{'
			Space,             // ' '
			KeyStart,          // '"'
			None,              // 'b'
			None,              // '"'
			KeyEnd,            // ':'
			ObjectStart,       // '{'
			None,              // '}'
			ObjectEnd | Space, // ' '
			None,              // ','
			Space,             // ' '
			KeyStart,          // '"'
			None,              // 'c'
			None,              // '"'
			KeyEnd | Space,    // ' '
			None,              // ':'
			Space,             // ' '
			NumberStart,       // '0'
			NumberEnd,         // '}'
			ObjectEnd | Space, // ' '
			None,              // '}'
			ObjectEnd,         // EOF
		},
	},
	{
		` [ ] `,
		[]Event{
			Space,            // ' '
			ArrayStart,       // '['
			Space,            // ' '
			None,             // ']'
			ArrayEnd | Space, // ' '
			None,             // EOF
		},
	},
	{
		`[[],[[]]]`,
		[]Event{
			ArrayStart, // '['
			ArrayStart, // '['
			None,       // ']'
			ArrayEnd,   // ','
			ArrayStart, // '['
			ArrayStart, // '['
			None,       // ']'
			ArrayEnd,   // ']'
			ArrayEnd,   // ']'
			ArrayEnd,   // EOF
		},
	},
	{
		`[0,1, 2 ,3 , 4]`,
		[]Event{
			ArrayStart,        // '['
			NumberStart,       // '0'
			NumberEnd,         // ','
			NumberStart,       // '1'
			NumberEnd,         // ','
			Space,             // ' '
			NumberStart,       // '2'
			NumberEnd | Space, // ' '
			None,              // ','
			NumberStart,       // '3'
			NumberEnd | Space, // ' '
			None,              // ','
			Space,             // ' '
			NumberStart,       // '4'
			NumberEnd,         // ']'
			ArrayEnd,          // EOF
		},
	},
	{
		`"foo"`,
		[]Event{
			StringStart, // '"'
			None,        // 'f'
			None,        // 'o'
			None,        // 'o'
			None,        // '"'
			StringEnd,   // EOF
		},
	},
	{
		" \" bar\"\n ",
		[]Event{
			Space,             // ' '
			StringStart,       // '"'
			None,              // ' '
			None,              // 'b'
			None,              // 'a'
			None,              // 'r'
			None,              // '"'
			StringEnd | Space, // '\n'
			Space,             // ' '
			None,              // EOF
		},
	},
	{
		`"\b\f\n\r\t\\\/\""`,
		[]Event{
			StringStart, // '"'
			None,        // '\\'
			None,        // 'b'
			None,        // '\\'
			None,        // 'f'
			None,        // '\\'
			None,        // 'n'
			None,        // '\\'
			None,        // 'r'
			None,        // '\\'
			None,        // 't'
			None,        // '\\'
			None,        // '\\'
			None,        // '\\'
			None,        // '/'
			None,        // '\\'
			None,        // '"'
			None,        // '"'
			StringEnd,   // EOF
		},
	},
	{
		`"\u2603 = ☃"`,
		[]Event{
			StringStart, // '"'
			None,        // '\\'
			None,        // 'u'
			None,        // '2'
			None,        // '6'
			None,        // '0'
			None,        // '3'
			None,        // ' '
			None,        // '='
			None,        // ' '
			None,        // '\xE2'
			None,        // '\x98'
			None,        // '\x83'
			None,        // '"'
			StringEnd,   // EOF
		},
	},
	{
		`0 `,
		[]Event{
			NumberStart,       // '0'
			NumberEnd | Space, // ' '
			None,              // EOF
		},
	},
	{
		`1 `,
		[]Event{
			NumberStart,       // '1'
			NumberEnd | Space, // ' '
			None,              // EOF
		},
	},
	{
		`2.5 `,
		[]Event{
			NumberStart,       // '2'
			None,              // '.'
			None,              // '5'
			NumberEnd | Space, // ' '
			None,              // EOF
		},
	},
	{
		`0.1e+2`,
		[]Event{
			NumberStart, // '0'
			None,        // '.'
			None,        // '1'
			None,        // 'e'
			None,        // '+'
			None,        // '2'
			NumberEnd,   // EOF
		},
	},
	{
		`-0.003`,
		[]Event{
			NumberStart, // '-'
			None,        // '0'
			None,        // '.'
			None,        // '0'
			None,        // '0'
			None,        // '3'
			NumberEnd,   // EOF
		},
	},
	{
		`-1000E10`,
		[]Event{
			NumberStart, // '-'
			None,        // '1'
			None,        // '0'
			None,        // '0'
			None,        // '0'
			None,        // 'E'
			None,        // '1'
			None,        // '0'
			NumberEnd,   // EOF
		},
	},
	{
		`true`,
		[]Event{
			BoolStart, // 't'
			None,      // 'r'
			None,      // 'u'
			None,      // 'e'
			BoolEnd,   // EOF
		},
	},
	{
		`false`,
		[]Event{
			BoolStart, // 'f'
			None,      // 'a'
			None,      // 'l'
			None,      // 's'
			None,      // 'e'
			BoolEnd,   // EOF
		},
	},
	{
		`null`,
		[]Event{
			NullStart, // 'n'
			None,      // 'u'
			None,      // 'l'
			None,      // 'l'
			NullEnd,   // EOF
		},
	},
	{
		`x`,
		[]Event{
			Error, // 'x'
		},
	},
	{
		` `,
		[]Event{
			Space, // ' '
			Error, // EOF
		},
	},
	{
		` 0 9`,
		[]Event{
			Space,             // ' '
			NumberStart,       // '0'
			NumberEnd | Space, // ' '
			Error,             // '9'
		},
	},
	{
		`{ : "foo" }`,
		[]Event{
			ObjectStart, // '{'
			Space,       // ' '
			Error,       // ':'
		},
	},
	{
		`{ "foo" "bar" }`,
		[]Event{
			ObjectStart,    // '{'
			Space,          // ' '
			KeyStart,       // '"'
			None,           // 'f'
			None,           // 'o'
			None,           // 'o'
			None,           // '"'
			KeyEnd | Space, // ' '
			Error,          // '"'
		},
	},
	{
		`{ "foo": }`,
		[]Event{
			ObjectStart, // '{'
			Space,       // ' '
			KeyStart,    // '"'
			None,        // 'f'
			None,        // 'o'
			None,        // 'o'
			None,        // '"'
			KeyEnd,      // ':'
			Space,       // ' '
			Error,       // '}'
		},
	},
	{
		`{ "foo": true false }`,
		[]Event{
			ObjectStart,     // '{'
			Space,           // ' '
			KeyStart,        // '"'
			None,            // 'f'
			None,            // 'o'
			None,            // 'o'
			None,            // '"'
			KeyEnd,          // ':'
			Space,           // ' '
			BoolStart,       // 't'
			None,            // 'r'
			None,            // 'u'
			None,            // 'e'
			BoolEnd | Space, // ' '
			Error,           // 'f'
		},
	},
	{
		`{ "foo": true, }`,
		[]Event{
			ObjectStart, // '{'
			Space,       // ' '
			KeyStart,    // '"'
			None,        // 'f'
			None,        // 'o'
			None,        // 'o'
			None,        // '"'
			KeyEnd,      // ':'
			Space,       // ' '
			BoolStart,   // 't'
			None,        // 'r'
			None,        // 'u'
			None,        // 'e'
			BoolEnd,     // ','
			Space,       // ' '
			Error,       // '}'
		},
	},
	{
		`[1 2]`,
		[]Event{
			ArrayStart,        // '['
			NumberStart,       // '1'
			NumberEnd | Space, // ' '
			Error,             // '2'
		},
	},
	{
		"\"\t\"",
		[]Event{
			StringStart, // '"'
			Error,       // '\t'
		},
	},
	{
		`"\U1234"`,
		[]Event{
			StringStart, // '"'
			None,        // '\\'
			Error,       // 'U'
		},
	},
	{
		`"\ux"`,
		[]Event{
			StringStart, // '"'
			None,        // '\\'
			None,        // 'u'
			Error,       // 'x'
		},
	},
	{
		`"\u1x"`,
		[]Event{
			StringStart, // '"'
			None,        // '\\'
			None,        // 'u'
			None,        // '1'
			Error,       // 'x'
		},
	},
	{
		`"\u12x"`,
		[]Event{
			StringStart, // '"'
			None,        // '\\'
			None,        // 'u'
			None,        // '1'
			None,        // '2'
			Error,       // 'x'
		},
	},
	{
		`"\u123x"`,
		[]Event{
			StringStart, // '"'
			None,        // '\\'
			None,        // 'u'
			None,        // '1'
			None,        // '2'
			None,        // '3'
			Error,       // 'x'
		},
	},
	{
		`-.5`,
		[]Event{
			NumberStart, // '-'
			Error,       // '.'
		},
	},
	{
		`0. `,
		[]Event{
			NumberStart, // '0'
			None,        // '.'
			Error,       // ' '
		},
	},
	{
		`0.0e`,
		[]Event{
			NumberStart, // '0'
			None,        // '.'
			None,        // '0'
			None,        // 'e'
			Error,       // EOF
		},
	},
	{
		`0.0e+ `,
		[]Event{
			NumberStart, // '0'
			None,        // '.'
			None,        // '0'
			None,        // 'e'
			None,        // '+'
			Error,       // ' '
		},
	},
	{
		`tx`,
		[]Event{
			BoolStart, // 't'
			Error,     // 'x'
		},
	},
	{
		`tr`,
		[]Event{
			BoolStart, // 't'
			None,      // 'r'
			Error,     // EOF
		},
	},
	{
		`truE`,
		[]Event{
			BoolStart, // 't'
			None,      // 'r'
			None,      // 'u'
			Error,     // 'E'
		},
	},
	{
		`fx`,
		[]Event{
			BoolStart, // 'f'
			Error,     // 'x'
		},
	},
	{
		`faL`,
		[]Event{
			BoolStart, // 'f'
			None,      // 'a'
			Error,     // 'L'
		},
	},
	{
		`falx`,
		[]Event{
			BoolStart, // 'f'
			None,      // 'a'
			None,      // 'l'
			Error,     // 'x'
		},
	},
	{
		`fals`,
		[]Event{
			BoolStart, // 'f'
			None,      // 'a'
			None,      // 'l'
			None,      // 's'
			Error,     // EOF
		},
	},
	{
		`nU`,
		[]Event{
			NullStart, // 'n'
			Error,     // 'U'
		},
	},
	{
		`nuL`,
		[]Event{
			NullStart, // 'n'
			None,      // 'u'
			Error,     // 'L'
		},
	},
	{
		`nul `,
		[]Event{
			NullStart, // 'n'
			None,      // 'u'
			None,      // 'l'
			Error,     // ' '
		},
	},
}

func TestScanner(t *testing.T) {
	for _, test := range scannerTests {
		var s = NewScanner()
		var ev Event

		for i, want := range test.out {
			if i < len(test.in) {
				ev = s.Scan(test.in[i])
			} else {
				ev = s.End()
			}

			if ev != want {
				t.Errorf("Scanner(%#q):", test.in)

				for j, prev := range test.out[:i] {
					if j < len(test.in) {
						t.Errorf("  %4q -> %s", test.in[j], prev)
					} else {
						t.Errorf("   EOF -> %s", prev)
					}
				}

				if i < len(test.in) {
					t.Errorf("  %4q -> %s (want %s)", test.in[i], ev, test.out[i])
				} else {
					t.Errorf("   EOF -> %s (want %s)", ev, test.out[i])
				}

				break
			}
		}
	}
}

func TestScannerErrors(t *testing.T) {
	var s = NewScanner()

	if s.Scan('1') != NumberStart || s.Scan('x') != Error {
		t.Fatalf("setup failed")
	}

	err := s.LastError()
	if err == nil {
		t.Fatalf("Scanner.LastError returned nil after Error event")
	}

	if s.Scan('2') != Error || s.LastError().Error() != err.Error() {
		t.Fatalf("Scanner.Scan did not remember previous error")
	}

	if s.End() != Error || s.LastError().Error() != err.Error() {
		t.Fatalf("Scanner.End did not remember previous error")
	}
}
