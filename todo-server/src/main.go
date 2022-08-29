package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const LISTEN_PORT = ":5000"


func main() {
	c := cors.AllowAll()

	if OpenDB() != nil {
		fmt.Println("Error: Failed to connect to the database")
		return
	}

	router := mux.NewRouter()

	server := &http.Server{}
	server.Handler = c.Handler(router)
	server.Addr = LISTEN_PORT

	router.Use(Middleware)
	SetRoutes(router)

	server.ListenAndServe()
}
