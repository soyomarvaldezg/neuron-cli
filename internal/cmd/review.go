// Package cmd implements the command line interface for Neuron CLI.
package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/soyomarvaldezg/neuron-cli/internal/db"
	"github.com/soyomarvaldezg/neuron-cli/internal/note" // <-- Make sure this is imported
	"github.com/soyomarvaldezg/neuron-cli/internal/study"
	"github.com/spf13/cobra"
)

// This variable will hold the value of the --any flag.
var reviewAny bool

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Start a spaced repetition review session",
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.GetDB()
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		var dueNote *note.Note // Declare the variable to hold the note

		// --- THIS IS THE NEW LOGIC ---
		if reviewAny {
			fmt.Println("Fetching a random note to review...")
			dueNote, err = db.GetAnyNote(database)
		} else {
			dueNote, err = db.GetDueNote(database)
		}
		// --- END OF NEW LOGIC ---

		if err != nil {
			if err == sql.ErrNoRows {
				if reviewAny {
					fmt.Println("You have no notes in your database to review!")
				} else {
					fmt.Println("🎉 No notes are due for review. Great job!")
				}
				return nil
			}
			return fmt.Errorf("failed to fetch note: %w", err)
		}

		// The rest of the command is exactly the same...
		fmt.Println("🧠 Generating question...")
		question, err := study.GenerateQuestion(dueNote)
		if err != nil {
			return fmt.Errorf("failed to generate question: %w", err)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\n🤔 Question: %s\n", question)
		fmt.Print("   (Press Enter to reveal concise answer)")
		_, _ = reader.ReadString('\n')

		fmt.Println("\n🤖 Generating concise answer...")
		conciseAnswer, err := study.GenerateAnswer(question, dueNote)
		if err != nil {
			return fmt.Errorf("failed to generate answer: %w", err)
		}

		fmt.Println("\n💡 Concise Answer:")
		fmt.Println("-----------------------------------------------------------")
		fmt.Println(conciseAnswer)
		fmt.Println("-----------------------------------------------------------")

		fmt.Print("   (Press Enter again to see the full note for context...)")
		_, _ = reader.ReadString('\n')
		fmt.Println("\n📖 Full Note Context:")
		fmt.Println("-----------------------------------------------------------")

		renderedContent, err := renderMarkdown(dueNote.Content)
		if err != nil {
			fmt.Println("Error rendering markdown, showing raw content:")
			fmt.Println(dueNote.Content)
		} else {
			fmt.Println(renderedContent)
		}

		fmt.Println("-----------------------------------------------------------")

		var rating int
		for {
			fmt.Print("\nHow well did you recall this? (1=Again, 2=Good, 3=Easy): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			rating, err = strconv.Atoi(input)
			if err == nil && (rating >= 1 && rating <= 3) {
				break
			}
			fmt.Println("Invalid input. Please enter 1, 2, or 3.")
		}

		study.UpdateSRSData(dueNote, rating)
		if err := db.UpdateNoteSRS(database, dueNote); err != nil {
			return fmt.Errorf("failed to update note schedule: %w", err)
		}
		nextReview := time.Until(dueNote.DueDate)
		days := int(math.Ceil(nextReview.Hours() / 24))
		fmt.Printf("✓ Good work! This note is scheduled for review in about %d day(s).\n", days)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
	// --- THIS IS THE NEW LINE ---
	// Here we define the --any flag, link it to our 'reviewAny' variable,
	// and provide a help message.
	reviewCmd.Flags().BoolVar(&reviewAny, "any", false, "Review any card, even if it's not due")
}
