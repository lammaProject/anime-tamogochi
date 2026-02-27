package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"anime-tamogochi/backend/internal/db"
	"anime-tamogochi/backend/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "server error", http.StatusInternalServerError)
		return
	}

	var id int
	err = db.DB.QueryRow(
		`INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`,
		req.Username, string(hash),
	).Scan(&id)
	if err != nil {
		jsonError(w, "username already taken", http.StatusConflict)
		return
	}

	token := newToken()
	db.DB.Exec(`INSERT INTO sessions (token, user_id) VALUES ($1, $2)`, token, id)
	setSessionCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"id": id, "username": req.Username})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	var id int
	var hash string
	err := db.DB.QueryRow(
		`SELECT id, password_hash FROM users WHERE username = $1`, req.Username,
	).Scan(&id, &hash)
	if err != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token := newToken()
	db.DB.Exec(`INSERT INTO sessions (token, user_id) VALUES ($1, $2)`, token, id)
	setSessionCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"id": id, "username": req.Username})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		db.DB.Exec(`DELETE FROM sessions WHERE token = $1`, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		MaxAge:  -1,
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
}

func Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	var username string
	db.DB.QueryRow(`SELECT username FROM users WHERE id = $1`, userID).Scan(&username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"id": userID, "username": username})
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 30, // 30 дней
	})
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
