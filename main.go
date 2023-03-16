package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "Hello World"
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
