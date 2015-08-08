package main

import "fmt"

func main() {
	cs := []Cell{Cell{0, 0}, Cell{1, 1}}
	fmt.Printf("hello, world %v\n", NewBoard(3, 2, cs))
}

type Board [][]bool
type Cell struct {
	row int
	col int
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
		b[c.row][c.col] = true
	}

	return b
}

func UnitRelativeToCell(unit Unit, cell Cell) Unit {
	newUnit := Unit{members: make([]Cell, len(unit.members)), pivot: cell}
	for i, member := range unit.members {
		newRow := cell.row + (member.row - unit.pivot.row)
		newCol := cell.col + (member.col - unit.pivot.col)
		newUnit.members[i] = Cell{newRow, newCol}
	}
	return newUnit
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

func MoveToTarget(board Board, unit Unit, target Unit) []Command {
	return []Command{E, SE}
}
