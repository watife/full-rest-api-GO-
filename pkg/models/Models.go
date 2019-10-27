package models

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ErrNoRecord
var (
	ErrNoRecord           = errors.New("no matching record found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("This email already exist")
	ErrEmailQueue         = errors.New("This email is not entered for future queue")
	ErrTransaction        = errors.New("Could not start transaction Error")
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
	Role     string `validate:"required"`
}

// UserLogin Struct
type UserLogin struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

//Token struct declaration
type Token struct {
	UserID     int    `json:"userId"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Authorized bool   `json:"Authorized"`
	*jwt.StandardClaims
}

// ForgetPassword struct
type ForgetPassword struct {
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required"`
}

// Inbox Struct
type Inbox struct {
	Email  string `validate:"required"`
	Send   int    `validate:"required"`
	UserID int    `json:"userID"`
}
