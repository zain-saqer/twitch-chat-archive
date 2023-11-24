package env

import (
	"errors"
	"github.com/rs/zerolog/log"
	"os"
)

var ErrNotExist = errors.New(`environment variable don't exist'`)

func GetEnv(name string) (string, error) {
	val, exist := os.LookupEnv(name)
	if !exist {
		return "", ErrNotExist
	}
	return val, nil
}

func MustGetEnv(name string) string {
	val, err := GetEnv(name)
	if err != nil {
		log.Fatal().Err(err).Msgf(`error getting env var: %s`, name)
	}
	return val
}
