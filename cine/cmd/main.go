package main

import (
	"fmt"
	"log"

	moviecontroller "github.com/DeepanshuChaid/Cine/tree/main/cine/controllers/movieController"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main () {

  err := godotenv.Load()

  if err != nil {
    log.Println(".env not found")
  }

  database.Connect()
  database.InitSchema()

  router := gin.Default()

  router.GET("/", func(c *gin.Context) {
    c.String(200, "Hello World")
  })

  router.GET("/movies", moviecontroller.GetAllMovies())

  if err := router.Run(":3000"); err != nil {
    fmt.Println("Error starting server:", err)
  }
  
}

