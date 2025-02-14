package middleware

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/api/auth" && c.Request().Method == "POST" {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, map[string]string{"error": "missing token"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return c.JSON(401, map[string]string{"error": "invalid token format"})
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(401, map[string]string{"error": "invalid token"})
			}

			claims := token.Claims.(jwt.MapClaims)
			c.Set("user_id", claims["user_id"])

			return next(c)
		}
	}
}
