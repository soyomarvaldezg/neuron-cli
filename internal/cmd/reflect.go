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

var reflectCmd = &cobra.Command{
	Use:   "reflect [topic]",
	Short: "Deepen understanding through Socratic questioning",
	Long: `Starts a reflection session where the AI acts as a devil's advocate
to challenge your assumptions and explore edge cases. This helps identify
weaknesses in your understanding and deepens your knowledge.

This implements your "Red Team Pattern" by:
1. Having you explain a concept in your own words
2. Challenging your assumptions and exploring edge cases
3. Encouraging critical thinking about limitations and alternatives`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		database, err := db.GetDB()
		if err != nil {
			return err
		}

		noteToReflect, err := db.GetNoteByTitleOrFilename(database, topic)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("Sorry, I couldn't find a note matching '%s'.\n", topic)
				return nil
			}
			return err
		}

		fmt.Printf("--- Starting Reflection Session on: %s ---\n", noteToReflect.Title)
		fmt.Println("I'll act as a devil's advocate to challenge your understanding.")
		fmt.Println("This helps identify assumptions, explore edge cases, and deepen your knowledge.")
		fmt.Println("---------------------------------------------------------------------------------")

		reader := bufio.NewReader(os.Stdin)

		// Show available commands at start
		helpColor := color.New(color.FgGreen)
		helpColor.Println("\n💡 Tip: Type 'help' anytime to see available commands\n")

		// First round: Get initial explanation
		fmt.Print("\n📝 Explain the concept in your own words: ")
		userExplanation, _ := reader.ReadString('\n')
		userExplanation = strings.TrimSpace(userExplanation)

		// Check for special commands
		if strings.ToLower(userExplanation) == "help" || strings.ToLower(userExplanation) == "?" {
			helpColor := color.New(color.FgGreen)
			helpColor.Println("\n🛠️  Available Commands:")
			fmt.Println("  • 'help' or '?' - Show this help message")
			fmt.Println("  • 'note' or 'show note' - Display the full note content")
			fmt.Println("  • 'quit' or 'exit' - End the session")
			fmt.Println("  • Type your explanation to begin reflection")
			fmt.Println()
			// Recursively call the function to get actual explanation
			return cmd.RunE(cmd, args)
		}

		if strings.ToLower(userExplanation) == "quit" || strings.ToLower(userExplanation) == "exit" {
			fmt.Println("Reflection session ended. Good work on critical thinking!")
			return nil
		}

		if strings.ToLower(userExplanation) == "note" || strings.ToLower(userExplanation) == "show note" {
			fmt.Println("\n📖 Full Note Content:")
			fmt.Println("-----------------------------------------------------------")
			rendered, err := renderMarkdown(noteToReflect.Content)
			if err != nil {
				fmt.Println(noteToReflect.Content)
			} else {
				fmt.Println(rendered)
			}
			fmt.Println("-----------------------------------------------------------")
			// Recursively call to get actual explanation
			return cmd.RunE(cmd, args)
		}

		if userExplanation == "" {
			fmt.Println("Please provide an explanation or type a command.")
			// Recursively call to get actual explanation
			return cmd.RunE(cmd, args)
		}

		// Now we have the initial explanation, start the reflection loop
		for {
			// Generate reflection challenges based on current explanation
			fmt.Println("\n🔍 Generating reflection challenges...")
			challenges, err := study.GenerateReflectionChallenges(userExplanation, noteToReflect.Content)
			if err != nil {
				return fmt.Errorf("failed to generate reflection challenges: %w", err)
			}

			// Display challenges
			fmt.Println("\n" + strings.Repeat("=", 60))
			fmt.Println("🎯 REFLECTION CHALLENGES")
			fmt.Println(strings.Repeat("=", 60))

			challengeColor := color.New(color.FgCyan)
			fmt.Print("\n")
			challengeColor.Println(challenges)

			fmt.Println(strings.Repeat("=", 60))

			// Ask if user wants to respond to challenges
			fmt.Print("\n💭 Would you like to respond to these challenges? (y/n): ")
			responseInput, _ := reader.ReadString('\n')
			responseInput = strings.TrimSpace(strings.ToLower(responseInput))

			if responseInput == "y" || responseInput == "yes" {
				fmt.Println("\n💬 Share your thoughts on these challenges:")
				userResponse, _ := reader.ReadString('\n')
				userResponse = strings.TrimSpace(userResponse)

				if userResponse != "" {
					fmt.Println("\n🤝 Your response has been noted. This deeper reflection strengthens your understanding!")
					// Update the explanation with the user's response for next round
					userExplanation = userExplanation + "\n\nReflection: " + userResponse
				}
			}

			// Ask if user wants to continue with another reflection round
			fmt.Print("\n🔄 Continue with another reflection round? (y/n): ")
			continueInput, _ := reader.ReadString('\n')
			continueInput = strings.TrimSpace(strings.ToLower(continueInput))

			if continueInput == "n" || continueInput == "no" {
				fmt.Println("Reflection session ended. Great work on critical thinking!")
				break
			}

			// For subsequent rounds, ask if they want to refine their explanation or explore new aspects
			fmt.Print("\n📝 Would you like to (1) refine your explanation or (2) explore new aspects? [1/2]: ")
			choiceInput, _ := reader.ReadString('\n')
			choiceInput = strings.TrimSpace(choiceInput)

			if choiceInput == "1" {
				fmt.Print("\n✏️ Refine your explanation based on the challenges: ")
				refinedExplanation, _ := reader.ReadString('\n')
				refinedExplanation = strings.TrimSpace(refinedExplanation)

				if refinedExplanation != "" {
					userExplanation = refinedExplanation
					fmt.Println("✅ Your explanation has been updated for the next reflection round.")
				}
			} else if choiceInput == "2" {
				fmt.Print("\n🔍 What aspect of the concept would you like to explore next? ")
				newAspect, _ := reader.ReadString('\n')
				newAspect = strings.TrimSpace(newAspect)

				if newAspect != "" {
					userExplanation = userExplanation + "\n\nNew aspect to explore: " + newAspect
					fmt.Println("✅ New aspect added for the next reflection round.")
				}
			} else {
				// Default to continuing with current explanation
				fmt.Println("✅ Continuing with your current explanation for the next reflection round.")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reflectCmd)
}
