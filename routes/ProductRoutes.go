package routes

import (
	"example/hello/controllers"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	router.POST("/createProduct", controllers.CreateProduct())
	router.GET("/getProducts", controllers.GetAllProducts())
	router.GET("/getProduct/:id", controllers.FindProductByID())
}
