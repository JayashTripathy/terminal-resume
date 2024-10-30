package templates

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"terminal-resume.jayash.space/models"
)

const (
	bodyPadding = 1
)

var useHighPerformanceRenderer = false

type SimpleModel struct {
	Content  models.JsonData
	Ready    bool
	Viewport viewport.Model
	Sess     ssh.Session
	Msg      tea.Msg
}

func (m SimpleModel) Init() tea.Cmd {
	return nil
}

func (m SimpleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.Msg = msg
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

		if !m.Ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.Viewport.YPosition = headerHeight
			m.Viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.Viewport.SetContent(m.contentView())
			m.Ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.Viewport.YPosition = headerHeight + 1
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.Viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m SimpleModel) View() string {
	if !m.Ready {
		return "\n  Initializing..."
	}
	// fmt.Printf("user %+v\n", m.user)

	var view = fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.Viewport.View(), m.footerView())

	// if(m.user != ""){
	// view = fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.Viewport.View(), m.footerView())
	// }else{
	// 	view = "No user provided"
	// }

	return view
}

func (m SimpleModel) headerView() string {
	title := titleStyle.Render(m.Content.Basics.Name)
	subtitle := m.Content.Basics.Headline

	titleWithSubtitle := titleBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Top, title, subtitle))

	line := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, titleWithSubtitle, line)
}

func (m SimpleModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m SimpleModel) AboutSection() string {
	summaryContent := sectionContentStyle(m).Render(m.Content.Sections.Summary.Content)
	summaryTitle := sectionTitleStyle.Render(m.Content.Sections.Summary.Name)
	return lipgloss.JoinVertical(lipgloss.Left, summaryTitle, summaryContent)
}

func (m SimpleModel) ExperienceItem(item models.ExperienceItem, isLast bool) string {
	company := item.Company
	position := item.Position
	location := item.Location
	date := item.Date
	summary := sectionContentStyle(m).PaddingTop(1).Render(item.Summary)
	halfSafeAreaWidth := (m.Viewport.Width / 2) - bodyPadding

	titlePositionBlock := lipgloss.NewStyle().
		Width(halfSafeAreaWidth).Align(lipgloss.Left).Render(lipgloss.JoinVertical(lipgloss.Left, company, position))
	locationDateBlock := lipgloss.NewStyle().
		Width(halfSafeAreaWidth).Align(lipgloss.Right).Render(lipgloss.JoinVertical(lipgloss.Right, location, date))
	experienceItemHeader := lipgloss.JoinHorizontal(lipgloss.Left, titlePositionBlock, locationDateBlock)

	return lipgloss.JoinVertical(lipgloss.Top, experienceItemHeader, summary) + func() string {
		if isLast {
			return ""
		}
		return "\n"
	}()
}

func (m SimpleModel) EducationItem(item models.EducationItem, isLast bool) string {
	institution := item.Institution
	studyType := item.StudyType
	area := item.Area
	date := item.Date
	summary := sectionContentStyle(m).PaddingTop(1).Render(item.Summary)
	halfSafeAreaWidth := (m.Viewport.Width / 2) - bodyPadding

	titlePositionBlock := lipgloss.NewStyle().
		Width(halfSafeAreaWidth).Align(lipgloss.Left).Render(lipgloss.JoinVertical(lipgloss.Left, institution, studyType))
	locationDateBlock := lipgloss.NewStyle().
		Width(halfSafeAreaWidth).Align(lipgloss.Right).Render(lipgloss.JoinVertical(lipgloss.Right, area, date))
	educationItemHeader := lipgloss.JoinHorizontal(lipgloss.Left, titlePositionBlock, locationDateBlock)

	return lipgloss.JoinVertical(lipgloss.Top, educationItemHeader, summary) + func() string {
		if isLast {
			return ""
		}
		return "\n"
	}()
}

func (m SimpleModel) ExperienceSection() string {
	title := sectionTitleStyle.Render(m.Content.Sections.Experience.Name)
	items := []string{}

	for index, item := range m.Content.Sections.Experience.Items {
		items = append(items, m.ExperienceItem(item, index == len(m.Content.Sections.Experience.Items)-1))
	}

	return lipgloss.JoinVertical(lipgloss.Top, title, lipgloss.JoinVertical(lipgloss.Top, items...))
}

