package prom

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/florianloch/roomloggo/pkg/hw"
)

const namespace = "roomloggo"

type PrometheusExporter struct {
	tempReadings *prometheus.GaugeVec
	humReadings  *prometheus.GaugeVec
}

func NewPrometheusExporter() (*PrometheusExporter, error) {
	p := &PrometheusExporter{
		tempReadings: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "temperature_celsius",
			Help:      "Temperature measured in Celsius, labelled with the sensor name",
		}, []string{"sensor"}),
		humReadings: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "humidity_ratio",
			Help:      "Humidity given as percentage between 0 and 1, labelled with the sensor name",
		}, []string{"sensor"}),
	}

	if err := prometheus.Register(p); err != nil {
		return nil, fmt.Errorf("registering prometheus collector: %w", err)
	}

	return p, nil
}

func (p *PrometheusExporter) Process(readings []*hw.Reading) {
	for _, reading := range readings {
		p.humReadings.WithLabelValues(reading.Sensor).Set(float64(reading.Humidity) / 100.0)
		p.tempReadings.WithLabelValues(reading.Sensor).Set(reading.Temperature)
	}
}

func (p *PrometheusExporter) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(p, descs)
}

func (p *PrometheusExporter) Collect(metrics chan<- prometheus.Metric) {
	p.humReadings.Collect(metrics)
	p.tempReadings.Collect(metrics)
}
