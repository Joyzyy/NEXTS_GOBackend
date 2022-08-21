package routes

import (
	"example/hello/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/auth/login", controllers.Login())
	router.POST("/auth/register", controllers.Register())
	router.POST("/auth/logout", controllers.Logout())
}
