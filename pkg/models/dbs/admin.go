package dbs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"marketplace/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

func (m *ClientModel) AuthenticateAdmin(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, password FROM admin WHERE email = ?"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *ClientModel) GetUserByIdAdmin(id string) ([]byte, error) {
	stmt := `SELECT * FROM admin WHERE id = ?`

	userRow := m.DB.QueryRow(stmt, id)

	c := &models.Client{}

	err := userRow.Scan(&c.Id, &c.Email, &c.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	convertedUser, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return convertedUser, nil
}
