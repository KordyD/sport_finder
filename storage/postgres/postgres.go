package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres() *Postgres {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connection := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", dbUsername, dbPassword, dbName)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatalln("Error in connection", err)
	}
	return &Postgres{db: db}
}

func (p *Postgres) AddUser(username string, password string, favoriteSport string) (int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return 0, err
	}
	result, err := tx.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, username, password)
	// TODO add favorite sport
	if err != nil {
		tx.Rollback()
		log.Println("Failed to execute insert query:", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Println("Rows affected:", rowsAffected)

	err = tx.Commit()
	if err != nil {
		log.Println("Failed to commit transaction:", err)
		return 0, err
	}
	return rowsAffected, nil
}

func (p *Postgres) GetPasswordByUsername(username string) (string, error) {
	var dbPassword string
	err := p.db.QueryRow(`SELECT password FROM users WHERE username = $1`, username).Scan(&dbPassword)
	if err != nil {
		log.Println("Failed to get password from database:", err)
		return "", err
	}
	return dbPassword, nil
}

func (p *Postgres) UpdateToken(username string, token string) (int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return 0, err
	}
	result, err := tx.Exec(`UPDATE users SET token = $1 WHERE username = $2`, token, username)
	if err != nil {
		tx.Rollback()
		log.Println("Failed to update token in database:", err)
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Println("Rows affected:", rowsAffected)
	err = tx.Commit()
	if err != nil {
		log.Println("Failed to commit transaction:", err)
		return 0, err
	}
	return rowsAffected, nil
}
