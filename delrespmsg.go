package ptp

import (
	"encoding/binary"
	"io"
	"time"
)

// DelRespMsg ...
type DelRespMsg struct {
	Header
	ReceiveTimestamp       time.Time
	RequestingPortIdentity uint64
	RequestingPortID       uint16
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t *DelRespMsg) MarshalBinary() ([]byte, error) {

	if t.Header.MessageType != DelayRespMsgType {
		return nil, ErrInvalidMsgType
	}

	b := make([]byte, HeaderLen+DelayRespPayloadLen)

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[:HeaderLen], headerSlice)
	offset := HeaderLen

	//TODO: add receiveTimestamp
	offset += 10

	portIdentitySlice := make([]byte, 8)
	binary.BigEndian.PutUint64(portIdentitySlice, t.RequestingPortIdentity)
	copy(b[offset:offset+8], portIdentitySlice)
	offset += 8

	portIDSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(portIDSlice, t.RequestingPortID)
	copy(b[offset:offset+2], portIDSlice)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a DelRespMsg.
//
// If the byte slice does not contain enough data to unmarshal a valid DelRespMsg,
// io.ErrUnexpectedEOF is returned.
func (t *DelRespMsg) UnmarshalBinary(b []byte) error {
	// Must contain type and length values
	if len(b) < HeaderLen+DelayRespPayloadLen {
		return io.ErrUnexpectedEOF
	}

	return ErrInvalidFrame
}
