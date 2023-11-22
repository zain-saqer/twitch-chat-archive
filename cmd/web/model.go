package main

import (
	"github.com/google/uuid"
	"github.com/zain-saqer/twitch-chat-archive/internal/chat"
	"strings"
)

type IndexView struct {
	Channels []*chat.Channel
}

type ChatLogView struct {
	Messages []*chat.Message
	Channel  string
	Username string
	Offset   int
	Limit    int
}

type GetChatLogParam struct {
	Channel  string `query:"channel"`
	Username string `query:"username"`
	Offset   int    `query:"offset"`
}

func (p GetChatLogParam) Validate() bool {
	return p.Offset >= 0
}

type AddChannel struct {
	Errors []string
	Name   string `form:"name"`
}

func (c *AddChannel) Trim() {
	c.Name = strings.TrimSpace(c.Name)
}

func (c *AddChannel) Validate() bool {
	errors := make([]string, 0)
	if c.Name == "" {
		errors = append(errors, "Name is required")
	}
	c.Errors = errors
	return len(errors) == 0
}

type DeleteChannel struct {
	ID   string `param:"id"`
	UUID uuid.UUID
}

func (c *DeleteChannel) Validate() bool {
	id, err := uuid.Parse(c.ID)
	if err != nil {
		return false
	}
	c.UUID = id
	return true
}
