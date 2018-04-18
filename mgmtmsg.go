package ptp

// ActionFieldType...
type ActionFiledType uint8

// ActionFieldType types codes
const (
	Get         ActionFiledType = 0
	Set         ActionFiledType = 1
	Response    ActionFiledType = 2
	Command     ActionFiledType = 3
	Acknowledge ActionFiledType = 4
)

// MgmtMsg ...
type MgmtMsg struct {
	Header
	ClockIdentity        uint64
	PortNumber           uint16
	StartingBoundaryHops uint8
	BoundaryHops         uint8
	actionField          ActionFiledType
	ManagementTlv
}
