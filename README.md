# Neuron CLI

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

- **AI-Powered Q\&A:** Uses a local LLM (via Ollama) to dynamically generate questions and concise answers.
- **Spaced Repetition System:** Automatically schedules future reviews based on your performance.
- **Markdown Zettelkasten Support:** Imports and syncs with your folder of Markdown notes.
- **Multiple Study Modes:** Use `review`, `mix`, `teach`, or `deep-dive` for different learning goals.
- **Beautiful Terminal UI:** Renders Markdown beautifully in the terminal with colors and formatting.
- **Clean, Centralized Data Storage:** Manages its database in your system's standard user config directory, keeping your project folders clean.
- **Local First:** Your notes and the AI run entirely on your local machine. No cloud services or subscriptions needed.

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

Note: Ensure your Go bin directory is in your shell's PATH. This is typically $(go env GOPATH)/bin. If the `neuron` command is not found after installation, add `export PATH=$PATH:$(go env GOPATH)/bin to your ~/.zshrc or ~/.bash_profile and restart your terminal.

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

**Daily Review (Spaced Repetition):**

```bash
neuron review
```

To review a random card even if none are due:

```bash
neuron review --any
```

**Interleaved Practice:**

```bash
neuron mix
```

**Deepen Understanding (Feynman Technique):**

```bash
# Use a unique keyword from the note's title
neuron teach "parquet"
```

**Explore Connections (Elaboration):**

```bash
neuron deep-dive "security"
```

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
```

---

## Credits

Neuron CLI is built with these amazing open-source libraries:

- Cobra for the powerful CLI structure.
- Goldmark for robust Markdown parsing.
- Glamour for beautiful terminal rendering.
- Fatih/Color for styled terminal output.
- go-sqlite3 for SQLite database access.

---

## License

MIT
