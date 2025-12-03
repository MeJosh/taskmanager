package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// reloadTasksMsg is sent when we need to reload the task list
type reloadTasksMsg struct{}

// taskFile represents a markdown file with its metadata
type taskFile struct {
	name      string       // filename
	modTime   time.Time    // last modification time
	fullPath  string       // absolute path to the file
	sourceDir string       // which directory this task came from
	metadata  TaskMetadata // parsed frontmatter metadata
}

// viewMode represents different states of the application
type viewMode int

const (
	listMode          viewMode = iota // Showing the list of tasks
	taskViewMode                      // Viewing a single task's content
	confirmDeleteMode                 // Confirming task deletion
	searchMode                        // Searching/filtering tasks
)

// model represents the application state
// In Bubble Tea, the model holds all the data your application needs
type model struct {
	tasks        []taskFile    // Our list of task files
	filteredTasks []taskFile   // Filtered list based on search
	cursor       int           // Which task our cursor is pointing at
	err          error         // Any error encountered while loading files
	configDirs   []string      // The configured task directories
	showDirInfo  bool          // Whether to show directory info for each task
	config       DisplayConfig // Display configuration
	mode         viewMode      // Current view mode
	taskContent  string        // Content of the task being viewed
	searchQuery  string        // Current search query
}

// visibleTasks returns the list of tasks that should be displayed
// (either filtered tasks if searching, or all tasks otherwise)
func (m model) visibleTasks() []taskFile {
	if m.mode == searchMode && len(m.filteredTasks) > 0 {
		return m.filteredTasks
	}
	if m.mode == searchMode && m.searchQuery != "" {
		return []taskFile{} // No matches
	}
	return m.tasks
}

// filterTasks filters the task list based on the search query
func (m *model) filterTasks() {
	if m.searchQuery == "" {
		m.filteredTasks = m.tasks
		return
	}

	query := strings.ToLower(m.searchQuery)
	m.filteredTasks = []taskFile{}

	for _, task := range m.tasks {
		// Search in filename
		if strings.Contains(strings.ToLower(task.name), query) {
			m.filteredTasks = append(m.filteredTasks, task)
			continue
		}
		// Search in title
		if strings.Contains(strings.ToLower(task.metadata.Title), query) {
			m.filteredTasks = append(m.filteredTasks, task)
			continue
		}
		// Search in status
		if strings.Contains(strings.ToLower(task.metadata.Status), query) {
			m.filteredTasks = append(m.filteredTasks, task)
			continue
		}
		// Search in tags
		for _, tag := range task.metadata.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				m.filteredTasks = append(m.filteredTasks, task)
				break
			}
		}
	}

	// Reset cursor if out of bounds
	if m.cursor >= len(m.filteredTasks) {
		m.cursor = 0
	}
}

