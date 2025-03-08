package components

import (
	"testing"

	"github.com/tnagatomi/gh-portrait/internal/github"
)

func TestNewRepositoryList(t *testing.T) {
	tests := []struct {
		name      string
		listType  string
		repos     []github.Repository
		wantTitle string
	}{
		{
			name:     "pinned repositories",
			listType: "pinned",
			repos: []github.Repository{
				{Name: "gh-portrait"},
				{Name: "cli"},
			},
			wantTitle: "Pinned repositories",
		},
		{
			name:     "owning repositories",
			listType: "owning",
			repos: []github.Repository{
				{Name: "gh-portrait"},
				{Name: "cli"},
			},
			wantTitle: "Most starred repositories",
		},
		{
			name:     "contributed repositories",
			listType: "contributed",
			repos: []github.Repository{
				{Name: "gh-portrait"},
				{Name: "cli"},
			},
			wantTitle: "Most starred contributed repositories (in the past year)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewRepositoryList(tt.repos, tt.listType)

			// Check title
			if list.list.Title != tt.wantTitle {
				t.Errorf("NewRepositoryList() title = %v, want %v", list.list.Title, tt.wantTitle)
			}

			// Check number of items
			if len(list.list.Items()) != len(tt.repos) {
				t.Errorf("NewRepositoryList() items count = %v, want %v", len(list.list.Items()), len(tt.repos))
			}

			// Check initial state
			if list.selected != nil {
				t.Error("NewRepositoryList() selected should be nil initially")
			}

			// Check list type
			if list.listType != tt.listType {
				t.Errorf("NewRepositoryList() listType = %v, want %v", list.listType, tt.listType)
			}
		})
	}
}

func TestRepositoryListSetSize(t *testing.T) {
	repos := []github.Repository{{Name: "gh-portrait"}}
	list := NewRepositoryList(repos, "owning")

	width, height := 100, 50
	list.SetSize(width, height)

	// Unfortunately, we can't directly access the width and height of the list.Model
	// as they are private fields. The best we can do is verify that SetSize doesn't panic
	// and the method exists for coverage purposes.
}

func TestRepositoryListSelected(t *testing.T) {
	repos := []github.Repository{{Name: "gh-portrait"}}
	list := NewRepositoryList(repos, "owning")

	// Initially, no repository should be selected
	if got := list.Selected(); got != nil {
		t.Errorf("Selected() = %v, want nil", got)
	}
}
