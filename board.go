package main

import "fmt"
import "math"
import "encoding/json"

func main() {
	u := Unit{Members: []Cell{Cell{0, 0}}, Pivot: Cell{0, 0}}
	b := NewBoard(3, 3, []Cell{})

	t := TargetLocation(b, u)
	s := b.StartLocation(u)
	m := b.MoveSequence(s, t)

	fmt.Printf("Found moves: %v for board: %v and unit: %v\n", m, b, u)
}

type Board [][]bool
type Cell struct {
	X int
	Y int
}
type Unit struct {
	Members []Cell
	Pivot   Cell
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

// todo: should we use float64 just cause json
type Program struct {
	Id           int
	Units        []Unit
	Width        int
	Height       int
	Filled       []Cell
	SourceLength int
	SourceSeeds  []int
}

func ReadProgram(data []byte) *Program {
	p := &Program{}
	json.Unmarshal(data, &p)

	return p
}

func NewBoard(rows int, cols int, cells []Cell) Board {
	b := make([][]bool, rows)
	for i := range b {
		b[i] = make([]bool, cols)
	}

	for _, c := range cells {
		b[c.Y][c.X] = true
	}

	return b
}

func (u Unit) MoveTo(cell Cell) Unit {
	unit := Unit{Members: make([]Cell, len(u.Members)), Pivot: cell}
	for i, member := range u.Members {
		x := cell.X + (member.X - u.Pivot.X)
		y := cell.Y + (member.Y - u.Pivot.Y)
		unit.Members[i] = Cell{Y: y, X: x}
	}
	return unit
}

func TargetLocation(b Board, u Unit) Unit {
	t := Unit{}
	for y := range b {
		for x := range b[y] {
			tmp := u.MoveTo(Cell{x, y})
			if tmp.isValid(b) {
				t = tmp
			}
		}
	}

	return t
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

func (c Cell) isValid(b Board) bool {
	if c.X < 0 ||
		c.X >= b.Width() ||
		c.Y < 0 ||
		c.Y >= b.Height() ||
		c.isFull(b) {
		return false
	}
	return true
}

func (u Unit) isValid(b Board) bool {
	if !u.Pivot.isValid(b) {
		return false
	}

	for _, c := range u.Members {
		if !c.isValid(b) {
			return false
		}
	}

	return true
}

func (c Cell) isFull(b Board) bool {
	return b[c.Y][c.X]
}

func (b Board) MoveSequence(s Unit, t Unit) []Command {

	commands := []Command{}

	// neg - left
	// pos - right
	// zero - down
	direction := t.Pivot.X - s.Pivot.X
	xSteps := direction
	if xSteps < 0 {
		xSteps = -xSteps
	}
	ySteps := t.Pivot.Y - s.Pivot.Y

	// move left / right
	for i := 0; i < xSteps; i++ {
		if direction < 0 {
			commands = append(commands, W)
		} else {
			commands = append(commands, E)
		}
	}

	// move down
	for i := 0; i < ySteps; i++ {
		switch {
		case i%2 == 0:
			commands = append(commands, SE)
		case i%2 == 1:
			commands = append(commands, SW)
		}
	}

	return commands
}

func (c Cell) ShiftX(offset int) Cell {
	return Cell{X: c.X + offset, Y: c.Y}
}

func (b Board) StartLocation(u Unit) Unit {
	offset := (b.Width() - u.Width()) / 2
	newPivot := u.Pivot.ShiftX(offset)

	return u.MoveTo(newPivot)
}

func (u Unit) Width() int {
	minX := math.MaxInt32
	maxX := -1

	for _, member := range u.Members {
		if minX > member.X {
			minX = member.X
		}
		if maxX < member.X {
			maxX = member.X
		}
	}

	return 1 + maxX - minX
}

func (u Unit) Height() int {
	minY := math.MaxInt32
	maxY := -1

	for _, member := range u.Members {
		if minY > member.Y {
			minY = member.Y
		}
		if maxY < member.Y {
			maxY = member.Y
		}
	}

	return 1 + maxY - minY
}
