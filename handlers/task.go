package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"task-management/models"
)

func TaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/task/"):]
		switch r.Method {
		case "GET":
			rows, err := db.Query("SELECT id, title, description, status FROM tasks")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var tasks []models.Task
			for rows.Next() {
				var task models.Task
				if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				tasks = append(tasks, task)
			}
			json.NewEncoder(w).Encode(tasks)

		case "POST":
			var task models.Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err := db.QueryRow(
				"INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id",
				task.Title, task.Description, task.Status,
			).Scan(&task.ID)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(task)

		case "PUT":
			var task models.Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			result, err := db.Exec(
				"UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4",
				task.Title, task.Description, task.Status, id,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rows, _ := result.RowsAffected()
			if rows == 0 {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)

		case "DELETE":
			result, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rows, _ := result.RowsAffected()
			if rows == 0 {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
