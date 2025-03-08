# gh-portrait

`gh-portrait` is a Terminal User Interface (TUI) extension for GitHub CLI (gh) that provides an interactive way to view GitHub user's profile and repositories.

## Installation

```bash
gh extension install tnagatomi/gh-portrait
```

### Requirements

- GitHub CLI (gh)

## Features

### User Profile View

- Display user information
- Show user's README

### Repository List

- Browse user repositories
  - Pinned, most starred repositories, most starred contributed repositories
- Open selected repository by browser

## Usage

```bash
gh portrait [username]
```

### Navigation

- Left/Right arrows or h/l: Switch between tabs
- Up/Down arrows or k/j: Navigate repositories
- Enter: Open repositories
- q: Quit application
