package api

import (
	"fmt"
	"github.com/Nabeegh-Ahmed/fund/api/controllers"
	"github.com/joho/godotenv"
	"log"
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
	server.Run()
}
