# Task Manager

A beautiful terminal-based task manager built with Go and Bubble Tea. Manage your tasks as markdown files across multiple project directories.

## Status

🚧 **Early Development** - Currently in Phase 4 Complete

## Features

### Current (Phase 4)
- ✅ Basic TUI interface using Bubble Tea
- ✅ Full-screen alternate mode (like lazygit)
- ✅ TOML configuration file support (`~/.config/taskmanager/config.toml`)
- ✅ **Multi-directory support** - track tasks across multiple project folders
- ✅ Configurable task directories (single or multiple)
- ✅ Load and display markdown files from all configured directories
- ✅ Show source directory for each task (when using multiple directories)
- ✅ Show last modification date for each task
- ✅ Automatic sorting by modification time (newest first)
- ✅ Keyboard navigation (↑/↓ or k/j)
- ✅ Backward compatible with single directory config

### Planned
- 🎨 Markdown frontmatter support for task metadata
- 🔍 Filter and search tasks
- 📝 View and edit tasks

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