// expandPath expands ~ to the user's home directory
func expandPath(path string) (string, error) {
	if len(path) >= 2 && path[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("couldn't get home directory: %w", err)
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	return path, nil
}

// loadTasksFromDirectory reads all .md files from the specified directory
func loadTasksFromDirectory(dir string) ([]taskFile, error) {
	// Expand the tilde (~) to the user's home directory
	expandedDir, err := expandPath(dir)
	if err != nil {
		return nil, err
	}

	// Read all files in the directory
	entries, err := os.ReadDir(expandedDir)
	if err != nil {
		return nil, fmt.Errorf("couldn't read directory %s: %w", dir, err)
	}

	// Collect all .md files
	var tasks []taskFile
	for _, entry := range entries {
		// Skip directories, only process files
		if entry.IsDir() {
			continue
		}

		// Only include .md files
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		// Get file info for modification time
		info, err := entry.Info()
		if err != nil {
			// Skip files we can't read, but don't fail entirely
			continue
		}

		fullPath := filepath.Join(expandedDir, entry.Name())

		// Parse frontmatter metadata
		metadata, _ := parseFrontmatter(fullPath)
		// We ignore errors here - files without frontmatter are valid

		tasks = append(tasks, taskFile{
			name:      entry.Name(),
			modTime:   info.ModTime(),
			fullPath:  fullPath,
			sourceDir: dir, // Store the original (unexpanded) directory
			metadata:  metadata,
		})
	}

	return tasks, nil
}

// loadTasksFromDirectories reads all .md files from multiple directories
func loadTasksFromDirectories(dirs []string) ([]taskFile, error) {
	var allTasks []taskFile
	var errors []string

	for _, dir := range dirs {
		tasks, err := loadTasksFromDirectory(dir)
		if err != nil {
			// Don't fail completely, just track the error
			errors = append(errors, fmt.Sprintf("%s: %v", dir, err))
			continue
		}
		allTasks = append(allTasks, tasks...)
	}

	// Sort all tasks by modification time (newest first)
	sort.Slice(allTasks, func(i, j int) bool {
		return allTasks[i].modTime.After(allTasks[j].modTime)
	})

	// If we had errors but still got some tasks, return tasks with a warning
	if len(errors) > 0 && len(allTasks) > 0 {
		// Just log the errors, don't fail
		fmt.Fprintf(os.Stderr, "Warning: Some directories couldn't be read:\n")
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
	}

	// If we had errors and no tasks, return the error
	if len(errors) > 0 && len(allTasks) == 0 {
		return nil, fmt.Errorf("couldn't read any directories: %s", errors[0])
	}

	return allTasks, nil
}

// initialModel creates the starting state of our application
func initialModel() model {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return model{
			tasks:       nil,
			cursor:      0,
			err:         fmt.Errorf("failed to load config: %w", err),
			configDirs:  []string{"~/.tasks"}, // fallback
			showDirInfo: false,
			config:      defaultConfig().Display,
			mode:        listMode,
		}
	}

	// Get all configured directories
	dirs := cfg.TaskManager.GetDirectories()

	// Load tasks from all configured directories
	tasks, loadErr := loadTasksFromDirectories(dirs)

	return model{
		tasks:       tasks,
		cursor:      0,
		err:         loadErr,
		configDirs:  dirs,
		showDirInfo: len(dirs) > 1, // Show directory info if multiple directories
		config:      cfg.Display,
		mode:        listMode,
	}
}

// getEditor returns the user's preferred editor
func getEditor() string {
	// Check EDITOR environment variable
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	// Check VISUAL environment variable
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	// Default to vim
	return "vim"
}

// editTask opens the current task in the user's editor
func (m model) editTask() tea.Cmd {
	editor := getEditor()
	taskPath := m.tasks[m.cursor].fullPath

	c := exec.Command(editor, taskPath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		// After editing, reload the task list to show updated content
		return reloadTasksMsg{}
	})
}

// createTask creates a new task file and opens it in the editor
func (m model) createTask() tea.Cmd {
	editor := getEditor()

	// Use the first configured directory for new tasks
	firstDir, err := expandPath(m.configDirs[0])
	if err != nil {
		return nil
	}

	// Generate a filename based on timestamp
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("task-%s.md", timestamp)
	taskPath := filepath.Join(firstDir, filename)

	// Create a template for the new task
	template := `---
title: "New Task"
status: todo
priority: medium
created: ` + time.Now().Format(time.RFC3339) + `
---

# New Task

Write your task description here...
`

	// Write the template to the file
	if err := os.WriteFile(taskPath, []byte(template), 0644); err != nil {
		return nil
	}

	// Open in editor
	c := exec.Command(editor, taskPath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		// Return a message to reload tasks and go back to list mode
		return reloadTasksMsg{}
	})
}

// deleteTask deletes the current task file after confirmation
func (m model) deleteTask() tea.Model {
	taskPath := m.tasks[m.cursor].fullPath

	// Delete the file
	if err := os.Remove(taskPath); err != nil {
		m.err = fmt.Errorf("failed to delete task: %w", err)
		m.mode = listMode
		return m
	}

	// Remove the task from the list
	m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)

	// Adjust cursor if needed
	if m.cursor >= len(m.tasks) && m.cursor > 0 {
		m.cursor--
	}

	// Return to list mode
	m.mode = listMode
	m.taskContent = ""

	return m
}

// Init is called once when the program starts
// It can return a command to run (we don't need any for now)
func (m model) Init() tea.Cmd {
	// No initial commands needed for this simple app
	return nil
}

