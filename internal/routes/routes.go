package routes

import (
	"database/sql"
	"net/http"
	"sport_finder/internal/controllers"
	"sport_finder/internal/services"
)

func RegisterRoutes(db *sql.DB) {
	authService := &services.AuthService{}

	http.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegistrationController(w, r, db, authService)
	})

	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		controllers.AuthorisationController(w, r, db, authService)
	})
}
