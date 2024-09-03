package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"log"
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
		log.Fatalln("Failed to begin transaction:", err)
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
		log.Fatalln("Failed to commit transaction:", err)
		return nil, err
	}

	return &models.User{Username: user.Username, Password: hashedPassword, FavoriteSport: user.FavoriteSport}, nil
}

func Authorisation(username string, password string, db *sql.DB) (bool, error) {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPassword := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	var dbPassword string
	err := db.QueryRow(`SELECT password FROM users WHERE username = $1`, username).Scan(&dbPassword)
	if err != nil {
		log.Println("Failed to get password from database:", err)
		return false, err
	}

	if hashedPassword == dbPassword {
		log.Println("Authorisation successful:", username, password)
		return true, nil
	}

	return false, nil
}
