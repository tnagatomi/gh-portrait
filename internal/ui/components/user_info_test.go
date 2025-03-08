package components

import (
	"strings"
	"testing"

	"github.com/tnagatomi/gh-portrait/internal/github"
)

func TestUserInfoView(t *testing.T) {
	tests := []struct {
		name string
		user *github.User
		want []string // Expected substrings in the output
		notWant []string // Substrings that should not appear in the output
	}{
		{
			name: "user with required fields only",
			user: &github.User{
				Name: "Takayuki Nagatomi",
			},
			want: []string{
				"Info",
				"Name: Takayuki Nagatomi",
			},
			notWant: []string{
				"Bio:",
				"Pronouns:",
				"Company:",
				"Location:",
				"Website:",
				"Social accounts",
			},
		},
		{
			name: "user with all fields",
			user: &github.User{
				Name:       "Takayuki Nagatomi",
				Bio:        "Software Engineer",
				Pronouns:   "they/them",
				Company:    "Example Inc.",
				Location:   "San Francisco",
				WebsiteURL: "https://example.com",
			},
			want: []string{
				"Info",
				"Name: Takayuki Nagatomi",
				"Bio: Software Engineer",
				"Pronouns: they/them",
				"Company: Example Inc.",
				"Location: San Francisco",
				"Website: https://example.com",
			},
			notWant: []string{
				"Social accounts",
			},
		},
		{
			name: "user with social accounts",
			user: &github.User{
				Name: "Takayuki Nagatomi",
				Social: []github.SocialAccount{
					{
						Provider: "BLUESKY",
						URL:      "https://bsky.app/profile/tnagatomi.okweird.net",
					},
					{
						Provider: "MASTODON",
						URL:      "https://hachyderm.io/@tnagatomi",
					},
				},
			},
			want: []string{
				"Info",
				"Name: Takayuki Nagatomi",
				"Social accounts",
				"BLUESKY: https://bsky.app/profile/tnagatomi.okweird.net",
				"MASTODON: https://hachyderm.io/@tnagatomi",
			},
			notWant: []string{
				"Bio:",
				"Pronouns:",
				"Company:",
				"Location:",
				"Website:",
			},
		},
		{
			name: "user with empty fields",
			user: &github.User{
				Name:       "Takayuki Nagatomi",
				Bio:        "",
				Pronouns:   "",
				Company:    "",
				Location:   "",
				WebsiteURL: "",
			},
			want: []string{
				"Info",
				"Name: Takayuki Nagatomi",
			},
			notWant: []string{
				"Bio:",
				"Pronouns:",
				"Company:",
				"Location:",
				"Website:",
				"Social accounts",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := NewUserInfo(tt.user)
			got := ui.View()

			// Check for expected substrings
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("UserInfo.View() = %v, want substring %v", got, want)
				}
			}

			// Check for unwanted substrings
			for _, notWant := range tt.notWant {
				if strings.Contains(got, notWant) {
					t.Errorf("UserInfo.View() = %v, should not contain substring %v", got, notWant)
				}
			}

			// Verify styling
			if !strings.Contains(got, userInfoTitleStyle.Render("Info")) {
				t.Error("UserInfo.View() does not contain styled Info title")
			}

			if len(tt.user.Social) > 0 && !strings.Contains(got, userInfoTitleStyle.Render("Social accounts")) {
				t.Error("UserInfo.View() does not contain styled Social accounts title")
			}
		})
	}
}