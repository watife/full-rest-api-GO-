package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// server error
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// respond JSON
func (app *application) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
	return
}

// respond Error
func (app *application) respondError(w http.ResponseWriter, code int, message string) {
	app.respondJSON(w, code, map[string]string{"error": message})
	return
}

// Generate JWT
func (app *application) GenerateJWT(userID int) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	// claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	return token.SignedString([]byte(jwtKey))
}
