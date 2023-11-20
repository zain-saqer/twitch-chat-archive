package chat

import (
	"context"
	"slices"
)

type GetMessageStream func(ctx context.Context, channels []string, messageTypes []MessageType) (<-chan *Message, error)

func FilterMessages(ctx context.Context, messageStream <-chan *Message, allowedTypes []MessageType) (<-chan *Message, error) {
	filteredMessageStream := make(chan *Message)

	go func() {
		defer close(filteredMessageStream)
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-messageStream:
				if slices.Contains(allowedTypes, message.MessageType) {
					filteredMessageStream <- message
				}
			}
		}
	}()

	return filteredMessageStream, nil
}

func CreateMessageStreamPipeline(ctx context.Context, getMessageStream GetMessageStream, channels []string, messageTypes []MessageType) (<-chan *Message, error) {
	messageStream, err := getMessageStream(ctx, channels, messageTypes)
	if err != nil {
		return nil, err
	}
	filteredMessageStream, err := FilterMessages(ctx, messageStream, messageTypes)
	if err != nil {
		return nil, err
	}

	return filteredMessageStream, nil
}
