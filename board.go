package main

import "fmt"
import "math"

func main() {
	cs := []Cell{Cell{0, 0}, Cell{1, 1}}
	fmt.Printf("hello, world %v\n", NewBoard(3, 2, cs))
}

type Board [][]bool
type Cell struct {
	x int
	y int
}
type Unit struct {
	members []Cell
	pivot   Cell
}
type Command int

const (
	E Command = iota
	W
	SE
	SW
	RC
	RCC
)

func NewBoard(rows int, cols int, cells []Cell) Board {
	b := make([][]bool, rows)
	for i := range b {
		b[i] = make([]bool, cols)
	}

	for _, c := range cells {
		b[c.y][c.x] = true
	}

	return b
}

func (u Unit) MoveTo(cell Cell) Unit {
	unit := Unit{members: make([]Cell, len(u.members)), pivot: cell}
	for i, member := range u.members {
		x := cell.x + (member.x - u.pivot.x)
		y := cell.y + (member.y - u.pivot.y)
		unit.members[i] = Cell{y: y, x: x}
	}
	return unit
}

func TargetLocation(board Board, unit Unit) Unit {
	// for y := range board {
	// 	for x := range board[y] {
	// 		for cell := range unit.members {

	// 		}
	// 	}
	// }

	return Unit{[]Cell{Cell{1, 1}}, Cell{1, 1}}
}

func (b Board) Width() int {
	if len(b) > 0 {
		return len(b[0])
	}

	return 0
}

func (b Board) Height() int {
	return len(b)
}

func MoveToTarget(board Board, unit Unit, target Unit) []Command {

	// get start location

	// move left / right

	// move down

	return []Command{E, SE}
}

func (c Cell) ShiftX(offset int) Cell {
	return Cell{x: c.x + offset, y: c.y}
}

func (b Board) StartLocation(u Unit) Unit {
	offset := (b.Width() - u.Width()) / 2
	newPivot := u.pivot.ShiftX(offset)

	return u.MoveTo(newPivot)
}

func (u Unit) Width() int {
	minX := math.MaxInt32
	maxX := -1

	for _, member := range u.members {
		if minX > member.x {
			minX = member.x
		}
		if maxX < member.x {
			maxX = member.x
		}
	}

	return 1 + maxX - minX
}

func (u Unit) Height() int {
	minY := math.MaxInt32
	maxY := -1

	for _, member := range u.members {
		if minY > member.y {
			minY = member.y
		}
		if maxY < member.y {
			maxY = member.y
		}
	}

	return 1 + maxY - minY
}
