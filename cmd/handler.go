package main

import (
	"encoding/json"
	utils "fakorede-bolu/full-rest-api/pkg/Utils"
	"fakorede-bolu/full-rest-api/pkg/models"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (app *application) register(w http.ResponseWriter, r *http.Request) {

	var user models.User

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user)

	if ok, errors := app.validateInputs(user); !ok {
		app.validationError(w, http.StatusUnprocessableEntity, errors)
		return
	}

	pass := app.hashPassword(w, user.Password)

	user.Password = string(pass)

	timein := time.Now().Add(time.Hour*5 + time.Minute + time.Second)

	t := timein.Hour()

	// i, err := app.inbox.Inbox(user.Email, t)

	if err != nil {
		app.respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	u, err := app.user.Register(user.Email, user.Password, user.Role, t)

	if err != nil {
		app.respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	json.NewEncoder(w).Encode(u)

}

// func (app *application) confirmRegister(w http.ResponseWriter, r *http.Request) {
// 	var user models.User

// 	decoder := json.NewDecoder(r.Body)

// 	err := decoder.Decode(&user)

// 	if ok, errors := app.validateInputs(user); !ok {
// 		app.validationError(w, http.StatusUnprocessableEntity, errors)
// 		return
// 	}

// 	pass := app.hashPassword(w, user.Password)

// 	user.Password = string(pass)

// 	u, err := app.user.Register(user.Email, user.Password, user.Role, t)

// 	if err != nil {
// 		app.respondError(w, http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	json.NewEncoder(w).Encode(u)
// }

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	var user models.UserLogin

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user)

	if ok, errors := app.validateInputs(user); !ok {
		app.validationError(w, http.StatusUnprocessableEntity, errors)
		return
	}

	u, err := app.findUser(user.Email, user.Password)

	if err != nil {
		app.respondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	json.NewEncoder(w).Encode(u)
}

func (app *application) findUser(email, password string) (map[string]interface{}, error) {

	resp := make(map[string]interface{})

	u, err := app.user.Login(email, password)

	if err != nil {
		fmt.Println(1, err)
		return nil, err
	}

	valideToken, err := app.generateJWT(u.ID, u.Role)

	if err != nil {
		fmt.Println(2, err)
		return nil, err
	}

	u = &models.User{
		ID:    u.ID,
		Email: u.Email,
		Role:  u.Role,
	}

	resp["user"] = u
	resp["token"] = valideToken

	return resp, nil
}

func (app *application) forgetPassword(w http.ResponseWriter, r *http.Request) {
	var pass models.ForgetPassword

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.respondError(w, http.StatusNotFound, "The specified User not found")
		return
	}

	decoder := json.NewDecoder(r.Body)

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

// send email to each user valid
func (app *application) outBox() {
	fmt.Println("I WAS CALLED")
	// timein := time.Now().Add(time.Hour*5 + time.Minute + time.Second)

	t := time.Now().Hour()

	inboxes, err := app.inbox.Outbox(t)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i := range inboxes {
		e := &utils.Email{
			URL:     "register.html",
			Name:    i.Email,
			Email:   "Please confirm your email",
			Subject: "Our register subject",
			ID:      i.UserID,
		}

		sent, err := app.sendEmail(e)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(sent)
	}
}
