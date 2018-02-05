package ptp

import (
	"encoding/binary"
	"io"
	"time"
)

// PDelReqMsg ...
type PDelReqMsg struct {
	Header
	OriginTimestamp time.Time
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t *PDelReqMsg) MarshalBinary() ([]byte, error) {

	if t.Header.MessageType != PDelayReqMsgType {
		return nil, ErrInvalidMsgType
	}

	if t.Header.MessageLength == 0 {
		t.Header.MessageLength = HeaderLen + PDelayReqPayloadLen
	}

	if t.Header.MessageLength != HeaderLen+PDelayReqPayloadLen {
		return nil, io.ErrUnexpectedEOF
	}

	b := make([]byte, HeaderLen+PDelayReqPayloadLen)

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[:HeaderLen], headerSlice)

	// Origin timestamp
	tslice := time2OriginTimestamp(t.OriginTimestamp)
	copy(b[HeaderLen:], tslice)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a PDelReqMsg.
//
// If the byte slice does not contain enough data to unmarshal a valid PDelReqMsg,
// io.ErrUnexpectedEOF is returned.
func (t *PDelReqMsg) UnmarshalBinary(b []byte) error {

	if len(b) != HeaderLen+PDelayReqPayloadLen {
		return io.ErrUnexpectedEOF
	}

	err := t.Header.UnmarshalBinary(b[:34])
	if err != nil {
		return err
	}

	secSlice := make([]byte, 8)
	copy(secSlice[2:], b[34:40])

	sec := binary.BigEndian.Uint64(secSlice)
	nsec := binary.BigEndian.Uint32(b[40:44])

	t.OriginTimestamp = time.Unix(int64(sec), int64(nsec))

	return nil
}
