package main

import (
	"fmt"
	"log"
	"os"

	moviecontroller "github.com/DeepanshuChaid/Cine/tree/main/cine/internal/controllers/movieController"
	usercontroller "github.com/DeepanshuChaid/Cine/tree/main/cine/internal/controllers/userController"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

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

	// MOVIE ROUTES
	router.GET("/movies", moviecontroller.GetAllMovies())
	router.GET("/movies/:id", moviecontroller.GetMovie())
	router.POST("/create/movie", moviecontroller.CreateMovie())

	// USER ROUTES
	router.POST("/register", usercontroller.Register())

	PORT := os.Getenv("PORT")

	if err := router.Run(":" + PORT); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
