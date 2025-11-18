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
	"github.com/soyomarvaldezg/neuron-cli/internal/note"
	"github.com/soyomarvaldezg/neuron-cli/internal/study"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workflowCmd)

	// Define the flags for workflow command
	workflowCmd.Flags().StringP("phase", "p", "foundational", "Phase of the workflow to run (foundational, verification, extension)")
	workflowCmd.Flags().StringP("question-type", "q", "mixed", "Type of questions to generate (factual, conceptual, application, mixed)")
}

var workflowCmd = &cobra.Command{
	Use:   "workflow [topic]",
	Short: "Guided learning through your three-phase framework",
	Long: `Implements your complete three-phase learning framework:
Phase 1: Build Foundational Competence
Phase 2: Metacognitive Verification
Phase 3: Use AI to Extend

Each phase provides specific activities to optimize learning.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		database, err := db.GetDB()
		if err != nil {
			return err
		}

		noteToWorkflow, err := db.GetNoteByTitleOrFilename(database, topic)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("Sorry, I couldn't find a note matching '%s'.\n", topic)
				return nil
			}
			return err
		}

		// Get flags
		phase, _ := cmd.Flags().GetString("phase")
		qType, _ := cmd.Flags().GetString("question-type")

		// Set defaults with explicit fallback
		var phaseTitle string
		switch strings.ToLower(phase) {
		case "foundational":
			phaseTitle = "Foundational"
		case "verification", "metacognitive":
			phaseTitle = "Verification"
		case "extension", "ai":
			phaseTitle = "Extension"
		default:
			phaseTitle = "Foundational"
		}

		if qType == "" {
			qType = "mixed"
		}

		fmt.Printf("--- Starting %s Phase for: %s ---\n", phaseTitle, noteToWorkflow.Title)
		fmt.Println("This is part of your three-phase learning framework.")
		fmt.Println("---------------------------------------------------------------------------------")

		reader := bufio.NewReader(os.Stdin)

		// Show available commands at start
		helpColor := color.New(color.FgGreen)
		helpColor.Println("\n💡 Tip: Type 'help' anytime to see available commands\n")

		// Run the appropriate phase
		switch strings.ToLower(phase) {
		case "foundational":
			return runFoundationalPhase(reader, noteToWorkflow, qType, database)
		case "verification", "metacognitive":
			return runVerificationPhase(reader, noteToWorkflow, qType, database)
		case "extension", "ai":
			return runExtensionPhase(reader, noteToWorkflow, qType, database)
		default:
			fmt.Printf("Unknown phase: %s. Valid phases are: foundational, verification, extension\n", phase)
			return nil
		}
	},
}

func runFoundationalPhase(reader *bufio.Reader, note *note.Note, qType string, database *sql.DB) error {
	fmt.Println("\n📚 PHASE 1: BUILD FOUNDATIONAL COMPETENCE")
	fmt.Println("Purpose: Develop baseline knowledge to evaluate AI output and reduce cognitive load")
	fmt.Println("Actions: Master fundamentals through traditional study without AI assistance")
	fmt.Println("---------------------------------------------------------------------------------")

	for {
		// Show menu options
		fmt.Println("\n🎯 Foundational Phase Options:")
		fmt.Println("  1. Review basic concepts")
		fmt.Println("  2. Test factual recall")
		fmt.Println("  3. Test conceptual understanding")
		fmt.Println("  4. Test application scenarios")
		fmt.Println("  5. Show full note")
		fmt.Println("  6. Help")
		fmt.Println("  7. Exit phase")

		fmt.Print("\nChoose an option (1-7): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Println("\n🧠 Reviewing basic concepts...")
			question, err := study.GenerateQuestion(note, study.QuestionType(qType))
			if err != nil {
				return fmt.Errorf("failed to generate question: %w", err)
			}

			questionColor := color.New(color.FgCyan)
			questionColor.Printf("\n🤔 Question: %s\n", question)

			fmt.Print("\nPress Enter to see answer...")
			_, _ = reader.ReadString('\n')

			fmt.Println("\n🤖 Generating answer...")
			answer, err := study.GenerateAnswer(question, note)
			if err != nil {
				return fmt.Errorf("failed to generate answer: %w", err)
			}

			answerColor := color.New(color.FgMagenta)
			answerColor.Println("\n💡 Answer:")
			fmt.Println("-----------------------------------------------------------")
			answerColor.Println(answer)
			fmt.Println("-----------------------------------------------------------")

		case "2":
			fmt.Println("\n📝 Testing factual recall...")
			return runSelfTestMode(reader, note, "factual", database)

		case "3":
			fmt.Println("\n🧠 Testing conceptual understanding...")
			return runSelfTestMode(reader, note, "conceptual", database)

		case "4":
			fmt.Println("\n🛠️ Testing application scenarios...")
			return runSelfTestMode(reader, note, "application", database)

		case "5":
			fmt.Println("\n📖 Full Note Content:")
			fmt.Println("-----------------------------------------------------------")
			rendered, err := renderMarkdown(note.Content)
			if err != nil {
				fmt.Println(note.Content)
			} else {
				fmt.Println(rendered)
			}
			fmt.Println("-----------------------------------------------------------")

		case "6":
			helpColor := color.New(color.FgGreen)
			helpColor.Println("\n🛠️  Foundational Phase Help:")
			fmt.Println("  • Options 1-4: Different ways to engage with the material")
			fmt.Println("  • Option 5: Review the full note content")
			fmt.Println("  • Option 6: Show this help message")
			fmt.Println("  • Option 7: Exit this phase and continue learning")
			fmt.Println("  • Type 'menu' to return to this menu")

		case "7":
			fmt.Println("\n✅ Foundational phase completed!")
			fmt.Println("You've built baseline knowledge. Ready for Phase 2: Metacognitive Verification.")
			return nil

		default:
			fmt.Println("\nInvalid option. Please choose 1-7.")
		}
	}
}

func runVerificationPhase(reader *bufio.Reader, note *note.Note, qType string, database *sql.DB) error {
	fmt.Println("\n🔍 PHASE 2: METACOGNITIVE VERIFICATION")
	fmt.Println("Purpose: Use AI as a challenging tutor that forces active thinking")
	fmt.Println("Actions: Reproduce solutions, practice explaining concepts, generate practice problems")
	fmt.Println("---------------------------------------------------------------------------------")

	for {
		// Show menu options
		fmt.Println("\n🎯 Verification Phase Options:")
		fmt.Println("  1. Self-test with factual questions")
		fmt.Println("  2. Self-test with conceptual questions")
		fmt.Println("  3. Self-test with application questions")
		fmt.Println("  4. Reflection mode (Red Team Pattern)")
		fmt.Println("  5. Review with mixed questions")
		fmt.Println("  6. Show full note")
		fmt.Println("  7. Help")
		fmt.Println("  8. Exit phase")

		fmt.Print("\nChoose an option (1-8): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			return runSelfTestMode(reader, note, "factual", database)

		case "2":
			return runSelfTestMode(reader, note, "conceptual", database)

		case "3":
			return runSelfTestMode(reader, note, "application", database)

		case "4":
			return runReflectionMode(reader, note)

		case "5":
			fmt.Println("\n🧠 Reviewing with mixed questions...")
			question, err := study.GenerateQuestion(note, study.QuestionType(qType))
			if err != nil {
				return fmt.Errorf("failed to generate question: %w", err)
			}

			questionColor := color.New(color.FgCyan)
			questionColor.Printf("\n🤔 Question: %s\n", question)

			fmt.Print("\nPress Enter to see answer...")
			_, _ = reader.ReadString('\n')

			fmt.Println("\n🤖 Generating answer...")
			answer, err := study.GenerateAnswer(question, note)
			if err != nil {
				return fmt.Errorf("failed to generate answer: %w", err)
			}

			answerColor := color.New(color.FgMagenta)
			answerColor.Println("\n💡 Answer:")
			fmt.Println("-----------------------------------------------------------")
			answerColor.Println(answer)
			fmt.Println("-----------------------------------------------------------")

		case "6":
			fmt.Println("\n📖 Full Note Content:")
			fmt.Println("-----------------------------------------------------------")
			rendered, err := renderMarkdown(note.Content)
			if err != nil {
				fmt.Println(note.Content)
			} else {
				fmt.Println(rendered)
			}
			fmt.Println("-----------------------------------------------------------")

		case "7":
			helpColor := color.New(color.FgGreen)
			helpColor.Println("\n🛠️  Verification Phase Help:")
			fmt.Println("  • Options 1-3: Test your knowledge with different question types")
			fmt.Println("  • Option 4: Reflection mode to challenge assumptions")
			fmt.Println("  • Option 5: Standard review with mixed questions")
			fmt.Println("  • Option 6: Review the full note content")
			fmt.Println("  • Option 7: Show this help message")
			fmt.Println("  • Option 8: Exit this phase and continue learning")
			fmt.Println("  • Type 'menu' to return to this menu")

		case "8":
			fmt.Println("\n✅ Verification phase completed!")
			fmt.Println("You've challenged your understanding and identified knowledge gaps.")
			fmt.Println("Ready for Phase 3: Use AI to Extend.")
			return nil

		default:
			fmt.Println("\nInvalid option. Please choose 1-8.")
		}
	}
}

func runExtensionPhase(reader *bufio.Reader, note *note.Note, qType string, database *sql.DB) error {
	fmt.Println("\n🚀 PHASE 3: USE AI TO EXTEND")
	fmt.Println("Purpose: Accelerate work while maintaining genuine competence")
	fmt.Println("Actions: Brainstorming, exploring alternatives, optimizing solutions")
	fmt.Println("---------------------------------------------------------------------------------")

	for {
		// Show menu options
		fmt.Println("\n🎯 Extension Phase Options:")
		fmt.Println("  1. Collaborative exploration")
		fmt.Println("  2. Generate edge cases")
		fmt.Println("  3. Optimize solution")
		fmt.Println("  4. Generate alternative approaches")
		fmt.Println("  5. Review with mixed questions")
		fmt.Println("  6. Show full note")
		fmt.Println("  7. Help")
		fmt.Println("  8. Exit phase")

		fmt.Print("\nChoose an option (1-8): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			return runCollaborativeExploration(reader, note)

		case "2":
			return runEdgeCaseGeneration(reader, note)

		case "3":
			return runSolutionOptimization(reader, note)

		case "4":
			return runAlternativeApproaches(reader, note)

		case "5":
			fmt.Println("\n🧠 Reviewing with mixed questions...")
			question, err := study.GenerateQuestion(note, study.QuestionType(qType))
			if err != nil {
				return fmt.Errorf("failed to generate question: %w", err)
			}

			questionColor := color.New(color.FgCyan)
			questionColor.Printf("\n🤔 Question: %s\n", question)

			fmt.Print("\nPress Enter to see answer...")
			_, _ = reader.ReadString('\n')

			fmt.Println("\n🤖 Generating answer...")
			answer, err := study.GenerateAnswer(question, note)
			if err != nil {
				return fmt.Errorf("failed to generate answer: %w", err)
			}

			answerColor := color.New(color.FgMagenta)
			answerColor.Println("\n💡 Answer:")
			fmt.Println("-----------------------------------------------------------")
			answerColor.Println(answer)
			fmt.Println("-----------------------------------------------------------")

		case "6":
			fmt.Println("\n📖 Full Note Content:")
			fmt.Println("-----------------------------------------------------------")
			rendered, err := renderMarkdown(note.Content)
			if err != nil {
				fmt.Println(note.Content)
			} else {
				fmt.Println(rendered)
			}
			fmt.Println("-----------------------------------------------------------")

		case "7":
			helpColor := color.New(color.FgGreen)
			helpColor.Println("\n🛠️  Extension Phase Help:")
			fmt.Println("  • Option 1: Collaborative exploration with AI")
			fmt.Println("  • Option 2: Generate edge cases to test understanding")
			fmt.Println("  • Option 3: Optimize solutions you already understand")
			fmt.Println("  • Option 4: Generate alternative approaches")
			fmt.Println("  • Option 5: Standard review with mixed questions")
			fmt.Println("  • Option 6: Review the full note content")
			fmt.Println("  • Option 7: Show this help message")
			fmt.Println("  • Option 8: Exit this phase and continue learning")
			fmt.Println("  • Type 'menu' to return to this menu")

		case "8":
			fmt.Println("\n✅ Extension phase completed!")
			fmt.Println("You've used AI to extend your knowledge while maintaining competence.")
			fmt.Println("Three-phase framework complete. Great work on comprehensive learning!")
			return nil

		default:
			fmt.Println("\nInvalid option. Please choose 1-8.")
		}
	}
}

// Helper functions for extension phase activities
func runCollaborativeExploration(reader *bufio.Reader, note *note.Note) error {
	fmt.Println("\n🤝 Collaborative Exploration")
	fmt.Println("Share your approach, and I'll help you explore alternatives and improvements.")

	fmt.Print("\nDescribe your current approach or solution: ")
	approach, _ := reader.ReadString('\n')
	approach = strings.TrimSpace(approach)

	if approach == "" {
		fmt.Println("Please provide an approach to explore.")
		return nil
	}

	fmt.Println("\n🤖 Generating collaborative exploration...")
	prompt := fmt.Sprintf(`I'm working on understanding this concept: %s

My current approach is: %s

Help me explore this by:
1. Suggesting alternative approaches I might not have considered
2. Identifying potential improvements or optimizations
3. Highlighting connections to related concepts
4. Pointing out potential edge cases or limitations

Focus on expanding my understanding while building on what I already know.`, note.Title, approach)

	messages := []study.OllamaMessage{
		{Role: "user", Content: prompt},
	}

	response, err := study.SendChatMessage(messages)
	if err != nil {
		return fmt.Errorf("failed to get collaborative exploration: %w", err)
	}

	responseColor := color.New(color.FgCyan)
	responseColor.Println("\n💡 Collaborative Exploration:")
	fmt.Println("-----------------------------------------------------------")
	responseColor.Println(response.Content)
	fmt.Println("-----------------------------------------------------------")

	return nil
}

func runEdgeCaseGeneration(reader *bufio.Reader, note *note.Note) error {
	fmt.Println("\n🔥 Edge Case Generation")
	fmt.Println("I'll generate challenging scenarios to test your understanding.")

	fmt.Println("\n🤖 Generating edge cases...")
	prompt := fmt.Sprintf(`Generate 4 challenging edge cases for this concept: %s

For each edge case, provide:
1. A brief scenario description
2. Why it challenges standard understanding
3. How to address it in practice

Focus on scenarios that might break common assumptions or require deeper thinking.`, note.Title)

	messages := []study.OllamaMessage{
		{Role: "user", Content: prompt},
	}

	response, err := study.SendChatMessage(messages)
	if err != nil {
		return fmt.Errorf("failed to generate edge cases: %w", err)
	}

	responseColor := color.New(color.FgCyan)
	responseColor.Println("\n🔥 Edge Cases:")
	fmt.Println("-----------------------------------------------------------")
	responseColor.Println(response.Content)
	fmt.Println("-----------------------------------------------------------")

	return nil
}

func runSolutionOptimization(reader *bufio.Reader, note *note.Note) error {
	fmt.Println("\n⚡ Solution Optimization")
	fmt.Println("Share your current solution, and I'll suggest improvements.")

	fmt.Print("\nDescribe your current solution: ")
	solution, _ := reader.ReadString('\n')
	solution = strings.TrimSpace(solution)

	if solution == "" {
		fmt.Println("Please provide a solution to optimize.")
		return nil
	}

	fmt.Println("\n🤖 Generating optimization suggestions...")
	prompt := fmt.Sprintf(`I have this solution for %s: %s

Help me optimize it by:
1. Suggesting more efficient approaches
2. Identifying potential bottlenecks or issues
3. Recommending best practices or patterns
4. Suggesting simplifications without losing functionality

Focus on practical improvements that maintain core value.`, note.Title, solution)

	messages := []study.OllamaMessage{
		{Role: "user", Content: prompt},
	}

	response, err := study.SendChatMessage(messages)
	if err != nil {
		return fmt.Errorf("failed to get optimization suggestions: %w", err)
	}

	responseColor := color.New(color.FgCyan)
	responseColor.Println("\n⚡ Optimization Suggestions:")
	fmt.Println("-----------------------------------------------------------")
	responseColor.Println(response.Content)
	fmt.Println("-----------------------------------------------------------")

	return nil
}

func runAlternativeApproaches(reader *bufio.Reader, note *note.Note) error {
	fmt.Println("\n🔄 Alternative Approaches")
	fmt.Println("I'll suggest different ways to approach this concept.")

	fmt.Println("\n🤖 Generating alternative approaches...")
	prompt := fmt.Sprintf(`Suggest 3 alternative approaches to understanding or implementing: %s

For each alternative, provide:
1. A brief description of the approach
2. When it might be more effective than standard methods
3. Potential limitations or trade-offs

Focus on diverse perspectives that might not be immediately obvious.`, note.Title)

	messages := []study.OllamaMessage{
		{Role: "user", Content: prompt},
	}

	response, err := study.SendChatMessage(messages)
	if err != nil {
		return fmt.Errorf("failed to get alternative approaches: %w", err)
	}

	responseColor := color.New(color.FgCyan)
	responseColor.Println("\n🔄 Alternative Approaches:")
	fmt.Println("-----------------------------------------------------------")
	responseColor.Println(response.Content)
	fmt.Println("-----------------------------------------------------------")

	return nil
}

// Helper function to run self-test mode
func runSelfTestMode(reader *bufio.Reader, note *note.Note, qType string, database *sql.DB) error {
	fmt.Printf("\n🧠 Self-Testing with %s questions...\n", qType)

	questionCount := 0
	for {
		questionCount++

		// Generate question with variation hint
		fmt.Printf("🧠 Generating %s question (#%d)...\n", qType, questionCount)

		// Add a small random element to prompt to force variation
		question, err := study.GenerateQuestionWithVariation(note, study.QuestionType(qType), questionCount)
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
			fmt.Println("  • 'quit' or 'exit' - End self-test and return to menu")
			fmt.Println("  • Type your answer to test your knowledge")
			fmt.Println()
			continue
		}

		if strings.ToLower(userInput) == "quit" || strings.ToLower(userInput) == "exit" {
			fmt.Println("Self-test completed. Returning to phase menu.")
			return nil
		}

		if strings.ToLower(userInput) == "note" || strings.ToLower(userInput) == "show note" {
			fmt.Println("\n📖 Full Note Content:")
			fmt.Println("-----------------------------------------------------------")
			rendered, err := renderMarkdown(note.Content)
			if err != nil {
				fmt.Println(note.Content)
			} else {
				fmt.Println(rendered)
			}
			fmt.Println("-----------------------------------------------------------")
			continue
		}

		if strings.ToLower(userInput) == "skip" {
			fmt.Println("Question skipped. Moving to next question.")
			continue
		}

		// User provided an answer, so let's compare it
		if userInput == "" {
			fmt.Println("Please provide an answer or type a command.")
			continue
		}

		// Generate AI answer
		fmt.Println("\n🤖 Generating AI answer for comparison...")
		aiAnswer, err := study.GenerateAnswer(question, note)
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
			fmt.Println("Self-test completed. Returning to phase menu.")
			return nil
		}
	}
}

// Helper function to run reflection mode
func runReflectionMode(reader *bufio.Reader, note *note.Note) error {
	fmt.Println("\n🔍 Reflection Mode (Red Team Pattern)")
	fmt.Println("I'll challenge your assumptions and explore edge cases.")

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
		fmt.Println("  • 'quit' or 'exit' - End reflection and return to menu")
		fmt.Println("  • Type your explanation to begin reflection")
		fmt.Println()
		return runReflectionMode(reader, note)
	}

	if strings.ToLower(userExplanation) == "quit" || strings.ToLower(userExplanation) == "exit" {
		fmt.Println("Reflection completed. Returning to phase menu.")
		return nil
	}

	if strings.ToLower(userExplanation) == "note" || strings.ToLower(userExplanation) == "show note" {
		fmt.Println("\n📖 Full Note Content:")
		fmt.Println("-----------------------------------------------------------")
		rendered, err := renderMarkdown(note.Content)
		if err != nil {
			fmt.Println(note.Content)
		} else {
			fmt.Println(rendered)
		}
		fmt.Println("-----------------------------------------------------------")
		return runReflectionMode(reader, note)
	}

	if userExplanation == "" {
		fmt.Println("Please provide an explanation or type a command.")
		return runReflectionMode(reader, note)
	}

	// Now we have the initial explanation, start the reflection loop
	for {
		// Generate reflection challenges based on current explanation
		fmt.Println("\n🔍 Generating reflection challenges...")
		challenges, err := study.GenerateReflectionChallenges(userExplanation, note.Content)
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
				// Update the explanation with the user's response for the next round
				userExplanation = userExplanation + "\n\nReflection: " + userResponse
			}
		}

		// Ask if user wants to continue with another reflection round
		fmt.Print("\n🔄 Continue with another reflection round? (y/n): ")
		continueInput, _ := reader.ReadString('\n')
		continueInput = strings.TrimSpace(strings.ToLower(continueInput))

		if continueInput == "n" || continueInput == "no" {
			fmt.Println("Reflection completed. Returning to phase menu.")
			return nil
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
			// Default to continuing with the current explanation
			fmt.Println("✅ Continuing with your current explanation for the next reflection round.")
		}
	}
}
