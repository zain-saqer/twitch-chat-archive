package chat

import (
	"time"
)

type Message struct {
	ID          string
	Username    string
	ChannelName string
	Message     string
	MessageType MessageType
	Time        time.Time
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
