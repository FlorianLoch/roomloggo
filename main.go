package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/florianloch/roomloggo/pkg/hw"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	readings, err := hw.Read()
	if err != nil {
		log.Fatal().Err(err).Msg("Process sensor data failed")
	}

	log.Info().Interface("readings", readings).Msg("Sensor data successfully processed")
}
