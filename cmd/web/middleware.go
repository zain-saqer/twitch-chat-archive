package main

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) middlewares() {
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.Logger())
}

func (s *Server) authMiddleware() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(s.App.Config.AuthUser)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(s.App.Config.AuthPass)) == 1 {
			return true, nil
		}
		return false, nil
	})
}
