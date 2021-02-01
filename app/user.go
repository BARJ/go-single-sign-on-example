package main

import "errors"

type User struct {
	ID      int    `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type UserManager struct {
	repository Repository
}

func NewUserManager(repository Repository) UserManager {
	return UserManager{repository: repository}
}

func (um UserManager) GetUserByID(id int) (User, error) {
	if id < 1 {
		return User{}, errors.New("invalid user id")
	}
	return um.repository.GetUserByID(id)
}

type ErrUserNotFound string

func (e ErrUserNotFound) Error() string {
	return string(e)
}
