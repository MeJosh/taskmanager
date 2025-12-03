package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// taskFile represents a markdown file with its metadata
type taskFile struct {
	name      string       // filename
	modTime   time.Time    // last modification time
	fullPath  string       // absolute path to the file
	sourceDir string       // which directory this task came from
	metadata  TaskMetadata // parsed frontmatter metadata
}

// model represents the application state
// In Bubble Tea, the model holds all the data your application needs
type model struct {
	tasks       []taskFile // Our list of task files
	cursor      int        // Which task our cursor is pointing at
	err         error      // Any error encountered while loading files
	configDirs  []string   // The configured task directories
	showDirInfo bool       // Whether to show directory info for each task
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
	}
}// Init is called once when the program starts
// It can return a command to run (we don't need any for now)
func (m model) Init() tea.Cmd {
	// No initial commands needed for this simple app
	return nil
}

// Update is called when something happens (like a key press)
// This is where we handle user input and update our model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {

		// Quit keys
		case "q", "ctrl+c":
			return m, tea.Quit

		// Move up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// Move down
		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		}
	}

	// Return the updated model (and no command)
	return m, nil
}

// View renders the UI based on the current model state
// This function is called after every Update
func (m model) View() string {
	// Build the UI string
	var title string
	if len(m.configDirs) == 1 {
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

	// Render each task in our list
	for i, task := range m.tasks {
		// Is the cursor pointing at this task?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Status emoji (if available)
		statusEmoji := getStatusEmoji(task.metadata.Status)

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
		row := fmt.Sprintf("%s %s %s%-40s  %s", cursor, statusEmoji, priorityEmoji, displayName, modTime)

		// If we have multiple directories, show which one this task is from
		if m.showDirInfo {
			row += fmt.Sprintf("  [%s]", task.sourceDir)
		}

		s += row + "\n"
	}

	// Footer with instructions
	s += "\n"
	s += fmt.Sprintf("Showing %d tasks", len(m.tasks))
	if len(m.configDirs) > 1 {
		s += fmt.Sprintf(" from %d directories", len(m.configDirs))
	}
	s += " • ↑/k up • ↓/j down • q quit\n"

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
