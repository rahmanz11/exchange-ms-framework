package api

import (
	"fmt"
	"github.com/Nabeegh-Ahmed/validate/api/controllers"
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
	server.Init()
	server.Run(":" + os.Getenv("PORT"))
}
