package ptp

import (
	"bytes"
	"testing"
)

func TestMarshalSignaling(t *testing.T) {
	var tests = []struct {
		desc string
		m    *SignalingMsg
		b    []byte
		err  error
	}{
		{
			desc: "Correct structure",
			m: &SignalingMsg{
				Header: Header{
					MessageType:      SignalingMsgType,
					VersionPTP:       Version2,
					CorrectionNs:     0,
					CorrectionSubNs:  0,
					ClockIdentity:    0x001d7ffffe80024a,
					PortNumber:       1,
					SequenceID:       27278,
					LogMessagePeriod: 127,
				},
				ClockIdentity: 0x78baf9fffe0a435e,
				PortNumber:    1,
				IntervalRequestTlv: IntervalRequestTlv{
					LinkDelayInterval:        1,
					TimeSyncInterval:         2,
					AnnounceInterval:         127,
					ComputeNeighborRateRatio: true,
					ComputeNeighborPropDelay: false,
				},
			},
			b: append([]byte{
				0x0c, 0x02, 0x0, 0x38, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1d,
				0x7f, 0xff, 0xfe, 0x80, 0x02, 0x4a, 0x00, 0x01, 0x6a, 0x8e, 0x05, 0x7f, 0x78, 0xba, 0xf9, 0xff,
				0xfe, 0x0a, 0x43, 0x5e, 0x00, 0x01, 0x00, 0x03, 0x00, 0x0c, 0x0, 0x80, 0xc2, 0x0, 0x0, 0x2, 0x1, 0x2, 0x7f,
				0x2, 0x0, 0x0,
			}),
		},
		{
			desc: "Invalid message type",
			m: &SignalingMsg{
				Header: Header{
					MessageType:      AnnounceMsgType,
					VersionPTP:       Version2,
					CorrectionNs:     0,
					CorrectionSubNs:  0,
					ClockIdentity:    0x001d7ffffe80024a,
					PortNumber:       1,
					SequenceID:       27278,
					LogMessagePeriod: 127,
				},
				ClockIdentity: 0x78baf9fffe0a435e,
				PortNumber:    1,
				IntervalRequestTlv: IntervalRequestTlv{
					LinkDelayInterval:        1,
					TimeSyncInterval:         2,
					AnnounceInterval:         127,
					ComputeNeighborRateRatio: true,
					ComputeNeighborPropDelay: false,
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
