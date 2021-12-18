package internal

import "github.com/florianloch/roomloggo/pkg/hw"

type MeasurementsProcessor interface {
	// Process is not allowed to operate on the referenced readings anymore after control flow has been returned
	// to the caller.
	Process(readings []*hw.Reading)
}

type ProcessFn func([]*hw.Reading)

func (f ProcessFn) Process(readings []*hw.Reading) {
	f(readings)
}
