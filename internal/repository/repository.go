package repository

import (
	"database/sql"
	"fmt"

	"github.com/yervsil/auth_service/domain"
)

var (
	usersTable = "users"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func(r *Repository) CreateUser(name, email, passwordHash string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, email, password_hash) values ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, name, email, passwordHash)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func(r *Repository) UserByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id, email, username, password_hash FROM %s WHERE email = $1", usersTable)

	row := r.db.QueryRow(query, email)
	if err := row.Scan(&user.Id, &user.Email, &user.Name, &user.Password_hash); err != nil {
		return nil, err
	}

	return &user, nil
}

func(r *Repository) UserById(id int) (*domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id, email, password_hash FROM %s WHERE id = $1", usersTable)

	row := r.db.QueryRow(query, id)
	if err := row.Scan(&user.Id, &user.Email, &user.Password_hash); err != nil {
		return nil, err
	}

	return &user, nil
}