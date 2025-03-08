package github

import (
	"context"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type User struct {
	Login      string
	Name       string
	Bio        string
	Pronouns   string
	Company    string
	Location   string
	WebsiteURL string
	Following  int
	Followers  int
	Social     []SocialAccount
	README     *string // Nullable README content
}

type SocialAccount struct {
	Provider string
	URL      string
}

func FetchUser(ctx context.Context, login string) (*User, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, err
	}

	var query struct {
		User struct {
			Login      graphql.String
			Name       graphql.String
			Bio        graphql.String
			Pronouns   graphql.String
			Company    graphql.String
			Location   graphql.String
			WebsiteUrl graphql.String
			Following  struct {
				TotalCount graphql.Int
			}
			Followers struct {
				TotalCount graphql.Int
			}
			SocialAccounts struct {
				Nodes []struct {
					Provider graphql.String
					URL      graphql.String
				}
			} `graphql:"socialAccounts(first: 10)"`
			Repository struct {
				Object *struct {
					Blob struct {
						Text graphql.String
					} `graphql:"... on Blob"`
				} `graphql:"object(expression: \"HEAD:README.md\")"`
			} `graphql:"repository(name: $login)"`
		} `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": graphql.String(login),
	}

	err = client.Query("FetchUser", &query, variables)
	if err != nil {
		if strings.Contains(err.Error(), "Could not resolve to a Repository") {
			// Ignore repository not found error
			err = nil
		} else {
			return nil, err
		}
	}

	// Convert social accounts
	social := make([]SocialAccount, 0, len(query.User.SocialAccounts.Nodes))
	for _, node := range query.User.SocialAccounts.Nodes {
		social = append(social, SocialAccount{
			Provider: string(node.Provider),
			URL:      string(node.URL),
		})
	}

	// Get README if it exists
	var readme *string
	if query.User.Repository.Object != nil {
		if text := string(query.User.Repository.Object.Blob.Text); text != "" {
			readme = &text
		}
	}

	return &User{
		Login:      string(query.User.Login),
		Name:       string(query.User.Name),
		Bio:        string(query.User.Bio),
		Pronouns:   string(query.User.Pronouns),
		Company:    string(query.User.Company),
		Location:   string(query.User.Location),
		WebsiteURL: string(query.User.WebsiteUrl),
		Following:  int(query.User.Following.TotalCount),
		Followers:  int(query.User.Followers.TotalCount),
		Social:     social,
		README:     readme,
	}, nil
}
