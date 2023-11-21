package main

import (
	"context"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chat-archive/internal/chat"
	"github.com/zain-saqer/twitch-chat-archive/internal/irc"
	"github.com/zain-saqer/twitch-chat-archive/internal/mongo"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	config := getConfigs()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	twitchIrcClient := twitch.NewAnonymousClient()
	channels := []string{`Smoke`, `summit1g`, `tarik`, `kaicenat`, `jynxzi`, `caseoh_`, `maximum`, `mizkif`, `casimito`, `xQc`, `montanablack88`}
	messageTypes := []chat.MessageType{chat.PrivMsg}
	mongoClient, err := mongo.NewClient(ctx, config.MongoHost, config.MongoPort, config.MongoUsername, config.MongoPassword, 3*time.Second)
	if err != nil {
		log.Fatal().Err(err).Stack().Msg(`error creating mongodb twitchIrcClient`)
	}
	err = mongo.PrepareDatabase(ctx, mongoClient, config.MongoDatabase, config.MongoCollection)
	if err != nil {
		log.Fatal().Err(err).Stack().Msg(`error preparing the database`)
	}
	logRepository := mongo.NewMongoLogRepository(mongoClient, config.MongoDatabase, config.MongoCollection)
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
