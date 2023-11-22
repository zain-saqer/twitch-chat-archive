package chat

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID          string    `ch:"id"`
	Username    string    `ch:"username"`
	ChannelName string    `ch:"channel"`
	Message     string    `ch:"message"`
	MessageType uint8     `ch:"message_type"`
	Time        time.Time `ch:"timestamp"`
}

type Channel struct {
	ID   uuid.UUID `ch:"id"`
	Name string    `ch:"name"`
	Time time.Time `ch:"timestamp"`
}

const (
	Unset           uint8 = iota
	Whisper         uint8 = iota
	PrivMsg         uint8 = iota
	ClearChat       uint8 = iota
	RoomState       uint8 = iota
	UserNotice      uint8 = iota
	UserState       uint8 = iota
	Notice          uint8 = iota
	Join            uint8 = iota
	Part            uint8 = iota
	Reconnect       uint8 = iota
	Names           uint8 = iota
	Ping            uint8 = iota
	Pong            uint8 = iota
	ClearMsg        uint8 = iota
	GlobalUserState uint8 = iota
)

type Repository interface {
	SaveMessage(ctx context.Context, message *Message) error
	GetMessages(ctx context.Context, channel, username string, limit, offset int) ([]*Message, error)
	GetChannels(ctx context.Context) ([]*Channel, error)
	GetChannel(ctx context.Context, uuid uuid.UUID) (*Channel, error)
	SaveChannel(ctx context.Context, channel *Channel) error
	DeleteChannel(ctx context.Context, channel *Channel) error
}
