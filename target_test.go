package main

import "testing"

func TestFindTargetLowerRight(t *testing.T) {
	board := NewBoard(2, 2, []Cell{Cell{1, 0}})
	unit := Unit{cells: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}}
	actual := TargetLocation(board, unit)
	expected := BoardUnit{cells: []Cell{Cell{1, 1}}, pivot: Cell{1, 1}}

	for i := range expected.cells {
		if actual.cells[i] != expected.cells[i] {
			t.Errorf("Failed to find target, got cell: %v expected %v", actual, expected)
		}
	}
}
