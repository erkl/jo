package jo

// NodeType indicates the type of the value stored in a Node.
type NodeType int

const (
	ObjectNode NodeType = iota
	ObjectKeyNode
	ArrayNode
	StringNode
	NumberNode
	BoolNode
	NullNode
)

// A Node is a part of a tree describing a JSON value. It contains
// information about the node's type, its original JSON representation
// (for literal values), and information about its position inside
// the tree.
type Node struct {
	Type NodeType
	Raw  []byte

	Parent                   *Node
	PrevSibling, NextSibling *Node
	FirstChild, LastChild    *Node
}

// AppendChild adds a child to the subject Node.
func (n *Node) AppendChild(c *Node) {
	if c.Parent != nil || c.PrevSibling != nil || c.NextSibling != nil {
		panic("jo: AppendChild called with an already attached Node")
	}

	if n.FirstChild == nil {
		n.FirstChild = c
	} else {
		n.LastChild.NextSibling = c
	}

	c.Parent = n
	c.PrevSibling = n.LastChild
	n.LastChild = c
}

// InsertChild adds a child to the subject Node, immediately before a
// reference node. Calling InsertChild with a nil reference node is
// functionally equivalent to just calling AppendChild.
func (n *Node) InsertChild(c, before *Node) {
	switch {
	case c.Parent != nil || c.PrevSibling != nil || c.NextSibling != nil:
		panic("jo: InsertChild called with an already attached Node")
	case before != nil && before.Parent != n:
		panic("jo: InsertChild called with unknown reference node")
	}

	c.Parent = n
	c.NextSibling = before

	if before == nil {
		if n.LastChild != nil {
			n.LastChild.NextSibling = c
		}
		c.PrevSibling = n.LastChild
		n.LastChild = c
	} else {
		if before.PrevSibling != nil {
			before.PrevSibling.NextSibling = c
		}
		c.PrevSibling = before.PrevSibling
		before.PrevSibling = c
	}

	if n.FirstChild == before {
		n.FirstChild = c
	}
}

// RemoveChild removes a child from the subject Node.
func (n *Node) RemoveChild(c *Node) {
	if c.Parent != n {
		panic("jo: RemoveChild called for non-child Node")
	}

	if n.FirstChild == c {
		n.FirstChild = c.NextSibling
	}
	if n.LastChild == c {
		n.LastChild = c.PrevSibling
	}

	if c.NextSibling != nil {
		c.NextSibling.PrevSibling = c.PrevSibling
	}
	if c.PrevSibling != nil {
		c.PrevSibling.NextSibling = c.NextSibling
	}

	c.Parent = nil
	c.PrevSibling = nil
	c.NextSibling = nil
}

// Parse parses raw JSON input and generates its tree representation.
func Parse(buf []byte) (*Node, error) {
	var p = parseState{buf: buf}

	root, err := p.close(p.scan())
	if err != nil {
		return nil, err
	}

	// consume all bytes after the top-level value
	if p.scan() == OpSyntaxError {
		return nil, p.s.LastError()
	}

	return root, nil
}

// parseState holds state when parsing a JSON value.
type parseState struct {
	s   Scanner
	buf []byte
	off int
}

// close completely parses the Node opened by op.
func (p *parseState) close(op Op) (*Node, error) {
	switch op {
	case OpSyntaxError:
		return nil, p.s.LastError()
	case OpObjectStart:
		return p.closeObject()
	case OpArrayStart:
		return p.closeArray()
	default:
		return p.closeLiteral()
	}
}

func (p *parseState) closeObject() (*Node, error) {
	obj := &Node{Type: ObjectNode}

	for {
		switch op := p.scan(); op {
		case OpSyntaxError:
			return nil, p.s.LastError()

		case OpObjectEnd:
			return obj, nil

		case OpObjectKeyStart:
			key, err := p.closeLiteral()
			if err != nil {
				return nil, err
			}
			obj.AppendChild(key)

		default:
			val, err := p.close(op)
			if err != nil {
				return nil, err
			}
			obj.LastChild.AppendChild(val)
		}
	}
}

func (p *parseState) closeArray() (*Node, error) {
	arr := &Node{Type: ArrayNode}

	for {
		switch op := p.scan(); op {
		case OpSyntaxError:
			return nil, p.s.LastError()

		case OpArrayEnd:
			return arr, nil

		default:
			elem, err := p.close(op)
			if err != nil {
				return nil, err
			}
			arr.AppendChild(elem)
		}
	}
}

func (p *parseState) closeLiteral() (*Node, error) {
	var typ NodeType
	var start = p.off - 1

	switch p.scan() {
	case OpSyntaxError:
		return nil, p.s.LastError()
	case OpObjectKeyEnd:
		typ = ObjectKeyNode
	case OpStringEnd:
		typ = StringNode
	case OpNumberEnd:
		typ = NumberNode
	case OpBoolEnd:
		typ = BoolNode
	case OpNullEnd:
		typ = NullNode
	default:
		panic("jo: unexpected opcode")
	}

	return &Node{Type: typ, Raw: p.buf[start:p.off]}, nil
}

// scan returns the next significant scanning opcode.
func (p *parseState) scan() Op {
	for p.off < len(p.buf) {
		op, n := p.s.Scan(p.buf[p.off])
		p.off += n

		if op != OpContinue && op != OpSpace {
			return op
		}
	}

	return p.s.Eof()
}
