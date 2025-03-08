package components

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	activeTabStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	tabGap = inactiveTabStyle.Render("  ")
)

// Tab represents a single tab
type Tab struct {
	Title    string
	Selected bool
}

// Tabs represents a collection of tabs
type Tabs struct {
	Tabs    []Tab
	Current int
}

// NewTabs creates a new Tabs instance
func NewTabs(titles []string) Tabs {
	tabs := make([]Tab, len(titles))
	for i, title := range titles {
		tabs[i] = Tab{
			Title:    title,
			Selected: i == 0,
		}
	}
	return Tabs{
		Tabs:    tabs,
		Current: 0,
	}
}

// Next selects the next tab
func (t *Tabs) Next() {
	t.Tabs[t.Current].Selected = false
	t.Current = (t.Current + 1) % len(t.Tabs)
	t.Tabs[t.Current].Selected = true
}

// Prev selects the previous tab
func (t *Tabs) Prev() {
	t.Tabs[t.Current].Selected = false
	t.Current = (t.Current - 1 + len(t.Tabs)) % len(t.Tabs)
	t.Tabs[t.Current].Selected = true
}

// View renders the tabs
func (t Tabs) View() string {
	var renderedTabs []string

	for i, tab := range t.Tabs {
		if tab.Selected {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(tab.Title))
		} else {
			renderedTabs = append(renderedTabs, inactiveTabStyle.Render(tab.Title))
		}
		if i < len(t.Tabs)-1 {
			renderedTabs = append(renderedTabs, tabGap)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}