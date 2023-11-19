package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/grafov/m3u8"
)

type Reader interface {
	GetMedia(ctx context.Context, path string) ([]byte, error)
}

type Uploader interface {
	UploadMedia(ctx context.Context, file io.Reader, fileExt string) error
}

type Repository interface {
	GetMediaPath(ctx context.Context, id string) (string, error)
}

type MediaService struct {
	mediaReader     Reader
	mediaRepository Repository
	mediaUploader   Uploader
}

func (ms MediaService) GetMedia(ctx context.Context, id string) ([]byte, error) {
	mediaPath, err := ms.mediaRepository.GetMediaPath(ctx, id)
	if err != nil {
		return nil, err
	}

	buff, err := ms.mediaReader.GetMedia(ctx, mediaPath)
	if err != nil {
		return nil, err
	}

	p, err := ms.appendPath(ctx, bytes.NewReader(buff))
	if err != nil {
		return nil, err
	}

	return p.Encode().Bytes(), nil
}

func (ms MediaService) GetMediaStream(ctx context.Context, path string) ([]byte, error) {
	buff, err := ms.mediaReader.GetMedia(ctx, fmt.Sprintf("assets/%v", path))
	if err != nil {
		return nil, err
	}

	return buff, nil
}

func (ms MediaService) UploadMedia(ctx context.Context, file io.Reader, ext string) error {
	// Read the first 512 bytes
	buffer := make([]byte, 512)
	_, err := io.ReadFull(file, buffer)
	if err != nil && err != io.EOF {
		return err
	}

	// Detect the file type based on its magic number or signature
	fileType := http.DetectContentType(buffer)
	if !ms.isVideoFile(fileType) {
		return errors.New("file is not video file")
	}

	return ms.mediaUploader.UploadMedia(ctx, file, ext)
}

func (ms MediaService) appendPath(ctx context.Context, r io.Reader) (m3u8.Playlist, error) {
	p, listType, err := m3u8.DecodeFrom(r, true)
	if err != nil {
		return nil, err
	}

	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		for i := 0; i < len(mediapl.Segments); i++ {
			if mediapl.Segments[i] == nil {
				break
			}

			mediapl.Segments[i].URI = fmt.Sprintf("%s/%s", Prefix, mediapl.Segments[i].URI)
		}
		p = mediapl
	}

	return p, nil
}

func (ms MediaService) isVideoFile(fileType string) bool {
	return fileType == "video/mp4" || fileType == "video/quicktime"
}

func NewService(mReader Reader, mUploader Uploader, mRepo Repository) MediaService {
	return MediaService{
		mediaUploader:   mUploader,
		mediaReader:     mReader,
		mediaRepository: mRepo,
	}
}
