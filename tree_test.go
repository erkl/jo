package jo

import (
	"errors"
	"reflect"
	"testing"
)

var parseTests = []struct {
	in   string
	want *Node
	err  error
}{
	{
		`"hello"`,
		&Node{
			String,
			[]byte(`"hello"`),
			nil,
			nil,
		},
		nil,
	},
	{
		`123`,
		&Node{
			Number,
			[]byte(`123`),
			nil,
			nil,
		},
		nil,
	},
	{
		" \t null ",
		&Node{
			Null,
			[]byte(`null`),
			nil,
			nil,
		},
		nil,
	},
	{
		`[1,2]`,
		&Node{
			Array,
			[]byte(`[1,2]`),
			&Node{
				Number,
				[]byte(`1`),
				nil,
				&Node{
					Number,
					[]byte(`2`),
					nil,
					nil,
				},
			},
			nil,
		},
		nil,
	},
	{
		`{ "foo": "bar" }`,
		&Node{
			Object,
			[]byte(`{ "foo": "bar" }`),
			&Node{
				Key,
				[]byte(`"foo"`),
				nil,
				&Node{
					String,
					[]byte(`"bar"`),
					nil,
					nil,
				},
			},
			nil,
		},
		nil,
	},
	{
		`[[0],1]`,
		&Node{
			Array,
			[]byte(`[[0],1]`),
			&Node{
				Array,
				[]byte(`[0]`),
				&Node{
					Number,
					[]byte(`0`),
					nil,
					nil,
				},
				&Node{
					Number,
					[]byte(`1`),
					nil,
					nil,
				},
			},
			nil,
		},
		nil,
	},
	{
		`{"foo":[1,"bar"]}`,
		&Node{
			Object,
			[]byte(`{"foo":[1,"bar"]}`),
			&Node{
				Key,
				[]byte(`"foo"`),
				nil,
				&Node{
					Array,
					[]byte(`[1,"bar"]`),
					&Node{
						Number,
						[]byte(`1`),
						nil,
						&Node{
							String,
							[]byte(`"bar"`),
							nil,
							nil,
						},
					},
					nil,
				},
			},
			nil,
		},
		nil,
	},
	{
		`{"a""b"}`,
		nil,
		errors.New("expected ':', found '\"' at offset 4"),
	},
	{
		`[]]`,
		nil,
		errors.New("expected end of input, found ']' at offset 2"),
	},
	{
		`[123e`,
		nil,
		errors.New("expected '-', '+' or digit, found end of input"),
	},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		actual, err := Parse([]byte(test.in))
		if !reflect.DeepEqual(actual, test.want) ||
			!reflect.DeepEqual(err, test.err) {
			t.Errorf("Parse(%#q):", test.in)
			t.Errorf("   got %v, %q", actual, err)
			t.Errorf("  want %v, %q", test.want, test.err)
		}
	}
}

func TestNodeMarshalJSON(t *testing.T) {
	in := ` [ { "a" : { "b" : true, "c" : false } } , [ null , 1.23 ] , [ ] ] `
	want := `[{"a":{"b":true,"c":false}},[null,1.23],[]]`

	root, err := Parse([]byte(in))
	if err != nil {
		t.Fatalf("setup failed")
	}

	actual, err := root.MarshalJSON()
	if string(actual) != want || err != nil {
		t.Errorf("Node.MarshalJSON():")
		t.Errorf("   got %#q, %v", actual, err)
		t.Errorf("  want %#q, <nil>", want)
	}
}
