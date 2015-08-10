package main

import "testing"

func equalsUnit(actual Unit, expected Unit) bool {
	if actual.Height() != expected.Height() ||
		actual.Width() != expected.Width() {
		return false
	}

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

func equalsBoard(actual Board, expected Board) bool {
	for y := range expected {
		for x := range expected[y] {
			if expected[y][x] != actual[y][x] {
				return false
			}
		}
	}

	return true
}

func TestFindTargetLowerRight(t *testing.T) {
	board := NewBoard(2, 2, []Cell{Cell{1, 0}})
	unit := Unit{[]Cell{Cell{0, 0}}, Cell{0, 0}}
	actual := TargetLocations(board, unit)
	expected := []Unit{
		Unit{[]Cell{Cell{0, 0}}, Cell{0, 0}},
		Unit{[]Cell{Cell{1, 1}}, Cell{1, 1}},
		Unit{[]Cell{Cell{0, 1}}, Cell{0, 1}}}

	if len(expected) != len(actual) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}

	for i := range expected {
		if !equalsUnit(actual[i], expected[i]) {
			t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
		}
	}

	board = NewBoard(2, 2, []Cell{
		Cell{0, 0}, Cell{1, 0},
		Cell{1, 1}})
	unit = Unit{[]Cell{Cell{0, 0}}, Cell{0, 0}}
	actual = TargetLocations(board, unit)
	expected = []Unit{Unit{[]Cell{Cell{0, 1}}, Cell{0, 1}}}

	if len(expected) != len(actual) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}

	for i := range expected {
		if !equalsUnit(actual[i], expected[i]) {
			t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
		}
	}

	board = NewBoard(3, 3, []Cell{
		Cell{0, 0}, Cell{1, 0}, Cell{2, 0},
		Cell{0, 1}, Cell{2, 1},
		Cell{0, 2}, Cell{1, 2}, Cell{2, 2}})
	unit = Unit{[]Cell{Cell{0, 0}}, Cell{0, 0}}
	actual = TargetLocations(board, unit)
	expected = []Unit{
		Unit{[]Cell{Cell{1, 1}}, Cell{1, 1}}}

	if len(expected) != len(actual) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}

	for i := range expected {
		if !equalsUnit(actual[i], expected[i]) {
			t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
		}
	}

	board = NewBoard(3, 3, []Cell{
		Cell{0, 0}, Cell{1, 0}, Cell{2, 0},
		Cell{0, 1},
		Cell{0, 2}, Cell{1, 2}})
	unit = Unit{[]Cell{Cell{0, 0}, Cell{1, 0}}, Cell{0, 0}}
	actual = TargetLocations(board, unit)
	expected = []Unit{
		Unit{[]Cell{Cell{1, 1}, Cell{2, 1}}, Cell{1, 1}}}

	if len(expected) != len(actual) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}

	for i := range expected {
		if !equalsUnit(actual[i], expected[i]) {
			t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
		}
	}

	board = NewBoard(3, 3, []Cell{
		Cell{0, 0}, Cell{1, 0}, Cell{2, 0},
		Cell{0, 1},
		Cell{0, 2}, Cell{1, 2}})
	unit = Unit{[]Cell{Cell{0, 0}}, Cell{0, 0}}
	actual = TargetLocations(board, unit)
	expected = []Unit{
		Unit{[]Cell{Cell{2, 2}}, Cell{2, 2}},
		Unit{[]Cell{Cell{2, 1}}, Cell{2, 1}},
		Unit{[]Cell{Cell{1, 1}}, Cell{1, 1}}}

	if len(expected) != len(actual) {
		t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
	}

	for i := range expected {
		if !equalsUnit(actual[i], expected[i]) {
			t.Errorf("Failed to find target, got unit: %v expected %v", actual, expected)
		}
	}
}

func TestUnitRelativeToCell(t *testing.T) {
	unit := Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}}
	cell := Cell{1, 1}
	actual := unit.MoveTo(cell, unit.Pivot)
	expected := Unit{Members: []Cell{Cell{1, 1}}, Pivot: Cell{1, 1}}

	if !equalsUnit(actual, expected) {
		t.Errorf("Failed to get relative cell, got cell: %v expected %v", actual, expected)
	}

	unit = Unit{Members: []Cell{Cell{0, 0}, Cell{2, 0}}, Pivot: Cell{1, 0}}
	cell = Cell{0, 3}
	actual = unit.MoveTo(cell, unit.Pivot)
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

func TestMovesToCommands(t *testing.T) {
	cs := []Move{E, SE, RC}
	expected := []string{"b", "l", "d"}
	actual := MovesToCommands(cs)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("Expected command sequence to be: %v but was: %v", expected, actual)
		}
	}
}

func TestCalcRandom(t *testing.T) {
	s := 17
	l := 10
	expected := []int{0, 24107, 16552, 12125, 9427, 13152, 21440, 3383, 6873, 16117}
	actual := CalcRandom(s, l)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("Expected random numbers to be: %v but was: %v", expected, actual)
		}
	}
}

