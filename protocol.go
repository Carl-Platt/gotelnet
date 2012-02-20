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

type readerState func(byte, []byte, int) (int, readerState)

func readAscii(c byte, b []byte, n int) (int, readerState) {
	if c == InterpretAsCommand {
		return n, readCommand
	}
	b[n] = c
	return n+1, readAscii
}

func readCommand(c byte, b []byte, n int) (int, readerState) {
	if c == InterpretAsCommand {
		b[n] = c
		return n+1, readAscii
	}
	return n, readAscii
}

func (p *telnetProtocol) Read(b []byte) (n int, err error) {
	for n < len(b) {
		c, er := p.in.ReadByte()
		if er == io.EOF {
			break
 		} else if er != nil {
			err = er
			break
		} else {
			n, p.state = p.state(c, b, n)
		}
	}
	return
}
