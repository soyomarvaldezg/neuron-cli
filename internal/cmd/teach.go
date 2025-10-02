// Package cmd implements the command line interface for Neuron CLI.
package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/soyomarvaldezg/neuron-cli/internal/db"
	"github.com/soyomarvaldezg/neuron-cli/internal/study"
	"github.com/spf13/cobra"
)

var teachCmd = &cobra.Command{
	Use:   "teach [topic]",
	Short: "Deepen your understanding of a topic using the Feynman Technique",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		database, err := db.GetDB()
		if err != nil {
			return err
		}

		noteToTeach, err := db.GetNoteByTitleOrFilename(database, topic)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("Sorry, I couldn't find a note matching '%s'.\n", topic)
				return nil
			}
			return err
		}

		fmt.Printf("--- Starting Feynman Session on: %s ---\n", noteToTeach.Title)
		fmt.Println("Explain the topic in simple terms. The AI will ask questions.")
		fmt.Println("---------------------------------------------------------------------------------")

		messages := []study.OllamaMessage{
			{
				Role:    "system",
				Content: "You are a curious but intelligent 10-year-old student. The user is your teacher and will explain a concept to you from their notes. Your job is to listen and ask simple questions to help you understand better. If the teacher uses a complex word or jargon, ask what it means. If you don't understand an explanation, ask for an analogy. Start the conversation by asking the teacher to explain the topic of their note.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("The topic of my note is: '%s'. Please ask your first question.", noteToTeach.Title),
			},
		}

		aiColor := color.New(color.FgCyan)
		userColor := color.New(color.FgYellow, color.Bold)

		reader := bufio.NewReader(os.Stdin)

		// Show available commands at start
		helpColor := color.New(color.FgGreen)
		helpColor.Println("\n💡 Tip: Type 'help' anytime to see available commands\n")

		for {
			aiResponse, err := study.SendChatMessage(messages)
			if err != nil {
				return err
			}
			messages = append(messages, aiResponse)

			aiColor.Printf("\n🤖 AI Student: %s\n", aiResponse.Content)
			userColor.Print("You: ")

			userInput, _ := reader.ReadString('\n')
			userInput = strings.TrimSpace(userInput)

			// Check for special commands
			isSpecial, shouldContinue, err := ProcessSpecialCommand(userInput, noteToTeach, &messages)
			if err != nil {
				return err
			}
			if isSpecial {
				if !shouldContinue {
					fmt.Println("Feynman session ended. Great work!")
					break
				}
				continue // Skip adding to messages, show prompt again
			}

			messages = append(messages, study.OllamaMessage{Role: "user", Content: userInput})
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(teachCmd)
}
