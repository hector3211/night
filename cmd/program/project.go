package program

import (
	"log"
	"os"

	"github.com/hector3211/night/cmd/flags"

	tea "github.com/charmbracelet/bubbletea"
)

type Project struct {
	DBDriver      flags.DataBaseDriver
	FilePath      flags.File
	ConnectionUrl flags.UrlConnection
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
