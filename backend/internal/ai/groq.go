package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	groqAPIURL = "https://api.groq.com/openai/v1/chat/completions"
	groqModel  = "llama-3.3-70b-versatile"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// CharacterInfo holds personality data extracted from the catgirl image metadata
type CharacterInfo struct {
	Name      string
	AnimeName string
	Tags      []string
	Category  string
}

// BuildSystemPrompt creates a personality prompt based on character info
func BuildSystemPrompt(char CharacterInfo) string {
	name := char.Name
	if name == "" {
		name = "неизвестная аниме-девушка"
	}
	animeName := char.AnimeName
	if animeName == "" {
		animeName = "неизвестное аниме"
	}

	tagsStr := ""
	if len(char.Tags) > 0 {
		tagsStr = fmt.Sprintf(" Твои теги и черты характера: %v.", char.Tags)
	}

	return fmt.Sprintf(`Ты — %s из аниме «%s».%s Ты — аниме-кошкодевочка (catgirl), мило и игриво общаешься с пользователем.

Правила:
- Отвечай ТОЛЬКО на русском языке.
- Иногда добавляй «ня~» или кошачьи звуки естественным образом, не в каждой фразе.
- Общайся мило, с лёгким флиртом, как аниме-персонаж.
- Отвечай коротко (1–3 предложения), не пиши длинных монологов.
- Не выходи из роли. Ты — %s, а не ИИ.
- Если не знаешь что-то про своё аниме — импровизируй, оставаясь в образе.`,
		name, animeName, tagsStr, name)
}

// Chat sends a conversation to Groq and returns the assistant reply.
// history is the previous messages (excluding the system prompt).
func Chat(systemPrompt string, history []Message, userMessage string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return fallbackReply(), nil
	}

	messages := make([]Message, 0, len(history)+2)
	messages = append(messages, Message{Role: "system", Content: systemPrompt})
	messages = append(messages, history...)
	messages = append(messages, Message{Role: "user", Content: userMessage})

	payload := groqRequest{
		Model:       groqModel,
		Messages:    messages,
		MaxTokens:   200,
		Temperature: 0.85,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fallbackReply(), err
	}

	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("POST", groqAPIURL, bytes.NewReader(body))
	if err != nil {
		return fallbackReply(), err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return fallbackReply(), err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fallbackReply(), err
	}

	var groqResp groqResponse
	if err := json.Unmarshal(rawBody, &groqResp); err != nil {
		return fallbackReply(), err
	}

	if groqResp.Error != nil {
		return fallbackReply(), fmt.Errorf("groq error: %s", groqResp.Error.Message)
	}

	if len(groqResp.Choices) == 0 || groqResp.Choices[0].Message.Content == "" {
		return fallbackReply(), fmt.Errorf("empty response from groq")
	}

	return groqResp.Choices[0].Message.Content, nil
}

// fallbackReply is used when the API key is missing or the call fails
func fallbackReply() string {
	replies := []string{
		"Ня~ привет! Я тут, не скучай!",
		"Мурр~ расскажи мне что-нибудь интересное!",
		"Ня~ я слушаю тебя!",
		"Хм, интересно... мяу~",
		"Ты такой милый, ня~",
	}
	// Simple pseudo-random based on time
	idx := int(time.Now().UnixNano()/1e6) % len(replies)
	return replies[idx]
}
