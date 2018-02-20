package ptp

import (
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

// ClockClass types codes
const (
	ClockAccuracy25ns    ClockAccuracyType = 32
	ClockAccuracy100ns   ClockAccuracyType = 33
	ClockAccuracy250ns   ClockAccuracyType = 34
	ClockAccuracy1mics   ClockAccuracyType = 35
	ClockAccuracy2_5mics ClockAccuracyType = 36
	ClockAccuracy10mics  ClockAccuracyType = 37
	ClockAccuracy25mics  ClockAccuracyType = 38
	ClockAccuracy100mics ClockAccuracyType = 39
	ClockAccuracy250mics ClockAccuracyType = 40
	ClockAccuracy1ms     ClockAccuracyType = 41
	ClockAccuracy2_5ms   ClockAccuracyType = 42
	ClockAccuracy10ms    ClockAccuracyType = 43
	ClockAccuracy25ms    ClockAccuracyType = 44
	ClockAccuracy100ms   ClockAccuracyType = 45
	ClockAccuracy250ms   ClockAccuracyType = 46
	ClockAccuracy1s      ClockAccuracyType = 47
	ClockAccuracy10s     ClockAccuracyType = 48
	ClockAccuracymore10s ClockAccuracyType = 49
)

// GMClockQuality defines Grand Master Clock Quality
type GMClockQuality struct {
	ClockClass              ClockClassType
	ClockAccuracy           ClockAccuracyType
	OffsetScaledLogVariance uint16
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
	TimeSourcehandSet     TimeSourceType = 96
	TimeSourceOther       TimeSourceType = 144
	TimeSourceInternalOsc TimeSourceType = 160
)

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
