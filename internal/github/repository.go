package github

import (
	"context"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/shurcooL-graphql"
)

// Repository represents a GitHub repository
type Repository struct {
	Name        string
	Description string
	URL         string
	StarCount   int
	Language    string
}

// FetchPinnedRepositories fetches a user's pinned repositories
func FetchPinnedRepositories(ctx context.Context, login string) ([]Repository, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, err
	}

	var query struct {
		User struct {
			PinnedItems struct {
				Nodes []struct {
					Repository struct {
						Name        graphql.String
						Description graphql.String
						URL         graphql.String
						StargazerCount graphql.Int
						PrimaryLanguage struct {
							Name graphql.String
						}
					} `graphql:"... on Repository"`
				}
			} `graphql:"pinnedItems(first: 6, types: REPOSITORY)"`
		} `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": graphql.String(login),
	}

	err = client.Query("FetchPinnedRepositories", &query, variables)
	if err != nil {
		return nil, err
	}

	repos := make([]Repository, 0, len(query.User.PinnedItems.Nodes))
	for _, node := range query.User.PinnedItems.Nodes {
		repos = append(repos, Repository{
			Name:        string(node.Repository.Name),
			Description: string(node.Repository.Description),
			URL:         string(node.Repository.URL),
			StarCount:   int(node.Repository.StargazerCount),
			Language:    string(node.Repository.PrimaryLanguage.Name),
		})
	}

	return repos, nil
}