// Update is called when something happens (like a key press)
// This is where we handle user input and update our model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle reload tasks message
	case reloadTasksMsg:
		// Reload tasks from all configured directories
		tasks, err := loadTasksFromDirectories(m.configDirs)
		m.tasks = tasks
		m.err = err
		m.mode = listMode
		m.taskContent = ""
		// Reset cursor to top
		m.cursor = 0
		return m, nil

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {

		// Quit keys
		case "q", "ctrl+c":
			return m, tea.Quit

		// Navigation and actions depend on current mode
		case "esc":
			if m.mode == taskViewMode {
				m.mode = listMode
				m.taskContent = ""
			} else if m.mode == confirmDeleteMode {
				// Cancel deletion
				m.mode = taskViewMode
			} else if m.mode == searchMode {
				// Exit search mode
				m.mode = listMode
				m.searchQuery = ""
				m.filteredTasks = nil
				m.cursor = 0
			}

		case "enter":
			if m.mode == listMode && len(m.tasks) > 0 {
				// Read the task file content
				visibleTasks := m.visibleTasks()
				if m.cursor < len(visibleTasks) {
					content, err := os.ReadFile(visibleTasks[m.cursor].fullPath)
					if err != nil {
						m.err = fmt.Errorf("failed to read task: %w", err)
					} else {
						m.mode = taskViewMode
						m.taskContent = string(content)
					}
				}
			} else if m.mode == searchMode && len(m.visibleTasks()) > 0 {
				// View selected task from search results
				visibleTasks := m.visibleTasks()
				if m.cursor < len(visibleTasks) {
					content, err := os.ReadFile(visibleTasks[m.cursor].fullPath)
					if err != nil {
						m.err = fmt.Errorf("failed to read task: %w", err)
					} else {
						m.mode = taskViewMode
						m.taskContent = string(content)
					}
				}
			}

		case "e":
			if m.mode == taskViewMode && len(m.tasks) > 0 {
				// Edit the current task
				return m, m.editTask()
			}

		case "n":
			if m.mode == listMode {
				// Create a new task
				return m, m.createTask()
			} else if m.mode == confirmDeleteMode {
				// Cancel deletion
				m.mode = taskViewMode
			}

		case "d":
			if m.mode == taskViewMode && len(m.tasks) > 0 {
				// Show delete confirmation
				m.mode = confirmDeleteMode
			}

		case "y":
			if m.mode == confirmDeleteMode && len(m.tasks) > 0 {
				// Confirm deletion
				return m.deleteTask(), nil
			}

		case "/":
			if m.mode == listMode {
				// Enter search mode
				m.mode = searchMode
				m.searchQuery = ""
				m.cursor = 0
			}

		case "backspace":
			if m.mode == searchMode && len(m.searchQuery) > 0 {
				// Remove last character from search query
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.filterTasks()
			}

		// Move up (only in list or search mode)
		case "up", "k":
			if (m.mode == listMode || m.mode == searchMode) && m.cursor > 0 {
				m.cursor--
			}

		// Move down (only in list or search mode)
		case "down", "j":
			visibleTasks := m.visibleTasks()
			if (m.mode == listMode || m.mode == searchMode) && m.cursor < len(visibleTasks)-1 {
				m.cursor++
			}

		default:
			// In search mode, add typed characters to search query
			if m.mode == searchMode {
				// Only allow printable characters
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
					m.filterTasks()
				}
			}
		}
	}

	// Return the updated model (and no command)
	return m, nil
}

// View renders the UI based on the current model state
// This function is called after every Update
func (m model) View() string {
	// If in confirmation mode, show confirmation dialog
	if m.mode == confirmDeleteMode {
		return m.renderDeleteConfirmation()
	}

	// If viewing a task, show task content
	if m.mode == taskViewMode {
		return m.renderTaskView()
	}

	// Otherwise, show the task list
	return m.renderListView()
}

// renderDeleteConfirmation shows a confirmation dialog for deleting a task
func (m model) renderDeleteConfirmation() string {
	s := "Delete Task\n"
	s += "======================\n\n"

	if m.cursor < len(m.tasks) {
		s += fmt.Sprintf("Are you sure you want to delete this task?\n\n")
		s += fmt.Sprintf("File: %s\n", m.tasks[m.cursor].name)
		if m.tasks[m.cursor].metadata.Title != "" {
			s += fmt.Sprintf("Title: %s\n", m.tasks[m.cursor].metadata.Title)
		}
		s += fmt.Sprintf("Path: %s\n", m.tasks[m.cursor].fullPath)
	}

	s += "\n----------------------\n"
	s += "This action cannot be undone!\n\n"
	s += "y: yes, delete • esc/n: cancel • q: quit\n"

	return s
}

