package controllers

import (
	"github.com/Nabeegh-Ahmed/validate/api/middlewares"
)

// initializeRoutes initializes the routes
func (server *Server) initializeRoutes() {
	server.router.HandleFunc("/api/v1/", middlewares.SetMiddlewareJSON(server.SubmitExchangeOrder)).Methods("POST")
}
