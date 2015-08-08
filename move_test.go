package main

import "testing"

func TestMoveToLowerRight(t *testing.T) {
	b := NewBoard(2, 3, []Cell{})
	atom := Unit{members: []Cell{Cell{0, 0}}, pivot: Cell{0, 0}}
	target := BoardUnit{members: []Cell{Cell{1, 2}}, pivot: Cell{1, 2}}
	actual := MoveToTarget(b, atom, target)
	expected := []Command{E, SE}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Failed to move to lower right, got moves: %v expected %v", actual, expected)
		}
	}
}

// func TestMoveToLowerLeft(t *testing.T) {
//   b := NewBoard(2, 2, []Cell{})
//   atom := Unit{members: []Cell{Cell{0,0}}, pivot: Cell{0,0}}
//   target := BoardUnit{members: []Cell{Cell{1,0}}, pivot: Cell{1,0}}
//   actual := MoveToTarget(b, atom, target)
//   expected := []Command{W, SE}

//   for i := range expected {
//     if actual[i] != expected[i] {
//       t.Errorf("Failed to move to lower right, got moves: %v expected %v", actual, expected)
//     }
//   }
// }

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
			atom:     Unit{members: []Cell{Cell{0, 0}, Cell{0, 2}}, pivot: Cell{0, 1}},
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
