package dbs

import (
	"database/sql"
	"marketplace/pkg/models"
)

type CartModel struct {
	DB *sql.DB
}

func (m *CartModel) Insert(cart *models.Cart) error {
	stmt := `
        INSERT INTO cart_item
        (client_id, product_id, quantity) 
        VALUES (?, ?, ?);`

	_, err := m.DB.Exec(stmt, cart.Client_id, cart.Product_id, cart.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (m *CartModel) GetByClientID(clientID int) ([]*models.Cart, error) {
	stmt := `
        SELECT product_id, client_id, quantity
        FROM cart
        WHERE client_id = ?;`

	rows, err := m.DB.Query(stmt, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []*models.Cart
	for rows.Next() {
		var cart models.Cart
		err = rows.Scan(&cart.Product_id, &cart.Client_id, &cart.Quantity)
		if err != nil {
			return nil, err
		}
		carts = append(carts, &cart)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return carts, nil
}

func (m *CartModel) Delete(clientID, productID int) error {
	stmt := `
        DELETE FROM cart_item
        WHERE client_id = ? AND product_id = ?;`

	_, err := m.DB.Exec(stmt, clientID, productID)
	if err != nil {
		return err
	}

	return nil
}
