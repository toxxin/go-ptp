package ptp

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestTime2OriginTimestamp(t *testing.T) {

	var tests = []struct {
		desc string
		t    time.Time
		b    []byte
	}{
		{
			desc: "Correct not null timestamp",
			t:    time.Unix(1169232201, 775045731),
			b:    append([]byte{0x0, 0x0, 0x45, 0xb1, 0x11, 0x49, 0x2e, 0x32, 0x42, 0x63}),
		},
		{
			desc: "Null timestamp",
			t:    time.Unix(0, 0),
			b:    append([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			b := make([]byte, 10)
			time2OriginTimestamp(tt.t, b)

			if want, got := tt.b, b; !bytes.Equal(want, got) {
				t.Fatalf("unexpected Frame bytes:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}

func TestOriginTimestamp2Time(t *testing.T) {

	var tests = []struct {
		desc string
		t    time.Time
		b    []byte
		err  error
	}{
		{
			desc: "Correct not null timestamp",
			t:    time.Unix(1169232201, 775045731),
			b:    append([]byte{0x0, 0x0, 0x45, 0xb1, 0x11, 0x49, 0x2e, 0x32, 0x42, 0x63}),
		},
		{
			desc: "Null timestamp",
			t:    time.Unix(0, 0),
			b:    append([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tmp, err := originTimestamp2Time(tt.b)
			if err != nil {
				if want, got := tt.err, err; want != got {
					t.Fatalf("unexpected error: %v != %v", want, got)
				}

				return
			}

			if want, got := tt.t, tmp; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Frame bytes:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}
