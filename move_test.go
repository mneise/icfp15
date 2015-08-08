package main

import "testing"

func TestMoveToLowerRight(t *testing.T) {

	b := NewBoard(2, 2, []Cell{})
	atom := Unit{members: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}}
	target := Unit{members: []Cell{Cell{1, 1}}, pivot: Cell{1, 1}}
	actual := MoveToTarget(b, atom, target)
	expected := []Command{E, SE}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Failed to move to lower right, got moves: %v expected %v", actual, expected)
		}
	}
}
