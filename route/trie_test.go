package route

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie(t *testing.T) {
	root := newNode("")

	root.insert("")
	root.insert("www.luoji.com/hello")
	root.insert("www.luoji.com/hello/m1")
	root.insert("www.luoji.com/hello/m1/")
	root.insert("www.luoji.com/hello/m1/sub1")
	root.insert("www.luoji.com/world")
	root.insert("www.luoji.com/v1/A/hello")

	root.insert("s2.luoji.com/hello")
	tprint(t, root)

	p := root.find("www.luoji.com/hello/m1/test/test")
	assert.Equal(t, "www.luoji.com/hello/m1", p)

	p = root.find("www.luoji.com/hello/test/test")
	assert.Equal(t, "www.luoji.com/hello", p)

	p = root.find("www.luoji.com/hello/")
	assert.Equal(t, "www.luoji.com/hello", p)

	p = root.find("www.luoji.com/hello/m1/sub1/test/test")
	assert.Equal(t, "www.luoji.com/hello/m1/sub1", p)

	p = root.find("www.luoji.com/world")
	assert.Equal(t, "www.luoji.com/world", p)

	p = root.find("www.luoji.com/v1/A/hello")
	assert.Equal(t, "www.luoji.com/v1/A/hello", p)

	p = root.find("www.luoji.com/not/exits")
	assert.Equal(t, "www.luoji.com/", p)

	p = root.find("www.luoji.com/")
	assert.Equal(t, "www.luoji.com/", p)

	p = root.find("s2.luoji.com/hello/sdfsdf")
	assert.Equal(t, "s2.luoji.com/hello", p)
}

func tprint(t *testing.T, n *node) {
	if n.parent != nil {
		t.Logf("%s -> %s", n.parent.name, n.name)
	}
	for _, item := range n.children {
		tprint(t, item)
	}
}
