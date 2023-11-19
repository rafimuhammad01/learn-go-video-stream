package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

const Prefix = "stream"

type File struct{}

func (fr File) GetMedia(ctx context.Context, path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (fr File) UploadMedia(ctx context.Context, file io.Reader, ext string) error {
	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./downloaded/%d%s", time.Now().UnixNano(), ext))
	if err != nil {
		return err
	}
	defer dst.Close()

	// copy file to dst
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}

func NewFile() File {
	return File{}
}
