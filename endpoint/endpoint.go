package endpoint

import (
	"context"
	"kitprc/service"

	"github.com/go-kit/kit/endpoint"
)

type PostRequest struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type Endpoints struct {
	PostEndpoint endpoint.Endpoint
}

func NewEndpoint(s service.Service, mdw endpoint.Middleware) Endpoints {
	eps := Endpoints{
		PostEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(PostRequest)
			res, err := s.Post(ctx, req.Key, req.Val)
			return res, err
		},
	}
	eps.PostEndpoint = mdw(eps.PostEndpoint)
	return eps
}
