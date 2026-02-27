package middleware

import (
	"context"
	"net/http"

	"anime-tamogochi/backend/internal/db"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// Auth проверяет куку session_token, кладёт user_id в контекст
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var userID int
		err = db.DB.QueryRow(
			`SELECT user_id FROM sessions WHERE token = $1`, cookie.Value,
		).Scan(&userID)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) int {
	v, _ := r.Context().Value(UserIDKey).(int)
	return v
}
