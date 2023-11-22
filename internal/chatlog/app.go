package chatlog

import (
	"context"
	twitchirc "github.com/gempir/go-twitch-irc/v4"
	"github.com/zain-saqer/twitch-chat-archive/internal/chat"
	"github.com/zain-saqer/twitch-chat-archive/internal/irc"
)

type Config struct {
	Debug          bool
	AuthUser       string
	AuthPass       string
	ServerAddress  string
	ClickhouseDB   string
	ClickhouseHost string
	ClickhousePort string
	ClickhouseUser string
	ClickhousePass string
}

type App struct {
	ChatRepository chat.Repository
	Config         *Config
	TwitchClient   *twitchirc.Client
}

func (a *App) JoinChannel(channel ...string) {
	a.TwitchClient.Join(channel...)
}

func (a *App) Depart(channel string) {
	a.TwitchClient.Depart(channel)
}

func (a *App) StartMessagePipeline(ctx context.Context) error {
	channels, err := a.ChatRepository.GetChannels(ctx)
	if err != nil {
		return err
	}
	for _, channel := range channels {
		a.JoinChannel(channel.Name)
	}
	messageTypes := []uint8{chat.PrivMsg}
	messageStream, err := irc.NewMessagePipeline(a.TwitchClient)(ctx, messageTypes)
	if err != nil {
		return err
	}
	filteredMessageStream := chat.FilterMessageStream(ctx, messageStream, messageTypes)
	chat.SaveMessageStream(ctx, filteredMessageStream, a.ChatRepository)
	return nil
}
