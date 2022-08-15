package controllers

import (
	"fmt"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/models"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/responses"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

/*
 * This file contains server management code
 *
 */

// Server represents the server, we can add more things in here to expand functionality
// A struct is used for better management of scope 
type Server struct {
	db     *gorm.DB
	router *mux.Router
}

// Init initializes the server
func (server *Server) Init(DBConnectionString string) {
	var err error
	// Make a connection to the database using GORM
	server.db, err = gorm.Open(postgres.Open(DBConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	// Migragte the tables
	err = server.db.AutoMigrate(&models.LinkedAccount{}, &models.SubAccount{})
	if err != nil {
		return
	}

	// Create a new mux router
	server.router = mux.NewRouter()
	// Initialize the routes
	server.initializeRoutes()
}

// Run starts the http server
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port " + addr)
	log.Fatal(http.ListenAndServe(addr, server.router))
}

// A simple test route
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome SubAccount API")

}
