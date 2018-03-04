package ptp

import (
	"io"
)

// PDelReqMsg ...
type PDelReqMsg struct {
	Header
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t *PDelReqMsg) MarshalBinary() ([]byte, error) {

	if t.Header.MessageType != PDelayReqMsgType {
		return nil, ErrInvalidMsgType
	}

	if t.Header.MessageLength == 0 {
		t.Header.MessageLength = HeaderLen + PDelayReqPayloadLen
	}

	b := make([]byte, HeaderLen+PDelayReqPayloadLen)

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[:HeaderLen], headerSlice)

	// All the rest 20 bytes are reserved. Keep them zero values.

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

	err := t.Header.UnmarshalBinary(b[:HeaderLen])
	if err != nil {
		return err
	}

	// All the rest 20 bytes are reserved. Keep them zero values.

	return nil
}
