package main

import "testing"

func equalsUnit(actual Unit, expected Unit) bool {
	for i := range expected.Members {
		if actual.Members[i] != expected.Members[i] {
			return false
		}
	}

	if expected.Pivot != actual.Pivot {
		return false
	}

	return true
}

func TestFindTargetLowerRight(t *testing.T) {
	board := NewBoard(2, 2, []Cell{Cell{1, 0}})
	unit := Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}}
	actual := TargetLocation(board, unit)
	expected := Unit{Members: []Cell{Cell{1, 1}}, Pivot: Cell{1, 1}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}

	board = NewBoard(2, 2, []Cell{
		Cell{0, 0}, Cell{1, 0},
		Cell{1, 1}})
	unit = Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}}
	actual = TargetLocation(board, unit)
	expected = Unit{Members: []Cell{Cell{0, 1}}, Pivot: Cell{0, 1}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}

	board = NewBoard(3, 3, []Cell{
		Cell{0, 0}, Cell{1, 0}, Cell{2, 0},
		Cell{0, 1}, Cell{2, 1},
		Cell{0, 2}, Cell{1, 2}, Cell{2, 2}})
	unit = Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}}
	actual = TargetLocation(board, unit)
	expected = Unit{Members: []Cell{Cell{1, 1}}, Pivot: Cell{1, 1}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}

	board = NewBoard(3, 3, []Cell{
		Cell{0, 0}, Cell{1, 0}, Cell{2, 0},
		Cell{0, 1},
		Cell{0, 2}, Cell{1, 2}})
	unit = Unit{Members: []Cell{Cell{0, 0}, Cell{1, 0}}, Pivot: Cell{0, 0}}
	actual = TargetLocation(board, unit)
	expected = Unit{Members: []Cell{Cell{1, 1}, Cell{2, 1}}, Pivot: Cell{1, 1}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}
}

func TestUnitRelativeToCell(t *testing.T) {
	unit := Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}}
	cell := Cell{1, 1}
	actual := unit.MoveTo(cell)
	expected := Unit{Members: []Cell{Cell{1, 1}}, Pivot: Cell{1, 1}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to get relative cell, got cell: %v expected %v", actual, expected)
	}

	unit = Unit{Members: []Cell{Cell{0, 0}, Cell{2, 0}}, Pivot: Cell{1, 0}}
	cell = Cell{0, 3}
	actual = unit.MoveTo(cell)
	expected = Unit{Members: []Cell{Cell{-1, 3}, Cell{1, 3}}, Pivot: Cell{0, 3}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to get relative cell, got cell: %v expected %v", actual, expected)
	}
}

func TestIsValidUnit(t *testing.T) {
	unit := Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}}
	board := NewBoard(2, 2, []Cell{})

	if !unit.isValid(board) {
		t.Errorf("Expected unit: %v to be valid on board %v, but was invalid", unit, board)
	}

	unit = Unit{Members: []Cell{Cell{0, 0}, Cell{1, 0}}, Pivot: Cell{0, 0}}
	board = NewBoard(2, 2, []Cell{Cell{1, 0}})

	if unit.isValid(board) {
		t.Errorf("Expected unit: %v to be invalid on board %v, but was valid", unit, board)
	}

	unit = Unit{Members: []Cell{Cell{0, 0}, Cell{-1, 0}}, Pivot: Cell{0, 0}}
	board = NewBoard(2, 2, []Cell{})

	if unit.isValid(board) {
		t.Errorf("Expected unit: %v to be invalid on board %v, but was valid", unit, board)
	}
}
