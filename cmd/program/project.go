package program

import (
	"github.com/hector3211/night/cmd/flags"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Project struct {
	DBDriver flags.DataBaseDriver
	FilePath flags.File
	// flags.SeedLanguage
	ConnectionUrl string
	Exit          bool
}

func (p *Project) ExitCli(tprogram *tea.Program) {
	if p.Exit {
		if err := tprogram.ReleaseTerminal(); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}
