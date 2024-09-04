package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"sport_finder/pkg/custom_errors"
	"sport_finder/pkg/models"
)

func Registration(user *models.User, db *sql.DB) (*models.User, error) {
	if user.Password == "" {
		log.Println("Empty password")
		return nil, custom_errors.ErrEmptyPassword
	}
	hash := sha256.New()
	hash.Write([]byte(user.Password))
	hashedPassword := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return nil, err
	}
	result, err := tx.Exec(`INSERT INTO users (username, password, favorite_sport) VALUES ($1, $2, $3)`, user.Username, hashedPassword, user.FavoriteSport)
	if err != nil {
		tx.Rollback()
		log.Println("Failed to execute insert query:", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	log.Println("Rows affected:", rowsAffected)

	err = tx.Commit()
	if err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &models.User{Username: user.Username, Password: hashedPassword, FavoriteSport: user.FavoriteSport}, nil
}

func Authorisation(username string, password string, db *sql.DB) (bool, string, error) {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPassword := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	var dbPassword string
	err := db.QueryRow(`SELECT password FROM users WHERE username = $1`, username).Scan(&dbPassword)
	if err != nil {
		log.Println("Failed to get password from database:", err)
		return false, "", err
	}

	if hashedPassword == dbPassword {
		log.Println("Authorisation successful:", username, password)
		key := []byte(os.Getenv("JWT_SECRET"))
		token := jwt.New(jwt.SigningMethodHS256)
		signedToken, err := token.SignedString(key)
		if err != nil {
			log.Println("Failed to sign token:", err)
			return false, "", err
		}
		_, err = db.Exec(`UPDATE users SET token = $1 WHERE username = $2`, signedToken, username)
		if err != nil {
			log.Println("Failed to update token in database:", err)
			return false, "", err
		}
		return true, signedToken, nil
	}

	return false, "", nil
}
