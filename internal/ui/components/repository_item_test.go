package components

import (
	"testing"

	"github.com/tnagatomi/gh-portrait/internal/github"
)

func TestRepositoryItemTitle(t *testing.T) {
	tests := []struct {
		name     string
		item     RepositoryItem
		expected string
	}{
		{
			name: "repository with language (owning)",
			item: RepositoryItem{
				repository: github.Repository{
					Name:     "gh-portrait",
					Language: "Go",
				},
				listType: "owning",
			},
			expected: "gh-portrait (Go)",
		},
		{
			name: "repository without language (owning)",
			item: RepositoryItem{
				repository: github.Repository{
					Name: "gh-portrait",
				},
				listType: "owning",
			},
			expected: "gh-portrait",
		},
		{
			name: "repository with language (contributed)",
			item: RepositoryItem{
				repository: github.Repository{
					Owner:    "cli",
					Name:     "cli",
					Language: "Go",
				},
				listType: "contributed",
			},
			expected: "cli/cli (Go)",
		},
		{
			name: "repository without language (contributed)",
			item: RepositoryItem{
				repository: github.Repository{
					Owner: "cli",
					Name:  "cli",
				},
				listType: "contributed",
			},
			expected: "cli/cli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.Title()
			if got != tt.expected {
				t.Errorf("Title() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRepositoryItemDescription(t *testing.T) {
	tests := []struct {
		name     string
		item     RepositoryItem
		expected string
	}{
		{
			name: "repository with description",
			item: RepositoryItem{
				repository: github.Repository{
					Description: "GitHub profile visualization tool",
					StarCount:   42,
				},
			},
			expected: "GitHub profile visualization tool (42 stars)",
		},
		{
			name: "repository without description",
			item: RepositoryItem{
				repository: github.Repository{
					StarCount: 10,
				},
			},
			expected: "No description (10 stars)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.Description()
			if got != tt.expected {
				t.Errorf("Description() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRepositoryItemFilterValue(t *testing.T) {
	item := RepositoryItem{
		repository: github.Repository{
			Name: "gh-portrait",
		},
	}

	got := item.FilterValue()
	expected := "gh-portrait"
	if got != expected {
		t.Errorf("FilterValue() = %v, want %v", got, expected)
	}
}