package components

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tnagatomi/gh-portrait/internal/github"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))
)

// RepositorySelectedMsg is sent when a repository is selected
type RepositorySelectedMsg struct {
	Repository *github.Repository
}

// RepositoryList represents a list of repositories
type RepositoryList struct {
	list     list.Model
	selected *github.Repository
	listType string
}

// NewRepositoryList creates a new RepositoryList
func NewRepositoryList(repositories []github.Repository, listType string) RepositoryList {
	items := make([]list.Item, len(repositories))
	for i, repo := range repositories {
		items[i] = RepositoryItem{repository: repo, listType: listType}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("86"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("243"))

	l := list.New(items, delegate, 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	switch listType {
	case "pinned":
		l.Title = "Pinned repositories"
	case "owning":
		l.Title = "Most starred repositories"
	case "contributed":
		l.Title = "Most starred contributed repositories (in the past year)"
	}

	l.Styles.FilterPrompt = lipgloss.NewStyle()
	l.Styles.FilterCursor = lipgloss.NewStyle()

	return RepositoryList{
		list:     l,
		selected: nil,
		listType: listType,
	}
}

// SetSize sets the size of the list
func (r *RepositoryList) SetSize(width, height int) {
	r.list.SetSize(width, height)
}

// Update handles list updates
func (r *RepositoryList) Update(msg tea.Msg) (*RepositoryList, tea.Cmd) {
	var cmd tea.Cmd
	r.list, cmd = r.list.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if i, ok := r.list.SelectedItem().(RepositoryItem); ok {
				r.selected = &i.repository
				return r, func() tea.Msg {
					return RepositorySelectedMsg{Repository: r.selected}
				}
			}
		}
	}

	return r, cmd
}

// View renders the list
func (r RepositoryList) View() string {
	return r.list.View()
}

// Selected returns the currently selected repository
func (r RepositoryList) Selected() *github.Repository {
	return r.selected
}
