package jo

import (
	"runtime"
	"testing"
)

func TestNodeAppendChild(t *testing.T) {
	r := &Node{Type: ArrayNode}
	a := &Node{Type: NumberNode, Raw: []byte(`123`)}
	b := &Node{Type: StringNode, Raw: []byte(`"hello"`)}
	c := &Node{Type: BoolNode, Raw: []byte(`true`)}

	// add all child nodes in order
	r.AppendChild(a)
	assert(t, a.Parent == r)
	assert(t, a.PrevSibling == nil)
	assert(t, a.NextSibling == nil)
	assert(t, r.FirstChild == a)
	assert(t, r.LastChild == a)

	r.AppendChild(b)
	assert(t, a.PrevSibling == nil)
	assert(t, a.NextSibling == b)
	assert(t, b.PrevSibling == a)
	assert(t, b.NextSibling == nil)
	assert(t, r.FirstChild == a)
	assert(t, r.LastChild == b)

	r.AppendChild(c)
	assert(t, a.PrevSibling == nil)
	assert(t, a.NextSibling == b)
	assert(t, b.PrevSibling == a)
	assert(t, b.NextSibling == c)
	assert(t, c.PrevSibling == b)
	assert(t, c.NextSibling == nil)
	assert(t, r.FirstChild == a)
	assert(t, r.LastChild == c)
}

func TestNodeInsertChild(t *testing.T) {
	r := &Node{Type: ArrayNode}
	a := &Node{Type: NumberNode, Raw: []byte(`123`)}
	b := &Node{Type: StringNode, Raw: []byte(`"hello"`)}
	c := &Node{Type: BoolNode, Raw: []byte(`true`)}

	// insert the children in order, before what is currently
	// the root node's last child
	r.InsertChild(a, r.LastChild)
	assert(t, a.Parent == r)
	assert(t, a.PrevSibling == nil)
	assert(t, a.NextSibling == nil)
	assert(t, r.FirstChild == a)
	assert(t, r.LastChild == a)

	r.InsertChild(b, r.LastChild)
	assert(t, b.PrevSibling == nil)
	assert(t, b.NextSibling == a)
	assert(t, a.PrevSibling == b)
	assert(t, a.NextSibling == nil)
	assert(t, r.FirstChild == b)
	assert(t, r.LastChild == a)

	r.InsertChild(c, r.LastChild)
	assert(t, b.PrevSibling == nil)
	assert(t, b.NextSibling == c)
	assert(t, c.PrevSibling == b)
	assert(t, c.NextSibling == a)
	assert(t, a.PrevSibling == c)
	assert(t, a.NextSibling == nil)
	assert(t, r.FirstChild == b)
	assert(t, r.LastChild == a)
}

func TestNodeRemoveChild(t *testing.T) {
	r := &Node{Type: ArrayNode}
	a := &Node{Type: NumberNode, Raw: []byte(`123`)}
	b := &Node{Type: StringNode, Raw: []byte(`"hello"`)}
	c := &Node{Type: BoolNode, Raw: []byte(`true`)}

	r.FirstChild = a
	r.LastChild = c

	a.Parent = r
	a.NextSibling = b

	b.Parent = r
	b.PrevSibling = a
	b.NextSibling = c

	c.Parent = r
	c.PrevSibling = b

	// remove the children one at a time
	r.RemoveChild(b)
	assert(t, b.Parent == nil)
	assert(t, b.PrevSibling == nil)
	assert(t, b.NextSibling == nil)
	assert(t, a.NextSibling == c)
	assert(t, c.PrevSibling == a)

	r.RemoveChild(a)
	assert(t, c.PrevSibling == nil)
	assert(t, c.NextSibling == nil)
	assert(t, r.FirstChild == c)
	assert(t, r.LastChild == c)

	r.RemoveChild(c)
	assert(t, r.FirstChild == nil)
	assert(t, r.LastChild == nil)
}

func assert(t *testing.T, cond bool) {
	if !cond {
		if _, _, line, ok := runtime.Caller(1); ok {
			t.Errorf("assertion failed on line %d", line)
		} else {
			t.Errorf("assertion failed on line ???")
		}
	}
}
