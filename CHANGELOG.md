# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2025-12-03

### Added
- TOML configuration file support
- Auto-creation of config directory and default config
- Config stored in `~/.config/taskmanager/config.toml`
- Configurable task directory path
- Config validation and error handling
- Helpful error messages when config issues occur

### Changed
- Task directory is now configurable via TOML config
- UI displays the configured directory path
- Default directory remains `~/.tasks`

### Phase 3 Complete ✅
- Full TOML configuration management
- User-configurable task directory
- Auto-creation of default config on first run

## [0.2.0] - 2025-12-03

### Added
- File system integration to read actual markdown files
- Load tasks from `~/.tasks` directory
- Display modification dates for each task file
- Automatic sorting by modification time (newest first)
- Error handling for directory access issues
- Empty state message when no markdown files found
- Task count display in footer

### Changed
- Replaced static task list with dynamic file loading
- Updated model to use `taskFile` struct instead of strings
- Enhanced UI with file metadata display

### Phase 2 Complete ✅
- Real markdown file loading from configured directory
- File metadata display (modification dates)
- Graceful error handling

## [0.1.1] - 2025-12-03

### Added
- Alternate screen mode (full-screen TUI like lazygit)
- Mouse support for future interaction features

### Changed
- App now takes over full terminal and restores on exit

## [0.1.0] - 2025-12-03

### Added
- Initial project setup with Go module
- Basic Bubble Tea TUI implementation
- Static list display with keyboard navigation
- Project documentation (README.md, project-plan.md)
- Basic keyboard controls (up/down/vim keys, quit)

### Phase 1 Complete ✅
- Working TUI application with static list
- Understanding of Bubble Tea's Model-View-Update pattern
