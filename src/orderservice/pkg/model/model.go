package model

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	Database *sql.DB
}

type orderItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type createOrderResponse struct {
	MenuItems []orderItem `json:"MenuItems"`
}

func (s Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotFound
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("read body error")
		status = http.StatusForbidden
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Info("close error")
			status = http.StatusForbidden
		}
	}(r.Body)
	var msg createOrderResponse
	err = json.Unmarshal(b, &msg)
	log.WithFields(log.Fields{
		"msg": msg,
	}).Info("debug")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("unmarshal error")
		status = http.StatusForbidden
	}
	orderId := uuid.NewString()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	cost := 999
	query := "INSERT INTO 'orders' (id, time, cost) VALUES (?, ?, ?)"
	_, err = s.Database.Query(query, orderId, timestamp, cost)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("mysql query error")
		status = http.StatusForbidden
	}
	status = http.StatusOK
	w.WriteHeader(status)
}
