package ptp

import (
	"io"
	"time"
)

// FollowUpMsg ...
type FollowUpMsg struct {
	Header
	PreciseOriginTimestamp time.Time
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t *FollowUpMsg) MarshalBinary() ([]byte, error) {

	if t.Header.MessageType != FollowUpMsgType {
		return nil, ErrInvalidMsgType
	}

	if t.Header.MessageLength == 0 {
		t.Header.MessageLength = HeaderLen + FollowUpPayloadLen
	}

	if t.Header.MessageLength != HeaderLen+FollowUpPayloadLen {
		return nil, io.ErrUnexpectedEOF
	}

	b := make([]byte, HeaderLen+FollowUpPayloadLen)

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[:HeaderLen], headerSlice)

	// Origin timestamp
	tslice := time2OriginTimestamp(t.PreciseOriginTimestamp)
	copy(b[HeaderLen:], tslice)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a FollowUpMsg.
//
// If the byte slice does not contain enough data to unmarshal a valid FollowUpMsg,
// io.ErrUnexpectedEOF is returned.
func (t *FollowUpMsg) UnmarshalBinary(b []byte) error {
	if len(b) < HeaderLen+FollowUpPayloadLen {
		return io.ErrUnexpectedEOF
	}
	err := t.Header.UnmarshalBinary(b[:HeaderLen])
	if err != nil {
		return err
	}

	if t.PreciseOriginTimestamp, err = originTimestamp2Time(b[HeaderLen:]); err != nil {
		return err
	}

	return nil
}
