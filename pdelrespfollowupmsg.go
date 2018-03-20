package ptp

import (
	"encoding/binary"
	"time"
)

// PDelRespFollowUpMsg ...
type PDelRespFollowUpMsg struct {
	Header
	OriginTimestamp time.Time
	ClockIdentity   uint64
	PortNumber      uint16
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t *PDelRespFollowUpMsg) MarshalBinary() ([]byte, error) {

	if t.Header.MessageType != PDelayRespFollowUpMsgType {
		return nil, ErrInvalidMsgType
	}

	b := make([]byte, HeaderLen+PDelayRespFollowUpPayloadLen)

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[:HeaderLen], headerSlice)
	offset := HeaderLen

	// Origin timestamp
	time2OriginTimestamp(t.OriginTimestamp, b[offset:offset+OriginTimestampFullLen])
	offset += OriginTimestampFullLen

	clockIdentitySlice := make([]byte, ClockIdentityLen)
	binary.BigEndian.PutUint64(clockIdentitySlice, t.ClockIdentity)
	copy(b[offset:offset+ClockIdentityLen], clockIdentitySlice)
	offset += ClockIdentityLen

	portIDSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(portIDSlice, t.PortNumber)
	copy(b[offset:offset+2], portIDSlice)

	return b, nil
}
