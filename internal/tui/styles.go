package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14")) // cyan

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")) // dim

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	greenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	redStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	yellowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11"))

	labelStyle = lipgloss.NewStyle().
			Width(18).
			Foreground(lipgloss.Color("7"))

	buttonStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2).
			Background(lipgloss.Color("14")).
			Foreground(lipgloss.Color("0"))

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("8"))
)
