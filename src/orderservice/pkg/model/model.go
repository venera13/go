package model

import (
	"database/sql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	Database *sql.DB
}

type Order struct {
	Id    string `json:"id"`
	Items []OrderItem
}

type Orders struct {
	Orders []Order
}

type OrderItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type CreateOrderResponse struct {
	MenuItems []OrderItem `json:"MenuItems"`
}

type OrderResponse struct {
	Order
	OrderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type OrderServiceInterface interface {
	CreateOrder(orderResponse CreateOrderResponse) int
	GetOrder(orderId string) (OrderResponse, int)
}

func (s *Server) CreateOrder(orderResponse CreateOrderResponse) int {
	status := http.StatusNotFound
	orderId := uuid.NewString()
	timestamp := time.Now().Unix()
	cost := 999
	query := "INSERT INTO orders(order_id, time, cost) VALUES (?, ?, ?)"
	_, err := s.Database.Exec(query, orderId, timestamp, cost)
	status = LogError(err)

	for _, menuItem := range orderResponse.MenuItems {
		query = "INSERT INTO order_item(menu_id, quantity) VALUES (?, ?)"
		quantity := 1
		if menuItem.Quantity != 0 {
			quantity = menuItem.Quantity
		}
		result, err := s.Database.Exec(query, menuItem.Id, quantity)
		status = LogError(err)
		lastInsertId, err := result.LastInsertId()
		status = LogError(err)
		query = "INSERT INTO item_in_order(order_id, item_id) VALUES (?, ?)"
		_, err = s.Database.Exec(query, orderId, lastInsertId)
		status = LogError(err)
	}
	return status
}

func (s *Server) GetOrder(orderId string) (OrderResponse, int) {
	status := http.StatusNotFound
	var order OrderResponse
	order.Id = orderId
	query := "SELECT time, cost FROM orders where order_id = ? "
	err := s.Database.QueryRow(query, orderId).Scan(&order.OrderedAtTimeStamp, &order.Cost)
	if err != nil {
		status = LogError(err)
		return OrderResponse{}, status
	}

	query = "SELECT item_id FROM item_in_order where order_id = ? "
	items, err := s.Database.Query(query, orderId)
	if err != nil {
		status = LogError(err)
		return OrderResponse{}, status
	}

	orderItemsId := make([]int, 0)

	for items.Next() {
		var orderId int
		err := items.Scan(&orderId)
		status = LogError(err)
		orderItemsId = append(orderItemsId, orderId)
	}

	query = "SELECT menu_id, quantity FROM order_item where id = ? "
	orderItems := make([]OrderItem, 0)
	for _, orderItemId := range orderItemsId {
		var orderItem OrderItem
		err := s.Database.QueryRow(query, orderItemId).Scan(&orderItem.Id, &orderItem.Quantity)
		if err != nil {
			status = LogError(err)
			return OrderResponse{}, status
		}
		orderItems = append(orderItems, orderItem)
	}

	order.Items = orderItems
	return order, status
}

func LogError(err error) int {
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("mysql query error")
		return http.StatusForbidden
	}
	return http.StatusOK
}
