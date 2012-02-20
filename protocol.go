package gotelnet

import (
	"io"
)

type telnetProtocol struct {
	in io.Reader
	out io.Writer
}

func (p *telnetProtocol) Read(b []byte) (n int, err error) {
	buf := make([]byte, len(b))
	n, err = p.in.Read(buf)
	buf = buf[0:n]
	for i := 0; len(buf) > 0; {
		switch buf[0] {
		case InterpretAsCommand:
			n--
			switch buf[1] {
			case InterpretAsCommand:
				buf = buf[1:]
			default :
				n--
				buf = buf[2:]
				continue
			}
			fallthrough
		default:
			b[i] = buf[0]
			buf = buf[1:]
			i++
		}
	}
	return
}
