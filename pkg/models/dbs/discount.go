package dbs

import (
	"database/sql"
	"marketplace/pkg/models"
)

type DiscountModel struct {
	DB *sql.DB
}

func (m *DiscountModel) Insert(discount *models.Discount) error {
	stmt := `
    INSERT INTO discount
    (product_id, name, description, discount_percent, active) 
    VALUES (?, ?, ?, ?, ?);
`

	_, err := m.DB.Exec(stmt, discount.Product_id, discount.Name, discount.Description, discount.Percent, discount.Active)
	if err != nil {
		return err
	}

	return nil
}
