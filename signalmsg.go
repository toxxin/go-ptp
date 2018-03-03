package ptp

import "encoding/binary"

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

	clockIdentitySlice := make([]byte, ClockIdentityLen)
	binary.BigEndian.PutUint64(clockIdentitySlice, t.ClockIdentity)
	copy(b[offset:offset+ClockIdentityLen], clockIdentitySlice)
	offset += ClockIdentityLen

	portNumberSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(portNumberSlice, t.PortNumber)
	copy(b[offset:offset+2], portNumberSlice)
	offset += 2

	copy(b[offset:offset+IntervalRequestTlvLen], tlvSlice)

	return b, nil
}
