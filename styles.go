package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	InputField        lipgloss.Style
	Table             table.Styles
	Error             lipgloss.Style
	Warning           lipgloss.Style
	Success           lipgloss.Style
	InactiveTabBorder lipgloss.Border
	ActiveTabBorder   lipgloss.Border
	DocStyle          lipgloss.Style
	HighlightColor    lipgloss.AdaptiveColor
	InactiveTabStyle  lipgloss.Style
	ActiveTabStyle    lipgloss.Style
	WindowStyle       lipgloss.Style
	VerySubduedColor  lipgloss.AdaptiveColor
	SubduedColor      lipgloss.AdaptiveColor
}

func DefaultStyles() *Styles {
	s := Styles{}
	s.InputField = lipgloss.NewStyle().
		BorderForeground(lipgloss.AdaptiveColor{
			Light: "63",
			Dark:  "63",
		}).
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1).Width(100)
	s.Error = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "203",
		Dark:  "204",
	})
	s.Success = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "39",
		Dark:  "86",
	})
	s.Warning = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "208",
		Dark:  "192",
	})

	s.Table = table.DefaultStyles()
	s.Table.Selected = s.Table.Selected.
		//Foreground(lipgloss.Color("#FFFFFF")).
		//Foreground(lipgloss.Color("229")).
		//Background(lipgloss.Color("17")).
		UnsetForeground().
		UnsetBackground().
		Bold(true)

	s.InactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	s.ActiveTabBorder = tabBorderWithBottom("┘", " ", "└")
	s.DocStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	s.HighlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	s.InactiveTabStyle = lipgloss.NewStyle().Border(s.InactiveTabBorder, true).BorderForeground(s.HighlightColor).Padding(0, 1)
	s.ActiveTabStyle = s.InactiveTabStyle.Copy().Border(s.ActiveTabBorder, true)
	s.WindowStyle = lipgloss.NewStyle().BorderForeground(s.HighlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()

	s.VerySubduedColor = lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	s.SubduedColor = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	return &s
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
