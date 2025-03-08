package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tnagatomi/gh-portrait/internal/github"
	"github.com/tnagatomi/gh-portrait/internal/ui"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: gh portrait <username>")
		os.Exit(1)
	}

	username := os.Args[1]
	user, err := github.FetchUser(context.Background(), username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := ui.Start(user); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
