package ptp

import (
	"encoding/binary"
	"io"
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

// UnmarshalBinary unmarshals a byte slice into a PDelRespFollowUpMsg.
//
// If the byte slice does not contain enough data to unmarshal a valid PDelRespFollowUpMsg,
// io.ErrUnexpectedEOF is returned.
func (t *PDelRespFollowUpMsg) UnmarshalBinary(b []byte) error {

	if len(b) != HeaderLen+PDelayRespFollowUpPayloadLen {
		return io.ErrUnexpectedEOF
	}

	err := t.Header.UnmarshalBinary(b[:HeaderLen])
	if err != nil {
		return err
	}

	if t.Header.MessageType != PDelayRespFollowUpMsgType {
		return ErrInvalidMsgType
	}

	if t.OriginTimestamp, err = originTimestamp2Time(b[HeaderLen : HeaderLen+OriginTimestampFullLen]); err != nil {
		return err
	}
	offset := HeaderLen + OriginTimestampFullLen

	t.ClockIdentity = binary.BigEndian.Uint64(b[offset : offset+ClockIdentityLen])
	offset += ClockIdentityLen

	t.PortNumber = binary.BigEndian.Uint16(b[offset : offset+SourcePortNumberLen])

	return nil
}
