package main

import "github.com/charmbracelet/lipgloss"

var (
	//hex codes
	colors = map[string]string{
		"primary": "#C77DFF",
	}

	titleBoxStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.MiddleRight = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
	titleStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Bold((true)).Foreground(lipgloss.Color(colors["primary"]))
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleBoxStyle.BorderStyle(b).Foreground(lipgloss.Color(colors["primary"]))
	}() 
)