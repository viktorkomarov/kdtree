package kdtree

import (
	"fmt"
	"sort"
)

func disc(level, k int) int {
	return level % k
}

type Node struct {
	left  *Node
	right *Node
	coord []int
}

type Root Node

func sortByKD(values [][]int, k, d int) [][]int {
	sort.Slice(values, func(i, j int) bool {
		d := disc(d, k)
		return values[i][d] < values[j][d]
	})

	return values
}

func NewKDTree(k int, values ...[]int) (Root, error) {
	for i := range values {
		if len(values[i]) != k {
			return Root{}, fmt.Errorf("unexpected dimension: %d != %d", len(values[i]), k)
		}
	}

	val := sortByKD(values, k, 0)
	md := len(val) / 2
	r := Root{
		coord: val[md],
	}

	r.left, r.right = buildChildren(k, 1, values[:md], values[md+1:])
	return r, nil
}

func buildChildren(k, d int, leftValues, rightValues [][]int) (*Node, *Node) {
	var left *Node
	if len(leftValues) != 0 {
		coord := sortByKD(leftValues, k, d)
		md := len(coord)
		left = &Node{
			coord: coord[md],
		}
		left.left, left.right = buildChildren(k, d+1, coord[:md], coord[md+1:])
	}

	var right *Node
	if len(rightValues) != 0 {
		coord := sortByKD(rightValues, k, d)
		md := len(coord)
		right = &Node{
			coord: coord[md],
		}
		right.left, right.right = buildChildren(k, d+1, coord[:md], coord[md+1:])
	}

	return left, right
}

func totalEqualCoord(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func exactWalk(node *Node, d int, k []int) *Node {
	if node == nil {
		return nil
	}

	if totalEqualCoord(node.coord, k) {
		return node
	}

	idx := disc(d, len(k))
	if node.coord[idx] < k[idx] {
		return exactWalk(node.left, d+1, k)
	}

	return exactWalk(node.right, d+1, k)
}

func ExactMatch(root Root, k []int) *Node {
	n := Node(root)
	return exactWalk(&n, 0, k)
}

func eqCoordMap(coord []int, keys map[int]int) bool {
	totalMatch := 0
	for i, val := range coord {
		if v, ok := keys[i]; ok && v == val {
			totalMatch++
		}
	}

	return totalMatch == len(keys)
}

// TODO :: check for slice bug
func multipleWalk(node *Node, d int, keys map[int]int, set [][]int) [][]int {
	if node == nil {
		return nil
	}

	if eqCoordMap(node.coord, keys) {
		set = append(set, node.coord)
	}

	idx := disc(d, len(node.coord))
	c, ok := keys[idx]
	switch {
	case !ok || c == node.coord[idx]:
		set = multipleWalk(node.left, d+1, keys, set)
		set = multipleWalk(node.right, d+1, keys, set)
	case c < node.coord[idx]:
		set = multipleWalk(node.left, d+1, keys, set)
	case c > node.coord[idx]:
		set = multipleWalk(node.right, d+1, keys, set)
	}

	return set
}

func MultipleMatch(root Root, keys map[int]int) [][]int {
	node := Node(root)
	var result [][]int
	return multipleWalk(&node, 0, keys, result)
}
