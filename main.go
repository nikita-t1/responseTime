package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := Styles{}
	s.BorderColor = "#00ff80"
	s.InputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1).Width(112)
	return &s
}

type model struct {
	inputField textinput.Model
	styles     *Styles
}

func initialModel() model {
	m := model{}

	ti := textinput.New()
	ti.Placeholder = "https://google.com"
	ti.Focus()

	m.styles = DefaultStyles()
	m.inputField = ti
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	m.inputField, cmd = m.inputField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.styles.InputField.Render(m.inputField.View())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
