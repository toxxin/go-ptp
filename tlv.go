package ptp

import (
	"encoding/binary"
	"io"
)

// TlvType Type
type TlvType uint16

// Tlv types codes
const (
	// Standard TLVs
	Management            TlvType = 0x0001
	ManagementErrorStatus TlvType = 0x0002
	OrganizationExtension TlvType = 0x0003

	// Optional unicast message negotiation TLVs
	RequestUnicastTransmission           TlvType = 0x0004
	GrantUnicastTransmission             TlvType = 0x0005
	CancelUnicastTransmission            TlvType = 0x0006
	AcknowledgeCancelUnicastTransmission TlvType = 0x0007

	// Optional path trace mechanism TLV
	PathTrace TlvType = 0x0008

	// Optional alternate timescale TLV
	AlternateTimeOffsetIndicator TlvType = 0x0009

	// Reserved for standard TLVs
	// 000A – 1FFF

	// Security TLVs
	Authentication            TlvType = 0x2000
	AuthenticationChallenge   TlvType = 0x2001
	SecurityAssociationUpdate TlvType = 0x2002

	// Cumulative frequency scale factor offset
	CumFreqScaleFactorOffset TlvType = 0x2003

	// Reserved for Experimental TLVs
	// 2004 – 3FFF

	// Reserved
	// 4000 – FFFF
)

// PathTraceTlv ...
type PathTraceTlv struct {
	pathSequence []uint64
}

// MarshalBinary allocates a byte slice and marshals a Header into binary form.
func (p *PathTraceTlv) MarshalBinary() ([]byte, error) {

	b := make([]byte, 4+8*len(p.pathSequence))

	// TLV type
	binary.BigEndian.PutUint16(b[:2], uint16(PathTrace))

	// TLV length
	binary.BigEndian.PutUint16(b[2:4], uint16(8*len(p.pathSequence)))

	for i, v := range p.pathSequence {
		binary.BigEndian.PutUint64(b[4+i*8:4+i*8+8], v)
	}

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a PathTraceTlv.
//
// If the byte slice does not contain enough data to unmarshal a valid PDelReqMsg,
// io.ErrUnexpectedEOF is returned.
func (p *PathTraceTlv) UnmarshalBinary(b []byte) error {
	if len(b) < (2 + 2 + ClockIdentityLen) {
		return io.ErrUnexpectedEOF
	}

	// Length of b must be 2 + 2 + 8N
	// Сheck the remainder of division by 8
	tlvLen := binary.BigEndian.Uint16(b[2:4])
	if int(tlvLen) != len(b[4:]) || tlvLen%8 != 0 {
		return io.ErrUnexpectedEOF
	}

	pathSeq := make([]uint64, tlvLen/8)
	for i := range pathSeq {
		pathSeq[i] = binary.BigEndian.Uint64(b[i*8+4 : i*8+8+4])
	}

	p.pathSequence = pathSeq

	return nil
}
