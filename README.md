# Neuron CLI

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An AI-powered, evidence-based study partner for your Markdown notes, right in your terminal.

---

## Why Neuron CLI? The Science of Learning

Neuron CLI isn't just a flashcard app; it's a powerful learning system built on proven cognitive science principles to maximize retention and understanding.

- **Active Recall (`review`, `self-test`):** An AI generates questions from your notes, forcing you to actively retrieve information from memory, which is dramatically more effective than passive re-reading.
- **Spaced Repetition (`review`, `mix`):** An integrated SRS algorithm schedules notes for review at the optimal moment—right before you're about to forget them.
- **Interleaving (`mix`):** Reviews a random assortment of notes from different topics, forcing your brain to switch contexts and build more flexible, robust knowledge.
- **The Feynman Technique (`teach`):** The AI acts as a curious student, asking you to explain concepts in simple terms. This is the ultimate test of true understanding and instantly reveals gaps in your knowledge.
- **Elaborative Interrogation (`deep-dive`):** The AI acts as a Socratic tutor, asking "why" and "how" questions to help you connect ideas and build deeper mental models.
- **Metacognitive Verification (`self-test`, `reflect`):** Test your understanding before seeing answers and challenge your assumptions with Socratic questioning to identify knowledge gaps.

---

## 🧠 Three-Phase Learning Framework

Neuron CLI implements a research-backed three-phase framework for optimal learning:

```
Phase 1: Foundational → Phase 2: Verification → Phase 3: Extension
```

### Phase 1: Build Foundational Competence

**Purpose:** Develop baseline knowledge to evaluate AI output and reduce cognitive load.

**When to use:**

- Learning new topics from scratch
- Building understanding of fundamentals
- Before using AI for complex tasks

**Commands:**

```
neuron workflow "topic" --phase foundational
neuron review --question-type factual
neuron review --question-type conceptual
```

### Phase 2: Metacognitive Verification

**Purpose:** Use AI as a challenging tutor to force active thinking and identify knowledge gaps.

**When to use:**

- After building foundational knowledge
- Testing your understanding
- Challenging your assumptions

**Commands:**

```
neuron workflow "topic" --phase verification
neuron self-test "topic" --question-type conceptual
neuron reflect "topic"
```

### Phase 3: Use AI to Extend

**Purpose:** Accelerate work while maintaining genuine competence.

**When to use:**

- After mastering a topic
- Exploring advanced concepts
- Brainstorming alternatives and optimizing solutions

**Commands:**

```
neuron workflow "topic" --phase extension
```

---

## 📚 Question Types

Neuron CLI supports four types of questions, each targeting different cognitive levels:

| Type            | What It Tests                              | Best For                             | Example                                                |
| --------------- | ------------------------------------------ | ------------------------------------ | ------------------------------------------------------ |
| **factual**     | Definitions, facts, specific details       | Building foundational knowledge      | "What is the definition of a binary search tree?"      |
| **conceptual**  | Relationships, principles, "why" questions | Understanding how things work        | "Why does a hash table provide O(1) lookup time?"      |
| **application** | Real-world scenarios, problem-solving      | Applying knowledge to new situations | "How would you design a caching system for a web API?" |
| **mixed**       | Combination of all types (default)         | Comprehensive review                 | Varies                                                 |

**Usage:**

```
# Test factual recall
neuron review --question-type factual

# Test conceptual understanding
neuron self-test "algorithms" --question-type conceptual

# Test application skills
neuron workflow "python" --phase verification --question-type application

# Mixed questions (default)
neuron review
```

---

## Features

- **AI-Powered Q&A:** Uses a local LLM (via Ollama) to dynamically generate deep, thoughtful questions and pedagogically sound answers.
- **Three-Phase Learning Framework:** Guides you from foundational learning through verification to AI-assisted extension.
- **Multiple Question Types:** Factual, conceptual, application, and mixed questions targeting different cognitive levels.
- **Spaced Repetition System:** Automatically schedules future reviews based on your performance.
- **Markdown Zettelkasten Support:** Imports and syncs with your folder of Markdown notes.
- **Multiple Study Modes:** Use `review`, `mix`, `teach`, `deep-dive`, `self-test`, `reflect`, or `workflow` for different learning goals.
- **Beautiful Terminal UI:** Renders Markdown beautifully in the terminal with colors and formatting.
- **Clean, Centralized Data Storage:** Manages its database in your system's standard user config directory, keeping your project folders clean.
- **Local First:** Your notes and the AI run entirely on your local machine. No cloud services or subscriptions needed.
- **Interactive Learning Tools:** Built-in commands during study sessions let you view notes, ask for explanations, and get help anytime.
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

### Step 2: Choose Your Learning Path

#### Option A: Guided Three-Phase Learning (Recommended for New Topics)

Use the `workflow` command for a structured learning experience:

```bash
# Phase 1: Build foundational knowledge
neuron workflow "python basics" --phase foundational

# Phase 2: Test and verify understanding
neuron workflow "python basics" --phase verification

# Phase 3: Extend with AI assistance
neuron workflow "python basics" --phase extension
```

**Phase Options:**

- `foundational` - Review concepts, test recall, understanding, and application
- `verification` - Self-test and reflection to identify gaps
- `extension` - Collaborative exploration, edge cases, optimization

**Interactive options in each phase:**

- Review basic concepts
- Test with different question types
- Show full note content
- Access help anytime
- Exit when ready

