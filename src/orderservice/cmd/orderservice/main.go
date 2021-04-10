package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderservice/pkg/orderservice/model"
	"orderservice/pkg/orderservice/transport"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	config, err := parseEnv()
	if err != nil {
		log.Fatal(err)
	}

	killSignalChat := getKillSignalChan()

	var srv *http.Server
	srv, err = startServer(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	waitForKillSignal(killSignalChat)
	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
}

func startServer(config *config) (*http.Server, error) {
	serverUrl := config.ServeRESTAddress
	log.WithFields(log.Fields{
		"url": serverUrl,
	}).Info("starting the server")

	db, err := createDBConn(config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	server := makeServer(db)
	router := transport.Router(server)
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv, nil
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}

func makeServer(db *sql.DB) *transport.Server {
	return &transport.Server{
		OrderService: model.NewServer(db),
	}
}

func createDBConn(config *config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@/%s", config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}
