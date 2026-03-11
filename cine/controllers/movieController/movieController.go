package moviecontroller 

import (
  "github.com/gin-gonic/gin"  
)

func GetAllMovies() gin.HandlerFunc {
  return func(c *gin.Context) {

    
    
    c.JSON(200, gin.H{
      "message": "I HAte nigga and i find sukhmeen really atttractive",
    })
  }
}

