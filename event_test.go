package jo

import (
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
