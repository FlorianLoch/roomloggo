package hw

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/karalabe/hid"
	"github.com/rs/zerolog/log"
)

const (
	vendorID = 0x483
	deviceID = 0x5750
)

var (
	ErrNoDeviceFound          = errors.New("no matching device found")
	ErrByteSliceInvalidLength = errors.New("length of slice invalid, expected blocks of 3 bytes plus header of length 1")
	ErrBadFirstByte           = errors.New("expected first byte to be 0x7b")
)

var readRequestBytes = func() []byte {
	b := make([]byte, 64)
	b[0] = 0x7B
	b[1] = 0x03
	b[2] = 0x40
	b[3] = 0x7D

	return b
}()

type Reading struct {
	Sensor string
	// float32 would be sufficient, but in most cases float64 will be needed eventually; as the conversion between
	// float32 and float64 is not really smooth we go with float64 straight away
	Temperature float64
	Humidity    int8
	Time        time.Time
}

func fromBytes(raw []byte) ([]*Reading, error) {
	if len(raw)%3 != 1 {
		return nil, ErrByteSliceInvalidLength
	}

	if raw[0] != 0x7b {
		return nil, ErrBadFirstByte
	}

	readings := make([]*Reading, 0, 7)

	for i := 0; i < len(raw)/3; i++ {
		j := i*3 + 1

		// Unused channels are set to 0x7f 0xff 0xff, let's just check for the humidity - which cannot be higher than 100%,
		// so an value of 0xff is obviously not in the range of valid ones.
		if raw[j+2] == 0xff {
			continue
		}

		readings = append(readings, &Reading{
			Sensor: strconv.Itoa(i + 1), // Derived from the channel used by the sensor
			// What happens here:
			// - take the two bytes containing the temperature
			// - decode them as Uint16
			// - convert to / interpret as signed 16-bit integer
			// - finally scale them to 32-bit float
			Temperature: float64(int16(binary.BigEndian.Uint16(raw[j:j+2]))) / 10,
			Humidity:    int8(raw[j+2]),
			Time:        time.Now(),
		})
	}

	return readings, nil
}

func Read() ([]*Reading, error) {
	deviceInfo, err := findDevice()
	if err != nil {
		return nil, err
	}

	device, err := deviceInfo.Open()
	if err != nil {
		return nil, fmt.Errorf("opening device: %w", err)
	}
	defer func() {
		if err := device.Close(); err != nil {
			log.Warn().Err(err).Msg("Could not close device after reading")
		}
	}()

	n, err := device.Write(readRequestBytes)
	if err != nil {
		return nil, fmt.Errorf("writing read request to device: %w", err)
	}

	if n != 64 {
		log.Warn().Int("bytesWritten", n).Msg("64 bytes should have been written to device")
	}

	buf := make([]byte, 25)

	n, err = device.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("reading bytes from device: %w", err)
	}

	if n != 25 {
		log.Warn().Int("bytesRead", n).Msg("25 bytes should have been read")
	}

	readings, err := fromBytes(buf)
	if err != nil {
		return nil, fmt.Errorf("parsing response bytes: %w", err)
	}

	return readings, nil
}

func IsHidSupported() bool {
	return hid.Supported()
}

func findDevice() (*hid.DeviceInfo, error) {
	deviceInfos := hid.Enumerate(vendorID, deviceID)

	switch len(deviceInfos) {
	case 0:
		return nil, ErrNoDeviceFound
	case 1:
		return &deviceInfos[0], nil
	default:
		log.Warn().Int("count", len(deviceInfos)).Msg("Found more than one matching device, simply took the first")
		return &deviceInfos[0], nil
	}
}
