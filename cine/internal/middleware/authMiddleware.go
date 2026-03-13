package middleware

import (
	"net/http"

	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/utils"
	"github.com/gin-gonic/gin"
)

func authMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    token, err := c.Cookie("access_token")
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
        "error": "Unauthorized",
      })
      c.Abort()
      return
    }

    claims, err := utils.ValidateToken(token)
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
        "error": "Unauthorized",
      })
      c.Abort()
      return
    }

    // store user info in request context
    c.Set("user_id", claims.UserId)
    c.Set("user_email", claims.Email)
    c.Set("user_role", claims.Role)

    c.Next()

    
    
  }
}