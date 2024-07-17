package dbs

import (
	"database/sql"
	"marketplace/pkg/models"
)

type FavModel struct {
	DB *sql.DB
}

func (m *FavModel) Insert(fav *models.Favorites) error {
	stmt := `
        INSERT INTO favorites
        (product_id, client_id) 
        VALUES (?, ?);`

	_, err := m.DB.Exec(stmt, fav.Product_id, fav.Client_id)
	if err != nil {
		return err
	}

	return nil
}

func (m *FavModel) GetByClientID(clientID int) ([]*models.Favorites, error) {
	stmt := `
        SELECT product_id, client_id
        FROM favorites
        WHERE client_id = ?;`

	rows, err := m.DB.Query(stmt, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []*models.Favorites
	for rows.Next() {
		var fav models.Favorites
		err = rows.Scan(&fav.Product_id, &fav.Client_id)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, &fav)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return favorites, nil
}
