package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// model represents the application state
// In Bubble Tea, the model holds all the data your application needs
type model struct {
	tasks  []string // Our list of tasks (currently static)
	cursor int      // Which task our cursor is pointing at
}

// initialModel creates the starting state of our application
func initialModel() model {
	return model{
		// Static list of tasks for Phase 1
		tasks: []string{
			"task1.md - Modified: 2025-12-01",
			"task2.md - Modified: 2025-12-02",
			"task3.md - Modified: 2025-12-03",
			"meeting-notes.md - Modified: 2025-11-30",
			"project-ideas.md - Modified: 2025-11-28",
		},
		cursor: 0, // Start at the first item
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
	s := "Task Manager - Phase 1\n"
	s += "======================\n\n"

	// Render each task in our list
	for i, task := range m.tasks {
		// Is the cursor pointing at this task?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, task)
	}

	// Footer with instructions
	s += "\n"
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
