package postgres

import (
	"database/sql"
	"fakorede-bolu/full-rest-api/pkg/models"
)

// InboxModel struct
type InboxModel struct {
	DB *sql.DB
}

// Outbox : for storing emails yet to be authenticated
func (m *InboxModel) Outbox(send int) ([]*models.Inbox, error) {
	stmt := `SELECT * FROM inbox WHERE send = $1;`

	rows, err := m.DB.Query(stmt, send)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	inbox := []*models.Inbox{}

	for rows.Next() {

		i := &models.Inbox{}

		err = rows.Scan(&i.Email, &i.UserID, &i.Send)

		if err != nil {
			return nil, err
		}

		inbox = append(inbox, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inbox, nil
}
