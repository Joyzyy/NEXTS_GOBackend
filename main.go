package main

import (
	"example/hello/configs"
	"example/hello/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	configs.Connect()
	routes.ProductRoutes(router)

	router.Run(":9999")
}
