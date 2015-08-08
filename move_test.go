package main

import "testing"

func TestMoveToLowerRight(t *testing.T) {
	b := NewBoard(2, 3, []Cell{})
	atom := Unit{Members: []Cell{Cell{X: 1, Y: 0}}, Pivot: Cell{X: 1, Y: 0}}
	target := Unit{Members: []Cell{Cell{X: 2, Y: 1}}, Pivot: Cell{X: 2, Y: 1}}

	actual := b.MoveSequence(atom, target)
	expected := []Move{E, SE}

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
	atom := Unit{Members: []Cell{Cell{X: 1, Y: 0}}, Pivot: Cell{X: 1, Y: 0}}
	target := Unit{Members: []Cell{Cell{X: 0, Y: 1}}, Pivot: Cell{X: 0, Y: 1}}

	actual := b.MoveSequence(atom, target)
	expected := []Move{W, SE}

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
	atom := Unit{Members: []Cell{Cell{X: 1, Y: 0}}, Pivot: Cell{X: 1, Y: 0}}
	target := Unit{Members: []Cell{Cell{X: 0, Y: 4}}, Pivot: Cell{X: 0, Y: 4}}

	actual := b.MoveSequence(atom, target)
	expected := []Move{W, SE, SW, SE, SW}

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
			atom:     Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}},
			expected: 1,
		},
		{
			atom:     Unit{Members: []Cell{Cell{0, 0}, Cell{2, 0}}, Pivot: Cell{1, 0}},
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
			atom:     Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}},
			expected: 1,
		},
		{
			atom:     Unit{Members: []Cell{Cell{0, 0}, Cell{2, 2}}, Pivot: Cell{1, 1}},
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
			atom:     Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}},
			expected: Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}},
		},
		{
			board:    NewBoard(2, 3, []Cell{}),
			atom:     Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}},
			expected: Unit{Members: []Cell{Cell{1, 0}}, Pivot: Cell{1, 0}},
		},
		// {
		// 	atom:     Unit{Members: []Cell{Cell{0, 0}, Cell{2, 2}}, Pivot: Cell{1, 1}},
		// 	expected: 3,
		// },
	}

	for _, data := range atoms {
		actual := data.board.StartLocation(data.atom)

		if len(actual.Members) != len(data.expected.Members) {
			t.Errorf("Not the same number of Members: %v expected %v", actual, data.expected)
			return
		}

		for i, member := range data.expected.Members {
			if actual.Members[i] != member {
				t.Errorf("Failed identify height: %v expected %v", actual, data.expected)
			}
		}
	}

}

func TestReadProgram(t *testing.T) {
	sample := `{"id": 23, "units": [], "width": 5, "height": 5, "filled": [], "sourceLength": 2, "sourceSeeds": []}`
	actual := ReadProgram([]byte(sample))
	if actual.Id != 23 || actual.Width != 5 || actual.Height != 5 || actual.SourceLength != 2 {
		t.Errorf("Failed to read program got: %v", actual)
	}

	sample = `{"id": 21, "units": [], "width": 10, "height": 5, "filled": [], "sourceLength": 2, "sourceSeeds": []}`
	actual = ReadProgram([]byte(sample))
	if actual.Id != 21 || actual.Width != 10 || actual.Height != 5 || actual.SourceLength != 2 {
		t.Errorf("Failed to read program got: %v", actual)
	}

	sample = `{"id": 21, "units": [{"members": [{"x": 1, "y": 1}], "pivot": {"x": 1, "y": 1}}], "width": 10, "height": 5, "filled": [{"x": 1, "y": 2}], "sourceLength": 2, "sourceSeeds": []}`
	actual = ReadProgram([]byte(sample))
	if len(actual.Units) != 1 || len(actual.Units[0].Members) != 1 || len(actual.Filled) != 1 {
		t.Errorf("Failed to read program got: %v", actual)
	}

}

func TestFillBoard(t *testing.T) {
	b := NewBoard(2, 2, []Cell{Cell{X: 1, Y: 1}})
	actual := b.FillCells([]Cell{Cell{X: 0, Y: 0}})

	if !actual[0][0] || !actual[1][1] || actual[0][1] || actual[1][0] {
		t.Errorf("Failed to read fill board got: %v", actual)
	}

}
