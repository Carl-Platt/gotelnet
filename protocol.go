package gotelnet

import (
	"bufio"
	"io"
)

type telnetProtocol struct {
	in *bufio.Reader
	out io.Writer

	iac bool
}

func makeTelnetProtocol(in io.Reader, out io.Writer) *telnetProtocol {
	bin := bufio.NewReader(in)
	return &telnetProtocol{bin, out, false}
}

func (p *telnetProtocol) Read(b []byte) (n int, err error) {
	var c byte
	for max := len(b); n < max; {
		if p.iac {
			p.iac = false
			switch c, err = p.in.ReadByte(); c {
			case InterpretAsCommand:
				b[n] = c
				n++
			default:
				continue
			}
		} else{
			switch c, err = p.in.ReadByte(); c {
			case InterpretAsCommand:
				p.iac = true
			default:
				b[n] = c
				n++
			}
		}
		if p.in.Buffered() == 0 { break }
	}
	return n, nil
}
