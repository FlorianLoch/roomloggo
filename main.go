package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/florianloch/roomloggo/internal"
	"github.com/florianloch/roomloggo/internal/config"
	"github.com/florianloch/roomloggo/internal/influx"
	"github.com/florianloch/roomloggo/internal/prom"
	"github.com/florianloch/roomloggo/internal/sensor"
	"github.com/florianloch/roomloggo/pkg/hw"
)

const envKeyConfig = "ROOMLOGGO_CONFIG"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if !hw.IsHidSupported() {
		log.Fatal().Msg("HID library does not support the current platform. Exiting.")
	}

	cli := &struct {
		Config config.Config `embed:""`
	}{}

	kong.ConfigureHelp(kong.HelpOptions{Compact: false, Summary: true})

	kong.Name("roomloggo")
	kong.Description("Small daemon reading temperature and humidity from a dnt RoomLogg Pro base station via USB and pushing it into an InfluxDB.")

	configPath := os.Getenv(envKeyConfig)

	if configPath == "" {
		configPath = "./roomloggo.config.yaml"
	}

	ctx := kong.Parse(cli, kong.Configuration(kongyaml.Loader, configPath))
	if ctx.Error != nil {
		log.Fatal().Err(ctx.Error).Msg("Failed to parse input parameters/commands")
	}

	cfg := &cli.Config

	timeoutCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	processors := []internal.MeasurementsProcessor{
		sensor.NewIDToNameMapper(cfg.Sensor.Names),
		internal.ProcessFn(sensor.LogReadings),
	}

	if cfg.PromExporter.Enabled {
		promExporter, err := prom.NewPrometheusExporter()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to setup Prometheus exporter")
		}

		go func() {
			http.Handle("/metrics", promhttp.Handler())

			log.Info().Msgf("Serving metrics at: http://%s/metrics", cfg.PromExporter.ListenAddr)
			if err := http.ListenAndServe(cfg.PromExporter.ListenAddr, http.DefaultServeMux); err != nil {
				log.Fatal().Err(err).Msg("HTTP server of Prometheus exporter crashed")
			}
		}()

		processors = append(processors, promExporter)
	}

	if cfg.Influx.Enabled {
		influxClient, err := influx.New(timeoutCtx, cfg.Influx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to setup client for InfluxDB")
		}

		log.Info().Msgf("Will push measurements to InfluxDB at: %s", cfg.Influx.BaseURL)

		processors = append(processors, influxClient)
	}

	log.Info().Msg("Setup completed, starting to log data...")

	sensor.StartLoop(cfg.Sensor.Interval, processors...)
}
