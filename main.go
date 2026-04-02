package main

import (
	"ciphera-api/db"
	"ciphera-api/env"
	"ciphera-api/routes"
	"ciphera-api/socket"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	err := db.ConnectDB()
	if err != nil {
		fmt.Println(`error connecting to database.`)
		return
	}

	r := gin.New()

	go socket.RunWebSocketServer()

	routes.HandleRoutes(r)

	fmt.Println("API + WebSocket running on :" + env.GetSystemPort())
	log.Fatal(r.Run(":" + env.GetSystemPort()))
}
