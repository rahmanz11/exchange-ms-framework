package api

import (
	"fmt"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/controllers"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Server controller
var server = controllers.Server{}

// Run the server
func Run() {
	var err error
	// Load the .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	// Run the server
	server.Init(os.Getenv("DB_CONNECTION_STRING"))
	server.Run(":" + os.Getenv("PORT"))
}
