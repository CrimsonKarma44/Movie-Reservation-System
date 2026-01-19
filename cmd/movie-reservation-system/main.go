package main

import (
	"fmt"
	server "movie-reservation-system/server"
)

var err error
var app *server.App

func main() {
	var refreshStore = make(map[uint]string)
	err = app.Init(refreshStore)
	if err != nil {
		fmt.Printf("Error initializing application: %v\n", err)
	}
}
