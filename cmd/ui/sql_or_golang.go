package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"night/cmd/program"
	"night/cmd/utils"
)

var (
	TitleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#15F5BA")).Foreground(lipgloss.Color("#F0F3FF")).Bold(true).Padding(0, 1, 0)
	ErrorStyle = lipgloss.NewStyle().Background(lipgloss.Color("#ff757f")).Bold(true).Padding(0, 0, 0)
)

//	type SelectedFilePath struct {
//		Choice string
//	}
//
//	func (s *SelectedFilePath) Update(value string) {
//		s.Choice = value
//	}
type SelectedLanguage struct {
	Choice string
}

func (s *SelectedLanguage) Update(value string) {
	s.Choice = value
}

func (s *SelectedLanguage) GetValue() string {
	return s.Choice
}

type ChoiceModel struct {
	header   string
	choices  []string       // items on the to-do list
	cursor   int            // which to-do list item our cursor is pointing at
	selected map[int]string // which to-do items are selectedfilePath   *SelectedFilePath
	language *SelectedLanguage
	err      error
	exit     *bool
}

func InitializeChoiceModel(selection *SelectedLanguage, program *program.Project) ChoiceModel {
	header := "Did you write your data in Go or SQL?"
	return ChoiceModel{
		header:   sqlFileTitleStyle.Render(header),
		choices:  []string{"Go", "SQL"},
		selected: make(map[int]string, 0),
		exit:     &program.Exit,
	}
}

func (f ChoiceModel) Init() tea.Cmd {
	return nil
}

func (f ChoiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "esc":
			*f.exit = true
			return f, tea.Quit

			// The "up" and "k" keys move the cursor up
		case "up", "k":
			if f.cursor > 0 {
				f.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if f.cursor < len(f.choices)-1 {
				f.cursor++
			}
		case "enter", " ":
			if selected, ok := f.selected[f.cursor]; ok != false {
				f.language.Update(selected)
			}
		}
	case utils.ErrMsg:
		f.err = msg
		*f.exit = true
		return f, nil
	}

	return f, cmd
}

func (f ChoiceModel) View() string {
	s := fmt.Sprintf("\n%s\n\n", f.header)
	// Iterate over our choices
	for i, choice := range f.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if f.cursor == i {
			cursor = "▶️" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
