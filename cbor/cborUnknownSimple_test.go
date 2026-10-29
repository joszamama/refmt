package cbor

import (
	"testing"

	. "github.com/polydawn/refmt/tok"
)

func TestCoerceUnknownSimpleToNull(t *testing.T) {
	unassigned := []struct {
		name    string
		payload []byte
	}{
		{"simple(0)", []byte{0xe0}},
		{"simple(19)", []byte{0xf3}},
		{"simple(32)", []byte{0xf8, 0x20}},
		{"simple(255)", []byte{0xf8, 0xff}},
	}
	for _, tc := range unassigned {
		t.Run(tc.name+" null when option set", func(t *testing.T) {
			tk, _, err := nextToken(t, DecodeOptions{CoerceUnknownSimpleToNull: true}, tc.payload)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tk.Type != TNull {
				t.Fatalf("expected null, got %v", tk)
			}
		})
		t.Run(tc.name+" rejected by default", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{}, tc.payload)
			if err == nil {
				t.Fatalf("expected an error")
			}
		})
	}

	malformed := []struct {
		name    string
		payload []byte
	}{
		{"f8 00", []byte{0xf8, 0x00}},
		{"f8 1f", []byte{0xf8, 0x1f}},
		{"f8 truncated", []byte{0xf8}},
		{"fc", []byte{0xfc}},
		{"fe", []byte{0xfe}},
	}
	for _, tc := range malformed {
		t.Run(tc.name+" still rejected", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{CoerceUnknownSimpleToNull: true}, tc.payload)
			if err == nil {
				t.Fatalf("expected an error")
			}
		})
	}

	t.Run("inside an array", func(t *testing.T) {
		var v []interface{}
		err := Unmarshal(DecodeOptions{CoerceUnknownSimpleToNull: true}, []byte{0x82, 0x01, 0xe0}, &v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(v) != 2 || v[1] != nil {
			t.Fatalf("expected [1 <nil>], got %v", v)
		}
	})
}
