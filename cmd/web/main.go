package main

import (
	"context"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chat-archive/internal/clickhouse"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	config := getConfigs()
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
	app := &App{
		ChatRepository: chatRepository,
		Config:         config,
		TwitchClient:   twitchIrcClient,
	}
	if err := app.StartMessagePipeline(ctx); err != nil {
		log.Fatal().Err(err).Stack().Msg(`error starting the message pipeline`)
	}
	channels := []string{`Smoke`, `summit1g`, `tarik`, `kaicenat`, `jynxzi`, `caseoh_`, `maximum`, `mizkif`, `casimito`, `xQc`, `montanablack88`, `anomaly`, `soursweet`, `ohnepixel`, `psp1g`, `handongsuk`, `raderaderader`}
	for _, channel := range channels {
		app.TwitchClient.Join(channel)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg(`shutting down...`)
				return
			}
		}
	}()

	wg.Wait()
}
