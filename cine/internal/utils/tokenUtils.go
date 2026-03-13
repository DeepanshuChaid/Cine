package utils

import (
	"os"
  "context"
	"time"

	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type SignedDetails struct {
  Email string
  Username string
  Role string
  UserId string
  jwt.RegisteredClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY string = os.Getenv("SECRET_REFRESH_KEY")

func GenerateAllTokens(email, username, role, userid string) (string, string, error) {
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