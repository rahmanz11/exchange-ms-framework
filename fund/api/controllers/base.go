package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	kafkaHandler "github.com/Nabeegh-Ahmed/fund/api/kafka"
	"github.com/Nabeegh-Ahmed/fund/api/models"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

/*
 * This file contains server management code
 *
 */

// Server struct
type Server struct {
	kafkaReader *kafka.Reader
	kafkaWriter *kafka.Writer
}

// Init function
func (server *Server) Init() {
	// Create the kafka reader
	server.kafkaReader = kafkaHandler.KafkaReader(os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_READ_GROUP_ID"), os.Getenv("KAFKA_READ_TOPIC"))
	// Create the kafka writer
	server.kafkaWriter = kafkaHandler.KafkaWriter(os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_WRITE_TOPIC"))
}

// Run the server
func (server *Server) Run() {
	// Start the kafka writer
	// Closure for error handling
	defer func(kafkaWriter *kafka.Writer) {
		err := kafkaWriter.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(server.kafkaWriter)
	// Start the kafka reader
	defer func(kafkaReader *kafka.Reader) {
		err := kafkaReader.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(server.kafkaReader)
	// Start the kafka reader
	fmt.Println("Listening to kafka")
	for {
		// Read the message from kafka
		msg, err := server.kafkaReader.ReadMessage(context.Background())
		fmt.Println("Received message: ", string(msg.Value))
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Unmarshal the message
		exchangeOrder := models.ExchangeOrder{}
		err = json.Unmarshal(msg.Value, &exchangeOrder)

		if err != nil {
			fmt.Println(err)
			continue
		}
		// Execute the order
		err = exchangeOrder.ExchangeRequestHandling()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Push the message to kafka
		exchangeOrderJson, err := json.Marshal(exchangeOrder)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = server.KafkaPush(context.Background(), []byte(exchangeOrder.TransactionId), exchangeOrderJson)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

// KafkaPush function to push the message to kafka
func (server *Server) KafkaPush(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return server.kafkaWriter.WriteMessages(parent, message)
}
