package main

import (
	"log"

	"github.com/0xMarvell/scissor/pkg/config"
	"github.com/0xMarvell/scissor/pkg/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	config.Connect()
	config.RunMigrations()
}

func main() {
	r := gin.Default()
	routes.RegisterRoutes(r)

	log.Println("server is up and running...")
	log.Fatal(r.Run())
}
