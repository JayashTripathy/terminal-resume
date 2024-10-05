package main

import "github.com/charmbracelet/lipgloss"

var (
	//hex codes
	colors = map[string]string{
		"primary": "#C77DFF",
		"secondary": "#F8F8F8", // off white
		"foreground": "#F8F8F2", // white
		"mutedForeground": "#BEBEBE",
		
	}

	titleBoxStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.MiddleRight = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
	titleStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colors["primary"]))
	}()

	sectionTitleStyle = func () lipgloss.Style {
		return lipgloss.NewStyle().Bold(true).PaddingBottom(1).Foreground(lipgloss.Color(colors["foreground"]))
	}()

	sectionContentStyle = func(m model) lipgloss.Style {
		return lipgloss.NewStyle().Width(m.viewport.Width - bodyPadding).Foreground(lipgloss.Color(colors["mutedForeground"]))
	}

	contactInfoStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colors["secondary"]))
	}()
	contactInfoItemStyle = func() lipgloss.Style {
		border := lipgloss.Border{
			Right: "|",
		}
		return lipgloss.NewStyle().
						BorderStyle(border).
						PaddingRight(1).
						MarginRight(1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleBoxStyle.BorderStyle(b).Foreground(lipgloss.Color(colors["primary"]))
	}() 
	
)