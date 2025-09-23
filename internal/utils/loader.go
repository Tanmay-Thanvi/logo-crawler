package utils

import (
	"fmt"
	"os"
	"time"
)

// Loader provides animated console loading indicators
type Loader struct {
	message string
	done    chan bool
}

// NewLoader creates a new loader with a message
func NewLoader(message string) *Loader {
	return &Loader{
		message: message,
		done:    make(chan bool),
	}
}

// Start begins the loading animation
func (l *Loader) Start() {
	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-l.done:
				return
			default:
				fmt.Printf("\r%s %s", frames[i], l.message)
				os.Stdout.Sync()
				time.Sleep(100 * time.Millisecond)
				i = (i + 1) % len(frames)
			}
		}
	}()
}

// Stop stops the loading animation and clears the line
func (l *Loader) Stop() {
	l.done <- true
	fmt.Printf("\r%s\r", "                                                                                ")
	os.Stdout.Sync()
}

// ProgressBar shows a progress bar for a specific task
type ProgressBar struct {
	total   int
	current int
	message string
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, message string) *ProgressBar {
	return &ProgressBar{
		total:   total,
		current: 0,
		message: message,
	}
}

// Update updates the progress bar
func (pb *ProgressBar) Update(current int) {
	pb.current = current
	percentage := float64(current) / float64(pb.total) * 100
	barLength := 30
	filledLength := int(float64(barLength) * percentage / 100)

	bar := ""
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	fmt.Printf("\r%s [%s] %.1f%% (%d/%d)", pb.message, bar, percentage, current, pb.total)
	os.Stdout.Sync()
}

// Complete marks the progress bar as complete
func (pb *ProgressBar) Complete() {
	pb.Update(pb.total)
	fmt.Println() // Move to next line
}
