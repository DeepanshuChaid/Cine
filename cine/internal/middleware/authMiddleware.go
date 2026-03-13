package middleware

import (
	"fmt"
	"net/http"

	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("AUTH MIDDLEWARE ENTERED")

		token, err := c.Cookie("access_token")

		fmt.Println("COOKIE VALUE:", token)
		fmt.Println("COOKIE ERROR:", err)

		if err != nil {
			fmt.Println("NO COOKIE FOUND")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		claims, err := utils.ValidateToken(token)

		fmt.Println("TOKEN VALIDATION ERROR:", err)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		fmt.Println("TOKEN VALID")

		c.Set("user_id", claims.UserId)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

