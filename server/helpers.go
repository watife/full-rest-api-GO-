package main

import (
	"bytes"
	"encoding/json"
	utils "fakorede-bolu/full-rest-api/server/pkg/Utils"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
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

// validate Error
func (app *application) validationError(w http.ResponseWriter, code int, payload interface{}) {
	app.respondJSON(w, code, payload)
	return
}

// Hash Passwords
func (app *application) hashPassword(w http.ResponseWriter, userPassword string) []byte {
	pass, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		app.respondError(w, http.StatusBadRequest, err.Error())
		return nil
	}

	return pass
}

// Generate JWT
func (app *application) generateJWT(userID int) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	// claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	return token.SignedString([]byte(jwtKey))
}

/*
 * --------------------------Validators----------------------------------------
 *
**/

func (app *application) validateInputs(dataSet interface{}) (bool, map[string]string) {
	err := validate.Struct(dataSet)

	if err != nil {

		//Validation syntax is invalid
		if err, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}

		//Validation errors occurred
		errors := make(map[string]string)
		//Use reflector to reverse engineer struct
		reflected := reflect.ValueOf(dataSet)

		for _, err := range err.(validator.ValidationErrors) {

			// Attempt to find field by name and get json tag name
			field, _ := reflected.Type().FieldByName(err.StructField())
			var name string

			//If json tag doesn't exist, use lower case of name
			if name = field.Tag.Get("json"); name == "" {
				name = strings.ToLower(err.StructField())
			}

			switch err.Tag() {
			case "required":
				errors[name] = "The " + name + " is required"
				break
			case "email":
				errors[name] = "The " + name + " should be a valid email"
				break
			case "eqfield":
				errors[name] = "The " + name + " should be equal to the " + err.Param()
				break
			default:
				errors[name] = "The " + name + " is invalid"
				break
			}
		}

		return false, errors
	}
	return true, nil
}

/*
 * ---------------------------Send Email----------------------------
 *
**/

// EmailUser struct
type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

// SMTPTemplateData struct
type SMTPTemplateData struct {
	Name  string
	Email string
	ID    int
}

func (app *application) sendEmail(m *utils.Email) (bool, error) {
	gmailUsername := os.Getenv("GMAIL_USERNAME")
	gmailPassword := os.Getenv("GMAIL_PASSWORD")
	gmailServer := "smtp.gmail.com"

	var doc bytes.Buffer

	smtpData := &SMTPTemplateData{
		Name:  gmailUsername,
		Email: m.Name,
		ID:    m.ID,
	}

	emailUser := &EmailUser{gmailUsername, gmailPassword, gmailServer, 587}

	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer,
	)

	t, err := template.ParseFiles("././templates/" + m.URL)

	if err != nil {
		return false, err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	doc.Write([]byte(fmt.Sprintf("Subject:"+m.Subject+"\n%s\n\n", mime)))

	err = t.Execute(&doc, smtpData)

	if err != nil {
		return false, err
	}

	err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
		auth,
		emailUser.Username,
		[]string{m.Name},
		doc.Bytes())
	if err != nil {
		return false, err
	}

	return true, nil

}
