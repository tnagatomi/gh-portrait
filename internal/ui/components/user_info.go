package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tnagatomi/gh-portrait/internal/github"
)

var (
	userInfoTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))
)

// UserInfo represents the user information view
type UserInfo struct {
	user           *github.User
	renderer       MarkdownRenderer
	cachedREADME   string
	viewWidth      int
	readmeRendered bool
}

// NewUserInfo creates a new UserInfo instance
func NewUserInfo(user *github.User, renderer MarkdownRenderer) UserInfo {
	return UserInfo{
		user:           user,
		renderer:       renderer,
		viewWidth:      80, // Default width
		readmeRendered: false,
	}
}

// SetWidth updates the view width and triggers README re-rendering if needed
func (u *UserInfo) SetWidth(width int) {
	if u.viewWidth != width {
		u.viewWidth = width
		u.readmeRendered = false // Force re-render on width change
	}
}

// renderREADME renders the README content with the current width
func (u *UserInfo) renderREADME() {
	if u.user.README == nil {
		u.cachedREADME = ""
		return
	}

	u.cachedREADME = u.renderer.Render(*u.user.README, u.viewWidth)
	u.readmeRendered = true
}

// View renders the user information
func (u *UserInfo) View() string {
	var content string

	// Info section
	content += userInfoTitleStyle.Render("Info") + "\n"
	content += "Name: " + u.user.Name + "\n"
	if u.user.Bio != "" {
		content += "Bio: " + u.user.Bio + "\n"
	}
	if u.user.Pronouns != "" {
		content += "Pronouns: " + u.user.Pronouns + "\n"
	}
	if u.user.Company != "" {
		content += "Company: " + u.user.Company + "\n"
	}
	if u.user.Location != "" {
		content += "Location: " + u.user.Location + "\n"
	}
	if u.user.WebsiteURL != "" {
		content += "Website: " + u.user.WebsiteURL + "\n"
	}
	content += "\n"

	// Social accounts section
	if len(u.user.Social) > 0 {
		content += userInfoTitleStyle.Render("Social accounts") + "\n"
		for _, account := range u.user.Social {
			content += fmt.Sprintf("%s: %s\n",
				account.Provider,
				account.URL,
			)
		}
		content += "\n"
	}

	// README section
	if u.user.README != nil {
		if !u.readmeRendered {
			u.renderREADME()
		}

		// Create a divider line using box-drawing characters
		divider := strings.Repeat("─", 50) + "\n\n"
		content += divider
		content += u.cachedREADME
	}

	return content
}
