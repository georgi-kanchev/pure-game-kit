package storage

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Compress(data []byte) []byte {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write(data)
	if err != nil {
		return data
	}
	if err := gw.Close(); err != nil {
		return data
	}
	return buf.Bytes()
}

func Decompress(data []byte) []byte {
	buf := bytes.NewReader(data)
	gr, err := gzip.NewReader(buf)
	if err != nil {
		return data
	}
	defer gr.Close()

	var result, err2 = io.ReadAll(gr)
	if err2 != nil {
		return data
	}

	return result
}
