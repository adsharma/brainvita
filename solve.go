package main

import (
	"github.com/golang-collections/collections/stack"
)

type Direction int

const (
	East = iota
	South
	West
	North
)

func (dir Direction) String() string {
	return []string{"East", "South", "West", "North"}[dir]
}

type Board struct {
	board [7][7]bool
}

type Pos struct {
	i int
	j int
}

type moveAction struct {
	pos Pos
	dir Direction
}

var invalid moveAction

func (b Board) isValid(p Pos) bool {
	return (p.i >= 2 && p.i <= 4 && p.j >= 0 && p.j < 7) ||
		(p.j >= 2 && p.j <= 4 && p.i >= 0 && p.i < 7)
}

func (b *Board) init() {
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			b.board[i][j] = b.isValid(Pos{i, j}) && !(i == 3 && j == 3)
		}
	}
}

func (b Board) incr(act moveAction) moveAction {
	newM := act
	if act.dir < 3 {
		newM.dir = act.dir + 1
		return newM
	} else if act.pos.i < 6 {
		newM.dir = 0
		newM.pos.i += 1
		return newM
	} else if act.pos.j < 6 {
		newM.dir = 0
		newM.pos.i = 0
		newM.pos.j += 1
		return newM
	}
	return invalid
}

func (b Board) isBoardCountOne() bool {
	c := 0
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			if b.board[i][j] {
				c += 1
				if c > 1 {
					return false
				}
			}
		}
	}
	return true
}

func (b *Board) move(action moveAction) {
	switch action.dir {
	case East:
		b.moveX(1, action.pos)
	case West:
		b.moveX(-1, action.pos)
	case South:
		b.moveY(1, action.pos)
	case North:
		b.moveY(-1, action.pos)
	}
}

func (b Board) movableX(i int, pos Pos) bool {
	x := pos.i
	y := pos.j
	return b.board[x][y] && b.board[x+i][y] && !b.board[x+2*i][y]
}

func (b Board) movableY(j int, pos Pos) bool {
	x := pos.i
	y := pos.j
	return b.board[x][y] && b.board[x][y+j] && !b.board[x][y+2*j]
}

func (b *Board) moveX(i int, pos Pos) bool {
	x := pos.i
	y := pos.j
	b.board[x][y] = false
	b.board[x+i][y] = false
	b.board[x+2*i][y] = true
	return true
}

func (b *Board) moveY(j int, pos Pos) bool {
	x := pos.i
	y := pos.j
	b.board[x][y] = false
	b.board[x][y+j] = false
	b.board[x][y+2*j] = true
	return true
}

func (b Board) isMovable(act moveAction) bool {
	if !b.isValid(act.pos) {
		return false
	}

	x := act.pos.i
	y := act.pos.j
	switch act.dir {
	case East:
		return b.isValid(Pos{x + 2, y}) && b.movableX(1, act.pos)
	case West:
		return b.isValid(Pos{x - 2, y}) && b.movableX(-1, act.pos)
	case South:
		return b.isValid(Pos{x, y + 2}) && b.movableY(1, act.pos)
	case North:
		return b.isValid(Pos{x, y - 2}) && b.movableY(-1, act.pos)
	}
	return false
}

func (b Board) Move(action moveAction) Board {
	newb := b
	newb.move(action)
	return newb
}

func (b Board) findFirstMove() moveAction {
	return b.findNextMove(moveAction{Pos{0, 0}, 0})
}

func (b Board) findNextMove(action moveAction) moveAction {
	next := b.incr(action)
	for next != invalid {
		if b.isMovable(next) {
			return next
		}
		next = b.incr(next)
	}
	return next
}

var moves *stack.Stack

func Solver(b Board) bool {
	if b.isBoardCountOne() {
		return true
	}
	next := b.findFirstMove()

	for next != invalid {
		if Solver(b.Move(next)) {
			moves.Push(next)
			return true
		}
		next = b.findNextMove(next)
	}

	return false
}

func printSolution() {

	for moves.Len() > 0 {
		m := moves.Pop().(moveAction)
		println(m.pos.i, ", ", m.pos.j, ", ", m.dir.String())
	}
}

func main() {
	b := Board{}
	invalid = moveAction{Pos{-1, -1}, -1}
	b.init()
	moves = stack.New()
	Solver(b)
	printSolution()
}
