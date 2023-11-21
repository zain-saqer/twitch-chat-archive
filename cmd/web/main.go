package main

import (
	"context"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chat-archive/internal/chat"
	"github.com/zain-saqer/twitch-chat-archive/internal/clickhouse"
	"github.com/zain-saqer/twitch-chat-archive/internal/irc"
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
	channels := []string{`Smoke`, `summit1g`, `tarik`, `kaicenat`, `jynxzi`, `caseoh_`, `maximum`, `mizkif`, `casimito`, `xQc`, `montanablack88`, `anomaly`, `soursweet`, `ohnepixel`, `psp1g`, `handongsuk`, `raderaderader`}
	messageTypes := []chat.MessageType{chat.PrivMsg}
	conn, err := clickhouse.NewConnection(ctx, config.ClickhouseHost, config.ClickhousePort, config.ClickhouseDB, config.ClickhouseUser, config.ClickhousePass)
	if err != nil {
		log.Fatal().Err(err).Stack().Msg(`error creating clickhouse connection`)
	}
	logRepository := clickhouse.NewRepository(conn, config.ClickhouseDB)
	if err := logRepository.PrepareDatabase(ctx); err != nil {
		log.Fatal().Err(err).Stack().Msg(`error while preparing clickhouse database`)
	}
	err = chat.CreateMessagePipeline(ctx, irc.NewMessagePipeline(twitchIrcClient), channels, messageTypes, logRepository)
	if err != nil {
		log.Fatal().Err(err).Stack().Msg(`error creating message pipeline`)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
	MsgLoop:
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg(`shutting down...`)
				break MsgLoop
			}
		}
	}()

	wg.Wait()
}
