package dbs

import (
	"database/sql"
	"errors"
	"marketplace/pkg/models"
)

type InformationModel struct {
	DB *sql.DB
}

func (m *InformationModel) Insert(info *models.Information) error {
	stmt := `
        INSERT INTO information
        (product_id, articul, brand, series, country, color, quantity, size, packing_size) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := m.DB.Exec(stmt, info.Product_id, info.Articul, info.Brand, info.Series, info.Country, info.Color, info.Quantity, info.Size, info.Packing_size)
	if err != nil {
		return err
	}

	return nil
}

func (m *InformationModel) GetInformation(id int) (*models.Information, error) {
	stmt := `SELECT id, name FROM information WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	var information models.Information
	err := row.Scan(&information.Product_id, &information.Articul, &information.Brand, &information.Series, &information.Country, &information.Color, &information.Quantity, &information.Size, &information.Packing_size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &information, nil
}
