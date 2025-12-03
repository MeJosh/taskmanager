package main

import (
	"os"
	"time"

	"github.com/adrg/frontmatter"
)

// TaskMetadata represents the frontmatter fields we care about
type TaskMetadata struct {
	Title    string    `yaml:"title"`
	Status   string    `yaml:"status"`   // todo, in-progress, done
	Priority string    `yaml:"priority"` // low, medium, high
	DueDate  time.Time `yaml:"due_date"`
	Tags     []string  `yaml:"tags"`
	Created  time.Time `yaml:"created"`
}

// parseFrontmatter extracts metadata from a markdown file's frontmatter
func parseFrontmatter(filePath string) (TaskMetadata, error) {
	var meta TaskMetadata

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return meta, err
	}
	defer file.Close()

	// Parse frontmatter (ignore the content body for now)
	_, err = frontmatter.Parse(file, &meta)
	if err != nil {
		// If there's no frontmatter or it's malformed, return empty metadata
		// This is not an error - files without frontmatter are valid
		return TaskMetadata{}, nil
	}

	return meta, nil
}

// getStatusEmoji returns an emoji for the task status
func getStatusEmoji(status string) string {
	switch status {
	case "done", "completed":
		return "[✓]" // Checkmark for completed
	case "in-progress", "doing":
		return "[~]" // Tilde for in-progress
	case "todo":
		return "[ ]" // Empty checkbox for not started
	default:
		return "   " // Three spaces for alignment when no status
	}
}

// getPriorityEmoji returns an emoji for the task priority
func getPriorityEmoji(priority string) string {
	switch priority {
	case "high":
		return "high"
	case "medium":
		return "med "
	case "low":
		return "low "
	default:
		return ""
	}
}
