package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main () {

  router := gin.Default()

  router.GET("/", func(c *gin.Context) {
    c.String(200, "Hello World")
  })

  if err := router.Run(":3000"); err != nil {
    fmt.Println("Error starting server:", err)
  }
  
}

