package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest)

	authUserMiddleware := alice.New(app.verifyUserToken)

	mux := pat.New()

	// User
	mux.Post("/register", http.HandlerFunc(app.register))
	mux.Post("/login", http.HandlerFunc(app.login))
	mux.Put("/forgetpassword/:id", http.HandlerFunc(app.forgetPassword))
	mux.Post("/todo", authUserMiddleware.ThenFunc(http.HandlerFunc(app.forgetPassword)))

	return standardMiddleware.Then(mux)
}
