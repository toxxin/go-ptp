package ptp

import (
	"encoding/binary"
	"io"
)

// ClockClassType Type
type ClockClassType uint8

// ClockClass types codes
const (
	PrimarySyncRefClass ClockClassType = 6
	LostSyncClass       ClockClassType = 7
	DefaultClass        ClockClassType = 248
	SlaveOnlyClass      ClockClassType = 255
)

// ClockAccuracyType Type
type ClockAccuracyType uint8

// ClockAccuracy types codes
const (
	ClockAccuracy25ns         ClockAccuracyType = 32
	ClockAccuracy100ns        ClockAccuracyType = 33
	ClockAccuracy250ns        ClockAccuracyType = 34
	ClockAccuracy1mics        ClockAccuracyType = 35
	ClockAccuracy2_5mics      ClockAccuracyType = 36
	ClockAccuracy10mics       ClockAccuracyType = 37
	ClockAccuracy25mics       ClockAccuracyType = 38
	ClockAccuracy100mics      ClockAccuracyType = 39
	ClockAccuracy250mics      ClockAccuracyType = 40
	ClockAccuracy1ms          ClockAccuracyType = 41
	ClockAccuracy2_5ms        ClockAccuracyType = 42
	ClockAccuracy10ms         ClockAccuracyType = 43
	ClockAccuracy25ms         ClockAccuracyType = 44
	ClockAccuracy100ms        ClockAccuracyType = 45
	ClockAccuracy250ms        ClockAccuracyType = 46
	ClockAccuracy1s           ClockAccuracyType = 47
	ClockAccuracy10s          ClockAccuracyType = 48
	ClockAccuracyMore10s      ClockAccuracyType = 49
	ClockAccuracyNotSupported ClockAccuracyType = 255
)

// ClockQuality defines Grand Master Clock Quality
type ClockQuality struct {
	ClockClass    ClockClassType
	ClockAccuracy ClockAccuracyType
	ClockVariance uint16
}

func isValidClockClass(class ClockClassType) bool {
	switch class {
	case
		PrimarySyncRefClass,
		LostSyncClass,
		DefaultClass,
		SlaveOnlyClass:
		return true
	}
	return false
}

func isValidClockAccuracy(c ClockAccuracyType) bool {
	switch c {
	case
		ClockAccuracy25ns,
		ClockAccuracy100ns,
		ClockAccuracy250ns,
		ClockAccuracy1mics,
		ClockAccuracy2_5mics,
		ClockAccuracy10mics,
		ClockAccuracy25mics,
		ClockAccuracy100mics,
		ClockAccuracy250mics,
		ClockAccuracy1ms,
		ClockAccuracy2_5ms,
		ClockAccuracy10ms,
		ClockAccuracy25ms,
		ClockAccuracy100ms,
		ClockAccuracy250ms,
		ClockAccuracy1s,
		ClockAccuracy10s,
		ClockAccuracyMore10s,
		ClockAccuracyNotSupported:
		return true
	}
	return false
}

// MarshalBinary allocates a byte slice and marshals a Header into binary form.
func (p ClockQuality) MarshalBinary() ([]byte, error) {

	b := make([]byte, 4)

	b[0] = uint8(p.ClockClass)
	b[1] = uint8(p.ClockAccuracy)

	binary.BigEndian.PutUint16(b[2:4], uint16(p.ClockVariance))

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a Header.
func (p *ClockQuality) UnmarshalBinary(b []byte) error {
	if len(b) != ClockQualityPayloadLen {
		return io.ErrUnexpectedEOF
	}

	p.ClockClass = ClockClassType(b[0])
	if !isValidClockClass(p.ClockClass) {
		return ErrInvalidClockClass
	}

	p.ClockAccuracy = ClockAccuracyType(b[1])
	if !isValidClockAccuracy(p.ClockAccuracy) {
		return ErrInvalidClockAccuracy
	}

	p.ClockVariance = binary.BigEndian.Uint16(b[2:])

	return nil
}
