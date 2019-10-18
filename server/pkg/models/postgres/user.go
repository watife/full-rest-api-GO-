package postgres

import (
	"database/sql"
	"errors"
	"fakorede-bolu/full-rest-api/server/pkg/models"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// UserModel struct
type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) createTable() error {
	env := os.Getenv("APP_ENV")

	d := `DROP TABLE users;`

	c := `CREATE TABLE users (
			id serial PRIMARY KEY,
			email text UNIQUE NOT NULL,
			password varchar(100) NOT NULL
		);`

	if env == "development" {
		_, err := m.DB.Exec(d)

		_, err = m.DB.Exec(c)

		if err != nil {
			return err
		}

		return nil
	}

	_, err := m.DB.Exec(c)

	if err != nil {
		return err
	}

	log.Println("table user create successfully")
	return nil
}

// Register : Create/save a new User.
//  Method: POST
func (m *UserModel) Register(email, password string) (*models.User, error) {
	m.createTable()
	stmt := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email;`

	u := &models.User{}

	row, err := m.DB.Query(stmt, email, password)

	if err != nil {
		var pgError *pq.Error

		if errors.As(err, &pgError) {
			if pgError.Code.Name() == "unique_violation" && strings.Contains(pgError.Message, "users_email_key") {
				return nil, models.ErrDuplicateEmail
			}
		}

		return nil, err
	}

	defer row.Close()

	for row.Next() {
		if err := row.Scan(&u.ID, &u.Email); err != nil {
			return nil, err
		}
	}

	return u, nil
}

// Login : Get a Registered User.
//  Method: POST
func (m *UserModel) Login(email, password string) (*models.User, error) {

	u, err := m.GetByEmail(email, password)

	if err != nil {
		return nil, models.ErrInvalidCredentials
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

// GetByEmail  : Get a Registered User.
func (m *UserModel) GetByEmail(email, password string) (*models.User, error) {
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

	return u, nil
}

// GetByID : Get a Registered User.
func (m *UserModel) GetByID(id int) (*models.User, error) {
	stmt := `SELECT * FROM users WHERE id = $1`

	u := &models.User{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&u.ID, &u.Email, &u.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err

	}

	return u, nil
}

// Update : Update a single User.
//  Method: PUT/PATCH
//  Params: True (id)
func (m *UserModel) Update(id int, oldPassword, password string) (string, error) {
	stmt := `UPDATE users SET password = $2 WHERE id = $1;`

	// u := &models.User{}

	u, err := m.GetByID(id)

	if err != nil {
		return "", models.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", models.ErrInvalidCredentials
		}
		return "", bcrypt.ErrHashTooShort
	}

	_, err = m.DB.Exec(stmt, id, password)

	if err != nil {
		return "", models.ErrInvalidCredentials
	}

	return "success: password updated successfully", nil
}

// Destroy : Remove a single todo.
//  Method: PUT/PATCH
//  Params: True (id)
func (m *UserModel) Destroy(title string) (string, error) {
	fmt.Println("not implemented")

	return "", nil
}
