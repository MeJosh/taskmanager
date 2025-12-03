# Task Manager

A beautiful terminal-based task manager built with Go and Bubble Tea. Manage your tasks as markdown files across multiple project directories.

## Status

🚧 **Early Development** - Currently in Phase 5 Complete

## Features

### Current (Phase 5)
- ✅ Basic TUI interface using Bubble Tea
- ✅ Full-screen alternate mode (like lazygit)
- ✅ TOML configuration file support (`~/.config/taskmanager/config.toml`)
- ✅ **Multi-directory support** - track tasks across multiple project folders
- ✅ **Markdown frontmatter parsing** - extract rich task metadata
- ✅ Status indicators (`[ ]` todo, `[~]` in-progress, `[✓]` done)
- ✅ Priority indicators (high, med, low)
- ✅ Display task titles from frontmatter
- ✅ Show source directory for each task (when using multiple directories)
- ✅ Show last modification date for each task
- ✅ Automatic sorting by modification time (newest first)
- ✅ Keyboard navigation (↑/↓ or k/j)
- ✅ Backward compatible with files without frontmatter

### Planned
- 🔍 Filter and search tasks
- 📝 View and edit tasks
- 📊 Sort by different criteria (status, priority, due date)

## Installation

### Prerequisites
- Go 1.21 or higher

### Building from Source

```bash
# Clone the repository
git clone https://github.com/MeJosh/taskmanager.git
cd taskmanager

# Build the application
go build -o taskmanager

# (Optional) Install to your PATH
go install
```

## Usage

```bash
# Run the task manager
./taskmanager
```

The application will display all `.md` files from your `~/.tasks` directory, sorted by modification date (newest first).

### Keyboard Controls
- `↑/k` - Move up
- `↓/j` - Move down
- `q` - Quit

## Configuration

Configuration is stored in the system's standard config directory:
- **macOS/Linux**: `~/.config/taskmanager/config.toml`
- **Windows**: `%AppData%\taskmanager\config.toml`

On first run, a default configuration file will be created automatically.

### Single Directory (Backward Compatible)

```toml
[taskmanager]
directory = "~/.tasks"
```

### Multiple Directories (Recommended)

Track tasks across multiple project directories:

```toml
[taskmanager]
directories = [
    "~/.tasks",
    "~/Projects/project-a/tasks",
    "~/Projects/project-b/tasks"
]
```

When using multiple directories, the app will:
- Load all `.md` files from all configured directories
- Sort them by modification time (newest first)
- Display the source directory for each task

## Task Files with Frontmatter

You can add YAML frontmatter to your markdown task files to include rich metadata:

```markdown
---
title: "Implement user authentication"
status: "in-progress"
priority: "high"
tags: ["security", "backend"]
due_date: 2025-12-31T00:00:00Z
created: 2025-12-01T10:00:00Z
---

# Task content goes here

Your markdown content...
```

### Supported Frontmatter Fields

- **title**: Display name for the task (shown instead of filename)
- **status**: Task status - `todo`, `in-progress`, or `done`
  - `todo` = `[ ]`, `in-progress` = `[~]`, `done` = `[✓]`
- **priority**: Task priority - `low`, `medium`, or `high`
  - Displays as: `low`, `med`, `high`
- **tags**: Array of tags for categorization
- **due_date**: When the task is due (ISO 8601 format)
- **created**: When the task was created (ISO 8601 format)

Tasks without frontmatter work perfectly fine - the app is fully backwards compatible.

## Project Structure

```
taskmanager/
├── docs/              # Project documentation
│   └── project-plan.md
├── main.go            # Application entry point
├── go.mod             # Go module definition
├── README.md          # This file
└── CHANGELOG.md       # Version history
```

## Development

See [docs/project-plan.md](docs/project-plan.md) for the detailed development roadmap.

## License

MIT License - See LICENSE file for details

## Contributing

This is a personal learning project, but suggestions and feedback are welcome!

## Acknowledgments

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling (future)
