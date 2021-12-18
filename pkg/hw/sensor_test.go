package hw

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_fromBytes(t *testing.T) {
	input, err := hex.DecodeString("7b00a03e00a54000d03400b53d7fffff00254b7fffff7fffff")
	require.NoError(t, err)

	readings, err := fromBytes(input)
	require.NoError(t, err)

	require.Len(t, readings, 5)

	expectedReadings := []*Reading{
		{
			Sensor:      "1",
			Temperature: 16,
			Humidity:    62,
		},
		{
			Sensor:      "2",
			Temperature: 16.5,
			Humidity:    64,
		},
		{
			Sensor:      "3",
			Temperature: 20.8,
			Humidity:    52,
		},
		{
			Sensor:      "4",
			Temperature: 18.1,
			Humidity:    61,
		},
		// Sensor 5 is not active
		{
			Sensor:      "6",
			Temperature: 3.7,
			Humidity:    75,
		},
	}

	for i := range readings {
		readings[i].Time = time.Time{} // "unset the time" in order to make the structs comparable
	}

	require.Equal(t, expectedReadings, readings)
}
