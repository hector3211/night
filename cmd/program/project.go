package program

import (
	"log"
	"night/cmd/flags"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type SeedLanguage int

const (
	UNKNWON SeedLanguage = iota
	GOLANG
	SQL
)

var DefaultSeedLanguage = SQL

func (s SeedLanguage) String() string {
	switch s {
	case GOLANG:
		return "Go"
	case SQL:
		return "SQL"
	default:
		return "SQL"
	}
}

type Project struct {
	DBDriver flags.DataBaseDriver
	FilePath flags.File
	SeedLanguage
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
