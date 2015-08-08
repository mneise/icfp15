package main

import "testing"

func TestMoveToLowerRight(t *testing.T) {
	b := NewBoard(2, 3, []Cell{})
	atom := Unit{members: []Cell{Cell{x: 1, y: 0}}, pivot: Cell{x: 1, y: 0}}
	target := Unit{members: []Cell{Cell{x: 2, y: 1}}, pivot: Cell{x: 2, y: 1}}

	actual := b.MoveSequence(atom, target)
	expected := []Command{E, SE}

	if len(actual) != len(expected) {
		t.Errorf("Not the same amount of moves: %v expected %v", actual, expected)
		return
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Failed to move to lower right, got moves: %v expected %v", actual, expected)
		}
	}
}

func TestMoveToLowerLeft(t *testing.T) {
	b := NewBoard(2, 3, []Cell{})
	atom := Unit{members: []Cell{Cell{x: 1, y: 0}}, pivot: Cell{x: 1, y: 0}}
	target := Unit{members: []Cell{Cell{x: 0, y: 1}}, pivot: Cell{x: 0, y: 1}}

	actual := b.MoveSequence(atom, target)
	expected := []Command{W, SE}

	if len(actual) != len(expected) {
		t.Errorf("Not the same amount of moves: %v expected %v", actual, expected)
		return
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Failed to move to lower right, got moves: %v expected %v", actual, expected)
		}
	}
}

func TestMoveFurtherToLowerLeft(t *testing.T) {
	b := NewBoard(5, 3, []Cell{})
	atom := Unit{members: []Cell{Cell{x: 1, y: 0}}, pivot: Cell{x: 1, y: 0}}
	target := Unit{members: []Cell{Cell{x: 0, y: 4}}, pivot: Cell{x: 0, y: 4}}

	actual := b.MoveSequence(atom, target)
	expected := []Command{W, SE, SW, SE, SW}

	if len(actual) != len(expected) {
		t.Errorf("Not the same amount of moves: %v expected %v", actual, expected)
		return
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Failed to move to lower right, got moves: %v expected %v", actual, expected)
		}
	}
}

func TestUnitWidth(t *testing.T) {
	atoms := []struct {
		atom     Unit
		expected int
	}{
		{
			atom:     Unit{members: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}},
			expected: 1,
		},
		{
			atom:     Unit{members: []Cell{Cell{0, 0}, Cell{2, 0}}, pivot: Cell{1, 0}},
			expected: 3,
		},
	}

	for _, data := range atoms {
		actual := data.atom.Width()

		if actual != data.expected {
			t.Errorf("Failed identify width: %v expected %v", actual, data.expected)
		}
	}

}

func TestUnitHeight(t *testing.T) {
	atoms := []struct {
		atom     Unit
		expected int
	}{
		{
			atom:     Unit{members: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}},
			expected: 1,
		},
		{
			atom:     Unit{members: []Cell{Cell{0, 0}, Cell{2, 2}}, pivot: Cell{1, 1}},
			expected: 3,
		},
	}

	for _, data := range atoms {
		actual := data.atom.Height()

		if actual != data.expected {
			t.Errorf("Failed identify height: %v expected %v", actual, data.expected)
		}
	}

}

func TestStartLocation(t *testing.T) {
	atoms := []struct {
		board    Board
		atom     Unit
		expected Unit
	}{
		{
			board:    NewBoard(2, 2, []Cell{}),
			atom:     Unit{members: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}},
			expected: Unit{members: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}},
		},
		{
			board:    NewBoard(2, 3, []Cell{}),
			atom:     Unit{members: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}},
			expected: Unit{members: []Cell{Cell{1, 0}}, pivot: Cell{1, 0}},
		},
		// {
		// 	atom:     Unit{members: []Cell{Cell{0, 0}, Cell{2, 2}}, pivot: Cell{1, 1}},
		// 	expected: 3,
		// },
	}

	for _, data := range atoms {
		actual := data.board.StartLocation(data.atom)

		if len(actual.members) != len(data.expected.members) {
			t.Errorf("Not the same number of members: %v expected %v", actual, data.expected)
			return
		}

		for i, member := range data.expected.members {
			if actual.members[i] != member {
				t.Errorf("Failed identify height: %v expected %v", actual, data.expected)
			}
		}
	}

}
