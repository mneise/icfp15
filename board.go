package main

import "flag"
import "fmt"
import "io/ioutil"
import "math"
import "encoding/json"

func main() {
	params := ParseArgs()

	solution := ""
	outs := make([]Output, len(params.Program.SourceSeeds))
	b := NewBoard(params.Program.Height, params.Program.Width, params.Program.Filled)

	// is := CalcUnitIndexes
	for i, seed := range params.Program.SourceSeeds {

		rs := CalcRandom(seed, params.Program.SourceLength)
		is := CalcUnitIndexes(rs, len(params.Program.Units))

		for _, i := range is {
			u := params.Program.Units[i]
			t := TargetLocation(b, u)
			s := b.StartLocation(u)
			m := b.MoveSequence(s, t)
			cs := MovesToCommands(m)
			for _, c := range cs {
				solution = solution + c
			}
			b = b.FillCells(t.Members)
			b = b.ClearFullRows()
		}

		out := Output{
			ProblemId: params.Program.Id,
			Seed:      seed,
			Tag:       "hippo rules.",
			Solution:  solution,
		}

		outs[i] = out
		solution = ""
	}

	o, err := json.Marshal(&outs)
	if err != nil {
		panic(fmt.Sprintf("can't marshal to json: %v", err))
	}
	fmt.Println(string(o))
}

