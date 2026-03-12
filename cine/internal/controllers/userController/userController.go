package usercontroller

import (
	"context"
	"errors"
	"net/http" // FIX: use HTTP status constants instead of raw numbers
	"time"

	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/database"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// VADLIDATOR VARIABLE
var validate = validator.New(validator.WithRequiredStructEnabled())


func HashPassword(password string) (string, error) {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

  if err != nil {
    return "", err
  }

  return string(hashedPassword), nil
}

func Register() gin.HandlerFunc {
  return func(c *gin.Context) {

    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{
        "error":   "Invalid request body",
        "details": err.Error(),
      })
      return
    }

    if err := validate.Struct(user); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{
        "error":   "Validation failed",
        "details": err.Error(),
      })
      return
    }

    // check if user already exists
    err := database.Pool.
      QueryRow(ctx, "SELECT id FROM users WHERE email = $1", user.Email).
      Scan(&user.ID)

    if err != nil {
      if errors.Is(err, pgx.ErrNoRows) {
        c.JSON(http.StatusInternalServerError, gin.H{
          "error":   "Database error",
          "details": err.Error(),
        })
        return
      }
    } else {
      c.JSON(http.StatusBadRequest, gin.H{
        "error": "User already exists",
      })
      return
    }

    hashedPassword, err := HashPassword(user.Password)

    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error":   "Failed to hash password",
        "details": err.Error(),
      })
      return
    }

    user.Password = hashedPassword

    err = database.Pool.QueryRow(
      ctx,
      "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id",
      user.Username,
      user.Email,
      user.Password,
      user.Role,
    ).Scan(&user.ID)

    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error":   "Failed to create user",
        "details": err.Error(),
      })
      return
    }

    user.Password = ""

    c.JSON(http.StatusCreated, gin.H{
      "message": "User created successfully",
      "user":    user,
    })
  }
}