#### Option B: Traditional Study Modes

##### Daily Review (Spaced Repetition)

```bash
# Standard review - shows due notes with optional full note display
neuron review

# Review with specific question type
neuron review --question-type conceptual

# Brief mode - skip the full note display for faster reviews
neuron review --brief

# Review any random note, even if not due
neuron review --any
```

##### Interleaved Practice

```bash
# Standard interleaved review
neuron mix

# Interleaved with specific question type
neuron mix --question-type application

# Brief mode for faster interleaved sessions
neuron mix --brief
```

##### Test Your Knowledge

```bash
# Self-test: answer before seeing AI response
neuron self-test "data structures" --question-type conceptual
```

**Interactive Commands Available:**

- `help` or `?` - Show available commands
- `note` or `show note` - Display the full note content
- `skip` - Skip current question
- `quit` or `exit` - End the session

##### Challenge Your Understanding

```bash
# Reflection mode with Socratic questioning
neuron reflect "algorithms"
```

The AI will challenge your assumptions and explore edge cases using the "Red Team Pattern."

##### Deepen Understanding (Feynman Technique)

```bash
# Use a unique keyword from the note's title
neuron teach "parquet"
```

**Interactive Commands Available:**

- `help` or `?` - Show available commands
- `note` or `show note` - Display the full note content
- `explain <topic>` - Ask the AI to explain a specific concept
- `quit` or `exit` - End the session

##### Explore Connections (Elaboration)

```bash
neuron deep-dive "security"
```

**Interactive Commands Available:**

- `help` or `?` - Show available commands
- `note` or `show note` - Display the full note content
- `explain <topic>` - Ask the AI to explain a specific concept
- `quit` or `exit` - End the session

---

## Learning Science Behind Neuron CLI

### Better Questions, Better Learning

Neuron CLI generates questions that require **application** and **analysis**, not just memorization. Questions are designed to make you think about:

- How concepts apply to real scenarios
- Relationships between ideas
- Why things work the way they do
- What would happen if conditions changed

### Four Cognitive Levels

1. **Factual** - Tests recall of definitions and facts (foundational)
2. **Conceptual** - Tests understanding of relationships and principles (deeper)
3. **Application** - Tests ability to apply knowledge to new scenarios (deepest)
4. **Mixed** - Combines all levels for comprehensive review

### Pedagogically Sound Answers

When you reveal answers, you get:

1. A direct, clear answer (1-2 sentences)
2. Explanation of the "why" or "how"
3. Concrete examples or analogies
4. Connection to broader principles

This structure helps you build deep, interconnected knowledge instead of isolated facts.

### Three-Phase Framework

1. **Phase 1 (Foundational)** - Build baseline competence without AI dependency
2. **Phase 2 (Verification)** - Use AI to test understanding and identify gaps
3. **Phase 3 (Extension)** - Leverage AI to accelerate learning while maintaining competence

This framework prevents AI dependency while maximizing AI benefits.

### Optional Context

You control when to see the full note, preventing the passive re-reading trap while still having context available when you need it.

---

## 💡 Usage Examples

### Daily Learning Routine

```bash
# Morning: Review due notes with mixed questions
neuron review

# Afternoon: Quick interleaved review
neuron mix --brief

# Evening: Deep dive on challenging topic
neuron reflect "algorithms"
```

### Learning a New Programming Language

```bash
# Week 1: Build foundations
neuron workflow "python basics" --phase foundational

# Week 2: Test understanding
neuron workflow "python basics" --phase verification
neuron self-test "python basics" --question-type application

# Week 3: Explore advanced topics
neuron workflow "python advanced" --phase extension
```

### Preparing for an Interview

```bash
# Test factual knowledge
neuron review --question-type factual

# Test conceptual understanding
neuron self-test "data structures" --question-type conceptual

# Practice application scenarios
neuron workflow "algorithms" --phase verification --question-type application
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

## Tips for Effective Learning

1. **Start with the workflow command:** Use `neuron workflow` for structured learning on new topics—it guides you through all three phases.

2. **Match question types to your goals:**
   - Learning basics? Use `factual`
   - Understanding concepts? Use `conceptual`
   - Preparing for real work? Use `application`
   - Comprehensive review? Use `mixed` (default)

3. **Review daily:** Even 5-10 minutes of daily review is more effective than long, infrequent sessions.

4. **Use brief mode for efficiency:** When you're confident and want to move quickly, use `--brief` to skip optional note displays.

5. **Progress through all three phases:** Don't skip Phase 2 (Verification)—it's where you identify and fix knowledge gaps.

6. **Self-test regularly:** Use `neuron self-test` to practice active recall before seeing answers.

7. **Reflect when stuck:** Use `neuron reflect` to challenge your assumptions and explore edge cases.

8. **Mix it up:** Use `neuron mix` regularly to practice context switching between topics.

9. **Teach what you learn:** Use `neuron teach` on topics you think you understand well—it will reveal gaps you didn't know existed.

10. **Go deep when needed:** If you're struggling with a concept, use `neuron deep-dive` to explore it from different angles.

11. **Trust the system:** If a note feels too easy, mark it "Easy." The SRS algorithm will adjust the interval appropriately.

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
- Mobile app for on-the-go reviews

---

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/soyomarvaldezg/neuron-cli/issues) page
2. Open a new issue with detailed information about your problem
3. Include your OS, Go version, and any error messages

---

Happy learning! 🧠
