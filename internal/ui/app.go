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
	tabIndex    int
}

// tabSelectedMsg is sent when a tab is selected
type tabSelectedMsg struct {
	index int
}

// Model represents the main application UI model
type Model struct {
	user            *github.User
	pinnedRepos     []github.Repository
	owningRepos     []github.Repository
	tabs            components.Tabs
	repoList        components.RepositoryList
	userInfo        components.UserInfo
	viewport        viewport.Model
	ready           bool
	width           int
	height          int
	loading         bool
	error           error
	pinnedLoaded    bool
	owningLoaded    bool
	currentTabIndex int
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
	tabs := components.NewTabs([]string{"Info", "Pinned", "Owning"})
	repoList := components.NewRepositoryList(nil, "pinned")
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
		owningLoaded: false,
	}
}

// Init initializes the Model
func (m Model) Init() tea.Cmd {
	return nil
}

// fetchRepositories fetches repositories based on the tab index
func fetchRepositories(username string, tabIndex int) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		var (
			repos []github.Repository
			err   error
		)

		switch tabIndex {
		case 1: // Pinned
			repos, err = github.FetchPinnedRepositories(ctx, username)
		case 2: // Owning
			repos, err = github.FetchOwningRepositories(ctx, username)
		}

		return fetchRepositoriesMsg{
			repositories: repos,
			err:         err,
			tabIndex:    tabIndex,
		}
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
			if m.error != nil {
				m.loading = true
				m.error = nil
				cmd = fetchRepositories(m.user.Login, m.currentTabIndex)
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
		m.currentTabIndex = msg.tabIndex
		if msg.err != nil {
			m.error = msg.err
			return m, nil
		}

		switch msg.tabIndex {
		case 1: // Pinned
			m.pinnedRepos = msg.repositories
			m.pinnedLoaded = true
			m.repoList = components.NewRepositoryList(msg.repositories, "pinned")
		case 2: // Owning
			m.owningRepos = msg.repositories
			m.owningLoaded = true
			m.repoList = components.NewRepositoryList(msg.repositories, "owning")
		}
		m.repoList.SetSize(m.width, m.height-4)

	case tabSelectedMsg:
		m.currentTabIndex = msg.index
		switch msg.index {
		case 1: // Pinned
			if !m.pinnedLoaded && !m.loading {
				m.loading = true
				m.error = nil
				cmd = fetchRepositories(m.user.Login, msg.index)
				cmds = append(cmds, cmd)
			} else if m.pinnedLoaded {
				m.repoList = components.NewRepositoryList(m.pinnedRepos, "pinned")
				m.repoList.SetSize(m.width, m.height-4)
			}
		case 2: // Owning
			if !m.owningLoaded && !m.loading {
				m.loading = true
				m.error = nil
				cmd = fetchRepositories(m.user.Login, msg.index)
				cmds = append(cmds, cmd)
			} else if m.owningLoaded {
				m.repoList = components.NewRepositoryList(m.owningRepos, "owning")
				m.repoList.SetSize(m.width, m.height-4)
			}
		}
	}

	if m.tabs.Current > 0 { // Repository tabs
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
	if m.tabs.Current > 0 { // Repository tabs
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
	if m.tabs.Current > 0 {
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
