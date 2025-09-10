// Package note defines the core data structure for a note and its parser.
package note

import "time"

// Note represents a single markdown note from your Zettelkasten.
type Note struct {
	ID        int       `db:"id"`
	Filename  string    `db:"filename"`
	Title     string    `db:"title"`
	Tags      []string  // Stored as JSON string in DB
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`

	// Fields for Spaced Repetition
	DueDate    time.Time `db:"due_date"`
	Interval   float64   `db:"interval"`
	EaseFactor float64   `db:"ease_factor"`
}
