package main

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	GetUserByID(id int) (User, error)
	GetUserByEmail(email string) (User, error)
	CreateUser(user User) (User, error)
}

var _ Repository = (*SqlRepository)(nil)

type SqlRepository struct {
	db *sql.DB
}

func NewSqlRepository(db *sql.DB) SqlRepository {
	return SqlRepository{db: db}
}

func (r SqlRepository) GetUserByID(id int) (User, error) {
	query := `SELECT "id", "email", "name", "picture" FROM "user" WHERE "id" = $1;`
	row := r.db.QueryRow(query, id)

	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Picture)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound(fmt.Sprintf("user with id %d not found", id))
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r SqlRepository) GetUserByEmail(email string) (User, error) {
	query := `SELECT "id", "email", "name", "picture" FROM "user" WHERE "email" = $1;`
	row := r.db.QueryRow(query, email)

	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Picture)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound(fmt.Sprintf("user with email %s not found", email))
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r SqlRepository) CreateUser(user User) (User, error) {
	query := `
		INSERT INTO "user" ("email", "name", "picture")
		VALUES ($1, $2, $3)
		RETURNING "id", "email", "name", "picture";
	`
	row := r.db.QueryRow(query, user.Email, user.Name, user.Picture)

	user = User{}
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Picture)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
