package config

import "time"

type Config struct {
	Influx       InfluxConfig       `embed:"" prefix:"influx-db-"`
	Sensor       SensorConfig       `embed:"" prefix:"sensor-"`
	PromExporter PromExporterConfig `embed:"" prefix:"prom-exporter-"`
}

type SensorConfig struct {
	Interval time.Duration `default:"1m" help:"Interval to check source dir for changes"`
	Names    []string      `help:"Set-up mappings from sensor IDs to sensor names"`
}

type InfluxConfig struct {
	Enabled      bool          `default:"true" help:"Turn writing to InfluxDB on or off"`
	APIToken     string        `name:"api-token" help:"API token for InfluxDB"`
	BaseURL      string        `name:"base-url" help:"Base URL of InfluxDB to push data to"`
	Org          string        `help:"Org to which the target bucket belongs"`
	Bucket       string        `help:"Bucket to which readings shall be written"`
	WriteTimeout time.Duration `default:"10s" help:"Timeout for writes towards InfluxDB"`
}

type PromExporterConfig struct {
	Enabled    bool   `default:"true" help:"Turn exporting of prometheus metrics on or off"`
	ListenAddr string `default:"localhost:9044" help:"Interface and port for Prometheus metrics exporter to listen on"`
}
