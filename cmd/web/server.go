package main

import (
	"github.com/labstack/echo/v4"
	"github.com/zain-saqer/twitch-chat-archive/internal/chatlog"
)

type Server struct {
	App    *chatlog.App
	Echo   *echo.Echo
	Config *Config
}

func NewServer(app *chatlog.App, e *echo.Echo, config *Config) *Server {
	return &Server{App: app, Echo: e, Config: config}
}
