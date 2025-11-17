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

var selfTestQuestionType string

var selfTestCmd = &cobra.Command{
	Use:   "self-test [topic]",
	Short: "Test your knowledge by answering questions before seeing the answer",
	Long: `Starts a self-test session where you answer questions in your own words
before seeing the AI-generated answer. This forces active recall and helps
identify knowledge gaps.

Use --question-type to specify the type of questions:
- factual: Questions about definitions, facts, and specific details
- conceptual: Questions about relationships, principles, and "why" things work
- application: Questions about applying concepts to real scenarios
- mixed: A mix of all question types (default)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		database, err := db.GetDB()
		if err != nil {
			return err
		}

		noteToTest, err := db.GetNoteByTitleOrFilename(database, topic)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("Sorry, I couldn't find a note matching '%s'.\n", topic)
				return nil
			}
			return err
		}

		// Convert string to QuestionType
		qType := study.QuestionType(selfTestQuestionType)
		if qType == "" {
			qType = study.QuestionTypeMixed // Default to mixed
		}

		fmt.Printf("--- Starting Self-Test Session on: %s ---\n", noteToTest.Title)
		fmt.Println("Answer the question in your own words before seeing the AI answer.")
		fmt.Println("This helps identify knowledge gaps and strengthens recall.")
		fmt.Println("---------------------------------------------------------------------------------")

		reader := bufio.NewReader(os.Stdin)

		// Show available commands at start
		helpColor := color.New(color.FgGreen)
		helpColor.Println("\n💡 Tip: Type 'help' anytime to see available commands\n")

		questionCount := 0
		for {
			questionCount++

			// Generate question with variation hint
			fmt.Printf("🧠 Generating %s question (#%d)...\n", qType, questionCount)

			// Add a small random element to prompt to force variation
			question, err := study.GenerateQuestionWithVariation(noteToTest, qType, questionCount)
			if err != nil {
				return fmt.Errorf("failed to generate question: %w", err)
			}

			questionColor := color.New(color.FgCyan)
			questionColor.Printf("\n🤔 Question: %s\n", question)

			// Check for special commands
			fmt.Print("\nType your answer (or 'help' for commands): ")
			userInput, _ := reader.ReadString('\n')
			userInput = strings.TrimSpace(userInput)

			// Check for special commands
			if strings.ToLower(userInput) == "help" || strings.ToLower(userInput) == "?" {
				helpColor := color.New(color.FgGreen)
				helpColor.Println("\n🛠️  Available Commands:")
				fmt.Println("  • 'help' or '?' - Show this help message")
				fmt.Println("  • 'note' or 'show note' - Display the full note content")
				fmt.Println("  • 'skip' - Skip this question")
				fmt.Println("  • 'quit' or 'exit' - End the session")
				fmt.Println("  • Type your answer to test your knowledge")
				fmt.Println()
				continue
			}

			if strings.ToLower(userInput) == "quit" || strings.ToLower(userInput) == "exit" {
				fmt.Println("Self-test session ended. Good work on practicing active recall!")
				break
			}

			if strings.ToLower(userInput) == "note" || strings.ToLower(userInput) == "show note" {
				fmt.Println("\n📖 Full Note Content:")
				fmt.Println("-----------------------------------------------------------")
				rendered, err := renderMarkdown(noteToTest.Content)
				if err != nil {
					fmt.Println(noteToTest.Content)
				} else {
					fmt.Println(rendered)
				}
				fmt.Println("-----------------------------------------------------------")
				continue
			}

			if strings.ToLower(userInput) == "skip" {
				fmt.Println("Question skipped. Moving to the next question.")
				continue
			}

			// User provided an answer, so let's compare it
			if userInput == "" {
				fmt.Println("Please provide an answer or type a command.")
				continue
			}

			// Generate AI answer
			fmt.Println("\n🤖 Generating AI answer for comparison...")
			aiAnswer, err := study.GenerateAnswer(question, noteToTest)
			if err != nil {
				return fmt.Errorf("failed to generate AI answer: %w", err)
			}

			// Compare answers
			fmt.Println("\n🔍 Analyzing your answer...")
			comparison, err := study.CompareAnswers(userInput, aiAnswer, question)
			if err != nil {
				return fmt.Errorf("failed to compare answers: %w", err)
			}

			// Display results
			fmt.Println("\n" + strings.Repeat("=", 60))
			fmt.Println("📊 COMPARISON RESULTS")
			fmt.Println(strings.Repeat("=", 60))

			userColor := color.New(color.FgYellow)
			aiColor := color.New(color.FgMagenta)
			feedbackColor := color.New(color.FgGreen)

			fmt.Print("\n💭 Your Answer: ")
			userColor.Println(userInput)

			fmt.Print("\n🤖 AI Answer: ")
			aiColor.Println(aiAnswer)

			fmt.Print("\n📝 Feedback: ")
			feedbackColor.Println(comparison)

			fmt.Println(strings.Repeat("=", 60))

			// Ask if user wants to continue
			fmt.Print("\nContinue with another question? (y/n): ")
			continueInput, _ := reader.ReadString('\n')
			continueInput = strings.TrimSpace(strings.ToLower(continueInput))

			if continueInput == "n" || continueInput == "no" {
				fmt.Println("Self-test session ended. Great work on practicing active recall!")
				break
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(selfTestCmd)
	selfTestCmd.Flags().StringVar(&selfTestQuestionType, "question-type", "mixed", "Type of question to generate: factual, conceptual, application, mixed")
}
