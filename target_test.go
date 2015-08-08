package main

import "testing"

func equalsUnit(actual Unit, expected Unit) bool {
	for i := range expected.members {
		if actual.members[i] != expected.members[i] {
			return false
		}
	}

	if expected.pivot != actual.pivot {
		return false
	}

	return true
}

func TestFindTargetLowerRight(t *testing.T) {
	board := NewBoard(2, 2, []Cell{Cell{1, 0}})
	unit := Unit{members: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}}
	actual := TargetLocation(board, unit)
	expected := Unit{members: []Cell{Cell{1, 1}}, pivot: Cell{1, 1}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}
}

func TestUnitRelativeToCell(t *testing.T) {
	unit := Unit{members: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}}
	cell := Cell{1, 1}
	actual := UnitRelativeToCell(unit, cell)
	expected := Unit{members: []Cell{Cell{1, 1}}, pivot: Cell{1, 1}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to get relative cell, got cell: %v expected %v", actual, expected)
	}

	unit = Unit{members: []Cell{Cell{0, 0}, Cell{2, 0}}, pivot: Cell{1, 0}}
	cell = Cell{0, 3}
	actual = UnitRelativeToCell(unit, cell)
	expected = Unit{members: []Cell{Cell{-1, 3}, Cell{1, 3}}, pivot: Cell{0, 3}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to get relative cell, got cell: %v expected %v", actual, expected)
	}
}
