# Task Manager

A beautiful terminal-based task manager built with Go and Bubble Tea. Manage your tasks as markdown files across multiple project directories.

## Status

🚧 **Early Development** - Currently in Phase 1 (Basic Setup)

## Features

### Current (Phase 1)
- ✅ Basic TUI interface using Bubble Tea
- ✅ Static list display

### Planned
- 📋 Load and display markdown files from configured directories
- 🔍 Filter and search tasks
- 📝 View and edit tasks
- 🎨 Markdown frontmatter support for task metadata
- 📂 Multi-directory support
- ⚙️ TOML-based configuration

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

### Keyboard Controls
- `↑/k` - Move up
- `↓/j` - Move down
- `q` - Quit

## Configuration

Configuration will be stored in `~/.config/taskmanager/config.toml` (Phase 3).

Example configuration (coming soon):
```toml
[taskmanager]
directory = "~/tasks"
```

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
