package ptp

import (
	"io"
	"time"
)

// DelReqMsg ...
type DelReqMsg struct {
	Header
	OriginTimestamp time.Time
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t *DelReqMsg) MarshalBinary() ([]byte, error) {

	if t.Header.MessageType != DelayReqMsgType {
		return nil, ErrInvalidMsgType
	}

	b := make([]byte, HeaderLen+DelayReqPayloadLen)

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[:HeaderLen], headerSlice)

	// Origin timestamp
	// Assume 0x0 value for testing

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a DelReqMsg.
//
// If the byte slice does not contain enough data to unmarshal a valid DelReqMsg,
// io.ErrUnexpectedEOF is returned.
func (t *DelReqMsg) UnmarshalBinary(b []byte) error {
	// Must contain type and length values
	if len(b) < HeaderLen+DelayReqPayloadLen {
		return io.ErrUnexpectedEOF
	}

	return ErrInvalidFrame
}
