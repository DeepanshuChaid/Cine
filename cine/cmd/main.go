package main

import (
	"fmt"
	"log"
	"os"

	moviecontroller "github.com/DeepanshuChaid/Cine/tree/main/cine/internal/controllers/movieController"
	usercontroller "github.com/DeepanshuChaid/Cine/tree/main/cine/internal/controllers/userController"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/database"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found")
	}

	database.Connect()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	auth := router.Group("/api")
	auth.Use(middleware.AuthMiddleware())

	// MOVIE ROUTES
  auth.GET("/movies", moviecontroller.GetAllMovies())
	auth.GET("/movies/:id", moviecontroller.GetMovie())
	auth.POST("/create/movie", moviecontroller.CreateMovie())

	// USER ROUTES
	router.POST("/api/register", usercontroller.Register())
	router.POST("/api/login", usercontroller.Login())

	PORT := os.Getenv("PORT")

	router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET","POST","PUT","DELETE"},
			AllowHeaders:     []string{"Content-Type"},
			AllowCredentials: true,
	}))

	if err := router.Run(":" + PORT); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
    