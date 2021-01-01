package game

import (
	"errors"
	"fmt"
	"strings"
)

// WIDTH is the size of the board
const WIDTH = 6

// STONE is the initial number per hole
const STONE = 10

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
	// validate hole has stones
	stones := p.near().Items[hole]
	if stones == 0 {
		return p, BadMove, errors.New("invalid move")
	}

	// create delta position
	d := deltaPosition(hole, stones)

	// combine
	result := p.add(d)
	return result, EndOfTurn, nil
}

func (p *Position) add(delta *Position) (pos *Position) {
	pos = ZeroPosition()
	for row := 0; row < 2; row++ {
		for hole := 0; hole <= WIDTH; hole++ {
			pos.Row[row].Items[hole] = p.Row[row].Items[hole] +
				delta.Row[row].Items[hole]
		}
	}
	return
}

func (p *Position) setHole(start int, count int, value int) (skip bool) {
	row := 0
	for count > 0 {
		offset := start - count
		if offset >= 0 {
			// found the correct row
			if offset == 0 && row%2 == 1 {
				// skip
				return true
			}
			// single round loop function
			p.Row[row%2].Items[offset] = value
			return false
		}
		count = count - start // reduce count
		row++                 // on to next row
		start = WIDTH + 1     // +1 for zero index
	}
	return false
}

func deltaPosition(h int, count int) (p *Position) {
	p = ZeroPosition()
	p.near().Items[h] = -count
	for i := 1; count > 0; i, count = i+1, count-1 {
		if skip := p.setHole(h, i, 1); skip {
			// we need to adjust our loop counters
			count = count + 1
		}
	}
	p.Show()
	return
}

// StartPosition creates the standard start
func StartPosition() (p *Position) {
	bar0 := make([]int, WIDTH+1)
	bar1 := make([]int, WIDTH+1)
	for i := range bar0 {
		bar0[i] = STONE
		// bar0[i] = i
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
	bar0 := make([]int, WIDTH+1)
	bar1 := make([]int, WIDTH+1)

	p = &Position{}
	p.Row[0] = Side{Items: bar0}
	p.Row[1] = Side{Items: bar1}
	return
}
