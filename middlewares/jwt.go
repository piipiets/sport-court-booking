package middlewares

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/piipiets/sport-court-booking/helpers/common"
	"github.com/spf13/viper"
)

type Claims struct {
	jwt.RegisteredClaims
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := GetJwtTokenFromHeader(c)
		if err != nil {
			common.GenerateErrorResponse(c, err.Error())
			c.Abort()
			return
		}

		// Parse dan validasi JWT
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Pastikan menggunakan algoritma HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(viper.GetString("jwt_secret_key")), nil
		},
		)

		if err != nil {
			common.GenerateErrorResponse(c, "token invalid or expired")
			c.Abort()
			return
		}

		if !token.Valid {
			common.GenerateErrorResponse(c, "token invalid")
			c.Abort()
			return
		}

		// Ambil claims dari token
		claims, ok := token.Claims.(*Claims)
		if !ok {
			common.GenerateErrorResponse(c, "invalid token claims")
			c.Abort()
			return
		}

		// Simpan claims ke Gin Context
		c.Set("auth", claims)

		c.Next()
	}
}

func GetJwtTokenFromHeader(c *gin.Context) (tokenString string, err error) {
	authHeader := c.Request.Header.Get("Authorization")

	if common.IsEmptyField(authHeader) {
		return tokenString, errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return tokenString, errors.New("invalid Authorization header format")
	}

	return parts[1], nil
}

func GenerateJwtToken() (token string, err error) {
	expirationTime := time.Now().Add(1 * time.Minute)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	generatedTokenJwt := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	token, err = generatedTokenJwt.SignedString(
		[]byte(viper.GetString("jwt_secret_key")),
	)

	if err != nil {
		return
	}

	return
}
