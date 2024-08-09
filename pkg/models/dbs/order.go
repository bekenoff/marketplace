package dbs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"marketplace/pkg/models"
)

type OrderModel struct {
	DB *sql.DB
}

func (m *OrderModel) Insert(userID int, status, address string, price, productID, quantity int) error {
	
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()


	stmtOrder := `INSERT INTO ` + "`order`" + ` (user_id, status, address, price, product_id, quantity) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(stmtOrder, userID, status, address, price, productID, quantity)
	if err != nil {
		return err
	}


	stmtInventory := `UPDATE product_inventory SET quantity = quantity - ? WHERE product_id = ?`
	_, err = tx.Exec(stmtInventory, quantity, productID)
	if err != nil {
		return err
	}

	
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *OrderModel) GetOrderById(id int) ([]byte, error) {
	stmt := `SELECT * FROM ` + "`order`" + ` WHERE id = ?`

	orderRow := m.DB.QueryRow(stmt, id)

	o := &models.Order{}

	err := orderRow.Scan(&o.Id, &o.User_id, &o.Status, &o.Address, &o.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	convertedOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return convertedOrder, nil
}

func (m *OrderModel) UpdateStatusByUserID(userID int, status string) error {
	stmt := `UPDATE ` + "`order`" + ` SET status = ? WHERE user_id = ?`

	result, err := m.DB.Exec(stmt, status, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	return nil
}

func (m *OrderModel) InsertOrderItem(orderID, productID, price, qty int) error {
	stmt := `
        INSERT INTO order_item (order_id, product_id, price, qty)
        VALUES (?, ?, ?, ?);`

	_, err := m.DB.Exec(stmt, orderID, productID, price, qty)
	if err != nil {
		return err
	}

	return nil
}

func (m *OrderModel) DecreaseProductInventory(productID, quantity int) error {
	stmt := `UPDATE product_inventory 
             SET quantity = quantity - ? 
             WHERE product_id = ?`
	_, err := m.DB.Exec(stmt, quantity, productID)
	if err != nil {
		return err
	}
	return nil
}

func (m *OrderModel) GetOrderItemById(id int) ([]byte, error) {
	stmt := `SELECT * FROM order_item WHERE order_id = ?`

	orderItemRow := m.DB.QueryRow(stmt, id)

	oi := &models.OrderItem{}

	err := orderItemRow.Scan(&oi.OrderID, &oi.ProductID, &oi.Price, &oi.Qty)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	convertedOrderItem, err := json.Marshal(oi)
	if err != nil {
		return nil, err
	}
	return convertedOrderItem, nil
}
