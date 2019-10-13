package models

import (
	"errors"
	"time"
)

// ErrNoRecord
var (
	ErrNoRecord           = errors.New("no matching record found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("This email already exist")
)

// Todo Struct
type Todo struct {
	ID      int
	Content string
	UserID  int
	Created time.Time
	Edited  time.Time
}

// User Struct
type User struct {
	ID       int
	Email    string
	Password string
}

//Token struct declaration
type Token struct {
	UserID int
	Email  string
	Token  string
}
