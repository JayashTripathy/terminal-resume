package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

type URL struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

type ExperienceItem struct {
	ID       string `json:"id"`
	Visible  bool   `json:"visible"`
	Company  string `json:"company"`
	Position string `json:"position"`
	Location string `json:"location"`
	Date     string `json:"date"`
	Summary  string `json:"summary"`
	Url      URL    `json:"url"`
}

type EducationItem struct {
	ID          string `json:"id"`
	Visible     bool   `json:"visible"`
	Institution string `json:"institution"`
	StudyType   string `json:"studyType"`
	Area        string `json:"area"`
	Score       string `json:"score"`
	Date        string `json:"date"`
	Summary     string `json:"summary"`
	URL         URL    `json:"url"`
}

type Education struct {
	Name          string          `json:"name"`
	Columns       int             `json:"columns"`
	SeparateLinks bool            `json:"separateLinks"`
	Visible       bool            `json:"visible"`
	ID            string          `json:"id"`
	Items         []EducationItem `json:"items"`
}

type Experience struct {
	Name          string           `json:"name"`
	Columns       int              `json:"columns"`
	SeparateLinks bool             `json:"separateLinks"`
	Visible       bool             `json:"visible"`
	ID            string           `json:"id"`
	Items         []ExperienceItem `json:"items"`
}

type SkillItem struct {
    ID       string `json:"id"`
    Visible  bool   `json:"visible"`
    Name     string `json:"name"`
    Level    int `json:"level"`
    Keywords []string `json:"keywords"`
}

type Skill struct {
    Name          string      `json:"name"`
    Columns       int         `json:"columns"`
    SeparateLinks bool        `json:"separateLinks"`
    Visible       bool        `json:"visible"`
    ID            string      `json:"id"`
    Items         []SkillItem `json:"items"`
}


type JsonData struct {
	Basics struct {
		Name     string `json:"name"`
		Headline string `json:"headline"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Location string `json:"location"`
		Url      URL    `json:"url"`
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
		Experience Experience `json:"experience"`
		Education  Education  `json:"education"`
		Skill      Skill      `json:"skills"`
	} `json:"sections"`
}
type model struct {
	content  JsonData
	ready    bool
	viewport viewport.Model
	msg      tea.Msg
}

func (m model) Init() tea.Cmd {

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.msg = msg
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

func (m model) AboutSection() string {
	summaryContent := sectionContentStyle.Render(m.content.Sections.Summary.Content)
	summaryTitle := sectionTitleStyle.Render(m.content.Sections.Summary.Name)
	return lipgloss.JoinVertical(lipgloss.Left, summaryTitle, summaryContent)
}

func (m model) ExperienceItem(item ExperienceItem) string {
	company := item.Company
	position := item.Position
	location := item.Location
	date := item.Date
	summary := sectionContentStyle.PaddingTop(1).Render(item.Summary)

	titlePositionBlock := lipgloss.NewStyle().
		Width(m.viewport.Width / 2).Align(lipgloss.Left).Render(lipgloss.JoinVertical(lipgloss.Left, company, position))
	locationDateBlock := lipgloss.NewStyle().
		Width(m.viewport.Width / 2).Align(lipgloss.Right).Render(lipgloss.JoinVertical(lipgloss.Right, location, date))
	experienceItemHeader := lipgloss.JoinHorizontal(lipgloss.Left, titlePositionBlock, locationDateBlock)

	return lipgloss.JoinVertical(lipgloss.Top, experienceItemHeader, summary) + "\n"
}

func (m model) EducationItem(item EducationItem) string {
	institution := item.Institution
	studyType := item.StudyType
	area := item.Area
	date := item.Date
	summary := sectionContentStyle.PaddingTop(1).Render(item.Summary)

	titlePositionBlock := lipgloss.NewStyle().
		Width(m.viewport.Width / 2).Align(lipgloss.Left).Render(lipgloss.JoinVertical(lipgloss.Left, institution, studyType))
	locationDateBlock := lipgloss.NewStyle().
		Width(m.viewport.Width / 2).Align(lipgloss.Right).Render(lipgloss.JoinVertical(lipgloss.Right, area, date))
	educationItemHeader := lipgloss.JoinHorizontal(lipgloss.Left, titlePositionBlock, locationDateBlock)

	return lipgloss.JoinVertical(lipgloss.Top, educationItemHeader, summary) + "\n"
}

func (m model) ExperienceSection() string {
	title := sectionTitleStyle.Render(m.content.Sections.Experience.Name)
	items := []string{}

	for _, item := range m.content.Sections.Experience.Items {
		items = append(items, m.ExperienceItem(item))
	}

	return lipgloss.JoinVertical(lipgloss.Top, title, lipgloss.JoinVertical(lipgloss.Top, items...))
}

func (m model) EducationSection() string {
	title := sectionTitleStyle.Render(m.content.Sections.Education.Name)
	items := []string{}

	for _ , item := range m.content.Sections.Education.Items {
		items = append(items, m.EducationItem(item))
	}
	return lipgloss.JoinVertical(lipgloss.Top, title) + "\n" + lipgloss.JoinVertical(lipgloss.Top, items...)

}

func (m model) SkillSection() string {
	title := sectionTitleStyle.Render(m.content.Sections.Skill.Name)
	items :=  make([]string, 0, len(m.content.Sections.Skill.Items))
	
	itemStyle := lipgloss.NewStyle().Width(m.viewport.Width / 5).Align(lipgloss.Left).PaddingBottom(1)


	for _, item := range m.content.Sections.Skill.Items {
		items = append(items, itemStyle.Render(item.Name))
	}
	
	rows := (len(items) + 4) / 5
	
	formattedRows := make([]string, 0, rows)

	for i := 0; i < len(items); i += 5 {
		end := i + 5
		if end > len(items) {
			end = len(items)
		}
		formattedRows = append(formattedRows, lipgloss.JoinHorizontal(lipgloss.Top, items[i:end]...))
	}

	return lipgloss.JoinVertical(lipgloss.Top, title, lipgloss.JoinHorizontal(lipgloss.Left, formattedRows...) )
}

func (m model) contentView() string {

	contactInfoItems := []string{m.content.Basics.Location, m.content.Basics.Email, m.content.Basics.Phone}
	for i, item := range contactInfoItems {
		if i != len(contactInfoItems)-1 {
			contactInfoItems[i] = contactInfoItemStyle.Render(item)
		} else {
			contactInfoItems[i] = item
		}
	}
	contactInfo := lipgloss.JoinHorizontal(lipgloss.Center, contactInfoItems...)
	return contactInfoStyle.Render(contactInfo) + "\n\n" +
								 m.AboutSection() + "\n\n" + 
								 m.ExperienceSection() + "\n\n" +
								 m.EducationSection() + "\n\n" +
								 m.SkillSection() + "\n\n" 
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
