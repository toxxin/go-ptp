package ptp

import (
	"bytes"
	"testing"
	"time"
)

func TestMarshalPDelResp(t *testing.T) {

	var tests = []struct {
		desc string
		m    *PDelRespMsg
		b    []byte
		err  error
	}{
		{
			desc: "Correct structure",
			m: &PDelRespMsg{
				Header: Header{
					MessageType:     PDelayRespMsgType,
					CorrectionNs:    0,
					CorrectionSubNs: 0,
					Flags: Flags{
						TwoSteps: true,
					},
					ClockIdentity:    0x0023aefffe5d688b,
					PortNumber:       1,
					SequenceID:       6365,
					LogMessagePeriod: 127,
				},
				ReceiveTimestamp: time.Unix(1312261115, 89388000),
				ClockIdentity:    0x000c29fffe08e6e8,
				PortNumber:       1,
			},
			b: append([]byte{0x3, 0x2, 0x0, 0x36, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x23, 0xae, 0xff, 0xfe, 0x5d, 0x68, 0x8b, 0x00, 0x01, 0x18, 0xdd,
				0x05, 0x7f, 0x00, 0x00, 0x4e, 0x37, 0x83, 0xfb, 0x05, 0x53, 0xf3, 0xe0, 0x00, 0x0c, 0x29, 0xff,
				0xfe, 0x08, 0xe6, 0xe8, 0x00, 0x01}),
		},
		{
			desc: "Invalid message type",
			m: &PDelRespMsg{
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
