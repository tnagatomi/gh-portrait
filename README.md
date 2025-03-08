# gh-portrait

`gh-portrait` is a Terminal User Interface (TUI) extension for GitHub CLI (gh) that provides an interactive way to view GitHub user's profile and repositories.

## Installation

```bash
gh extension install tnagatomi/gh-portrait
```

### Requirements

- GitHub CLI (gh)

## Usage

```bash
gh portrait [username]
```

## Features

### User Profile View

- Display user information
- Show user's README

<img width="469" alt="Info tab" src="https://github.com/user-attachments/assets/93c7df43-5c64-4c27-bb76-ae27821f8975" />

### Repository List

- Browse user repositories
  - Pinned, most starred repositories, most starred contributed repositories
- Open selected repository by browser

<img width="833" alt="Pinned tab" src="https://github.com/user-attachments/assets/31ab8237-8ac0-447b-9e38-cd350a472cab" />
<img width="1099" alt="Ownning tab" src="https://github.com/user-attachments/assets/df30cc88-d592-4c36-97fd-5414cbea9b91" />
<img width="862" alt="Contributed tab" src="https://github.com/user-attachments/assets/52128e3c-5b39-4bce-a90a-37b53494cbc8" />

### Navigation

- Left/Right arrows or h/l: Switch between tabs
- Up/Down arrows or k/j: Navigate repositories
- Enter: Open repositories
- q: Quit application
