package ui

import (
	"fmt"
	"night/cmd/program"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Flag, Title, Desc string
}

var Items = []Item{
	{
		Title: "Sqlite",
		Desc:  "Sqlite3 Driver",
		Flag:  "sqlite3",
	},
	{
		Title: "Postgres",
		Desc:  "Postgres Driver",
		Flag:  "postgres",
	},
}

var (
	driverTitleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#15F5BA")).Foreground(lipgloss.Color("#F0F3FF")).Bold(true).Padding(0, 1, 0)
	driverErrorStyle = lipgloss.NewStyle().Background(lipgloss.Color("#ff757f")).Bold(true).Padding(0, 0, 0)
	driverCursorStle = lipgloss.NewStyle().Foreground(lipgloss.Color("#836FFF"))
)

type SelectedDriver struct {
	Choice string
}

func (s *SelectedDriver) Update(value string) {
	s.Choice = value
}

type DriverModel struct {
	header   string
	cursor   int
	choices  []Item
	selected map[int]struct{}
	choice   *SelectedDriver
	exit     *bool
}

func (m DriverModel) Init() tea.Cmd {
	return nil
}

func InitialModel(selection *SelectedDriver, program *program.Project) DriverModel {
	header := "Select which database you'll be using"
	return DriverModel{
		header:   driverTitleStyle.Render(header),
		choices:  Items,
		selected: make(map[int]struct{}),
		choice:   selection,
		exit:     &program.Exit,
	}
}

func (m DriverModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			*m.exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.choice.Update(m.choices[m.cursor].Flag)
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DriverModel) View() string {
	s := fmt.Sprintf("\n\n%s\n\n", m.header)
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		// checked := " "
		// if _, ok := m.selected[i]; ok {
		// 	checked = "x"
		// }

		s += fmt.Sprintf("%s %s\n%s\n\n", driverCursorStle.Render(cursor), choice.Title, choice.Desc)
	}
	return s
}
