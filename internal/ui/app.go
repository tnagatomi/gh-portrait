package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tnagatomi/gh-portrait/internal/github"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))

	infoStyle = lipgloss.NewStyle()

	dividerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
)

// Model represents the main application UI model
type Model struct {
	user     *github.User
	viewport viewport.Model
	ready    bool
}

// Start initializes and starts the TUI application
func Start(user *github.User) error {
	m := New(user)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// New creates a new Model instance
func New(user *github.User) Model {
	return Model{
		user:  user,
		ready: false,
	}
}

// Init initializes the Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles UI updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.YPosition = 0
			m.viewport.SetContent(m.createContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "\nInitializing..."
	}
	return m.viewport.View()
}

// createContent generates the content to be displayed
func (m Model) createContent() string {
	var content string

	// Info section
	content += titleStyle.Render("Info") + "\n"
	content += infoStyle.Render("Name: " + m.user.Name) + "\n"
	if m.user.Bio != "" {
		content += infoStyle.Render("Bio: " + m.user.Bio) + "\n"
	}
	if m.user.Pronouns != "" {
		content += infoStyle.Render("Pronouns: " + m.user.Pronouns) + "\n"
	}
	if m.user.Company != "" {
		content += infoStyle.Render("Company: " + m.user.Company) + "\n"
	}
	if m.user.Location != "" {
		content += infoStyle.Render("Location: " + m.user.Location) + "\n"
	}
	if m.user.WebsiteURL != "" {
		content += infoStyle.Render("Website: " + m.user.WebsiteURL) + "\n"
	}
	content += "\n"

	// Social accounts section
	if len(m.user.Social) > 0 {
		content += titleStyle.Render("Social accounts") + "\n"
		for _, account := range m.user.Social {
			content += infoStyle.Render(
				fmt.Sprintf("%s: %s",
					account.Provider,
					account.URL,
				),
			) + "\n"
		}
	}

	// Help section
	content += "\n" + dividerStyle.Render("Press q to quit") + "\n"

	return content
}