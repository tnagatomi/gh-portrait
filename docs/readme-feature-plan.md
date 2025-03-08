# README Feature Implementation Plan

## 1. Data Fetching Implementation

```mermaid
sequenceDiagram
    participant App
    participant GitHub API
    participant User Component
    participant Glamour

    App->>GitHub API: FetchUser(login)
    GitHub API-->>App: User Info
    App->>GitHub API: FetchUserReadme(login)
    GitHub API-->>App: README Content (or nil)
    opt README exists
        App->>Glamour: Render README
        Glamour-->>App: Rendered Content
        App->>User Component: Display Info, Divider & README
    else
        App->>User Component: Display Info only
    end
```

### 1.1 GitHub Package Extension
- Add README fetching GraphQL query to `internal/github/user.go`
- Add README field to User struct (as nullable)

### 1.2 Dependency Addition
- Execute `go get github.com/charmbracelet/glamour`

## 2. UI Implementation

### 2.1 UserInfo Component Extension
- When README exists:
  - Display divider
  - Display README section
  - Render Markdown using glamour
- When README does not exist:
  - Display user information only (current behavior)

### 2.2 Layout Adjustments
- Adjust viewport height (for README display)
- Verify scrolling functionality

## 3. Error Handling
- Handle Markdown rendering errors
- Handle GraphQL query errors