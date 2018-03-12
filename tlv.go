package ptp

import (
	"encoding/binary"
	"io"
)

// TLV payload length
const (
	FollowUpTlvLen        = 28
	IntervalRequestTlvLen = 12
	CsnTlvLen             = 46
)

var organizationID = []byte{0x0, 0x80, 0xc2}

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

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
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
// If the byte slice does not contain enough data to unmarshal a valid PathTraceTlv,
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

// IntervalRequestTlv ...
type IntervalRequestTlv struct {
	// OrganizationSubType = 2
	LinkDelayInterval        int8
	TimeSyncInterval         int8
	AnnounceInterval         int8
	ComputeNeighborRateRatio bool
	ComputeNeighborPropDelay bool
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (p *IntervalRequestTlv) MarshalBinary() ([]byte, error) {

	b := make([]byte, IntervalRequestTlvLen+4)

	// TLV type
	binary.BigEndian.PutUint16(b[:2], uint16(OrganizationExtension))

	// TLV length
	binary.BigEndian.PutUint16(b[2:4], uint16(IntervalRequestTlvLen))

	copy(b[4:7], organizationID)

	// organizationSubType
	copy(b[7:10], []byte{0x0, 0x0, 0x2})

	b[10] = uint8(p.LinkDelayInterval)

	b[11] = uint8(p.TimeSyncInterval)

	b[12] = uint8(p.AnnounceInterval)

	b[13] = uint8(b2i(p.ComputeNeighborRateRatio)<<1 | b2i(p.ComputeNeighborPropDelay)<<2)

	return b, nil
}

// FollowUpTlv ...
type FollowUpTlv struct {
	// OrganizationSubType = 3
	CumulativeScaledRateOffset int32
	GmTimeBaseIndicator        uint16
	lastGmPhaseChange          UScaledNs
	scaledLastGmFreqChange     int32
}

// UnmarshalBinary unmarshals a byte slice into a IntervalRequestTlv.
//
// If the byte slice does not contain enough data to unmarshal a valid IntervalRequestTlv,
// io.ErrUnexpectedEOF is returned.
func (p *IntervalRequestTlv) UnmarshalBinary(b []byte) error {
	if len(b) != (IntervalRequestTlvLen + 4) {
		return io.ErrUnexpectedEOF
	}

	tlvLen := binary.BigEndian.Uint16(b[2:4])
	if int(tlvLen) != IntervalRequestTlvLen {
		return io.ErrUnexpectedEOF
	}

	// TODO: check message type

	p.LinkDelayInterval = int8(b[10])

	p.TimeSyncInterval = int8(b[11])

	p.AnnounceInterval = int8(b[12])

	p.ComputeNeighborRateRatio = (b[13] & 0x2) != 0
	p.ComputeNeighborPropDelay = (b[13] & 0x4) != 0

	return nil
}

// CsnTlv ...
type CsnTlv struct {
	UpstreamTxTime    UScaledNs
	NeighborRateRatio int32
	NeighborPropDelay UScaledNs
	DelayAsymmetry    UScaledNs
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (p *CsnTlv) MarshalBinary() ([]byte, error) {

	b := make([]byte, CsnTlvLen+4)

	// TLV type
	binary.BigEndian.PutUint16(b[:2], uint16(OrganizationExtension))

	// TLV length
	binary.BigEndian.PutUint16(b[2:4], uint16(IntervalRequestTlvLen))

	copy(b[4:7], organizationID)

	// organizationSubType
	copy(b[7:10], []byte{0x0, 0x0, 0x3})

	tx, err := p.UpstreamTxTime.MarshalBinary()
	if err != nil {
		return nil, err
	}

	offset := 10

	copy(b[offset:offset+UScaledNsLen], tx)
	offset += UScaledNsLen

	binary.BigEndian.PutUint32(b[22:26], uint32(p.NeighborRateRatio))
	offset += 4

	nd, err := p.NeighborPropDelay.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[offset:offset+UScaledNsLen], nd)
	offset += UScaledNsLen

	da, err := p.DelayAsymmetry.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(b[offset:offset+UScaledNsLen], da)

	return b, nil
}
