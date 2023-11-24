package main

import (
	"github.com/zain-saqer/twitch-chat-archive/internal/env"
	"os"
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
	SentryDsn      string
}

func getConfigs() *Config {
	_, debug := os.LookupEnv(`DEBUG`)
	return &Config{
		Debug:          debug,
		ServerAddress:  env.MustGetEnv(`SERVER_ADDRESS`),
		ClickhouseDB:   env.MustGetEnv(`CLICKHOUSE_DB`),
		ClickhouseHost: env.MustGetEnv(`CLICKHOUSE_HOST`),
		ClickhousePort: env.MustGetEnv(`CLICKHOUSE_PORT`),
		ClickhouseUser: env.MustGetEnv(`CLICKHOUSE_USER`),
		ClickhousePass: env.MustGetEnv(`CLICKHOUSE_PASS`),
		AuthUser:       env.MustGetEnv(`AUTH_USER`),
		AuthPass:       env.MustGetEnv(`AUTH_PASS`),
		SentryDsn:      env.MustGetEnv(`SENTRY_DSN`),
	}
}
