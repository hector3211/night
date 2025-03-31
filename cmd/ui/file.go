package ui

import (
	"fmt"
	"os"

	"github.com/hector3211/night/cmd/program"
	"github.com/hector3211/night/cmd/utils"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	sqlFileTitleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#15F5BA")).Foreground(lipgloss.Color("#F0F3FF")).Bold(true).Padding(0, 1, 0)
	sqlFileErrorStyle = lipgloss.NewStyle().Background(lipgloss.Color("#ff757f")).Bold(true).Padding(0, 0, 0)
)

type SelectedFilePath struct {
	Choice string
}

func (s *SelectedFilePath) Update(value string) {
	s.Choice = value
}

type FileModel struct {
	header     string
	filePicker filepicker.Model
	filePath   *SelectedFilePath
	err        error
	exit       *bool
}

func InitializeFileModel(selection *SelectedFilePath, program *program.Project) FileModel {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".sql", ".go"}
	fp.CurrentDirectory, _ = os.Getwd()
	header := "Search for your seed file"
	return FileModel{
		header:     sqlFileTitleStyle.Render(header),
		filePicker: fp,
		filePath:   selection,
		exit:       &program.Exit,
	}
}

func (f FileModel) Init() tea.Cmd {
	return f.filePicker.Init()
}

func (f FileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "esc", "q":
			*f.exit = true
			return f, tea.Quit
		}
	case utils.ErrMsg:
		f.err = msg
		*f.exit = true
		return f, nil
	}

	f.filePicker, cmd = f.filePicker.Update(msg)

	if didSelect, path := f.filePicker.DidSelectFile(msg); didSelect {
		f.filePath.Update(path)
		return f, tea.Batch(tea.ClearScreen, tea.Quit)
	}

	return f, cmd
}

func (f FileModel) View() string {
	s := fmt.Sprintf("\n%s\n\n", f.header)
	s += f.filePicker.View()
	return s
}