func (m SimpleModel) EducationSection() string {
	title := sectionTitleStyle.Render(m.Content.Sections.Education.Name)
	items := []string{}

	for index, item := range m.Content.Sections.Education.Items {
		items = append(items, m.EducationItem(item, index == len(m.Content.Sections.Education.Items)-1))
	}
	return lipgloss.JoinVertical(lipgloss.Top, title) + "\n" + lipgloss.JoinVertical(lipgloss.Top, items...)

}

func (m SimpleModel) SkillSection() string {
	title := sectionTitleStyle.Render(m.Content.Sections.Skill.Name)
	items := make([]string, 0, len(m.Content.Sections.Skill.Items))

	colLen := 5

	if m.Content.Sections.Skill.Columns > 0 {
		colLen = m.Content.Sections.Skill.Columns
	}

	// itemsContainerStyle := lipgloss.NewStyle()
	itemStyle := lipgloss.NewStyle().Width(m.Viewport.Width / colLen).Align(lipgloss.Left).PaddingBottom(1)

	for _, item := range m.Content.Sections.Skill.Items {
		items = append(items, itemStyle.Render(item.Name))
	}

	// Calculate the number of rows
	cols := (len(items) + colLen - 1) / colLen

	formattedRows := make([]string, 0, cols)

	for i := 0; i < len(items); i += colLen {
		end := i + colLen
		if end > len(items) {
			end = len(items)
		}
		formattedRows = append(formattedRows, lipgloss.JoinHorizontal(lipgloss.Top, items[i:end]...))
	}

	return lipgloss.JoinVertical(lipgloss.Top, title, lipgloss.JoinVertical(lipgloss.Left, formattedRows...))
}

func (m SimpleModel) ProjectSection() string {
	title := sectionTitleStyle.Render(m.Content.Sections.Projects.Name)
	projects := make([]string, 0, len(m.Content.Sections.Projects.Items))
	halfSafeAreaWidth := (m.Viewport.Width / 2) - bodyPadding

	for index, project := range m.Content.Sections.Projects.Items {
		projectName := project.Name
		projectSummary := project.Summary
		projectUrl := project.URL

		projectNameStyle := lipgloss.NewStyle().Width(halfSafeAreaWidth).Align(lipgloss.Left).Render(projectName)
		projectUrlStyle := lipgloss.NewStyle().Width(halfSafeAreaWidth).Align(lipgloss.Right).Render(projectUrl.Href)

		projectHeader := lipgloss.JoinHorizontal(lipgloss.Left, projectNameStyle, projectUrlStyle)
		projectDescriptionStyle := sectionContentStyle(m).PaddingTop(1).Render(projectSummary)

		projectItem := lipgloss.JoinVertical(lipgloss.Top, projectHeader, projectDescriptionStyle)

		isLast := index == len(m.Content.Sections.Projects.Items)-1
		projects = append(projects, projectItem)
		if !isLast {
			projects = append(projects, "\n")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Top, title) + "\n" + lipgloss.JoinVertical(lipgloss.Top, projects...)
}

func (m SimpleModel) contentView() string {

	contactInfoItems := []string{m.Content.Basics.Location, m.Content.Basics.Email, m.Content.Basics.Phone}
	for i, item := range contactInfoItems {
		if i != len(contactInfoItems)-1 {
			contactInfoItems[i] = contactInfoItemStyle.Render(item)
		} else {
			contactInfoItems[i] = item
		}
	}
	contactInfo := lipgloss.JoinHorizontal(lipgloss.Center, contactInfoItems...)
	return lipgloss.NewStyle().Padding(bodyPadding).Render(lipgloss.JoinHorizontal(lipgloss.Left, contactInfoStyle.Render(contactInfo)+"\n\n"+
		m.AboutSection()+"\n\n\n"+
		m.ExperienceSection()+"\n\n\n"+
		m.ProjectSection()+"\n\n\n"+
		m.EducationSection()+"\n\n\n"+
		m.SkillSection()+"\n\n\n"))
}
