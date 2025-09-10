// Package db handles all database interactions for Neuron CLI.
package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/soyomarvaldezg/neuron-cli/internal/note"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

// GetDatabasePath determines the correct, centralized path for the database file.
func GetDatabasePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not get user config directory: %w", err)
	}
	appDataDir := filepath.Join(configDir, "neuron-cli")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return "", fmt.Errorf("could not create app data directory: %w", err)
	}
	return filepath.Join(appDataDir, "neuron.db"), nil
}

// GetDB establishes a singleton connection to the SQLite database.
func GetDB() (*sql.DB, error) {
	once.Do(func() {
		dbPath, err := GetDatabasePath()
		if err != nil {
			log.Fatalf("FATAL: Could not determine database path: %v", err)
		}
		dbInstance, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("FATAL: Could not open database at %s: %v", dbPath, err)
		}
		if err = dbInstance.Ping(); err != nil {
			log.Fatalf("FATAL: Could not connect to database at %s: %v", dbPath, err)
		}
		log.Println("Database connection established at:", dbPath)

		if err = createTables(dbInstance); err != nil {
			log.Fatalf("FATAL: Could not create database tables: %v", err)
		}
	})
	return dbInstance, nil
}

func createTables(db *sql.DB) error {
	notesTableSQL := `CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY, filename TEXT NOT NULL UNIQUE, title TEXT NOT NULL, tags TEXT, content TEXT NOT NULL, created_at TIMESTAMP, due_date TIMESTAMP NOT NULL, interval REAL, ease_factor REAL);`
	_, err := db.Exec(notesTableSQL)
	return err
}

func InsertNote(db *sql.DB, n *note.Note) error {
	tagsJSON, _ := json.Marshal(n.Tags)
	query := `INSERT INTO notes (filename, title, tags, content, created_at, due_date, interval, ease_factor) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(filename) DO UPDATE SET title=excluded.title, tags=excluded.tags, content=excluded.content, created_at=excluded.created_at;`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(n.Filename, n.Title, string(tagsJSON), n.Content, n.CreatedAt, n.DueDate, n.Interval, n.EaseFactor)
	return err
}

func GetDueNote(db *sql.DB) (*note.Note, error) {
	query := `SELECT id, filename, title, tags, content, created_at, due_date, interval, ease_factor FROM notes WHERE due_date <= ? ORDER BY due_date ASC LIMIT 1;`
	row := db.QueryRow(query, time.Now())
	return scanNote(row)
}

func GetDueNotes(db *sql.DB, limit int) ([]*note.Note, error) {
	query := `SELECT id, filename, title, tags, content, created_at, due_date, interval, ease_factor FROM notes WHERE due_date <= ? ORDER BY RANDOM() LIMIT ?;`
	rows, err := db.Query(query, time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []*note.Note
	for rows.Next() {
		note, err := scanNote(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func GetAnyNote(db *sql.DB) (*note.Note, error) {
	query := `SELECT id, filename, title, tags, content, created_at, due_date, interval, ease_factor FROM notes ORDER BY RANDOM() LIMIT 1;`
	row := db.QueryRow(query)
	return scanNote(row)
}

func GetNoteByTitleOrFilename(db *sql.DB, searchTerm string) (*note.Note, error) {
	query := `SELECT id, filename, title, tags, content, created_at, due_date, interval, ease_factor FROM notes WHERE title LIKE ? OR filename LIKE ? LIMIT 1;`
	row := db.QueryRow(query, "%"+searchTerm+"%", "%"+searchTerm+"%")
	return scanNote(row)
}

func UpdateNoteSRS(db *sql.DB, n *note.Note) error {
	query := `UPDATE notes SET due_date = ?, interval = ?, ease_factor = ? WHERE id = ?;`
	_, err := db.Exec(query, n.DueDate, n.Interval, n.EaseFactor, n.ID)
	return err
}

// scanNote is a helper to reduce code duplication when scanning a single row into a Note struct.
type scannable interface {
	Scan(dest ...any) error
}

func scanNote(row scannable) (*note.Note, error) {
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
