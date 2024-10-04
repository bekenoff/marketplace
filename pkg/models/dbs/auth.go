package dbs

import (
	"context"
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

func (m *ClientModel) Insert(telephone int, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
        INSERT INTO user
        (telephone, password) 
        VALUES (?, ?);`

	_, err = m.DB.Exec(stmt, telephone, string(hashedPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return models.ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *ClientModel) GetPasswordByTelephone(telephone int) (string, error) {
	var password string
	stmt := `SELECT password FROM user WHERE telephone = ?`
	err := m.DB.QueryRow(stmt, telephone).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (m *ClientModel) InsertLaw(client *models.ClientLaw) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Password), 12)
	if err != nil {
		return err
	}

	stmt := `
        INSERT INTO user_law
        (company_name, contact_name, password, law_address, email, phone, bin, bik, iik, bank) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = m.DB.Exec(stmt, client.CompanyName, client.ContactName, string(hashedPassword), client.LawAddress, client.Email, client.Phone, client.Bin, client.Bik, client.Iik, client.Bank)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return models.ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *ClientModel) GetUserById(id string) ([]byte, error) {
	stmt := `SELECT * FROM user WHERE id = ?`

	userRow := m.DB.QueryRow(stmt, id)

	c := &models.Client{}

	err := userRow.Scan(&c.Id, &c.Telephone, &c.Password, &c.Created_at, &c.Modified_at)
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

func (m *ClientModel) Authenticate(password string, telephone int) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, password FROM user WHERE telephone = ?"
	row := m.DB.QueryRow(stmt, telephone)
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

func (m *ClientModel) GetClientPhoneById(idclient string) (string, error) {
	stmt := `SELECT telephone FROM user WHERE id = ?`
	var clientphone string
	err := m.DB.QueryRow(stmt, idclient).Scan(&clientphone)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("no client found with the given id")
		}
		return "", err
	}
	return clientphone, nil
}

func (m *ClientModel) ChangePassword(id int, oldPassword, newPassword string) error {

	var hashedPassword string
	stmt := "SELECT password FROM user WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&hashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oldPassword))
	if err != nil {
		// Если пароль не совпадает, возвращаем ошибку
		return errors.New("incorrect old password")
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	updateStmt := `
		UPDATE client
		SET clientpass = ?
		WHERE idclient = ?`

	_, err = m.DB.Exec(updateStmt, string(hashedNewPassword), id)
	if err != nil {
		return err
	}

	return nil
}

func (m *ClientModel) SetSession(ctx context.Context, id string, session models.Session) error {

	query := `
		UPDATE users 
		SET refresh_token = ?, expires_at = ? 
		WHERE id = ?
	`

	result, err := m.DB.ExecContext(ctx, query, session.RefreshToken, session.ExpiresAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}
