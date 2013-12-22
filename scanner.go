package jo

// Valid scanner states.
const (
	sClean = iota
)

// Scanner is a JSON scanning state machine. The zero value for Scanner is
// ready to be used without any further initialization.
type Scanner struct {
	state int
	err   error
	stack []int
	size  int
}

// Scan feeds another byte to the scanner.
//
// The first return value, refered to as an opcode, tells the caller about
// significant parsing events like beginning and ending values, so that the
// caller can follow along if it wishes.
//
// The second return value is used to indicate whether or not, which is
// necessary because of number literals (is 123 a whole value, or just the
// beginning of 12345?). Even though it will only ever be either 1 or 0, it is
// an integer rather than a boolean to allow for these kinds of constructs:
//
//     for i := 0; i < len(input); {
//         _, n := s.Scan(input[i])
//         i += n
//     }
func (s *Scanner) Scan(c byte) (Op, int) {
	return OpContinue, 1
}

// Eof signals the scanner that the end of input has been reached. It returns
// an opcode, just as s.Scan does, which will be either OpEOF, OpNumberEnd or
// OpSyntaxError.
func (s *Scanner) Eof() Op {
	return OpEOF
}

// LastError returns the last error raised by the Scan or Eof methods.
func (s *Scanner) LastError() error {
	return s.err
}

// Reset resets the scanner's internal state.
func (s *Scanner) Reset() {
	s.state = 0
	s.err = nil
	s.stack = s.stack[:0]
	s.size = 0
}
