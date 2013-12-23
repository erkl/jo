package jo

import (
	"fmt"
	"testing"
)

var scannerTests = []struct {
	input string
	steps []step
}{
	{
		`[123, null]`,
		[]step{
			ret(OpArrayStart, 1),  // '['
			ret(OpNumberStart, 1), // '1'
			ret(OpContinue, 1),    // '2'
			ret(OpContinue, 1),    // '3'
			ret(OpNumberEnd, 0),   // ','
			ret(OpContinue, 1),    // ',' (again)
			ret(OpSpace, 1),       // ' '
			ret(OpNullStart, 1),   // 'n'
			ret(OpContinue, 1),    // 'u'
			ret(OpContinue, 1),    // 'l'
			ret(OpNullEnd, 1),     // 'l'
			ret(OpArrayEnd, 1),    // ']'
			eof(OpEOF),            // EOF
		},
	},
	{
		`-0.123e+456`,
		[]step{
			ret(OpNumberStart, 1), // '-'
			ret(OpContinue, 1),    // '0'
			ret(OpContinue, 1),    // '.'
			ret(OpContinue, 1),    // '1'
			ret(OpContinue, 1),    // '2'
			ret(OpContinue, 1),    // '3'
			ret(OpContinue, 1),    // 'e'
			ret(OpContinue, 1),    // '+'
			ret(OpContinue, 1),    // '4'
			ret(OpContinue, 1),    // '5'
			ret(OpContinue, 1),    // '6'
			eof(OpNumberEnd),      // EOF
			eof(OpEOF),            // EOF (again)
		},
	},
	{
		`true`,
		[]step{
			ret(OpBoolStart, 1), // 't'
			ret(OpContinue, 1),  // 'r'
			ret(OpContinue, 1),  // 'u'
			ret(OpBoolEnd, 1),   // 'e'
			eof(OpEOF),          // EOF
		},
	},
	{
		`false`,
		[]step{
			ret(OpBoolStart, 1), // 'f'
			ret(OpContinue, 1),  // 'a'
			ret(OpContinue, 1),  // 'l'
			ret(OpContinue, 1),  // 's'
			ret(OpBoolEnd, 1),   // 'e'
			eof(OpEOF),          // EOF
		},
	},
	{
		`null`,
		[]step{
			ret(OpNullStart, 1), // 'n'
			ret(OpContinue, 1),  // 'u'
			ret(OpContinue, 1),  // 'l'
			ret(OpNullEnd, 1),   // 'l'
			eof(OpEOF),          // EOF
		},
	},
}

func TestScanner(t *testing.T) {
	for _, test := range scannerTests {
		var s Scanner
		var l []string
		var i int

		for _, step := range test.steps {
			n, desc := step(&s, test.input, i)

			l = append(l, "  "+desc)
			i += n

			if n < 0 {
				t.Errorf("- %#q", test.input)
				for _, s := range l {
					t.Error(s)
				}
				break
			}
		}
	}
}

// A step function describes a step in a scanner test case.
type step func(s *Scanner, in string, i int) (int, string)

func ret(wantOp Op, wantN int) step {
	return func(s *Scanner, in string, i int) (int, string) {
		op, n := s.Scan(in[i])
		desc := fmt.Sprintf(".Scan(%q) -> %s, %d", in[i], op, n)

		if op != wantOp || n != wantN {
			return -1, fmt.Sprintf("%s (want %s, %d)", desc, wantOp, wantN)
		} else {
			return n, desc
		}
	}
}

func eof(wantOp Op) step {
	return func(s *Scanner, in string, i int) (int, string) {
		op := s.Eof()
		desc := fmt.Sprintf(".Eof() -> %s", op)

		if op != wantOp {
			return -1, fmt.Sprintf("%s (want %s)", desc, wantOp)
		} else {
			return 0, desc
		}
	}
}
