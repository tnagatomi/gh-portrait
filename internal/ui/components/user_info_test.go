package components

import (
	"strings"
	"testing"

	"github.com/tnagatomi/gh-portrait/internal/github"
)

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}

func TestUserInfoView(t *testing.T) {
	tests := []struct {
		name     string
		user     *github.User
		width    int
		want     []string // Expected substrings in the output
		notWant  []string // Substrings that should not appear in the output
		setWidth bool     // Whether to test width changes
	}{
		{
			name: "user with required fields only",
			user: &github.User{
				Name: "Takayuki Nagatomi",
			},
			width: 80,
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
				"─", // No divider without README
			},
			setWidth: false,
		},
		{
			name: "user with README",
			user: &github.User{
				Name:   "Takayuki Nagatomi",
				README: stringPtr("# Test README\nThis is a test README file."),
			},
			width: 80,
			want: []string{
				"Info",
				"Name: Takayuki Nagatomi",
				"─", // Divider should be present
				"Test README",
				"This is a test README file.",
			},
			notWant: []string{
				"Error rendering README:",
			},
			setWidth: false,
		},
		{
			name: "user with README and width change",
			user: &github.User{
				Name:   "Takayuki Nagatomi",
				README: stringPtr("# Test README\nThis is a test README file."),
			},
			width: 40,
			want: []string{
				"Info",
				"Name: Takayuki Nagatomi",
				"─",
				"Test README",
				"This is a test README file.", // Width change should affect wrapping
			},
			notWant: []string{
				"Error rendering README:",
			},
			setWidth: true,
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
			width: 80,
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
				"─", // No divider without README
			},
			setWidth: false,
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
			testRenderer := NewTestRenderer()
			ui := NewUserInfo(tt.user, testRenderer)
			if tt.setWidth {
				ui.SetWidth(tt.width)
			}
			got := ui.View()

			// Remove ANSI escape sequences for comparison
			got = strings.ReplaceAll(got, "\x1b[0m", "")
			got = strings.ReplaceAll(got, "\x1b[1m", "")

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

			// Test width change if required
			if tt.setWidth {
				// Get view again to verify caching
				secondView := ui.View()
				if secondView != got {
					t.Error("UserInfo.View() content changed without width change")
				}

				// Change width and verify content changes
				ui.SetWidth(tt.width + 10)
				thirdView := ui.View()
				if thirdView == got {
					t.Error("UserInfo.View() content did not change after width change")
				}
			}
		})
	}
}
