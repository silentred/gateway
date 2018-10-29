package route

import (
	"strings"
)

const (
	sep = "/"
)

type node struct {
	name     string
	parent   *node
	children children
}

type children []*node

func newNode(name string) *node {
	return &node{name: name}
}

// insert path prefix, for example: www.luoji.com/hello/m1
func (n *node) insert(path string) {
	var parts []string
	var currNode *node
	currNode = n

	// remove trailing /
	if strings.HasSuffix(path, sep) {
		path = path[:len(path)-1]
	}
	// root path
	if path == "" {
		return
	}

	parts = strings.Split(path, sep)

	for len(parts) > 0 {
		var found bool
		// create children for current node
		if currNode.children == nil {
			currNode.children = make(children, 0)
		}

		for i := 0; i < len(currNode.children) && !found; i++ {
			// found matched child
			if currNode.children[i].name == parts[0] {
				parts = parts[1:]
				currNode = currNode.children[i]
				found = true
			}
		}

		if !found {
			// not found matched child
			c := newNode(parts[0])
			c.parent = currNode
			currNode.children = append(currNode.children, c)
			currNode = c
			parts = parts[1:]
		}
	}
}

// find the prefix which is mostly matched in the tree
func (n *node) find(fullpath string) (routeStr string) {
	var matched []*node
	var parts []string
	var currNode *node
	if strings.HasPrefix(fullpath, sep) {
		fullpath = fullpath[1:]
	}
	parts = strings.Split(fullpath, sep)

	currNode = n

	for currNode != nil && currNode.children != nil && len(parts) > 0 {
		var found bool
		for i := 0; i < len(currNode.children) && !found; i++ {
			// match one node
			if currNode.children[i].name == parts[0] {
				matched = append(matched, currNode.children[i])
				currNode = currNode.children[i]
				parts = parts[1:]
				found = true
			}
		}
		// no match
		if !found {
			break
		}
	}

	// only matched domain
	if len(matched) == 1 {
		routeStr = matched[0].name + "/"
	} else if len(matched) > 1 {
		var matchedStrings = make([]string, len(matched))
		for i := 0; i < len(matched); i++ {
			matchedStrings[i] = matched[i].name
		}
		routeStr = strings.Join(matchedStrings, sep)
	}

	return
}
