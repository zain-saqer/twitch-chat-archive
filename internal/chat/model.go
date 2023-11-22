package chat

import (
	"context"
	"github.com/google/uuid"
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
	ID   uuid.UUID `ch:"id"`
	Name string    `ch:"name"`
	Time time.Time `ch:"timestamp"`
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
	GetChannel(ctx context.Context, uuid uuid.UUID) (*Channel, error)
	SaveChannel(ctx context.Context, channel *Channel) error
	DeleteChannel(ctx context.Context, channel *Channel) error
}
