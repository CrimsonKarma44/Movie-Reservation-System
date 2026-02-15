package main

import (
	"fmt"
	server "movie-reservation-system/server"
)

var err error
var app *server.App

func main() {
	err = app.Init()
	if err != nil {
		fmt.Printf("Error initializing application: %v\n", err)
	}
}
