package main

import (
	"github.com/google/uuid"
	"strings"
)

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
