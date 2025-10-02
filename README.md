# Neuron CLI

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An AI-powered, evidence-based study partner for your Markdown notes, right in your terminal.

---

## Why Neuron CLI? The Science of Learning

Neuron CLI isn't just a flashcard app; it's a powerful learning system built on proven cognitive science principles to maximize retention and understanding.

- **Active Recall (`review`):** An AI generates questions from your notes, forcing you to actively retrieve information from memory, which is dramatically more effective than passive re-reading.
- **Spaced Repetition (`review`, `mix`):** An integrated SRS algorithm schedules notes for review at the optimal moment—right before you're about to forget them.
- **Interleaving (`mix`):** Reviews a random assortment of notes from different topics, forcing your brain to switch contexts and build more flexible, robust knowledge.
- **The Feynman Technique (`teach`):** The AI acts as a curious student, asking you to explain concepts in simple terms. This is the ultimate test of true understanding and instantly reveals gaps in your knowledge.
- **Elaborative Interrogation (`deep-dive`):** The AI acts as a Socratic tutor, asking "why" and "how" questions to help you connect ideas and build deeper mental models.

---

## Features

- **AI-Powered Q&A:** Uses a local LLM (via Ollama) to dynamically generate deep, thoughtful questions and pedagogically sound answers.
- **Spaced Repetition System:** Automatically schedules future reviews based on your performance.
- **Markdown Zettelkasten Support:** Imports and syncs with your folder of Markdown notes.
- **Multiple Study Modes:** Use `review`, `mix`, `teach`, or `deep-dive` for different learning goals.
- **Beautiful Terminal UI:** Renders Markdown beautifully in the terminal with colors and formatting.
- **Clean, Centralized Data Storage:** Manages its database in your system's standard user config directory, keeping your project folders clean.
- **Local First:** Your notes and the AI run entirely on your local machine. No cloud services or subscriptions needed.
- **Interactive Learning Tools:** Built-in commands during `teach` and `deep-dive` sessions let you view notes, ask for explanations, and get help anytime.
- **Flexible Review Options:** Control your study flow with optional note display and brief mode for efficient reviews.

---

## Getting Started

### Prerequisites

1.  **Go:** Version 1.18 or higher.

2.  **SQLite:** The command-line tools. On macOS: `brew install sqlite3`.

