package usercontroller

import (
	"context"
	"errors"
	"net/http" // FIX: use HTTP status constants instead of raw numbers
	"time"

	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/database"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/models"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/utils"
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
        "error":   "User already Exists.",
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


// LOGIN
func Login() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    var RequestData models.UserLogin
    var foundUser models.FoundUser

    if err := c.ShouldBindJSON(&RequestData); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{
        "error":   "Invalid request body",
        "details": err.Error(),
      })
      return
    }

    if err := validate.Struct(RequestData); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{
        "error":   "Validation failed",
        "details": err.Error(),
      })
      return
    }

    err := database.Pool.QueryRow(ctx, "SELECT id, username, email, password, role, favouritegeneres FROM users WHERE email = $1", RequestData.Email).Scan(
      &foundUser.UserId,
      &foundUser.Username,
      &foundUser.Email,
      &foundUser.Password,
      &foundUser.Role,
      &foundUser.Favouritegeneres,
      )
    if err != nil {
      if errors.Is(err, pgx.ErrNoRows) {
        c.JSON(http.StatusNotFound, gin.H{
          "error": "User not found",
        })
        return
      }

      c.JSON(http.StatusInternalServerError, gin.H{
        "error":   "Database error",
        "details": err.Error(),
      })
      return
    }

    err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(RequestData.Password))
    if err != nil {
      c.JSON(http.StatusUnauthorized, gin.H{
        "error": "Invalid password",
      })
      return
    }

    token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.Username, foundUser.Role, foundUser.UserId)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error":   "Failed to generate tokens",
        "details": err.Error(),
      })
      return
    }

    err = utils.UpdateAllTokens(foundUser.UserId, token, refreshToken)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error":   "Failed to update tokens",
        "details": err.Error(),
      })
      return
    }


    c.JSON(http.StatusOK, models.FoundUser{
      UserId: foundUser.UserId,
      Username: foundUser.Username,
      Email: foundUser.Email,
      Role: foundUser.Role,
      Token: token,
      Refreshtoken: refreshToken,
      Favouritegeneres: foundUser.Favouritegeneres,
    })
    
    
  }
}
