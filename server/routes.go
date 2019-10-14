package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest)

	mux := pat.New()

	// User
	mux.Post("/register", http.HandlerFunc(app.register))
	mux.Post("/login", http.HandlerFunc(app.login))
	mux.Put("/forgetpassword/:id", http.HandlerFunc(app.forgetPassword))
	// mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	// mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	// mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	// // Serve static files
	// fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	// mux.Get("./ui/static/", http.NotFoundHandler())
	// mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
