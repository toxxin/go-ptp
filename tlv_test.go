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
