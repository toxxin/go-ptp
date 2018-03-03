package ptp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMarshalHeader(t *testing.T) {

	var tests = []struct {
		desc string
		h    *Header
		b    []byte
		err  error
	}{
		{
			desc: "Correct structure",
			h: &Header{
				MessageType:   PDelayReqMsgType,
				MessageLength: 44,
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
			b: append([]byte{0x2, 0x2, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			b, err := tt.h.MarshalBinary()
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

func TestUnmarshalHeader(t *testing.T) {

	var tests = []struct {
		desc string
		h    *Header
		b    []byte
		err  error
	}{
		{
			desc: "Correct structure",
			h: &Header{
				MessageType:   PDelayReqMsgType,
				MessageLength: 44,
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
			b: append([]byte{0x2, 0x2, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc}),
		},
		{
			desc: "Invalid message type",
			b: append([]byte{0x4, 0x2, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc}),
			err: ErrInvalidMsgType,
		},
		{
			desc: "Unsupported protocol version",
			b: append([]byte{0x2, 0x1, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc}),
			err: ErrUnsupportedVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			h := new(Header)
			err := h.UnmarshalBinary(tt.b)
			if err != nil {
				if want, got := tt.err, err; want != got {
					t.Fatalf("unexpected error: %v != %v", want, got)
				}

				return
			}

			if want, got := tt.h, h; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Frame bytes:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}
