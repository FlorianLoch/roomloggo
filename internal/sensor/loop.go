package sensor

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/florianloch/roomloggo/internal"
	"github.com/florianloch/roomloggo/pkg/hw"
)

func StartLoop(interval time.Duration, downstream ...internal.MeasurementsProcessor) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		if readings, err := hw.Read(); err != nil {
			log.Error().Err(err).Msg("Failed to read data from station")
		} else {
			// Run this asynchronously, just as a defensive measure in order to avoid a varying reading frequency
			go func() {
				for _, u := range downstream {
					u.Process(readings)
				}
			}()
		}

		<-t.C
	}
}

func LogReadings(readings []*hw.Reading) {
	for _, r := range readings {
		log.Info().Float64("temperature", r.Temperature).Int8("humidity", r.Humidity).Msgf("%s:", r.Sensor)
	}
}
