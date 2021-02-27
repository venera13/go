package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Order struct {
	Id string `json:"id"`
	Quantity string `json:"quantity"`

}

type Orders struct {
	Id string `json:"id"`
	MenuItems []Order `json:"menuItems"`

}

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/orders", orders).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID}", order).Methods(http.MethodGet)
	return logMiddleware(r)
}

func helloWorld(w http.ResponseWriter, _ *http.Request){
	fmt.Fprint(w, "Hello World!")
}

func orders(w http.ResponseWriter, _ *http.Request){
	orders := Orders{
		Id: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		MenuItems: []Order{{
			Id: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			Quantity: "1",
		}},
	}
	b, err := json.Marshal(orders)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("orders marshal error")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(b))
}

func order(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	order := Order{
		Id: vars["ID"],
		Quantity: "1",
	}
	b, err := json.Marshal(order)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("order marshal error")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(b))
}

func logMiddleware(h http.Handler) http.Handler{
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request){
		log.WithFields(log.Fields{
			"method": r.Method,
			"url":    r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent": r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}
