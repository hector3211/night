package ui

import (
	"errors"
	"fmt"
	"night/cmd/program"
	"night/cmd/utils"
	"regexp"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Background(lipgloss.Color("#15F5BA")).Foreground(lipgloss.Color("#F0F3FF")).Bold(true).Padding(0, 1, 0)
	errorStyle  = lipgloss.NewStyle().Background(lipgloss.Color("#ff757f")).Bold(true).Padding(0, 0, 0)
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#836FFF"))
)

type SelectedUrl struct {
	url string
}

func (s *SelectedUrl) Update(value string) {
	s.url = value
}

func (s *SelectedUrl) GetValue() string {
	return s.url
}

type model struct {
	header        string
	input         textinput.Model
	connectionUrl *SelectedUrl
	err           error
	exit          *bool
}

func validateInput(input string) error {
	pattern := `^[a-zA-Z0-9_\-:@/\\]+$`
	regex := regexp.MustCompile(pattern)

	matched := regex.MatchString(input)
	if !matched {
		return fmt.Errorf("Input is invalid failed")
	}

	return nil
}

func InitializeConnModel(selection *SelectedUrl, program *program.Project) model {
	input := textinput.New()
	input.Focus()
	input.CharLimit = 250
	input.Width = 200
	input.Validate = validateInput
	input.Cursor.Style.Foreground(lipgloss.Color("#836FFF"))
	header := "Provide your postgres url to connect."
	return model{
		header:        titleStyle.Render(header),
		input:         input,
		connectionUrl: selection,
		exit:          &program.Exit,
	}
}

func CreateErrorInputModel(err error) model {
	input := textinput.New()
	input.Focus()
	input.CharLimit = 150
	input.Width = 100
	exit := true

	return model{
		header:        "",
		input:         input,
		connectionUrl: nil,
		err:           errors.New(errorStyle.Render(err.Error())),
		exit:          &exit,
	}
}

func (m model) Err() string {
	return m.err.Error()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if len(m.input.Value()) > 1 {
				m.connectionUrl.Update(m.input.Value())
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			*m.exit = true
			return m, tea.Quit
		}
	case utils.ErrMsg:
		{
			m.err = msg
			*m.exit = true
			return m, nil
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := fmt.Sprintf("%s\n\n%s\n\n", m.header, m.input.View())
	return s
}
