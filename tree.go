package jo

import (
	"bytes"
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

// Marshals the value back to JSON.
func (n *Node) MarshalJSON() ([]byte, error) {
	if n.Type != Object && n.Type != Array {
		return n.Bytes, nil
	}

	// composite values require slightly more work, so we'll
	// have to use a buffer
	buf := &bytes.Buffer{}
	if _, err := n.WriteJSON(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Duplicate of io.Writer.
type Writer interface {
	Write(data []byte) (int, error)
}

// WriteJSON serializes the node and its children, just like MarshalJSON()
// would, but writes its output straight to a Writer.
func (n *Node) WriteJSON(w Writer) (int, error) {
	var total int
	var err error

	if n.Type == Object {
		total, err = writeObject(n, w)
	} else if n.Type == Array {
		total, err = writeArray(n, w)
	} else {
		total, err = w.Write(n.Bytes)
	}

	return total, err
}

// Marshals the object node n into w.
func writeObject(n *Node, w Writer) (int, error) {
	// write the opening brace
	total, err := w.Write([]byte{'{'})
	if err != nil {
		return total, err
	}

	c := n.Child
	i := 0

	for c != nil {
		// insert colons after keys, and commas between each value
		// and the next key
		if i > 0 {
			if i&1 != 0 {
				_, err = w.Write([]byte{':'})
			} else {
				_, err = w.Write([]byte{','})
			}

			if err != nil {
				return total, err
			}
			total++
		}

		// write the child element (either a key or value)
		n, err := c.WriteJSON(w)
		if err != nil {
			return total + n, err
		}

		total += n
		c = c.Next
		i++
	}

	// write the closing brace
	_, err = w.Write([]byte{'}'})
	if err != nil {
		return total, err
	}

	return total + 1, nil
}

// Marshals the array node n into w.
func writeArray(n *Node, w Writer) (int, error) {
	// write the opening bracket
	total, err := w.Write([]byte{'['})
	if err != nil {
		return total, err
	}

	c := n.Child
	i := 0

	for c != nil {
		// insert commas between elements in the array
		if i > 0 {
			_, err = w.Write([]byte{','})
			if err != nil {
				return total, err
			}
			total++
		}

		// write the child element
		n, err := c.WriteJSON(w)
		if err != nil {
			return total + n, err
		}

		total += n
		c = c.Next
		i++
	}

	// write the closing bracket
	_, err = w.Write([]byte{']'})
	if err != nil {
		return total, err
	}

	return total + 1, nil
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
