package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunWebserver() error {
	server := &http.Server{
		Addr: "0.0.0.0:9292",
	}

	http.HandleFunc("/", HandleProducts)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		<-c

		Logger.Info("stopping webserver...")
		ctxShutdown, cancel := context.WithTimeout(context.TODO(), 5*time.Second)

		if err := server.Shutdown(ctxShutdown); err != nil {
			Logger.Fatal(err)
		}

		cancel()
		Logger.Info("webserver terminated")
	}()

	Logger.Info("starting webserver...")

	return server.ListenAndServe()
}
