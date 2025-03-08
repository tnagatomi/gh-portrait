package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tnagatomi/gh-portrait/internal/github"
	"github.com/tnagatomi/gh-portrait/internal/ui"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: gh portrait <username>")
		os.Exit(1)
	}

	username := os.Args[1]
	ctx := context.Background()

	// Fetch user information
	user, err := github.FetchUser(ctx, username)
	if err != nil {
		if strings.Contains(err.Error(), "Could not resolve to a User") {
			fmt.Fprintf(os.Stderr, "Error: User '%s' not found\n", username)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}

	if err := ui.Start(user); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
