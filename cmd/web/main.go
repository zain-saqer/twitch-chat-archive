package main

import (
	"context"
	"errors"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chat-archive/internal/chatlog"
	"github.com/zain-saqer/twitch-chat-archive/internal/clickhouse"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	config := getConfigs()
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	twitchIrcClient := twitch.NewAnonymousClient()

	conn, err := clickhouse.NewConnection(ctx, config.ClickhouseHost, config.ClickhousePort, config.ClickhouseDB, config.ClickhouseUser, config.ClickhousePass)
	if err != nil {
		log.Fatal().Err(err).Stack().Msg(`error creating clickhouse connection`)
	}
	chatRepository := clickhouse.NewRepository(conn, config.ClickhouseDB)
	if err := chatRepository.PrepareDatabase(ctx); err != nil {
		log.Fatal().Err(err).Stack().Msg(`error while preparing clickhouse database`)
	}
	app := &chatlog.App{
		ChatRepository: chatRepository,
		Config:         config,
		TwitchClient:   twitchIrcClient,
	}
	if err := app.StartMessagePipeline(ctx); err != nil {
		log.Fatal().Err(err).Stack().Msg(`error starting the message pipeline`)
	}
	e := echo.New()
	e.Debug = config.Debug
	server := NewServer(app, e)
	server.middlewares()
	server.setupRoutes()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := e.Start(config.ServerAddress)
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatal().Err(err).Msg(`shutting down server error`)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg(`shutting down...`)
				if err := e.Shutdown(ctx); err != nil {
					log.Error().Err(err).Msg(`error while shutting down the web server`)
				}
				return
			}
		}
	}()

	wg.Wait()
}
