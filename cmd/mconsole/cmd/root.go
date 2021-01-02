package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/EFX-PXT1/mancala-go/pkg/game"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var width int
var stones int
var showDelta bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mconsole",
	Short: "Mancala Console",
	Long: `Mancala in the Console
enables development and playing positions.
For example:

mconsole --width <width> --stones <start stones> ...moves`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		game.DefineGame(
			viper.GetInt("game.width"),
			viper.GetInt("game.stones"),
		)

		pos := game.StartPosition()
		pos.Show()

		var x string
		for len(args) > 0 {
			x, args = args[0], args[1:]
			if hole, err := strconv.Atoi(x); err == nil {
				var delta *game.Position
				if pos, delta, _, err = pos.Move(hole); err == nil {
					if viper.GetBool("show.delta") {
						delta.Show()
					}
					pos.Show()
				} else {
					fmt.Printf(" %s %s\n---\n", x, err)
				}
			}
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
	rootCmd.Flags().BoolVar(&showDelta, "delta", false, "show delta position")

	viper.BindPFlag("game.width", rootCmd.Flags().Lookup("width"))
	viper.BindPFlag("game.stones", rootCmd.Flags().Lookup("stones"))
	viper.BindPFlag("show.delta", rootCmd.Flags().Lookup("delta"))
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
