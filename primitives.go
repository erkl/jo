package jo

import (
	"errors"
	"strconv"
)

var (
	ErrSyntax   = errors.New("syntax error")
	ErrRange    = errors.New("out of range")
	ErrSigned   = errors.New("out of range")
	ErrFraction = errors.New("not an integer")
)

// Parses a byte slice into a bool. Fails if the input is not a
// valid JSON boolean value.
func ParseBool(bytes []byte) (bool, error) {
	bytes = trim(bytes)

	if len(bytes) == 4 && bytes[0] == 't' && bytes[1] == 'r' &&
		bytes[2] == 'u' && bytes[3] == 'e' {
		return true, nil
	}
	if len(bytes) == 5 && bytes[0] == 'f' && bytes[1] == 'a' &&
		bytes[2] == 'l' && bytes[3] == 's' && bytes[4] == 'e' {
		return false, nil
	}

	return false, ErrSyntax
}

// Parses a byte slice as a 64-bit signed integer. Fails if the input
// is not a valid JSON number, or won't fit in an int64.
func ParseInt(bytes []byte) (int64, error) {
	bytes = trim(bytes)

	if len(bytes) == 0 {
		return 0, ErrSyntax
	}

	// for negative integers, ignore the minus sign
	if bytes[0] == '-' {
		unsigned, err := ParseUint(bytes[1:])
		if err != nil {
			if err == ErrSigned {
				return 0, ErrSyntax
			}
			return 0, err
		}

		// make sure the value isn't too great to fit in 63 bits
		if unsigned > 1<<63 {
			return 0, ErrRange
		}

		return -int64(unsigned), nil
	}

	unsigned, err := ParseUint(bytes)
	if err != nil {
		return 0, err
	}

	// watch out for integer overflows
	if unsigned > 1<<63-1 {
		return 0, ErrRange
	}

	return int64(unsigned), nil
}

// Parses a byte slice as a 64-bit unsigned integer. Fails if the input
// is not a valid JSON number, or won't fit in a uint64.
func ParseUint(bytes []byte) (uint64, error) {
	bytes = trim(bytes)

	// very rudimentary error checking
	if len(bytes) == 0 {
		return 0, ErrSyntax
	}
	if bytes[0] == '-' {
		return 0, ErrSigned
	}

	// parse the base number
	base, pow10, rest, err := parseNumBase(bytes)
	if err != nil {
		return 0, err
	}

	// any remaining bytes must be part of an exponent sequence (i.e.
	// begin with e or E)
	if len(rest) > 0 {
		// the shortest legal exponent is 2 bytes long, so anything
		// shorter than that could not possibly be valid
		if len(rest) < 2 {
			return 0, ErrSyntax
		}

		exp, err := parseNumExp(rest)
		if err != nil {
			return 0, err
		}

		// make sure the exponent isn't large enough to cause an
		// integer overflow
		if (exp < 0 && pow10+exp > pow10) || (exp > 0 && pow10+exp < pow10) {
			return 0, ErrSyntax
		}

		pow10 += exp
	}

	// let's not make this more complicated than it has to be
	if base == 0 || pow10 == 0 {
		return base, nil
	}

	// a negative pow10 means we've got a non-integer value on
	// our hands; bail!
	if pow10 < 0 {
		return 0, ErrFraction
	}

	// adjust the base value to account for pow10
	for ; pow10 > 0; pow10-- {
		if base*10 < base {
			return 0, ErrRange
		}
		base *= 10
	}

	return base, nil
}

// Parses a JSON number's base (e.g. "123" in "123e10").
func parseNumBase(bytes []byte) (uint64, int, []byte, error) {
	n, e, r, d := uint64(0), 0, 0, -1

	// read all digits (and up to one dot)
	for ; r < len(bytes); r++ {
		b := bytes[r]

		if b < '0' || b > '9' {
			if b == '.' && d == -1 {
				d = r
				continue
			}

			// we've found all valid characters, time to bail
			break
		}

		// when enconutering a zero we don't immediately multiply n by 10;
		// instead we keep track of how many zeros we've encountered in a row
		// (dealing with zeros the naive way would mean that we wouldn't be
		// able to parse numbers like "100000000000000000000e-20")
		if b == '0' {
			e++
			continue
		}

		// time to pay for all those zeroes we put off accounting for
		for ; e > 0; e-- {
			if n*10 < n {
				return 0, 0, nil, ErrRange
			}
			n *= 10
		}

		// guard against integer overflows
		tmp := n*10 + uint64(b-'0')
		if tmp < n {
			return 0, 0, nil, ErrRange
		}

		n = tmp
	}

	// time for some retrospective error checking
	switch {
	case r == 0:
		return 0, 0, nil, ErrSyntax
	case d == 0:
		return 0, 0, nil, ErrSyntax
	case d == r-1:
		return 0, 0, nil, ErrSyntax
	}

	// if we encountered a dot, adjust e accordingly
	if d > 0 {
		e += 1 + d - r
	}

	return n, e, bytes[r:], nil
}

// Parses a JSON number's exponent (e.g. "e10" in "123e10").
func parseNumExp(bytes []byte) (int, error) {
	if e := bytes[0]; e != 'e' && e != 'E' {
		return 0, ErrSyntax
	}

	bytes = bytes[1:]

	// check for the optional sign
	neg := bytes[0] == '-'
	if neg || bytes[0] == '+' {
		bytes = bytes[1:]
	}

	// there's at least a chance of finding a digit, right?
	if len(bytes) == 0 {
		return 0, ErrSyntax
	}

	r, n := 0, 0

	for ; r < len(bytes); r++ {
		b := bytes[r]

		// only digits are allowed at this point
		if b < '0' || b > '9' {
			return 0, ErrSyntax
		}

		// simple integer overflow detection
		tmp := n*10 + int(b-'0')
		if tmp < n {
			return 0, ErrSyntax
		}

		n = tmp
	}

	// negate the result if we found a '-' earlier
	if neg {
		n = -n
	}

	return n, nil
}

// Parses a byte slice as a 64-bit float. Fails if the input is not
// a valid JSON number, or won't fit in a float64.
func ParseFloat(bytes []byte) (float64, error) {
	s := string(trim(bytes))

	if len(s) == 0 {
		return 0, ErrSyntax
	}

	// quickly check for values which strconv.ParseFloat will happily parse
	// for us, but which aren't necessarily valid JSON numbers
	if f, l := s[0], s[len(s)-1]; f == '+' || l < '0' || l > '9' {
		return 0, ErrSyntax
	}

	// float parsing is difficult to get right, so let's reuse the
	// work of much cleverer people
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		if err.(*strconv.NumError).Err == strconv.ErrRange {
			return 0, ErrRange
		}
		return 0, ErrSyntax
	}

	return f, nil
}

// Removes leading and trailing whitespace from a JSON value.
func trim(bytes []byte) []byte {
	for len(bytes) > 0 && isSpace(bytes[0]) {
		bytes = bytes[1:]
	}

	for len(bytes) > 0 {
		if last := len(bytes) - 1; isSpace(bytes[last]) {
			bytes = bytes[:last]
		} else {
			break
		}
	}

	return bytes
}
