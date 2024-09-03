package routes

import (
	"database/sql"
	"net/http"
	"sport_finder/internal/controllers"
)

func RegisterRoutes(db *sql.DB) {
	http.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegistrationController(w, r, db)
	})

	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		controllers.AuthorisationController(w, r, db)
	})
}
