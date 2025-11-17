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

// QuestionType represents different types of questions that can be generated
type QuestionType string

const (
	QuestionTypeFactual     QuestionType = "factual"
	QuestionTypeConceptual  QuestionType = "conceptual"
	QuestionTypeApplication QuestionType = "application"
	QuestionTypeMixed       QuestionType = "mixed"
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

// GenerateQuestion asks the LLM to generate a review question based on a note's content and question type.
func GenerateQuestion(n *note.Note, questionType QuestionType) (string, error) {
	promptContent := extractSummary(n.Content)

	var prompt string
	switch questionType {
	case QuestionTypeFactual:
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in factual recall.
Generate ONE factual recall question that tests specific knowledge from this material.

RULES:
1. Focus on definitions, facts, dates, names, or specific details
2. Questions should have clear, objective answers
3. Use these question types:
   - "What is the definition of [concept]?"
   - "Who developed [theory/method]?"
   - "When did [event] occur?"
   - "What are the key components of [system]?"
   - "List the main features of [concept]"
4. VARY your questions - don't always ask the same thing
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent)

	case QuestionTypeConceptual:
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in conceptual understanding.
Generate ONE conceptual question that tests deep understanding of relationships and principles.

RULES:
1. Focus on "why" and "how" questions
2. Test understanding of underlying principles
3. Use these question types:
   - "Why does [concept] work the way it does?"
   - "What is the relationship between [A] and [B]?"
   - "How does [principle] explain [phenomenon]?"
   - "What are the underlying assumptions of [theory]?"
   - "Compare and contrast [concept] with [similar concept]"
4. VARY your questions - don't always ask the same thing
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent)

	case QuestionTypeApplication:
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in practical application.
Generate ONE application question that tests ability to apply concepts to real scenarios.

RULES:
1. Create scenarios that require applying the concept
2. Focus on problem-solving and implementation
3. Use these question types:
   - "How would you apply [concept] to solve [problem]?"
   - "What would happen if [condition] changed in [scenario]?"
   - "Design a solution using [principle] for [situation]"
   - "Given [context], how would you use [method]?"
   - "Create a scenario where [concept] would be useful"
4. VARY your questions - don't always ask the same thing
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent)

	case QuestionTypeMixed:
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in comprehensive understanding.
Generate ONE high-quality question that tests understanding of this material.

RULES:
1. Create questions that require thinking, not just memorization
2. Mix of factual, conceptual, and application approaches
3. Use these question types (choose the most appropriate):
   - "How would you apply [concept] to [scenario]?"
   - "What's the relationship between [A] and [B]?"
   - "Why does [X] lead to [Y]?"
   - "What would happen if [condition changed]?"
   - "Compare [concept] with [alternative]"
   - "What are the limitations of [method]?"
4. VARY your questions - don't always ask the same thing
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent)

	default:
		// Fallback to original behavior
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in active recall. Generate ONE high-quality question that tests deep understanding of this material.

RULES:
1. Create questions that require APPLICATION or ANALYSIS, not just memorization
2. Use these question types (choose the most appropriate):
   - "How would you apply [concept] to [scenario]?"
   - "What's the relationship between [A] and [B]?"
   - "Why does [X] lead to [Y]?"
   - "What would happen if [condition changed]?"
3. VARY your questions - don't always ask the same thing
4. The question should make the learner THINK, not just recall facts
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent)
	}

	payload := OllamaRequest{Model: "llama3:8b-instruct-q4_K_M", Prompt: prompt, Stream: false}
	return sendOllamaRequest(payload)
}

// GenerateAnswer asks the LLM to provide a concise answer to a specific question.
func GenerateAnswer(question string, n *note.Note) (string, error) {
	promptContent := extractSummary(n.Content)
	prompt := fmt.Sprintf(`You are a learning coach providing pedagogically effective answers.

QUESTION: %s

YOUR TASK: Provide an answer that helps deep learning:
1. Start with a direct 1-2 sentence answer
2. Then explain the "why" or "how" behind it
3. If applicable, give a concrete example or analogy
4. End with a connection to a broader principle (if relevant)

Keep it concise but insightful (3-5 sentences total).

SOURCE MATERIAL:
---
%s
---`, question, promptContent)
	payload := OllamaRequest{Model: "llama3:8b-instruct-q4_K_M", Prompt: prompt, Stream: false}
	return sendOllamaRequest(payload)
}

