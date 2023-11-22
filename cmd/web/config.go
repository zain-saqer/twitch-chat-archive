package main

import (
	"github.com/zain-saqer/twitch-chat-archive/internal/chatlog"
	"os"
)

func getConfigs() *chatlog.Config {
	_, debug := os.LookupEnv(`DEBUG`)
	return &chatlog.Config{
		Debug:          debug,
		ServerAddress:  os.Getenv(`SERVER_ADDRESS`),
		ClickhouseDB:   os.Getenv(`CLICKHOUSE_DB`),
		ClickhouseHost: os.Getenv(`CLICKHOUSE_HOST`),
		ClickhousePort: os.Getenv(`CLICKHOUSE_PORT`),
		ClickhouseUser: os.Getenv(`CLICKHOUSE_USER`),
		ClickhousePass: os.Getenv(`CLICKHOUSE_PASS`),
		AuthUser:       os.Getenv(`AUTH_USER`),
		AuthPass:       os.Getenv(`AUTH_PASS`),
	}
}
