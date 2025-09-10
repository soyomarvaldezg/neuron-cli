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

var deepDiveCmd = &cobra.Command{
	Use:   "deep-dive [topic]",
	Short: "Explore a topic's connections using Socratic questioning",
	Long: `Starts an interactive session where the AI acts as a Socratic tutor.
It will ask you "why" and "how" questions about a specific note to help you
explore its connections and deepen your understanding.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		database, err := db.GetDB("neuron.db")
		if err != nil {
			return err
		}

		noteToExplore, err := db.GetNoteByTitleOrFilename(database, topic)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("Sorry, I couldn't find a note matching '%s'.\n", topic)
				return nil
			}
			return err
		}

		fmt.Printf("--- Starting Deep Dive Session on: %s ---\n", noteToExplore.Title)
		fmt.Println("The AI tutor will ask you questions. Explore your thoughts freely. Type 'quit' to end.")
		fmt.Println("---------------------------------------------------------------------------------")

		// --- THIS IS THE KEY DIFFERENCE: The System Prompt ---
		messages := []study.OllamaMessage{
			{
				Role:    "system",
				Content: "You are a Socratic tutor. Your goal is to help the user think more deeply about their notes. Read the initial text provided by the user. Then, ask one insightful 'why' or 'how' question at a time to encourage them to elaborate and make connections. Do not provide answers, only ask questions based on their text and their responses.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Here is the content of my note titled '%s'. Please read it and ask me your first insightful question to begin our deep dive.\n\n---\n%s\n---", noteToExplore.Title, noteToExplore.Content),
			},
		}

		aiColor := color.New(color.FgCyan)
		userColor := color.New(color.FgYellow, color.Bold)

		reader := bufio.NewReader(os.Stdin)
		for {
			aiResponse, err := study.SendChatMessage(messages)
			if err != nil {
				return err
			}
			messages = append(messages, aiResponse)

			aiColor.Printf("\n🤔 Tutor: %s\n", aiResponse.Content)
			userColor.Print("Your Thoughts: ")

			userInput, _ := reader.ReadString('\n')
			userInput = strings.TrimSpace(userInput)

			if strings.ToLower(userInput) == "quit" {
				fmt.Println("Deep dive session ended. Excellent reflection!")
				break
			}

			messages = append(messages, study.OllamaMessage{Role: "user", Content: userInput})
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deepDiveCmd)
}
