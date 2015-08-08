package main

import "fmt"

func main() {
  cs := []Cell{Cell{0,0}, Cell{1,1}}
  fmt.Printf("hello, world %v\n", NewBoard(3, 2, cs))
}

type Board [][]bool
type Cell struct {
  row int
  col int
}
type Unit struct {
  cells []Cell
  pivot Cell
}
type BoardUnit Unit
type Command int
const (
  E Command = iota
  W
  SE
  SW
  RC
  RCC
)

func NewBoard(rows int, cols int, cells []Cell) Board {
  b := make([][]bool, rows)
  for i := range b {
    b[i] = make([]bool, cols)
  }

  for _, c := range cells {
    b[c.row][c.col] = true
  }

  return b
}

func TargetLocation(board Board, unit Unit) BoardUnit {
  return BoardUnit{[]Cell{}, Cell{}}
}

func MoveToTarget(board Board, unit Unit, target BoardUnit) []Command {
  return []Command{E, SW}
}
