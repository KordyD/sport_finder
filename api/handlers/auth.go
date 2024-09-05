package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"sport_finder/models"
	"time"
)

type Registrator interface {
	AddUser(username string, password string, favoriteSport string) (int64, error)
}
type Authenticator interface {
	GetPasswordByUsername(username string) (string, error)
	UpdateToken(username string, token string) (int64, error)
}

func NewRegHandler(registrator Registrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Failed to parse request body:", err)
			return
		}
		if user.Password == "" {
			log.Println("Empty password")
			http.Error(w, "Empty password", http.StatusBadRequest)
			return
		}
		hash := sha256.New()
		hash.Write([]byte(user.Password))
		hashedPassword := base64.StdEncoding.EncodeToString(hash.Sum(nil))
		_, err = registrator.AddUser(user.Username, hashedPassword, user.FavoriteSport)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Failed to register user:", err)
			return
		}
		registeredUserWithoutPassword := models.User{
			Username:      user.Username,
			FavoriteSport: user.FavoriteSport,
		}
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(registeredUserWithoutPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Failed to encode response:", err)
			return
		}
	}
}

func NewAuthHandler(authenticator Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Failed to parse request body:", err)
			return
		}

		hash := sha256.New()
		hash.Write([]byte(user.Password))
		hashedPassword := base64.StdEncoding.EncodeToString(hash.Sum(nil))

		dbPassword, err := authenticator.GetPasswordByUsername(user.Username)

		if hashedPassword != dbPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("Authorisation successful: ", user.Username)
		key := []byte(os.Getenv("JWT_SECRET"))
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user.Username
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		signedToken, err := token.SignedString(key)
		if err != nil {
			log.Println("Failed to sign token:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = authenticator.UpdateToken(user.Username, signedToken)
		if err != nil {
			log.Println("Failed to update token in database:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Token",
			Value:    signedToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			// TODO add secure
			// Secure: true,
		})
		w.WriteHeader(http.StatusOK)

	}
}
