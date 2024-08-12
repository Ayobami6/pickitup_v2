package main

type UserRepo interface {
	CreateUser(user *User) error
	GetUserByID(id uint) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
}