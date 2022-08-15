package controllers

import (
	"context"
	"fmt"
	"github.com/Nabeegh-Ahmed/validate/api/kafka"
	"github.com/Nabeegh-Ahmed/validate/api/responses"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"os"
	"time"
)

/*
 * This file contains server management code
 *
 */

// Server struct
type Server struct {
	router      *mux.Router
	kafkaWriter *kafka.Writer
}

// Initialize the server
func (server *Server) Init() {
	// Initialize the router
	server.router = mux.NewRouter()
	// Initialize the routes
	server.initializeRoutes()
	// Initialize the kafka writer
	server.kafkaWriter = kafkaHandler.KafkaWriter(os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_TOPIC"))
}

// Run the server on port addr
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port " + addr)
	// Closure for error handling
	defer func(kafkaWriter *kafka.Writer) {
		err := kafkaWriter.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(server.kafkaWriter)
	// Run the server
	log.Fatal(http.ListenAndServe(addr, server.router))
}

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to Validate API")
}

// Push the message to kafka topic
func (server *Server) KafkaPush(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return server.kafkaWriter.WriteMessages(parent, message)
}
