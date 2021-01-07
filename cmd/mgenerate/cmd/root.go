package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/EFX-PXT1/mancala-go/pkg/game"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var width int
var stones int
var filename string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mgenerate",
	Short: "Mancala Position Generator",
	Long: `Mancala Position Generator
creates all possible position for a mancala game
For example:

mgenerate --width <width> --stones <start stones> --file <filename>`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		game.DefineGame(
			viper.GetInt("game.width"),
			viper.GetInt("game.stones"),
		)

		filename := viper.GetString("generator.filename")

		p := game.StartPosition()

		// initial seed position and moves
		set := make(map[string]bool) // empty set of processed or to process
		moves := p.ValidMoves()
		for _, m := range moves {
			set[createKey(p, m)] = false
		}

		file, err := os.Create(filename)
		if err != nil {
			return
		}
		defer file.Close()

		count := 0

		// create a todo slice
		todo := toGenerate(set)

		for len(todo) > 0 {
			for _, k := range todo {
				// evaluate this single todo move
				s := strings.Split(k, ";")
				p = game.CreatePositionCsv(s[0])
				move, _ := strconv.Atoi(s[1])
				moves := p.ValidMoves()
				e, _, result, _ := p.Move(move)
				writeLine(file, p, moves, move, result, e)
				set[k] = true
				count = count + 1
				fmt.Printf("\r%d", count)
				// generate the todo list for the result
				if result == game.EndOfTurn {
					e = e.ChangePlayer()
				}
				moves = e.ValidMoves()
				for _, move = range moves {
					key := createKey(e, move)
					if !set[key] {
						set[key] = false
					}
				}
			}
			// create new todo slice
			todo = toGenerate(set)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mancala.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().IntVarP(&width, "width", "w", 6, "width of board")
	rootCmd.Flags().IntVarP(&stones, "stones", "s", 4, "intial number of stones")
	rootCmd.Flags().StringVarP(&filename, "filename", "f", "output.txt", "filename")

	viper.BindPFlag("game.width", rootCmd.Flags().Lookup("width"))
	viper.BindPFlag("game.stones", rootCmd.Flags().Lookup("stones"))
	viper.BindPFlag("generator.filename", rootCmd.Flags().Lookup("filename"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ksatctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mancala")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func createKey(p *game.Position, move int) string {
	return fmt.Sprintf("%s;%d", p.AsCsv(), move)
}

func writeLine(f *os.File, p *game.Position, moves []int, move int, result game.MoveResult, e *game.Position) {
	var m []string
	for _, i := range moves {
		m = append(m, strconv.Itoa(i))
	}
	s := fmt.Sprintf("%s;%s;%d;%s\n",
		createKey(p, move),
		strings.Join(m, ","),
		result,
		e.AsCsv(),
	)
	//	fmt.Printf(s)
	f.WriteString(s)
}

func toGenerate(set map[string]bool) (todo []string) {
	todo = make([]string, 0)
	for k, v := range set {
		if v == false {
			todo = append(todo, k)
		}
	}
	return
}
