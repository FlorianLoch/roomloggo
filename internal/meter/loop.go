package meter

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/florianloch/roomloggo/internal"
	"github.com/florianloch/roomloggo/pkg/hw"
)

func StartLoop(upstream internal.MeasurementProcessor, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		if readings, err := hw.Read(); err != nil {
			log.Error().Err(err).Msg("Failed to read data from station")
		} else {
			// Run this asynchronously, just as a defensive measure in order to avoid a varying reading frequency
			go upstream.Process(readings)
		}

		<-t.C
	}
}
