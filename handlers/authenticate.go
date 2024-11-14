package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"task-management/models"
	"task-management/utils"
)

func AuthenticateHandler(db *sql.DB, jwtKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var login models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user models.User
		var hashedPassword string
		err := db.QueryRow("SELECT id, username, password FROM users WHERE username = $1",
			login.Username).Scan(&user.ID, &user.Username, &hashedPassword)

		if err != nil || !utils.CheckPasswordHash(login.Password, hashedPassword) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateToken(user.Username, jwtKey)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"token":    token,
			"username": user.Username,
		})
	}
}
