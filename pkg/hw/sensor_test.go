package hw

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_fromBytes(t *testing.T) {
	input, err := hex.DecodeString("7b00ac4200a03b00c13500b43f7ffffffffe5a7fffff7fffff")
	require.NoError(t, err)

	readings, err := fromBytes(input)
	require.NoError(t, err)

	require.Len(t, readings, 5)

	expectedReadings := []*Reading{
		{
			Sensor:      "1",
			Temperature: 17.2,
			Humidity:    66,
		},
		{
			Sensor:      "2",
			Temperature: 16,
			Humidity:    59,
		},
		{
			Sensor:      "3",
			Temperature: 19.3,
			Humidity:    53,
		},
		{
			Sensor:      "4",
			Temperature: 18,
			Humidity:    63,
		},
		// Sensor 5 is not active
		{
			Sensor:      "6",
			Temperature: -0.2,
			Humidity:    90,
		},
	}

	for i := range readings {
		readings[i].Time = time.Time{} // "unset the time" in order to make the structs comparable
	}

	require.Equal(t, expectedReadings, readings)
}
