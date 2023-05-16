package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello")

	r := gin.Default()
	log.Println("server is up and running...")
	log.Fatal(r.Run())
}
