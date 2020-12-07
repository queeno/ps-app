package main

import "net/http"

func main() {
	Logger.Debug("starting app")

	if err := RunWebserver(); err != http.ErrServerClosed {
		Logger.Fatal(err)
	}
}
