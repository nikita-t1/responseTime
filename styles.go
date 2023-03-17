package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	InputField lipgloss.Style
	Table      table.Styles
	Error      lipgloss.Style
	Success    lipgloss.Style
}

func DefaultStyles() *Styles {
	s := Styles{}
	s.InputField = lipgloss.NewStyle().
		BorderForeground(lipgloss.AdaptiveColor{
			Light: "63",
			Dark:  "63",
		}).
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1).Width(112)
	s.Error = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "203",
		Dark:  "204",
	})
	s.Success = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "39",
		Dark:  "86",
	})

	s.Table = table.DefaultStyles()
	s.Table.Selected = s.Table.Selected.
		//Foreground(lipgloss.Color("#FFFFFF")).
		//Foreground(lipgloss.Color("229")).
		//Background(lipgloss.Color("17")).
		UnsetForeground().
		UnsetBackground().
		Bold(true)

	return &s
}
