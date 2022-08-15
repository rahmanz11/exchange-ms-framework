package api

import (
	"fmt"
	"github.com/Nabeegh-Ahmed/exchange_api/api/controllers"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var server = controllers.Server{}

// Run the server
func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Init()
	server.Run(":" + os.Getenv("PORT"))
}
