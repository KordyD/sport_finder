package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sport_finder/internal/services"
	"sport_finder/pkg/custom_errors"
	"sport_finder/pkg/models"
	"time"
)

func RegistrationController(w http.ResponseWriter, r *http.Request, db *sql.DB, userService services.UserService) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Failed to parse request body:", err)
		return
	}

	registeredUser, err := userService.Registration(&user, db)
	if errors.Is(err, custom_errors.ErrEmptyPassword) {
		http.Error(w, "Empty password", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Failed to register user:", err)
		return
	}

	registeredUserWithoutPassword := models.User{
		Username:      registeredUser.Username,
		FavoriteSport: registeredUser.FavoriteSport,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(registeredUserWithoutPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Failed to encode response:", err)
		return
	}
}

func AuthorisationController(w http.ResponseWriter, r *http.Request, db *sql.DB, userService services.UserService) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Failed to parse request body:", err)
		return
	}
	status, token, err := userService.Authorisation(user.Username, user.Password, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Failed to authorise user:", err)
		return
	}
	if status {
		w.WriteHeader(http.StatusOK)
		// Add token to cookies response
		http.SetCookie(w, &http.Cookie{
			Name:    "Authorization",
			Value:   fmt.Sprintf("Bearer %s", token),
			Expires: time.Now().Add(24 * time.Hour),
			// TODO add https
			// Secure:     false,
			HttpOnly: true,
		})
		_, err = w.Write([]byte("Authorisation successful"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Failed to write response:", err)
			return
		}
		return
	} else {
		http.Error(w, "Authorisation failed", http.StatusUnauthorized)
		return
	}
}
