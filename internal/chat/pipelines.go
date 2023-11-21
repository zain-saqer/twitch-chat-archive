package chat

import (
	"context"
	"github.com/rs/zerolog/log"
	"slices"
)

type GetMessageStream func(ctx context.Context, messageTypes []MessageType) (<-chan *Message, error)

func FilterMessageStream(ctx context.Context, messageStream <-chan *Message, allowedTypes []MessageType) <-chan *Message {
	filteredMessageStream := make(chan *Message)

	go func() {
		defer close(filteredMessageStream)
		for {
			select {
			case <-ctx.Done():
				return
			case message, ok := <-messageStream:
				if !ok {
					return
				}
				if slices.Contains(allowedTypes, message.MessageType) {
					filteredMessageStream <- message
				}
			}
		}
	}()

	return filteredMessageStream
}

func SaveMessageStream(ctx context.Context, messagesStream <-chan *Message, repository Repository) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case messages := <-messagesStream:
				if err := repository.SaveMessage(ctx, messages); err != nil {
					log.Error().Err(err).Msg(`error while saving a message`)
				}
			}
		}
	}()
}
