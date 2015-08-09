package main

import "flag"
import "time"
import "fmt"
import "io/ioutil"
import "math"
import "encoding/json"

func logBoard(p Params, m string, b Board) {
	if p.Debug {
		fmt.Printf("%v:\n%v\n", m, b)
	}
}

func logMsg(p Params, m string) {
	if p.Debug {
		fmt.Printf("%v\n", m)
	}
}

func main() {
	params := ParseArgs()

	b := NewBoard(params.Program.Height, params.Program.Width, params.Program.Filled)
	solution := ""
	outs := make([]Output, len(params.Program.SourceSeeds))
	totalScore := 0

	if params.LogBoard {
		logBoard(params, fmt.Sprintf("Board for problem: %v", params.Program.Id), b)
		return
	}

	for i, seed := range params.Program.SourceSeeds {

		moveScores := 0
		cleared := 0
		clearedOld := 0

		rs := CalcRandom(seed, params.Program.SourceLength)
		is := CalcUnitIndexes(rs, len(params.Program.Units))
		count := 0

		for _, i := range is {
			clearedOld = cleared

			count++
			// if count > 20 { // last: 20-borked
			// 	break
			// }

			u := params.Program.Units[i]
			s := b.StartLocation(u)
			if !s.isValid(b) {
				logMsg(params, fmt.Sprintf("couldn't place unit %v %v! GAME OVER BABY", count, u))
				break
			}

			logMsg(params, "======================================================")
			logBoard(params, fmt.Sprintf("trying to place unit %v (%vth) on board", u, count), b.FillCells(s.Members))

			ts := TargetLocations(b, u)
			m := []Move{}
			t := Unit{}
			for _, t = range ts {
				m = b.MoveSequence(s, t)
				if len(m) > 0 {
					logMsg(params, fmt.Sprintf("found moves %v", m))
					xxx := s
					for ms, x := range m {
						if ms < len(m)-1 {
							xxx = xxx.Move(x)
							// logBoard(params, fmt.Sprintf("move %v step %v", x, ms), b.FillCells(xxx.Members))
						}
					}
					break
				} else {
					// logMsg(params, fmt.Sprintf("found no moves for target %v", t))
				}
			}

			if len(m) == 0 {
				logMsg(params, fmt.Sprintf("found no moves! GAME OVER BABY"))
				break
			}

			logMsg(params, fmt.Sprintf("found moves: %v", m))

			cs := MovesToCommands(m)
			for _, c := range cs {
				solution = solution + c
			}
			b = b.FillCells(t.Members)
			logBoard(params, fmt.Sprintf("unit %v placed on board", i), b)
			b, cleared = b.ClearFullRows()
			if cleared > 0 {
				logBoard(params, fmt.Sprintf("cleared full rows"), b)
			}

			points := len(u.Members) + 100*(1+cleared)*cleared/2
			lineBonus := 0
			if clearedOld > 1 {
				lineBonus = (clearedOld - 1) * points / 10
			}
			moveScores += points + lineBonus
		}

		logMsg(params, fmt.Sprintf("Game score: %v\n", moveScores))
		totalScore += moveScores

		out := Output{
			ProblemId: params.Program.Id,
			Seed:      seed,
			Tag:       fmt.Sprintf("hippo rules @ %v", time.Now()),
			Solution:  solution,
		}

		outs[i] = out
		solution = ""
		b = NewBoard(params.Program.Height, params.Program.Width, params.Program.Filled)
	}

	logMsg(params, fmt.Sprintf("Total score: %v\n", totalScore/len(outs)))

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

func (b Board) String() (s string) {

	for ri, r := range b {
		if ri%2 == 1 {
			s += " "
		}

		for ci, c := range r {
			if c {
				s += "⬢"
			} else {
				s += "⬡"
			}
			if ci < len(r)-1 {
				s += " "
			}
		}

		if ri < len(b)-1 {
			s += "\n"
		}
	}

	return s
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

type Cube struct {
	X int
	Y int
	Z int
}

func (c Cell) cube() Cube {
	// http://www.redblobgames.com/grids/hexagons/
	// # convert odd-r offset to cube
	// x = col - (row - (row&1)) / 2
	// z = row
	// y = -x-z
	x := c.X - (c.Y-(c.Y&1))/2
	z := c.Y
	y := -x - z

	return Cube{X: x, Y: y, Z: z}
}

func (u Unit) eq(o Unit) bool {
	if u.Pivot != o.Pivot {
		return false
	}

	if len(u.Members) != len(o.Members) {
		return false
	}

	for _, m := range u.Members {
		foundM := false
		for _, om := range o.Members {
			if om == m {
				foundM = true
				break
			}
		}

		if !foundM {
			return false
		}
		// if o.Members[i] != m {
		// 	return false
		// }
	}

	return true
}

func (c Cube) cell() Cell {
	// http://www.redblobgames.com/grids/hexagons/
	// # convert cube to odd-r offset
	// col = x + (z - (z&1)) / 2
	// row = z
	return Cell{X: c.X + (c.Z-(c.Z&1))/2, Y: c.Z}
}

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
	case m == RC:
		return "RC"
	case m == RCC:
		return "RCC"
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
	Debug                bool
	LogBoard             bool
}

func ParseArgs() Params {
	var f = flag.String("f", "", "input file name")
	var t = flag.Int("t", 0, "time limit in seconds")
	var m = flag.Int("m", 0, "memory limit in megabytes")
	var c = flag.Int("c", 0, "number of cores available")
	var p = flag.String("p", "Ei!", "phrase of power")
	var d = flag.Bool("d", false, "print debug output")
	var b = flag.Bool("b", false, "print start board only")

	flag.Parse()

	in, err := ioutil.ReadFile(*f)
	if err != nil {
		panic(fmt.Sprintf("can't open file %v", f))
	}
	return Params{
		Program:              *ReadProgram(in),
		TimeLimitSeconds:     *t,
		MemoryLimitMegaBytes: *m,
		Cores:                *c,
		PhraseOfPower:        *p,
		Debug:                *d,
		LogBoard:             *b,
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

func (c Cell) Move(m Move, p Cell) Cell { // TODO
	// E: y-1 x+1 z
	// SE: y-1 x z+1
	// W: y+1 x-1 z
	// SW: y x-1 z+1
	q := c.cube()
	nq := Cube{-1, -1, -1}

	switch {
	case m == E:
		nq = Cube{
			X: q.X + 1,
			Y: q.Y - 1,
			Z: q.Z,
		}
	case m == SE:
		nq = Cube{
			X: q.X,
			Y: q.Y - 1,
			Z: q.Z + 1,
		}
	case m == W:
		nq = Cube{
			X: q.X - 1,
			Y: q.Y + 1,
			Z: q.Z,
		}
	case m == SW:
		nq = Cube{
			X: q.X - 1,
			Y: q.Y,
			Z: q.Z + 1,
		}
	case m == RC:
		pq := p.cube()

		// right
		//     [ x,  y,  z]
		//     [-z, -x, -y]
		// left
		// [ x,  y,  z]
		// [-y, -z, -x]

		xd := pq.X - q.X
		yd := pq.Y - q.Y
		zd := pq.Z - q.Z

		// fmt.Printf("\napply %v to %v ", m, c)
		// fmt.Printf("xd %2d yd %2d zd %2d\n", xd, yd, zd)

		nq = Cube{
			X: q.X - yd,
			Y: q.Y - zd,
			Z: q.Z - xd,
		}

		// fmt.Printf("moved c%v q%v to q%v or c%v\n", c, c.cube(), nq, nq.cell())
	case m == RCC:
		pq := p.cube()

		// right
		//     [ x,  y,  z]
		//     [-z, -x, -y]
		// left
		// [ x,  y,  z]
		// [-y, -z, -x]

		xd := pq.X - q.X
		yd := pq.Y - q.Y
		zd := pq.Z - q.Z

		// fmt.Printf("\napply %v to %v ", m, c)
		// fmt.Printf("xd %2d yd %2d zd %2d\n", xd, yd, zd)

		nq = Cube{
			X: q.X - zd,
			Y: q.Y - xd,
			Z: q.Z - yd,
		}

		// fmt.Printf("moved c%v q%v to q%v or c%v\n", c, c.cube(), nq, nq.cell())
	}

	return nq.cell()
}

func (u Unit) Move(m Move) (nu Unit) {

	if m != RC && m != RCC {
		nu.Pivot = u.Pivot.Move(m, u.Pivot)
	} else {
		nu.Pivot = u.Pivot
	}

	for _, x := range u.Members {
		nu.Members = append(nu.Members, x.Move(m, u.Pivot))
	}
	return nu
}

func (u Unit) MoveTo(new Cell, old Cell) Unit {
	tu := Unit{}

	xd := new.cube().X - old.cube().X
	yd := new.cube().Y - old.cube().Y
	zd := new.cube().Z - old.cube().Z

	for _, om := range u.Members {
		oc := om.cube()
		nm := Cube{
			X: oc.X + xd,
			Y: oc.Y + yd,
			Z: oc.Z + zd,
		}
		tu.Members = append(tu.Members, nm.cell())
	}

	op := u.Pivot.cube()
	np := Cube{
		X: op.X + xd,
		Y: op.Y + yd,
		Z: op.Z + zd,
	}
	tu.Pivot = np.cell()

	return tu
}

func TargetLocations(b Board, u Unit) []Unit {
	ts := []Unit{}
	for y := range b {
		for x := range b[y] {
			t := u.MoveTo(Cell{x, y}, u.Pivot)
			if t.isValid(b) {
				ts = append([]Unit{t}, ts...)
			}
		}
	}

	return ts
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
		return []Move{W, SW, SE, E, RC, RCC}
	case xd == 0 && yd > 0: // down
		return []Move{SE, SW, E, W, RC, RCC}
	case xd > 0 && yd > 0: // right down
		return []Move{E, SE, SW, W}
	case xd < 0 && yd == 0: // left
		return []Move{W}
	case xd > 0 && yd == 0: // right
		return []Move{E}
	case xd == 0 && yd == 0: // done
	case yd < 0: // cant move up
	}

	return []Move{}
}

type MoveWithBacktrack struct {
	move    Move
	options []Move
}

func NewMoveWithBacktrack(xd, yd int) MoveWithBacktrack {
	return MoveWithBacktrack{
		move:    W,
		options: moves(xd, yd), // TODO should bias this
	}
}

func (b Board) MoveSequence(s Unit, t Unit) []Move {
	// msb := []MoveWithBacktrack{}
	// pls := []Unit{s}
	// seen := func(u Unit) bool {
	// 	for _, pl := range pls {
	// 		if u.eq(pl) {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	// movesString := func() string {
	// 	r := ""
	// 	for _, mb := range msb {
	// 		r += mb.move.String() + " "
	// 	}
	// 	return r
	// }

	// 	xd, yd := direction(s.Pivot, t.Pivot)

	// 	count := 0

	// 	mu := s
	// 	next := NewMoveWithBacktrack(xd, yd)
	// 	m := next.move
	// 	tu := mu.Move(m)

	// XXX:
	// 	for !mu.eq(t) {
	// 		count++
	// 		if count > 100 {
	// 			// fmt.Printf("failing with upper limit 10k\n")
	// 			break XXX
	// 		}

	// 		// found one!
	// 		if !seen(tu) && tu.isValid(b) {
	// 			// fmt.Printf("found valid move m %v adding it to msb %v.\n", m, movesString())
	// 			// fmt.Printf("found move %v, board:\n%v\n", m, b.FillCells(tu.Members))

	// 			mu = tu
	// 			msb = append(msb, next)
	// 			pls = append(pls, mu)
	// 			xd, yd = direction(mu.Pivot, t.Pivot)
	// 			next = NewMoveWithBacktrack(xd, yd)
	// 			m = next.move
	// 			tu = mu.Move(m)

	// 			continue XXX
	// 		}

	// 		// try next avail command
	// 		if len(next.options) > 0 {
	// 			// fmt.Printf("move m %v is invalid (unseen %v, valid %v), trying next %v.\n",
	// 			// 	m,
	// 			// 	!seen(tu),
	// 			// 	tu.isValid(b),
	// 			// 	next.move,
	// 			// )
	// 			next.move = next.options[0]
	// 			next.options = next.options[1:]
	// 			m = next.move
	// 			tu = mu.Move(m)
	// 			continue XXX
	// 		}

	// 		for len(msb) > 1 {
	// 			// fmt.Printf("backtracking.\n")
	// 			// fmt.Printf("len(msb) %v len(pls) %v.\n", len(msb), len(pls))
	// 			mu = pls[len(pls)-2]
	// 			pls = pls[:len(pls)-1]
	// 			next = msb[len(pls)-2]
	// 			msb = msb[:len(msb)-1]

	// 			if len(next.options) > 0 {
	// 				// fmt.Printf("restored board:\n%v\n", b.FillCells(mu.Members))

	// 				next.move = next.options[0]
	// 				next.options = next.options[1:]
	// 				m = next.move
	// 				tu = mu.Move(m)
	// 				continue XXX
	// 			}

	// 		}

	// 		// fmt.Printf("no options left\n")
	// 		break XXX
	// 	}

	// 	ms := []Move{}
	// 	for _, m := range msb {
	// 		ms = append(ms, m.move)
	// 	}
	// 	// fmt.Printf("new moves: %v\n", ms)

	///////////////////////////////////////////////////////////////////
	mp := s.Pivot
	xd, yd := direction(s.Pivot, t.Pivot)
	ms := []Move{}
	ls := []Unit{s}
	mu := s

	for true {
		before := len(ms)
	findNextMove:
		for _, m := range moves(xd, yd) {
			tu := mu.Move(m)

			if len(ms) > 0 &&
				((m == W && ms[len(ms)-1] == E) || (m == E && ms[len(ms)-1] == W)) {
				continue findNextMove
			}

			for _, pl := range ls {
				if tu.eq(pl) {
					// fmt.Printf("been here, skipping %v already got %v\n", m, ms)
					continue findNextMove
				}
			}

			if tu.isValid(b) { // found valid one,yay!

				// fmt.Printf("found move %v, board:\n%v\n", m, b.FillCells(tu.Members))
				mu = tu
				mp = tu.Pivot
				xd, yd = direction(mp, t.Pivot)
				ms = append(ms, m)
				ls = append(ls, mu)
				break
			}
		}

		if before == len(ms) {
			break // couldnt find move, skip out
		}
	}

	if mp.X == t.Pivot.X && mp.Y == t.Pivot.Y {
		// TODO: improve this.
		switch {
		case !t.Move(SE).isValid(b):
			return append(ms, SE)
		case !t.Move(SW).isValid(b):
			return append(ms, SW)
		case !t.Move(E).isValid(b):
			return append(ms, E)
		case !t.Move(W).isValid(b):
			return append(ms, W)
		}
	}
	// fmt.Printf("found moves, but not on target: %v\n", ms)

	return []Move{} // can't find a legal way
}

func (c Cell) ShiftX(offset int) Cell {
	return Cell{X: c.X + offset, Y: c.Y}
}

func (c Cell) ShiftY(offset int) Cell {
	return Cell{X: c.X, Y: c.Y + offset}
}

func (b Board) StartLocation(u Unit) Unit {
	// move to right
	minXCell := u.MinXCell()
	offset := (b.Width() - u.Width()) / 2
	nc := minXCell.ShiftX(0 - minXCell.X + offset)
	u = u.MoveTo(nc, minXCell)

	// move up
	minYCell := u.MinYCell()
	offset = 0 - minYCell.Y
	nc = minYCell.ShiftY(offset)
	u = u.MoveTo(nc, minYCell)

	return u
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

func (u Unit) MinYCell() Cell {
	minY := math.MaxInt32
	c := Cell{}

	for _, member := range u.Members {
		if minY > member.Y {
			minY = member.Y
			c = member
		}
	}

	return c
}

func (u Unit) MinXCell() Cell {
	minX := math.MaxInt32
	c := Cell{}

	for _, member := range u.Members {
		if minX > member.X {
			minX = member.X
			c = member
		}
	}

	return c
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

func (b Board) ClearFullRows() (Board, int) {
	var cleared int
	nb := NewBoard(b.Height(), b.Width(), []Cell{})
	copy(nb, b)
	for i := b.Height() - 1; i >= 0; i-- {
		if nb.IsRowFull(i) {
			cleared++
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
	return nb, cleared
}
