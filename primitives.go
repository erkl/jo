package jo

import (
	"errors"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

var (
	// returned by the ParseFoo functions on failure
	ErrSyntax   = errors.New("syntax error")
	ErrRange    = errors.New("out of range")
	ErrSigned   = errors.New("out of range")
	ErrFraction = errors.New("not an integer")
)

// ParseBool interprets a byte slice as a boolean value. Fails if the input
// is not a valid JSON boolean.
func ParseBool(data []byte) (bool, error) {
	data = trim(data)

	if len(data) == 4 && data[0] == 't' && data[1] == 'r' &&
		data[2] == 'u' && data[3] == 'e' {
		return true, nil
	}
	if len(data) == 5 && data[0] == 'f' && data[1] == 'a' &&
		data[2] == 'l' && data[3] == 's' && data[4] == 'e' {
		return false, nil
	}

	return false, ErrSyntax
}

// ParseUint interprets a byte slice as a 64-bit signed integer. Fails if
// the input is not a valid JSON number, or won't fit in an int64.
func ParseInt(data []byte) (int64, error) {
	data = trim(data)

	if len(data) == 0 {
		return 0, ErrSyntax
	}

	// for negative integers, ignore the minus sign
	if data[0] == '-' {
		unsigned, err := ParseUint(data[1:])
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

	unsigned, err := ParseUint(data)
	if err != nil {
		return 0, err
	}

	// watch out for integer overflows
	if unsigned > 1<<63-1 {
		return 0, ErrRange
	}

	return int64(unsigned), nil
}

// ParseUint interprets a byte slice as a 64-bit unsigned integer. Fails if
// the input is not a valid JSON number, or won't fit in a uint64.
func ParseUint(data []byte) (uint64, error) {
	data = trim(data)

	// very rudimentary error checking
	if len(data) == 0 {
		return 0, ErrSyntax
	}
	if data[0] == '-' {
		return 0, ErrSigned
	}

	// parse the base number
	base, pow10, rest, err := parseNumBase(data)
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

// parseNumBase interprets a JSON number's base (e.g. "123" in "123e10").
func parseNumBase(data []byte) (uint64, int, []byte, error) {
	n, e, r, d := uint64(0), 0, 0, -1

	// read all digits (and up to one dot)
	for ; r < len(data); r++ {
		b := data[r]

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

	return n, e, data[r:], nil
}

// parseNumExp interprets a JSON number's exponent (e.g. "e10" in "123e10").
func parseNumExp(data []byte) (int, error) {
	if e := data[0]; e != 'e' && e != 'E' {
		return 0, ErrSyntax
	}

	data = data[1:]

	// check for the optional sign
	neg := data[0] == '-'
	if neg || data[0] == '+' {
		data = data[1:]
	}

	// there's at least a chance of finding a digit, right?
	if len(data) == 0 {
		return 0, ErrSyntax
	}

	r, n := 0, 0

	for ; r < len(data); r++ {
		b := data[r]

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

// ParseFloat interprets the byte slice as a 64-bit float. Fails if the input
// is not a valid JSON number, or won't fit in a float64.
func ParseFloat(data []byte) (float64, error) {
	s := string(trim(data))

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

// Unquote interprets the input byte slice as a JSON string, returning the
// actual string value. Fails if the input is not a valid JSON number.
func Unquote(data []byte) (string, error) {
	data = trim(data)

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return "", ErrSyntax
	}

	data = data[1 : len(data)-1]
	r := 0

	// first, check for unusual characters
	for r < len(data) {
		b := data[r]

		if b == '\\' || b == '"' || b < ' ' {
			break
		}

		if b < utf8.RuneSelf {
			r++
			continue
		}

		rr, size := utf8.DecodeRune(data[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}

	// if no tricky characters were found we only have to cast the
	// input byte slice to a string
	if r == len(data) {
		return string(data), nil
	}

	buf := make([]byte, len(data)+2*utf8.UTFMax)
	w := copy(buf, data[:r])

	for r < len(data) {
		// are we out of room?
		if w >= len(buf)-2*utf8.UTFMax {
			next := make([]byte, 2*(len(buf)+utf8.UTFMax))
			copy(next, buf[:w])
			buf = next
		}

		switch b := data[r]; {
		case b == '\\':
			r++
			if r == len(data) {
				return "", ErrSyntax
			}

			b = data[r]

			// unicode escape sequences ("\u1234")
			if b == 'u' {
				rr := unquoteRune(data[r-1:])
				if rr < 0 {
					return "", ErrSyntax
				}
				r += 5

				if utf16.IsSurrogate(rr) {
					pair := utf16.DecodeRune(rr, unquoteRune(data[r:]))
					if pair != unicode.ReplacementChar {
						r += 6
						w += utf8.EncodeRune(buf[w:], pair)
						break
					}

					// invalid surrogate, fall back to the replacement rune
					rr = unicode.ReplacementChar
				}

				w += utf8.EncodeRune(buf[w:], rr)
				break
			}

			switch b {
			case '"', '\\', '/':
				buf[w] = b
			case 'b':
				buf[w] = '\b'
			case 'f':
				buf[w] = '\f'
			case 'n':
				buf[w] = '\n'
			case 'r':
				buf[w] = '\r'
			case 't':
				buf[w] = '\t'
			default:
				return "", ErrSyntax
			}

			r++
			w++

		case b == '"', b < ' ':
			return "", ErrSyntax

		case b < utf8.RuneSelf:
			buf[w] = b
			r++
			w++

		default:
			rr, size := utf8.DecodeRune(data[r:])
			r += size
			w += utf8.EncodeRune(buf[w:], rr)
		}
	}

	return string(buf[:w]), nil
}

// unquoteRune reads a unicode escape sequence ("\u1234") from the beginning of
// data, returning the rune it represents. Returns -1 if one cannot be found.
func unquoteRune(data []byte) (r rune) {
	if len(data) < 6 || data[0] != '\\' || data[1] != 'u' {
		return -1
	}

	for _, b := range data[2:6] {
		switch {
		case '0' <= b && b <= '9':
			r = r<<4 + rune(b-'0')
		case 'a' <= b && b <= 'f':
			r = r<<4 + rune(10+b-'a')
		case 'A' <= b && b <= 'F':
			r = r<<4 + rune(10+b-'A')
		default:
			return -1
		}
	}

	return r
}

// trim removes leading and trailing whitespace from a JSON value.
func trim(data []byte) []byte {
	for len(data) > 0 && isSpace(data[0]) {
		data = data[1:]
	}

	for len(data) > 0 {
		if last := len(data) - 1; isSpace(data[last]) {
			data = data[:last]
		} else {
			break
		}
	}

	return data
}
