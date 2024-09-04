package services

import (
	"database/sql"
	"sport_finder/pkg/models"
)

type UserService interface {
	Registration(user *models.User, db *sql.DB) (*models.User, error)
	Authorisation(username string, password string, db *sql.DB) (bool, string, error)
}
