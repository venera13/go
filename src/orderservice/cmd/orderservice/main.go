package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"net/http"
	model "orderservice/pkg/model"
	transport "orderservice/pkg/orderservice"
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

	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	serverUrl := config.ServeRESTAddress
	killSignalChat := getKillSignalChan()
	dataSourceName := fmt.Sprintf("%d:%d@/%d", config.Database.DBUser, config.Database.DBPass, config.Database.DBName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
		return
	}
	server := model.Server{Database: db}
	srv := startServer(serverUrl, server)
	log.WithFields(log.Fields{
		"url": serverUrl,
	}).Info("starting the server")
	waitForKillSignal(killSignalChat)
	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
}

func startServer(serverUrl string, server model.Server) *http.Server {
	router := transport.Router(&server)
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
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
