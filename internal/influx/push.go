package influx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/rs/zerolog/log"

	"github.com/florianloch/roomloggo/internal/config"
	"github.com/florianloch/roomloggo/pkg/hw"
)

var ErrInfluxDBNotReachable = errors.New("could not ping InfluxDB; seems down")

type InfluxPush struct {
	writer       api.WriteAPIBlocking
	writeTimeout time.Duration
}

func New(ctx context.Context, cfg config.InfluxConfig) (*InfluxPush, error) {
	client := influxdb2.NewClient(cfg.BaseURL, cfg.APIToken)

	ok, err := client.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging InfluxDB: %w", err)
	}
	if !ok {
		return nil, ErrInfluxDBNotReachable
	}

	return &InfluxPush{
		writer:       client.WriteAPIBlocking(cfg.Org, cfg.Bucket),
		writeTimeout: cfg.WriteTimeout,
	}, nil
}

func (i *InfluxPush) Process(readings []*hw.Reading) {
	timeoutCtx, cancelFn := context.WithTimeout(context.Background(), i.writeTimeout)
	defer cancelFn()

	points := make([]*write.Point, len(readings))

	for i, reading := range readings {
		points[i] = influxdb2.NewPointWithMeasurement("climate").
			AddField("channel", reading.Channel).
			AddField("temperature", reading.Temperature).
			AddField("humidity", reading.Humidity).
			SetTime(reading.Time)
	}

	if err := i.writer.WritePoint(timeoutCtx, points...); err != nil {
		log.Error().Err(err).Msg("Failed to write readings to InfluxDB")
	}
}
