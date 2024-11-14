package middleware

import (
	"net/http"
	"strings"
	"task-management/utils"
)

func Authenticate(next http.HandlerFunc, jwtKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ValidateToken(tokenString, jwtKey)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User", claims.Username)
		next(w, r)
	}
}
