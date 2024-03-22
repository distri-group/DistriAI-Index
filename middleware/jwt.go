package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// TokenExpireDuration JWT expire duration
const TokenExpireDuration = time.Hour * 24 * 7

// CustomSecret Signature Secret
var CustomSecret = []byte("ai-distri-ai")

// CustomClaims Create custom Claims
type CustomClaims struct {
	Account string
	jwt.RegisteredClaims
}

// GenToken Generate JWT
func GenToken(userId string) (string, error) {
	claims := CustomClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "distri-ai", // 签发人
		},
	}
	// Create a signature object using the specified signature method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign using the specified secret and obtain the complete encoded string token
	return token.SignedString(CustomSecret)
}

// ParseToken parse JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// Perform type assertion on Claim in token object
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// Jwt Authorization middleware based JWT
func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "Authorization in header is null",
			})
			c.Abort()
			return
		}

		claims, err := ParseToken(authHeader)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "Invalid token",
			})
			c.Abort()
			return
		}

		c.Set("account", claims.Account)
		c.Next()
	}
}
