package main

import (
	"example/hello/configs"
	"example/hello/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://nextjs-app-charka-frontend.herokuapp.com"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))

	configs.Connect()
	routes.ProductRoutes(router)
	routes.UserRoutes(router)

	gin.SetMode(gin.ReleaseMode)

	router.Run()
}
