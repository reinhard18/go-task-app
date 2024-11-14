package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"task-management/models"
	"task-management/utils"
)

func RegisterHandler(db *sql.DB, jwtKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var register models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate input
		if register.Username == "" || register.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		// Check if username already exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)",
			register.Username).Scan(&exists)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

		// Hash the password
		hashedPassword, err := utils.HashPassword(register.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Insert new user
		var user models.User
		err = db.QueryRow(
			"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username",
			register.Username, hashedPassword,
		).Scan(&user.ID, &user.Username)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate token for the new user
		token, err := utils.GenerateToken(user.Username, jwtKey)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		// Return the token and username
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token":    token,
			"username": user.Username,
		})
	}
}
