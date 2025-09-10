// Package db handles all database interactions for Neuron CLI.
package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/soyomarvaldezg/neuron-cli/internal/note"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

func GetDB(dataSourceName string) (*sql.DB, error) {
	var err error
	once.Do(func() {
		dbInstance, err = sql.Open("sqlite3", dataSourceName)
		if err != nil {
			return
		}
		if err = dbInstance.Ping(); err != nil {
			return
		}
		log.Println("Database connection established.")
		err = createTables(dbInstance)
	})
	return dbInstance, err
}

func createTables(db *sql.DB) error {
	notesTableSQL := `
    CREATE TABLE IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        filename TEXT NOT NULL UNIQUE,
        title TEXT NOT NULL,
        tags TEXT,
        content TEXT NOT NULL,
        created_at TIMESTAMP,
        due_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        interval REAL DEFAULT 1.0,
        ease_factor REAL DEFAULT 2.5
    );`
	_, err := db.Exec(notesTableSQL)
	if err != nil {
		return err
	}
	log.Println("Notes table created or already exists.")
	return nil
}

func InsertNote(db *sql.DB, n *note.Note) error {
	tagsJSON, err := json.Marshal(n.Tags)
	if err != nil {
		return err
	}
	query := `
    INSERT INTO notes (filename, title, tags, content, created_at, due_date, interval, ease_factor)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    ON CONFLICT(filename) DO UPDATE SET
        title=excluded.title,
        tags=excluded.tags,
        content=excluded.content,
        created_at=excluded.created_at;`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(n.Filename, n.Title, string(tagsJSON), n.Content, n.CreatedAt, n.DueDate, n.Interval, n.EaseFactor)
	return err
}

func GetDueNote(db *sql.DB) (*note.Note, error) {
	query := `SELECT id, filename, title, tags, content, created_at, due_date, interval, ease_factor
    FROM notes WHERE due_date <= ? ORDER BY due_date ASC LIMIT 1;`
	row := db.QueryRow(query, time.Now())
	var n note.Note
	var tagsJSON string
	err := row.Scan(&n.ID, &n.Filename, &n.Title, &tagsJSON, &n.Content, &n.CreatedAt, &n.DueDate, &n.Interval, &n.EaseFactor)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(tagsJSON), &n.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags for note %d: %w", n.ID, err)
	}
	return &n, nil
}

// UpdateNoteSRS saves the updated spaced repetition data for a note.
// Note that this function is EXPORTED (starts with a capital U).
func UpdateNoteSRS(db *sql.DB, n *note.Note) error {
	query := `UPDATE notes SET due_date = ?, interval = ?, ease_factor = ? WHERE id = ?;`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(n.DueDate, n.Interval, n.EaseFactor, n.ID)
	return err
}

// GetDueNotes fetches a specified number of random notes that are due for review.
func GetDueNotes(db *sql.DB, limit int) ([]*note.Note, error) {
	query := `SELECT id, filename, title, tags, content, created_at, due_date, interval, ease_factor
    FROM notes
    WHERE due_date <= ?
    ORDER BY RANDOM()
    LIMIT ?;`

	rows, err := db.Query(query, time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*note.Note

	for rows.Next() {
		var n note.Note
		var tagsJSON string

		if err := rows.Scan(&n.ID, &n.Filename, &n.Title, &tagsJSON, &n.Content, &n.CreatedAt, &n.DueDate, &n.Interval, &n.EaseFactor); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(tagsJSON), &n.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags for note %d: %w", n.ID, err)
		}
		notes = append(notes, &n)
	}

	return notes, nil
}

// GetNoteByTitleOrFilename fetches a single note that matches a title or filename.
func GetNoteByTitleOrFilename(db *sql.DB, searchTerm string) (*note.Note, error) {
	// The '%' are wildcards, so we can search for "parquet" instead of the full title.
	query := `SELECT id, filename, title, tags, content, created_at, due_date, interval, ease_factor
    FROM notes
    WHERE title LIKE ? OR filename LIKE ?
    LIMIT 1;`

	searchTerm = "%" + searchTerm + "%" // Wrap search term in wildcards
	row := db.QueryRow(query, searchTerm, searchTerm)

	var n note.Note
	var tagsJSON string

	err := row.Scan(&n.ID, &n.Filename, &n.Title, &tagsJSON, &n.Content, &n.CreatedAt, &n.DueDate, &n.Interval, &n.EaseFactor)
	if err != nil {
		return nil, err // Will be sql.ErrNoRows if nothing is found
	}

	if err := json.Unmarshal([]byte(tagsJSON), &n.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags for note %d: %w", n.ID, err)
	}
	return &n, nil
}

// GetAnyNote fetches a single random note from the entire collection, ignoring the due date.
func GetAnyNote(db *sql.DB) (*note.Note, error) {
	// This query is the same as GetDueNote, but without the WHERE clause.
	query := `SELECT id, filename, title, tags, content, created_at, due_date, interval, ease_factor
    FROM notes
    ORDER BY RANDOM()
    LIMIT 1;`

	row := db.QueryRow(query)

	var n note.Note
	var tagsJSON string

	err := row.Scan(&n.ID, &n.Filename, &n.Title, &tagsJSON, &n.Content, &n.CreatedAt, &n.DueDate, &n.Interval, &n.EaseFactor)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(tagsJSON), &n.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags for note %d: %w", n.ID, err)
	}
	return &n, nil
}
