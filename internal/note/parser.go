// Package note defines the core data structure for a note and its parser.
package note

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/text"
)

// ParseFile reads a markdown file, parses its frontmatter and content, and returns a Note struct.
func ParseFile(path string) (*Note, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
	)

	var buf bytes.Buffer
	reader := text.NewReader(contentBytes)
	doc := md.Parser().Parse(reader)

	md.Renderer().Render(&buf, contentBytes, doc)

	metaData := doc.OwnerDocument().Meta()

	note := &Note{
		Filename:   path,
		Content:    string(contentBytes),
		EaseFactor: 2.5,
		Interval:   1.0,
		DueDate:    time.Now(),
	}

	if title, ok := metaData["title"].(string); ok {
		note.Title = title
	} else {
		note.Title = findFirstH1(string(contentBytes))
	}

	if tags, ok := metaData["Tags"].([]any); ok { // Changed to 'any'
		for _, t := range tags {
			if tagStr, ok := t.(string); ok {
				note.Tags = append(note.Tags, tagStr)
			}
		}
	}

	if createdStr, ok := metaData["Created"].(string); ok {
		t, err := time.Parse("2006-01-02", createdStr)
		if err == nil {
			note.CreatedAt = t
		}
	}

	return note, nil
}

// findFirstH1 scans content for the first line starting with "# ".
func findFirstH1(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if trimmedLine, found := strings.CutPrefix(line, "# "); found { // Changed to CutPrefix
			return strings.TrimSpace(trimmedLine)
		}
	}
	return "Untitled"
}
