package jo

import (
	"testing"
)

var parseBoolTests = []struct {
	in   string
	want bool
	err  error
}{
	{`true`, true, nil},
	{` true `, true, nil},
	{`false`, false, nil},
	{` false `, false, nil},
	{``, false, ErrSyntax},
	{`0`, false, ErrSyntax},
	{`tru`, false, ErrSyntax},
	{`alse`, false, ErrSyntax},
}

func TestParseBool(t *testing.T) {
	for _, test := range parseBoolTests {
		got, err := ParseBool([]byte(test.in))
		if got != test.want || err != test.err {
			t.Errorf("ParseBool(%q):", test.in)
			t.Errorf("   got %v, %v", got, err)
			t.Errorf("  want %v, %v", test.want, test.err)
		}
	}
}

var parseIntTests = []struct {
	in   string
	want int64
	err  error
}{
	{`0`, 0, nil},
	{`-0`, 0, nil},
	{`-10.0e2`, -1000, nil},
	{`0.5`, 0, ErrFraction},
	{`100e-3`, 0, ErrFraction},
	{`0e-2`, 0, nil},
	{`-1000`, -1000, nil},
	{`123456789`, 123456789, nil},
	{`-0.9223372036854775808e19`, -9223372036854775808, nil},
	{`0.9223372036854775807e19`, 9223372036854775807, nil},
	{`-0.2147483648e10`, -2147483648, nil},
	{`0.2147483647e10`, 2147483647, nil},
	{`-32768`, -32768, nil},
	{`32767`, 32767, nil},
	{`-128`, -128, nil},
	{`127`, 127, nil},
	{`-000000000000.00000000000000000000000000000000000000000000`, 0, nil},
	{`10000000000000000000000000000000000000000000000000000000e-55`, 1, nil},
	{`-9223372036854775808`, -9223372036854775808, nil},
	{`-9223372036854775809`, 0, ErrRange},
	{`9223372036854775807`, 9223372036854775807, nil},
	{`9223372036854775808`, 0, ErrRange},
	{`  123  `, 123, nil},
	{` 1 2 3 `, 0, ErrSyntax},
	{``, 0, ErrSyntax},
	{`--10`, 0, ErrSyntax},
	{`10a`, 0, ErrSyntax},
	{`0x20`, 0, ErrSyntax},
}

func TestParseInt(t *testing.T) {
	for _, test := range parseIntTests {
		got, err := ParseInt([]byte(test.in))
		if got != test.want || err != test.err {
			t.Errorf("ParseInt(%q):", test.in)
			t.Errorf("   got %d, %v", got, err)
			t.Errorf("  want %d, %v", test.want, test.err)
		}
	}
}

var parseUintTests = []struct {
	in   string
	want uint64
	err  error
}{
	{`0`, 0, nil},
	{`10`, 10, nil},
	{`101.00`, 101, nil},
	{`0.0101e4`, 101, nil},
	{`0.0101e3`, 0, ErrFraction},
	{`0e-2`, 0, nil},
	{`18446744073709551615`, 18446744073709551615, nil},
	{`1.8446744073709551615e19`, 18446744073709551615, nil},
	{`0.0018446744073709551615e22`, 18446744073709551615, nil},
	{`0.0018446744073709551615e21`, 0, ErrFraction},
	{`0.18446744073709551615e20`, 18446744073709551615, nil},
	{`0.4294967295e10`, 4294967295, nil},
	{`0.65535e5`, 65535, nil},
	{`0.255e3`, 255, nil},
	{`000000000000.000000000000000000000000000000000000000000000`, 0, nil},
	{`10000000000000000000000000000000000000000000000000000000e-55`, 1, nil},
	{` 1 2 3 `, 0, ErrSyntax},
	{``, 0, ErrSyntax},
	{`10a`, 0, ErrSyntax},
	{`0x20`, 0, ErrSyntax},
}

func TestParseUint(t *testing.T) {
	for _, test := range parseUintTests {
		got, err := ParseUint([]byte(test.in))
		if got != test.want || err != test.err {
			t.Errorf("ParseUint(%q):", test.in)
			t.Errorf("   got %d, %v", got, err)
			t.Errorf("  want %d, %v", test.want, test.err)
		}
	}
}

var parseFloatTests = []struct {
	in   string
	want float64
	err  error
}{
	{``, 0, ErrSyntax},
	{`10.0`, 10.0, nil},
	{`12345e10`, 12345e10, nil},
	{`-102030`, -102030, nil},
	{`1.7976931348623157e308`, 1.7976931348623157e308, nil},
	{`-1.7976931348623157e308`, -1.7976931348623157e308, nil},
	{"1.7976931348623159e308", 0, ErrRange},
	{"-1.7976931348623159e308", 0, ErrRange},
	{"1.7976931348623158e308", 1.7976931348623157e+308, nil},
	{"-1.7976931348623158e308", -1.7976931348623157e+308, nil},
	{"1.797693134862315808e308", 0, ErrRange},
	{"-1.797693134862315808e308", 0, ErrRange},
	{"1e308", 1e+308, nil},
	{"2e308", 0, ErrRange},
	{"1e-4294967296", 0, nil},
	{"1e+4294967296", 0, ErrRange},
	{"1e-18446744073709551616", 0, nil},
	{"1e+18446744073709551616", 0, ErrRange},
	{``, 0, ErrSyntax},
	{`1.2.3`, 0, ErrSyntax},
	{`.e1`, 0, ErrSyntax},
	{"1e", 0, ErrSyntax},
	{"1e-", 0, ErrSyntax},
	{".e-1", 0, ErrSyntax},
	{`+10`, 0, ErrSyntax},
	{`NaN`, 0, ErrSyntax},
	{`Inf`, 0, ErrSyntax},
	{`Infinity`, 0, ErrSyntax},
}

func TestParseFloat(t *testing.T) {
	for _, test := range parseFloatTests {
		got, err := ParseFloat([]byte(test.in))
		if got != test.want || err != test.err {
			t.Errorf("ParseFloat(%q):", test.in)
			t.Errorf("   got %f, %v", got, err)
			t.Errorf("  want %f, %v", test.want, test.err)
		}
	}
}
