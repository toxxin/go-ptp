package ptp

import (
	"encoding/binary"
	"io"
	"time"
)

// PDelRespMsg ...
type PDelRespMsg struct {
	Header
	ReceiveTimestamp time.Time
	ClockIdentity    uint64
	PortNumber       uint16
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t *PDelRespMsg) MarshalBinary() ([]byte, error) {

	if t.Header.MessageType != PDelayRespMsgType {
		return nil, ErrInvalidMsgType
	}

	if t.Header.MessageLength == 0 {
		t.Header.MessageLength = HeaderLen + PDelayRespPayloadLen
	}

	b := make([]byte, HeaderLen+PDelayRespPayloadLen)

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[:HeaderLen], headerSlice)
	offset := HeaderLen

	time2OriginTimestamp(t.ReceiveTimestamp, b[offset:offset+OriginTimestampFullLen])
	offset += OriginTimestampFullLen

	binary.BigEndian.PutUint64(b[offset:offset+ClockIdentityLen], t.ClockIdentity)
	offset += ClockIdentityLen

	binary.BigEndian.PutUint16(b[offset:offset+2], t.PortNumber)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a PDelRespMsg.
//
// If the byte slice does not contain enough data to unmarshal a valid PDelRespMsg,
// io.ErrUnexpectedEOF is returned.
func (t *PDelRespMsg) UnmarshalBinary(b []byte) error {

	if len(b) != HeaderLen+PDelayRespPayloadLen {
		return io.ErrUnexpectedEOF
	}

	err := t.Header.UnmarshalBinary(b[:HeaderLen])
	if err != nil {
		return err
	}

	if t.Header.MessageType != PDelayRespMsgType {
		return ErrInvalidMsgType
	}

	if t.ReceiveTimestamp, err = originTimestamp2Time(b[HeaderLen : HeaderLen+OriginTimestampFullLen]); err != nil {
		return err
	}
	offset := HeaderLen + OriginTimestampFullLen

	t.ClockIdentity = binary.BigEndian.Uint64(b[offset : offset+ClockIdentityLen])
	offset += ClockIdentityLen

	t.PortNumber = binary.BigEndian.Uint16(b[offset : offset+SourcePortNumberLen])

	return nil
}
