package main

import (
	"os"
)

type Config struct {
	MongoHost       string
	MongoUsername   string
	MongoPassword   string
	MongoPort       string
	MongoDatabase   string
	MongoCollection string
}

func getConfigs() *Config {
	return &Config{
		MongoHost:       os.Getenv(`MONGO_HOST`),
		MongoUsername:   os.Getenv(`MONGO_USERNAME`),
		MongoPassword:   os.Getenv(`MONGO_PASSWORD`),
		MongoPort:       os.Getenv(`MONGO_PORT`),
		MongoDatabase:   os.Getenv(`MONGO_DATABASE`),
		MongoCollection: os.Getenv(`MONGO_COLLECTION`),
	}
}
