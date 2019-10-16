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
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

//Token struct declaration
type Token struct {
	UserID int    `json:"userId"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

// ForgetPassword struct
type ForgetPassword struct {
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required"`
}

// SendEmail Struct
type SendEmail struct {
}
