package sensor

import (
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/florianloch/roomloggo/pkg/hw"
)

type IDToNameMapper struct {
	names []string
}

func NewIDToNameMapper(names []string) *IDToNameMapper {
	return &IDToNameMapper{
		names: names,
	}
}

func (i *IDToNameMapper) Process(readings []*hw.Reading) {
	for _, reading := range readings {
		if reading == nil {
			continue
		}

		idx, err := strconv.Atoi(reading.Sensor)
		if err != nil {
			log.Warn().Str("id", reading.Sensor).Msg("Invalid sensor ID, cannot map to name")

			continue
		}

		// Sensor IDs are derived from the channel used by the sensor, these are starting from 1
		idx--

		if idx < 0 || idx >= len(i.names) {
			continue
		}

		reading.Sensor = i.names[idx]
	}
}
