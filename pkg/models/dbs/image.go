package dbs

import (
	"database/sql"
	"marketplace/pkg/models"
)

type ImageModel struct {
	DB *sql.DB
}

func (i *ImageModel) Insert(image *models.Image) error {

	stmt := `
	INSERT INTO image
	(product_id, image_url)
	VALUES
	(?, ?);
	`

	_, err := i.DB.Exec(stmt, image.Product_id, image.Image_url)

	if err != nil {
		return err
	}

	return nil

}
