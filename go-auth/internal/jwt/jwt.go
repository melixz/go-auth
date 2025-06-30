package jwt

import (
	"errors"
	"go-auth/internal/models"
	"os"
	"time"

	"github.com/kataras/jwt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

const tokenTTL = 15 * time.Minute

func GenerateAccessToken(userID string) (string, error) {
	claims := models.AccessTokenClaims{
		UserID:   userID,
		IssuedAt: time.Now().Unix(),
		Expiry:   time.Now().Add(tokenTTL).Unix(),
	}
	token, err := jwt.Sign(jwt.HS512, secretKey, claims, jwt.MaxAge(tokenTTL))
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func ValidateAccessToken(token string) (*models.AccessTokenClaims, error) {
	verified, err := jwt.Verify(jwt.HS512, secretKey, []byte(token))
	if err != nil {
		return nil, err
	}
	var claims models.AccessTokenClaims
	err = verified.Claims(&claims)
	if err != nil {
		return nil, err
	}
	if claims.Expiry < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}
