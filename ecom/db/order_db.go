package db

import (
	"database/sql"

	"github.com/dkr290/go-advanced-projects/ecom/types"
)

type OrderDatabaseInt interface {
	CreateOrder(types.Order) (int, error)
	CreateOrderItem(types.OrderItem) error
}

type OrderMysqlDB struct {
	DB *sql.DB
}

func (o *OrderMysqlDB) CreateOrder(order types.Order) (int, error) {

	res, err := o.DB.Exec("INSERT INTO orders(userID,total,status,address) VALUES (?,?,?,?)",
		order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (o *OrderMysqlDB) CreateOrderItem(ord types.OrderItem) error {

	_, err := o.DB.Exec("INSERT INTO order_items(orderID,productID,quantity,price) VALUES (?,?,?,?)",
		ord.OrderID, ord.ProductID, ord.Quantity, ord.Price)

	if err != nil {
		return err
	}

	return nil
}
