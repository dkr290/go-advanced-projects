package db

import (
	"database/sql"

	"github.com/dkr290/go-advanced-projects/ecom/types"
)

type ProductDatabase interface {
	GetProducts() ([]types.Product, error)
}

type ProductMysql struct {
	DB *sql.DB
}

func (p *ProductMysql) GetProducts() ([]types.Product, error) {
	rows, err := p.DB.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)
	var product types.Product
	for rows.Next() {
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Image,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
