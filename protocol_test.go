package gotelnet

import (
	"bytes"
	"testing"
)

func processBytes(t *testing.T, b []byte) (r, w []byte) {
	r = make([]byte, len(b)) 	// At most we'll read all the bytes
	in := bytes.NewBuffer(b)
	out := &bytes.Buffer{}
	protocol := &telnetProtocol{in, out}
	if n, err := protocol.Read(r); err != nil {
		t.Fatalf("Read error %q", err)
	} else {
		r = r[0:n] 				// Truncate to the length actually read
		t.Logf("Read %d bytes %q", n, r)
	}
	w = out.Bytes()
	t.Logf("Wrote %d bytes %q", len(w), w)
	return
}

func assertEqual(t *testing.T, a, b []byte) {
	if !bytes.Equal(a, b) {
		t.Fatalf("Expected %q to be %q", a, b)
	}
}

func TestAsciiText(t *testing.T) {
	r, w := processBytes(t, []byte("hello"))
	assertEqual(t, r, []byte("hello"))
	assertEqual(t, w, []byte{})
}

func TestStripTelnetCommands(t *testing.T) {
	r, w := processBytes(t, []byte{'h', InterpretAsCommand, NoOperation, 'i'})
	assertEqual(t, r, []byte("hi"))
	assertEqual(t, w, []byte{})
}

func TestEscapedIAC(t *testing.T) {
	r, w := processBytes(t, []byte{'h', InterpretAsCommand, InterpretAsCommand, 'i'})
	assertEqual(t, r, []byte("h\xffi"))
	assertEqual(t, w, []byte{})
}
