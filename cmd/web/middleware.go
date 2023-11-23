package main

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chat-archive/internal/chatlog"
)

func (s *Server) middlewares() {
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(loggerMiddleware())
}

func loggerMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("method", v.Method).
				Str("host", v.Host).
				Str("remote_ip", v.RemoteIP).
				Str("user_agent", v.UserAgent).
				Dur("latency", v.Latency).
				Str("latency_human", v.Latency.String()).
				Msg("request")

			return nil
		},
	})
}

func authMiddleware(config *chatlog.Config) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(config.AuthUser)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(config.AuthPass)) == 1 {
			return true, nil
		}
		return false, nil
	})
}