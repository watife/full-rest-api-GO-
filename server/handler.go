package main

import (
	"encoding/json"
	"fakorede-bolu/full-rest-api/server/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user models.User

	err := decoder.Decode(&user)

	if ok, errors := app.validateInputs(user); !ok {
		app.validationError(w, http.StatusUnprocessableEntity, errors)
		return
	}

	pass := app.hashPassword(w, user.Password)

	user.Password = string(pass)

	u, err := app.user.Register(user.Email, user.Password)

	if err != nil {
		app.respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	json.NewEncoder(w).Encode(u)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var user models.User

	err := decoder.Decode(&user)

	if ok, errors := app.validateInputs(user); !ok {
		app.validationError(w, http.StatusUnprocessableEntity, errors)
		return
	}

	u, err := app.user.Login(user.Email, user.Password)

	if err != nil {
		app.respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	valideToken, err := app.generateJWT(u.ID)

	tk := &models.Token{
		UserID: u.ID,
		Email:  u.Email,
		Token:  valideToken,
	}

	json.NewEncoder(w).Encode(tk)
}

func (app *application) forgetPassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.respondError(w, http.StatusNotFound, "The specified User not found")
		return
	}

	decoder := json.NewDecoder(r.Body)

	var pass models.ForgetPassword

	err = decoder.Decode(&pass)

	if ok, errors := app.validateInputs(pass); !ok {
		app.validationError(w, http.StatusUnprocessableEntity, errors)
		return
	}

	newPasswordHash := app.hashPassword(w, pass.NewPassword)

	pass.NewPassword = string(newPasswordHash)

	resp, err := app.user.Update(id, pass.OldPassword, pass.NewPassword)

	if err != nil {
		app.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	app.respondJSON(w, http.StatusOK, resp)
	return
}
