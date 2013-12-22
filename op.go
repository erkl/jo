package jo

// The Op type holds opcodes returned by a Scanner.
type Op int

const (
	OpContinue Op = iota
	OpSpace       = iota
	OpEOF         = iota

	OpSyntaxError = iota

	OpObjectStart    = iota | OpFlagValueStart
	OpObjectEnd      = iota | OpFlagValueEnd
	OpObjectKeyStart = iota | OpFlagValueStart
	OpObjectKeyEnd   = iota | OpFlagValueEnd
	OpArrayStart     = iota | OpFlagValueStart
	OpArrayEnd       = iota | OpFlagValueEnd

	OpStringStart = iota | OpFlagValueStart
	OpStringEnd   = iota | OpFlagValueEnd
	OpNumberStart = iota | OpFlagValueStart
	OpNumberEnd   = iota | OpFlagValueEnd
	OpBoolStart   = iota | OpFlagValueStart
	OpBoolEnd     = iota | OpFlagValueEnd
	OpNullStart   = iota | OpFlagValueStart
	OpNullEnd     = iota | OpFlagValueEnd
)

// The OpFlagValueStart and OpFlagValueEnd opcode bit flags make it possible
// to detect arbitrary start or end events without testing every OpXStart or
// OpXEnd opcode individually.
const (
	OpFlagValueStart Op = 1 << (8 + iota)
	OpFlagValueEnd      = 1 << (8 + iota)
)

// String returns the string representation of any valid opcode.
func (op Op) String() string {
	switch op {
	case OpContinue:
		return "Continue"
	case OpSpace:
		return "Space"
	case OpEOF:
		return "EOF"
	case OpSyntaxError:
		return "SyntaxError"
	case OpStringStart:
		return "StringStart"
	case OpStringEnd:
		return "StringEnd"
	case OpNumberStart:
		return "NumberStart"
	case OpNumberEnd:
		return "NumberEnd"
	case OpBoolStart:
		return "BoolStart"
	case OpBoolEnd:
		return "BoolEnd"
	case OpNullStart:
		return "NullStart"
	case OpNullEnd:
		return "NullEnd"
	case OpObjectStart:
		return "ObjectStart"
	case OpObjectEnd:
		return "ObjectEnd"
	case OpObjectKeyStart:
		return "ObjectKeyStart"
	case OpObjectKeyEnd:
		return "ObjectKeyEnd"
	case OpArrayStart:
		return "ArrayStart"
	case OpArrayEnd:
		return "ArrayEnd"
	}
	return "illegal"
}
