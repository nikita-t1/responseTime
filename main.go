package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
)

type model struct {
	inputField textinput.Model
	feedback   string
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
		case "enter":
			m.feedback = "Loading"
			_, err := ExecuteRequest(m.inputField.Value())
			if err != nil {
				m.feedback = m.styles.Error.Render(err.Error())
			} else {
				m.feedback = m.styles.Success.Render("Successfully")
			}
		}
	}
	m.inputField, cmd = m.inputField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.styles.InputField.Render(m.inputField.View()),
		m.feedback,
	)
}

func main() {
	f, _ := tea.LogToFile("debug.log", "")
	defer f.Close()
	log.Printf("Start")

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
