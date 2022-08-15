package controllers

import (
	"github.com/Nabeegh-Ahmed/exchange_api/api/middlewares"
)

func (server *Server) initializeRoutes() {
	// Auth Routes
	server.router.HandleFunc("/api/v1/link-account", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.LinkAccount))).Methods("POST")
	server.router.HandleFunc("/api/v1/exchange", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.Exchange))).Methods("POST")
	server.router.HandleFunc("/api/v1/wire-in", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.WireIn))).Methods("POST")
	server.router.HandleFunc("/api/v1/wire-out", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.WireOut))).Methods("POST")
	server.router.HandleFunc("/api/v1/status/", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.GetStatus))).Methods("GET")
	server.router.HandleFunc("/api/v1/enable/", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.EnableAccount))).Methods("POST")
	server.router.HandleFunc("/api/v1/disable/", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.DisableAccount))).Methods("POST")
	server.router.HandleFunc("/api/v1/confirmation", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.GetTransactionConfirmation))).Methods("GET")
	server.router.HandleFunc("/api/v1/balance", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.GetBalance))).Methods("GET")
	server.router.HandleFunc("/api/v1/transaction", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.GetTransaction))).Methods("GET")
	// Open Routes
	server.router.HandleFunc("/api/v1/create-sub-account", middlewares.SetMiddlewareJSON(server.CreateSubAccount)).Methods("POST")
	server.router.HandleFunc("/api/v1/register", middlewares.SetMiddlewareJSON(server.TransactionsRegister)).Methods("GET")
	server.router.HandleFunc("/api/v1/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")

}
