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

	binary.BigEndian.PutUint64(b[offset:offset+ClockIdentityLen], t.ClockIdentity)
	offset += ClockIdentityLen

	binary.BigEndian.PutUint16(b[offset:offset+SourcePortNumberLen], t.PortNumber)

	return b, nil
}
