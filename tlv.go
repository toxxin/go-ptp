package ptp

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