func TestCalcUnitIndexes(t *testing.T) {
	rands := []int{0, 34, 1000, 3, 76, 93}
	l := 5
	expected := []int{0, 4, 0, 3, 1, 3}
	actual := CalcUnitIndexes(rands, l)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("Expected unit indexes to be: %v but was: %v", expected, actual)
		}
	}
}

func TestIsRowFull(t *testing.T) {
	b := NewBoard(2, 2, []Cell{Cell{0, 0}, Cell{1, 0}})

	if !b.IsRowFull(0) {
		t.Errorf("Expected row to be full: %v but wasn't.", b[0])
	}

	if b.IsRowFull(1) {
		t.Errorf("Expected row not to be full: %v but was.", b[1])
	}
}

func TestClearFullRows(t *testing.T) {
	// bottom row is full
	b := NewBoard(2, 2, []Cell{
		Cell{0, 0},
		Cell{0, 1}, Cell{1, 1}})
	expected := NewBoard(2, 2, []Cell{Cell{0, 1}})
	actual, _ := b.ClearFullRows()

	if !equalsBoard(expected, actual) {
		t.Errorf("Expected cleared board to be: %v, but was: %v",
			expected, actual)
	}

	// top row is full
	b = NewBoard(3, 2, []Cell{
		Cell{0, 0}, Cell{1, 0},
		Cell{1, 1},
		Cell{0, 2}})
	expected = NewBoard(3, 2, []Cell{
		Cell{1, 1}, Cell{0, 2}})
	actual, _ = b.ClearFullRows()

	if !equalsBoard(expected, actual) {
		t.Errorf("Expected cleared board to be: %v, but was: %v",
			expected, actual)
	}

	// two rows are full
	b = NewBoard(3, 2, []Cell{
		Cell{0, 0}, Cell{1, 0},
		Cell{1, 1},
		Cell{0, 2}, Cell{1, 2}})
	expected = NewBoard(3, 2, []Cell{Cell{1, 2}})
	actual, _ = b.ClearFullRows()

	if !equalsBoard(expected, actual) {
		t.Errorf("Expected cleared board to be: %v, but was: %v",
			expected, actual)
	}

	// all rows are full
	b = NewBoard(3, 2, []Cell{
		Cell{0, 0}, Cell{1, 0},
		Cell{0, 1}, Cell{1, 1},
		Cell{0, 2}, Cell{1, 2}})
	expected = NewBoard(3, 2, []Cell{})
	actual, _ = b.ClearFullRows()

	if !equalsBoard(expected, actual) {
		t.Errorf("Expected cleared board to be: %v, but was: %v",
			expected, actual)
	}

	// all cells are empty
	b = NewBoard(3, 2, []Cell{})
	expected = NewBoard(3, 2, []Cell{})
	actual, _ = b.ClearFullRows()

	if !equalsBoard(expected, actual) {
		t.Errorf("Expected cleared board to be: %v, but was: %v",
			expected, actual)
	}

	// all rows are empty
	b = NewBoard(3, 2, []Cell{Cell{1, 0}, Cell{1, 2}})
	expected = NewBoard(3, 2, []Cell{Cell{1, 0}, Cell{1, 2}})
	actual, _ = b.ClearFullRows()

	if !equalsBoard(expected, actual) {
		t.Errorf("Expected cleared board to be: %v, but was: %v",
			expected, actual)
	}

	// some rows are full
	b = NewBoard(4, 2, []Cell{
		Cell{0, 0},
		Cell{0, 2}, Cell{1, 2},
		Cell{1, 3}})
	expected = NewBoard(4, 2, []Cell{Cell{0, 1}, Cell{1, 3}})
	actual, _ = b.ClearFullRows()

	if !equalsBoard(expected, actual) {
		t.Errorf("Expected cleared board to be: %v, but was: %v",
			expected, actual)
	}
}

func TestInsertPowerPhrases(t *testing.T) {
	s := "bap"
	expected := "ei!"
	actual := InsertPowerPhrases(s)

	if expected != actual {
		t.Errorf("Expected solution to be: %v, but was: %v",
			expected, actual)
	}

	s = "bapbap"
	expected = "ei!ei!"
	actual = InsertPowerPhrases(s)

	if expected != actual {
		t.Errorf("Expected solution to be: %v, but was: %v",
			expected, actual)
	}

	s = "bapaaplaapbap"
	expected = "ei!ia! ia!ei!"
	actual = InsertPowerPhrases(s)

	if expected != actual {
		t.Errorf("Expected solution to be: %v, but was: %v",
			expected, actual)
	}
}

func TestCalcPowerScore(t *testing.T) {
	s := "ei!"
	expected := 306
	actual := CalcPowerScore(s)

	if expected != actual {
		t.Errorf("Expected solution to be: %v, but was: %v",
			expected, actual)
	}

	s = "ei!ei!"
	expected = 312
	actual = CalcPowerScore(s)

	if expected != actual {
		t.Errorf("Expected solution to be: %v, but was: %v",
			expected, actual)
	}

	s = "ei!ia! ia!ei!"
	expected = 626
	actual = CalcPowerScore(s)

	if expected != actual {
		t.Errorf("Expected solution to be: %v, but was: %v",
			expected, actual)
	}
}
