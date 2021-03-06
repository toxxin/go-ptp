package ptp

import (
	"encoding/binary"
	"io"
)

// ProtoVersion is PTP protocol version
type ProtoVersion uint8

// Version numbers
const (
	_ ProtoVersion = iota
	Version1
	Version2
	Version3
)

// Flags is header's field to indicate status
type Flags struct {
	LI61               bool
	LI59               bool
	UtcReasonable      bool
	TimeScale          bool
	TimeTraceable      bool
	FrequencyTraceable bool
	AlternateMaster    bool
	TwoSteps           bool
	Unicast            bool
	ProfileSpecific1   bool
	ProfileSpecific2   bool
	Security           bool
}

const (
	lI61Bit               uint16 = 1 << 0
	lI59Bit               uint16 = 1 << 1
	utcReasonableBit      uint16 = 1 << 2
	timeScaleBit          uint16 = 1 << 3
	timeTraceableBit      uint16 = 1 << 4
	frequencyTraceableBit uint16 = 1 << 5
	alternateMasterBit    uint16 = 1 << 8
	twoStepsBit           uint16 = 1 << 9
	unicastBit            uint16 = 1 << 10
	profileSpecific1Bit   uint16 = 1 << 13
	profileSpecific2Bit   uint16 = 1 << 14
	securityBit           uint16 = 1 << 15
)

func b2i(b bool) uint16 {
	if b {
		return 1
	}
	return 0
}

// MarshalBinary returns Flags as uint16 value.
func (f *Flags) MarshalBinary() uint16 {
	return (b2i(f.LI61)<<0 |
		b2i(f.LI59)<<1 |
		b2i(f.UtcReasonable)<<2 |
		b2i(f.TimeScale)<<3 |
		b2i(f.TimeTraceable)<<4 |
		b2i(f.FrequencyTraceable)<<5 |
		b2i(f.AlternateMaster)<<8 |
		b2i(f.TwoSteps)<<9 |
		b2i(f.Unicast)<<10 |
		b2i(f.ProfileSpecific1)<<13 |
		b2i(f.ProfileSpecific2)<<14 |
		b2i(f.Security)<<15)
}

func (f *Flags) UnmarshalBinary(b []byte) error {
	if len(b) != 2 {
		return io.ErrUnexpectedEOF
	}

	flags := binary.BigEndian.Uint16(b[:])

	f.LI61 = flags&lI61Bit != 0
	f.LI59 = flags&lI59Bit != 0
	f.UtcReasonable = flags&utcReasonableBit != 0
	f.TimeScale = flags&timeScaleBit != 0
	f.TimeTraceable = flags&timeTraceableBit != 0
	f.FrequencyTraceable = flags&frequencyTraceableBit != 0
	f.AlternateMaster = flags&alternateMasterBit != 0
	f.TwoSteps = flags&twoStepsBit != 0
	f.Unicast = flags&unicastBit != 0
	f.ProfileSpecific1 = flags&profileSpecific1Bit != 0
	f.ProfileSpecific2 = flags&profileSpecific2Bit != 0
	f.Security = flags&securityBit != 0

	return nil
}

// Header struct describes the header of a PTP message.
type Header struct {
	Flags
	MessageType      MsgType
	MessageLength    uint16
	VersionPTP       ProtoVersion
	CorrectionNs     uint64
	CorrectionSubNs  uint16
	ClockIdentity    uint64
	PortNumber       uint16
	SequenceID       uint16
	LogMessagePeriod int8
}

// MarshalBinary allocates a byte slice and marshals a Header into binary form.
func (h *Header) MarshalBinary() ([]byte, error) {

	var correction uint64

	b := make([]byte, HeaderLen)
	offset := 0

	// Transport specific, messageId
	b[0] = 0x0 | uint8(h.MessageType)
	offset++

	// PTP proto version
	b[1] = byte(Version2)
	offset++

	// Message length
	b[offset] = byte(h.MessageLength >> 8)
	offset++
	b[offset] = byte(h.MessageLength)
	offset++

	// Subdomain number
	b[offset] = 0x0
	offset++

	// Skip reserved byte
	offset++

	flags := (&h.Flags).MarshalBinary()

	binary.BigEndian.PutUint16(b[offset:offset+FlagsLen], flags)
	offset += FlagsLen

	// Correction Ns & SubNs
	correction = (h.CorrectionNs << 2) | (uint64)(h.CorrectionSubNs)
	binary.BigEndian.PutUint64(b[offset:offset+CorrectionFullLen], correction)
	offset += CorrectionFullLen

	// Skip 4 reserved bytes
	offset += 4

	// Clock identity
	binary.BigEndian.PutUint64(b[offset:offset+ClockIdentityLen], h.ClockIdentity)
	offset += ClockIdentityLen

	// Source port
	binary.BigEndian.PutUint16(b[offset:offset+SourcePortNumberLen], h.PortNumber)
	offset += SourcePortNumberLen

	// Sequence ID
	binary.BigEndian.PutUint16(b[offset:offset+SequenceIDLen], h.SequenceID)
	offset += SequenceIDLen

	var msgCtrl MsgCtrlType
	switch h.MessageType {
	case (SyncMsgType):
		msgCtrl = SyncMsgCtrlType
	case (DelayReqMsgType):
		msgCtrl = DelayReqMsgCtrlType
	case (FollowUpMsgType):
		msgCtrl = FollowUpMsgCtrlType
	case (DelayRespMsgType):
		msgCtrl = DelayRespMsgCtrlType
	default:
		msgCtrl = OtherMsgCtrlType
	}

	b[offset] = byte(msgCtrl)
	offset++

	b[offset] = (byte)(h.LogMessagePeriod)
	offset++

	return b, nil
}

func isValidMsgType(msgtype MsgType) bool {
	switch msgtype {
	case
		SyncMsgType,
		DelayReqMsgType,
		PDelayReqMsgType,
		PDelayRespMsgType,
		FollowUpMsgType,
		DelayRespMsgType,
		PDelayRespFollowUpMsgType,
		AnnounceMsgType,
		SignalingMsgType,
		MgmtMsgType:
		return true
	}
	return false
}

// UnmarshalBinary unmarshals a byte slice into a Header.
func (h *Header) UnmarshalBinary(b []byte) error {
	if len(b) != HeaderLen {
		return io.ErrUnexpectedEOF
	}

	h.MessageType = MsgType(b[0] & 0x0f)
	if !isValidMsgType(h.MessageType) {
		return ErrInvalidMsgType
	}

	// TODO: Add implementation another versions
	h.VersionPTP = ProtoVersion(0xf & b[1])
	if h.VersionPTP != Version2 {
		return ErrUnsupportedVersion
	}

	h.MessageLength = binary.BigEndian.Uint16(b[2:4])

	h.Flags.UnmarshalBinary(b[6:8])

	// Correct Ns & SubNs
	tmpSlice := make([]byte, 8)
	copy(tmpSlice[2:], b[8:14])

	h.CorrectionNs = binary.BigEndian.Uint64(tmpSlice)
	h.CorrectionSubNs = binary.BigEndian.Uint16(b[14:16])

	h.ClockIdentity = binary.BigEndian.Uint64(b[20:28])
	h.PortNumber = binary.BigEndian.Uint16(b[28:30])

	h.SequenceID = binary.BigEndian.Uint16(b[30:32])
	h.LogMessagePeriod = int8(b[33])

	return nil
}
