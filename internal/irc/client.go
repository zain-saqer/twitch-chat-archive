package irc

import (
	"context"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chat-archive/internal/chat"
)

func NewMessagePipeline(client *twitch.Client) chat.GetMessageStream {
	return func(ctx context.Context, channels []string, messageTypes []chat.MessageType) (<-chan *chat.Message, error) {
		messageStream := make(chan *chat.Message)
		for _, channel := range channels {
			client.Join(channel)
		}
		clientDone := make(chan any)
		go func() {
			err := client.Connect()
			if err != nil {
				log.Err(err).Msg(`error while connect twitch irc client`)
			}
			clientDone <- struct{}{}
		}()
		client.OnPrivateMessage(func(message twitch.PrivateMessage) {
			select {
			case <-ctx.Done():
				return
			default:
				privateMessage := &chat.Message{
					ID:          message.ID,
					Username:    message.User.Name,
					ChannelName: message.Channel,
					Message:     message.Message,
					MessageType: mapToOurMessageType(message.Type),
					Time:        message.Time,
				}
				select {
				case <-ctx.Done():
					return
				case messageStream <- privateMessage:
				}
			}
		})
		go func() {
			defer close(messageStream)
			select {
			case <-ctx.Done():
			case <-clientDone:
				client.OnPrivateMessage(nil)
			}
		}()

		return messageStream, nil
	}
}

func mapToOurMessageType(messageType twitch.MessageType) chat.MessageType {
	switch messageType {
	case twitch.PRIVMSG:
		return chat.PrivMsg
	case twitch.CLEARCHAT:
		return chat.ClearChat
	case twitch.JOIN:
		return chat.Join
	case twitch.GLOBALUSERSTATE:
		return chat.GlobalUserState
	case twitch.NAMES:
		return chat.Names
	case twitch.PART:
		return chat.Part
	case twitch.PING:
		return chat.Ping
	case twitch.PONG:
		return chat.Pong
	case twitch.NOTICE:
		return chat.Notice
	case twitch.RECONNECT:
		return chat.Reconnect
	case twitch.ROOMSTATE:
		return chat.RoomState
	case twitch.USERNOTICE:
		return chat.UserNotice
	case twitch.WHISPER:
		return chat.Whisper
	case twitch.CLEARMSG:
		return chat.ClearMsg
	case twitch.USERSTATE:
		return chat.UserState
	case twitch.UNSET:
	default:
		return chat.Unset
	}
	return chat.Unset
}