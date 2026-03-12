package usercontroller

import (
	"context"
	"time"

	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/database"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(HashPassword), nil
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		validate := validator.New()

		if err := validate.Struct(user); err != nil {
			c.JSON(400, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		err := database.Pool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", user.Email).Scan(&user.ID)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "The bitchass developer could not handle database properly!",
			})
			return
		}

		// if user actually exists, then return error
		if user.ID != "" {
			c.JSON(400, gin.H{
				"error": "User already exists",
			})
			return
		}

		hashedPassword, err := HashPassword(user.Password)

		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Failed to hash password",
				"details": err.Error(),
			})
			return
		}

		// Update the user's password with the hashed version
		user.Password = hashedPassword

		err = database.Pool.QueryRow(ctx, "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id", user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)

		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Failed to create user",
				"details": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"message": "User created successfully",
			"user":    user,
		})
	}
}
