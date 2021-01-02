package game

import (
	"errors"
	"fmt"
	"strings"
)

// Game represents the dimensions of this game
type Game struct {
	// Width is the size of the board
	Width int
	// Stone is the initial number per hole
	Stone int
}

var game *Game

// DefineGame sets the dimensions of the game
func DefineGame(width int, stone int) {
	game = &Game{
		Width: width,
		Stone: stone,
	}
}

// WIDTH retrieve universal game width
func WIDTH() int { return game.Width }

// STONE retrieve universal game stone number
func STONE() int { return game.Stone }

// Side represent one players side
type Side struct {
	Items []int
}

func (s *Side) home() int {
	return s.Items[0]
}

func (s *Side) holes() []int {
	return s.Items[1:]
}

// Position is the state for a single round
type Position struct {
	Row [2]Side // 0 is near, 1 is far
}

func (p *Position) near() *Side {
	return &p.Row[0]
}

func (p *Position) far() *Side {
	return &p.Row[1]
}

// Show displays position to the console
func (p *Position) Show() {
	// Far row
	far := ""
	for _, v := range p.far().holes() {
		// left to right
		far = far + fmt.Sprintf(" %2d", v)
	}

	// Near row
	near := ""
	for _, v := range p.near().holes() {
		// right to left
		near = fmt.Sprintf(" %2d", v) + near
	}

	// middle
	padding := max(len(far), len(near))

	// output over 3 lines
	fmt.Printf("   %s\n", far)
	fmt.Printf("%2d %s %2d\n", p.far().home(), strings.Repeat(" ", padding), p.near().home())
	fmt.Printf("  %s\n", near)
	fmt.Println(strings.Repeat("-", padding+6))
}

// max finds largest of two ints
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MoveResult captures who goes next
type MoveResult int8

const (
	// BadMove if not valie
	BadMove MoveResult = -1
	// EndOfTurn to swap player
	EndOfTurn MoveResult = 0
	// RepeatTurn for same player
	RepeatTurn MoveResult = 1
)

// Move creates a new position given a players move
func (p *Position) Move(hole int) (*Position, MoveResult, error) {
	// validate in range
	if hole < 1 || hole > WIDTH() {
		return p, BadMove, errors.New("hole not in range")
	}

	// validate hole has stones
	stones := p.near().Items[hole]
	if stones == 0 {
		return p, BadMove, errors.New("invalid move")
	}

	// create delta position
	d, lastRow, lastHole := deltaPosition(hole, stones)

	// combine
	result := p.add(d)

	// determina result from last position
	moveResult := EndOfTurn
	if lastHole == 0 {
		moveResult = RepeatTurn
	}

	// check for steal
	// need some steel processing
	if p.isSteal(lastRow, lastHole) {
		// todo
	}

	return result, moveResult, nil
}

func (p *Position) isSteal(row int, hole int) bool {
	return false
}

func (p *Position) add(delta *Position) (pos *Position) {
	pos = ZeroPosition()
	for row := 0; row < 2; row++ {
		for hole := 0; hole <= WIDTH(); hole++ {
			pos.Row[row].Items[hole] = p.Row[row].Items[hole] +
				delta.Row[row].Items[hole]
		}
	}
	return
}

// setHole computes the correct offset on the correct row
// from the starting hole and step count.
// if the computed hole is the opponent home, skip is returned
func (p *Position) setHole(start int, count int, value int) (skip bool, row int, offset int) {
	row = 0
	for count > 0 {
		offset := start - count
		if offset >= 0 {
			// found the correct row
			if offset == 0 && row%2 == 1 {
				// skip
				return true, 1, 0
			}
			// single round loop function
			p.Row[row%2].Items[offset] = value
			return false, row % 2, offset
		}
		count = count - start // reduce count
		row++                 // on to next row
		start = WIDTH() + 1   // +1 for zero index
	}
	return false, row % 2, count
}

// deltaPosition creates a position with each hole
// having the change of stones required
// it return the final row and hole populated
func deltaPosition(h int, count int) (p *Position, row int, hole int) {
	p = ZeroPosition()
	p.near().Items[h] = -count
	for i := 1; count > 0; i, count = i+1, count-1 {
		var skip bool
		if skip, row, hole = p.setHole(h, i, 1); skip {
			// we need to adjust our loop counters
			count = count + 1
		}
	}
	p.Show()
	return
}

// StartPosition creates the standard start
func StartPosition() (p *Position) {
	bar0 := make([]int, WIDTH()+1)
	bar1 := make([]int, WIDTH()+1)
	for i := range bar0 {
		bar0[i] = STONE()
	}
	bar0[0] = 0
	copy(bar1, bar0)

	p = &Position{}
	p.Row[0] = Side{Items: bar0}
	p.Row[1] = Side{Items: bar1}
	return
}

// ZeroPosition creates an empty position
func ZeroPosition() (p *Position) {
	bar0 := make([]int, WIDTH()+1)
	bar1 := make([]int, WIDTH()+1)

	p = &Position{}
	p.Row[0] = Side{Items: bar0}
	p.Row[1] = Side{Items: bar1}
	return
}

// DiagnosticPosition creates a position with
// the number of stones the value of the hole
func DiagnosticPosition() (p *Position) {
	bar0 := make([]int, WIDTH()+1)
	bar1 := make([]int, WIDTH()+1)
	for i := range bar0 {
		bar0[i] = i
	}
	bar0[0] = 0
	copy(bar1, bar0)

	p = &Position{}
	p.Row[0] = Side{Items: bar0}
	p.Row[1] = Side{Items: bar1}
	return
}
