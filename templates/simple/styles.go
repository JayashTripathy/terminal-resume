package templates

import (
	"github.com/charmbracelet/lipgloss"
)



var (
	//hex codes
	colors = map[string]string{
		"primary": "198",
		"secondary": "255", 
		"foreground": "231", 
		"mutedForeground": "249",
		
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

	sectionContentStyle = func(m SimpleModel) lipgloss.Style {
		return lipgloss.NewStyle().Width(m.Viewport.Width - bodyPadding).Foreground(lipgloss.Color(colors["mutedForeground"]))
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