// Package study contains logic related to the learning process, like SRS and LLM interaction.
package study

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/soyomarvaldezg/neuron-cli/internal/note"
)

// OllamaRequest represents the JSON payload for the Ollama /api/generate endpoint.
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse is not exported, so it doesn't need a special comment format.
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// OllamaMessage represents a single message in a chat conversation.
type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaChatRequest represents the JSON payload for the Ollama /api/chat endpoint.
type OllamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []OllamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

// OllamaChatResponse is not exported.
type OllamaChatResponse struct {
	Message OllamaMessage `json:"message"`
	Done    bool          `json:"done"`
}

// GenerateQuestion asks the LLM to generate a review question based on a note's content.
func GenerateQuestion(n *note.Note) (string, error) {
	promptContent := extractSummary(n.Content)
	prompt := fmt.Sprintf(`You are a helpful study assistant. Your goal is to help me review my notes using active recall. Based ONLY on the following text from my notes, generate one single, clear, and concise question that tests the main concept. RULES: 1. Ask only ONE question. 2. Do NOT provide the answer in your question. 3. The question must be answerable from the provided text. 4. Do not add any introductory text like "Here is a question for you:". Just provide the question. TEXT: --- %s ---`, promptContent)
	payload := OllamaRequest{Model: "llama3:8b-instruct-q4_K_M", Prompt: prompt, Stream: false}
	return sendOllamaRequest(payload)
}

// GenerateAnswer asks the LLM to provide a concise answer to a specific question.
func GenerateAnswer(question string, n *note.Note) (string, error) {
	promptContent := extractSummary(n.Content)
	prompt := fmt.Sprintf(`You are a helpful study assistant. Your goal is to provide a concise, direct answer to a question, using ONLY the provided text as your source of truth. RULES: 1. Base your answer STRICTLY on the text provided. Do not add outside information. 2. Answer the question directly and concisely. A few sentences or a bulleted list is best. 3. Do not add any introductory text like "The answer is:". Just provide the answer. QUESTION: %s TEXT TO USE: --- %s ---`, question, promptContent)
	payload := OllamaRequest{Model: "llama3:8b-instruct-q4_K_M", Prompt: prompt, Stream: false}
	return sendOllamaRequest(payload)
}

// sendOllamaRequest is a private helper to reduce code duplication for the /api/generate endpoint.
func sendOllamaRequest(payload OllamaRequest) (string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to send request to ollama: %w. Is Ollama running?", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal ollama response: %w. Response was: %s", err, string(body))
	}
	return strings.TrimSpace(ollamaResp.Response), nil
}

// SendChatMessage sends a list of messages to the Ollama chat endpoint and returns the AI's response.
func SendChatMessage(messages []OllamaMessage) (OllamaMessage, error) {
	payload := OllamaChatRequest{
		Model:    "llama3:8b-instruct-q4_K_M",
		Messages: messages,
		Stream:   false,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return OllamaMessage{}, err
	}
	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return OllamaMessage{}, fmt.Errorf("failed to send chat request to ollama: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OllamaMessage{}, err
	}
	var ollamaResp OllamaChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return OllamaMessage{}, fmt.Errorf("failed to unmarshal ollama chat response: %w. Response was: %s", err, string(body))
	}
	return ollamaResp.Message, nil
}

// extractSummary is a private helper function.
func extractSummary(fullContent string) string {
	var summary, takeaways strings.Builder
	inSummary := false
	inTakeaways := false
	scanner := bufio.NewScanner(strings.NewReader(fullContent))
	for scanner.Scan() {
		line := scanner.Text()
		lowerLine := strings.ToLower(line)
		if strings.HasPrefix(lowerLine, "## summary") {
			inSummary = true
			inTakeaways = false
			continue
		}
		if strings.HasPrefix(lowerLine, "## key takeaways") {
			inSummary = false
			inTakeaways = true
			continue
		}
		if strings.HasPrefix(lowerLine, "##") {
			inSummary = false
			inTakeaways = false
		}
		if inSummary {
			if !strings.HasPrefix(lowerLine, "## summary") {
				summary.WriteString(line + "\n")
			}
		}
		if inTakeaways {
			if !strings.HasPrefix(lowerLine, "## key takeaways") {
				takeaways.WriteString(line + "\n")
			}
		}
	}
	combined := summary.String() + takeaways.String()
	if len(strings.TrimSpace(combined)) > 10 {
		return combined
	}
	return fullContent
}
