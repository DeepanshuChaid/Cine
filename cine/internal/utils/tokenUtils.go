package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SignedDetails struct {
  Email string
  Username string
  Role string
  UserId string
  jwt.RegisteredClaims
}

func GenerateAllTokens(email, username, role, userid string) (string, string, error) {
  SECRET_KEY := os.Getenv("SECRET_KEY")
  SECRET_REFRESH_KEY := os.Getenv("SECRET_REFRESH_KEY")
   
  claims := &SignedDetails{
    Email: email,
    Username: username,
    Role: role,
    UserId: userid,
    RegisteredClaims: jwt.RegisteredClaims{
      Issuer: "NiggaNo1",
      IssuedAt: jwt.NewNumericDate(time.Now()),
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signedToken, err := token.SignedString([]byte(SECRET_KEY))

  if err != nil {
    return "", "", err
  }

  refreshClaims := &SignedDetails{
    Email: email,
    Username: username,
    Role: role,
    UserId: userid,
    RegisteredClaims: jwt.RegisteredClaims{
      Issuer: "NiggaNo1",
      IssuedAt: jwt.NewNumericDate(time.Now()),
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(7*24*time.Hour)),
    },
  }

  refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
  signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))

  if err != nil {
    return "", "", err
  }

  return signedToken, signedRefreshToken, nil
}

func UpdateAllTokens(userId, token, refreshToken string)(err error) {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    query := `
      UPDATE users
      SET token = $1,
          refreshtoken = $2
      WHERE id = $3
    `

    _, err = database.Pool.Exec(
      ctx,
      query,
      token,
      refreshToken,
      userId,
    )

    return err
}


func ValidateToken(signedToken string) (*SignedDetails, error) {
  jwtSecret := os.Getenv("SECRET_KEY")

  token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
      return []byte(jwtSecret), nil
  })

  if err != nil {
    return nil, err
  }

  claims, ok := token.Claims.(*SignedDetails)
  if !ok || !token.Valid {
    return nil, errors.New("invalid token")
  }

  return claims, nil
}


// GET USER ID FROM CONTEXT
func GetUserIdFromContext(c *gin.Context)(string, error) {
  userId, exists := c.Get("user_id")
  if !exists {
    return "", errors.New("user ID not found in context")
  }

  id, ok := userId.(string)
  if !ok {
    return "", errors.New("user ID is not a string")
  }

  return id, nil
}

// GET USER FAV GENRES FROM CONTEXT
func GetUserFavGenresFromContext(userId string,c *gin.Context)([]string, error) {
  ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
  defer cancel()

  var genres []string
  rows, err := database.Pool.Query(ctx, "SELECT genre_id FROM user_favourite_genres WHERE user_id = $1", userId)
  
  if err != nil {
    return nil, err
  }

  defer rows.Close()


  for rows.Next() {
    var genreID string
    if err := rows.Scan(&genreID); err != nil {
      return nil, err
    }

    genres = append(genres, genreID)
  }

  return genres, nil
}