package main

import "os"

type Config struct {
	ClickhouseDB   string
	ClickhouseHost string
	ClickhousePort string
	ClickhouseUser string
	ClickhousePass string
}

func getConfigs() *Config {
	return &Config{
		ClickhouseDB:   os.Getenv(`CLICKHOUSE_DB`),
		ClickhouseHost: os.Getenv(`CLICKHOUSE_HOST`),
		ClickhousePort: os.Getenv(`CLICKHOUSE_PORT`),
		ClickhouseUser: os.Getenv(`CLICKHOUSE_USER`),
		ClickhousePass: os.Getenv(`CLICKHOUSE_PASS`),
	}
}
