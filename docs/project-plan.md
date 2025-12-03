# Task Manager - Project Plan

## Overview
A TUI (Terminal User Interface) application built with Go and Bubble Tea for managing tasks stored as markdown files across multiple directories.

## Goals
- Manage tasks stored as individual markdown (.md) files
- Support multiple task directories (e.g., tasks within project folders)
- Use markdown frontmatter for task metadata
- Provide an intuitive terminal-based interface
- Store configuration in `~/.config/taskmanager/config.toml`

## Technology Stack
- **Language**: Go
- **TUI Framework**: Bubble Tea
- **Config Format**: TOML (fallback: YAML if needed)
- **Task Storage**: Markdown files with frontmatter

## Implementation Phases

### Phase 1: Basic Setup ✅
**Goal**: Set up Go module and create a basic Bubble Tea TUI with a static list

**Tasks**:
- [x] Initialize Go module
- [x] Install Bubble Tea dependency
- [x] Create main.go with basic TUI structure
- [x] Display a static list of items
- [x] Understand Bubble Tea's Model-View-Update pattern

**Deliverables**:
- Working `taskmanager` command that displays a static list
- Basic understanding of Bubble Tea architecture

### Phase 2: File System Integration ✅
**Goal**: Read and display actual markdown files from a directory

**Tasks**:
- [x] Add file system scanning functionality
- [x] Filter for `.md` files
- [x] Display actual filenames in the list
- [x] Show last modification date for each file
- [x] Handle errors gracefully (directory doesn't exist, no permissions, etc.)

**Deliverables**:
- TUI displays real markdown files from `~/.tasks` directory
- Files sorted by modification date (newest first)

### Phase 3: Configuration Management ✅
**Goal**: Make the application configurable via TOML config file

**Tasks**:
- [x] Add TOML parsing library
- [x] Create config structure
- [x] Load config from `~/.config/taskmanager/config.toml`
- [x] Create default config if it doesn't exist
- [x] Support configurable directory path
- [x] Add config validation

**Deliverables**:
- Configuration file support
- Auto-creation of config directory and default config
- Documentation on config options

### Phase 4: Multi-Directory Support ✅
**Goal**: Support multiple task directories

**Tasks**:
- [x] Update config to support multiple directories
- [x] Aggregate files from all configured directories
- [x] Display directory source for each task
- [x] Handle duplicate filenames across directories
- [x] Maintain backward compatibility with single directory config

**Deliverables**:
- Support for multiple task directories in config
- Clear indication of which directory each task belongs to
- Backward compatible with existing single-directory configs

### Phase 5: Markdown Frontmatter Parsing ✅
**Goal**: Extract and display task metadata from markdown frontmatter

**Tasks**:
- [x] Add markdown/frontmatter parsing library
- [x] Define frontmatter schema (status, priority, due date, etc.)
- [x] Parse frontmatter from each markdown file
- [x] Display parsed metadata in the list
- [x] Add status and priority emojis

**Deliverables**:
- Rich task information from frontmatter
- Visual indicators for status and priority
- Backwards compatible with non-frontmatter files

### Phase 6: Task Interaction ✅
**Goal**: Enable viewing and basic task operations

**Tasks**:
- [x] Implement task selection
- [x] View task content (full markdown)
- [x] Edit task (open in $EDITOR)
- [x] Create new task
- [x] Delete task (with confirmation)

**Deliverables**:
- Full CRUD operations on tasks
- Integration with system editor
- Delete confirmation dialog
- Automatic task list reload after create/edit/delete
- Smooth transitions between list and view modes

**Notes**:
- Tasks use $EDITOR environment variable (fallback to vim)
- New tasks created with timestamp-based filenames
- Template includes frontmatter with common fields
- Delete requires 'y' confirmation to prevent accidents
- All operations return to list view instead of quitting

### Phase 7: Advanced Features
**Goal**: Polish and additional features

**Tasks**:
- [x] Search/filter functionality
- [ ] Task completion tracking
- [x] Keyboard shortcuts reference
- [x] Color theming
- [ ] Performance optimization for large task lists

**Deliverables**:
- Production-ready task manager
- Comprehensive documentation

**Completed Features**:
- Real-time search/filter with `/` key
- Search across filename, title, status, and tags
- Case-insensitive matching with live results
- Help screen (`?` or `h`) with all keyboard shortcuts
- Organized by context (List View, Search Mode, Task View, etc.)
- Color theming with Lip Gloss library
  - Color-coded status (gray/orange/green)
  - Color-coded priority (red/orange/gray)
  - Professional visual hierarchy
  - Improved scannability and aesthetics

## Configuration File Structure

### Initial (Phase 3)
```toml
[taskmanager]
directory = "~/tasks"
```

### Future (Phase 4+)
```toml
[taskmanager]
directories = [
    "~/tasks",
    "~/Projects/project-a/tasks",
    "~/Projects/project-b/tasks"
]

[display]
sort_by = "modified"  # modified, name, status
show_path = true
date_format = "2006-01-02"

[editor]
command = "$EDITOR"  # Uses $EDITOR env var by default
```

## Frontmatter Schema (Future)

```yaml
---
title: "Task Title"
status: "todo"  # todo, in-progress, done
priority: "high"  # low, medium, high
due_date: "2025-12-31"
tags: ["urgent", "feature"]
created: "2025-12-03"
---
```

## Learning Resources
- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea)
- [Go File I/O](https://gobyexample.com/reading-files)
- [TOML in Go](https://github.com/BurntSushi/toml)
