package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/fztea/flipperzero"
	"github.com/jon4hz/fztea/internal/version"
	"github.com/muesli/coral"
	mcoral "github.com/muesli/mango-coral"
	"github.com/muesli/roff"
)

var rootFlags struct {
	port string
}

var rootCmd = &coral.Command{
	Use:     "fztea",
	Short:   "TUI to interact with your flipper zero",
	Version: version.Version,
	Run:     root,
}

func init() {
	rootCmd.Flags().StringVarP(&rootFlags.port, "port", "p", "", "port to connect to")
	rootCmd.AddCommand(serverCmd, versionCmd, manCmd)
}

func root(cmd *coral.Command, args []string) {
	fz, err := flipperzero.NewFlipperZero(flipperzero.WithPort(rootFlags.port))
	if err != nil {
		log.Fatal(err)
	}
	m := model{
		flipper: flipperzero.New(fz),
	}
	if err := tea.NewProgram(m, tea.WithMouseCellMotion()).Start(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var manCmd = &coral.Command{
	Use:                   "man",
	Short:                 "generates the manpages",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Hidden:                true,
	Args:                  coral.NoArgs,
	RunE: func(cmd *coral.Command, args []string) error {
		manPage, err := mcoral.NewManPage(1, rootCmd)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
		return err
	},
}

var versionCmd = &coral.Command{
	Use:   "version",
	Short: "Print the version info",
	Run: func(cmd *coral.Command, args []string) {
		fmt.Printf("Version: %s\n", version.Version)
		fmt.Printf("Commit: %s\n", version.Commit)
		fmt.Printf("Date: %s\n", version.Date)
		fmt.Printf("Build by: %s\n", version.BuiltBy)
	},
}
