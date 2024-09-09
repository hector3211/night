package ui

import (
	"github.com/hector3211/night/cmd/program"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	prgressTitleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#15F5BA")).Foreground(lipgloss.Color("#F0F3FF")).Bold(true).Padding(0, 1, 0)
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
)

type tickMsg time.Time

type ProgressModel struct {
	// header   string
	percent  float64
	progress progress.Model
	exit     *bool
}

func InitializeProgressModel(program *program.Project) ProgressModel {
	pg := progress.New(progress.WithScaledGradient("#15F5BA", "#836FFF"))
	return ProgressModel{
		// header:   header,
		percent:  0.0,
		progress: pg,
		exit:     &program.Exit,
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return tickCmd()
}
func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			*m.exit = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		m.percent += 0.50
		if m.percent > 1.0 {
			m.percent = 1.0
			return m, tea.Batch(tea.ClearScreen, tea.Quit)
		}
		return m, tickCmd()

	default:
		return m, nil
	}
	return m, nil
}

func (m ProgressModel) View() string {
	// s := prgressTitleStyle.Render(m.header)
	s := ""
	// pad := strings.Repeat(" ", padding)
	s += m.progress.ViewAs(m.percent) + "\n\n"
	s += helpStyle.Render("Press [esc] [ctrl-c] [q] to quit")

	return s
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
