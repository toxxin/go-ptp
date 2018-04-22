package ptp

import (
	"bytes"
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

	tlvType := TlvType(binary.BigEndian.Uint16(b[0:2]))
	if tlvType != PathTrace {
		return ErrInvalidTlvType
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

	tlvType := TlvType(binary.BigEndian.Uint16(b[0:2]))
	if tlvType != OrganizationExtension {
		return ErrInvalidTlvType
	}

	if !bytes.Equal(b[4:7], organizationID) {
		return ErrInvalidTlvOrgId
	}

	// The value of organizationSubType is 2
	if !bytes.Equal([]byte{0x0, 0x0, 0x2}, b[7:10]) {
		return ErrInvalidTlvOrgSubType
	}

	p.LinkDelayInterval = int8(b[10])

	p.TimeSyncInterval = int8(b[11])

	p.AnnounceInterval = int8(b[12])

	p.ComputeNeighborRateRatio = (b[13] & 0x2) != 0
	p.ComputeNeighborPropDelay = (b[13] & 0x4) != 0

	return nil
}

// FollowUpTlv ...
type FollowUpTlv struct {
	// OrganizationSubType = 3
	CumulativeScaledRateOffset int32
	GmTimeBaseIndicator        uint16
	LastGmPhaseChange          UScaledNs
	ScaledLastGmFreqChange     int32
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (p *FollowUpTlv) MarshalBinary() ([]byte, error) {

	b := make([]byte, FollowUpTlvLen+4)

	// TLV type
	binary.BigEndian.PutUint16(b[:2], uint16(OrganizationExtension))

	// TLV length
	binary.BigEndian.PutUint16(b[2:4], uint16(FollowUpTlvLen))

	copy(b[4:7], organizationID)

	// organizationSubType
	copy(b[7:10], []byte{0x0, 0x0, 0x1})
	offset := 10

	binary.BigEndian.PutUint32(b[offset:offset+4], uint32(p.CumulativeScaledRateOffset))
	offset += 4

	binary.BigEndian.PutUint16(b[offset:offset+2], p.GmTimeBaseIndicator)
	offset += 2

	lastGM, err := p.LastGmPhaseChange.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(b[offset:offset+UScaledNsLen], lastGM)
	offset += UScaledNsLen

	binary.BigEndian.PutUint32(b[offset:offset+4], uint32(p.ScaledLastGmFreqChange))

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a FollowUpTlv.
//
// If the byte slice does not contain enough data to unmarshal a valid FollowUpTlv,
// io.ErrUnexpectedEOF is returned.
func (p *FollowUpTlv) UnmarshalBinary(b []byte) error {

	var err error

	if len(b) != (FollowUpTlvLen + 4) {
		return io.ErrUnexpectedEOF
	}

	tlvLen := binary.BigEndian.Uint16(b[2:4])
	if int(tlvLen) != FollowUpTlvLen {
		return io.ErrUnexpectedEOF
	}

	tlvType := TlvType(binary.BigEndian.Uint16(b[0:2]))
	if tlvType != OrganizationExtension {
		return ErrInvalidTlvType
	}

	if !bytes.Equal(b[4:7], organizationID) {
		return ErrInvalidTlvOrgId
	}

	// The value of organizationSubType is 1
	if !bytes.Equal([]byte{0x0, 0x0, 0x1}, b[7:10]) {
		return ErrInvalidTlvOrgSubType
	}

	p.CumulativeScaledRateOffset = int32(binary.BigEndian.Uint32(b[10:14]))

	p.GmTimeBaseIndicator = binary.BigEndian.Uint16(b[14:16])

	p.LastGmPhaseChange, err = NewUScaledNs(b[16:28])
	if err != nil {
		return err
	}

	p.ScaledLastGmFreqChange = int32(binary.BigEndian.Uint32(b[28:32]))

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
	binary.BigEndian.PutUint16(b[2:4], uint16(CsnTlvLen))

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

	binary.BigEndian.PutUint32(b[offset:offset+4], uint32(p.NeighborRateRatio))
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

// UnmarshalBinary unmarshals a byte slice into a CsnTlv frame.
//
// If the byte slice does not contain enough data to unmarshal a valid CsnTlv frame,
// io.ErrUnexpectedEOF is returned.
func (p *CsnTlv) UnmarshalBinary(b []byte) error {

	var err error

	if len(b) != (CsnTlvLen + 4) {
		return io.ErrUnexpectedEOF
	}

	tlvLen := binary.BigEndian.Uint16(b[2:4])
	if int(tlvLen) != CsnTlvLen {
		return io.ErrUnexpectedEOF
	}

	tlvType := TlvType(binary.BigEndian.Uint16(b[0:2]))
	if tlvType != OrganizationExtension {
		return ErrInvalidTlvType
	}

	if !bytes.Equal(b[4:7], organizationID) {
		return ErrInvalidTlvOrgId
	}

	// The value of organizationSubType is 3
	if !bytes.Equal([]byte{0x0, 0x0, 0x3}, b[7:10]) {
		return ErrInvalidTlvOrgSubType
	}

	offset := 10

	err = p.UpstreamTxTime.UnmarshalBinary(b[offset : offset+UScaledNsLen])
	if err != nil {
		return err
	}

	offset += UScaledNsLen

	p.NeighborRateRatio = int32(binary.BigEndian.Uint32(b[offset : offset+4]))
	offset += 4

	err = p.NeighborPropDelay.UnmarshalBinary(b[offset : offset+UScaledNsLen])
	if err != nil {
		return err
	}

	offset += UScaledNsLen

	err = p.DelayAsymmetry.UnmarshalBinary(b[offset : offset+UScaledNsLen])
	if err != nil {
		return err
	}

	return nil
}

type ManagementIdType uint16

const (
	// Applicable to all node types 0000 – 1FFF
	NullManagement           ManagementIdType = 0x0000
	ClockDescription         ManagementIdType = 0x0001
	UserDescription          ManagementIdType = 0x0002
	SaveInNonVolatileStorage ManagementIdType = 0x0003
	ResetNonVolatileStorage  ManagementIdType = 0x0004
	Initialize               ManagementIdType = 0x0005
	FaultLog                 ManagementIdType = 0x0006
	FaultLogReset            ManagementIdType = 0x0007

	// Reserved 0008 – 1FFF

	// Applicable to ordinary and boundary clocks 2000 – 3FFF
	DefaultDataSet                ManagementIdType = 0x2000
	CurrentDataSet                ManagementIdType = 0x2001
	ParentDataSet                 ManagementIdType = 0x2002
	TimePropertiesDataSet         ManagementIdType = 0x2003
	PortDataSet                   ManagementIdType = 0x2004
	Priority1                     ManagementIdType = 0x2005
	Priority2                     ManagementIdType = 0x2006
	Domain                        ManagementIdType = 0x2007
	SlaveOnly                     ManagementIdType = 0x2008
	LogAnnounceInterval           ManagementIdType = 0x2009
	AnnounceReceiptTimeout        ManagementIdType = 0x200a
	LogSyncInterval               ManagementIdType = 0x200b
	VersionNumber                 ManagementIdType = 0x200c
	EneablePort                   ManagementIdType = 0x200d
	DisablePort                   ManagementIdType = 0x200e
	Time                          ManagementIdType = 0x200f
	ClockAccuracy                 ManagementIdType = 0x2010
	UtcProperties                 ManagementIdType = 0x2011
	TraceabilityProperties        ManagementIdType = 0x2012
	TimescaleProperties           ManagementIdType = 0x2013
	UnicastNegotiationEnable      ManagementIdType = 0x2014
	PathTraceList                 ManagementIdType = 0x2015
	PathTraceEnable               ManagementIdType = 0x2016
	GrandMasterClusterTable       ManagementIdType = 0x2017
	UnicastMasterTable            ManagementIdType = 0x2018
	UnicastMasterMaxTableSize     ManagementIdType = 0x2019
	AcceptableMasterTable         ManagementIdType = 0x201a
	AcceptableMasterTableEnabled  ManagementIdType = 0x201b
	AcceptableMasterMaxTableSize  ManagementIdType = 0x201c
	AlternateMaster               ManagementIdType = 0x201d
	AlternateTimeOffsetEnable     ManagementIdType = 0x201e
	AlternateTimeOffsetName       ManagementIdType = 0x201f
	AlternateTimeOffsetMaxKey     ManagementIdType = 0x2020
	AlternateTimeOffsetProperties ManagementIdType = 0x2021

	// Reserved 2022 – 3FFF

	// Applicable to transparent clocks 4000 – 5FFF
	TransparentClockDefaultDataSet ManagementIdType = 0x4000
	TransparentClockPortDataSet    ManagementIdType = 0x4001
	PrimaryDomain                  ManagementIdType = 0x4002

	// Reserved 4003 – 5FFF

	// Applicable to ordinary, boundary, and transparent clocks 6000 – 7FFF
	DelayMechanism          ManagementIdType = 0x6000
	LogMinPdelayReqInterval ManagementIdType = 0x6001

	// Reserved 6002 – BFFF

	// This range is to be used for implementation-specific identifiers C000 – DFFF
	// This range is to be assigned by an alternate PTP profile E000 – FFFE

	// Reserved FFFF

)

type ClockType uint16

const (
	OrdinaryClock              ClockType = 0
	boundaryClock              ClockType = 1
	PeerToPeerTransparentClock ClockType = 2
	EndToEndTransparentClock   ClockType = 3
	ManagementNode             ClockType = 4
	// Reserved 5–F
)

// SeverityCode is FaultRecord.severityCode
type SeverityCode uint8

const (
	Emergency SeverityCode = 0
	Alert
	Critical
	Error
	Warning
	Notice
	Informational
	Debug
	// Reserved 08–FF
)

type ClockDescriptionTlv struct {
}

type UserDescriptionTlv struct {
}

type SaveInNonVolatileStorageTlv struct {
}

type ResetNonVolatileStorageTlv struct {
}

type InitializeTlv struct {
	InitializationKey uint16
}

type FaultLogTlv struct {
}

type FaultLogResetTlv struct {
}

type DefaultDataSetTlv struct {
	TSC         bool
	SO          bool
	NumberPorts uint16
	Priority1   uint8
	ClockQuality
	Priority2     uint8
	ClockIdentity uint64
	DomainNumber  uint8
}

type CurrentDataSetTlv struct {
}

type ParentDataSetTlv struct {
}

type TimePropertiesDataSetTlv struct {
}

type PortDataSetTlv struct {
}

type Priority1Tlv struct {
	Priority1 uint8
	// Reserved 1byte
}

type Priority2Tlv struct {
	Priority2 uint8
	// Reserved 1byte
}

type DomainTlv struct {
	DomainNumber uint8
	// Reserved 1byte
}

type SlaveOnlyTlv struct {
}

type LogAnnounceIntervalTlv struct {
}

type AnnounceReceiptTimeoutTlv struct {
}

type LogSyncIntervalTlv struct {
}

type VersionNumberTlv struct {
}

type TimeTlv struct {
}

type ClockAccuracyTlv struct {
	clockAccuracy ClockAccuracyType
}

type UtcPropertiesTlv struct {
}

type TraceabilityPropertiesTlv struct {
}

type TimescalePropertiesTlv struct {
}

type EneablePortTlv struct {
}

type DisablePortTlv struct {
}

type TransparentClockDefaultDataSetTlv struct {
}

type TransparentClockPortDataSetTlv struct {
}

type PrimaryDomainTlv struct {
}

type DelayMechanismTlv struct {
}

type LogMinPdelayReqIntervalTlv struct {
}

type ManagementTlv struct {
}
