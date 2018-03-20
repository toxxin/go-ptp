package ptp

import (
	"io"
	"time"
)

// SyncMsg ...
type SyncMsg struct {
	Header
	OriginTimestamp time.Time
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t *SyncMsg) MarshalBinary() ([]byte, error) {

	if t.Header.MessageType != SyncMsgType {
		return nil, ErrInvalidMsgType
	}

	if t.Header.MessageLength == 0 {
		t.Header.MessageLength = HeaderLen + SyncPayloadLen
	}

	if t.Header.MessageLength != HeaderLen+SyncPayloadLen {
		return nil, io.ErrUnexpectedEOF
	}

	b := make([]byte, HeaderLen+SyncPayloadLen)

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[:HeaderLen], headerSlice)

	// Origin timestamp
	time2OriginTimestamp(t.OriginTimestamp, b[HeaderLen:])

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a SyncMsg.
//
// If the byte slice does not contain enough data to unmarshal a valid SyncMsg,
// io.ErrUnexpectedEOF is returned.
func (t *SyncMsg) UnmarshalBinary(b []byte) error {
	if len(b) < HeaderLen+SyncPayloadLen {
		return io.ErrUnexpectedEOF
	}

	err := t.Header.UnmarshalBinary(b[:HeaderLen])
	if err != nil {
		return err
	}

	if t.Header.MessageType != SyncMsgType {
		return ErrInvalidMsgType
	}

	if t.OriginTimestamp, err = originTimestamp2Time(b[HeaderLen:]); err != nil {
		return err
	}

	return nil
}
