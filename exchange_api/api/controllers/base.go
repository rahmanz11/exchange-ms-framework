package controllers

import (
	"fmt"
	"github.com/Nabeegh-Ahmed/exchange_api/api/responses"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

/*
 * This file contains server management code
 *
 */

type Server struct {
	router *mux.Router
}

func (server *Server) Init() {
	// Create a new instance of the router
	server.router = mux.NewRouter()
	// Initialize the routes
	server.initializeRoutes()
}

// Run the server
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port " + addr)
	log.Fatal(http.ListenAndServe(addr, server.router))
}

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to Exchange API")
}
