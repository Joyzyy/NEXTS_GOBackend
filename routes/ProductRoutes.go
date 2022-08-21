package routes

import (
	"example/hello/controllers"
	"example/hello/middlewares"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	router.POST("/createProduct", controllers.CreateProduct())
	router.GET("/getProducts", middlewares.AuthMiddleware(), controllers.GetAllProducts())
	router.GET("/getProduct/:id", controllers.FindProductByID())
}
