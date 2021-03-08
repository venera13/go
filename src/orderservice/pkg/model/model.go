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

type orderItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type CreateOrderResponse struct {
	MenuItems []orderItem `json:"MenuItems"`
}

type OrderServiceInterface interface {
	CreateOrder(orderResponse CreateOrderResponse) int
}

func (s *Server) CreateOrder(orderResponse CreateOrderResponse) int {
	orderId := uuid.NewString()
	timestamp := time.Now().Unix()
	cost := 999
	query := "INSERT INTO orders(order_id, time, cost) VALUES (?, ?, ?)"
	_, err := s.Database.Exec(query, orderId, timestamp, cost)
	status := LogError(err)

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

func LogError(err error) int {
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("mysql query error")
		return http.StatusForbidden
	}
	return http.StatusOK
}
