package ui

import (
	"context"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh/v2/pkg/browser"
	"github.com/tnagatomi/gh-portrait/internal/github"
	"github.com/tnagatomi/gh-portrait/internal/ui/components"
)

var (
	dividerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

	errorHelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))
)

// fetchRepositoriesMsg is sent when repositories are fetched
type fetchRepositoriesMsg struct {
	repositories []github.Repository
	err         error
}

// tabSelectedMsg is sent when a tab is selected
type tabSelectedMsg struct {
	index int
}

// Model represents the main application UI model
type Model struct {
	user         *github.User
	repositories []github.Repository
	tabs         components.Tabs
	repoList     components.RepositoryList
	userInfo     components.UserInfo
	viewport     viewport.Model
	ready        bool
	width        int
	height       int
	loading      bool
	error        error
	pinnedLoaded bool
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
	tabs := components.NewTabs([]string{"Info", "Pinned"})
	repoList := components.NewRepositoryList(nil)
	userInfo := components.NewUserInfo(user)

	return Model{
		user:         user,
		tabs:         tabs,
		repoList:     repoList,
		userInfo:     userInfo,
		ready:        false,
		loading:      false,
		error:        nil,
		pinnedLoaded: false,
	}
}

// Init initializes the Model
func (m Model) Init() tea.Cmd {
	return nil
}

// fetchRepositories fetches pinned repositories
func fetchRepositories(username string) tea.Cmd {
	return func() tea.Msg {
		repos, err := github.FetchPinnedRepositories(context.Background(), username)
		return fetchRepositoriesMsg{repositories: repos, err: err}
	}
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
		case "right", "l":
			m.tabs.Next()
			return m, func() tea.Msg {
				return tabSelectedMsg{index: m.tabs.Current}
			}
		case "left", "h":
			m.tabs.Prev()
			return m, func() tea.Msg {
				return tabSelectedMsg{index: m.tabs.Current}
			}
		case "r":
			if m.tabs.Current == 1 && m.error != nil {
				m.loading = true
				m.error = nil
				cmd = fetchRepositories(m.user.Login)
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-4) // 4 for tabs and help
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 4
		}
		m.repoList.SetSize(msg.Width, msg.Height-4)

	case components.RepositorySelectedMsg:
		if msg.Repository != nil {
			cmd := openURL(msg.Repository.URL)
			cmds = append(cmds, cmd)
		}

	case fetchRepositoriesMsg:
		m.loading = false
		if msg.err != nil {
			m.error = msg.err
			return m, nil
		}
		m.repositories = msg.repositories
		m.repoList = components.NewRepositoryList(msg.repositories)
		m.repoList.SetSize(m.width, m.height-4)
		m.pinnedLoaded = true

	case tabSelectedMsg:
		if msg.index == 1 && !m.pinnedLoaded && !m.loading {
			m.loading = true
			m.error = nil
			cmd = fetchRepositories(m.user.Login)
			cmds = append(cmds, cmd)
		}
	}

	if m.tabs.Current == 1 { // Pinned tab
		newRepoList, cmd := m.repoList.Update(msg)
		m.repoList = *newRepoList
		cmds = append(cmds, cmd)
	} else {
		m.viewport.SetContent(m.userInfo.View())
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var content string

	// Tabs
	content += m.tabs.View() + "\n\n"

	// Content
	if m.tabs.Current == 1 { // Pinned tab
		if m.loading {
			content += "Loading..."
		} else if m.error != nil {
			errMsg := m.error.Error()
			if strings.Contains(errMsg, "Could not resolve") {
				content += errorStyle.Render("Authentication error") + "\n"
				content += errorHelpStyle.Render("Please run 'gh auth login' to authenticate with GitHub")
			} else if strings.Contains(errMsg, "connect:") || strings.Contains(errMsg, "timeout") {
				content += errorStyle.Render("Network error") + "\n"
				content += errorHelpStyle.Render("Please check your internet connection")
			} else {
				content += errorStyle.Render("Error: "+errMsg) + "\n"
				content += errorHelpStyle.Render("An unexpected error occurred")
			}
			content += "\n\n" + errorHelpStyle.Render("Press r to retry")
		} else {
			content += m.repoList.View()
		}
	} else {
		content += m.viewport.View()
	}

	// Help
	var help string
	if m.tabs.Current == 1 {
		if m.error != nil {
			help = "r: Retry • ←/→: Switch tabs • q: Quit"
		} else {
			help = "↑/↓: Navigate • enter: Open in browser • q: Quit"
		}
	} else {
		help = "←/→: Switch tabs • q: Quit"
	}
	content += "\n" + dividerStyle.Render(help)

	return content
}

// openURL opens the given URL in the default browser
func openURL(url string) tea.Cmd {
	return func() tea.Msg {
		b := browser.New("", os.Stdout, os.Stdin)
		_ = b.Browse(url)
		return nil
	}
}
