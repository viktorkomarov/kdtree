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
