package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"anime-tamogochi/backend/internal/db"
	"anime-tamogochi/backend/internal/handler"
	"anime-tamogochi/backend/internal/middleware"

	"github.com/gorilla/mux"
)

// getAllowedOrigins читает ALLOWED_ORIGINS из env (через запятую)
// Если не задано — разрешает localhost для локальной разработки
func getAllowedOrigins() map[string]bool {
	raw := os.Getenv("ALLOWED_ORIGINS")
	result := map[string]bool{
		"http://localhost:5173":    true,
		"http://192.168.0.37:5173": true,
	}
	if raw != "" {
		for _, o := range strings.Split(raw, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				result[o] = true
			}
		}
	}
	return result
}

func main() {
	db.Connect()
	db.Migrate()

	r := mux.NewRouter()

	// CORS для фронта
	r.Use(corsMiddleware)

	// Auth (публичные)
	r.HandleFunc("/api/register", handler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/login", handler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/logout", handler.Logout).Methods("POST", "OPTIONS")

	// Защищённые маршруты
	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth)

	protected.HandleFunc("/api/me", handler.Me).Methods("GET", "OPTIONS")
	protected.HandleFunc("/api/liked", handler.GetLiked).Methods("GET", "OPTIONS")
	protected.HandleFunc("/api/liked", handler.AddLiked).Methods("POST", "OPTIONS")
	protected.HandleFunc("/api/liked/{catgirlId}", handler.RemoveLiked).Methods("DELETE", "OPTIONS")

	// Chat REST API
	protected.HandleFunc("/api/chat/{catgirlId}", handler.GetChatHistory).Methods("GET", "OPTIONS")
	protected.HandleFunc("/api/chat/{catgirlId}", handler.SendMessage).Methods("POST", "OPTIONS")

	log.Println("🚀 Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := getAllowedOrigins()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" {
			if allowedOrigins[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Add("Vary", "Origin")
			} else {
				http.Error(w, "CORS: origin not allowed", http.StatusForbidden)
				return
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
