package middleware

import (
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		authHeader := c.GetHeader("Authorization")

		if authHeader == ""{
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})

			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":"invalid authorization header",
			})
			c.Abort()
			return
		}
		tokenString := parts[1]
		userID, err := utils.ValidToken(tokenString)

		if err != nil{
			c.JSON(http.StatusUnauthorized, gin.H{

				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}