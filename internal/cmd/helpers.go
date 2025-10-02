// Package cmd implements the command line interface for Neuron CLI.
package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/soyomarvaldezg/neuron-cli/internal/note"
	"github.com/soyomarvaldezg/neuron-cli/internal/study"
)

// ProcessSpecialCommand checks if the user input is a special command
// Returns: (isSpecialCommand, shouldContinue, error)
func ProcessSpecialCommand(input string, currentNote *note.Note, messages *[]study.OllamaMessage) (bool, bool, error) {
	input = strings.TrimSpace(strings.ToLower(input))

	switch {
	case input == "quit" || input == "exit":
		return true, false, nil

	case input == "note" || input == "show note":
		// Show the full note
		fmt.Println("\n📖 Current Note Content:")
		fmt.Println("-----------------------------------------------------------")
		rendered, err := renderMarkdown(currentNote.Content)
		if err != nil {
			fmt.Println(currentNote.Content)
		} else {
			fmt.Println(rendered)
		}
		fmt.Println("-----------------------------------------------------------")
		return true, true, nil

	case strings.HasPrefix(input, "explain "):
		// User wants AI to explain something specific
		topic := strings.TrimPrefix(input, "explain ")
		explainMsg := study.OllamaMessage{
			Role:    "user",
			Content: fmt.Sprintf("Please explain this concept clearly: %s", topic),
		}
		*messages = append(*messages, explainMsg)

		aiResponse, err := study.SendChatMessage(*messages)
		if err != nil {
			return true, true, err
		}
		*messages = append(*messages, aiResponse)

		aiColor := color.New(color.FgMagenta)
		aiColor.Printf("\n🧠 Explanation: %s\n\n", aiResponse.Content)
		return true, true, nil

	case input == "help" || input == "?":
		helpColor := color.New(color.FgGreen)
		helpColor.Println("\n🛠️  Available Commands:")
		fmt.Println("  • 'note' or 'show note' - Display the full note content")
		fmt.Println("  • 'explain <topic>' - Ask the AI to explain a specific concept")
		fmt.Println("  • 'help' or '?' - Show this help message")
		fmt.Println("  • 'quit' or 'exit' - End the session")
		fmt.Println()
		return true, true, nil
	}

	return false, true, nil
}
