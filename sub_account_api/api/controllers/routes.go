package controllers

import (
	"github.com/Nabeegh-Ahmed/sub_account_api/api/middlewares"
)

// initializeRoutes initializes the routes
// Functions use the middleware to convert request body to JSON for easier processing
func (server *Server) initializeRoutes() {
	// Auth Routes
	server.router.HandleFunc("/api/v1/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")
	server.router.HandleFunc("/api/v1/register", middlewares.SetMiddlewareJSON(server.Register)).Methods("POST")
	// Basic Home Route
	server.router.HandleFunc("/api/v1/", middlewares.SetMiddlewareJSON(server.Home)).Methods("GET")
	// SubAccount routes
	server.router.HandleFunc("/api/v1/{account-number}", middlewares.SetMiddlewareJSON(server.FindSubAccount)).Methods("GET")
	server.router.HandleFunc("/api/v1/id/{account-id}", middlewares.SetMiddlewareJSON(server.FindSubAccountById)).Methods("GET")
	server.router.HandleFunc("/api/v1/update/{account-number}", middlewares.SetMiddlewareJSON(server.UpdateSubAccount)).Methods("POST")
	// Linked Account routes
	server.router.HandleFunc("/api/v1/link-account", middlewares.SetMiddlewareJSON(server.CreateLinkedAccount)).Methods("POST")
}
