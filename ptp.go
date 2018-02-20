package ptp

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

// Errors
var (
	ErrInvalidFrame         = errors.New("Invalid frame")
	ErrInvalidHeader        = errors.New("Invalid header")
	ErrInvalidMsgType       = errors.New("Invalid message type")
	ErrUnsupportedVersion   = errors.New("Unsupported protocol version")
	ErrInvalidClockClass    = errors.New("Invalid clock class")
	ErrInvalidClockAccuracy = errors.New("Invalid clock accuracy")
)

// MsgType Type
type MsgType uint8

// Message types codes
const (
	SyncMsgType               MsgType = 0x0
	DelayReqMsgType           MsgType = 0x1
	PDelayReqMsgType          MsgType = 0x2
	PDelayRespMsgType         MsgType = 0x3
	FollowUpMsgType           MsgType = 0x8
	DelayRespMsgType          MsgType = 0x9
	PDelayRespFollowUpMsgType MsgType = 0xA
	AnnonceMsgType            MsgType = 0xB
	SignalingMsgType          MsgType = 0xC
	MgmtMsgType               MsgType = 0xD
)

// MsgCtrlType Control Type
type MsgCtrlType uint8

// Message control codes
const (
	SyncMsgCtrlType      MsgCtrlType = 0
	DelayReqMsgCtrlType  MsgCtrlType = 1
	FollowUpMsgCtrlType  MsgCtrlType = 2
	DelayRespMsgCtrlType MsgCtrlType = 3
	OtherMsgCtrlType     MsgCtrlType = 5
)

// Length in octets of main fields
const (
	MessageLengthLen          = 2
	FlagsLen                  = 2
	CorrectionNanoSecLen      = 6
	CorrectionSubNanoSecLen   = 2
	CorrectionFullLen         = CorrectionNanoSecLen + CorrectionSubNanoSecLen
	ClockIdentityLen          = 8
	SourcePortIDLen           = 2
	PortIdentityLen           = ClockIdentityLen + SourcePortIDLen
	SequenceIDLen             = 2
	OriginTimestampSecLen     = 6
	OriginTimestampNanoSecLen = 4
	OriginTimestampFullLen    = OriginTimestampSecLen + OriginTimestampNanoSecLen
	CurrentUtcOffsetLen       = 2
	GrandMasterClockQuality   = 4
	GrandMasterIdentityLen    = 8
	StepsRemovedLen           = 2
	Reserved4                 = 4
	Reserved10                = 10
)

// Length in octets of header and payloads
const (
	HeaderLen                    = 34
	SyncPayloadLen               = OriginTimestampFullLen
	DelayReqPayloadLen           = OriginTimestampFullLen
	FollowUpPayloadLen           = OriginTimestampFullLen
	DelayRespPayloadLen          = OriginTimestampFullLen + PortIdentityLen
	PDelayReqPayloadLen          = OriginTimestampFullLen + Reserved10
	PDelayRespPayloadLen         = OriginTimestampFullLen + PortIdentityLen
	PDelayRespFollowUpPayloadLen = OriginTimestampFullLen + PortIdentityLen
	AnnouncePayloadLen           = 30
	// SignalingPayloadLen depends on TLVs
)

// time2OriginTimestamp allocates 6+4(sec+nanosec) bytes slice
// and converts Time into binary form accordingly with ptp timestamp format
func time2OriginTimestamp(t time.Time) []byte {

	sec := t.Unix()
	nanosec := t.UnixNano()

	secHexSlice := make([]byte, 8)
	nanoHexSlice := make([]byte, 4)

	binary.BigEndian.PutUint64(secHexSlice, uint64(sec))
	binary.BigEndian.PutUint32(nanoHexSlice, uint32(nanosec-sec*1000000000))

	res := make([]byte, 10)
	copy(res[:6], secHexSlice[2:])
	copy(res[6:], nanoHexSlice)

	return res
}

// originTimestamp2Time converts 6+4(sec+nanosec) bytes slice into Time
func originTimestamp2Time(b []byte) (time.Time, error) {

	if len(b) != OriginTimestampSecLen+OriginTimestampNanoSecLen {
		return time.Now(), io.ErrUnexpectedEOF
	}

	sec := binary.BigEndian.Uint64(append([]byte{0, 0}, b[:6]...))
	nsec := binary.BigEndian.Uint32(append(b[6:10], []byte{0, 0, 0, 0}...))

	return time.Unix(int64(sec), int64(nsec)), nil
}
