package main

import (
	"ciphera-api/db"
	"ciphera-api/env"
	"ciphera-api/routes"
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

	//go socket.RunWebSocketServer()

	routes.HandleRoutes(r)

	fmt.Println("API + Socket.IO running on :" + env.GetSystemPort())
	log.Fatal(r.Run(":" + env.GetSystemPort()))
}
