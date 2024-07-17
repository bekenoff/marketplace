package dbs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"marketplace/pkg/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type ClientModel struct {
	DB *sql.DB
}

func (m *ClientModel) Insert(username, password, email, first_name, last_name, telephone string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
        INSERT INTO user
        (username, password, email, first_name, last_name, telephone) 
        VALUES (?, ?, ?, ?, ?, ?);`

	_, err = m.DB.Exec(stmt, username, string(hashedPassword), email, first_name, last_name, telephone)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return models.ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *ClientModel) GetUserById(id string) ([]byte, error) {
	stmt := `SELECT * FROM astana.client WHERE id = ?`

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

func (m *ClientModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, password FROM astana.client WHERE email = ?"
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
