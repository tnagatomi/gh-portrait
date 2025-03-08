package components

import (
	"fmt"

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
	user *github.User
}

// NewUserInfo creates a new UserInfo instance
func NewUserInfo(user *github.User) UserInfo {
	return UserInfo{
		user: user,
	}
}

// View renders the user information
func (u UserInfo) View() string {
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
	}

	return content
}