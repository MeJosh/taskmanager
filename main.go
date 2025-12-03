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
	name     string    // filename
	modTime  time.Time // last modification time
	fullPath string    // absolute path to the file
}

// model represents the application state
// In Bubble Tea, the model holds all the data your application needs
type model struct {
	tasks     []taskFile // Our list of task files
	cursor    int        // Which task our cursor is pointing at
	err       error      // Any error encountered while loading files
	configDir string     // The configured task directory
}

// loadTasksFromDirectory reads all .md files from the specified directory
func loadTasksFromDirectory(dir string) ([]taskFile, error) {
	// Expand the tilde (~) to the user's home directory
	if dir[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("couldn't get home directory: %w", err)
		}
		dir = filepath.Join(homeDir, dir[2:])
	}

	// Read all files in the directory
	entries, err := os.ReadDir(dir)
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

		tasks = append(tasks, taskFile{
			name:     entry.Name(),
			modTime:  info.ModTime(),
			fullPath: filepath.Join(dir, entry.Name()),
		})
	}

	// Sort by modification time (newest first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].modTime.After(tasks[j].modTime)
	})

	return tasks, nil
}

// initialModel creates the starting state of our application
func initialModel() model {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return model{
			tasks:     nil,
			cursor:    0,
			err:       fmt.Errorf("failed to load config: %w", err),
			configDir: "~/.tasks", // fallback
		}
	}

	// Load tasks from configured directory
	tasks, loadErr := loadTasksFromDirectory(cfg.TaskManager.Directory)

	return model{
		tasks:     tasks,
		cursor:    0,
		err:       loadErr,
		configDir: cfg.TaskManager.Directory,
	}
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
	s := fmt.Sprintf("Task Manager - %s\n", m.configDir)
	s += "======================\n\n"

	// If there was an error loading tasks, display it
	if m.err != nil {
		s += fmt.Sprintf("Error: %v\n\n", m.err)
		s += fmt.Sprintf("Make sure the %s directory exists.\n", m.configDir)
		s += "\nPress 'q' to quit\n"
		return s
	}

	// If no tasks found, show a helpful message
	if len(m.tasks) == 0 {
		s += fmt.Sprintf("No markdown files found in %s\n\n", m.configDir)
		s += "Add some .md files to get started!\n"
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

		// Format the modification time nicely
		modTime := task.modTime.Format("2006-01-02 15:04")

		// Render the row
		s += fmt.Sprintf("%s %-40s  %s\n", cursor, task.name, modTime)
	}

	// Footer with instructions
	s += "\n"
	s += fmt.Sprintf("Showing %d tasks • ", len(m.tasks))
	s += "↑/k up • ↓/j down • q quit\n"

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
