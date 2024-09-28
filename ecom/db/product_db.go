package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/dkr290/go-advanced-projects/ecom/types"
)

type ProductDatabaseInt interface {
	GetProducts() ([]types.Product, error)
	CreateProduct(types.CreateProductPayload) error
	UpdateProduct(types.CreateProductPayload, int) error
	GetProductById(id int) (*types.Product, error)
}

type ProductMysqlDB struct {
	DB *sql.DB
}

func (p *ProductMysqlDB) GetProducts() ([]types.Product, error) {
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

func (p *ProductMysqlDB) CreateProduct(product types.CreateProductPayload) error {

	// Check if product name already exists
	var count int
	err := p.DB.QueryRow("SELECT COUNT(*) FROM products WHERE name = ?", product.Name).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("product with name %s already exists", product.Name)

	}

	_, err = p.DB.Exec("INSERT INTO products(name,description,image,price,quantity) VALUES(?,?,?,?,?)",
		product.Name, product.Description, product.Image, product.Price, product.Quantity)

	if err != nil {
		return err
	}
	return nil
}

func (p *ProductMysqlDB) UpdateProduct(product types.CreateProductPayload, id int) error {
	_, err := p.DB.Exec("UPDATE products SET name = ?, description = ?, image = ?, price = ?, quantity = ? WHERE id = ?",
		product.Name, product.Description, product.Image, product.Price, product.Quantity, id)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductMysqlDB) GetProductById(id int) (*types.Product, error) {
	row := p.DB.QueryRow("SELECT * FROM products WHERE id = ?", id)

	product := &types.Product{}
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
	}
	return product, nil

}
