package hw

import (
	"encoding/hex"
	"testing"

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
			Channel:     1,
			Temperature: 16,
			Humidity:    62,
			Present:     true,
		},
		{
			Channel:     2,
			Temperature: 16.5,
			Humidity:    64,
			Present:     true,
		},
		{
			Channel:     3,
			Temperature: 20.8,
			Humidity:    52,
			Present:     true,
		},
		{
			Channel:     4,
			Temperature: 18.1,
			Humidity:    61,
			Present:     true,
		},
		{
			Channel:     6,
			Temperature: 3.7,
			Humidity:    75,
			Present:     true,
		},
	}

	require.Equal(t, expectedReadings, readings)
}
