package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/grafov/m3u8"
	"io"
)

type Reader interface {
	GetMedia(ctx context.Context, path string) ([]byte, error)
}

type Repository interface {
	GetMediaPath(ctx context.Context, id string) (string, error)
}

type MediaService struct {
	mediaReader     Reader
	mediaRepository Repository
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

func NewService(mReader Reader, mRepo Repository) MediaService {
	return MediaService{
		mediaReader:     mReader,
		mediaRepository: mRepo,
	}
}
