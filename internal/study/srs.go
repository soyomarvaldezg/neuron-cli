// Package study contains logic related to the learning process, like SRS and LLM interaction.
package study

import (
	"math"
	"time"

	"github.com/soyomarvaldezg/neuron-cli/internal/note"
)

// Performance ratings
const (
	RatingAgain = 1 // Knew nothing, reset.
	RatingGood  = 2 // Recalled correctly.
	RatingEasy  = 3 // Recalled with no effort.
)

// UpdateSRSData calculates the next review date for a note based on user performance.
// Note that this function is EXPORTED (starts with a capital U).
func UpdateSRSData(n *note.Note, rating int) {
	// 1. If rating is "Again", reset the interval.
	if rating == RatingAgain {
		n.Interval = 1 // Reset to 1 day
		// We slightly decrease the ease factor to acknowledge difficulty
		n.EaseFactor = math.Max(1.3, n.EaseFactor-0.2)
	} else {
		// 2. For "Good" or "Easy", calculate the new interval.
		if n.Interval < 1 {
			n.Interval = 1
		} else if n.Interval < 6 {
			n.Interval = math.Ceil(n.Interval * 1.6)
		} else {
			n.Interval = math.Ceil(n.Interval * n.EaseFactor)
		}

		// 3. Adjust the ease factor. Only "Easy" increases it.
		if rating == RatingEasy {
			n.EaseFactor += 0.15
		}
	}

	// 4. Set the next due date.
	// Interval is in days, so we multiply by 24 hours.
	duration := time.Hour * 24 * time.Duration(n.Interval)
	n.DueDate = time.Now().Add(duration)
}
