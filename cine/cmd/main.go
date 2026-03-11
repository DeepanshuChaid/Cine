package main

import (
	"fmt"
	"log"
  "os"

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

  PORT := os.Getenv("PORT")

  if err := router.Run(":" + PORT); err != nil {
    fmt.Println("Error starting server:", err)
  }
  
}

