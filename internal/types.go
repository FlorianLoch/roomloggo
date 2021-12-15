package internal

import "github.com/florianloch/roomloggo/pkg/hw"

type MeasurementProcessor interface {
	Process(readings []*hw.Reading)
}