func (b Board) FillCells(cells []Cell) Board {
	nb := NewBoard(b.Height(), b.Width(), cells)
	for y, row := range b {
		for x, cell := range row {
			if cell {
				nb[y][x] = true
			}
		}
	}

	return nb
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
type Move int

func (m Move) String() string {
	switch {
	case m == E:
		return "E"
	case m == SE:
		return "SE"
	case m == W:
		return "W"
	case m == SW:
		return "SW"
	}
	return "?"
}

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

type Output struct {
	ProblemId int    `json:"problemId"`
	Seed      int    `json:"seed"`
	Tag       string `json:"tag"`
	Solution  string `json:"solution"`
}

type Params struct {
	Program              Program
	TimeLimitSeconds     int
	MemoryLimitMegaBytes int
	Cores                int
	PhraseOfPower        string
}

func ParseArgs() Params {
	var f = flag.String("f", "", "input file name")
	var t = flag.Int("t", 0, "time limit in seconds")
	var m = flag.Int("m", 0, "memory limit in megabytes")
	var c = flag.Int("c", 0, "number of cores available")
	var p = flag.String("p", "Ei!", "phrase of power")

	flag.Parse()

	d, err := ioutil.ReadFile(*f)
	if err != nil {
		panic(fmt.Sprintf("can't open file %v", f))
	}
	return Params{
		Program:              *ReadProgram(d),
		TimeLimitSeconds:     *t,
		MemoryLimitMegaBytes: *m,
		Cores:                *c,
		PhraseOfPower:        *p,
	}
}

const (
	E Move = iota
	W
	SE
	SW
	RC
	RCC
)

var commands = map[Move][]string{
	E:   []string{"b", "c", "e", "f", "y", "2"},
	W:   []string{"p", "'", "!", ".", "0", "3"},
	SE:  []string{"l", "m", "n", "o", " ", "5"},
	SW:  []string{"a", "g", "h", "i", "j", "4"},
	RC:  []string{"d", "q", "r", "v", "z", "1"},
	RCC: []string{"k", "s", "t", "u", "w", "x"},
}

func MovesToCommands(ms []Move) []string {
	cs := []string{}
	for _, m := range ms {
		cs = append(cs, commands[m][0])
	}
	return cs
}

func ReadProgram(data []byte) *Program {
	p := &Program{}
	json.Unmarshal(data, &p)

	return p
}

func NewBoard(height int, width int, cells []Cell) Board {
	b := make([][]bool, height)
	for i := range b {
		b[i] = make([]bool, width)
	}

	for _, c := range cells {
		b[c.Y][c.X] = true
	}

	return b
}

func (c Cell) Move(m Move, b Board) Cell {
	switch {
	case m == W:
		return Cell{X: c.X - 1, Y: c.Y}
	case m == E:
		return Cell{X: c.X + 1, Y: c.Y}

	case m == SW && c.Y%2 == 1 && (c.X > 0 || c.X == 0):
		return Cell{X: c.X, Y: c.Y + 1}
	case m == SW && c.Y%2 == 0 && (c.X > 0):
		return Cell{X: c.X - 1, Y: c.Y + 1}

	case m == SE && c.Y%2 == 0 && (c.X < b.Width()-1 || c.X == b.Width()-1):
		return Cell{X: c.X, Y: c.Y + 1}
	case m == SE && c.Y%2 == 1 && c.X < b.Width()-1:
		return Cell{X: c.X + 1, Y: c.Y + 1}
	}

	return Cell{X: -1, Y: -1}
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
	return c.X >= 0 &&
		c.X < b.Width() &&
		c.Y >= 0 &&
		c.Y < b.Height() &&
		!b.isCellFull(c)
}

// A unit is in a valid location if all of its cells are on empty board
// cells. Note that a unit's pivot point need not be on a board cell.
func (u Unit) isValid(b Board) bool {
	for _, c := range u.Members {
		if !c.isValid(b) {
			return false
		}
	}

	return true
}

func (b Board) isCellFull(c Cell) bool {
	return b[c.Y][c.X]
}

func (b Board) IsRowFull(r int) bool {
	for _, c := range b[r] {
		if !c {
			return false
		}
	}
	return true
}

func direction(s, t Cell) (xd, yd int) {
	yd = t.Y - s.Y
	// ys = yd
	// if ys < 0 {
	// 	ys = -ys
	// }

	xd = t.X - s.X
	// xs = xd
	// if xs < 0 {
	// 	xs = -xs
	// }

	return xd, yd //, xs, ys
}

func moves(xd, yd int) []Move {
	// neg - left
	// pos - right
	// zero - down

	switch {
	case xd < 0 && yd > 0: // left down
		return []Move{W, SW, SE, E}
	case xd == 0 && yd > 0: // down
		return []Move{SE, SW, E, W}
	case xd > 0 && yd > 0: // right down
		return []Move{E, SE, SW, W}
	case xd < 0 && yd == 0: // left
		return []Move{W}
	case xd > 0 && yd == 0: // right
		return []Move{E}
	case xd == 0 && yd == 0: // done
	case yd < 0: // can't move up
	}

	return []Move{}
}

// TODO: never EW or WE
func (b Board) MoveSequence(s Unit, t Unit) []Move {
	// fmt.Printf("move from %v to %v\n", s.Pivot, t.Pivot)
	mu := s
	mp := s.Pivot
	xd, yd := direction(s.Pivot, t.Pivot)
	ms := []Move{}

	for true {
		before := len(ms)
		// fmt.Printf("main loop mp %v xd %v yd %v ms %v\n", mp, xd, yd, ms)

		for _, m := range moves(xd, yd) {
			// fmt.Printf("found move %v\n", m)
			// try to move pivot / unit
			tp := mp.Move(m, b)
			tu := mu.MoveTo(tp)
			if tp.isValid(b) && tu.isValid(b) { // found valid one,yay!
				mu = tu
				mp = tp
				xd, yd = direction(mp, t.Pivot)
				ms = append(ms, m)
				break
			}
		}

		if before == len(ms) {
			break // couldn't find move, skip out
		}
	}

	if mp.X == t.Pivot.X && mp.Y == t.Pivot.Y {
		return append(ms, SE) // lock in move
	}

	return []Move{} // can't find a legal way
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

func CalcRandom(s int, l int) []int {

	rands := make([]int, l)
	m := 1 << 31
	a := 1103515245
	c := 12345
	x := s

	for i := 0; i < l; i++ {
		out := (x >> 16) & 0x7fff
		x = ((a*x + c) % m)
		if x < 0 {
			x += 4294967296
		}
		rands[i] = out
	}

	return rands
}

func CalcUnitIndexes(rands []int, l int) []int {
	idxs := make([]int, len(rands))
	for i, rand := range rands {
		idxs[i] = rand % l
	}
	return idxs
}

func (b Board) ClearFullRows() Board {
	nb := NewBoard(b.Height(), b.Width(), []Cell{})
	copy(nb, b)
	for i := b.Height() - 1; i >= 0; i-- {
		if nb.IsRowFull(i) {
			if i == b.Height()-1 {
				nb = nb[:i]
			} else {
				nb = append(nb[:i], nb[i+1:]...)
			}
			r := make([]bool, b.Width())
			nb = append([][]bool{r}, nb...)
			i++
		}
	}
	return nb
}
