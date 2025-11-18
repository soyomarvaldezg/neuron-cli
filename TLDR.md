# Neuron CLI - TL;DR Quick Start

## 🚀 5-Minute Setup

```bash
# 1. Install prerequisites
brew install sqlite3 go  # macOS
# Or: apt-get install sqlite3 golang  # Linux

# 2. Install Ollama & pull model
curl -fsSL https://ollama.com/install.sh | sh
ollama pull llama3:8b-instruct-q4_K_M
ollama serve  # Run in background

# 3. Install Neuron CLI
go install github.com/soyomarvaldezg/neuron-cli@latest
export PATH=$PATH:$(go env GOPATH)/bin

# 4. Import your notes
neuron import /path/to/your/notes

# 5. Start learning!
neuron review
```

## 📚 Daily Commands

```bash
# Quick daily review (5-10 mins)
neuron review

# Fast review without full notes
neuron review --brief

# Mix different topics
neuron mix

# Explain concepts (Feynman Technique)
neuron teach "topic"

# Deep exploration with Socratic questions
neuron deep-dive "topic"

# Test yourself before seeing answers
neuron self-test "topic"

# Challenge your assumptions
neuron reflect "topic"

# Structured 3-phase learning
neuron workflow "topic" --phase foundational
neuron workflow "topic" --phase verification
neuron workflow "topic" --phase extension
```

**📖 3-Phase TL;DR:** Learn basics → Test understanding → Use AI to expand

## 🎯 Question Types

```bash
# Basic facts
neuron review --question-type factual

# Deep understanding
neuron review --question-type conceptual

# Real-world application
neuron review --question-type application

# Mixed (default)
neuron review
```

## 💡 Pro Tips

- **Explain to learn?** Use `neuron teach "topic"`
- **Deep dive?** Use `neuron deep-dive "topic"`
- **New topic?** Start with `neuron workflow`
- **Daily habit?** Use `neuron review --brief`
- **Stuck?** Try `neuron reflect "topic"`
- **Interview prep?** Use `--question-type application`
- **Quick session?** Add `--brief` to any command

## 🎮 Interactive Commands

During `teach`, `deep-dive`, `self-test`, and `reflect` sessions:

- `help` or `?` - Show commands
- `note` - View full note
- `quit` - Exit session
- `explain <topic>` - Ask AI to explain (teach/deep-dive only)
- `skip` - Skip question (self-test only)

That's it! You're ready to learn effectively. 🧠
