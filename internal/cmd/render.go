// internal/cmd/render.go

// Package cmd implements the command line interface for Neuron CLI.
package cmd

import (
	"github.com/charmbracelet/glamour"
)

// renderMarkdown takes a string of markdown and returns a string
// of beautifully rendered terminal-ready output.
func renderMarkdown(content string) (string, error) {
	// glamour.WithAutoStyle() will automatically detect if the terminal
	// has a light or dark background and choose colors accordingly.
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	if err != nil {
		return "", err
	}

	// Render the markdown content.
	out, err := renderer.Render(content)
	if err != nil {
		return "", err
	}

	return out, nil
}
