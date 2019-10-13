package postgres

import (
	"database/sql"
	"errors"
	"fakorede-bolu/full-rest-api/server/pkg/models"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// UserModel struct
type UserModel struct {
	DB *sql.DB
}

// Register : Create/save a new User.
//  Method: POST
func (m *UserModel) Register(email, password string) (*models.User, error) {
	stmt := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email;`

	u := &models.User{}

	result, err := m.DB.Query(stmt, email, password)

	if err != nil {
		var pgError *pq.Error

		if errors.As(err, &pgError) {
			if pgError.Code.Name() == "unique_violation" && strings.Contains(pgError.Message, "users_email_key") {
				return nil, models.ErrDuplicateEmail
			}
		}

		return nil, err
	}

	defer result.Close()

	for result.Next() {
		if err := result.Scan(&u.ID, &u.Email); err != nil {
			return nil, err
		}
	}

	return u, nil
}

// Login : Get a Registered User.
//  Method: POST
func (m *UserModel) Login(email, password string) (*models.User, error) {
	stmt := `SELECT * FROM users WHERE email = $1`

	u := &models.User{}

	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&u.ID, &u.Email, &u.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err

	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err

	}

	return u, nil
}

// Update : Update a single User.
//  Method: PUT/PATCH
//  Params: True (id)
func (m *UserModel) Update(title string) (*models.User, error) {
	fmt.Println("not implemented")

	return nil, nil
}

// Destroy : Remove a single todo.
//  Method: PUT/PATCH
//  Params: True (id)
func (m *UserModel) Destroy(title string) (string, error) {
	fmt.Println("not implemented")

	return "", nil
}
