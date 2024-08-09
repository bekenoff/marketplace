package dbs

import (
	"database/sql"
	"errors"
	"marketplace/pkg/models"
)

type ProductModel struct {
	DB *sql.DB
}

func (m *ProductModel) InsertProduct(name, description string, category_id, inventory_id, discount_id int, price float64) error {
	stmt := `
        INSERT INTO products (name, description, category_id, inventory_id, price, discount_id)
        VALUES (?, ?, ?, ?, ?, ?);
    `
	_, err := m.DB.Exec(stmt, name, description, category_id, inventory_id, price, discount_id)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProductModel) InsertProductInventory(quantity int) error {
	stmt := `
        INSERT INTO product_inventory
		(quantity)
        VALUES (?);
    `
	_, err := m.DB.Exec(stmt, quantity)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProductModel) InsertRating(productID, rating int, review string) error {
	stmt := `INSERT INTO reviews (product_id, rating, review) VALUES (?, ?, ?)`
	_, err := m.DB.Exec(stmt, productID, rating, review)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProductModel) GetAllProducts() ([]models.Product, error) {
	stmt := `SELECT id, name FROM products`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.Id, &product.Name)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (m *ProductModel) GetProductByID(id int) (*models.Product, error) {
	stmt := `SELECT * FROM products WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	var product models.Product
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Category_id, &product.Inventory_id, &product.Price, &product.Discount_id, &product.Created_at, &product.Modified_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &product, nil
}

func (m *ProductModel) GetProductByCategoryID(categoryID int) ([]models.Product, error) {
	stmt := `SELECT id, name, description, category_id, inventory_id, price, discount_id, created_at, modified_at 
             FROM products 
             WHERE category_id = ?`
	rows, err := m.DB.Query(stmt, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Category_id, &product.Inventory_id, &product.Price, &product.Discount_id, &product.Created_at, &product.Modified_at)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (m *ProductModel) GetReviewsByProductID(productID int) ([]models.Review, error) {
	stmt := `SELECT id, product_id, rating, review FROM reviews WHERE product_id = ?`
	rows, err := m.DB.Query(stmt, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var review models.Review
		err = rows.Scan(&review.ID, &review.ProductID, &review.Rating, &review.Review)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}
