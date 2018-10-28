package kdtree

import (
	"image/color"
)

type Kdtree struct {
	root *node
}

type node struct {
	rgb         [3]int
	left, right *node
}

func New() *Kdtree {
	return &Kdtree{}
}

// Adds the color to the tree. Panics if t is nil.
func (t *Kdtree) Add(c color.Color) {
	if t == nil {
		panic("cannot Add() to nil Kdtree!")
	}

	rgb := unpackColor(c)
	if t.root == nil {
		t.root = &node{rgb: rgb}
	} else {
		t.root.add(rgb, 0)
	}
}

func (n *node) add(rgb [3]int, splitAxis int) {
	next := &n.left
	if rgb[splitAxis] > n.rgb[splitAxis] {
		next = &n.right
	}

	if *next == nil {
		*next = &node{rgb: rgb}
	} else {
		(*next).add(rgb, (splitAxis+1)%3)
	}
}

// Returns the nearest color to tgt in the tree. If tree is empty, returns black.
func (t *Kdtree) Nearest(tgt color.Color) color.Color {
	if t == nil || t.root == nil {
		return color.RGBA{}
	}

	n, _ := t.root.nearest(unpackColor(tgt), 0)
	return color.RGBA{uint8(n.rgb[0]), uint8(n.rgb[1]), uint8(n.rgb[2]), 0xff}
}

func (n *node) nearest(tgt [3]int, splitAxis int) (*node, int) {
	if n == nil {
		return nil, -1
	}

	// Easy case: exact match at this node.
	curDist := sqDist(tgt, n.rgb)
	if curDist == 0 {
		return n, 0
	}

	// Choose next subtree to search based on splitting axis RGB component.
	n1, n2 := n.left, n.right
	if tgt[splitAxis] > n.rgb[splitAxis] {
		n1, n2 = n2, n1
	}

	bestNode, bestDist := n1.nearest(tgt, (splitAxis+1)%3)
	if bestNode == nil || bestDist > (tgt[splitAxis]-n.rgb[splitAxis])*(tgt[splitAxis]-n.rgb[splitAxis]) {
		// Try the other subtree if there could be a closer node on that side of tree.
		r, d := n2.nearest(tgt, (splitAxis+1)%3)
		if r != nil && (bestNode == nil || d < bestDist) {
			bestNode = r
			bestDist = d
		}
	}

	// Compare best result from subtrees to current node.
	if bestNode == nil || curDist < bestDist {
		bestNode = n
		bestDist = curDist
	}
	return bestNode, bestDist
}

func unpackColor(c color.Color) [3]int {
	r, g, b, _ := c.RGBA()
	return [3]int{int(r >> 8), int(g >> 8), int(b >> 8)}
}

func sqDist(c1, c2 [3]int) int {
	d := 0
	for i := 0; i < 3; i++ {
		d += (c2[i] - c1[i]) * (c2[i] - c1[i])
	}
	return d
}
