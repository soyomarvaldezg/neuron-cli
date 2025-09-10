// Package cmd implements the command line interface for Neuron CLI.
package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color" // <-- NEW IMPORT
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
		fmt.Println("Explain the topic in simple terms. The AI will ask questions. Type 'quit' to end.")
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

		// --- NEW: Create colored printers ---
		aiColor := color.New(color.FgCyan)
		userColor := color.New(color.FgYellow, color.Bold)

		reader := bufio.NewReader(os.Stdin)
		for {
			aiResponse, err := study.SendChatMessage(messages)
			if err != nil {
				return err
			}
			messages = append(messages, aiResponse)

			// --- NEW: Use the colored printers ---
			aiColor.Printf("\n🤖 AI Student: %s\n", aiResponse.Content)
			userColor.Print("You: ")

			userInput, _ := reader.ReadString('\n')
			userInput = strings.TrimSpace(userInput)

			if strings.ToLower(userInput) == "quit" {
				fmt.Println("Feynman session ended. Great work!")
				break
			}

			messages = append(messages, study.OllamaMessage{Role: "user", Content: userInput})
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(teachCmd)
}
