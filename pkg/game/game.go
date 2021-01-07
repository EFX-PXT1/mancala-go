package game

import (
	"errors"
	"fmt"
	"strconv"
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

// CreatePosition creates and initialise a position
// values supplied are used - with zero being default
func CreatePosition(vals ...int) (p *Position) {
	bar0 := make([]int, WIDTH()+1)
	bar1 := make([]int, WIDTH()+1)
	p = &Position{}
	p.Row[0] = Side{Items: bar0}
	p.Row[1] = Side{Items: bar1}
	for i, v := range vals {
		r := i / (WIDTH() + 1)
		c := i % (WIDTH() + 1)
		p.Row[r%2].Items[c] = v
	}
	return
}

// CreatePositionCsv creates and initialise a position
// values supplied as a csv
func CreatePositionCsv(csv string) (p *Position) {
	bar0 := make([]int, WIDTH()+1)
	bar1 := make([]int, WIDTH()+1)
	p = &Position{}
	p.Row[0] = Side{Items: bar0}
	p.Row[1] = Side{Items: bar1}
	for i, s := range strings.Split(csv, ",") {
		r := i / (WIDTH() + 1)
		c := i % (WIDTH() + 1)
		if v, err := strconv.Atoi(s); err == nil {
			p.Row[r%2].Items[c] = v
		}
	}
	return
}

// AsCsv returns a string representation of a Position
func (p *Position) AsCsv() string {
	var s []string
	for r := range []int{0, 1} {
		for i := 0; i < WIDTH()+1; i++ {
			s = append(s, strconv.Itoa(p.Row[r].Items[i]))
		}
	}
	return strings.Join(s, ",")
}

// near is a convenience helper
func (p *Position) near() *Side {
	return &p.Row[0]
}

// far is a convenience helper
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

// IsValid confirm there is no corruption
func (p *Position) IsValid() (bool, int) {
	// Check a position is valid, specifically
	// that the sum of stones is correct
	correct := STONE() * WIDTH() * 2
	sum := 0
	for _, v := range p.near().Items {
		sum += v
	}
	for _, v := range p.far().Items {
		sum += v
	}
	return sum == correct, correct - sum
}

// ValidMoves returns an array of all valid moves
func (p *Position) ValidMoves() (holes []int) {
	holes = make([]int, 0, WIDTH()-1)
	for i := 1; i <= WIDTH(); i++ {
		if p.near().Items[i] > 0 {
			holes = append(holes, i)
		}
	}
	return
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
	// EndOfGame when all done
	EndOfGame MoveResult = 2
)

// Move creates a new position given a players move
func (p *Position) Move(hole int) (*Position, *Position, MoveResult, error) {
	// validate in range
	if hole < 1 || hole > WIDTH() {
		return p, nil, BadMove, errors.New("hole not in range")
	}

	// validate hole has stones
	stones := p.near().Items[hole]
	if stones == 0 {
		return p, nil, BadMove, errors.New("invalid move")
	}

	// create delta position
	delta, lastRow, lastHole := deltaPosition(hole, stones)
	// fmt.Printf("deltaPosition lastRow:%d, lastHole:%d\n", lastRow, lastHole)
	// combine
	result := p.add(delta)

	// determina result from last position
	moveResult := EndOfTurn
	if lastHole == 0 {
		moveResult = RepeatTurn
	}

	// check for steal
	if isSteal, opRow, opHole, opCount := result.IsSteal(lastRow, lastHole); isSteal {
		// create steal position
		steal := stealPosition(lastRow, lastHole, opRow, opHole, opCount)
		// apply
		result = result.add(steal)
	}

	if result.IsGameEnd() {
		moveResult = EndOfGame
	}

	return result, delta, moveResult, nil
}

// IsGameEnd checks for end of game
func (p *Position) IsGameEnd() bool {
	// end of game if all holes on either side zero
	near := 0
	for _, v := range p.near().holes() {
		near += v
	}
	far := 0
	for _, v := range p.far().holes() {
		far += v
	}
	return (near == 0 || far == 0)
}

// IsSteal determines if last position is a steal
func (p *Position) IsSteal(row int, hole int) (steal bool, opRow int, opHole int, opCount int) {
	if hole == 0 {
		return
	}
	// check if last position resulted in a single stone
	// and opposite isn't empty
	opRow = (row + 1) % 2
	opHole = WIDTH() + 1 - hole
	opCount = p.Row[opRow].Items[opHole]
	if opCount > 0 && p.Row[row].Items[hole] == 1 {
		steal = true
	}
	return
}

// add one position to another
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

// addHole computes the correct offset on the correct row
// from the starting hole and step count.
// if the computed hole is the opponent home, skip is returned
func (p *Position) addHole(start int, count int, value int) (skip bool, row int, offset int) {
	row = 0
	for count > 0 {
		offset := start - count
		if offset >= 0 {
			// found the correct row
			if offset == 0 && row%2 == 1 {
				// skip
				return true, 1, 0
			}
			// add value to existing value
			p.Row[row%2].Items[offset] = p.Row[row%2].Items[offset] + value
			return false, row % 2, offset
		}
		count = count - start // reduce count
		row++                 // on to next row
		start = WIDTH() + 1   // +1 for zero index
	}
	return false, row % 2, count
}

// ChangePlayer create a new position from other perspective
func (p *Position) ChangePlayer() (s *Position) {
	bar0 := make([]int, WIDTH()+1)
	bar1 := make([]int, WIDTH()+1)
	copy(bar0, p.far().Items)
	copy(bar1, p.near().Items)

	s = &Position{}
	s.Row[0] = Side{Items: bar0}
	s.Row[1] = Side{Items: bar1}
	return
}

// deltaPosition creates a position with each hole
// having the change of stones required
// it return the final row and hole populated
func deltaPosition(h int, count int) (p *Position, row int, hole int) {
	p = ZeroPosition()
	p.near().Items[h] = -count
	for i := 1; count > 0; i, count = i+1, count-1 {
		var skip bool
		if skip, row, hole = p.addHole(h, i, 1); skip {
			// we need to adjust our loop counters
			count = count + 1
		}
	}
	return
}

// stealPosition creates a position with each hole
// having the change of stones required for a steal
func stealPosition(r int, h int, opRow int, opHole int, opCount int) (p *Position) {
	p = ZeroPosition()
	p.Row[opRow].Items[opHole] = -opCount
	p.Row[r].Items[h] = -1
	p.near().Items[0] = opCount + 1
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
