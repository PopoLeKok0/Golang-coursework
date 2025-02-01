package main

import "fmt"

type Point struct {
	x int
	y int
}

type PtTree struct {
	pt    Point
	left  *PtTree
	right *PtTree
}

func (t *PtTree) postOrder() {
	if t.left != nil {
		t.left.postOrder()
	}
	if t.right != nil {
		t.right.postOrder()
	}
	fmt.Printf("(%d,%d)", t.pt.x, t.pt.y)
}

func (t *PtTree) Find(x, y int) bool {
	if t.pt.x == x && t.pt.y == y {
		return true
	}
	if t.left != nil && t.left.Find(x, y) {
		return true
	}
	if t.right != nil && t.right.Find(x, y) {
		return true
	}
	return false
}

type PointSearcher interface {
	Find(int, int) bool
}

func main() {
	tree := PtTree{Point{2, 3},
		&PtTree{Point{5, 1},
			&PtTree{Point{2, 2}, nil, nil},
			&PtTree{Point{8, 3}, &PtTree{Point{1, 6}, nil, nil}, nil}},
		&PtTree{Point{4, 7},
			&PtTree{Point{7, 2},
				&PtTree{Point{6, 4}, nil, nil},
				&PtTree{Point{0, 9}, nil, nil}},
			&PtTree{Point{3, 6}, nil, nil}}}

	tree.postOrder()
	fmt.Println("")

	// Utilisation de l'interface
	var ps PointSearcher = &tree

	u, v := 7, 2
	if ps.Find(u, v) {
		fmt.Printf("Found: %d %d \n", u, v)
	} else {
		fmt.Printf("Not Found\n")
	}

	x, y := 8, 6
	if ps.Find(x, y) {
		fmt.Printf("Found: %d %d \n", x, y)
	} else {
		fmt.Printf("Not Found\n")
	}
}
