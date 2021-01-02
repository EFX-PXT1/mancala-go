# mancala-go

Mancala written in Go as a bit of fun to think about game strategy
and port to Python / JavaScript for some GCSE computer science

## mconsole

The initial console game which simply implements the rules.

### build

```
cd cmd/mconsole
go build
```

### run

For a standard game with a repl

```
mconsole -r
```

#### options

* --width to change the width of the board from the usual 6.
* --stones to change the initial number of stones from the usual 4.

### playing

Stones move anti-clockwise with the aim to get the most into your home hole.
The active player is at the bottom - just like sitting across a table.
Enter the hole number you wish to use, 1 being next to home on the right
and 6 being furthest to left on normal game.

If your last stone lanes in your home, you get another turn,
otherwise the board is displayed ready for player 2.

#### automated initial moves

Moves can be entered as part of the command-line and are played before
the repl is entered 

### stealing

If your last stone lands in an empty hole and opposite a hole
which is occupied, you gain your single stone and _all_ the opposite stones.
Your turn is over with a steal.

### players

As an alternative to repl mode, you can specify a player type with a *-t*.
Players of various types will be developed but currently we have

* console - input is needed just like repl
* random - a valid random hole is chosen.

Thus we start to have the games played automatically.

Any moves on the command line are still played first.

