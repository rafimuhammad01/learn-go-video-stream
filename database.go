package main

import (
	"context"
	"fmt"
	"sync"
)

var memDB = map[string]string{
	"1": "assets/rick-astley-never-gonna-give-you-up-official-music-video.m3u8",
}

type DB struct {
	sync.Mutex
	storage map[string]string
}

func (db *DB) GetMediaPath(ctx context.Context, id string) (string, error) {
	db.Mutex.Lock()
	path, ok := db.storage[id]
	db.Mutex.Unlock()
	if !ok {
		return "", fmt.Errorf("file not found with id=%s", id)
	}

	return path, nil
}

func NewDB() *DB {
	return &DB{
		Mutex:   sync.Mutex{},
		storage: memDB,
	}
}
