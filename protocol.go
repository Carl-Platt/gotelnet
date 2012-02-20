package gotelnet

import (
	"bufio"
	"io"
)

type telnetProtocol struct {
	in *bufio.Reader
	out io.Writer
	state readerState
}

func makeTelnetProtocol(in io.Reader, out io.Writer) *telnetProtocol {
	bin := bufio.NewReader(in)
	return &telnetProtocol{bin, out, readAscii}
}

type readerState func(*telnetProtocol, byte) (readerState, bool)

func readAscii(_ *telnetProtocol, c byte) (readerState, bool) {
	if c == InterpretAsCommand {
		return readCommand, false
	}
	return readAscii, true
}

func readCommand(_ *telnetProtocol, c byte) (readerState, bool) {
	switch c {
	case InterpretAsCommand:
		return readAscii, true
	case Do, Dont:
		return wontOption, false
	case Will, Wont:
		return dontOption, false
	}
	return readAscii, false
}

func wontOption(p *telnetProtocol, c byte) (readerState, bool) {
	p.out.Write([]byte{InterpretAsCommand, Wont, c})
	return readAscii, false
}

func dontOption(p *telnetProtocol, c byte) (readerState, bool) {
	p.out.Write([]byte{InterpretAsCommand, Dont, c})
	return readAscii, false
}

func (p *telnetProtocol) Read(b []byte) (n int, err error) {
	for n < len(b) {
		c, er := p.in.ReadByte()
 		if er != nil {
			err = er
			break
		} else {
			var ok bool
			if p.state, ok = p.state(p, c); ok {
				b[n] = c
				n++
			}
		}
		if p.in.Buffered() == 0 { break }
	}
	return
}
