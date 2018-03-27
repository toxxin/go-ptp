package ptp

import (
	"encoding/binary"
	"io"
)

// SignalingMsg ...
type SignalingMsg struct {
	Header
	ClockIdentity uint64
	PortNumber    uint16
	IntervalRequestTlv
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t *SignalingMsg) MarshalBinary() ([]byte, error) {

	if t.Header.MessageType != SignalingMsgType {
		return nil, ErrInvalidMsgType
	}

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	tlvSlice, err := t.IntervalRequestTlv.MarshalBinary()
	if err != nil {
		return nil, err
	}

	b := make([]byte, HeaderLen+SignalingPayloadLen+IntervalRequestTlvLen)

	copy(b[:HeaderLen], headerSlice)
	offset := HeaderLen

	binary.BigEndian.PutUint64(b[offset:offset+ClockIdentityLen], t.ClockIdentity)
	offset += ClockIdentityLen

	binary.BigEndian.PutUint16(b[offset:offset+2], t.PortNumber)
	offset += 2

	copy(b[offset:offset+IntervalRequestTlvLen], tlvSlice)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a Frame.
func (t *SignalingMsg) UnmarshalBinary(b []byte) error {
	if len(b) < HeaderLen+AnnouncePayloadLen {
		return io.ErrUnexpectedEOF
	}
	err := t.Header.UnmarshalBinary(b[:HeaderLen])
	if err != nil {
		return err
	}

	if t.Header.MessageType != SignalingMsgType {
		return ErrInvalidMsgType
	}

	offset := HeaderLen

	t.ClockIdentity = binary.BigEndian.Uint64(b[offset : offset+ClockIdentityLen])
	offset += ClockIdentityLen
	t.PortNumber = binary.BigEndian.Uint16(b[offset : offset+SourcePortNumberLen])
	offset += SourcePortNumberLen

	if err = t.IntervalRequestTlv.UnmarshalBinary(b[offset:]); err != nil {
		return err
	}

	return nil
}
