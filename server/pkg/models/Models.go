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
	ID      int    `json:"id"`
	Content string `json:"content"`
	UserID  int    `json:"userId"`
	Created time.Time
	Edited  time.Time
}

// User Struct
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Token struct declaration
type Token struct {
	UserID int    `json:"userId"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}
