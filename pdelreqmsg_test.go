package ptp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMarshalPDelReq(t *testing.T) {

	var tests = []struct {
		desc string
		m    *PDelReqMsg
		b    []byte
		err  error
	}{
		{
			desc: "Correct structure",
			m: &PDelReqMsg{
				Header: Header{
					MessageType:     PDelayReqMsgType,
					CorrectionNs:    0,
					CorrectionSubNs: 0,
					Flags: Flags{
						Security:           false,
						ProfileSpecific2:   false,
						ProfileSpecific1:   true,
						Unicast:            false,
						TwoSteps:           false,
						AlternateMaster:    false,
						FrequencyTraceable: false,
						TimeTraceable:      false,
						UtcReasonable:      false,
						LI59:               false,
						LI61:               false,
					},
					ClockIdentity:    0x000af7fffe42a753,
					PortNumber:       2,
					SequenceID:       55330,
					LogMessagePeriod: -4,
				},
			},
			b: append([]byte{0x2, 0x2, 0x0, 0x36, 0x0, 0x0, 0x20, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc,
				// Reserved 20 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
		},
		{
			desc: "Invalid message type",
			m: &PDelReqMsg{
				Header: Header{
					MessageType:      SyncMsgType,
					CorrectionNs:     0,
					CorrectionSubNs:  0,
					ClockIdentity:    0x000af7fffe42a753,
					PortNumber:       2,
					SequenceID:       55330,
					LogMessagePeriod: -4,
				},
			},
			err: ErrInvalidMsgType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			b, err := tt.m.MarshalBinary()
			if err != nil {
				if want, got := tt.err, err; want != got {
					t.Fatalf("unexpected error: %v != %v", want, got)
				}

				return
			}

			if want, got := tt.b, b; !bytes.Equal(want, got) {
				t.Fatalf("unexpected Frame bytes:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}

func TestUnmarshalPDelReq(t *testing.T) {

	var tests = []struct {
		desc string
		m    PDelReqMsg
		b    []byte
		err  error
	}{
		{
			desc: "Correct structure",
			m: PDelReqMsg{
				Header: Header{
					MessageType:   PDelayReqMsgType,
					MessageLength: HeaderLen + PDelayReqPayloadLen,
					VersionPTP:    Version2,
					Flags: Flags{
						Security:           false,
						ProfileSpecific2:   false,
						ProfileSpecific1:   false,
						Unicast:            false,
						TwoSteps:           false,
						AlternateMaster:    false,
						FrequencyTraceable: false,
						TimeTraceable:      false,
						UtcReasonable:      false,
						LI59:               false,
						LI61:               false,
					},
					CorrectionNs:     0,
					CorrectionSubNs:  0,
					ClockIdentity:    0x000af7fffe42a753,
					PortNumber:       2,
					SequenceID:       55330,
					LogMessagePeriod: -4,
				},
			},
			b: append([]byte{0x2, 0x2, 0x0, 0x36, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc,
				// Reserved 20 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var m PDelReqMsg
			err := m.UnmarshalBinary(tt.b)
			if err != nil {
				if want, got := tt.err, err; want != got {
					t.Fatalf("unexpected error: %v != %v", want, got)
				}

				return
			}
			if want, got := tt.m, m; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Frame bytes:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}