// renderTaskView displays the content of a single task
func (m model) renderTaskView() string {
	s := "Task Viewer\n"
	s += "======================\n\n"

	if m.cursor < len(m.tasks) {
		s += fmt.Sprintf("File: %s\n", m.tasks[m.cursor].name)
		s += fmt.Sprintf("Path: %s\n", m.tasks[m.cursor].fullPath)
		s += "----------------------\n\n"
	}

	s += m.taskContent
	s += "\n\n----------------------\n"
	s += "esc: back • e: edit • d: delete • q: quit\n"

	return s
}

// renderListView displays the list of tasks
func (m model) renderListView() string {
	// Build the UI string
	var title string
	if m.mode == searchMode {
		title = fmt.Sprintf("Search: %s", m.searchQuery)
		if m.searchQuery == "" {
			title = "Search (type to filter)"
		}
	} else if len(m.configDirs) == 1 {
		title = fmt.Sprintf("Task Manager - %s", m.configDirs[0])
	} else {
		title = fmt.Sprintf("Task Manager - %d directories", len(m.configDirs))
	}
	s := title + "\n"
	s += "======================\n\n"

	// If there was an error loading tasks, display it
	if m.err != nil {
		s += fmt.Sprintf("Error: %v\n\n", m.err)
		s += "Make sure the configured directories exist:\n"
		for _, dir := range m.configDirs {
			s += fmt.Sprintf("  - %s\n", dir)
		}
		s += "\nPress 'q' to quit\n"
		return s
	}

	// Get visible tasks (filtered or all)
	visibleTasks := m.visibleTasks()

	// If no tasks found, show a helpful message
	if len(m.tasks) == 0 {
		s += "No markdown files found in:\n"
		for _, dir := range m.configDirs {
			s += fmt.Sprintf("  - %s\n", dir)
		}
		s += "\nAdd some .md files to get started!\n"
		s += "\nPress 'q' to quit\n"
		return s
	}

	// If in search mode and no results
	if m.mode == searchMode && len(visibleTasks) == 0 {
		s += "No tasks match your search.\n"
		s += "\nesc: clear search • q: quit\n"
		return s
	}

	// Render each visible task in our list
	for i, task := range visibleTasks {
		// Is the cursor pointing at this task?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Get status, using default if not set
		status := task.metadata.Status
		if status == "" {
			status = m.config.GetDefaultStatus()
		}

		// Status indicator (if available)
		statusIndicator := m.config.GetStatusIndicator(status)

		// Priority emoji (if available)
		priorityEmoji := getPriorityEmoji(task.metadata.Priority)
		if priorityEmoji != "" {
			priorityEmoji = priorityEmoji + " "
		}

		// Use title from frontmatter if available, otherwise use filename
		displayName := task.name
		if task.metadata.Title != "" {
			displayName = task.metadata.Title
		}

		// Format the modification time nicely
		modTime := task.modTime.Format("2006-01-02 15:04")

		// Build the row with status and priority
		row := fmt.Sprintf("%s %s %s%-40s  %s", cursor, statusIndicator, priorityEmoji, displayName, modTime)

		// If we have multiple directories, show which one this task is from
		if m.showDirInfo {
			row += fmt.Sprintf("  [%s]", task.sourceDir)
		}

		s += row + "\n"
	}

	// Footer with instructions
	s += "\n"
	
	if m.mode == searchMode {
		s += fmt.Sprintf("Showing %d of %d tasks", len(visibleTasks), len(m.tasks))
		s += " • esc: clear search • enter: view • q: quit\n"
	} else {
		s += fmt.Sprintf("Showing %d tasks", len(m.tasks))
		if len(m.configDirs) > 1 {
			s += fmt.Sprintf(" from %d directories", len(m.configDirs))
		}
		s += " • / search • ↑/k up • ↓/j down • enter view • n new • q quit\n"
	}

	return s
}

func main() {
	// Create a new Bubble Tea program with our model
	// WithAltScreen() enables alternate screen mode - the app takes over
	// the full terminal and restores it when you quit (like vim, lazygit, etc.)
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support (optional, but nice!)
	)

	// Start the program and handle any errors
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
