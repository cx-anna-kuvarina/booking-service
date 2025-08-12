package users

import "context"

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type Store interface {
	CreateUser(context.Context, *User) error
	GetByID(context.Context, string) (*User, error)
	GetUserIdByEmail(context.Context, string) (string, error)
	GetUser(context.Context, string) (*User, error)
	UpdateUser(context.Context, *User) error
	DeleteUser(context.Context, string) error
}
