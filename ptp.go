package ptp

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

const (
	PtpEtherType  uint16 = 0x88f7
	AvtpEtherType uint16 = 0x22f0
)

// Errors
var (
	ErrInvalidFrame         = errors.New("Invalid frame")
	ErrInvalidHeader        = errors.New("Invalid header")
	ErrInvalidMsgType       = errors.New("Invalid message type")
	ErrUnsupportedVersion   = errors.New("Unsupported protocol version")
	ErrInvalidClockClass    = errors.New("Invalid clock class")
	ErrInvalidClockAccuracy = errors.New("Invalid clock accuracy")
	ErrInvalidTimeSource    = errors.New("Invalid time source")
	ErrInvalidTlvType       = errors.New("Invalid TLV type")
	ErrInvalidTlvOrgId      = errors.New("Invalid TLV organizationId")
	ErrInvalidTlvOrgSubType = errors.New("Invalid organization sub type")
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
	AnnounceMsgType           MsgType = 0xB
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

// MulticastType
type MulticastType uint8

const (
	McastNone MulticastType = iota
	McastPdelay
	McastTestStatus
	McastOther
)

// PortSate
type PortState uint8

const (
	Initializing PortState = iota
	Faulty
	Disabled
	Listening
	PreMaster
	Master
	Passive
	Uncalibrated
	Slave
)

// Length in octets of main fields
const (
	MessageLengthLen           = 2
	FlagsLen                   = 2
	CorrectionNanoSecLen       = 6
	CorrectionSubNanoSecLen    = 2
	CorrectionFullLen          = CorrectionNanoSecLen + CorrectionSubNanoSecLen
	ClockIdentityLen           = 8
	SourcePortNumberLen        = 2
	PortIdentityLen            = ClockIdentityLen + SourcePortNumberLen
	SequenceIDLen              = 2
	OriginTimestampSecLen      = 6
	OriginTimestampNanoSecLen  = 4
	OriginTimestampFullLen     = OriginTimestampSecLen + OriginTimestampNanoSecLen
	CurrentUtcOffsetLen        = 2
	GrandMasterClockQualityLen = 4
	GrandMasterIdentityLen     = 8
	StepsRemovedLen            = 2
	Reserved4                  = 4
	Reserved10                 = 10
)

// Length in octets of header and payloads
const (
	HeaderLen                    = 34
	SyncPayloadLen               = OriginTimestampFullLen
	DelayReqPayloadLen           = OriginTimestampFullLen
	FollowUpPayloadLen           = OriginTimestampFullLen
	DelayRespPayloadLen          = OriginTimestampFullLen + PortIdentityLen
	PDelayReqPayloadLen          = Reserved10 + Reserved10
	PDelayRespPayloadLen         = OriginTimestampFullLen + PortIdentityLen
	PDelayRespFollowUpPayloadLen = OriginTimestampFullLen + PortIdentityLen
	AnnouncePayloadLen           = 30
	SignalingPayloadLen          = 10
	GMClockQualityPayloadLen     = 4
	// SignalingPayloadLen depends on TLVs
)

// time2OriginTimestamp converts time.Time into bytes slice 6+4(sec+nanosec)
// accordingly with ptp timestamp format.
func time2OriginTimestamp(t time.Time, b []byte) error {

	if len(b) != OriginTimestampFullLen {
		return io.ErrUnexpectedEOF
	}

	sec := t.Unix()
	nanosec := t.UnixNano()

	secHexSlice := make([]byte, 8)

	binary.BigEndian.PutUint64(secHexSlice, uint64(sec))
	binary.BigEndian.PutUint32(b[6:], uint32(nanosec-sec*1000000000))

	copy(b[:6], secHexSlice[2:])

	return nil
}

// originTimestamp2Time converts 6+4(sec+nanosec) bytes slice into Time
func originTimestamp2Time(b []byte) (time.Time, error) {

	if len(b) != OriginTimestampSecLen+OriginTimestampNanoSecLen {
		return time.Now(), io.ErrUnexpectedEOF
	}

	sec := binary.BigEndian.Uint64(append([]byte{0, 0}, b[:6]...))

	nsecSlice := append([]byte{}, b[6:10]...)
	nsecSlice = append(nsecSlice, []byte{0, 0, 0, 0}...)
	nsec := binary.BigEndian.Uint32(nsecSlice)

	return time.Unix(int64(sec), int64(nsec)), nil
}

const UScaledNsLen = 12

type UScaledNs struct {
	ms int32
	ls uint64
}

func NewUScaledNs(b []byte) (UScaledNs, error) {
	if len(b) != UScaledNsLen {
		return UScaledNs{}, io.ErrUnexpectedEOF
	}

	return UScaledNs{
		ms: int32(binary.BigEndian.Uint32(b[:4])),
		ls: binary.BigEndian.Uint64(b[4:]),
	}, nil
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (p *UScaledNs) MarshalBinary() ([]byte, error) {
	b := make([]byte, UScaledNsLen)

	binary.BigEndian.PutUint32(b[:4], uint32(p.ms))

	binary.BigEndian.PutUint64(b[4:], p.ls)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a UScaledNs.
//
// If the byte slice does not contain enough data to unmarshal a valid UScaledNs,
// io.ErrUnexpectedEOF is returned.
func (p *UScaledNs) UnmarshalBinary(b []byte) error {
	if len(b) != UScaledNsLen {
		return io.ErrUnexpectedEOF
	}

	p.ms = int32(binary.BigEndian.Uint32(b[:4]))
	p.ls = binary.BigEndian.Uint64(b[4:])

	return nil
}

// GetClockIdByMac takes MAC address as a slice and converts it
// into slice of bytes(EUI-64) in accordance with IEEE 1588v2 spec.
func GetClockIdByMac(b []byte) ([]byte, error) {
	if len(b) != 6 {
		return []byte{}, io.ErrUnexpectedEOF
	}

	res := make([]byte, ClockIdentityLen)

	copy(res[:3], b[:3])
	b[3], b[4] = 0xff, 0xfe
	copy(res[5:], b[3:])
	return res, nil
}
