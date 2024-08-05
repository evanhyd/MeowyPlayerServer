package logger

import (
	"bytes"
	"fmt"
	"io"
)

type circularBuffer struct {
	buffer []byte
	pos    int
	len    int
}

func makeCircularBuffer(size int) circularBuffer {
	return circularBuffer{buffer: make([]byte, size)}
}

func (b *circularBuffer) Write(msg []byte) (n int, err error) {
	if len(msg) > len(b.buffer) {
		return 0, fmt.Errorf("buffer size is too small, want %v, have %v", len(msg), len(b.buffer))
	}
	for _, c := range msg {
		b.buffer[b.pos] = c
		b.pos++
		b.len = max(b.len, b.pos)
		if b.pos >= len(b.buffer) {
			b.pos = 0
		}
	}
	return len(msg), nil
}

func (b *circularBuffer) WriteTo(w io.Writer) (int64, error) {
	n1, err := io.Copy(w, bytes.NewReader(b.buffer[b.pos:b.len]))
	if err != nil {
		return n1, err
	}
	n2, err := io.Copy(w, bytes.NewReader(b.buffer[0:b.pos]))
	return n1 + n2, err
}
