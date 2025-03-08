package components

import (
	"fmt"

	"github.com/tnagatomi/gh-portrait/internal/github"
)

// RepositoryItem represents a repository in the list
type RepositoryItem struct {
	repository github.Repository
	listType   string
}

// Title returns the repository name and language
func (r RepositoryItem) Title() string {
	var name string
	if r.listType == "contributed" {
		name = fmt.Sprintf("%s/%s", r.repository.Owner, r.repository.Name)
	} else {
		name = r.repository.Name
	}

	if r.repository.Language != "" {
		return fmt.Sprintf("%s (%s)", name, r.repository.Language)
	}
	return name
}

// Description returns the repository description and star count
func (r RepositoryItem) Description() string {
	desc := r.repository.Description
	if desc == "" {
		desc = "No description"
	}
	return fmt.Sprintf("%s (%d stars)", desc, r.repository.StarCount)
}

// FilterValue returns the value to use for filtering
func (r RepositoryItem) FilterValue() string {
	return r.repository.Name
}