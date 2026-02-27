package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"anime-tamogochi/backend/internal/ai"
	"anime-tamogochi/backend/internal/db"
	"anime-tamogochi/backend/internal/middleware"

	"github.com/gorilla/mux"
)

type MessageResponse struct {
	Content   string `json:"content"`
	FromUser  bool   `json:"from_user"`
	CreatedAt string `json:"created_at"`
}

// GET /api/chat/{catgirlId} — история сообщений
func GetChatHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	catgirlID := mux.Vars(r)["catgirlId"]

	rows, err := db.DB.Query(
		`SELECT content, from_user, created_at FROM messages
		 WHERE user_id = $1 AND catgirl_id = $2
		 ORDER BY created_at ASC LIMIT 100`,
		userID, catgirlID,
	)
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	msgs := []MessageResponse{}
	for rows.Next() {
		var m MessageResponse
		var t time.Time
		if err := rows.Scan(&m.Content, &m.FromUser, &t); err != nil {
			continue
		}
		m.CreatedAt = t.Format(time.RFC3339)
		msgs = append(msgs, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

// POST /api/chat/{catgirlId} — отправить сообщение, получить ответ AI
func SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	catgirlID := mux.Vars(r)["catgirlId"]

	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Content == "" {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	now := time.Now()

	// Сохраняем сообщение пользователя
	db.DB.Exec(
		`INSERT INTO messages (user_id, catgirl_id, content, from_user) VALUES ($1, $2, $3, TRUE)`,
		userID, catgirlID, body.Content,
	)

	// Загружаем персонажа и историю для AI
	charInfo := loadCharacterInfo(userID, catgirlID)
	systemPrompt := ai.BuildSystemPrompt(charInfo)
	aiHistory := loadAIHistory(userID, catgirlID, 20)

	// Получаем ответ AI
	reply, err := ai.Chat(systemPrompt, aiHistory, body.Content)
	if err != nil {
		log.Printf("ai chat error: %v", err)
	}

	replyTime := time.Now()

	// Сохраняем ответ AI
	db.DB.Exec(
		`INSERT INTO messages (user_id, catgirl_id, content, from_user) VALUES ($1, $2, $3, FALSE)`,
		userID, catgirlID, reply,
	)

	// Возвращаем оба сообщения
	type sendResponse struct {
		UserMessage MessageResponse `json:"user_message"`
		BotMessage  MessageResponse `json:"bot_message"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sendResponse{
		UserMessage: MessageResponse{
			Content:   body.Content,
			FromUser:  true,
			CreatedAt: now.Format(time.RFC3339),
		},
		BotMessage: MessageResponse{
			Content:   reply,
			FromUser:  false,
			CreatedAt: replyTime.Format(time.RFC3339),
		},
	})
}

// loadCharacterInfo и loadAIHistory — общие утилиты
func loadCharacterInfo(userID int, catgirlID string) ai.CharacterInfo {
	var dataStr string
	err := db.DB.QueryRow(
		`SELECT data FROM liked_catgirls WHERE user_id = $1 AND catgirl_id = $2 LIMIT 1`,
		userID, catgirlID,
	).Scan(&dataStr)
	if err != nil || dataStr == "" {
		return ai.CharacterInfo{}
	}

	var data struct {
		Anime struct {
			Character string `json:"character"`
			Title     string `json:"title"`
		} `json:"anime"`
		Tags     []string `json:"tags"`
		Category string   `json:"category"`
	}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return ai.CharacterInfo{}
	}

	return ai.CharacterInfo{
		Name:      data.Anime.Character,
		AnimeName: data.Anime.Title,
		Tags:      data.Tags,
		Category:  data.Category,
	}
}

func loadAIHistory(userID int, catgirlID string, limit int) []ai.Message {
	rows, err := db.DB.Query(
		`SELECT content, from_user FROM (
			SELECT content, from_user, created_at
			FROM messages
			WHERE user_id = $1 AND catgirl_id = $2
			ORDER BY created_at DESC
			LIMIT $3
		) sub ORDER BY created_at ASC`,
		userID, catgirlID, limit,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var history []ai.Message
	for rows.Next() {
		var content string
		var fromUser bool
		if err := rows.Scan(&content, &fromUser); err != nil {
			continue
		}
		role := "assistant"
		if fromUser {
			role = "user"
		}
		history = append(history, ai.Message{Role: role, Content: content})
	}
	return history
}
