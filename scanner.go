package jo

// A Scanner is a state machine which emits a series of Events when fed JSON
// input.
type Scanner struct {
}

// NewScanner initializes a new Scanner.
func NewScanner() *Scanner {
	return new(Scanner)
}

// Reset restores a Scanner to its initial state.
func (s *Scanner) Reset() {
}

// Scan accepts a byte of input and returns an Event.
func (s *Scanner) Scan(c byte) Event {
	return None
}

// End signals the Scanner that the end of input has been reached. It returns
// an event just as Scan does.
func (s *Scanner) End() Event {
	return None
}

// LastError returns a syntax error description after either Scan or End has
// returned an Error event.
func (s *Scanner) LastError() error {
	return nil
}
