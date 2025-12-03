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

### Phase 2: File System Integration
**Goal**: Read and display actual markdown files from a directory

**Tasks**:
- [ ] Add file system scanning functionality
- [ ] Filter for `.md` files
- [ ] Display actual filenames in the list
- [ ] Show last modification date for each file
- [ ] Handle errors gracefully (directory doesn't exist, no permissions, etc.)

**Deliverables**:
- TUI displays real markdown files from a hardcoded directory
- Files sorted by modification date

### Phase 3: Configuration Management
**Goal**: Make the application configurable via TOML config file

**Tasks**:
- [ ] Add TOML parsing library
- [ ] Create config structure
- [ ] Load config from `~/.config/taskmanager/config.toml`
- [ ] Create default config if it doesn't exist
- [ ] Support configurable directory path
- [ ] Add config validation

**Deliverables**:
- Configuration file support
- Auto-creation of config directory and default config
- Documentation on config options

### Phase 4: Multi-Directory Support
**Goal**: Support multiple task directories

**Tasks**:
- [ ] Update config to support multiple directories
- [ ] Aggregate files from all configured directories
- [ ] Display directory source for each task
- [ ] Handle duplicate filenames across directories

**Deliverables**:
- Support for multiple task directories in config
- Clear indication of which directory each task belongs to

### Phase 5: Markdown Frontmatter Parsing
**Goal**: Extract and display task metadata from markdown frontmatter

**Tasks**:
- [ ] Add markdown/frontmatter parsing library
- [ ] Define frontmatter schema (status, priority, due date, etc.)
- [ ] Parse frontmatter from each markdown file
- [ ] Display parsed metadata in the list
- [ ] Add sorting/filtering by metadata fields

**Deliverables**:
- Rich task information from frontmatter
- Sortable/filterable task list

### Phase 6: Task Interaction
**Goal**: Enable viewing and basic task operations

**Tasks**:
- [ ] Implement task selection
- [ ] View task content (full markdown)
- [ ] Edit task (open in $EDITOR)
- [ ] Create new task
- [ ] Delete task (with confirmation)

**Deliverables**:
- Full CRUD operations on tasks
- Integration with system editor

### Phase 7: Advanced Features
**Goal**: Polish and additional features

**Tasks**:
- [ ] Search/filter functionality
- [ ] Task completion tracking
- [ ] Keyboard shortcuts reference
- [ ] Color theming
- [ ] Performance optimization for large task lists

**Deliverables**:
- Production-ready task manager
- Comprehensive documentation

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
