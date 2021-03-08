package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"orderservice/pkg/model"
	_ "orderservice/pkg/model"
	"strconv"
	"time"
)

type order struct {
	Id    string `json:"id"`
	Items []orderItem
}

type orderResponse struct {
	order
	OrderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type orderItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type orders struct {
	Orders []order
}

func Router(srv *model.Server) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/orders", getOrders).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID}", getOrder).Methods(http.MethodGet)
	s.HandleFunc("/order", CreateOrder(srv)).Methods(http.MethodPost)
	return logMiddleware(r)
}

func getOrders(w http.ResponseWriter, _ *http.Request) {
	orders := orders{
		Orders: []order{
			{
				Id: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				Items: []orderItem{
					{
						Id:       "3fa85f64-5717-4562-b3fc-2c963f66afa6",
						Quantity: 1,
					},
				},
			},
		},
	}
	b, err := json.Marshal(orders)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("orders marshal error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = io.WriteString(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["ID"] == "" {
		http.Error(w, "order not found", http.StatusBadRequest)
	}
	order := orderResponse{
		order: order{
			Id: vars["ID"],
			Items: []orderItem{
				{
					Id:       vars["ID"],
					Quantity: 1,
				},
			},
		},
		OrderedAtTimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		Cost:               999,
	}

	b, err := json.Marshal(order)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("order marshal error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = io.WriteString(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}

func CreateOrder(serviceInterface model.OrderServiceInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusNotFound
		b, err := ioutil.ReadAll(r.Body)
		status = model.LogError(err)
		defer func(Body io.ReadCloser) {
			status = model.LogError(err)
		}(r.Body)
		var msg model.CreateOrderResponse
		err = json.Unmarshal(b, &msg)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Info("unmarshal error")
			status = http.StatusForbidden
		}
		status = serviceInterface.CreateOrder(msg)
		w.WriteHeader(status)
	}
}
