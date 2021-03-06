package main

import (
	"context"
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"net/http"
	transport "orderservice/pkg/orderservice"
	"os"
	"os/signal"
	"syscall"
)

//type Server struct {
//	db *sql.DB
//}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	serverUrl := ":8000"
	killSignalChat := getKillSignalChan()
	//db, err := sql.Open("mysql", "root:1234@/orders")
	if err != nil {
		log.Fatal(err)
	}
	//server := Server{db}
	//srv := server.startServer(serverUrl)
	log.WithFields(log.Fields{
		"url": serverUrl,
	}).Info("starting the server")
	srv := startServer(serverUrl)
	waitForKillSignal(killSignalChat)
	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func startServer(serverUrl string) *http.Server {
	router := transport.Router()
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