3.  **Ollama:** Follow the instructions at [ollama.com](https://ollama.com) to install and run a model. This project was built with `llama3:8b-instruct-q4_K_M`.
    ```bash
    # Pull the model
    ollama pull llama3:8b-instruct-q4_K_M
    # Make sure the Ollama server is running in another terminal
    ollama serve
    ```

### Installation

The recommended method is to use `go install`:

```bash
go install github.com/soyomarvaldezg/neuron-cli@latest
```

Note: Ensure your Go bin directory is in your shell's PATH. This is typically `$(go env GOPATH)/bin`. If the `neuron` command is not found after installation, add `export PATH=$PATH:$(go env GOPATH)/bin` to your `~/.zshrc` or `~/.bash_profile` and restart your terminal.

---

## Your Workflow

### Step 1: Import Your Notes

This is the most important first step. Run the import command and point it to the folder containing your Markdown notes. This will create and populate your central study database.

```bash
neuron import /path/to/your/zettelkasten
```

Neuron CLI will store its database in the standard location for your OS (e.g., `~/.config/neuron-cli` on Linux, `~/Library/Application Support/neuron-cli` on macOS). Run import again anytime you add or change your notes to keep everything in sync.

### Step 2: Start Studying

Use one of the powerful study commands from any directory in your terminal:

#### Daily Review (Spaced Repetition)

```bash
# Standard review - shows due notes with optional full note display
neuron review

# Brief mode - skip the full note display for faster reviews
neuron review --brief

# Review any random note, even if not due
neuron review --any
```

#### Interleaved Practice

```bash
# Standard interleaved review
neuron mix

# Brief mode for faster interleaved sessions
neuron mix --brief
```

#### Deepen Understanding (Feynman Technique)

```bash
# Use a unique keyword from the note's title
neuron teach "parquet"
```

**Interactive Commands Available During `teach` Sessions:**

- `help` or `?` - Show available commands
- `note` or `show note` - Display the full note content
- `explain <topic>` - Ask the AI to explain a specific concept
- `quit` or `exit` - End the session

#### Explore Connections (Elaboration)

```bash
neuron deep-dive "security"
```

**Interactive Commands Available During `deep-dive` Sessions:**

- `help` or `?` - Show available commands
- `note` or `show note` - Display the full note content
- `explain <topic>` - Ask the AI to explain a specific concept
- `quit` or `exit` - End the session

---

## Command Reference

### `neuron import [path]`

Import and sync notes from a directory of Markdown files.

### `neuron review [flags]`

Start a spaced repetition review session with notes that are due.

**Flags:**

- `--any` - Review any random note, even if it's not due
- `--brief` - Skip showing full note, only show question and answer

### `neuron mix [flags]`

Start an interleaved review session with multiple random notes.

**Flags:**

- `--brief` - Skip showing full note, only show question and answer

### `neuron teach [topic]`

Practice the Feynman Technique by explaining a topic to an AI student.

**Interactive commands:** `help`, `note`, `explain <topic>`, `quit`

### `neuron deep-dive [topic]`

Explore a topic deeply with Socratic questioning from an AI tutor.

**Interactive commands:** `help`, `note`, `explain <topic>`, `quit`

---

## Learning Science Behind Neuron CLI

### Better Questions, Better Learning

Neuron CLI generates questions that require **application** and **analysis**, not just memorization. Questions are designed to make you think about:

- How concepts apply to real scenarios
- Relationships between ideas
- Why things work the way they do
- What would happen if conditions changed

### Pedagogically Sound Answers

When you reveal answers, you get:

1. A direct, clear answer (1-2 sentences)
2. Explanation of the "why" or "how"
3. Concrete examples or analogies
4. Connection to broader principles

This structure helps you build deep, interconnected knowledge instead of isolated facts.

### Optional Context

You control when to see the full note, preventing the passive re-reading trap while still having context available when you need it.

---

## Development

The Makefile contains all necessary commands for development.

```bash
# Build the binary into the ./build directory
make build

# Install the binary for global use
make install

# Clean build artifacts
make clean

# Run tests
make test
```

---

## Project Structure

```
neuron-cli/
├── main.go                    # Entry point
├── internal/
│   ├── cmd/                   # CLI commands
│   │   ├── root.go           # Root command
│   │   ├── import.go         # Import command
│   │   ├── review.go         # Review command
│   │   ├── mix.go            # Mix command
│   │   ├── teach.go          # Teach command
│   │   ├── deepdive.go       # Deep-dive command
│   │   ├── render.go         # Markdown rendering
│   │   └── helpers.go        # Interactive command helpers
│   ├── db/                    # Database operations
│   │   └── database.go
│   ├── note/                  # Note data structures
│   │   ├── note.go
│   │   └── parser.go
│   └── study/                 # Learning logic
│       ├── srs.go            # Spaced repetition
│       └── llm.go            # AI interactions
├── Makefile
├── go.mod
└── README.md
```

---

## Tips for Effective Learning

1. **Review daily:** Even 5-10 minutes of daily review is more effective than long, infrequent sessions.

2. **Use brief mode for efficiency:** When you're confident and want to move quickly, use `--brief` to skip optional note displays.

3. **Mix it up:** Use `neuron mix` regularly to practice context switching between topics.

4. **Teach what you learn:** Use `neuron teach` on topics you think you understand well—it will reveal gaps you didn't know existed.

5. **Go deep when stuck:** If you're struggling with a concept, use `neuron deep-dive` to explore it from different angles.

6. **Trust the system:** If a note feels too easy, mark it "Easy." The SRS algorithm will adjust the interval appropriately.

---

## Credits

Neuron CLI is built with these amazing open-source libraries:

- [Cobra](https://github.com/spf13/cobra) for the powerful CLI structure
- [Goldmark](https://github.com/yuin/goldmark) for robust Markdown parsing
- [Glamour](https://github.com/charmbracelet/glamour) for beautiful terminal rendering
- [Fatih/Color](https://github.com/fatih/color) for styled terminal output
- [go-sqlite3](https://github.com/mattn/go-sqlite3) for SQLite database access

---

## License

MIT

---

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

## Roadmap

Future enhancements we're considering:

- Support for different LLM backends
- Custom SRS algorithms
- Statistics and progress tracking
- Export/import of review history

---

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/soyomarvaldezg/neuron-cli/issues) page
2. Open a new issue with detailed information about your problem
3. Include your OS, Go version, and any error messages

---

Happy learning! 🧠
