package gzipx

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func Gunzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip reader: %w", err)
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
