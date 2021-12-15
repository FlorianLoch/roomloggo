package config

import "time"

type Config struct {
	Influx       InfluxConfig  `embed:"" prefix:"influx-db-"`
	ReadInterval time.Duration `yaml:"interval" default:"5m" help:"Interval to check source dir for changes."`
}

type InfluxConfig struct {
	APIToken     string        `yaml:"api-token" required:"" help:"API token for InfluxDB"`
	BaseURL      string        `yaml:"base-url" required:"" help:"Base URL of InfluxDB to push data to"`
	Org          string        `required:"" help:"Org to which the target bucket belongs"`
	Bucket       string        `required:"" help:"Bucket to which readings shall be written"`
	WriteTimeout time.Duration `default:"10s" help:"Timeout for writes towards InfluxDB"`
}
