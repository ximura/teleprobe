package sink

import (
	"bytes"
	"log"
	"os"
	"sync"
)

// Buffer collects telemetry messages in memory and writes them to a file
// when the buffer exceeds a maximum byte size or is flushed manually.
// It is safe for concurrent use.
type Buffer struct {
	mu       sync.Mutex
	buf      *bytes.Buffer
	maxBytes int
	writer   *os.File
}

func NewBuffer(logFilePath string, maxBytes int) (*Buffer, error) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Buffer{
		buf:      bytes.NewBuffer(make([]byte, 0, maxBytes)),
		writer:   file,
		maxBytes: maxBytes,
	}, nil
}

func (b *Buffer) Append(line string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.buf.Len()+len(line) > b.maxBytes {
		// flush before write
		log.Println("buffer limit reached")
		if err := b.flushLocked(); err != nil {
			return err
		}
	}

	_, err := b.buf.WriteString(line)
	return err
}

func (b *Buffer) Flush() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.flushLocked()
}

func (b *Buffer) flushLocked() error {
	if b.buf.Len() == 0 {
		return nil
	}

	if _, err := b.writer.Write(b.buf.Bytes()); err != nil {
		return err
	}
	if err := b.writer.Sync(); err != nil {
		return err
	}

	log.Println("buffer flushed")
	b.buf.Reset()
	return nil
}
