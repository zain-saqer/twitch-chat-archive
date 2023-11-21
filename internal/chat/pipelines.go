package chat

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"slices"
	"time"
)

type GetMessageStream func(ctx context.Context, channels []string, messageTypes []MessageType) (<-chan *Message, error)

func FilterMessages(ctx context.Context, messageStream <-chan *Message, allowedTypes []MessageType) <-chan *Message {
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

func BufferMessages(ctx context.Context, messageStream <-chan *Message, interval time.Duration) <-chan []*Message {
	bufferedMessageStream := make(chan []*Message)
	bufferStream := make(chan *Message, 50000)
	go func() {
		defer close(bufferedMessageStream)
		for {
			select {
			case <-ctx.Done():
				return
			case message, ok := <-messageStream:
				if !ok {
					return
				}
				bufferStream <- message
			}
		}
	}()

	go func() {
		for range time.Tick(interval) {
			buffer := make([]*Message, 0)
		loop:
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.Tick(time.Millisecond):
					break loop
				case message := <-bufferStream:
					buffer = append(buffer, message)
				}
			}
			go func() {
				select {
				case <-ctx.Done():
					return
				case bufferedMessageStream <- buffer:
				}
			}()
		}
	}()
	return bufferedMessageStream
}

func SaveMessages(ctx context.Context, messagesStream <-chan []*Message, repository LogRepository) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case messages := <-messagesStream:
				fmt.Printf("saving %d messages\n", len(messages))
				now := time.Now()
				err := repository.SaveAll(ctx, messages)
				elapsed := time.Since(now)
				if err != nil {
					log.Error().Err(err).Msg(`error while saving message`)
				}
				fmt.Printf("time elapsed: %s\n", elapsed)
			}
		}
	}()
}

func CreateMessagePipeline(ctx context.Context, getMessageStream GetMessageStream, channels []string, messageTypes []MessageType, repository LogRepository) error {
	messageStream, err := getMessageStream(ctx, channels, messageTypes)
	if err != nil {
		return err
	}
	filteredMessageStream := FilterMessages(ctx, messageStream, messageTypes)
	bufferedMessageStream := BufferMessages(ctx, filteredMessageStream, 3*time.Second)
	SaveMessages(ctx, bufferedMessageStream, repository)
	return nil
}
