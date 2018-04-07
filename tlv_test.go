package ptp

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestMarshalPathTraceTlv(t *testing.T) {
	var tests = []struct {
		desc string
		m    *PathTraceTlv
		b    []byte
		err  error
	}{
		{
			desc: "Empty pathSequence",
			m: &PathTraceTlv{
				pathSequence: []uint64{},
			},
			b: append([]byte{0x0, 0x8, 0x0, 0x0}),
		},
		{
			desc: "Single clockId in pathSequence",
			m: &PathTraceTlv{
				pathSequence: []uint64{0x0011223344556677},
			},
			b: append([]byte{0x0, 0x8, 0x0, 0x8,
				0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}),
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

func TestUnmarshalPathTraceTlv(t *testing.T) {
	var tests = []struct {
		desc string
		m    *PathTraceTlv
		b    []byte
		err  error
	}{
		{
			desc: "Empty pathSequence",
			b:    append([]byte{0x0, 0x8, 0x0, 0x0}),
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "Single clockId in pathSequence",
			m: &PathTraceTlv{
				pathSequence: []uint64{0x0011223344556677},
			},
			b: append([]byte{0x0, 0x8, 0x0, 0x8,
				0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}),
		},
		{
			desc: "Multiple clockId in pathSequence",
			m: &PathTraceTlv{
				pathSequence: []uint64{0x0011223344556677,
					0x5544226677889911},
			},
			b: append([]byte{0x0, 0x8, 0x0, 0x10,
				0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
				0x55, 0x44, 0x22, 0x66, 0x77, 0x88, 0x99, 0x11}),
		},
		{
			desc: "Invalid TLV type",
			b: append([]byte{0x1, 0x8, 0x0, 0x8,
				0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}),
			err: ErrInvalidTlvType,
		},
		{
			desc: "Mismatch lengthField and actual amount of bytes",
			b: append([]byte{0x0, 0x8, 0x0, 0x8,
				0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
				0x55, 0x44, 0x22, 0x66, 0x77, 0x88, 0x99, 0x11}),
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "The number of bytes is not a multiple of 8",
			b: append([]byte{0x0, 0x8, 0x0, 0x9,
				0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
				0x55}),
			err: io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			m := new(PathTraceTlv)
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

func TestMarshalIntervalRequestTlv(t *testing.T) {
	var tests = []struct {
		desc string
		m    *IntervalRequestTlv
		b    []byte
		err  error
	}{
		{
			desc: "Correct TLV values",
			m: &IntervalRequestTlv{
				LinkDelayInterval:        127,
				TimeSyncInterval:         127,
				AnnounceInterval:         127,
				ComputeNeighborRateRatio: true,
				ComputeNeighborPropDelay: false,
			},
			b: append([]byte{0x0, 0x3, 0x0, 0xc,
				0x0, 0x80, 0xc2, 0x0, 0x0, 0x2,
				0x7f, 0x7f, 0x7f,
				// Flags
				0x2,
				// Reserved
				0x0, 0x0}),
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

func TestUnmarshalIntervalRequestTlv(t *testing.T) {
	var tests = []struct {
		desc string
		m    *IntervalRequestTlv
		b    []byte
		err  error
	}{
		{
			desc: "Correct TLV values",
			m: &IntervalRequestTlv{
				LinkDelayInterval:        127,
				TimeSyncInterval:         127,
				AnnounceInterval:         127,
				ComputeNeighborRateRatio: true,
				ComputeNeighborPropDelay: false,
			},
			b: append([]byte{0x0, 0x3, 0x0, 0xc,
				0x0, 0x80, 0xc2, 0x0, 0x0, 0x2,
				0x7f, 0x7f, 0x7f,
				// Flags
				0x2,
				// Reserved
				0x0, 0x0}),
		},
		{
			desc: "Invalid TLV type",
			b: append([]byte{0x1, 0x3, 0x0, 0xc,
				0x0, 0x80, 0xc2, 0x0, 0x0, 0x2,
				0x7f, 0x7f, 0x7f,
				// Flags
				0x2,
				// Reserved
				0x0, 0x0}),
			err: ErrInvalidTlvType,
		},
		{
			desc: "Invalid organizationId",
			b: append([]byte{0x0, 0x3, 0x0, 0xc,
				0x1, 0x80, 0xc2, 0x0, 0x0, 0x2,
				0x7f, 0x7f, 0x7f,
				// Flags
				0x2,
				// Reserved
				0x0, 0x0}),
			err: ErrInvalidTlvOrgId,
		},
		{
			desc: "Invalid organizationSubType",
			b: append([]byte{0x0, 0x3, 0x0, 0xc,
				0x0, 0x80, 0xc2, 0x0, 0x0, 0x1,
				0x7f, 0x7f, 0x7f,
				// Flags
				0x2,
				// Reserved
				0x0, 0x0}),
			err: ErrInvalidTlvOrgSubType,
		},
		{
			desc: "Mismatch lengthField and actual amount of bytes",
			b: append([]byte{0x0, 0x3, 0x0, 0xa,
				0x0, 0x80, 0xc2, 0x0, 0x0, 0x2,
				0x7f, 0x7f, 0x7f,
				// Flags
				0x2,
				// Reserved
				0x0, 0x0}),
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "Invalid amount of bytes",
			b: append([]byte{0x0, 0x3, 0x0, 0xa,
				0x0, 0x80, 0xc2, 0x0, 0x0, 0x2,
				0x7f, 0x7f, 0x7f,
				// Flags
				0x2,
				// Reserved
				0x0}),
			err: io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			m := new(IntervalRequestTlv)
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

func TestMarshalFollowUpTlv(t *testing.T) {
	var tests = []struct {
		desc string
		m    *FollowUpTlv
		b    []byte
		err  error
	}{
		{
			desc: "Correct TLV values",
			m: &FollowUpTlv{
				CumulativeScaledRateOffset: 1,
				GmTimeBaseIndicator:        2,
				LastGmPhaseChange:          UScaledNs{1, 2},
				ScaledLastGmFreqChange:     7,
			},
			b: append([]byte{0x0, 0x3, 0x0, 0x1c,
				0x0, 0x80, 0xc2, 0x0, 0x0, 0x1,
				// cumulativeScaledRateOffset
				0x0, 0x0, 0x0, 0x1,
				// gmTimeBaseIndicator
				0x0, 0x2,
				// lastGmPhaseChange
				0x0, 0x0, 0x0, 0x1,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2,
				// scaledLastGmFreqChange
				0x0, 0x0, 0x0, 0x7}),
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

func TestUnmarshalFollowUpTlv(t *testing.T) {
	var tests = []struct {
		desc string
		m    *FollowUpTlv
		b    []byte
		err  error
	}{
		{
			desc: "Correct TLV values",
			m: &FollowUpTlv{
				CumulativeScaledRateOffset: 1,
				GmTimeBaseIndicator:        2,
				LastGmPhaseChange:          UScaledNs{1, 2},
				ScaledLastGmFreqChange:     7,
			},
			b: append([]byte{0x0, 0x3, 0x0, 0x1c,
				0x0, 0x80, 0xc2, 0x0, 0x0, 0x1,
				// cumulativeScaledRateOffset
				0x0, 0x0, 0x0, 0x1,
				// gmTimeBaseIndicator
				0x0, 0x2,
				// lastGmPhaseChange
				0x0, 0x0, 0x0, 0x1,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2,
				// scaledLastGmFreqChange
				0x0, 0x0, 0x0, 0x7}),
		},
		{
			desc: "Invalid TLV type",
			b: append([]byte{0x0, 0x2, 0x0, 0x1c,
				0x0, 0x80, 0xc2, 0x0, 0x0, 0x1,
				// cumulativeScaledRateOffset
				0x0, 0x0, 0x0, 0x1,
				// gmTimeBaseIndicator
				0x0, 0x2,
				// lastGmPhaseChange
				0x0, 0x0, 0x0, 0x1,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2,
				// scaledLastGmFreqChange
				0x0, 0x0, 0x0, 0x7}),
			err: ErrInvalidTlvType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			m := new(FollowUpTlv)
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
