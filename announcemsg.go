package ptp

import (
	"encoding/binary"
	"io"
	"time"
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

// GMClockQuality defines Grand Master Clock Quality
type GMClockQuality struct {
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

// UnmarshalBinary unmarshals a byte slice into a Header.
func (p *GMClockQuality) UnmarshalBinary(b []byte) error {
	if len(b) != GMClockQualityPayloadLen {
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

// TimeSourceType Type
type TimeSourceType uint8

// TimeSource types codes
const (
	TimeSourceAtomic      TimeSourceType = 16
	TimeSourceGPS         TimeSourceType = 32
	TimeSourceTRadio      TimeSourceType = 48
	TimeSourcePTP         TimeSourceType = 64
	TimeSourceNTP         TimeSourceType = 80
	TimeSourceHandSet     TimeSourceType = 96
	TimeSourceOther       TimeSourceType = 144
	TimeSourceInternalOsc TimeSourceType = 160
)

func isValidTimeSource(t TimeSourceType) bool {
	switch t {
	case
		TimeSourceAtomic,
		TimeSourceGPS,
		TimeSourceTRadio,
		TimeSourcePTP,
		TimeSourceNTP,
		TimeSourceHandSet,
		TimeSourceOther,
		TimeSourceInternalOsc:
		return true
	}
	return false
}

// AnnounceMsg ...
type AnnounceMsg struct {
	Header
	GMClockQuality
	OriginTimestamp  time.Time
	CurrentUtcOffset int16
	GMPriority1      uint8
	GMPriority2      uint8
	GMIdentity       uint64
	StepsRemoved     uint16
	TimeSource       TimeSourceType
}

// UnmarshalBinary unmarshals a byte slice into a Header.
func (t *AnnounceMsg) UnmarshalBinary(b []byte) error {
	if len(b) != HeaderLen+AnnouncePayloadLen {
		return io.ErrUnexpectedEOF
	}
	err := t.Header.UnmarshalBinary(b[:34])
	if err != nil {
		return err
	}

	offset := 34

	t.OriginTimestamp, err = originTimestamp2Time(b[offset : offset+10])
	if err != nil {
		return err
	}
	offset += 10

	utcoffset := binary.BigEndian.Uint16(b[offset : offset+2])
	t.CurrentUtcOffset = int16(utcoffset)
	offset += 2

	// Reserved byte
	offset++

	t.GMPriority1 = b[offset]
	offset++

	err = t.GMClockQuality.UnmarshalBinary(b[offset : offset+GMClockQualityPayloadLen])
	if err != nil {
		return err
	}
	offset += GMClockQualityPayloadLen

	t.GMPriority2 = b[offset]
	offset++

	t.GMIdentity = binary.BigEndian.Uint64(b[offset : offset+GrandMasterIdentityLen])
	offset += GrandMasterIdentityLen

	t.StepsRemoved = binary.BigEndian.Uint16(b[offset : offset+StepsRemovedLen])
	offset += StepsRemovedLen

	t.TimeSource = TimeSourceType(b[offset])
	if !isValidTimeSource(t.TimeSource) {
		return ErrInvalidTimeSource
	}

	return nil
}
