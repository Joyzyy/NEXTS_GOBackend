package middlewares

import (
	"example/hello/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		clientToken := g.Request.Header.Get("Authorization")
		cookieToken := g.Request.Header.Get("Cookies")

		if clientToken == "" || cookieToken == "" {
			g.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			g.Abort()
			return
		}

		// Authorization: Bearer jwt -> get jwt after 7 characters.
		headerJWT := clientToken[7:]
		if headerJWT == "nil" {
			// if cookie token != "" && e valid atunci update jwt and send back a cookie
			g.JSON(http.StatusBadGateway, gin.H{"message": "eroare"})
			g.Abort()
			return
		}
		claims, returnStatus, errMsg := utils.Verify(headerJWT)
		if errMsg != "" {
			g.JSON(returnStatus, gin.H{"message": errMsg})
			g.Abort()
			return
		}

		g.Set("id", claims.Id)
		g.Next()
	}
}
