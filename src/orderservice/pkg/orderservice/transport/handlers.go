package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"orderservice/pkg/orderservice/model"
)

type createOrderResponse struct {
	MenuItems []model.OrderItem `json:"MenuItems"`
}

func Router(serviceInterface model.OrderServiceInterface) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/orders", getOrders(serviceInterface)).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID}", getOrder(serviceInterface)).Methods(http.MethodGet)
	s.HandleFunc("/order", createOrder(serviceInterface)).Methods(http.MethodPost)
	return logMiddleware(r)
}

func getOrders(serviceInterface model.OrderServiceInterface) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//orders := model.Orders{
		//	Orders: []model.Order{
		//		{
		//			Id: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		//			Items: []model.OrderItem{
		//				{
		//					Id:       "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		//					Quantity: 1,
		//				},
		//			},
		//		},
		//	},
		//}
		orders, err := serviceInterface.GetOrders()
		if err != nil {
			logError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var b []byte
		b, err = json.Marshal(orders)
		if err != nil {
			logError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeResponse(w, http.StatusOK, string(b))
	}
}

func getOrder(serviceInterface model.OrderServiceInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		vars := mux.Vars(r)
		if vars["ID"] == "" {
			http.Error(w, "order not found", http.StatusBadRequest)
			return
		}
		order, err := serviceInterface.GetOrder(vars["ID"])
		if err == nil {
			b, err := json.Marshal(order)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Info("order marshal error")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeResponse(w, status, string(b))
		} else {
			writeResponse(w, status, "Order not found")
		}
	}
}

func createOrder(serviceInterface model.OrderServiceInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logError(err)
			status = http.StatusInternalServerError
			writeResponse(w, status, "Error")
			return
		}
		defer func(Body io.ReadCloser) {
			if err != nil {
				logError(err)
				status = http.StatusInternalServerError
				writeResponse(w, status, "Error")
				return
			}
		}(r.Body)
		var msg createOrderResponse
		err = json.Unmarshal(b, &msg)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Info("unmarshal error")
			status = http.StatusForbidden
			writeResponse(w, status, "Error")
			return
		}
		err = serviceInterface.CreateOrder(&msg.MenuItems)
		if err != nil {
			logError(err)
			return
		}
		writeResponse(w, status, "Success")
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

func writeResponse(w http.ResponseWriter, status int, response string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, err := io.WriteString(w, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func logError(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		})
	}
}
