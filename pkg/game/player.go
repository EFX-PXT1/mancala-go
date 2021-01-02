package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"
)

// Player is the general base interface
type Player interface {
	Move(pos *Position) int
}

// RandomPlayer picks a valid random move
type RandomPlayer struct{}

func newRandomPlayer(conf map[string]string) (Player, error) {
	return &RandomPlayer{}, nil
}

// Move chooses a valid random move
func (p *RandomPlayer) Move(pos *Position) (hole int) {
	moves := pos.ValidMoves()
	i := rand.Intn(len(moves))
	hole = moves[i]
	fmt.Printf("random > %d\n", hole)
	return
}

// ConsolePlayer gets a value from the console
type ConsolePlayer struct {
	Name string
}

func newConsolePlayer(conf map[string]string) (Player, error) {
	return &ConsolePlayer{
		Name: conf["name"],
	}, nil
}

// Move chooses a valid random move
func (p *ConsolePlayer) Move(pos *Position) int {
	fmt.Printf("%s > ", p.Name)
	moves := pos.ValidMoves()
	reader := bufio.NewReader(os.Stdin)
	for {
		x, _ := reader.ReadString('\n')
		x = strings.TrimRight(x, "\r\n")
		if hole, err := strconv.Atoi(x); err == nil {
			// check value is valid
			for _, m := range moves {
				if hole == m {
					// valid move
					return hole
				}
			}
		} else {
			fmt.Printf(" %s not valid\n---\n", x)
		}
	}
}

// PlayerFactory types a function which takes a map and creates a Player
type PlayerFactory func(conf map[string]string) (Player, error)

var playerFactories = make(map[string]PlayerFactory)

// RegisterPlayer registers a name to a PlayerFactory
func RegisterPlayer(name string, factory PlayerFactory) {
	if factory == nil {
		log.Panicf("Player factory %s does not exist.", name)
	}
	_, registered := playerFactories[name]
	if registered {
		log.Errorf("Player factory %s already registered. Ignoring.", name)
	}
	playerFactories[name] = factory
}

// CreatePlayer creates a Player by registered name
func CreatePlayer(conf map[string]string) (Player, error) {
	// Query configuration for datastore defaulting to "memory".
	playerName := conf["type"]

	playerFactory, ok := playerFactories[playerName]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availablePlayers := make([]string, 0, len(playerFactories))
		for k := range playerFactories {
			availablePlayers = append(availablePlayers, k)
		}
		return nil, fmt.Errorf("Invalid Player name. Must be one of: %s", strings.Join(availablePlayers, ", "))
	}

	// Run the factory with the configuration.
	return playerFactory(conf)
}

func init() {
	RegisterPlayer("random", newRandomPlayer)
	RegisterPlayer("console", newConsolePlayer)
}
