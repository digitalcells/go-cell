package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

type node struct {
	isLast   bool
	segment  string
	handler  ControllerHandler
	children []*node
}

func NewNode() *node {
	return &node{
		isLast:   false,
		segment:  "",
		children: []*node{},
	}
}

func NewTree() *Tree {
	root := NewNode()
	return &Tree{root: root}
}

func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

func (n *node) filterChildNodes(segment string) []*node {
	if len(n.children) == 0 {
		return nil
	}

	if isWildSegment(segment) {
		return n.children
	}

	nodes := make([]*node, 0, len(n.children))

	for _, inode := range n.children {
		if isWildSegment(inode.segment) {
			nodes = append(nodes, inode)
			continue
		}

		if inode.segment == segment {
			nodes = append(nodes, inode)
			continue
		}
	}

	return nodes
}

func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]

	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}

	inodes := n.filterChildNodes(segment)

	if inodes == nil || len(inodes) == 0 {
		return nil
	}

	if len(segments) == 1 {
		for _, n := range inodes {
			if n.isLast {
				return n
			}
		}

		return nil
	}

	for _, n := range inodes {
		_match := n.matchNode(segments[1])
		if _match != nil {
			return _match
		}
	}

	return nil
}

func (tree *Tree) AddRouter(uri string, handler ControllerHandler) error {
	n := tree.root

	if n.matchNode(uri) != nil {
		return errors.New("route exist: " + uri)
	}

	segments := strings.Split(uri, "/")

	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}

		isLast := index == len(segments)-1

		var object *node

		childNodes := n.filterChildNodes(segment)

		if len(childNodes) > 0 {
			for _, inode := range childNodes {
				if inode.segment == segment {
					object = inode
					break
				}
			}
		}

		if object == nil {
			temp := NewNode()
			temp.segment = segment

			if isLast {
				temp.isLast = true
				temp.handler = handler
			}

			n.children = append(n.children, temp)

			object = temp
		}

		n = object
	}

	return nil
}

func (tree *Tree) FindHandler(uri string) ControllerHandler {
	_match := tree.root.matchNode(uri)

	if _match == nil {
		return nil
	}

	return _match.handler
}
