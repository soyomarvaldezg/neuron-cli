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
	"github.com/soyomarvaldezg/neuron-cli/internal/study"
	"github.com/spf13/cobra"
)

const reviewLimit = 3 // Number of notes to review in one 'mix' session

var mixBrief bool

var mixCmd = &cobra.Command{
	Use:   "mix",
	Short: "Start an interleaved review session with random due notes",
	Long: `Starts a review session with a small number of randomly selected notes
that are currently due. This helps improve memory by forcing context switching.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.GetDB()
		if err != nil {
			return err
		}

		notes, err := db.GetDueNotes(database, reviewLimit)
		if err != nil {
			if err == sql.ErrNoRows || len(notes) == 0 {
				fmt.Println("🎉 No notes are due for review. Great job!")
				return nil
			}
			return err
		}

		fmt.Printf("--- Starting Interleaved Review Session (%d notes) ---\n", len(notes))
		reader := bufio.NewReader(os.Stdin)

		// Loop through each randomly selected note
		for i, dueNote := range notes {
			fmt.Printf("\n--- Card %d of %d ---\n", i+1, len(notes))

			fmt.Println("🧠 Generating question...")
			question, err := study.GenerateQuestion(dueNote)
			if err != nil {
				fmt.Printf("Error generating question for %s: %v. Skipping.\n", dueNote.Title, err)
				continue
			}

			fmt.Printf("\n🤔 Question: %s\n", question)
			fmt.Print("   (Press Enter to reveal concise answer)")
			_, _ = reader.ReadString('\n')

			fmt.Println("\n🤖 Generating concise answer...")
			conciseAnswer, err := study.GenerateAnswer(question, dueNote)
			if err != nil {
				fmt.Printf("Error generating answer for %s: %v. Skipping.\n", dueNote.Title, err)
				continue
			}

			fmt.Println("\n💡 Concise Answer:")
			fmt.Println("-----------------------------------------------------------")
			fmt.Println(conciseAnswer)
			fmt.Println("-----------------------------------------------------------")

			// Only ask about showing the full note if not in brief mode
			if !mixBrief {
				fmt.Print("\n📖 See full note? (y/n): ")
				showNote, _ := reader.ReadString('\n')
				showNote = strings.TrimSpace(strings.ToLower(showNote))

				if showNote == "y" || showNote == "yes" {
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
				}
			}

			var rating int
			for {
				fmt.Print("\nHow well did you recall this? (1=Again, 2=Good, 3=Easy): ")
				input, _ := reader.ReadString('\n')
				rating, err = strconv.Atoi(strings.TrimSpace(input))
				if err == nil && (rating >= 1 && rating <= 3) {
					break
				}
				fmt.Println("Invalid input. Please enter 1, 2, or 3.")
			}

			study.UpdateSRSData(dueNote, rating)
			if err := db.UpdateNoteSRS(database, dueNote); err != nil {
				return fmt.Errorf("failed to update note schedule: %w", err)
			}
			days := int(math.Ceil(time.Until(dueNote.DueDate).Hours() / 24))
			fmt.Printf("✓ Scheduled for review in about %d day(s).\n", days)
		}

		fmt.Println("\n--- Interleaved session complete! ---")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mixCmd)
	mixCmd.Flags().BoolVar(&mixBrief, "brief", false, "Skip showing full note, only show Q&A")
}
