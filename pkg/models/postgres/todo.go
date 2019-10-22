package postgres

import (
	"database/sql"
	"errors"
	"fakorede-bolu/full-rest-api/pkg/models"
	"fmt"
)

// TodoModel Struct
type TodoModel struct {
	DB *sql.DB
}

// Index : Show a list of all todos.
//  Method: GET
func (m *TodoModel) Index(userID int) ([]*models.Todo, error) {
	stmt := `SELECT * FROM todos WHERE userId = ?`

	rows, err := m.DB.Query(stmt, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []*models.Todo{}

	for rows.Next() {
		t := &models.Todo{}

		err = rows.Scan(&t.ID, &t.Content, &t.Created, &t.Edited)

		if err != nil {
			return nil, err
		}

		todos = append(todos, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil

}

// Create : Create/save a new todo.
//  Method: POST
func (m *TodoModel) Create(content string) (*models.Todo, error) {
	stmt := `INSERT INTO todos (content, created, edited)
			VALUES(?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	_, err := m.DB.Exec(stmt, content)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Show : Display a single todo.
//  Method: GET
//  Params: True (id)
func (m *TodoModel) Show(userID, id int) (*models.Todo, error) {
	stmt := `SELECT id, title, content, created, edited FROM snippets WHERE userId = ? AND id = ?`

	t := &models.Todo{}

	err := m.DB.QueryRow(stmt, id).Scan(&t.ID, &t.Content, &t.Created, &t.Edited)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return t, nil
}

// Update : Update a single todo.
//  Method: PUT/PATCH
//  Params: True (id)
func (m *TodoModel) Update(title string) (*models.Todo, error) {
	fmt.Println("not implemented")

	return nil, nil
}

// Destroy : Remove a single todo.
//  Method: PUT/PATCH
//  Params: True (id)
func (m *TodoModel) Destroy(title string) (string, error) {
	fmt.Println("not implemented")

	return "", nil
}
