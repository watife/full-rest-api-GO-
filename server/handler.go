package main

import (
	"encoding/json"
	"fakorede-bolu/full-rest-api/server/pkg/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	json.NewDecoder(r.Body).Decode(user)

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		app.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user.Password = string(pass)

	u, err := app.user.Register(user.Email, user.Password)

	if err != nil {
		app.respondError(w, http.StatusUnauthorized, err.Error())
	}

	json.NewEncoder(w).Encode(u)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	json.NewDecoder(r.Body).Decode(user)

	u, err := app.user.Login(user.Email, user.Password)

	if err != nil {
		app.respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	valideToken, err := app.GenerateJWT(u.ID)

	tk := &models.Token{
		UserID: u.ID,
		Email:  u.Email,
		Token:  valideToken,
	}

	json.NewEncoder(w).Encode(tk)
}