// CompareAnswers compares user's answer with the correct answer and provides feedback.
func CompareAnswers(userAnswer, correctAnswer, question string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert learning coach comparing a student's answer with the correct answer.

QUESTION: %s

STUDENT'S ANSWER: %s

CORRECT ANSWER: %s

YOUR TASK: Provide constructive feedback in this format:
1. ✅ What they got right (acknowledge correct parts)
2. 🤔 What they missed or misunderstood (gaps in understanding)
3. 💡 How to improve their understanding (specific suggestions)
4. 📚 Key concepts they should review (if applicable)

Be encouraging but precise. Focus on helping them understand, not just pointing out mistakes.`, question, userAnswer, correctAnswer)

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

// GenerateQuestionWithVariation generates a question with a variation hint to avoid repetition.
func GenerateQuestionWithVariation(n *note.Note, questionType QuestionType, attempt int) (string, error) {
	promptContent := extractSummary(n.Content)

	var prompt string
	switch questionType {
	case QuestionTypeFactual:
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in factual recall.
Generate ONE factual recall question that tests specific knowledge from this material.

RULES:
1. Focus on definitions, facts, dates, names, or specific details
2. Questions should have clear, objective answers
3. Use these question types:
   - "What is the definition of [concept]?"
   - "Who developed [theory/method]?"
   - "When did [event] occur?"
   - "What are the key components of [system]?"
   - "List the main features of [concept]"
4. VARY your questions - this is attempt #%d, so ask something DIFFERENT
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent, attempt)

	case QuestionTypeConceptual:
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in conceptual understanding.
Generate ONE conceptual question that tests deep understanding of relationships and principles.

RULES:
1. Focus on "why" and "how" questions
2. Test understanding of underlying principles
3. Use these question types:
   - "Why does [concept] work the way it does?"
   - "What is the relationship between [A] and [B]?"
   - "How does [principle] explain [phenomenon]?"
   - "What are the underlying assumptions of [theory]?"
   - "Compare and contrast [concept] with [similar concept]"
4. VARY your questions - this is attempt #%d, so ask something DIFFERENT
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent, attempt)

	case QuestionTypeApplication:
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in practical application.
Generate ONE application question that tests ability to apply concepts to real scenarios.

RULES:
1. Create scenarios that require applying the concept
2. Focus on problem-solving and implementation
3. Use these question types:
   - "How would you apply [concept] to solve [problem]?"
   - "What would happen if [condition] changed in [scenario]?"
   - "Design a solution using [principle] for [situation]"
   - "Given [context], how would you use [method]?"
   - "Create a scenario where [concept] would be useful"
4. VARY your questions - this is attempt #%d, so ask something DIFFERENT
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent, attempt)

	case QuestionTypeMixed:
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in comprehensive understanding.
Generate ONE high-quality question that tests understanding of this material.

RULES:
1. Create questions that require thinking, not just memorization
2. Mix of factual, conceptual, and application approaches
3. Use these question types (choose the most appropriate):
   - "How would you apply [concept] to [scenario]?"
   - "What's the relationship between [A] and [B]?"
   - "Why does [X] lead to [Y]?"
   - "What would happen if [condition changed]?"
   - "Compare [concept] with [alternative]"
   - "What are the limitations of [method]?"
4. VARY your questions - this is attempt #%d, so ask something DIFFERENT
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent, attempt)

	default:
		// Fallback to original behavior
		prompt = fmt.Sprintf(`You are an expert learning coach specializing in active recall. Generate ONE high-quality question that tests deep understanding of this material.

RULES:
1. Create questions that require APPLICATION or ANALYSIS, not just memorization
2. Use these question types (choose the most appropriate):
   - "How would you apply [concept] to [scenario]?"
   - "What's the relationship between [A] and [B]?"
   - "Why does [X] lead to [Y]?"
   - "What would happen if [condition changed]?"
3. VARY your questions - this is attempt #%d, so ask something DIFFERENT
4. The question should make the learner THINK, not just recall facts
5. Output ONLY the question, no preamble

MATERIAL:
---
%s
---`, promptContent, attempt)
	}

	payload := OllamaRequest{Model: "llama3:8b-instruct-q4_K_M", Prompt: prompt, Stream: false}
	return sendOllamaRequest(payload)
}
