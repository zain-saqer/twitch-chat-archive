package chat

import (
	"context"
	"time"
)

type Message struct {
	ID          string      `json:"id"`
	Username    string      `json:"username"`
	ChannelName string      `json:"channelName"`
	Message     string      `json:"message"`
	MessageType MessageType `json:"messageType"`
	Time        time.Time   `json:"timestamp"`
}

type Channel struct {
	Name string    `json:"name"`
	Time time.Time `json:"timestamp"`
}

type MessageType int

const (
	Unset           MessageType = iota
	Whisper         MessageType = iota
	PrivMsg         MessageType = iota
	ClearChat       MessageType = iota
	RoomState       MessageType = iota
	UserNotice      MessageType = iota
	UserState       MessageType = iota
	Notice          MessageType = iota
	Join            MessageType = iota
	Part            MessageType = iota
	Reconnect       MessageType = iota
	Names           MessageType = iota
	Ping            MessageType = iota
	Pong            MessageType = iota
	ClearMsg        MessageType = iota
	GlobalUserState MessageType = iota
)

type Repository interface {
	SaveMessage(ctx context.Context, message *Message) error
	GetChannels(ctx context.Context) ([]*Channel, error)
}
