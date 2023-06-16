package service

import (
	"context"
	"kitprc/encode"
	"kitprc/repository"
)

type Service interface {
	Post(ctx context.Context, key, val string) (res encode.Response, err error)
}

type service struct {
	repository repository.Repository
}

func (s *service) Post(ctx context.Context, key, val string) (res encode.Response, err error) {
	return s.repository.Post(key, val)
}

func New(repository repository.Repository) Service {
	return &service{repository: repository}
}
