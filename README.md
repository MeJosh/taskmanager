# Task Manager

A beautiful terminal-based task manager built with Go and Bubble Tea. Manage your tasks as markdown files across multiple project directories.

## Status

🚧 **Early Development** - Currently in Phase 3 Complete

## Features

### Current (Phase 3)
- ✅ Basic TUI interface using Bubble Tea
- ✅ Full-screen alternate mode (like lazygit)
- ✅ TOML configuration file support (`~/.config/taskmanager/config.toml`)
- ✅ Configurable task directory
- ✅ Load and display markdown files from configured directory
- ✅ Show last modification date for each task
- ✅ Automatic sorting by modification time (newest first)
- ✅ Keyboard navigation (↑/↓ or k/j)

### Planned
- � Multi-directory support
- 🔍 Filter and search tasks
- 📝 View and edit tasks
- 🎨 Markdown frontmatter support for task metadata

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

Configuration is stored in `~/.config/taskmanager/config.toml`.

On first run, a default configuration file will be created automatically:

```toml
[taskmanager]
directory = "~/.tasks"
```

You can edit this file to change where your tasks are stored:

```toml
[taskmanager]
directory = "~/Documents/my-tasks"
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
