package jo

import "testing"

type eventMethodTest struct {
	ev Event
	is bool
}

var literals = []eventMethodTest{
	{Continue, false},
	{Done, false},
	{ObjectStart, false},
	{ObjectEnd, false},
	{KeyStart, false},
	{KeyEnd, false},
	{ArrayStart, false},
	{ArrayEnd, false},
	{SyntaxError, false},
	{StringStart, true},
	{StringEnd, true},
	{NumberStart, true},
	{NumberEnd, true},
	{BoolStart, true},
	{BoolEnd, true},
	{NullStart, true},
	{NullEnd, true},
}

var starts = []eventMethodTest{
	{Continue, false},
	{Done, false},
	{ObjectStart, true},
	{ObjectEnd, false},
	{KeyStart, true},
	{KeyEnd, false},
	{ArrayStart, true},
	{ArrayEnd, false},
	{SyntaxError, false},
	{StringStart, true},
	{StringEnd, false},
	{NumberStart, true},
	{NumberEnd, false},
	{BoolStart, true},
	{BoolEnd, false},
	{NullStart, true},
	{NullEnd, false},
}

var ends = []eventMethodTest{
	{Continue, false},
	{Done, false},
	{ObjectStart, false},
	{ObjectEnd, true},
	{KeyStart, false},
	{KeyEnd, true},
	{ArrayStart, false},
	{ArrayEnd, true},
	{SyntaxError, false},
	{StringStart, false},
	{StringEnd, true},
	{NumberStart, false},
	{NumberEnd, true},
	{BoolStart, false},
	{BoolEnd, true},
	{NullStart, false},
	{NullEnd, true},
}

func TestEventMethods(t *testing.T) {
	for _, test := range literals {
		if test.ev.IsLiteral() != test.is {
			t.Errorf("%s.IsLiteral() != %v", test.ev, test.is)
		}
	}

	for _, test := range starts {
		if test.ev.IsStart() != test.is {
			t.Errorf("%s.IsStart() != %v", test.ev, test.is)
		}
	}

	for _, test := range ends {
		if test.ev.IsEnd() != test.is {
			t.Errorf("%s.IsEnd() != %v", test.ev, test.is)
		}
	}
}
