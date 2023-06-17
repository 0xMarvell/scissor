package main

import (
	"log"

	_ "github.com/0xMarvell/scissor/docs"
	"github.com/0xMarvell/scissor/pkg/config"
	"github.com/0xMarvell/scissor/pkg/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	config.Connect()
	config.RunMigrations()
}

// @title Scissor API
// @description Shorten urls for users
// @version 1.0
// @host localhost:8080
// @BasePath /api/v1
func main() {
	r := gin.Default()
	routes.RegisterRoutes(r)

	log.Println("server is up and running...")
	log.Fatal(r.Run())
}
