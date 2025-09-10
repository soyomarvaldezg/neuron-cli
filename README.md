# Neuron CLI

An AI-powered, evidence-based study partner for your Markdown notes, right in your terminal.

## Why Neuron CLI? The Science of Learning

Neuron CLI isn't just a flashcard app; it's a powerful learning system built on proven cognitive science principles to maximize retention and understanding.

- **Active Recall (`review`):** The AI generates questions, forcing you to actively retrieve information from memory, which is dramatically more effective than passive re-reading.
- **Spaced Repetition (`review`, `mix`):** An integrated SRS algorithm schedules notes for review at the optimal moment—right before you're about to forget them.
- **Interleaving (`mix`):** The `mix` command reviews a random assortment of due notes from different topics, forcing your brain to switch contexts and build more flexible, robust knowledge.
- **The Feynman Technique (`teach`):** The AI acts as a curious student, asking you to explain concepts in simple terms. This is the ultimate test of true understanding and instantly reveals gaps in your knowledge.
- **Elaborative Interrogation (`deep-dive`):** The AI acts as a Socratic tutor, asking "why" and "how" questions to help you connect ideas and build deeper mental models.

## Features

- **AI-Powered Q&A:** Uses a local LLM (via Ollama) to dynamically generate questions and concise answers from your notes.
- **Spaced Repetition System:** Automatically schedules future reviews based on your performance.
- **Markdown Zettelkasten Support:** Imports and syncs with your folder of Markdown notes.
- **Multiple Study Modes:** Use `review`, `mix`, `teach`, or `deep-dive` for different learning goals.
- **Beautiful Terminal UI:** Renders Markdown beautifully in the terminal with colors and formatting.
- **Local First:** Your notes and the AI run entirely on your local machine. No cloud services or subscriptions needed.

## Getting Started

### Prerequisites

You must have the following installed and running on your system:

1.  **Go:** Version 1.18 or higher.
2.  **SQLite:** The command-line tools. On macOS, this is easy: `brew install sqlite3`.
3.  **Ollama:** Follow the instructions at [ollama.com](https://ollama.com) to install and run a model. This project was built with `llama3:8b-instruct-q4_K_M`.
    ```bash
    # Pull the model
    ollama pull llama3:8b-instruct-q4_K_M
    # Make sure the Ollama server is running in another terminal
    ollama serve
    ```

### Installation

#### Using Go (Recommended)

```bash
go install github.com/soyomarvaldezg/neuron-cli@latest
```

#### From Source

```bash
# Clone the repository
git clone https://github.com/soyomarvaldezg/neuron-cli.git
cd neuron-cli
# Build and install
make install
```

## Usage

**1. Import Your Notes**

First, import your folder of Markdown notes. This creates a local `neuron.db` file and syncs your notes. Run this anytime you add or change your notes.

```bash
neuron import /path/to/your/zettelkasten
```

**2. Start a Study Session**

Use one of the powerful study commands:

- **Daily Review (Spaced Repetition):**

  ```bash
  neuron review
  ```

  To review a random card even if none are due:

  ```bash
  neuron review --any
  ```

- **Interleaved Practice:**

  ```bash
  neuron mix
  ```

- **Deepen Understanding (Feynman Technique):**

  ```bash
  neuron teach "parquet"
  ```

- **Explore Connections (Elaboration):**
  ```bash
  neuron deep-dive "security"
  ```

## Development

### Requirements

- Go 1.18+
- Ollama
- SQLite

### Build from Source

The `Makefile` contains all necessary commands.

```bash
# Build the binary into the ./build directory
make build

# Run tests (not yet implemented)
make test

# Clean build artifacts
make clean
```

## Credits

Neuron CLI is built with these amazing open-source libraries:

- [Cobra](https://github.com/spf13/cobra) for the powerful CLI structure.
- [Goldmark](https://github.com/yuin/goldmark) for robust Markdown parsing.
- [Glamour](https://github.com/charmbracelet/glamour) for beautiful terminal rendering.
- [Fatih/Color](https://github.com/fatih/color) for styled terminal output.
- [go-sqlite3](https://github.com/mattn/go-sqlite3) for SQLite database access.

## License

MIT
