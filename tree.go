package jo

import (
	"fmt"
)

// Indicates the type of value stored in a Node.
type Type int

const (
	Object = Type(ObjectStart)
	Array  = Type(ArrayStart)
	Key    = Type(KeyStart)
	String = Type(StringStart)
	Number = Type(NumberStart)
	Bool   = Type(BoolStart)
	Null   = Type(NullStart)
)

func (t Type) String() string {
	switch t {
	case Object:
		return "object"
	case Array:
		return "array"
	case Key:
		return "key"
	case String:
		return "string"
	case Number:
		return "number"
	case Bool:
		return "bool"
	case Null:
		return "null"
	}
	return "invalid type"
}

// A Node is an element in a tree of JSON values.
//
// A tree might look something like this (vertical arrows symbolize Node.Child
// pointers, and horizontal arrows symbolize Node.Next pointers):
//
//     {"foo":[1,null],"baz":{"a":"b"}}
//                    ↓
//                  "foo" → [1,"bar"] → "baz" → {"a":"b"}
//                               ↓                  ↓
//                               1 → null          "a" → "b"
type Node struct {
	Type  Type
	Bytes []byte
	Child *Node
	Next  *Node
}

// Temporary struct for storing state while parsing.
type marker struct {
	up    *marker
	start int
	node  *Node
	tail  **Node
}

// Generates a Node tree representation of the JSON value in data.
func Parse(data []byte) (*Node, error) {
	root := (*Node)(nil)
	top := &marker{nil, 0, nil, &root}

	p := Parser{}
	r := 0

	for r < len(data) {
		n, ev, err := p.Next(data[r:])
		r += n

		switch {
		case ev == SyntaxError:
			return nil, fmt.Errorf("%s at offset %d", err.Error(), r)

		case ev&Start != 0:
			v := &Node{Type: Type(ev)}
			*top.tail = v
			top.tail = &v.Next
			top = &marker{top, r - 1, v, &v.Child}

		case ev&End != 0:
			top.node.Bytes = data[top.start:r]
			top = top.up
		}
	}

	ev, err := p.End()

	switch {
	case ev == SyntaxError:
		return nil, err

	case ev != Done:
		top.node.Bytes = data[top.start:r]
		top = top.up
	}

	return root, nil
}
