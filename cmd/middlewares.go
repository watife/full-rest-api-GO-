package main

import (
	"context"
	"fakorede-bolu/full-rest-api/pkg/models"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// secure header
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)

	})
}

// log resquest
func (app *application) logRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)

	}

	return http.HandlerFunc(fn)
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Authenticate
type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

func (app *application) verifyUserToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header

		accessJWT := headers.Get("Cf-Access-Jwt-Assertion")

		accessJWT = strings.TrimSpace(accessJWT)

		if accessJWT == "" {
			app.respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tk := &models.Token{}

		token, err := jwt.Parse(accessJWT, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil {
			app.respondError(w, http.StatusForbidden, err.Error())
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if ok && token.Valid {
			role, err := claims["role"].(string)

			if err == false || role != "user" {
				app.respondError(w, http.StatusForbidden, "Unauthorized")
				return
			}
		}

		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (app *application) verifyAdminToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header

		accessJWT := headers.Get("Cf-Access-Jwt-Assertion")

		if accessJWT == "" {
			app.respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
	}
	return http.HandlerFunc(fn)
}
