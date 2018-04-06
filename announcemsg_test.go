package ptp

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestMarshalAnnounce(t *testing.T) {
	var tests = []struct {
		desc string
		m    *AnnounceMsg
		b    []byte
		err  error
	}{
		{
			desc: "Correct structure",
			m: &AnnounceMsg{
				Header: Header{
					MessageType:      AnnounceMsgType,
					MessageLength:    HeaderLen + AnnouncePayloadLen,
					VersionPTP:       Version2,
					CorrectionNs:     0,
					CorrectionSubNs:  0,
					ClockIdentity:    0x000af7fffe42a753,
					PortNumber:       2,
					SequenceID:       55330,
					LogMessagePeriod: 0,
				},
				GMClockQuality: GMClockQuality{
					ClockClass:    PrimarySyncRefClass,
					ClockAccuracy: ClockAccuracy100ns,
					ClockVariance: 200,
				},
				CurrentUtcOffset: 36,
				GMPriority1:      128,
				GMPriority2:      128,
				GMIdentity:       0x001d7ffffe80024a,
				StepsRemoved:     0,
				TimeSource:       TimeSourceGPS,
				PathTraceTlv:     PathTraceTlv{},
			},
			b: append([]byte{0xb, 0x2, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				// ClockIdentity
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0x0,
				// Message body
				// Reserved 10 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x24,
				0x0, 0x80,
				0x6, 0x21, 0x0, 0xc8, 0x80,
				// GM Identity
				0x0, 0x1d, 0x7f, 0xff, 0xfe, 0x80, 0x2, 0x4a,
				0x0, 0x0, 0x20,
				// PathTraceTlv
				0x0, 0x8, 0x0, 0x0,
			}),
		},
		{
			desc: "Invalid message type",
			m: &AnnounceMsg{
				Header: Header{
					MessageType:      SyncMsgType,
					MessageLength:    HeaderLen + AnnouncePayloadLen,
					VersionPTP:       Version2,
					CorrectionNs:     0,
					CorrectionSubNs:  0,
					ClockIdentity:    0x000af7fffe42a753,
					PortNumber:       2,
					SequenceID:       55330,
					LogMessagePeriod: 0,
				},
				GMClockQuality: GMClockQuality{
					ClockClass:    PrimarySyncRefClass,
					ClockAccuracy: ClockAccuracy100ns,
					ClockVariance: 200,
				},
				CurrentUtcOffset: 36,
				GMPriority1:      128,
				GMPriority2:      128,
				GMIdentity:       0x001d7ffffe80024a,
				StepsRemoved:     0,
				TimeSource:       TimeSourceGPS,
				PathTraceTlv:     PathTraceTlv{},
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

func TestUnmarshalAnnounce(t *testing.T) {
	var tests = []struct {
		desc string
		m    *AnnounceMsg
		b    []byte
		err  error
	}{
		{
			desc: "Correct structure",
			m: &AnnounceMsg{
				Header: Header{
					MessageType:      AnnounceMsgType,
					MessageLength:    HeaderLen + AnnouncePayloadLen,
					VersionPTP:       Version2,
					CorrectionNs:     0,
					CorrectionSubNs:  0,
					ClockIdentity:    0x000af7fffe42a753,
					PortNumber:       2,
					SequenceID:       55330,
					LogMessagePeriod: 0,
				},
				GMClockQuality: GMClockQuality{
					ClockClass:    PrimarySyncRefClass,
					ClockAccuracy: ClockAccuracy100ns,
					ClockVariance: 200,
				},
				CurrentUtcOffset: 36,
				GMPriority1:      128,
				GMPriority2:      128,
				GMIdentity:       0x001d7ffffe80024a,
				StepsRemoved:     0,
				TimeSource:       TimeSourceGPS,
			},
			b: append([]byte{0xb, 0x2, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				// ClockIdentity
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0x0,
				// Message body
				// Reserved 10 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x24,
				0x0, 0x80,
				0x6, 0x21, 0x0, 0xc8, 0x80,
				// GM Identity
				0x0, 0x1d, 0x7f, 0xff, 0xfe, 0x80, 0x2, 0x4a,
				0x0, 0x0, 0x20,
			}),
		},
		{
			desc: "Invalid clock class",
			b: append([]byte{0xb, 0x2, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc,
				// Message body
				// Reserved 10 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x24,
				0x0, 0x80,
				// Wrong class - 0x8 instead of {6,7,248,255}
				0x8, 0x21, 0x0, 0xc8, 0x80,
				// GM Identity
				0x0, 0x1d, 0x7f, 0xff, 0xfe, 0x80, 0x2, 0x4a,
				0x0, 0x0, 0x20,
			}),
			err: ErrInvalidClockClass,
		},
		{
			desc: "Invalid clock accuracy",
			b: append([]byte{0xb, 0x2, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc,
				// Message body
				// Reserved 10 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x24,
				0x0, 0x80,
				// Wrong class accuracy - 0x1
				0x6, 0x1, 0x0, 0xc8, 0x80,
				// GM Identity
				0x0, 0x1d, 0x7f, 0xff, 0xfe, 0x80, 0x2, 0x4a,
				0x0, 0x0, 0x20,
			}),
			err: ErrInvalidClockAccuracy,
		},
		{
			desc: "Invalid time source",
			b: append([]byte{0xb, 0x2, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc,
				// Message body
				// Reserved 10 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x24,
				0x0, 0x80,
				// Wrong class accuracy - 0x1
				0x6, 0x21, 0x0, 0xc8, 0x80,
				// GM Identity
				0x0, 0x1d, 0x7f, 0xff, 0xfe, 0x80, 0x2, 0x4a,
				0x0, 0x0, 0x21,
			}),
			err: ErrInvalidTimeSource,
		},
		{
			desc: "Invalid time source",
			b: append([]byte{0xb, 0x2, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0xfc,
				// Message body
				// Reserved 10 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x24,
				0x0, 0x80,
				// Wrong class accuracy - 0x1
				0x6, 0x21, 0x0, 0xc8, 0x80,
				// GM Identity
				0x0, 0x1d, 0x7f, 0xff, 0xfe, 0x80, 0x2, 0x4a,
			}),
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "Invalid message type",
			b: append([]byte{0x1, 0x2, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				// ClockIdentity
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0x0,
				// Message body
				// Reserved 10 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x24,
				0x0, 0x80,
				0x6, 0x21, 0x0, 0xc8, 0x80,
				// GM Identity
				0x0, 0x1d, 0x7f, 0xff, 0xfe, 0x80, 0x2, 0x4a,
				0x0, 0x0, 0x20,
			}),
			err: ErrInvalidMsgType,
		},
		{
			desc: "Invalid length",
			b: append([]byte{0x1, 0x2, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				// ClockIdentity
				0x0, 0xa, 0xf7, 0xff, 0xfe, 0x42, 0xa7, 0x53, 0x0, 0x2, 0xd8, 0x22, 0x5, 0x0,
				// Message body
				// Reserved 10 bytes
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x24,
				0x0, 0x80,
				0x6, 0x21, 0x0, 0xc8, 0x80,
				// GM Identity
				0x0, 0x1d, 0x7f, 0xff, 0xfe, 0x80, 0x2, 0x4a,
				0x0, 0x0, 0x20, 0x0,
			}),
			err: ErrInvalidMsgType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			m := new(AnnounceMsg)
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
