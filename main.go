package main

// An example program demonstrating the pager component from the Bubbles
// component library.

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().
const useHighPerformanceRenderer = false

type JsonData struct {
    Basics struct {
        Name     string `json:"name"`
        Headline string `json:"headline"`
        Email    string `json:"email"`
        Phone    string `json:"phone"`
        Location string `json:"location"`
        Url      struct {
            Label string `json:"label"`
            Href  string `json:"href"`
        } `json:"url"`
    } `json:"basics"`
    Sections struct {
        Summary struct {
            Name          string `json:"name"`
            Columns       int    `json:"columns"`
            SeparateLinks bool   `json:"separateLinks"`
            Visible       bool   `json:"visible"`
            ID            string `json:"id"`
            Content       string `json:"content"`
        } `json:"summary"`
    } `json:"sections"`
}

type model struct {
	content  JsonData
	ready    bool
	viewport viewport.Model
}

func (m model) Init() tea.Cmd {
	
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)


	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.contentView())
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render(m.content.Basics.Name)
	subtitle := m.content.Basics.Headline

	titleWithSubtitle := titleBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Top, title, subtitle))

	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, titleWithSubtitle, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m  model) contentView() string {
	summaryContent := m.content.Sections.Summary.Content
	contactInfoItems := []string{m.content.Basics.Location, m.content.Basics.Email, m.content.Basics.Phone}

	for i, item := range contactInfoItems {
		if(i != len(contactInfoItems) - 1) {
			contactInfoItems[i] = contactInfoItemStyle.Render(item)
		} else {
			contactInfoItems[i] = item
		}
	}


	contactInfo := lipgloss.JoinHorizontal(lipgloss.Center, contactInfoItems...)

	return contactInfoStyle.Render(contactInfo) + "\n\n" + summaryContent
}


func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Load some text for our viewport
	jsonContent, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	var jsonData JsonData
	err = json.Unmarshal(jsonContent, &jsonData)

	if err != nil {
		fmt.Println("error unmarshaling JSON:", err)
		os.Exit(1)
	}

	



	p := tea.NewProgram(
		model{content: jsonData},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err) 
		os.Exit(1)
	}
}
