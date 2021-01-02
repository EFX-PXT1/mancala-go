package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/EFX-PXT1/mancala-go/pkg/game"
)

func main() {
	args := os.Args[1:]

	game.DefineGame(6, 4)

	pos := game.StartPosition()
	pos.Show()

	var x string
	for len(args) > 0 {
		x, args = args[0], args[1:]
		if hole, err := strconv.Atoi(x); err == nil {
			if pos, _, _, err = pos.Move(hole); err == nil {
				pos.Show()
			} else {
				fmt.Printf(" %s %s\n---\n", x, err)
			}
		}
	}
}
