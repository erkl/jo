package jo

// A Node is a part of a tree describing a JSON value. It contains
// information about the node's type, its original JSON representation
// (for literal values), and information about its position inside
// the tree.
//
// The tree representation of the following JSON value
//
//     { "foo": [123, 456], "bar": false }
//
// looks like this:
//
//     ObjectNode
//     ├─ ObjectKeyNode ("foo")
//     │  └─ ArrayNode
//     │     ├─ NumberNode (123)
//     │     └─ NumberNode (456)
//     └─ ObjectKeyNode ("bar")
//        └─ BoolNode (false)
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
		panic("AppendChild called with an already attached Node")
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

// RemoveChild removes a child from the subject Node.
func (n *Node) RemoveChild(c *Node) {
	if c.Parent != n {
		panic("RemoveChild called for non-child Node")
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
