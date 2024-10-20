package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0" // Listen on all interfaces
	port = "2222"   // Default SSH port

	// global style variables

	bodyPadding = 1
)

const useHighPerformanceRenderer = false

type model struct {
	content  JsonData
	ready    bool
	viewport viewport.Model
	sess     ssh.Session
	msg      tea.Msg
	user string
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
	fmt.Printf("user %+v\n", m.user)

	var view string

	if(m.user != ""){
		view = fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	}else{
		view = "No user provided"
	}
	

	return view
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
	summaryContent := sectionContentStyle(m).Render(m.content.Sections.Summary.Content)
	summaryTitle := sectionTitleStyle.Render(m.content.Sections.Summary.Name)
	return lipgloss.JoinVertical(lipgloss.Left, summaryTitle, summaryContent)
}

func (m model) ExperienceItem(item ExperienceItem, isLast bool) string {
	company := item.Company
	position := item.Position
	location := item.Location
	date := item.Date
	summary := sectionContentStyle(m).PaddingTop(1).Render(item.Summary)
	halfSafeAreaWidth := (m.viewport.Width / 2) - bodyPadding

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

func (m model) EducationItem(item EducationItem, isLast bool) string {
	institution := item.Institution
	studyType := item.StudyType
	area := item.Area
	date := item.Date
	summary := sectionContentStyle(m).PaddingTop(1).Render(item.Summary)
	halfSafeAreaWidth := (m.viewport.Width / 2) - bodyPadding
	


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

func (m model) ExperienceSection() string {
	title := sectionTitleStyle.Render(m.content.Sections.Experience.Name)
	items := []string{}

	for index, item := range m.content.Sections.Experience.Items {
		items = append(items, m.ExperienceItem(item, index == len(m.content.Sections.Experience.Items)-1))
	}

	return lipgloss.JoinVertical(lipgloss.Top, title, lipgloss.JoinVertical(lipgloss.Top, items...))
}

func (m model) EducationSection() string {
	title := sectionTitleStyle.Render(m.content.Sections.Education.Name)
	items := []string{}

	for index, item := range m.content.Sections.Education.Items {
		items = append(items, m.EducationItem(item, index == len(m.content.Sections.Education.Items)-1))
	}
	return lipgloss.JoinVertical(lipgloss.Top, title) + "\n" + lipgloss.JoinVertical(lipgloss.Top, items...)

}

func (m model) SkillSection() string {
	title := sectionTitleStyle.Render(m.content.Sections.Skill.Name)
	items := make([]string, 0, len(m.content.Sections.Skill.Items))

	colLen := 5

	if m.content.Sections.Skill.Columns > 0 {
		colLen = m.content.Sections.Skill.Columns
	}

	// itemsContainerStyle := lipgloss.NewStyle()
	itemStyle := lipgloss.NewStyle().Width(m.viewport.Width / colLen).Align(lipgloss.Left).PaddingBottom(1)

	for _, item := range m.content.Sections.Skill.Items {
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

func (m model) ProjectSection() string {
	title := sectionTitleStyle.Render(m.content.Sections.Projects.Name)
	projects := make([]string, 0, len(m.content.Sections.Projects.Items))
	halfSafeAreaWidth := (m.viewport.Width / 2) - bodyPadding


	for index, project := range m.content.Sections.Projects.Items {
		projectName := project.Name
		projectSummary := project.Summary
		projectUrl := project.URL

		projectNameStyle := lipgloss.NewStyle().Width(halfSafeAreaWidth).Align(lipgloss.Left).Render(projectName)
		projectUrlStyle := lipgloss.NewStyle().Width(halfSafeAreaWidth).Align(lipgloss.Right).Render(projectUrl.Href)

		projectHeader := lipgloss.JoinHorizontal(lipgloss.Left, projectNameStyle, projectUrlStyle)
		projectDescriptionStyle := sectionContentStyle(m).PaddingTop(1).Render(projectSummary)

		projectItem := lipgloss.JoinVertical(lipgloss.Top, projectHeader, projectDescriptionStyle)

		isLast := index == len(m.content.Sections.Projects.Items)-1
		projects = append(projects, projectItem)
		if !isLast {
			projects = append(projects, "\n")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Top, title) + "\n" + lipgloss.JoinVertical(lipgloss.Top, projects...)
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
	return lipgloss.NewStyle().Padding(bodyPadding).Render(lipgloss.JoinHorizontal(lipgloss.Left, contactInfoStyle.Render(contactInfo)+"\n\n"+
		m.AboutSection()+"\n\n\n"+
		m.ExperienceSection()+"\n\n\n"+
		m.ProjectSection()+"\n\n\n"+
		m.EducationSection()+"\n\n\n"+
		m.SkillSection()+"\n\n\n"))
}

//go:embed data.json
var jsonContent []byte

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {

	var jsonData JsonData
	err := json.Unmarshal(jsonContent, &jsonData)
	var username string

	if err != nil {
		fmt.Println("error unmarshaling JSON:", err)
		os.Exit(1)
	}

	err = extractUsername(s, &username)

	if err != nil {
		log.Warn(err)
	}

	log.Info("Username is " + username)


	m := model{
		sess:    s,
		content: jsonData,
		user: username,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}



func main() {

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		// Allocate a pty.
		// This creates a pseudoconsole on windows, compatibility is limited in
		// that case, see the open issues for more details.
		ssh.AllocatePty(),
		wish.WithMiddleware(
			// run our Bubble Tea handler
			bubbletea.Middleware(teaHandler),
			// ensure the user has requested a tty
			activeterm.Middleware(),
			logging.Middleware(),

		),
		ssh.HostKeyFile(".ssh/id_ed25519"),
	)

	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}

}
