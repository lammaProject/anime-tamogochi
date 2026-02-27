package handler

import (
	"encoding/json"
	"net/http"

	"anime-tamogochi/backend/internal/db"
	"anime-tamogochi/backend/internal/middleware"

	"github.com/gorilla/mux"
)

// GET /api/liked — список лайкнутых для текущего пользователя
func GetLiked(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	rows, err := db.DB.Query(
		`SELECT catgirl_id, data FROM liked_catgirls WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type item struct {
		CatgirlID string          `json:"catgirl_id"`
		Data      json.RawMessage `json:"data"`
	}
	result := []item{}

	for rows.Next() {
		var it item
		var dataStr string
		if err := rows.Scan(&it.CatgirlID, &dataStr); err != nil {
			continue
		}
		it.Data = json.RawMessage(dataStr)
		result = append(result, it)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// POST /api/liked — добавить в лайки
func AddLiked(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var body struct {
		CatgirlID string          `json:"catgirl_id"`
		Data      json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.CatgirlID == "" {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(
		`INSERT INTO liked_catgirls (user_id, catgirl_id, data) VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, catgirl_id) DO NOTHING`,
		userID, body.CatgirlID, string(body.Data),
	)
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
}

// DELETE /api/liked/{catgirlId}
func RemoveLiked(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	catgirlID := mux.Vars(r)["catgirlId"]

	db.DB.Exec(
		`DELETE FROM liked_catgirls WHERE user_id = $1 AND catgirl_id = $2`,
		userID, catgirlID,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
}
