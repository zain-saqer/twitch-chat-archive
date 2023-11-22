package main

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chat-archive/internal/chat"
	"github.com/zain-saqer/twitch-chat-archive/web"
	"html/template"
	"net/http"
	"sync"
	"time"
)

func (s *Server) setupRoutes() {
	adminGroup := s.Echo.Group(`/admin`, s.authMiddleware())
	adminGroup.GET(``, s.getAdminIndex)
	adminGroup.GET(`/channels`, s.getAdminChannels)
	adminGroup.GET(`/channels/add`, s.getAdminAddChannel)
	adminGroup.POST(`/channels/add`, s.postAdminAddChannel)
	adminGroup.DELETE(`/channels/:id`, s.deleteAdminDeleteChannel)
}

func (s *Server) getAdminIndex(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, `/admin/channels`)
}

func (s *Server) getAdminChannels(c echo.Context) error {
	var t *template.Template
	sync.OnceFunc(func() {
		var err error
		t, err = template.ParseFS(web.F, `templates/layout.gohtml`, `templates/nav.gohtml`, `templates/admin/channels.gohtml`)
		if err != nil {
			log.Fatal().Err(err).Stack().Msg(`error parsing templates`)
		}
	})()
	channels, err := s.App.ChatRepository.GetChannels(c.Request().Context())
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(c.Response(), `base`, channels)
}

func (s *Server) getAdminAddChannel(c echo.Context) error {
	var t *template.Template
	sync.OnceFunc(func() {
		var err error
		t, err = template.ParseFS(web.F, `templates/layout.gohtml`, `templates/nav.gohtml`, `templates/admin/add_channel.gohtml`)
		if err != nil {
			log.Fatal().Err(err).Stack().Msg(`error parsing templates`)
		}
	})()
	return t.ExecuteTemplate(c.Response(), `base`, AddChannel{})
}

func (s *Server) postAdminAddChannel(c echo.Context) error {
	var t *template.Template
	sync.OnceFunc(func() {
		var err error
		t, err = template.ParseFS(web.F, `templates/layout.gohtml`, `templates/nav.gohtml`, `templates/admin/add_channel.gohtml`)
		if err != nil {
			log.Fatal().Err(err).Stack().Msg(`error parsing templates`)
		}
	})()
	addChannel := &AddChannel{}
	err := c.Bind(addChannel)
	if err != nil {
		return err
	}
	addChannel.Trim()
	if !addChannel.Validate() {
		return t.ExecuteTemplate(c.Response(), `base`, addChannel)
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	if err = s.App.ChatRepository.SaveChannel(c.Request().Context(), &chat.Channel{ID: id, Name: addChannel.Name, Time: time.Now()}); err != nil {
		return err
	}
	s.App.JoinChannel(addChannel.Name)
	return c.Redirect(http.StatusSeeOther, `/admin/channels`)
}

func (s *Server) deleteAdminDeleteChannel(c echo.Context) error {
	deleteChannel := &DeleteChannel{}
	err := c.Bind(deleteChannel)
	if err != nil {
		return err
	}
	if !deleteChannel.Validate() {
		return errors.New(`invalid request`)
	}
	channel, err := s.App.ChatRepository.GetChannel(c.Request().Context(), deleteChannel.UUID)
	if err != nil {
		return err
	}
	if err = s.App.ChatRepository.DeleteChannel(c.Request().Context(), channel); err != nil {
		return err
	}
	s.App.Depart(channel.Name)
	return c.HTML(200, ``)
}
