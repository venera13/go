package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Server struct {
	Database *sql.DB
}

type Orders struct {
	Orders []Order
}

type OrderItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	Id                 string `json:"id"`
	Items              []OrderItem
	OrderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type OrderServiceInterface interface {
	CreateOrder(orderItems *[]OrderItem) error
	GetOrder(orderId string) (Order, error)
	GetOrders() ([]Order, error)
}

func (s *Server) CreateOrder(orderItems *[]OrderItem) error {
	orderId := uuid.NewString()
	timestamp := time.Now().Unix()
	cost := 999
	query := "INSERT INTO orders(order_id, time, cost) VALUES (?, ?, ?)"
	_, err := s.Database.Exec(query, orderId, timestamp, cost)

	for _, menuItem := range *orderItems {
		query = "INSERT INTO order_item(menu_id, quantity) VALUES (?, ?)"
		quantity := 1
		if menuItem.Quantity != 0 {
			quantity = menuItem.Quantity
		}
		result, err := s.Database.Exec(query, menuItem.Id, quantity)
		if err != nil {
			return err
		}
		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		query = "INSERT INTO item_in_order(order_id, item_id) VALUES (?, ?)"
		_, err = s.Database.Exec(query, orderId, lastInsertId)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Server) GetOrder(orderId string) (Order, error) {
	var order Order
	order.Id = orderId
	query := "SELECT time, cost FROM orders where order_id = ? "
	err := s.Database.QueryRow(query, orderId).Scan(&order.OrderedAtTimeStamp, &order.Cost)
	if err != nil {
		return Order{}, err
	}

	query = "SELECT item_id FROM item_in_order where order_id = ? "
	var items *sql.Rows
	items, err = s.Database.Query(query, orderId)
	if err != nil {
		return Order{}, err
	}

	orderItemsId := make([]int, 0)

	for items.Next() {
		var orderId int
		err := items.Scan(&orderId)
		if err != nil {
			return Order{}, err
		}
		orderItemsId = append(orderItemsId, orderId)
	}

	query = "SELECT menu_id, quantity FROM order_item where id = ? "
	orderItems := make([]OrderItem, 0)
	for _, orderItemId := range orderItemsId {
		var orderItem OrderItem
		err := s.Database.QueryRow(query, orderItemId).Scan(&orderItem.Id, &orderItem.Quantity)
		if err != nil {
			return Order{}, err
		}
		orderItems = append(orderItems, orderItem)
	}

	order.Items = orderItems
	return order, err
}

func (s *Server) GetOrders() ([]Order, error) {
	query := "SELECT order_id FROM orders"
	ordersId, err := s.Database.Query(query)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)

	for ordersId.Next() {
		var orderId string
		err := ordersId.Scan(&orderId)
		if err != nil {
			return nil, err
		}
		var order Order
		order, err = s.GetOrder(orderId)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, err
}
