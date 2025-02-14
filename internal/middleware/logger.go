package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestLogger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Info("Incoming request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("ip", c.RealIP()),
			)

			err := next(c)

			log.Info("Request completed",
				zap.Int("status", c.Response().Status),
			)
			return err
		}
	}
}
