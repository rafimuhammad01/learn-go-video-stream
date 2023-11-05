package main

import (
	"context"
	"os"
)

const Prefix = "stream"

type FileReader struct{}

func (fr FileReader) GetMedia(ctx context.Context, path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func NewFileReader() FileReader {
	return FileReader{}
}
