package config

import "time"

type Config struct {
	Influx InfluxConfig `embed:"" prefix:"influx-db-"`
	Sensor SensorConfig `embed:"" prefix:"sensor-"`
}

type SensorConfig struct {
	Interval time.Duration `default:"1m" help:"Interval to check source dir for changes"`
	Names    []string      `help:"Set-up mappings from sensor IDs to sensor names"`
}

type InfluxConfig struct {
	APIToken     string        `name:"api-token" required:"" help:"API token for InfluxDB"`
	BaseURL      string        `name:"base-url" required:"" help:"Base URL of InfluxDB to push data to"`
	Org          string        `required:"" help:"Org to which the target bucket belongs"`
	Bucket       string        `required:"" help:"Bucket to which readings shall be written"`
	WriteTimeout time.Duration `default:"10s" help:"Timeout for writes towards InfluxDB"`
}
