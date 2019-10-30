package odf

import (
	"bytes"
	"testing"
)

// emptyZip:
//	> _
// 	zip -r _.zip _
//	zip -d _.zip _
//	od -x _.zip
var emptyZip = [22]byte{0x50, 0x4b, 0x05, 0x06}

func TestNewReader(t *testing.T) {

	// Test whether NewReader properly fails with an error,
	// when provided with a zip file that is not a well-formed ODF file
	t.Run("NilCloser", func(t *testing.T) {
		z := emptyZip[:]
		r, err := NewReader(bytes.NewReader(z), int64(len(z)))
		if r != nil || err == nil {
			t.Fatal("expected an error")
		}
	})
}
