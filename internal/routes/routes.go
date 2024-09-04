package routes

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"net/http"
	"sport_finder/internal/controllers"
	"sport_finder/internal/services"
	"sport_finder/internal/services/map_service"
)

func RegisterRoutes(db *sql.DB, rdb *redis.Client) {
	authService := &services.AuthService{}

	http.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegistrationController(w, r, db, authService)
	})

	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		controllers.AuthorisationController(w, r, db, authService)
	})

	http.HandleFunc("GET /objects", func(w http.ResponseWriter, r *http.Request) {
		controllers.MapController(w, r, &map_service.MapService{}, rdb)
	})

}
