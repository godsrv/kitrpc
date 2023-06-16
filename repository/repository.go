package repository

import (
	"kitprc/encode"
	"sync"
	"time"
)

type Store struct {
	Key       string
	Val       string
	CreatedAt time.Time
}

type Repository interface {
	Post(key, val string) (encode.Response, error)
}

type store struct {
	mtx sync.RWMutex
}

func (s *store) Post(key, val string) (encode.Response, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return encode.Response{
		Code:    0,
		Data:    map[string]interface{}{"key": key, "val": val},
		Message: "success",
	}, nil
}

func New() Repository {
	return &store{}
}
