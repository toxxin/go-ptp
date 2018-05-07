package ptp

import (
	"encoding/binary"
	"io"
)

// AnnounceMsg ...
type AnnounceMsg struct {
	Header
	GMClockQuality   ClockQuality
	CurrentUtcOffset int16
	GMPriority1      uint8
	GMPriority2      uint8
	GMIdentity       uint64
	StepsRemoved     uint16
	TimeSource       TimeSourceType
	PathTraceTlv
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (t AnnounceMsg) MarshalBinary() ([]byte, error) {
	if t.Header.MessageType != AnnounceMsgType {
		return nil, ErrInvalidMsgType
	}

	tlvSlice, err := t.PathTraceTlv.MarshalBinary()
	if err != nil {
		return nil, err
	}

	headerSlice, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	b := make([]byte, HeaderLen+AnnouncePayloadLen+len(tlvSlice))

	copy(b[:HeaderLen], headerSlice)
	offset := HeaderLen

	// Reserved 10 bytes
	offset += Reserved10

	binary.BigEndian.PutUint16(b[offset:offset+CurrentUtcOffsetLen], uint16(t.CurrentUtcOffset))
	offset += CurrentUtcOffsetLen

	// Reserved byte
	offset++

	b[offset] = t.GMPriority1
	offset++

	gmClockQualitySlice, err := t.GMClockQuality.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(b[offset:offset+ClockQualityPayloadLen], gmClockQualitySlice)
	offset += ClockQualityPayloadLen

	b[offset] = t.GMPriority2
	offset++

	binary.BigEndian.PutUint64(b[offset:offset+ClockIdentityLen], t.GMIdentity)
	offset += ClockIdentityLen

	binary.BigEndian.PutUint16(b[offset:offset+StepsRemovedLen], uint16(t.StepsRemoved))
	offset += StepsRemovedLen

	b[offset] = uint8(t.TimeSource)
	offset++

	copy(b[offset:offset+len(tlvSlice)], tlvSlice)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a Frame.
func (t *AnnounceMsg) UnmarshalBinary(b []byte) error {
	if len(b) < HeaderLen+AnnouncePayloadLen {
		return io.ErrUnexpectedEOF
	}
	err := t.Header.UnmarshalBinary(b[:HeaderLen])
	if err != nil {
		return err
	}

	if t.Header.MessageType != AnnounceMsgType {
		return ErrInvalidMsgType
	}

	offset := HeaderLen

	// Reserved 10 bytes
	offset += Reserved10

	utcoffset := binary.BigEndian.Uint16(b[offset : offset+2])
	t.CurrentUtcOffset = int16(utcoffset)
	offset += 2

	// Reserved byte
	offset++

	t.GMPriority1 = b[offset]
	offset++

	err = t.GMClockQuality.UnmarshalBinary(b[offset : offset+ClockQualityPayloadLen])
	if err != nil {
		return err
	}
	offset += ClockQualityPayloadLen

	t.GMPriority2 = b[offset]
	offset++

	t.GMIdentity = binary.BigEndian.Uint64(b[offset : offset+GrandMasterIdentityLen])
	offset += GrandMasterIdentityLen

	t.StepsRemoved = binary.BigEndian.Uint16(b[offset : offset+StepsRemovedLen])
	offset += StepsRemovedLen

	t.TimeSource = TimeSourceType(b[offset])
	if !isValidTimeSource(t.TimeSource) {
		return ErrInvalidTimeSource
	}

	return nil
}
