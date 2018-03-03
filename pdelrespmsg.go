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

	b := make([]byte, HeaderLen+PDelayRespPayloadLen)

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[:HeaderLen], headerSlice)
	offset := HeaderLen

	tmsSlice := time2OriginTimestamp(t.ReceiveTimestamp)
	copy(b[offset:offset+10], tmsSlice)
	offset += 10

	clockIdentitySlice := make([]byte, 8)
	binary.BigEndian.PutUint64(clockIdentitySlice, t.ClockIdentity)
	copy(b[offset:offset+8], clockIdentitySlice)
	offset += 8

	portIDSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(portIDSlice, t.PortNumber)
	copy(b[offset:offset+2], portIDSlice)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a PDelRespMsg.
//
// If the byte slice does not contain enough data to unmarshal a valid PDelRespMsg,
// io.ErrUnexpectedEOF is returned.
func (t *PDelRespMsg) UnmarshalBinary(b []byte) error {
	// Must contain type and length values
	if len(b) < HeaderLen+PDelayRespPayloadLen {
		return io.ErrUnexpectedEOF
	}

	return ErrInvalidFrame
}
